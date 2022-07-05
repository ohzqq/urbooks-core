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
	"strconv"
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
		var b = make(map[string]Meta)
		for key, val := range book {
			field := Field{Field: lib.DB.GetField(key), query: response.books.query}
			cat := Category{Field: field}
			item := Item{Field: field}
			col := Column{Field: field}
			var err error
			switch key {
			case "cover", "series", "publishers":
				err = json.Unmarshal(val, &item.meta)
				u := &url.URL{Path: item.Get("uri"), RawQuery: response.books.query.Encode()}
				item.Set("url", u.String())
				b[key] = item
			case "authors", "narrators", "identifiers", "formats", "languages", "tags":
				err = json.Unmarshal(val, &cat.items)

				for _, item := range cat.items {
					u := &url.URL{Path: item.Get("uri"), RawQuery: response.books.query.Encode()}
					item.Set("url", u.String())
				}

				b[key] = cat
			default:
				err = json.Unmarshal(val, &col.meta)
				b[key] = col
			}
			if err != nil {
				fmt.Println("error:", err)
			}
		}
		response.books.Add(&Book{meta: b, query: response.books.query})
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
	field Field
	label string
}

type BookMeta map[string]Meta

type Field struct {
	*calibredb.Field
	query url.Values
	null  bool
}

func (f Field) IsMultiple() bool {
	return f.IsMultiple()
}

func (f Field) IsCategory() bool {
	return f.IsCategory()
}

func (f Field) IsCustom() bool {
	return f.IsCustom()
}

func NewBook() BookMeta {
	return make(BookMeta)
}

func (meta BookMeta) Get(k string) Meta {
	return meta[k]
}

func (meta *BookMeta) Set(k string, v Meta) *BookMeta {
	m := *meta
	m[k] = v
	return meta
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
	if field := b.GetField(f); field.FieldMeta().IsCategory() && !field.FieldMeta().IsMultiple() {
		return field.(Item)
	}
	return Item{}
}

func (b Book) GetCategory(f string) Category {
	if field := b.GetField(f); field.FieldMeta().IsCategory() && field.FieldMeta().IsMultiple() {
		return field.(Category)
	}
	return Category{}
}

func (b Book) SetCategory(f, val string) *Book {
	b.meta[f] = Category{}
	return &b
}

