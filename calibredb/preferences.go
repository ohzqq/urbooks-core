package calibredb

import (
	"encoding/json"
	"log"
	"strings"

	"golang.org/x/exp/slices"
)

type Preferences struct {
	raw              calibrePref
	HiddenCategories []string          `json:"hiddenCategories"`
	FieldMeta        Fields            `json:"fieldMetadata"`
	SavedSearches    map[string]string `json:"savedSearches"`
}

type calibrePref struct {
	dbPreferences
	library   string
	meta      Fields
	AllFields []string
}

type dbPreferences struct {
	HiddenCategories json.RawMessage `json:"tag_browser_hidden_categories"`
	DisplayFields    json.RawMessage `json:"book_display_fields"`
	SavedSearches    json.RawMessage `json:"saved_searches"`
	FieldMeta        json.RawMessage `json:"field_metadata"`
}

var custColstmt = `
SELECT 
label label, 
lower(id) id, 
CASE IFNULL(JSON_EXTRACT(display, "$.is_names"), 0)
WHEN 0 THEN "false"
WHEN 1 THEN "true"
END is_names,
"custom_column_" || id 'table',
CASE is_multiple
WHEN true THEN "true"
ELSE "false"
END is_multiple,
CASE is_multiple
WHEN true THEN "books_custom_column_" || id || "_link"
ELSE ""
END join_table,
CASE is_multiple
WHEN true THEN 'value'
ELSE ""
END link_column
FROM custom_columns;
`

func (lib *Lib) getCustCols() {
	rows, err := lib.db.Queryx(custColstmt)
	if err != nil {
		log.Fatalf("cust col query failed %v\n", err)
	}

	var cols []map[string]interface{}
	for rows.Next() {
		results := make(map[string]interface{})
		err = rows.MapScan(results)
		if err != nil {
			log.Fatalf("something happened when scanning db results %v\n", err)
		}
		cols = append(cols, results)
	}

	for _, val := range cols {
		results := make(map[string]string)
		for k, v := range val {
			results[k] = v.(string)
		}
		results["label"] = "#" + results["label"]
		lib.CustCols = append(lib.CustCols, results)
		lib.Fields.CustomCol = append(lib.Fields.CustomCol, results["label"])
	}
	//fmt.Printf("%+v\n", lib.CustCols)
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

func GetFieldMeta(lib *Lib, f, v string) string {
	return lib.getFieldMeta(f, v)
}

func (lib *Lib) getFieldMeta(f, v string) string {
	if slices.Contains(lib.Fields.CustomCol, "#"+f) {
		f = "#" + f
	}

	if v == "custom_column" {
		if strings.HasPrefix(f, "#") {
			return "true"
		}
	}

	jmeta := lib.fieldMeta[f]

	if v == "join_table" {
		if f == "series" || f == "publisher" || len(jmeta["is_multiple"].(map[string]interface{})) != 0 {
			if col := jmeta["link_column"]; col != nil {
				return "books_" + jmeta["table"].(string) + "_link"
			}
		}
	}

	if val := jmeta[v]; val != nil {
		return val.(string)
	}

	if v == "table" {
		switch f {
		case "comments":
			return "comments"
		case "formats":
			return "data"
		case "identifiers":
			return "identifiers"
		}
	}
	return ""
}

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
		raw:              pref,
	}

	err = json.Unmarshal(pref.FieldMeta, &lib.fieldMeta)
	if err != nil {
		log.Fatalf("getDBfieldMeta json unmarshal failed: %v\n", err)
	}
}

func (p Preferences) toJSON() []byte {
	json, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	return json
}

