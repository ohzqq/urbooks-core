package book

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/exp/slices"
)

type Fields struct {
	displayFields [][]interface{}
	lib           string
	DbMeta        map[string]*Field
	data          fields
	customColumns []string
}

type fields map[string]*Field

type Field struct {
	idx          int
	Library      string `json:"-"`
	IsDisplayed  bool   `json:"-"`
	IsNames      bool   `json:"-"`
	IsMultiple   bool   `json:"-"`
	CatID        string `json:"-"`
	JsonLabel    string `json:"-"`
	jsonData     []byte `json:"-"`
	stringData   string `json:"-"`
	data         any    `json:"-"`
	Data         []byte `json:"-"`
	Meta         Meta   `json:"-"`
	CalibreLabel string `json:"label"`
	IsCategory   bool   `json:"is_category"`
	IsCustom     bool   `json:"is_custom"`
}

func NewFields() *Fields {
	return &Fields{
		data: defaultCalibreFields(),
	}
}

func (f *Fields) NewField(label string) *Field {
	field := NewField(label)
	f.AddField(field)
	return field
}

func (f *Fields) GetField(name string) *Field {
	if !strings.HasPrefix(name, "#") && slices.Contains(f.customColumns, "#"+name) {
		name = "#" + name
	}
	return f.data[name]
}

//func (f *Fields) GetMeta(name string) string {
//  if !strings.HasPrefix(name, "#") && slices.Contains(f.customColumns, "#"+name) {
//    name = "#" + name
//  }
//  return f.data[name]
//}

func (f *Fields) AddField(field *Field) *Field {
	f.data[field.Label()] = field
	return field
}

func (f *Fields) SetField(name string, field *Field) *Fields {
	f.data[name] = field
	return f
}

func (f *Fields) SetMeta(name string, meta Meta) *Fields {
	f.GetField(name).Meta = meta
	return f
}

func (f *Fields) GetMeta(name string) string {
	return f.GetField(name).String()
}

func (f *Fields) EachField() fields {
	return f.data
}

func (f *Fields) ListFields() []string {
	var fields []string
	for _, field := range f.EachField() {
		fields = append(fields, field.Label())
	}
	return fields
}

func NewField(label string) *Field {
	return &Field{
		JsonLabel:    label,
		CalibreLabel: label,
		IsCustom:     strings.HasPrefix(label, "#"),
	}
}

func NewCollection(label string) *Field {
	field := NewField(label).SetIsMultiple()
	field.Meta = NewMetaCollection()
	return field
}

func NewItem(label string) *Field {
	field := NewField(label)
	field.Meta = NewMetaItem()
	return field
}

func NewColumn(label string) *Field {
	field := NewField(label)
	field.Meta = NewMetaColumn()
	return field
}

func (f *Field) rawQuery() string {
	q := url.Values{}
	q.Set("library", f.Library)
	return q.Encode()
}

func (f Field) Label() string {
	return f.JsonLabel
}

func (f *Field) SetCalibreLabel(label string) *Field {
	f.CalibreLabel = label
	return f
}

func (f *Field) setJsonData(data []byte) *Field {
	f.jsonData = data
	return f
}

func (f *Field) SetData(data any) *Field {
	//f.data = data
	//switch d := data.(type) {
	//case string:
	//  f.SetStringMeta(d)
	//  //case json.RawMessage:
	//}
	//f.ParseData()
	f.SetMeta(data)
	//f.Meta = f.Meta.ParseMeta(f)
	return f
}

func (f *Field) ParseData() *Field {
	f.Meta.ParseData(f)
	return f
}

func (f *Field) setStringData(data string) *Field {
	f.stringData = data
	return f
}

func (f *Field) SetIsMultiple() *Field {
	f.IsMultiple = true
	return f
}

func (f *Field) SetIsNames() *Field {
	f.IsNames = true
	return f
}

func (f *Field) SetIsCustom() *Field {
	f.IsCustom = true
	return f
}

