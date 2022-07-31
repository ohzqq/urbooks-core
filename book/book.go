package book

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Books []*Book

type Book struct {
	lib string
	fmt Fmt
	*Fields
}

func ParseBooks(r []byte) (Books, error) {
	var err error

	var resp map[string]json.RawMessage
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return nil, fmt.Errorf("pkg book unmarshal response error: %v\n", err)
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

			if key != field.JsonLabel {
				return fmt.Errorf("json: %v\n field meta: %v\n", key, field.JsonLabel)
			}

			switch key {
			case "customColumns":
				var custom = make(map[string]map[string]json.RawMessage)
				err = json.Unmarshal(value, &custom)
				if err != nil {
					return fmt.Errorf("custom column parsing error: %v\n", err)
				}

				for name, cdata := range custom {
					field.Collection().AddItem().Set("value", name)
					book.customColumns = append(book.customColumns, name)

					meta := make(map[string]string)
					err = json.Unmarshal(cdata["meta"], &meta)
					if err != nil {
						return fmt.Errorf("custom column parsing error: %v\n", err)
					}

					var col *Field
					switch meta["is_multiple"] {
					case "true":
						col = book.AddField(NewCollection(name))
					case "false":
						col = book.AddField(NewColumn(name))
					}
					col.SetIsCustom().SetIsEditable()

					col.SetMeta(cdata["data"])

					if meta["is_names"] == "true" {
						col.SetIsNames()
					}
				}
			default:
				field.SetMeta(value)
			}
		}
		books.AddBook(book)
	}

	for _, b := range books.EachBook() {
		for _, field := range b.EachField() {
			field.Library = b.GetField("library").String()
		}
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

type Meta interface {
	String(f *Field) string
	URL(f *Field) string
	IsNull() bool
	ParseMeta(f *Field) Meta
	RawData() interface{}
}

func NewBook() *Book {
	return &Book{Fields: NewFields()}
}

func (b *Book) GetSeriesString() string {
	s := b.GetField("series")
	if !s.IsNull() {
		p := "1.0"
		if s.String() != "" {
			if pos := b.GetField("position").String(); pos != "" {
				p = pos
			}
			if pos := s.GetMeta().Item().Get("position"); pos != "" {
				p = pos
			}
		}
		return s.String() + `, Book ` + p
	}
	return ""
}

func (b *Book) GetTitleAndSeries() string {
	title := b.GetField("title").String()
	if t := b.GetField("titleAndSeries"); !t.IsNull() {
		return t.String()
	}
	if s := b.GetSeriesString(); s != "" {
		return title + " [" + s + "]"
	}
	return title
}

func (b Book) GetFile(f string) *Item {
	formats := b.GetField("formats")
	item := NewMetaItem()
	q := url.Values{}
	q.Set("library", b.GetField("library").String())
	switch f {
	case "cover":
		item = b.GetField("cover").Item()
	case "audio":
		for _, i := range formats.Collection().EachItem() {
			if slices.Contains(AudioFormats(), i.Get("extension")) {
				q.Set("format", i.Get("extension"))
				item = i
				break
			}
		}
	default:
		for _, i := range formats.Collection().EachItem() {
			if i.Get("extension") == f {
				item = i
			}
		}
	}
	if uri := item.Get("uri"); uri != "" {
		u := url.URL{Path: uri, RawQuery: q.Encode()}
		item.Set("url", u.String())
	}
	return item
}

func (b Book) FilterValue() string {
	var filter []string
	for _, field := range []string{"title", "authors", "series"} {
		filter = append(filter, b.GetField(field).String())
	}
	return strings.Join(filter, " ")
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

func (c *Collection) StringSlice() []string {
	var items []string
	for _, i := range c.EachItem() {
		items = append(items, i.Get("value"))
	}
	return items
}

func (c *Collection) RawData() interface{} {
	var items []interface{}
	for _, i := range c.EachItem() {
		items = append(items, i.RawData())
	}
	return items
}

func (c *Collection) splitString(value string, isNames bool) *Collection {
	sep := itemSep
	if isNames {
		sep = nameSep
	}
	for _, val := range strings.Split(value, sep) {
		c.AddItem().Set("value", val)
	}
	return c
}

func (c *Collection) URL(f *Field) string {
	u := url.URL{Path: f.Label(), RawQuery: f.rawQuery()}
	return u.String()
}

func (c *Collection) IsNull() bool {
	return len(c.data) == 0
}

func (c *Collection) ParseMeta(f *Field) Meta {
	switch d := f.data.(type) {
	case string:
		c.splitString(d, f.IsNames)
	case []string:
		for _, val := range d {
			c.AddItem().Set("value", val)
		}
	case json.RawMessage:
		if len(d) > 0 {
			if err := json.Unmarshal(d, &c.data); err != nil {
				log.Fatalf("poot failed: %v\n", err)
			}
		}
	}
	return c
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
	u := url.URL{Path: i.Get("uri"), RawQuery: f.rawQuery()}
	return u.String()
}

func (i *Item) IsNull() bool {
	return len(i.data) == 0
}

func (i *Item) RawData() interface{} {
	return i.data
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

func (i *Item) ParseMeta(f *Field) Meta {
	switch d := f.data.(type) {
	case string:
		i = i.Set("value", d)
	case json.RawMessage:
		err := i.UnmarshalJSON(d)
		if err != nil {
			fmt.Printf("item failed: %v\n", err)
		}
	}
	return i
}

func (i *Item) UnmarshalJSON(b []byte) error {
	if len(b) > 0 {
		if err := json.Unmarshal(b, &i.data); err != nil {
			fmt.Printf("collection failed: %v\n", err)
			return err
		}
	}
	return nil
}

type Column struct {
	data string
}

func NewMetaColumn() *Column {
	//ms := Column("")
	//return &ms
	return &Column{}
}

func (c *Column) String(f *Field) string {
	if f.Label() == "uri" {
		u := url.URL{Path: c.data, RawQuery: f.rawQuery()}
		return u.String()
	}
	return c.data
	//return string(*c)
}

func (c *Column) URL(f *Field) string {
	return ""
}

func (c *Column) IsNull() bool {
	return c.data == ""
	//return string(*c) == ""
}

func (c *Column) RawData() interface{} {
	return c.data
	//return string(*c)
}

func (c *Column) ParseMeta(f *Field) Meta {
	switch d := f.data.(type) {
	case string:
		c.data = d
		//s := Column(d)
		//return &s
	case json.RawMessage:
		if len(d) > 0 {
			if err := json.Unmarshal(d, &c.data); err != nil {
				fmt.Printf("%v failed: %v\n", d, err)
			}
		}
	}
	return c
}

func (c *Column) Set(v string) *Column {
	c.data = v
	//s := Column(v)
	return c
}

var EditableFields = []string{
	"id",
	"authorSort",
	"authors",
	"description",
	"identifiers",
	"languages",
	"published",
	"publisher",
	"rating",
	"series",
	"position",
	"sortAs",
	"tags",
	"added",
	"title",
}
