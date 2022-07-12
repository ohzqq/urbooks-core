package urbooks

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"text/template"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gosimple/slug"
	"github.com/ohzqq/urbooks-core/calibredb"
)

type Meta interface {
	Value() string
	String() string
	URL() string
	FieldMeta() *calibredb.Field
	IsNull() bool
}

type BookMeta map[string]Meta

func NewBookMeta(m map[string]string) BookMeta {
	meta := make(BookMeta)
	for key, val := range m {
		meta[key] = MetaString(val)
	}
	return meta
}

func (bm BookMeta) Get(k string) Meta {
	if _, ok := bm[k]; ok {
		return bm[k]
	}
	return NewColumn()
}

func (bm BookMeta) FieldMeta(f string) *calibredb.Field {
	return bm.Get(f).FieldMeta()
}

func (bm BookMeta) GetItem(f string) *Item {
	if field := bm.Get(f); field.FieldMeta().Type() == "item" {
		return field.(*Item)
	}
	return &Item{}
}

func (bm BookMeta) GetCategory(f string) *Category {
	if field := bm.Get(f); field.FieldMeta().Type() == "category" {
		return field.(*Category)
	}
	return &Category{}
}

func (bm BookMeta) GetColumn(f string) *Column {
	if field := bm.Get(f); field.FieldMeta().Type() == "column" {
		return field.(*Column)
	}
	return &Column{}
}

func (bm BookMeta) series() string {
	position := "1.0"
	series := bm.GetItem("series")
	if series.Value() != "" {
		if pos := bm.Get("position").String(); pos != "" {
			position = pos
		}
		if pos := series.Get("position"); pos != "" {
			position = pos
		}
	}
	return series.Value() + `, Book ` + position
}

func (bm BookMeta) String(meta string) string {
	field := bm.Get(meta)
	if field.IsNull() {
		return ""
	}

	switch meta {
	case "formats":
		return bm.GetCategory(meta).Join("extension")
	case "series":
		return bm.series()
	}

	if field.FieldMeta().Type() == "category" && !field.IsNull() {
		f := field.(*Item)
		return f.Value()
	}

	return field.Value()
}

func (bm BookMeta) StringMap() map[string]string {
	m := make(map[string]string)
	lib := Lib(bm.Get("library").String())
	for k, val := range bm {
		field := lib.DB.GetField(k)
		var key string
		switch {
		case field.IsCustom:
			key = "#" + k
		default:
			key = k
		}
		m[key] = val.String()
		if key == "series" {
			if pos := bm.Get("series").(*Item).Get("position"); pos != "" {
				m["position"] = pos
			}
		}
	}
	return m
}

func (bm BookMeta) StringMapToBook() *Book {
	lib := DefaultLib()
	if l := bm["library"].Value(); l == "" {
		lib = Lib(l)
	}
	book := NewBook(lib.Name)
	for key, val := range bm {
		field := lib.DB.GetField(key)
		switch {
		case field.IsCategory:
			switch field.IsMultiple {
			case true:
				cat := book.NewCategory(key)
				switch {
				case field.IsNames:
					cat.Split(val.String(), true)
				default:
					cat.Split(val.String(), false)
				}
			case false:
				book.NewItem(key).SetValue(val.String())
			}
		default:
			book.NewColumn(key).SetValue(val.String())
		}
	}
	return book
}

func (b *Book) opfFields() []string {
	fields := []string{"authors", "tags", "languages", "identifiers", "title", "published", "description", "series"}
	for _, f := range Lib(b.Get("library").Value()).DB.CustomColumns() {
		fields = append(fields, f)
	}
	return fields
}

