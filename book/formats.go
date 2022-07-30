package book

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"

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

func (b *Book) StringMap() map[string]string {
	return StringMap(b, false)
}

func DataMap(b *Book, hash bool) map[string]interface{} {
	m := make(map[string]interface{})
	for _, field := range b.EachField() {
		key := field.Label()
		if hash {
			key = strings.TrimPrefix(key, "#")
		}

		if key != "customColumns" && !field.IsNull() {
			m[key] = field.String()
		}

		if key == "titleAndSeries" && field.IsNull() {
			m["titleAndSeries"] = b.GetTitleAndSeries()
		}
	}
	return m
}

func StringMap(b *Book, hash bool) map[string]string {
	m := make(map[string]string)
	for _, field := range b.EachField() {
		key := field.Label()
		if hash {
			key = strings.TrimPrefix(key, "#")
		}

		if key != "customColumns" && !field.IsNull() {
			m[key] = field.String()
		}

		if key == "titleAndSeries" && field.IsNull() {
			m["titleAndSeries"] = b.GetTitleAndSeries()
		}
	}
	return m
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
	"ToIni":        ToIni,
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
	metaFmt{
		name: "ini",
		ext:  ".ini",
	},
	metaFmt{
		name: "toml",
		ext:  ".toml",
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
	file, err := os.Create(slug.Make(b.GetMeta("title")) + b.fmt.ext)
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
	case "ini":
		return b.ToIni()
	default:
		err := m.tmpl.Execute(&buf, StringMap(b, true))
		if err != nil {
			log.Fatal(err)
		}
	}
	return buf.Bytes()
}

//func (b *Book) toToml() []byte {
//}

func ToIni(b map[string]string) string {
	ini.PrettyFormat = false
	file := ini.Empty(ini.LoadOptions{
		AllowNonUniqueSections: true,
	})
	sec, err := file.GetSection("")
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range b {
		_, err := sec.NewKey(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	var buf bytes.Buffer
	_, err = file.WriteTo(&buf)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (b *Book) ToIni() []byte {
	return []byte(ToIni(b.StringMap()))
}

func MediaMetaToBook(lib string, m *avtools.Media) *Book {
	b := NewBook()
	titleRegex := regexp.MustCompile(`(?P<title>.*) \[(?P<series>.*), Book (?P<position>.*)\]$`)
	titleAndSeries := titleRegex.FindStringSubmatch(m.GetTag("title"))

	b.GetField("title").SetMeta(titleAndSeries[titleRegex.SubexpIndex("title")])
	b.GetField("series").
		SetMeta(titleAndSeries[titleRegex.SubexpIndex("series")])
	b.GetField("series").
		SetMeta(titleAndSeries[titleRegex.SubexpIndex("position")])
	b.GetField("authors").SetMeta(m.GetTag("artist"))
	b.AddField(NewCollection("#narrators")).SetIsNames().SetMeta(m.GetTag("composer"))
	b.GetField("description").SetMeta(m.GetTag("comment"))
	b.GetField("tags").SetMeta(m.GetTag("genre"))
	return b
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
