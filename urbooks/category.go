package urbooks

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ohzqq/urbooks-core/calibredb"
)

type Category struct {
	Field *calibredb.Field
	items []*Item
	value string
	item  Item
}

func NewCategory(label string) *Category {
	cat := &Category{Field: &calibredb.Field{}}
	cat.Field.Label = label
	return cat
}

//func (c *Category) AddItem(i *Item) {
//  c.items = append(c.items, i)
//}

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
	switch c.Field.IsNames {
	case true:
		return strings.Join(meta, nameSep)
	default:
		return strings.Join(meta, itemSep)
	}
}

func (c *Category) Split() {
	sep := itemSep
	if c.Field.IsNames {
		sep = nameSep
	}
	for _, val := range strings.Split(c.value, sep) {
		c.AddItem().SetValue(val)
	}
}

func (c Category) IsNull() bool {
	return len(c.Items()) == 0
}

func (c Category) FieldMeta() *calibredb.Field {
	return c.Field
}

func (c Category) Value() string {
	return c.Join("value")
}

func (c Category) Items() []*Item {
	return c.items
}

func (c Category) URL() string {
	return c.Field.Label + "/"
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

	//response.data.Field.query = url.Values{}
	//response.data.Field.query.Set("library", response.Meta["library"])

	lib := Lib(response.Meta["library"])

	cats := Category{
		Field: lib.DB.GetField(response.Meta["endpoint"]),
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
