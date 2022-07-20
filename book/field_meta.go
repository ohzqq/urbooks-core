package book

import (
	"encoding/json"
	"log"
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

func (f *Fields) AddField(field *Field) *Fields {
	f.idx[field.CalibreLabel] = len(f.meta)
	f.meta = append(f.meta, field)
	return f
}

func (f *Fields) GetField(name string) *Field {
	return f.meta[f.GetFieldIndex(name)]
}

func (f *Fields) GetFieldIndex(name string) int {
	return f.idx[name]
}

func (f *Fields) Each() []*Field {
	return defaultFields()
}

func (f *Field) Index() int {
	return f.idx
}

const (
	nameSep    = " & "
	itemSep    = ", "
	cliItemSep = `,`
	cliNameSep = `&`
)

func (f *Field) String() string {
	return f.Meta.String(f)
}

func (f *Field) URL() string {
	return f.Meta.URL(f)
}

func (f *Field) IsNull() bool {
	return f.Meta.IsNull()
}

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
		&Field{
			idx:          authors,
			CalibreLabel: "authors",
			JsonLabel:    "authors",
			Meta:         NewCollection(),
			IsCollection: true,
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
			IsNames:      true,
		},
		&Field{
			idx:          authorSort,
			CalibreLabel: "author_sort",
			JsonLabel:    "authorSort",
			Meta:         NewColumn(),
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          description,
			CalibreLabel: "comments",
			JsonLabel:    "description",
			Meta:         NewColumn(),
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          cover,
			CalibreLabel: "cover",
			JsonLabel:    "cover",
			Meta:         NewItem(),
			IsItem:       true,
			IsEditable:   true,
		},
		&Field{
			idx:          formats,
			CalibreLabel: "formats",
			JsonLabel:    "formats",
			Meta:         NewCollection(),
			IsCollection: true,
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
			CategorySort: "format",
		},
		&Field{
			idx:          id,
			CalibreLabel: "id",
			JsonLabel:    "id",
			Meta:         NewColumn(),
			IsColumn:     true,
		},
		&Field{
			idx:          identifiers,
			CalibreLabel: "identifiers",
			JsonLabel:    "identifiers",
			Meta:         NewCollection(),
			IsCollection: true,
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
			CategorySort: "type",
		},
		&Field{
			idx:          languages,
			CalibreLabel: "languages",
			JsonLabel:    "languages",
			Meta:         NewCollection(),
			IsCollection: true,
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
		},
		&Field{
			idx:          library,
			CalibreLabel: "library",
			JsonLabel:    "library",
			Meta:         NewColumn(),
			IsColumn:     true,
		},
		&Field{
			idx:          modified,
			CalibreLabel: "last_modified",
			JsonLabel:    "modified",
			Meta:         NewColumn(),
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          path,
			CalibreLabel: "path",
			JsonLabel:    "path",
			Meta:         NewColumn(),
			IsColumn:     true,
		},
		&Field{
			idx:          published,
			CalibreLabel: "pubdate",
			JsonLabel:    "published",
			Meta:         NewColumn(),
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          publisher,
			CalibreLabel: "publisher",
			JsonLabel:    "publisher",
			Meta:         NewItem(),
			IsItem:       true,
			IsCategory:   true,
			IsEditable:   true,
		},
		&Field{
			idx:          rating,
			CalibreLabel: "rating",
			JsonLabel:    "rating",
			Meta:         NewColumn(),
			IsColumn:     true,
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
		},
		&Field{
			idx:          series,
			CalibreLabel: "series",
			JsonLabel:    "series",
			Meta:         NewItem(),
			IsItem:       true,
			IsCategory:   true,
			IsEditable:   true,
		},
		&Field{
			idx:          position,
			CalibreLabel: "series_index",
			JsonLabel:    "position",
			Meta:         NewColumn(),
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          sortAs,
			CalibreLabel: "sort",
			JsonLabel:    "sortAs",
			Meta:         NewColumn(),
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          tags,
			CalibreLabel: "tags",
			JsonLabel:    "tags",
			Meta:         NewCollection(),
			IsCollection: true,
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
		},
		&Field{
			idx:          added,
			CalibreLabel: "timestamp",
			JsonLabel:    "added",
			Meta:         NewColumn(),
			IsColumn:     true,
		},
		&Field{
			idx:          title,
			CalibreLabel: "title",
			JsonLabel:    "title",
			Meta:         NewColumn(),
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          titleAndSeries,
			CalibreLabel: "titleAndSeries",
			JsonLabel:    "titleAndSeries",
			Meta:         NewColumn(),
			IsColumn:     true,
		},
		&Field{
			idx:          uri,
			CalibreLabel: "uri",
			JsonLabel:    "uri",
			Meta:         NewColumn(),
			IsColumn:     true,
		},
		&Field{
			idx:          uuid,
			CalibreLabel: "uuid",
			JsonLabel:    "uuid",
			Meta:         NewColumn(),
			IsColumn:     true,
		},
		&Field{
			idx:          custCols,
			CalibreLabel: "custom_columns",
			JsonLabel:    "customColumns",
		},
	}
}
