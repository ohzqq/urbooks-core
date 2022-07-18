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
		//"duration",
		//"narrators",
	}
}

type Fields struct {
	displayFields [][]interface{}
	lib           string
	DbMeta        map[string]*Field
	//calibre       []string
	//json          []string
	meta []*Field
	idx  map[string]int
}

type Field struct {
	idx          int
	Library      string            `json:"-"`
	TableColumns map[string]string `json:"-"`
	IsDisplayed  bool              `json:"-"`
	IsNames      bool              `json:"-"`
	HasJoin      bool              `json:"-"`
	IsMultiple   bool              `json:"-"`
	CatID        string            `json:"-"`
	CategorySort string            `json:"category_sort"`
	Colnum       int               `json:"colnum"`
	Column       string            `json:"column"`
	Datatype     string            `json:"datatype"`
	Display      Display           `json:"display"`
	IsCategory   bool              `json:"is_category"`
	IsCustom     bool              `json:"is_custom"`
	IsCsp        bool              `json:"is_csp"`
	IsEditable   bool              `json:"is_editable"`
	Multiple     Multiple          `json:"is_multiple"`
	Kind         string            `json:"kind"`
	Label        string            `json:"label"`
	LinkColumn   string            `json:"link_column"`
	Name         string            `json:"name"`
	RecIndex     int               `json:"rec_index"`
	SearchTerms  []string          `json:"search_terms"`
	Table        string            `json:"table"`
	IsJson       bool              `json:"-"`
	JsonLabel    string            `json:"-"`
	Value        string            `json:"-"`
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
		//calibre: defaultCalibreFields(),
		//json:    defaultJsonFields(),
		meta: defaultFields(),
		idx:  defaultFieldsIdx(),
	}
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
		if ff := f.GetField(name); ff.Label == name {
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
	f.idx[field.Label] = len(f.meta)
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
			idx:        authors,
			Label:      "authors",
			JsonLabel:  "authors",
			IsCategory: true,
			IsEditable: true,
			IsMultiple: true,
			IsNames:    true,
			LinkColumn: "author",
			Multiple: Multiple{
				CacheToList: ",",
				ListToUi:    " & ",
				UiToList:    "&",
			},
			Table: "authors",
			TableColumns: map[string]string{
				"value": "name",
				"uri":   `"authors/"` + " || id",
			},
		},
		&Field{
			idx:        authorSort,
			Label:      "author_sort",
			JsonLabel:  "authorSort",
			IsEditable: true,
		},
		&Field{
			idx:        description,
			Label:      "comments",
			JsonLabel:  "description",
			Table:      "comments",
			IsEditable: true,
		},
		&Field{
			idx:        cover,
			Label:      "cover",
			JsonLabel:  "cover",
			IsEditable: true,
		},
		&Field{
			idx:          formats,
			Label:        "formats",
			JsonLabel:    "formats",
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
			Table:        "data",
			CategorySort: "format",
			Column:       "format",
			Multiple: Multiple{
				CacheToList: ",",
				ListToUi:    ", ",
				UiToList:    ",",
			},
			TableColumns: map[string]string{
				"basename":  "name",
				"extension": "lower(format)",
				"value":     `name || '.' || lower(format)`,
				"size":      "lower(uncompressed_size)",
				"uri":       `"books/" || books.id`,
				//"path":      `"` + p.library + `" || "/" || path || "/" || name || '.' || lower(format)`,
			},
		},
		&Field{
			idx:       id,
			Label:     "id",
			JsonLabel: "id",
		},
		&Field{
			idx:          identifiers,
			Label:        "identifiers",
			JsonLabel:    "identifiers",
			IsCategory:   true,
			IsEditable:   true,
			IsMultiple:   true,
			Column:       "val",
			CategorySort: "type",
			Table:        "identifiers",
			TableColumns: map[string]string{
				"value": "val",
				"type":  "type",
			},
		},
		&Field{
			idx:        languages,
			Label:      "languages",
			JsonLabel:  "languages",
			IsCategory: true,
			IsEditable: true,
			IsMultiple: true,
			LinkColumn: "lang_code",
			Table:      "languages",
			Multiple: Multiple{
				CacheToList: ",",
				ListToUi:    ", ",
				UiToList:    ",",
			},
			TableColumns: map[string]string{
				"value": "lang_code",
				"uri":   `"languages/"` + " || id",
			},
		},
		&Field{
			idx:       library,
			Label:     "library",
			JsonLabel: "library",
		},
		&Field{
			idx:        modified,
			Label:      "last_modified",
			JsonLabel:  "modified",
			IsEditable: true,
		},
		&Field{
			idx:        path,
			Label:      "path",
			JsonLabel:  "path",
			IsEditable: true,
		},
		&Field{
			idx:        published,
			Label:      "pubdate",
			JsonLabel:  "published",
			IsEditable: true,
		},
		&Field{
			idx:        publisher,
			Label:      "publisher",
			JsonLabel:  "publisher",
			IsCategory: true,
			IsEditable: true,
			LinkColumn: "publisher",
			Table:      "publishers",
			TableColumns: map[string]string{
				"value": "name",
				"uri":   `"publisher/"` + " || id",
			},
		},
		&Field{
			idx:        rating,
			Label:      "rating",
			JsonLabel:  "rating",
			IsCategory: true,
			IsEditable: true,
			IsMultiple: true,
			LinkColumn: "rating",
			Table:      "ratings",
			TableColumns: map[string]string{
				"value": "rating",
			},
		},
		&Field{
			idx:        series,
			Label:      "series",
			JsonLabel:  "series",
			IsCategory: true,
			IsEditable: true,
			LinkColumn: "series",
			Table:      "series",
			TableColumns: map[string]string{
				"value":    "name",
				"position": "lower(series_index)",
				"uri":      `"series/"` + " || id",
			},
		},
		&Field{
			idx:        position,
			Label:      "series_index",
			JsonLabel:  "position",
			IsEditable: true,
		},
		&Field{
			idx:        sortAs,
			Label:      "sort",
			JsonLabel:  "sortAs",
			IsEditable: true,
		},
		&Field{
			idx:        tags,
			Label:      "tags",
			JsonLabel:  "tags",
			IsCategory: true,
			IsEditable: true,
			IsMultiple: true,
			LinkColumn: "tag",
			Table:      "tags",
			Multiple: Multiple{
				CacheToList: ",",
				ListToUi:    ", ",
				UiToList:    ",",
			},
			TableColumns: map[string]string{
				"value": "name",
				"uri":   `"tags/"` + " || id",
			},
		},
		&Field{
			idx:        added,
			Label:      "timestamp",
			JsonLabel:  "added",
			IsEditable: true,
		},
		&Field{
			idx:        title,
			Label:      "title",
			JsonLabel:  "title",
			IsEditable: true,
		},
		&Field{
			idx:       titleAndSeries,
			Label:     "titleAndSeries",
			JsonLabel: "titleAndSeries",
			Name:      "titleAndSeries",
		},
		&Field{
			idx:        uri,
			Label:      "uri",
			JsonLabel:  "uri",
			Column:     "uri",
			Name:       "uri",
			IsEditable: true,
		},
		&Field{
			idx:        uuid,
			Label:      "uuid",
			JsonLabel:  "uuid",
			IsEditable: true,
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
