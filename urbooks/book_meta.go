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

func (bm BookMeta) String(meta string) string {
	field := bm.Get(meta)

	switch meta {
	case "formats":
		return bm.GetCategory(meta).Join("extension")
	case "position":
		if series := bm.GetItem("series"); series.IsNull() {
			return series.Get("position")
		}
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

//func (b *Book) ToFFmeta() {
//  meta, err := os.Create(slug.Make(b.Get("title").String()) + ".ini")
//  if err != nil {
//    log.Fatal(err)
//  }
//  defer meta.Close()

//  err = MetaFmt.FFmeta.Execute(meta, b)
//  if err != nil {
//    log.Fatal(err)
//  }
//}

//func (b *Book) ToPlain() string {
//  var buf bytes.Buffer
//  err := MetaFmt.Plain.Execute(&buf, b)
//  if err != nil {
//    log.Fatal(err)
//  }
//  return buf.String()
//}

//func (b *Book) ToMarkdown() string {
//  var buf bytes.Buffer
//  err := MetaFmt.MD.Execute(&buf, b)
//  if err != nil {
//    log.Fatal(err)
//  }
//  //fmt.Println(markdown)
//  return buf.String()
//}

const ffmetaTmpl = `;FFMETADATA
{{$title := .Get "titleAndSeries" -}}
title={{$title.String}}
album={{$title.String}}
artist=
{{- with $authors := .Get "authors" -}}
	{{- $authors.String -}}
{{- end}}
composer=
{{- with $narrators := .Get "narrators" -}}
	{{- $narrators.String -}}
{{- end}}
genre=
{{- with $tags := .Get "tags" -}}
	{{- $tags.String -}}
{{- end}}
comment=
{{- with $description := .Get "description" -}}
	{{- $description.String -}}
{{- end -}}
`

const mdTmpl = `{{if .Title}}# {{.Title}}   
{{end}}{{if .HasSeries}}**Series:** {{.SeriesString}}   
{{end}}{{if .Authors}}**Authors:** {{.Authors.Join}}   
{{end}}{{if .Narrators}}**Narrators:** {{.Narrators.Join}}   
{{end}}{{if .Tags}}**Tags:** {{.Tags.Join}}   
{{end}}{{if .Rating}}**Rating:** {{.Rating}}   
{{end}}{{if .Description}}**Description:** {{toMarkdown .Description}}{{end}}`

const plainTmpl = `{{if .Title}}{{.Title}}   
{{end}}{{if .HasSeries}}Series: {{.SeriesString}}   
{{end}}{{if .Authors}}Authors: {{.Authors.Join}}   
{{end}}{{if .Narrators}}Narrators: {{.Narrators.Join}}   
{{end}}{{if .Tags}}Tags: {{.Tags.Join}}   
{{end}}{{if .Rating}}Rating: {{.Rating}}   
{{end}}{{if .Description}}Description: {{.Description}}{{end}}`
