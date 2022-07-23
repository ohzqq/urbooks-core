package book

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gosimple/slug"
	"golang.org/x/exp/slices"
)

type Books []*Book

func ParseBooks(r []byte) (Books, error) {
	var err error

	var resp map[string]json.RawMessage
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response error: %v\n", err)
	}

	var rmeta map[string]string
	err = json.Unmarshal(resp["meta"], &rmeta)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response meta error: %v\n", err)
	}
	lib := rmeta["library"]

	var books Books
	err = json.Unmarshal(resp["data"], &books)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal books: %v\n", err)
	}

	for _, b := range books.EachBook() {
		b.lib = lib
		for _, field := range b.EachField() {
			field.Library = lib
		}
	}
	return books, nil
}

func (books *Books) UnmarshalJSON(r []byte) error {
	var err error

	var rawbooks []map[string]json.RawMessage
	err = json.Unmarshal(r, &rawbooks)
	if err != nil {
		return fmt.Errorf("book parsing error: %v\n", err)
	}

	for _, b := range rawbooks {
		book := NewBook()
		for key, value := range b {
			field := book.GetField(key)

			if key != "customColumns" {
				err := field.Meta.UnmarshalJSON(value)
				if err != nil {
					return err
				}
			}

			field.Data = value

			if key != field.JsonLabel {
				return fmt.Errorf("json: %v\n field meta: %v\n", key, field.JsonLabel)
			}

			if key == "customColumns" {
				var custom = make(map[string]map[string]json.RawMessage)
				err = json.Unmarshal(value, &custom)
				if err != nil {
					return fmt.Errorf("custom column parsing error: %v\n", err)
				}

				for name, cdata := range custom {
					col := &Field{
						IsCustom:     true,
						Data:         cdata["data"],
						CalibreLabel: name,
						JsonLabel:    name,
						IsEditable:   true,
					}

					meta := make(map[string]string)
					err = json.Unmarshal(cdata["meta"], &meta)
					if err != nil {
						return fmt.Errorf("custom column parsing error: %v\n", err)
					}

					switch meta["is_multiple"] {
					case "true":
						col.IsMultiple = true
						col.IsCollection = true
						col.Meta = &Collection{}
					case "false":
						col.IsColumn = true
						col.Meta = NewMetaColumn()
					}
					err = col.Meta.UnmarshalJSON(col.Data)
					if err != nil {
						return err
					}

					if meta["is_names"] == "true" {
						col.IsNames = true
					}

					book.AddField(col)
				}
			}
		}
		books.AddBook(book)
	}
	return nil
}

func (b *Books) AddBook(book *Book) *Books {
	*b = append(*b, book)
	return b
}

func (b *Books) EachBook() []*Book {
	return *b
}

type Book struct {
	lib string
	fmt metaFmt
	*Fields
}

type Meta interface {
	String(f *Field) string
	URL(f *Field) string
	IsNull() bool
	UnmarshalJSON(b []byte) error
}

func NewBook() *Book {
	return &Book{Fields: NewFields()}
}

func (b Book) GetFile(f string) *Item {
	formats := b.GetField("formats")
	switch f {
	case "cover":
		return b.GetField("cover").Item()
	case "audio":
		for _, item := range formats.Collection().EachItem() {
			if slices.Contains(AudioFormats(), item.Get("extension")) {
				q := url.Values{}
				q.Set("library", b.lib)
				q.Set("format", item.Get("extension"))
				u := url.URL{Path: item.Get("uri"), RawQuery: q.Encode()}
				item.Set("url", u.String())
				return item
			}
		}
	default:
		for _, item := range formats.Collection().EachItem() {
			if item.Get("extension") == f {
				return item
			}
		}
	}
	return &Item{}
}

type Collection struct {
	data []*Item
}

func NewMetaCollection() *Collection {
	return &Collection{}
}

func (c *Collection) AddItem() *Item {
	item := NewMetaItem()
	c.data = append(c.data, item)
	return item
}

func (c *Collection) EachItem() []*Item {
	return c.data
}

const (
	nameSep    = " & "
	itemSep    = ", "
	cliItemSep = `,`
	cliNameSep = `&`
)

func (c *Collection) String(f *Field) string {
	return c.Join(f.IsNames)
}

func (c *Collection) Join(isNames bool) string {
	var meta []string
	for _, item := range c.data {
		meta = append(meta, item.data["value"])
	}
	switch isNames {
	case true:
		return strings.Join(meta, nameSep)
	default:
		return strings.Join(meta, itemSep)
	}
}

func (c *Collection) Split(value string, names bool) *Collection {
	sep := itemSep
	if names {
		sep = nameSep
	}
	for _, val := range strings.Split(value, sep) {
		c.AddItem().Set("value", val)
	}
	return c
}

func (c *Collection) URL(f *Field) string {
	q := url.Values{}
	q.Set("library", f.Library)
	u := url.URL{Path: f.JsonLabel, RawQuery: q.Encode()}
	return u.String()
}

func (c *Collection) IsNull() bool {
	return len(c.data) == 0
}

func (c *Collection) UnmarshalJSON(b []byte) error {
	if len(b) > 0 {
		if err := json.Unmarshal(b, &c.data); err != nil {
			fmt.Printf("poot failed: %v\n", err)
			return err
		}
	}
	return nil
}

type Item struct {
	data map[string]string
}

func NewMetaItem() *Item {
	return &Item{data: make(map[string]string)}
}

func (i *Item) Get(val string) string {
	if v := i.data[val]; v != "" {
		return v
	}
	return ""
}

func (i *Item) Set(k, v string) *Item {
	i.data[k] = v
	return i
}

func (i *Item) String(f *Field) string {
	return i.Get("value")
}

func (i *Item) URL(f *Field) string {
	return i.Get("uri")
}

func (i *Item) IsNull() bool {
	return len(i.data) == 0
}

func (i Item) TotalBooks() int {
	if t := i.Get("books"); t != "" {
		b, err := strconv.Atoi(t)
		if err != nil {
			return 0
		}
		return b
	}
	return 0
}

func (i *Item) UnmarshalJSON(b []byte) error {
	if len(b) > 0 {
		i.data = make(map[string]string)
		if err := json.Unmarshal(b, &i.data); err != nil {
			fmt.Printf("collection failed: %v\n", err)
			return err
		}
	}
	return nil
}

type Column string

func NewMetaColumn() *Column {
	ms := Column("")
	return &ms
}

func (c *Column) String(f *Field) string {
	if len(f.Data) > 0 {
		var s string
		if err := json.Unmarshal(f.Data, &s); err != nil {
			fmt.Printf("%v failed: %v\n", f.JsonLabel, err)
		}
		return s
	}

	return string(*c)
}

func (c *Column) URL(f *Field) string {
	return ""
}

func (c *Column) IsNull() bool { return string(*c) == "" }

func (c *Column) UnmarshalJSON(b []byte) error { return nil }

func (c *Column) Set(v string) *Column {
	s := Column(v)
	return &s
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

func (b *Book) StringMap() map[string]string {
	m := make(map[string]string)
	for _, field := range b.EachField() {
		key := field.JsonLabel
		//if field.IsCustom {
		//  key = "#" + key
		//}

		if key != "customColumns" {
			m[key] = field.Meta.String(field)
		}

		if key == "series" {
			if pos := field.GetMeta().Item().Get("position"); pos != "" {
				m["position"] = pos
			}
		}
	}
	return m
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
	//case "opf":
	//return b.ToOpf()
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
