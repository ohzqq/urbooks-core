package book

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"log"
	"os"
	"regexp"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gosimple/slug"
	"github.com/ohzqq/avtools/avtools"
	"gopkg.in/ini.v1"
)

func ListFormats() []string {
	var fmts []string
	for _, f := range MetaFmt {
		fmts = append(fmts, f.name)
	}
	return fmts
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
	"toMarkdown":   toMarkdown,
	"stringToHTML": stringToHTML,
}

func stringToHTML(s string) template.HTML {
	return template.HTML(html.UnescapeString(s))
}

var MetaFmt = []metaFmt{
	metaFmt{
		name: "ffmeta",
		ext:  ".ini",
		tmpl: template.Must(template.New("ffmeta").Funcs(funcMap).Parse(ffmetaTmpl)),
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

func (b *Book) ConvertTo(f string) *Book {
	for _, fmt := range MetaFmt {
		if fmt.name == f {
			b.fmt = fmt
		}
	}
	return b
}

func getFmt(n string) (metaFmt, error) {
	return metaFmt{}, fmt.Errorf("Not a format")
}

func (b *Book) Print() {
	fmt.Println(string(b.fmt.Render(b)))
}

func (b *Book) Tmp() *os.File {
	file, err := os.CreateTemp("", b.fmt.ext)
	if err != nil {
		log.Fatal(err)
	}
	m := b.fmt.Render(b)
	fmt.Println(string(m))
	_, err = file.Write(m)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func (b *Book) Write() {
	file, err := os.Create(slug.Make(b.GetField("title").String()) + b.fmt.ext)
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
		return ToOpf(b).Marshal()
	default:
		err := m.tmpl.Execute(&buf, b.StringMap())
		if err != nil {
			log.Fatal(err)
		}
	}
	return buf.Bytes()
}

func (b *Book) StringMap() map[string]string {
	m := make(map[string]string)
	for _, field := range b.EachField() {
		key := field.Label()
		//key := strings.TrimPrefix(field.Label(), "#")

		if key != "customColumns" && !field.IsNull() {
			m[key] = field.String()
		}

		if key == "titleAndSeries" && field.IsNull() {
			m["titleAndSeries"] = b.GetTitleAndSeries()
		}
	}
	return m
}

func (b *Book) ToIni() {
	ini.PrettyFormat = false
	file := ini.Empty(ini.LoadOptions{
		AllowNonUniqueSections: true,
	})
	sec, err := file.GetSection("")
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range b.StringMap() {
		_, err := sec.NewKey(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}
	file.WriteTo(os.Stdout)
}

func MediaMetaToBook(lib string, m *avtools.Media) *Book {
	book := NewBook()
	titleRegex := regexp.MustCompile(`(?P<title>.*) \[(?P<series>.*), Book (?P<position>.*)\]$`)
	titleAndSeries := titleRegex.FindStringSubmatch(m.GetTag("title"))

	book.GetField("title").SetData(titleAndSeries[titleRegex.SubexpIndex("title")])
	book.GetField("series").
		SetData(titleAndSeries[titleRegex.SubexpIndex("series")])
	book.GetField("series").
		SetData(titleAndSeries[titleRegex.SubexpIndex("position")])
	book.GetField("authors").Collection().Split(m.GetTag("artist"), true)
	//book.GetField("#narrators").Collection().Split(m.GetTag("composer"), true)
	book.GetField("description").SetData(m.GetTag("comment"))
	book.GetField("tags").Collection().Split(m.GetTag("genre"), false)
	return book
}

func toMarkdown(str string) string {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(str)
	if err != nil {
		log.Fatal(err)
	}
	return markdown
}

const ffmetaTmpl = `;FFMETADATA
title={{stringToHTML .titleAndSeries}}
album={{stringToHTML .titleAndSeries}}
artist={{stringToHTML .authors}}
composer={{stringToHTML .narrators}}
genre={{stringToHTML .tags}}
comment={{stringToHTML .description}}
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
