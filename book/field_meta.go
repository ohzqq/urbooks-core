package book

import (
	"encoding/json"
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
	Library      string   `json:"-"`
	IsDisplayed  bool     `json:"-"`
	IsNames      bool     `json:"-"`
	IsMultiple   bool     `json:"-"`
	CatID        string   `json:"-"`
	IsJson       bool     `json:"-"`
	JsonLabel    string   `json:"-"`
	jsonData     []byte   `json:"-"`
	stringData   string   `json:"-"`
	data         any      `json:"-"`
	Data         []byte   `json:"-"`
	Value        string   `json:"-"`
	Meta         Meta     `json:"-"`
	IsColumn     bool     `json:"-"`
	IsCollection bool     `json:"-"`
	IsItem       bool     `json:"-"`
	CategorySort string   `json:"category_sort"`
	Colnum       int      `json:"colnum"`
	Column       string   `json:"column"`
	Datatype     string   `json:"datatype"`
	Display      Display  `json:"display"`
	IsCategory   bool     `json:"is_category"`
	IsCustom     bool     `json:"is_custom"`
	IsCsp        bool     `json:"is_csp"`
	IsEditable   bool     `json:"is_editable"`
	Multiple     Multiple `json:"is_multiple"`
	Kind         string   `json:"kind"`
	CalibreLabel string   `json:"label"`
	LinkColumn   string   `json:"link_column"`
	Name         string   `json:"name"`
	RecIndex     int      `json:"rec_index"`
	SearchTerms  []string `json:"search_terms"`
	Table        string   `json:"table"`
}

type Display struct {
	Description     string `json:"description"`
	HeadingPosition string `json:"heading_position"`
	InterpretAs     string `json:"long-text"`
	IsNames         bool   `json:"is_names"`
}

type Multiple struct {
	CacheToList string `json:"cache_to_list"`
	ListToUi    string `json:"list_to_ui"`
	UiToList    string `json:"ui_to_list"`
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
	field.IsCollection = true
	field.Meta = NewMetaCollection()
	return field
}

func NewItem(label string) *Field {
	field := NewField(label)
	field.IsItem = true
	field.Meta = NewMetaItem()
	return field
}

func NewColumn(label string) *Field {
	field := NewField(label)
	field.IsColumn = true
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

func (f *Field) SetIsEditable() *Field {
	f.IsEditable = true
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
		f.IsCollection = true
		f.Meta = NewMetaCollection()
	case "item":
		f.IsItem = true
		f.Meta = NewMetaItem()
	case "column":
		f.IsColumn = true
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
	case f.IsCollection:
		meta = f.Collection().Split(data, f.IsNames)
	case f.IsItem:
		meta = f.Item().Set("value", data)
	case f.IsColumn:
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

func (f *Field) Item() *Item {
	return f.Meta.(*Item)
}

func (f *Field) Collection() *Collection {
	return f.Meta.(*Collection)
}

func (f *Field) Col() *Column {
	return f.Meta.(*Column)
}

//func (f *Field) UnmarshalJSON(d []byte) error {
//  err := f.Meta.UnmarshalJSON(d)
//  if err != nil {
//    return err
//  }
//  return nil
//}

func (f *Fields) ParseDBFieldMeta(meta, display json.RawMessage) {
	err := json.Unmarshal(display, &f.displayFields)
	if err != nil {
		log.Fatal(err)
	}

	for _, field := range f.displayFields {
		name := field[0].(string)
		if ff := f.GetField(name); ff.CalibreLabel == name {
			ff.IsDisplayed = field[1].(bool)
		}
	}

	err = json.Unmarshal(meta, &f.DbMeta)
	if err != nil {
		log.Fatal(err)
	}

	for name, meta := range f.DbMeta {
		if meta.IsCustom {
			meta.JsonLabel = name
			f.AddField(meta)

			if meta.Multiple != (Multiple{}) {
				meta.IsMultiple = true
				if del := meta.Multiple.UiToList; del == "&" {
					meta.IsNames = true
				}
			}
		}
		f.DbMeta[name] = meta
	}
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
			SetIsEditable().
			SetIsNames(),
		NewColumn("authorSort").
			SetCalibreLabel("author_sort").
			SetIndex(authorSort).
			SetIsEditable(),
		NewColumn("description").
			SetCalibreLabel("comments").
			SetIndex(description).
			SetIsEditable(),
		NewItem("cover").
			SetIndex(cover).
			SetIsEditable(),
		NewCollection("formats").
			SetIndex(formats).
			SetIsCategory().
			SetIsMultiple().
			SetIsEditable(),
		NewColumn("id").SetIndex(id),
		NewCollection("identifiers").
			SetIndex(identifiers).
			SetIsCategory().
			SetIsMultiple().
			SetIsEditable(),
		NewCollection("languages").
			SetIndex(languages).
			SetIsCategory().
			SetIsMultiple().
			SetIsEditable(),
		NewColumn("library").SetIndex(library),
		NewColumn("modified").
			SetCalibreLabel("last_modified").
			SetIndex(modified),
		NewColumn("path").SetIndex(path),
		NewColumn("published").
			SetCalibreLabel("pubdate").
			SetIndex(published).
			SetIsEditable(),
		NewItem("publisher").
			SetIndex(publisher).
			SetIsEditable().
			SetIsCategory(),
		NewColumn("rating").
			SetIndex(rating).
			SetIsCategory().
			SetIsEditable(),
		NewItem("series").
			SetIndex(series).
			SetIsCategory().
			SetIsEditable(),
		NewColumn("position").
			SetCalibreLabel("series_index").
			SetIndex(position).
			SetIsEditable(),
		NewColumn("sortAs").
			SetCalibreLabel("sort").
			SetIndex(sortAs).
			SetIsEditable(),
		NewCollection("tags").
			SetIndex(tags).
			SetIsCategory().
			SetIsMultiple().
			SetIsEditable(),
		NewColumn("added").
			SetCalibreLabel("timestamp").
			SetIndex(added),
		NewColumn("title").
			SetIndex(title).
			SetIsEditable(),
		NewColumn("titleAndSeries").
			SetIndex(titleAndSeries).
			SetIsEditable(),
		NewColumn("uri").SetIndex(uri),
		NewColumn("uuid").SetIndex(uuid),
		NewField("customColumns").
			SetCalibreLabel("custom_columns").
			SetIndex(custCols),
	}
}