func (f *Field) SetIsCategory() *Field {
	f.IsCategory = true
	return f
}

func (f *Field) SetIndex(idx int) *Field {
	f.idx = idx
	return f
}

func (f *Field) SetMeta(data any) *Field {
	f.data = data
	f.Meta = f.Meta.ParseMeta(f)
	return f
}

func (f *Field) SetStringMeta(data string) *Field {
	var meta Meta
	switch {
	case f.IsCollection():
		meta = f.Collection().Split(data, f.IsNames)
	case f.IsItem():
		meta = f.Item().Set("value", data)
	case f.IsColumn():
		meta = f.Col().Set(data)
	}
	f.Meta = meta
	return f
}

func (f *Field) String() string {
	return f.Meta.String(f)
}

func (f *Field) URL() string {
	return f.Meta.URL(f)
}

func (f *Field) IsNull() bool {
	return f.Meta.IsNull()
}

func (f *Field) GetMeta() *Field {
	return f
}

func (f *Field) IsCollection() bool {
	if _, ok := f.Meta.(*Collection); ok {
		return true
	}
	return false
}

func (f *Field) IsItem() bool {
	if _, ok := f.Meta.(*Item); ok {
		return true
	}
	return false
}

func (f *Field) IsColumn() bool {
	if _, ok := f.Meta.(*Column); ok {
		return true
	}
	return false
}

func (f *Field) Collection() *Collection {
	return f.Meta.(*Collection)
}

func (f *Field) Item() *Item {
	return f.Meta.(*Item)
}

func (f *Field) Col() *Column {
	return f.Meta.(*Column)
}

func UnmarshalField(d []byte) (Field, error) {
	field := Field{}
	err := json.Unmarshal(d, &field)
	if err != nil {
		return field, fmt.Errorf("unmarshal field %v\n", err)
	}
	return field, nil
}

const (
	sort     = iota
	order    = iota
	limit    = iota
	category = iota
)

func GetDefaultField(name string) *Field {
	return defaultCalibreFields()[name]
}

func ListDefaultFields() []string {
	var fields []string
	for label, _ := range defaultCalibreFields() {
		fields = append(fields, label)
	}
	return fields
}

func defaultCalibreFields() fields {
	return fields{
		"authors": NewCollection("authors").
			SetIsCategory().
			SetIsMultiple().
			SetIsNames(),
		"authorSort": NewColumn("authorSort").
			SetCalibreLabel("author_sort"),
		"description": NewColumn("description").
			SetCalibreLabel("comments"),
		"cover": NewItem("cover"),
		"formats": NewCollection("formats").
			SetIsCategory().
			SetIsMultiple(),
		"id": NewColumn("id"),
		"identifiers": NewCollection("identifiers").
			SetIsCategory().
			SetIsMultiple(),
		"languages": NewCollection("languages").
			SetIsCategory().
			SetIsMultiple(),
		"library": NewColumn("library"),
		"modified": NewColumn("modified").
			SetCalibreLabel("last_modified"),
		"path": NewColumn("path"),
		"published": NewColumn("published").
			SetCalibreLabel("pubdate"),
		"publisher": NewItem("publisher"),
		"rating": NewColumn("rating").
			SetIsCategory(),
		"series": NewItem("series").
			SetIsCategory(),
		"position": NewColumn("position").
			SetCalibreLabel("series_index"),
		"sortAs": NewColumn("sortAs").
			SetCalibreLabel("sort"),
		"tags": NewCollection("tags").
			SetIsCategory().
			SetIsMultiple(),
		"added": NewColumn("added").
			SetCalibreLabel("timestamp"),
		"title":          NewColumn("title"),
		"titleAndSeries": NewColumn("titleAndSeries"),
		"uri":            NewColumn("uri"),
		"uuid":           NewColumn("uuid"),
		"customColumns": NewCollection("customColumns").
			SetCalibreLabel("custom_columns"),
	}
}
