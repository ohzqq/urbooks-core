package urbooks

import (
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
	if meta == "titleAndSeries" {
		field = bm.Get("series")
	}

	switch meta {
	case "formats":
		return bm.GetCategory(meta).Join("extension")
	case "position":
		if series := bm.GetItem("series"); series.IsNull() {
			return series.Get("position")
		}
	case "titleAndSeries":
		title := bm.Get("title").Value()
		if series := bm.Get("series"); !field.IsNull() {
			return title + " [" + series.String() + "]"
		}
		return title
	}

	if field.FieldMeta().Type() == "category" && !field.IsNull() {
		f := field.(*Item)
		if meta == "series" {
			return f.Value() + ", Book " + f.Get("position")
		}
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

type MetaString string

func NewMetaString() *MetaString {
	ms := MetaString("")
	return &ms
}

func (ms *MetaString) SetValue(v string) *MetaString {
	s := MetaString(v)
	ms = &s
	return ms
}

func (ms MetaString) URL() string                 { return "" }
func (ms MetaString) IsNull() bool                { return ms == "" }
func (ms MetaString) Value() string               { return string(ms) }
func (ms MetaString) String() string              { return string(ms) }
func (ms MetaString) FieldMeta() *calibredb.Field { return &calibredb.Field{} }

func (b *Book) ToFFmeta() {
	meta, err := os.Create(slug.Make(b.Get("title").String()) + ".ini")
	if err != nil {
		log.Fatal(err)
	}
	defer meta.Close()

	err = MetaFmt.FFmeta.Execute(meta, b)
	if err != nil {
		log.Fatal(err)
	}
}

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

type metadataFormats struct {
	FFmeta *template.Template
	MD     *template.Template
	Plain  *template.Template
	Cue    *template.Template
}

var funcMap = template.FuncMap{
	"toMarkdown": toMarkdown,
}

func toMarkdown(str string) string {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(str)
	if err != nil {
		log.Fatal(err)
	}
	return markdown
}

var MetaFmt = metadataFormats{
	FFmeta: template.Must(template.New("ffmeta").Parse(ffmetaTmpl)),
	MD:     template.Must(template.New("html").Funcs(funcMap).Parse(mdTmpl)),
	Plain:  template.Must(template.New("plain").Funcs(funcMap).Parse(plainTmpl)),
}

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
