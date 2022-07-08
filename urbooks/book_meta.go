package urbooks

import (
	"bytes"
	"fmt"
	"io"
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
	return bm[k]
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
	for key, val := range bm {
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
				item := book.NewItem(key).SetValue(val.String())
				if key == "series" {
					if pos := bm.Get("position").String(); pos != "" {
						item.Set("position", pos)
					}
				}
			}
		default:
			book.NewColumn(key).SetValue(val.String())
		}
	}
	return book
}

type metaFmt struct {
	tmpl   *template.Template
	ext    string
	name   string
	save   bool
	buffer bytes.Buffer
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
			b.tmpl = fmt
		}
	}
	return b
}

func (b *Book) Print() {
	b.tmpl.Render(os.Stdout, b)
}

func (b *Book) Write() {
	file, err := os.Create(slug.Make(b.Get("title").String()) + b.tmpl.ext)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	b.tmpl.Render(file, b)
}

func (m metaFmt) Render(wr io.Writer, b *Book) {
	err := m.tmpl.Execute(wr, b.StringMap())
	if err != nil {
		log.Fatal(err)
	}
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
