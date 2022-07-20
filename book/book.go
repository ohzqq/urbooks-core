package book

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Books []*Book

type Book struct {
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

type Collection struct {
	data []*Item
}

func NewCollection() *Collection {
	return &Collection{}
}

func (c *Collection) AddItem() *Item {
	item := NewItem()
	c.data = append(c.data, item)
	return item
}

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
	if err := json.Unmarshal(b, &c.data); err != nil {
		fmt.Printf("collection failed: %v\n", err)
		return err
	}
	return nil
}

type Item struct {
	data map[string]string
}

func NewItem() *Item {
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

func (i *Item) UnmarshalJSON(b []byte) error {
	i.data = make(map[string]string)
	if err := json.Unmarshal(b, &i.data); err != nil {
		fmt.Printf("collection failed: %v\n", err)
		return err
	}
	return nil
}

type Column string

func NewColumn() *Column {
	ms := Column("")
	return &ms
}
func (c *Column) String(f *Field) string {
	return string(f.Data)
}

func (c *Column) URL(f *Field) string {
	return ""
}

func (c *Column) IsNull() bool {
	return string(*c) == ""
}

func (c *Column) UnmarshalJSON(b []byte) error {
	return nil
}

func (c *Column) SetValue(v string) *Column {
	s := Column(v)
	c = &s
	return c
}
