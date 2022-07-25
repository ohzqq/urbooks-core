package calibredb

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"text/template"

	"golang.org/x/exp/slices"
)

var custColstmt = `
SELECT 
label label, 
(SELECT 
JSON_EXTRACT(val, "$." || "#" || label)
FROM preferences 
WHERE key = 'field_metadata'
) meta,
(SELECT 
CASE IFNULL(JSON_EXTRACT(val, "$." || "#" || label || ".is_category"), 0)
WHEN 0 THEN "false"
WHEN 1 then "true"
END
FROM preferences 
WHERE key = 'field_metadata'
) is_category,
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

const fieldMetaSql = `
SELECT
JSON_OBJECT(
{{range .}}
"{{.}}", JSON_EXTRACT(val, "$." || {{.}}) 
{{end}}
) fieldMeta
FROM preferences 
WHERE key = "field_metadata"
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

type Preferences struct {
	raw              calibrePref
	HiddenCategories []string          `json:"hiddenCategories"`
	SavedSearches    map[string]string `json:"savedSearches"`
}

func (lib *Lib) GetPref(p string) json.RawMessage {
	var stmt string
	switch p {
	case "field_meta":
		stmt = lib.renderFieldMetaTmpl()
	}

	row := lib.db.QueryRowx(stmt)
	var dbPref []byte
	row.Scan(&dbPref)

	return json.RawMessage(dbPref)
}

func (lib *Lib) renderFieldMetaTmpl() string {
	var buf bytes.Buffer
	err := lib.bookTmpl.ExecuteTemplate(&buf, "rangeFieldMeta", lib.Request.Fields)
	if err != nil {
		log.Println("executing template:", err)
	}

	return buf.String()
}

func (lib *Lib) GetPreferences() json.RawMessage {
	stmt := lib.renderSqlTmpl("Prefs")
	row := lib.db.QueryRowx(stmt)
	var dbPref []byte
	row.Scan(&dbPref)

	return json.RawMessage(dbPref)
}

type calibrePref struct {
	dbPreferences
	library   string
	AllFields []string
}

type dbPreferences struct {
	HiddenCategories json.RawMessage `json:"tag_browser_hidden_categories"`
	DisplayFields    json.RawMessage `json:"book_display_fields"`
	SavedSearches    json.RawMessage `json:"saved_searches"`
	FieldMeta        json.RawMessage `json:"field_metadata"`
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

	lib.Preferences = &Preferences{
		raw: pref,
	}

	err = json.Unmarshal(pref.FieldMeta, &lib.fieldMeta)
	if err != nil {
		log.Fatalf("getDBfieldMeta json unmarshal failed: %v\n", err)
	}
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

func (lib *Lib) AllFields() []string {
	fields := DefaultCalibreFieldList()
	for _, c := range lib.Fields.CustomCol {
		fields = append(fields, c)
	}
	return fields
}

func (lib *Lib) RenderFieldMetaSql() string {
	var meta []string
	for _, field := range lib.AllFields() {
		var sql bytes.Buffer
		err := lib.bookTmpl.ExecuteTemplate(&sql, "fieldMeta", field)
		if err != nil {
			log.Println("executing template:", err)
		}
		meta = append(meta, sql.String())
	}
	return strings.Join(meta, ",")
}

func (lib *Lib) fieldMetaStmt() string {

	tmpl := template.Must(template.New("").Parse(fieldMetaTmpl))
	var stmt bytes.Buffer
	err := tmpl.Execute(&stmt, lib.RenderFieldMetaSql())
	if err != nil {
		log.Println("executing template:", err)
	}

	return stmt.String()
}

const fieldMetaTmpl = `SELECT
JSON_OBJECT(
{{.}}
) fieldMeta

FROM preferences 
WHERE key = 'field_metadata'
`

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

func DefaultCalibreFieldList() []string {
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
	fields := DefaultCalibreFieldList()
	for _, f := range defaultFields.urCols {
		fields = append(fields, f)
	}
	return fields
}
