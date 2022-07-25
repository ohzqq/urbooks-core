package book

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Fields struct {
	displayFields [][]interface{}
	lib           string
	DbMeta        map[string]*Field
	meta          []*Field
	idx           map[string]int
}

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
		meta: defaultFields(),
		idx:  defaultFieldsIdx(),
	}
}

func (f *Fields) NewField(label string) *Field {
	field := NewField(label)
	f.AddField(field)
	return field
}

func (f *Fields) AddField(field *Field) *Fields {
	f.idx[field.CalibreLabel] = len(f.meta)
	f.meta = append(f.meta, field)
	return f
}

func (f *Fields) GetSeriesString() string {
	s := f.GetField("series")
	if !s.IsNull() {
		p := "1.0"
		if s.String() != "" {
			if pos := f.GetField("position").String(); pos != "" {
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

func (f *Fields) SetField(name string, field *Field) *Fields {
	old := f.GetField(name)
	field.Meta = old.Meta
	field.Data = old.Data
	field.data = old.data
	field.JsonLabel = old.JsonLabel
	field.IsNames = old.IsNames
	field.IsMultiple = old.IsMultiple
	field.idx = old.idx
	f.meta[f.GetFieldIndex(name)] = field
	return f
}

func (f *Fields) SetFieldsFromDB(data []byte) *Fields {
	meta := make(map[string]*Field)
	err := json.Unmarshal(data, &meta)
	if err != nil {
		log.Fatal(err)
	}

	for name, field := range meta {
		f.SetField(name, field)
	}

	return f
}

func (f *Fields) GetField(name string) *Field {
	idx := f.GetFieldIndex(name)
	return f.meta[idx]
}

func (f *Fields) SetMeta(name string, meta Meta) *Fields {
	f.GetField(name).Meta = meta
	return f
}

func (f *Fields) GetFieldIndex(name string) int {
	return f.idx[name]
}

func (f *Fields) EachField() []*Field {
	return f.meta
}

func (f *Fields) ListFields() []string {
	var fields []string
	for _, field := range f.EachField() {
		fields = append(fields, field.JsonLabel)
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
	field := NewField(label)
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

func (f *Field) SetCalibreLabel(label string) *Field {
	f.CalibreLabel = label
	return f
}

func (f *Field) setJsonData(data []byte) *Field {
	f.jsonData = data
	return f
}

func (f *Field) SetData(data any) *Field {
	f.data = data
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

func (f *Field) SetKind(kind string) *Field {
	switch kind {
	case "collection":
		f.Meta = NewMetaCollection()
	case "item":
		f.Meta = NewMetaItem()
	case "column":
		f.Meta = NewMetaColumn()
	}
	return f
}

func (f *Field) GetIndex() int {
	return f.idx
}

func (f *Field) SetMeta(m Meta) *Field {
	f.Meta = m
	return f
}

func (f *Field) ParseData() *Field {
	f.Meta.ParseData(f)
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

const (
	authors        = iota
	authorSort     = iota
	description    = iota
	cover          = iota
	formats        = iota
	id             = iota
	identifiers    = iota
	languages      = iota
	library        = iota
	modified       = iota
	path           = iota
	published      = iota
	publisher      = iota
	rating         = iota
	series         = iota
	position       = iota
	sortAs         = iota
	tags           = iota
	added          = iota
	title          = iota
	titleAndSeries = iota
	uri            = iota
	uuid           = iota
	custCols       = iota
)

func GetDefaultField(name string) *Field {
	return defaultFields()[defaultFieldsIdx()[name]]
}

func ListDefaultFields() []string {
	var fields []string
	for _, field := range defaultFields() {
		fields = append(fields, field.JsonLabel)
	}
	return fields
}

func defaultFieldsIdx() map[string]int {
	return map[string]int{
		"authors":        authors,
		"authorSort":     authorSort,
		"description":    description,
		"cover":          cover,
		"formats":        formats,
		"id":             id,
		"identifiers":    identifiers,
		"languages":      languages,
		"library":        library,
		"modified":       modified,
		"path":           path,
		"published":      published,
		"publisher":      publisher,
		"rating":         rating,
		"series":         series,
		"position":       position,
		"sortAs":         sortAs,
		"tags":           tags,
		"added":          added,
		"title":          title,
		"titleAndSeries": titleAndSeries,
		"uri":            uri,
		"uuid":           uuid,
		"customColumns":  custCols,
	}
}

func defaultFields() []*Field {
	return []*Field{
		NewCollection("authors").
			SetIndex(authors).
			SetIsCategory().
			SetIsMultiple().
			SetIsNames(),
		NewColumn("authorSort").
			SetCalibreLabel("author_sort").
			SetIndex(authorSort),
		NewColumn("description").
			SetCalibreLabel("comments").
			SetIndex(description),
		NewItem("cover").SetIndex(cover),
		NewCollection("formats").
			SetIndex(formats).
			SetIsCategory().
			SetIsMultiple(),
		NewColumn("id").SetIndex(id),
		NewCollection("identifiers").
			SetIndex(identifiers).
			SetIsCategory().
			SetIsMultiple(),
		NewCollection("languages").
			SetIndex(languages).
			SetIsCategory().
			SetIsMultiple(),
		NewColumn("library").SetIndex(library),
		NewColumn("modified").
			SetCalibreLabel("last_modified").
			SetIndex(modified),
		NewColumn("path").SetIndex(path),
		NewColumn("published").
			SetCalibreLabel("pubdate").
			SetIndex(published),
		NewItem("publisher").SetIndex(publisher),
		NewColumn("rating").
			SetIndex(rating).
			SetIsCategory(),
		NewItem("series").
			SetIndex(series).
			SetIsCategory(),
		NewColumn("position").
			SetCalibreLabel("series_index").
			SetIndex(position),
		NewColumn("sortAs").
			SetCalibreLabel("sort").
			SetIndex(sortAs),
		NewCollection("tags").
			SetIndex(tags).
			SetIsCategory().
			SetIsMultiple(),
		NewColumn("added").
			SetCalibreLabel("timestamp").
			SetIndex(added),
		NewColumn("title").SetIndex(title),
		NewColumn("titleAndSeries").SetIndex(titleAndSeries),
		NewColumn("uri").SetIndex(uri),
		NewColumn("uuid").SetIndex(uuid),
		NewField("customColumns").
			SetCalibreLabel("custom_columns").
			SetIndex(custCols),
	}
}
