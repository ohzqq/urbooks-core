package book

import (
	"encoding/json"
	"log"
)

const (
	sort     = iota
	order    = iota
	limit    = iota
	category = iota
)

func defaultCalibreFields() []string {
	return []string{
		"authors",
		"author_sort",
		"comments",
		"cover",
		"formats",
		"id",
		"identifiers",
		"languages",
		"library",
		"last_modified",
		"path",
		"pubdate",
		"publisher",
		"rating",
		"series",
		"series_index",
		"sort",
		"tags",
		"timestamp",
		"title",
		"titleAndSeries",
		"uri",
		"uuid",
		"customColumns",
		//"duration",
		//"narrators",
	}
}

type Fields struct {
	displayFields [][]interface{}
	lib           string
	DbMeta        map[string]*Field
	meta          []*Field
	idx           map[string]int
}

type Field struct {
	idx          int
	Library      string          `json:"-"`
	IsDisplayed  bool            `json:"-"`
	IsNames      bool            `json:"-"`
	IsMultiple   bool            `json:"-"`
	CatID        string          `json:"-"`
	IsJson       bool            `json:"-"`
	JsonLabel    string          `json:"-"`
	Data         json.RawMessage `json:"-"`
	Value        string          `json:"-"`
	Meta         Meta            `json:"-"`
	IsColumn     bool            `json:"-"`
	IsCollection bool            `json:"-"`
	IsItem       bool            `json:"-"`
	CategorySort string          `json:"category_sort"`
	Colnum       int             `json:"colnum"`
	Column       string          `json:"column"`
	Datatype     string          `json:"datatype"`
	Display      Display         `json:"display"`
	IsCategory   bool            `json:"is_category"`
	IsCustom     bool            `json:"is_custom"`
	IsCsp        bool            `json:"is_csp"`
	IsEditable   bool            `json:"is_editable"`
	Multiple     Multiple        `json:"is_multiple"`
	Kind         string          `json:"kind"`
	CalibreLabel string          `json:"label"`
	LinkColumn   string          `json:"link_column"`
	Name         string          `json:"name"`
	RecIndex     int             `json:"rec_index"`
	SearchTerms  []string        `json:"search_terms"`
	Table        string          `json:"table"`
	//Value        any               `json:"#value#"`
	//Extra        any               `json:"#extra#"`
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

func NewField() *Field {
	return &Field{}
}

func (f *Fields) Each() []*Field {
	return defaultFields()
	//return f.meta
}

//func (f *Fields) AllJson() []string {
//  return f.json
//}

//func (f *Fields) AllCalibre() []string {
//  return f.calibre
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
		//println(name)
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

func (f *Fields) AddField(field *Field) *Fields {
	f.idx[field.CalibreLabel] = len(f.meta)
	f.meta = append(f.meta, field)
	//f.calibre = append(f.calibre, field.Label)
	//f.json = append(f.json, field.Label)
	return f
}

func (f *Fields) GetField(name string) *Field {
	idx := f.GetFieldIndex(name)
	return f.meta[idx]
}

func (f *Fields) GetFieldIndex(name string) int {
	return f.idx[name]
}

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
	//duration       = iota
	//narrators      = iota
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
		//"duration":       duration,
		//"narrators":      narrators,
	}
}

func defaultJsonFields() []string {
	return []string{
		"authors",
		"authorSort",
		"description",
		"cover",
		"formats",
		"id",
		"identifiers",
		"languages",
		"library",
		"modified",
		"path",
		"published",
		"publisher",
		"rating",
		"series",
		"position",
		"sortAs",
		"tags",
		"added",
		"title",
		"titleAndSeries",
		"uri",
		"uuid",
		//"duration",
		//"narrators",
	}
}

func defaultFields() []*Field {
	return []*Field{
		&Field{
			idx:          authors,
			CalibreLabel: "authors",
			JsonLabel:    "authors",
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
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          description,
			CalibreLabel: "comments",
			JsonLabel:    "description",
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          cover,
			CalibreLabel: "cover",
			JsonLabel:    "cover",
			IsItem:       true,
			IsEditable:   true,
		},
		&Field{
			idx:          formats,
			CalibreLabel: "formats",
			JsonLabel:    "formats",
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
			IsColumn:     true,
		},
		&Field{
			idx:          identifiers,
			CalibreLabel: "identifiers",
			JsonLabel:    "identifiers",
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
			IsCollection: true,
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
		},
		&Field{
			idx:          library,
			CalibreLabel: "library",
			JsonLabel:    "library",
			IsColumn:     true,
		},
		&Field{
			idx:          modified,
			CalibreLabel: "last_modified",
			JsonLabel:    "modified",
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          path,
			CalibreLabel: "path",
			JsonLabel:    "path",
			IsColumn:     true,
		},
		&Field{
			idx:          published,
			CalibreLabel: "pubdate",
			JsonLabel:    "published",
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          publisher,
			CalibreLabel: "publisher",
			JsonLabel:    "publisher",
			IsItem:       true,
			IsCategory:   true,
			IsEditable:   true,
		},
		&Field{
			idx:          rating,
			CalibreLabel: "rating",
			JsonLabel:    "rating",
			IsColumn:     true,
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
		},
		&Field{
			idx:          series,
			CalibreLabel: "series",
			JsonLabel:    "series",
			IsItem:       true,
			IsCategory:   true,
			IsEditable:   true,
		},
		&Field{
			idx:          position,
			CalibreLabel: "series_index",
			JsonLabel:    "position",
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          sortAs,
			CalibreLabel: "sort",
			JsonLabel:    "sortAs",
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          tags,
			CalibreLabel: "tags",
			JsonLabel:    "tags",
			IsCollection: true,
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
		},
		&Field{
			idx:          added,
			CalibreLabel: "timestamp",
			JsonLabel:    "added",
			IsColumn:     true,
		},
		&Field{
			idx:          title,
			CalibreLabel: "title",
			JsonLabel:    "title",
			IsColumn:     true,
			IsEditable:   true,
		},
		&Field{
			idx:          titleAndSeries,
			CalibreLabel: "titleAndSeries",
			JsonLabel:    "titleAndSeries",
			IsColumn:     true,
		},
		&Field{
			idx:          uri,
			CalibreLabel: "uri",
			JsonLabel:    "uri",
			IsColumn:     true,
		},
		&Field{
			idx:          uuid,
			CalibreLabel: "uuid",
			JsonLabel:    "uuid",
			IsColumn:     true,
		},
		&Field{
			idx:          custCols,
			CalibreLabel: "custom_columns",
			JsonLabel:    "customColumns",
		},
		//"seriesSort": &Field{
		//  idx:        seriesSort,
		//  Label:      "series_sort",
		//  IsEditable: true,
		//},
		//"duration": &Field{
		//  idx: duration,
		//},
		//"narrators": &Field{
		//  idx: narrators,
		//},
	}
}
