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
			fmt.book = b
			return fmt
			//b.fmt = fmt
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
	buffer bytes.Buffer
}

func (f Fmt) Render() Fmt {
	switch f.name {
	case "opf":
		f.data = ToOpf(f.book).Marshal()
	case "rss":
		f.data = f.ToRss()
	case "ini":
		f.data = f.ToIni()
	case "toml":
		//f.ToToml()
		err := toml.NewEncoder(&f.buffer).Encode(f.book.StringMap(f.hash))
		if err != nil {
			log.Fatal(err)
		}
		f.data = f.buffer.Bytes()
	default:
		err := f.tmpl.Execute(&f.buffer, f.book)
		if err != nil {
			log.Fatal(err)
		}
		f.data = f.buffer.Bytes()
	}
	return f
}

func (f Fmt) String() string {
	return string(f.data)
}

func (f Fmt) Bytes() []byte {
	return f.data
}

func getFmt(n string) (Fmt, error) {
	return Fmt{}, fmt.Errorf("Not a format")
}

func (f Fmt) Print() {
	fmt.Println(f.Render().String())
}

func (f Fmt) Tmp() *os.File {
	file, err := os.CreateTemp("", f.ext)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(f.Render().Bytes())
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func (f Fmt) Write() {
	file, err := os.Create(slug.Make(f.book.GetMeta("title")) + f.ext)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(f.Render().Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func (f Fmt) ToToml() []byte {
	err := toml.NewEncoder(&f.buffer).Encode(f.book.StringMap(f.hash))
	if err != nil {
		log.Fatal(err)
	}
	return f.buffer.Bytes()
}

func (f Fmt) ToRss() []byte {
	return BookToRssChannel(f.book).Marshal()
}

var iniOpts = ini.LoadOptions{
	IgnoreInlineComment:    true,
	AllowNonUniqueSections: true,
}

func ToIni(b map[string]string) []byte {
	//ini.PrettyFormat = false
	file := ini.Empty(iniOpts)
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

	return buf.Bytes()
}

func (f Fmt) ToIni() []byte {
	return ToIni(f.book.StringMap(f.hash))
}

var funcMap = template.FuncMap{
	"toMarkdown":   toMarkdown,
	"stringToHTML": stringToHTML,
	"ToIni":        ToIni,
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

var MetaFmt = []Fmt{
	Fmt{
		name: "ffmeta",
		ext:  ".ini",
		tmpl: template.Must(template.New("ffmeta").Funcs(funcMap).Parse(ffmetaTmpl)),
	},
	Fmt{
		name: "markdown",
		ext:  ".md",
		hash: true,
		tmpl: template.Must(template.New("md").Funcs(funcMap).Parse(mdTmpl)),
	},
	Fmt{
		name: "md",
		ext:  ".md",
		hash: true,
		tmpl: template.Must(template.New("md").Funcs(funcMap).Parse(mdTmpl)),
	},
	Fmt{
		name: "plain",
		ext:  ".txt",
		tmpl: template.Must(template.New("plain").Funcs(funcMap).Parse(plainTmpl)),
	},
	Fmt{
		name: "opf",
		ext:  ".opf",
	},
	Fmt{
		name: "ini",
		ext:  ".ini",
	},
	Fmt{
		name: "toml",
		ext:  ".toml",
		hash: true,
	},
	Fmt{
		name: "rss",
	},
}

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
