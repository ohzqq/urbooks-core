package urbooks

import (
	"fmt"
	//"log"
	"encoding/json"
	"net/url"

	//"os"
	//"time"
	//"path"
	//"bytes"

	"strings"

	"github.com/ohzqq/urbooks-core/calibredb"

	"golang.org/x/exp/slices"
)

var _ = fmt.Sprintf("%v", "poot")

type BookResponse struct {
	Response
	books Books
}

func (b BookResponse) Books() []*Book {
	return b.books.Books
}

func (b BookResponse) GetMeta(k string) string {
	if m := b.Response.Meta[k]; m != "" {
		return m
	}
	return ""
}

func (b BookResponse) GetLink(k string) string {
	if m := b.Response.Links[k]; m != "" {
		return m
	}
	return ""
}

type Books struct {
	query url.Values
	Books []*Book
}

func ParseBooks(r []byte) *BookResponse {
	var (
		response BookResponse
		err      error
	)

	var resp map[string]json.RawMessage
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println("error:", err)
	}

	response.Response = ParseResponse(resp)

	var books []map[string]json.RawMessage
	err = json.Unmarshal(resp["data"], &books)
	if err != nil {
		fmt.Println("error:", err)
	}

	response.books.query = url.Values{}
	response.books.query.Set("library", response.GetMeta("library"))
	lib := Lib(response.GetMeta("library"))

	for _, book := range books {
		bb := NewBook(lib.Name)
		for key, val := range book {
			var err error
			switch key {
			case "cover", "series", "publishers":
				item := bb.meta.NewItem(key)
				err = json.Unmarshal(val, &item.meta)
				u := &url.URL{Path: item.Get("uri"), RawQuery: response.books.query.Encode()}
				item.Set("url", u.String())
			case "authors", "narrators", "identifiers", "formats", "languages", "tags":
				cat := bb.meta.NewCategory(key)
				err = json.Unmarshal(val, &cat.items)

				for _, item := range cat.items {
					u := &url.URL{Path: item.Get("uri"), RawQuery: response.books.query.Encode()}
					item.Set("url", u.String())
				}
			default:
				col := bb.meta.NewColumn(key)
				err = json.Unmarshal(val, &col.meta)
			}
			if err != nil {
				fmt.Printf("%v: %v\n", key, err)
			}
		}
		response.books.Add(bb)
	}
	return &response
}

func (i *Item) UnmarshalJSON(b []byte) error {
	i.meta = make(map[string]string)
	if err := json.Unmarshal(b, &i.meta); err != nil {
		return err
	}
	return nil
}

func (b *Books) Add(book *Book) {
	b.Books = append(b.Books, book)
}

type Book struct {
	query url.Values
	meta  BookMeta
	field *calibredb.Field
	label string
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
	cat := NewCategory(k)
	book := *b
	if lib := book.meta["library"]; !lib.IsNull() {
		cat.Field = Lib(lib.Value()).DB.GetField(k)
	}
	book.meta[k] = cat
	return cat
}

type Meta interface {
	Value() string
	String() string
	URL() string
	FieldMeta() *calibredb.Field
	IsNull() bool
}

type BookMeta map[string]Meta

func (meta BookMeta) Get(k string) Meta {
	return meta[k]
}

func (b Book) Get(f string) Book {
	b.label = f
	return b
}

func (b *Book) Set(k string, v Meta) *Book {
	b.meta[k] = v
	return b
}

func (b Book) GetField(f string) Meta {
	return b.meta[f]
}

func (b Book) GetItem(f string) Item {
	if field := b.GetField(f); field.FieldMeta().IsCategory && !field.FieldMeta().IsMultiple {
		return field.(Item)
	}
	return Item{}
}

func (b Book) GetCategory(f string) Category {
	if field := b.GetField(f); field.FieldMeta().IsCategory && field.FieldMeta().IsMultiple {
		return field.(Category)
	}
	return Category{}
}

func (b Book) SetCategory(f, val string) *Book {
	b.meta[f] = Category{}
	return &b
}

func (b Book) GetColumn(f string) Column {
	if field := b.GetField(f); !field.FieldMeta().IsCategory {
		return field.(Column)
	}
	return Column{}
}

func (b Book) SetColumn(f, val string) *Book {
	b.meta[f] = Column{
		meta: val,
	}
	return &b
}

func (b Book) URL() string {
	var u string
	switch b.label == "" {
	case false:
		if !b.field.IsMultiple && b.field.IsCategory {
			u = b.meta[b.label].URL()
		}
	case true:
		bu := &url.URL{Path: b.Get("uri").String(), RawQuery: b.query.Encode()}
		u = bu.String()
	}
	return u
}

func (b Book) Value() string {
	return b.Get(b.label).String()
}

func (b Book) FieldMeta() *calibredb.Field {
	return b.GetField(b.label).FieldMeta()
}

func (b Book) FilterValue() string {
	var filter []string
	for _, field := range []string{"title", "authors", "series"} {
		filter = append(filter, b.Get(field).String())
	}
	return strings.Join(filter, " ")
}

func (b Book) Items() []*Item {
	col := b.GetCategory(b.label)
	return col.items
}

func (b Book) String() string {
	field := b.GetField(b.label)

	switch b.label {
	case "formats":
		return b.GetCategory(b.label).Join("extension")
	case "position":
		if series := b.GetItem("series"); series.IsNull() {
			return series.Get("position")
		}
	case "seriesAndTitle":
		title := b.Get("title").Value()
		if series := b.GetItem("series"); series.IsNull() {
			return title + " [" + b.Get("series").String() + "]"
		}
		return title
	}

	if field.FieldMeta().IsCategory && !field.FieldMeta().IsMultiple && field.IsNull() {
		f := field.(Item)
		if b.label == "series" {
			return f.Value() + ", Book " + f.Get("position")
		}
		return f.Value()
	}

	return field.Value()
}

type BookFile map[string]string

func (b Book) GetFile(f string) BookFile {
	var bfile *Item
	switch f {
	case "cover":
		i := b.GetItem("cover")
		bfile = &i
	case "audio":
		for _, item := range b.Get("formats").Items() {
			if slices.Contains(AudioFormats(), item.Get("extension")) {
				b.query.Set("format", item.Get("extension"))
				u := url.URL{Path: item.Get("uri"), RawQuery: b.query.Encode()}
				item.Set("url", u.String())
				bfile = item
				break
			}
		}
	default:
		for _, item := range b.Get("formats").Items() {
			if item.Get("extension") == f {
				bfile = item
			}
		}
	}
	return BookFile(bfile.meta)
}

func (f BookFile) Path() string {
	return f.Get("path")
}

func (f BookFile) Get(v string) string {
	return f[v]
}

func (f BookFile) Ext() string {
	return f.Get("extension")
}

func (f BookFile) URL() string {
	return f.Get("url")
}
