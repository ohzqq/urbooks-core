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

	"github.com/BurntSushi/toml"
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
	b.AddField(NewCollection("#narrators")).SetIsNames().SetIsCustom().SetMeta(m.GetTag("composer"))
	b.GetField("description").SetMeta(m.GetTag("comment"))
	b.GetField("tags").SetMeta(m.GetTag("genre"))
	return b
}

func (b *Book) ConvertTo(f string) Fmt {
	for _, fmt := range MetaFmt {
		if fmt.name == f {
			b.tmpl = fmt.tmpl
			fmt.book = b
			return fmt
		}
	}
	return Fmt{}
}

func (b *Book) StringMap(hash bool) map[string]string {
	m := make(map[string]string)
	for key, field := range b.EachField() {
		//if key == "titleAndSeries" {
		//  m["titleAndSeries"] = b.GetTitleAndSeries()
		//}

		if field.IsEditable && !field.IsNull() {
			if !hash {
				key = strings.TrimPrefix(key, "#")
			}
			m[key] = field.String()
		}
	}
	return m
}

func (b *Book) DataMap(hash bool) map[string]interface{} {
	m := make(map[string]interface{})
	for key, field := range b.EachField() {
		if field.IsEditable && !field.IsNull() {
			if !hash {
				key = strings.TrimPrefix(key, "#")
			}
			switch field.IsMultiple {
			case true:
				m[key] = field.Collection().StringSlice()
			default:
				m[key] = field.String()
			}
		}
	}
	return m
}

type Fmt struct {
	tmpl   *template.Template
	book   *Book
	ext    string
	name   string
	hash   bool
	data   []byte
	render func(b *Book, hash bool) *bytes.Buffer
}

func (f Fmt) String() string {
	return f.render(f.book, f.hash).String()
}

func (f Fmt) Bytes() []byte {
	return f.render(f.book, f.hash).Bytes()
}

func (f Fmt) Print() {
	fmt.Println(f.String())
}

func (f Fmt) Write() {
	file, err := os.Create(slug.Make(f.book.GetMeta("title")) + f.ext)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(f.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func (f Fmt) Tmp() *os.File {
	file, err := os.CreateTemp("", f.ext)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(f.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func renderTmpl(b *Book, hash bool) *bytes.Buffer {
	var buf bytes.Buffer
	err := b.tmpl.Execute(&buf, b)
	if err != nil {
		log.Fatal(err)
	}
	return &buf
}

func ToToml(b *Book, hash bool) *bytes.Buffer {
	var buf bytes.Buffer
	err := toml.NewEncoder(&buf).Encode(b.StringMap(hash))
	if err != nil {
		log.Fatal(err)
	}
	return &buf
}

var iniOpts = ini.LoadOptions{
	IgnoreInlineComment:    true,
	AllowNonUniqueSections: true,
}

func ToIni(b *Book, hash bool) *bytes.Buffer {
	book := b.StringMap(hash)
	//ini.PrettyFormat = false
	file := ini.Empty(iniOpts)
	sec, err := file.GetSection("")
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range book {
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

	return &buf
}

func stringToHTML(s string) template.HTML {
	return template.HTML(html.UnescapeString(s))
}

func toMarkdown(str string) string {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(str)
	if err != nil {
		log.Fatal(err)
	}
	return markdown
}

var (
	funcMap = template.FuncMap{
		"toMarkdown":   toMarkdown,
		"stringToHTML": stringToHTML,
		"ToIni":        ToIni,
	}

	MetaFmt = []Fmt{
		Fmt{
			name:   "ffmeta",
			ext:    ".ini",
			tmpl:   template.Must(template.New("ffmeta").Funcs(funcMap).Parse(ffmetaTmpl)),
			render: renderTmpl,
		},
		Fmt{
			name:   "markdown",
			ext:    ".md",
			hash:   true,
			tmpl:   template.Must(template.New("md").Funcs(funcMap).Parse(mdTmpl)),
			render: renderTmpl,
		},
		Fmt{
			name:   "md",
			ext:    ".md",
			hash:   true,
			tmpl:   template.Must(template.New("md").Funcs(funcMap).Parse(mdTmpl)),
			render: renderTmpl,
		},
		Fmt{
			name:   "plain",
			ext:    ".txt",
			tmpl:   template.Must(template.New("plain").Funcs(funcMap).Parse(plainTmpl)),
			render: renderTmpl,
		},
		Fmt{
			name:   "opf",
			ext:    ".opf",
			render: func(b *Book, hash bool) *bytes.Buffer { return buildOPF(b).Marshal() },
		},
		Fmt{
			name:   "ini",
			ext:    ".ini",
			render: ToIni,
		},
		Fmt{
			name:   "toml",
			ext:    ".toml",
			hash:   true,
			render: ToToml,
		},
		Fmt{
			name:   "rss",
			ext:    ".xml",
			render: func(b *Book, hash bool) *bytes.Buffer { return BookToRssChannel(b).Marshal() },
		},
	}
)

const ffmetaTmpl = `;FFMETADATA
title={{with .GetTitleAndSeries}}{{stringToHTML .}}{{end}}
album={{with .GetTitleAndSeries}}{{stringToHTML .}}{{end}}
artist={{with .GetMeta "authors"}}{{stringToHTML .}}{{end}}
composer={{with .GetMeta "#narrators"}}{{stringToHTML .}}{{end}}
genre={{with .GetMeta "tags"}}{{stringToHTML .}}{{end}}
comment={{with .GetMeta "description"}}{{stringToHTML .}}{{end}}`

const mdTmpl = `
{{- with .GetMeta "title"}}# {{stringToHTML .}}{{end}}
**Series:** {{with .GetSeriesString}}{{stringToHTML .}}{{end}}
**Authors:** {{with .GetMeta "authors"}}{{stringToHTML .}}{{end}}
**Narrators:** {{with .GetMeta "narrators"}}{{stringToHTML .}}{{end}}
**Tags:** {{with .GetMeta "tags"}}{{stringToHTML .}}{{end}}
**Rating:** {{with .GetMeta "rating"}}{{stringToHTML .}}{{end}}
**Description:** {{with .GetMeta "description"}}{{toMarkdown .}}{{end}}`

const plainTmpl = `
{{- with .GetMeta "title"}}{{stringToHTML .}}{{end}}
Series: {{with .GetSeriesString}}{{stringToHTML .}}{{end}}
Authors: {{with .GetMeta "authors"}}{{stringToHTML .}}{{end}}
Narrators: {{with .GetMeta "narrators"}}{{stringToHTML .}}{{end}}
Tags: {{with .GetMeta "tags"}}{{stringToHTML .}}{{end}}
Rating: {{with .GetMeta "rating"}}{{stringToHTML .}}{{end}}
Description: {{with .GetMeta "description"}}{{toMarkdown .}}{{end}}`
