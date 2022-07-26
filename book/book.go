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
			field.SetData(value)

			if key != field.JsonLabel {
				return fmt.Errorf("json: %v\n field meta: %v\n", key, field.JsonLabel)
			}

			if key != "customColumns" {
				field.ParseData()
			}

			if key == "customColumns" {
				var custom = make(map[string]map[string]json.RawMessage)
				err = json.Unmarshal(value, &custom)
				if err != nil {
					return fmt.Errorf("custom column parsing error: %v\n", err)
				}

				for name, cdata := range custom {
					col := NewField(name).SetIsCustom().SetData(cdata["data"])

					meta := make(map[string]string)
					err = json.Unmarshal(cdata["meta"], &meta)
					if err != nil {
						return fmt.Errorf("custom column parsing error: %v\n", err)
					}

					switch meta["is_multiple"] {
					case "true":
						col.SetIsMultiple()
						col.SetMeta(NewMetaCollection())
					case "false":
						col.SetMeta(NewMetaColumn())
					}
					col.ParseData()
					//if err != nil {
					//  return err
					//}

					if meta["is_names"] == "true" {
						col.SetIsNames()
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
	ParseData(f *Field)
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
				q.Set("library", b.GetField("library").String())
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

func (c *Collection) Split(value string, isNames bool) *Collection {
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
	q := url.Values{}
	q.Set("library", f.Library)
	u := url.URL{Path: f.JsonLabel, RawQuery: q.Encode()}
	return u.String()
}

func (c *Collection) IsNull() bool {
	return len(c.data) == 0
}

func (c *Collection) ParseData(f *Field) {
	switch d := f.data.(type) {
	case string:
		c.Split(d, f.IsNames)
	case json.RawMessage:
		if len(d) > 0 {
			if err := json.Unmarshal(d, &c.data); err != nil {
				log.Fatalf("poot failed: %v\n", err)
			}
		}
	}
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

func (i *Item) ParseData(f *Field) {
	switch d := f.data.(type) {
	case string:
		i = i.Set("value", d)
	case json.RawMessage:
		err := i.UnmarshalJSON(d)
		if err != nil {
			fmt.Printf("item failed: %v\n", err)
		}
	}
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

type Column string

func NewMetaColumn() *Column {
	ms := Column("")
	return &ms
}

func (c *Column) String(f *Field) string {
	return string(*c)
}

func (c *Column) URL(f *Field) string {
	return ""
}

func (c *Column) IsNull() bool {
	return string(*c) == ""
}

func (c *Column) ParseData(f *Field) {
	switch d := f.data.(type) {
	case string:
		c.Set(d)
	case json.RawMessage:
		if len(d) > 0 {
			if err := json.Unmarshal(d, &c); err != nil {
				fmt.Printf("%v failed: %v\n", d, err)
			}
		}
	}
}

func (c *Column) Set(v string) *Column {
	s := Column(v)
	return &s
}