func (p Preferences) ToJSON() []byte {
	pref := map[string]json.RawMessage{
		"HiddenCategories": p.raw.HiddenCategories,
		"SavedSearches":    p.raw.SavedSearches,
		"DisplayFields":    p.raw.DisplayFields,
		"FieldMeta":        p.raw.FieldMeta,
	}

	json, err := json.Marshal(pref)
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
	Library      string            `json:"-"`
	TableColumns map[string]string `json:"-"`
	//IsDisplayed  bool              `json:"-"`
	IsNames      bool     `json:"-"`
	HasJoin      bool     `json:"-"`
	IsMultiple   bool     `json:"-"`
	CatID        string   `json:"-"`
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
	Label        string   `json:"label"`
	LinkColumn   string   `json:"link_column"`
	Name         string   `json:"name"`
	RecIndex     int      `json:"rec_index"`
	SearchTerms  []string `json:"search_terms"`
	Table        string   `json:"table"`
	Value        any      `json:"#value#"`
	Extra        any      `json:"#extra#"`
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

type dbFieldTypes struct {
	Lib        string
	MultiCats  []string
	SingleCats []string
	JoinCats   []string
	BookCol    []string
	DateCol    []string
	CustomCol  []string
	Special    []string
	urCols     []string
}

var defaultFields = dbFieldTypes{
	MultiCats:  []string{"formats", "identifiers"},
	JoinCats:   []string{"authors", "languages", "tags"},
	SingleCats: []string{"publisher", "series"},
	BookCol:    []string{"author_sort", "path", "sort", "title", "uuid"},
	DateCol:    []string{"last_modified", "pubdate", "timestamp"},
	Special:    []string{"comments", "id", "series_index", "cover", "rating"},
	urCols:     []string{"library", "titleAndSeries", "uri"},
}

func newLibFields(lib string) *dbFieldTypes {
	return &dbFieldTypes{
		Lib:        lib,
		MultiCats:  defaultFields.MultiCats,
		JoinCats:   defaultFields.JoinCats,
		SingleCats: defaultFields.SingleCats,
		BookCol:    defaultFields.BookCol,
		DateCol:    defaultFields.DateCol,
		Special:    defaultFields.Special,
		urCols:     defaultFields.urCols,
	}
}

func GetTableColumns(f, lib string) map[string]string {
	fields := map[string]map[string]string{
		"authors": map[string]string{
			"value": "name",
			"uri":   `"authors/"` + " || id",
		},
		"formats": map[string]string{
			"basename":  "name",
			"extension": "lower(format)",
			"value":     `name || '.' || lower(format)`,
			"size":      "lower(uncompressed_size)",
			"uri":       `"books/" || books.id`,
			"path":      `"` + lib + `" || "/" || books.path || "/" || name || '.' || lower(format)`,
		},
		"identifiers": map[string]string{
			"value": "val",
			"type":  "type",
		},
		"languages": map[string]string{
			"value": "lang_code",
			"uri":   `"languages/"` + " || id",
		},
		"publisher": map[string]string{
			"value": "name",
			"uri":   `"publisher/"` + " || id",
		},
		"rating": map[string]string{
			"value": "lower(rating)",
		},
		"series": map[string]string{
			"value":    "name",
			"position": "lower(series_index)",
			"uri":      `"series/"` + " || id",
		},
		"tags": map[string]string{
			"value": "name",
			"uri":   `"tags/"` + " || id",
		},
	}
	return fields[f]
}

func CalibreFieldList() []string {
	var fields []string
	for _, f := range defaultFields.MultiCats {
		fields = append(fields, f)
	}
	for _, f := range defaultFields.JoinCats {
		fields = append(fields, f)
	}
	for _, f := range defaultFields.SingleCats {
		fields = append(fields, f)
	}
	for _, f := range defaultFields.BookCol {
		fields = append(fields, f)
	}
	for _, f := range defaultFields.DateCol {
		fields = append(fields, f)
	}
	for _, f := range defaultFields.Special {
		fields = append(fields, f)
	}
	return fields
}

func FieldList() []string {
	fields := CalibreFieldList()
	for _, f := range defaultFields.urCols {
		fields = append(fields, f)
	}
	return fields
}

//func FieldType() string {
//}

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

func (f Field) ToJson() []byte {
	js, err := json.Marshal(f)
	if err != nil {
		log.Fatalf("json failed to marshal a calibred.Field, %v", err)
	}
	return js
}

func (f Field) IsCat() bool {
	return f.IsCategory && f.IsMultiple
}

func (f Field) IsItem() bool {
	return f.IsCategory && !f.IsMultiple
}

func (f Field) IsCol() bool {
	return !f.IsCategory && !f.IsMultiple
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

	//var dFields [][]interface{}
	//err = json.Unmarshal(p.DisplayFields, &dFields)
	//if err != nil {
	//  log.Fatal(err)
	//}

	//for _, f := range dFields {
	//  name := f[0].(string)
	//  if slices.Contains(p.AllFields, name) {
	//    fields[name].IsDisplayed = f[1].(bool)
	//  }
	//}

	var fmeta = make(Fields)
	for name, meta := range fields {
		meta.Library = p.library

		if strings.Contains(name, "#") {
			name = strings.Replace(name, "#", "", 1)
		}

		if meta.Multiple != (Multiple{}) {
			meta.IsMultiple = true
			//if del := meta.Multiple.UiToList; del == "&" {
			//meta.IsNames = true
			//}
		}

		switch name {
		case "authors":
			meta.TableColumns = map[string]string{
				"value": "name",
				"uri":   `"` + GetJsonField(name) + `/"` + " || id",
			}
		case "languages":
			meta.TableColumns = map[string]string{
				"value": "lang_code",
				"uri":   `"` + GetJsonField(name) + `/"` + " || id",
			}
		case "tags":
			meta.TableColumns = map[string]string{
				"value": "name",
				"uri":   `"` + GetJsonField(name) + `/"` + " || id",
			}
		case "formats":
			meta.TableColumns = map[string]string{
				"basename":  "name",
				"extension": "lower(format)",
				"value":     `name || '.' || lower(format)`,
				"size":      "lower(uncompressed_size)",
				"uri":       `"books/" || books.id`,
				"path":      `"` + p.library + `" || "/" || books.path || "/" || name || '.' || lower(format)`,
			}
			meta.CategorySort = "format"
			meta.Table = "data"
			meta.Column = "format"
		case "identifiers":
			meta.TableColumns = map[string]string{
				"value": "val",
				"type":  "type",
			}
			meta.Column = "val"
			meta.CategorySort = "type"
			meta.Table = "identifiers"
		case "comments":
			meta.Table = "comments"
		case "publisher":
			meta.TableColumns = map[string]string{
				"value": "name",
				"uri":   `"` + GetJsonField(name) + `/"` + " || id",
			}
		case "rating":
			meta.TableColumns = map[string]string{
				"value": "rating",
			}
		case "series":
			meta.TableColumns = map[string]string{
				"value":    "name",
				"position": "lower(series_index)",
				"uri":      `"` + GetJsonField(name) + `/"` + " || id",
			}
		case "cover":
			meta.IsCustom = false
		case "library":
			meta.IsCustom = false
		}

		fmeta[GetJsonField(name)] = meta
	}

	fmeta["uri"] = &Field{
		Column: "uri",
		Name:   "uri",
		Label:  "uri",
	}

	return fmeta
}
