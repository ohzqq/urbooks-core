package urbooks

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/ohzqq/urbooks-core/calibredb"

	"golang.org/x/exp/slices"
)

type Books struct {
	query url.Values
	Books []*Book
}

type Book struct {
	meta  BookMeta
	label string
	fmt   metaFmt
}

func (b *Books) Add(book *Book) {
	b.Books = append(b.Books, book)
}

func NewBook(lib string) *Book {
	book := Book{meta: make(BookMeta)}
	if lib == "" {
		lib = DefaultLib().Name
	}
	book.Set("library", NewColumn().SetValue(lib))
	return &book
}

func (b *Book) NewColumn(k string) *Column {
	col := NewColumn()
	book := *b
	if lib := book.meta["library"]; !lib.IsNull() {
		col.Field = Lib(lib.Value()).DB.GetField(k)
	}
	book.meta[k] = col
	return col
}

func (b *Book) NewItem(k string) *Item {
	item := NewCategoryItem()
	book := *b
	if lib := book.meta["library"]; !lib.IsNull() {
		item.Field = Lib(lib.Value()).DB.GetField(k)
	}
	book.meta[k] = item
	return item
}

func (b *Book) NewCategory(k string) *Category {
	cat := NewCategory()
	book := *b
	if lib := book.meta["library"]; !lib.IsNull() {
		cat.Field = Lib(lib.Value()).DB.GetField(k)
	}
	book.meta[k] = cat
	return cat
}

func (b Book) Get(f string) Book {
	b.label = f
	return b
}

func (b *Book) Set(k string, v Meta) *Book {
	b.meta[k] = v
	return b
}

func (b Book) GetMeta() Meta {
	return b.meta[b.label]
}

func (b Book) FieldMeta() (*calibredb.Field, error) {
	f := b.meta.Get(b.label)
	if f.IsNull() {
		return &calibredb.Field{}, fmt.Errorf("field is null")
	}
	return f.FieldMeta(), nil
}

func (b Book) GetItem() *Item {
	return b.meta.GetItem(b.label)
}

func (b Book) GetCategory() *Category {
	return b.meta.GetCategory(b.label)
}

func (b Book) GetColumn() *Column {
	return b.meta.GetColumn(b.label)
}

func (b Book) URL() string {
	var u string
	field := b.GetMeta()
	switch {
	case b.label == "cover":
		if ur := b.meta[b.label].URL(); ur != "" {
			u = ur
		}
	case b.label == "":
		q := url.Values{}
		q.Set("library", b.Get("library").String())
		bu := &url.URL{Path: b.Get("uri").String(), RawQuery: q.Encode()}
		u = bu.String()
	default:
		if field.FieldMeta().Type() == "item" {
			u = b.meta[b.label].URL()
		}
	}
	return u
}

func (b Book) Value() string {
	return b.meta.Get(b.label).String()
}

func (b Book) String() string {
	return b.meta.String(b.label)
}

func (b *Book) StringMap() map[string]string {
	return b.meta.StringMap()
}

func (b Book) FilterValue() string {
	var filter []string
	for _, field := range []string{"title", "authors", "series"} {
		filter = append(filter, b.meta.Get(field).String())
	}
	return strings.Join(filter, " ")
}

func (b Book) Items() []*Item {
	return b.GetCategory().Items()
}

func (b Book) GetFile(f string) *Item {
	var bfile *Item
	formats := b.Get("formats").GetCategory()
	switch f {
	case "cover":
		f = ".jpg"
		fallthrough
	case "audio":
		for _, item := range formats.Items() {
			if slices.Contains(AudioFormats(), item.Get("extension")) {
				formats.query.Set("format", item.Get("extension"))
				u := url.URL{Path: item.Get("uri"), RawQuery: formats.query.Encode()}
				item.Set("url", u.String())
				bfile = item
				break
			}
		}
	default:
		for _, item := range formats.Items() {
			if item.Get("extension") == f {
				bfile = item
			}
		}
	}
	return bfile
}

type Category struct {
	Field *calibredb.Field
	query url.Values
	items []*Item
	value string
	item  Item
}

const (
	nameSep    = " & "
	itemSep    = ", "
	cliItemSep = `,`
	cliNameSep = `&`
)

func NewCategory() *Category {
	cat := &Category{}
	return cat
}

func (c *Category) AddItem() *Item {
	item := NewCategoryItem()
	c.items = append(c.items, item)
	return item
}

func (c Category) String() string {
	return c.Join("value")
}

func (c Category) Join(v string) string {
	var meta []string
	for _, field := range c.items {
		meta = append(meta, field.Get(v))
	}
	switch c.Field.IsNames {
	case true:
		return strings.Join(meta, nameSep)
	default:
		return strings.Join(meta, itemSep)
	}
}

func (c Category) Cli(v string) string {
	var meta []string
	for _, field := range c.items {
		meta = append(meta, field.Get(v))
	}
	switch c.Field.IsNames {
	case true:
		return strings.Join(meta, cliNameSep)
	default:
		return strings.Join(meta, cliItemSep)
	}
}

func (c *Category) Split(value string, names bool) {
	sep := itemSep
	if names {
		sep = nameSep
	}
	for _, val := range strings.Split(value, sep) {
		c.AddItem().SetValue(val)
	}
}

func (c Category) ItemStringSlice() []string {
	var i []string
	for _, item := range c.Items() {
		i = append(i, item.String())
	}
	return i
}

func (c *Category) SetFieldMeta(k, v string) *Category {
	switch k {
	case "isNames":
		if v == "true" {
			c.Field.IsNames = true
		}
	}
	return c
}

func (c Category) IsNull() bool                { return len(c.Items()) == 0 }
func (c Category) FieldMeta() *calibredb.Field { return c.Field }
func (c Category) Value() string               { return c.Join("value") }
func (c Category) Items() []*Item              { return c.items }
func (c Category) URL() string                 { return c.Field.Label + "/" }

type Item struct {
	Field *calibredb.Field
	query url.Values
	meta  map[string]string
}

func NewCategoryItem() *Item {
	return &Item{meta: make(map[string]string)}
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

func (i *Item) Set(k, v string) *Item {
	i.meta[k] = v
	return i
}

func (i *Item) SetValue(v string) *Item {
	i.meta["value"] = v
	return i
}

func (i Item) Value() string               { return i.Get("value") }
func (i Item) String() string              { return i.Get("value") }
func (i Item) FieldMeta() *calibredb.Field { return i.Field }
func (i Item) IsNull() bool                { return len(i.meta) == 0 }
func (i Item) ID() string                  { return i.Get("id") }
func (i Item) URL() string                 { return i.Get("url") }
func (i Item) Get(v string) string         { return i.meta[v] }

type Column struct {
	Field *calibredb.Field
	query url.Values
	meta  string
}

func NewColumn() *Column {
	return &Column{}
}

func (c *Column) SetValue(v string) *Column {
	c.meta = v
	return c
}

func (c Column) URL() string                 { return "" }
func (c Column) Value() string               { return c.meta }
func (c Column) String() string              { return c.meta }
func (c Column) FieldMeta() *calibredb.Field { return c.Field }
func (c Column) IsNull() bool                { return c.meta == "" }

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
