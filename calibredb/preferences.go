package calibredb

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"golang.org/x/exp/slices"
)

var _ = fmt.Sprintf("%v", "poot")

type Preferences struct {
	HiddenCategories []string          `json:"hiddenCategories"`
	FieldMeta        Fields            `json:"fieldMetadata"`
	SavedSearches    map[string]string `json:"savedSearches"`
}

type calibrePref struct {
	library          string
	HiddenCategories json.RawMessage `json:"tag_browser_hidden_categories"`
	SavedSearches    json.RawMessage `json:"saved_searches"`
	FieldMeta        json.RawMessage `json:"field_metadata"`
	meta             Fields
	DisplayFields    json.RawMessage `json:"book_display_fields"`
	AllFields        []string
}

const prefSql = `
SELECT
JSON_GROUP_OBJECT(
key, JSON(val)
) as pref
FROM preferences 
WHERE key 
IN (
	'saved_searches', 
	'field_metadata', 
	'book_display_fields', 
	'tag_browser_hidden_categories'
)
`

func (lib *Lib) getPreferences() {
	row := lib.db.QueryRowx(prefSql)
	var dbPref []byte
	row.Scan(&dbPref)

	var pref calibrePref
	err := json.Unmarshal(dbPref, &pref)
	if err != nil {
		log.Fatal(err)
	}
	pref.library = lib.Name

	meta := pref.parseFieldMeta()
	lib.Preferences = &Preferences{
		HiddenCategories: pref.parseHiddenCategories(),
		SavedSearches:    pref.parseSavedSearches(),
		FieldMeta:        meta,
	}
}

func (p Preferences) toJSON() []byte {
	json, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	return json
}

func (p calibrePref) parseHiddenCategories() []string {
	var hidden []string
	err := json.Unmarshal(p.HiddenCategories, &hidden)
	if err != nil {
		log.Fatal(err)
	}
	return hidden
}

func (p calibrePref) parseSavedSearches() map[string]string {
	var searches map[string]string
	err := json.Unmarshal(p.SavedSearches, &searches)
	if err != nil {
		log.Fatal(err)
	}
	return searches
}

type Fields map[string]*Field

type Field struct {
	Library      string
	CategorySort string `json:"category_sort"`
	Column       string `json:"column"`
	IsDisplayed  bool   `json:"is_displayed"`
	IsNames      bool   `json:"is_names"`
	IsCategory   bool   `json:"is_category"`
	IsCustom     bool   `json:"is_custom"`
	IsEditable   bool   `json:"is_editable"`
	HasJoin      bool
	IsMultiple   bool
	TableColumns []string
	LinkColumn   string            `json:"link_column"`
	Multiple     map[string]string `json:"is_multiple"`
	Label        string            `json:"label"`
	Name         string            `json:"name"`
	Table        string            `json:"table"`
	CatID        string
}

func (f Field) Type() string {
	switch {
	case f.IsCategory && !f.IsMultiple:
		return "item"
	case f.IsCategory && f.IsMultiple:
		return "category"
	default:
		return "column"
	}
}

