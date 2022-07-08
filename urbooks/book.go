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

type Book struct {
	meta  BookMeta
	label string
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

func (b Book) FieldMeta() *calibredb.Field {
	return b.meta.Get(b.label).FieldMeta()
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
