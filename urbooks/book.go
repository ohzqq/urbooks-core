package urbooks

import (
	"net/url"
	"strings"

	"github.com/ohzqq/urbooks-core/calibredb"

	"golang.org/x/exp/slices"
)

type Books struct {
	query url.Values
	Books []*Book
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

func (b *Book) StringMap() map[string]string {
	return b.meta.StringMap()
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

type Meta interface {
	Value() string
	String() string
	URL() string
	FieldMeta() *calibredb.Field
	IsNull() bool
}

type BookMeta map[string]Meta

func NewBookMeta(m map[string]string) BookMeta {
	meta := make(BookMeta)
	for key, val := range m {
		meta[key] = MetaString(val)
	}
	return meta
}

func (meta BookMeta) Get(k string) Meta {
	return meta[k]
}

func (bm BookMeta) StringMap() map[string]string {
	m := make(map[string]string)
	for key, val := range bm {
		m[key] = val.String()
		if key == "series" {
			if pos := bm.Get("series").(*Item).Get("position"); pos != "" {
				m["position"] = pos
			}
		}
	}
	return m
}

func (bm BookMeta) StringMapToBook() *Book {
	lib := DefaultLib()
	if l := bm["library"].Value(); l == "" {
		lib = Lib(l)
	}
	book := NewBook(lib.Name)
	for key, val := range bm {
		field := lib.DB.GetField(key)
		switch {
		case field.IsCategory:
			switch field.IsMultiple {
			case true:
				cat := book.NewCategory(key)
				switch {
				case field.IsNames:
					cat.Split(val.String(), true)
				default:
					cat.Split(val.String(), false)
				}
			case false:
				item := book.NewItem(key).SetValue(val.String())
				if key == "series" {
					if pos := bm.Get("position").String(); pos != "" {
						item.Set("position", pos)
					}
				}
			}
		default:
			book.NewColumn(key).SetValue(val.String())
		}
	}
	return book
}

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

func (b Book) GetItem(f string) *Item {
	if field := b.GetField(f); field.FieldMeta().IsCategory && !field.FieldMeta().IsMultiple {
		return field.(*Item)
	}
	return &Item{}
}

func (b Book) GetCategory(f string) *Category {
	if field := b.GetField(f); field.FieldMeta().IsCategory && field.FieldMeta().IsMultiple {
		return field.(*Category)
	}
	return &Category{}
}

func (b Book) GetColumn(f string) *Column {
	if field := b.GetField(f); !field.FieldMeta().IsCategory {
		return field.(*Column)
	}
	return &Column{}
}

func (b Book) URL() string {
	var u string
	switch {
	case b.label == "cover":
		if ur := b.meta[b.label].URL(); ur != "" {
			u = ur
		}
	case b.label == "":
		bu := &url.URL{Path: b.Get("uri").String(), RawQuery: b.query.Encode()}
		u = bu.String()
	default:
		if !b.field.IsMultiple && b.field.IsCategory {
			u = b.meta[b.label].URL()
		}
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
	if b.label == "titleAndSeries" {
		field = b.GetField("series")
	}

	switch b.label {
	case "formats":
		return b.GetCategory(b.label).Join("extension")
	case "position":
		if series := b.GetItem("series"); series.IsNull() {
			return series.Get("position")
		}
	case "titleAndSeries":
		title := b.Get("title").Value()
		if series := b.Get("series"); !field.IsNull() {
			return title + " [" + series.String() + "]"
		}
		return title
	}

	if field.FieldMeta().IsCategory && !field.FieldMeta().IsMultiple && !field.IsNull() {
		f := field.(*Item)
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
		f = ".jpg"
		fallthrough
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

func (f BookFile) Get(v string) string {
	return f[v]
}

func (f BookFile) Path() string {
	return f.Get("path")
}

func (f BookFile) Ext() string {
	return f.Get("extension")
}

func (f BookFile) URL() string {
	return f.Get("url")
}

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