func GetFields(f string) *Field {
	switch f {
	case "added":
		return &Field{
			//Library:    lib.Name,
			Label:      "timestamp",
			IsEditable: true,
		}
	case "authorSort":
		return &Field{
			//Library:    lib.Name,
			IsEditable: true,
			Label:      "author_sort",
		}
	case "authors":
		return &Field{
			//Library:    lib.Name,
			Label:      "authors",
			IsCategory: true,
			IsEditable: true,
			IsMultiple: true,
			IsNames:    true,
			LinkColumn: "author",
			Table:      "authors",
		}
	case "description":
		return &Field{
			//Library:    lib.Name,
			Label:      "comments",
			IsEditable: true,
		}
	case "cover":
		return &Field{
			//Library:    lib.Name,
			Label:      "cover",
			IsEditable: true,
		}
	case "formats":
		return &Field{
			//Library:    lib.Name,
			Label:      "formats",
			IsCategory: true,
			IsEditable: true,
			IsMultiple: true,
			Table:      "data",
		}
	case "id":
		return &Field{
			//Library: lib.Name,
			Label: "id",
		}
	case "identifiers":
		return &Field{
			//Library:    lib.Name,
			Label:      "identifiers",
			IsCategory: true,
			IsEditable: true,
			IsMultiple: true,
			Table:      "identifiers",
		}
	case "languages":
		return &Field{
			//Library:    lib.Name,
			Label:      "languages",
			IsCategory: true,
			IsEditable: true,
			IsMultiple: true,
			LinkColumn: "lang_code",
			Table:      "languages",
		}
	case "lastModified":
		return &Field{
			//Library:    lib.Name,
			Label:      "last_modified",
			IsEditable: true,
		}
	case "path":
		return &Field{
			//Library:    lib.Name,
			Label:      "path",
			IsEditable: true,
		}
	case "position":
		return &Field{
			//Library:    lib.Name,
			IsEditable: true,
			Label:      "series_index",
		}
	case "published":
		return &Field{
			//Library:    lib.Name,
			Label:      "pubdate",
			IsEditable: true,
		}
	case "publisher":
		return &Field{
			//Library:    lib.Name,
			Label:      "publisher",
			IsCategory: true,
			IsEditable: true,
			LinkColumn: "publisher",
			Table:      "publishers",
		}
	case "rating":
		return &Field{
			//Library:    lib.Name,
			Label:      "rating",
			IsCategory: true,
			IsEditable: true,
			IsMultiple: true,
			LinkColumn: "rating",
			Table:      "ratings",
		}
	case "series":
		return &Field{
			//Library:    lib.Name,
			Label:      "series",
			IsCategory: true,
			IsEditable: true,
			LinkColumn: "series",
			Table:      "series",
		}
	case "seriesSort":
		return &Field{
			//Library:    lib.Name,
			Label:      "series_sort",
			IsEditable: true,
		}
	case "sort":
		return &Field{
			//Library:    lib.Name,
			Label:      "cover",
			IsEditable: true,
		}
	case "tags":
		return &Field{
			//Library:    lib.Name,
			Label:      "tags",
			IsCategory: true,
			IsEditable: true,
			IsMultiple: true,
			LinkColumn: "tag",
			Table:      "tags",
		}
	case "title":
		return &Field{
			//Library:    lib.Name,
			Label:      "title",
			IsEditable: true,
		}
	case "uuid":
		return &Field{
			//Library:    lib.Name,
			Label:      "uuid",
			IsEditable: true,
		}
	}
	return nil
}

func (p *calibrePref) parseFieldMeta() Fields {
	var fields Fields
	err := json.Unmarshal(p.FieldMeta, &fields)
	if err != nil {
		log.Fatal(err)
	}
	delete(fields, "au_map")
	delete(fields, "size")
	delete(fields, "marked")
	delete(fields, "news")
	delete(fields, "ondevice")
	delete(fields, "search")
	delete(fields, "series_sort")

	for key, _ := range fields {
		p.AllFields = append(p.AllFields, key)
	}

	var dFields [][]interface{}
	err = json.Unmarshal(p.DisplayFields, &dFields)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dFields {
		name := f[0].(string)
		if slices.Contains(p.AllFields, name) {
			fields[name].IsDisplayed = f[1].(bool)
		}
	}

	var fmeta = make(Fields)
	for name, meta := range fields {
		meta.Library = p.library

		if strings.Contains(name, "#") {
			name = strings.Replace(name, "#", "", 1)
		}

		if len(meta.Multiple) != 0 {
			meta.IsMultiple = true
			if del := meta.Multiple["ui_to_list"]; del == "&" {
				meta.IsNames = true
			}
		}

		switch name {
		case "authors":
			meta.TableColumns = []string{"name"}
		case "languages":
			meta.TableColumns = []string{"lang_code"}
		case "tags":
			meta.TableColumns = []string{"name"}
		case "formats":
			meta.TableColumns = []string{}
			meta.CategorySort = "format"
			meta.Table = "data"
			meta.Column = "format"
		case "identifiers":
			meta.TableColumns = []string{"type", "val"}
			meta.Column = "val"
			meta.CategorySort = "type"
			meta.Table = "identifiers"
		case "comments":
			meta.Table = "comments"
		case "publisher":
			meta.TableColumns = []string{"name"}
		case "rating":
			meta.TableColumns = []string{"rating"}
		case "series":
			meta.TableColumns = []string{"name"}
		case "cover":
		}

		fmeta[getJsonField(name)] = meta
	}

	fmeta["uri"] = &Field{
		Column: "uri",
		Name:   "uri",
		Label:  "uri",
	}

	return fmeta
}