func (b *Book) ToOpf() []byte {
	opf := NewOpfMetadata()
	for _, f := range b.opfFields() {
		meta := b.Get(f)
		field := meta.FieldMeta()
		switch {
		case field.IsCustom:
			cField := field
			var val interface{}
			if cField.IsMultiple {
				val = meta.GetCategory().ItemStringSlice()
			} else {
				val = meta.String()
			}
			cField.Value = val
			opf.AddMeta("#"+f, string(cField.ToJson()))
		case field.IsCat():
			if cat := meta.GetCategory(); !cat.IsNull() {
				for _, item := range cat.Items() {
					switch f {
					case "authors":
						opf.AddAuthor(item.String())
					case "tags":
						opf.AddSubject(item.String())
					case "languages":
						opf.AddLanguage(item.String())
					case "identifiers":
						fmt.Printf("%+v\n", item)
						opf.AddIdentifier(item.Get("value"), item.Get("type"))
					}
				}
			}
		case field.IsItem():
			if i := meta.GetItem(); !i.IsNull() {
				switch f {
				case "series":
					opf.SetSeries(i.Get("value"))
					opf.SetSeriesIndex(i.Get("position"))
				}
			}
		case field.IsCol():
			if col := meta.GetColumn(); !col.IsNull() {
				switch f {
				case "title":
					opf.SetTitle(col.String())
				case "published":
					opf.SetDate(col.String())
				case "description":
					opf.SetDescription(col.String())
				}
			}
		}
	}

	pkg := bytes.NewBufferString(xml.Header)
	enc := xml.NewEncoder(pkg)
	enc.Indent("", "  ")
	err := enc.Encode(opf.BuildPackage())
	if err != nil {
		log.Fatal(err)
	}
	return pkg.Bytes()
}

type metaFmt struct {
	tmpl   *template.Template
	ext    string
	name   string
	save   bool
	data   []byte
	buffer *bytes.Buffer
}

var funcMap = template.FuncMap{
	"toMarkdown": toMarkdown,
}

var MetaFmt = []metaFmt{
	metaFmt{
		name: "ffmeta",
		ext:  ".ini",
		tmpl: template.Must(template.New("ffmeta").Parse(ffmetaTmpl)),
	},
	metaFmt{
		name: "markdown",
		ext:  ".md",
		tmpl: template.Must(template.New("md").Funcs(funcMap).Parse(mdTmpl)),
	},
	metaFmt{
		name: "plain",
		ext:  ".txt",
		tmpl: template.Must(template.New("plain").Funcs(funcMap).Parse(plainTmpl)),
	},
	metaFmt{
		name: "opf",
		ext:  ".opf",
	},
}

func toMarkdown(str string) string {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(str)
	if err != nil {
		log.Fatal(err)
	}
	return markdown
}

func ListFormats() []string {
	var fmts []string
	for _, f := range MetaFmt {
		fmts = append(fmts, f.name)
	}
	return fmts
}

func getFmt(n string) (metaFmt, error) {
	return metaFmt{}, fmt.Errorf("Not a format")
}

func (b *Book) ConvertTo(f string) *Book {
	for _, fmt := range MetaFmt {
		if fmt.name == f {
			b.fmt = fmt
		}
	}
	return b
}

func (b *Book) Print() {
	fmt.Println(string(b.fmt.Render(b)))
}

func (b *Book) Write() {
	file, err := os.Create(slug.Make(b.Get("title").String()) + b.fmt.ext)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(b.fmt.Render(b))
	if err != nil {
		log.Fatal(err)
	}
}

func (m metaFmt) Render(b *Book) []byte {
	var buf bytes.Buffer
	switch m.name {
	case "opf":
		return b.ToOpf()
	default:
		err := m.tmpl.Execute(&buf, b.StringMap())
		if err != nil {
			log.Fatal(err)
		}
	}
	return buf.Bytes()
}

const ffmetaTmpl = `;FFMETADATA
title={{.titleAndSeries}}
album={{.titleAndSeries}}
artist={{.authors}}
composer={{.narrators}}
genre={{.tags}}
comment={{.description}}
`

const mdTmpl = `
{{- with .title}}# {{.}}{{end}}
**Series:** {{with .series}}{{.}}{{end}}
**Authors:** {{with .authors}}{{.}}{{end}}
**Narrators:** {{with .narrators}}{{.}}{{end}}
**Tags:** {{with .tags}}{{.}}{{end}}
**Rating:** {{with .rating}}{{.}}{{end}}
**Description:** {{with .description}}{{toMarkdown .}}{{end}}`

const plainTmpl = `
{{- with .title}}{{.}}{{end}}
Series: {{with .series}}{{.}}{{end}}
Authors: {{with .authors}}{{.}}{{end}}
Narrators: {{with .narrators}}{{.}}{{end}}
Tags: {{with .tags}}{{.}}{{end}}
Rating: {{with .rating}}{{.}}{{end}}
Description: {{with .description}}{{toMarkdown .}}{{end}}`