func (b Book) GetColumn(f string) Column {
	if field := b.GetField(f); !field.FieldMeta().IsCategory() {
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
		if !b.field.IsMultiple() && b.field.IsCategory() {
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

	if field.FieldMeta().IsCategory() && !field.FieldMeta().IsMultiple() && field.IsNull() {
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
		bfile = b.Get("cover").Items()[0]
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

type Meta interface {
	Value() string
	String() string
	URL() string
	FieldMeta() Field
	//IsMultiple() bool
	//IsCategory() bool
	//IsCustom() bool
	IsNull() bool
}

type Category struct {
	Field
	items []*Item
	value string
	item  Item
}

func NewCategory() *Category {
	return &Category{
		Field: Field{Field: &calibredb.Field{}},
	}
}

func (c *Category) AddItem(i *Item) {
	c.items = append(c.items, i)
}

func (c Category) String() string {
	return c.Join("value")
}

const (
	nameSep = " & "
	itemSep = ", "
)

func (c Category) Join(v string) string {
	var meta []string
	for _, field := range c.items {
		meta = append(meta, field.Get(v))
	}
	switch c.IsNames {
	case true:
		return strings.Join(meta, nameSep)
	default:
		return strings.Join(meta, itemSep)
	}
}

func (c *Category) Split() {
	sep := itemSep
	if c.IsNames {
		sep = nameSep
	}
	for _, val := range strings.Split(c.value, sep) {
		i := Item{
			meta: map[string]string{"value": val},
		}
		c.AddItem(&i)
	}
}

func (c Category) IsNull() bool {
	return len(c.Items()) == 0
}

func (c Category) FieldMeta() Field {
	return c.Field
}

func (c Category) Value() string {
	return c.Join("value")
}

func (c Category) Items() []*Item {
	return c.items
}

func (c Category) URL() string {
	return c.Label + "/"
}

func (c *Category) SetField(k, v string) *Category {
	switch k {
	case "isNames":
		if v == "true" {
			c.IsNames = true
		}
	}
	return c
}

type Item struct {
	Field
	meta map[string]string
}

func NewCategoryItem() *Item {
	return &Item{meta: make(map[string]string)}
}

func (i Item) Value() string {
	return i.Get("value")
}

func (i Item) String() string {
	return i.Get("value")
}

func (i Item) FieldMeta() Field {
	return i.Field
}

//func (i Item) IsMultiple() bool {
//  return i.field.IsMultiple
//}

//func (i Item) IsCustom() bool {
//  return i.field.IsCustom
//}

//func (i Item) IsCategory() bool {
//  return i.field.IsCategory
//}

func (i Item) IsNull() bool {
	return len(i.meta) == 0
}

func (i Item) ID() string {
	return i.Get("id")
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

func (i Item) URL() string {
	return i.Get("url")
}

//func (i *Item) Query() url.Values {
//  return i.Query()
//}

func (i Item) Get(v string) string {
	return i.meta[v]
}

func (i *Item) Set(k, v string) *Item {
	i.meta[k] = v
	return i
}

type Column struct {
	Field
	meta string
}

func NewColumn(v string) Column {
	return Column{meta: v}
}

func (c Column) URL() string {
	return ""
}

func (c Column) Value() string {
	return c.meta
}

func (c Column) String() string {
	return c.meta
}

func (c Column) FieldMeta() Field {
	return c.Field
}

//func (c Column) IsMultiple() bool {
//  return c.field.IsMultiple
//}

//func (c Column) IsCustom() bool {
//  return c.field.IsCustom
//}

//func (c Column) IsCategory() bool {
//  return c.field.IsCategory
//}

func (c Column) IsNull() bool {
	return c.meta == ""
}

type CatResponse struct {
	Response
	data Category
}

func ParseCategory(r []byte) *CatResponse {
	var (
		response CatResponse
		err      error
	)

	var resp map[string]json.RawMessage
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println("error:", err)
	}

	response.Response = ParseResponse(resp)

	response.data.Field.query = url.Values{}
	response.data.Field.query.Set("library", response.Meta["library"])

	lib := Lib(response.Meta["library"])

	cats := Category{
		Field: Field{
			Field: lib.DB.GetField(response.Meta["endpoint"]),
			query: response.data.Field.query,
		},
	}
	err = json.Unmarshal(resp["data"], &cats.items)
	if err != nil {
		fmt.Println("error:", err)
	}

	return &response
}

func (c CatResponse) Items() []*Item {
	return c.data.Items()
}

func (c CatResponse) Label() string {
	return c.data.Field.Label
}

//func (b *Book) ToFFmeta() {
//  err := MetaFmt.FFmeta.Execute(os.Stdout, b)
//  if err != nil {
//    log.Fatal(err)
//  }
//}

//func (b *Book) ToPlain() string {
//  var buf bytes.Buffer
//  err := MetaFmt.Plain.Execute(&buf, b)
//  if err != nil {
//    log.Fatal(err)
//  }
//  return buf.String()
//}

//func (b *Book) ToMarkdown() string {
//  var buf bytes.Buffer
//  err := MetaFmt.MD.Execute(&buf, b)
//  if err != nil {
//    log.Fatal(err)
//  }
//  //fmt.Println(markdown)
//  return buf.String()
//}

func AudioFormats() []string {
	return []string{"m4b", "m4a", "mp3", "opus", "ogg"}
}

func AudioMimeType(ext string) string {
	switch ext {
	case "m4b", "m4a":
		return "audio/mp4"
	case "mp3":
		return "audio/mpeg"
	case "ogg", "opus":
		return "audio/ogg"
	}
	return ""
}

func BookSortFields() []string {
	return []string{
		"added",
		"sortAs",
	}
}

func BookSortTitle(idx int) string {
	titles := []string{
		"by Newest",
		"by Title",
	}
	return titles[idx]
}

func BookCats() []string {
	return []string{
		"authors",
		"tags",
		"series",
		"languages",
		"rating",
	}
}

func BookCatsTitle(idx int) string {
	titles := []string{
		"by Authors",
		"by Tags",
		"by Series",
		"by Languages",
		"by Rating",
	}
	return titles[idx]
}
