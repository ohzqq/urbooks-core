package calibredb

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slices"
)

type Lib struct {
	Name           string
	Path           string
	dbPath         string
	Fields         *dbFieldTypes
	fieldMeta      map[string]map[string]interface{}
	shortFieldMeta map[string]map[string]bool
	Preferences    *Preferences
	CustCols       []map[string]string
	db             *sqlx.DB
	Request        *request
	response       *response
	mtx            sync.Mutex
	bookTmpl       *template.Template
}

//go:embed sql/*
var sqlTmpl embed.FS

func NewLib(path string) *Lib {
	lib := Lib{}
	lib.Path = path
	lib.dbPath = "file:" + filepath.Join(path, "metadata.db") + "?cache=shared&mode=ro"
	lib.Name = filepath.Base(path)
	lib.db = lib.connectDB()
	lib.Fields = newLibFields(lib.Name)
	lib.getPreferences()
	lib.getCustCols()
	lib.bookTmpl = template.Must(template.New("book").Funcs(bookTmplFuncs).ParseFS(sqlTmpl, "sql/*"))

	//fmt.Printf("field meta %+v\n", lib.fieldMeta)
	return &lib
}

var (
	bookTmplFuncs = map[string]any{
		"GetCalibreField":         GetCalibreField,
		"GetJsonField":            GetJsonField,
		"GetTableColumns":         GetTableColumns,
		"FieldList":               FieldList,
		"DefaultCalibreFieldList": DefaultCalibreFieldList,
		"GetFieldMeta":            GetFieldMeta,
	}
)

func (lib *Lib) Get(u string) []byte {
	lib.mtx.Lock()
	defer lib.mtx.Unlock()
	//fmt.Printf("%+v\n", u)
	lib.response = newResponse()

	var err error
	lib.Request, err = lib.newRequest(u)
	if err != nil {
		lib.response.addErr(err)
		return lib.response.json()
	}

	//fmt.Printf("%+v\n", lib.Request)
	lib.numberOfItems()

	lib.setResponseURL()

	lib.setResponseMeta()

	if !lib.Request.allItems {
		if lib.response.numberOfItems > lib.Request.itemsPerPage {
			if lib.Request.bookQuery {
				lib.calculatePagination()
			}
		}
	}

	var data any
	if lib.Request.cat == "preferences" {
		if lib.Request.HasFields {
			data = lib.GetPref("field_meta")
		} else {
			data = lib.GetPreferences()
		}
		//return lib.GetPreferences()
	} else {
		data = lib.queryDB()
	}

	lib.setResponseData(data)
	json := lib.response.json()
	return json
}

func (lib *Lib) validEndpoint(point string) bool {
	end := lib.Categories()
	end = append(end, "preferences", "customColumns", "books")
	return slices.Contains(end, point)
}

func (lib *Lib) Categories() []string {
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
	for _, f := range lib.CustCols {
		if f["is_category"] == "true" {
			fields = append(fields, f["label"])
		}
	}
	return fields
}

func (lib *Lib) connectDB() *sqlx.DB {
	database, err := sqlx.Open("sqlite3", lib.dbPath)
	if err != nil {
		log.Fatal(err)
	}
	return database
}

type dbData []map[string]map[string]string

func (lib *Lib) queryDB() any {
	var (
		query, args = lib.queryStmt()
	)
	//fmt.Println(query)

	rows, err := lib.db.Queryx(query, args...)
	if err != nil {
		log.Fatal(query)
	}
	defer rows.Close()

	var data []map[string]field
	for rows.Next() {
		m := make(map[string]interface{})
		if err := rows.MapScan(m); err != nil {
			fmt.Errorf("Association %s", err)
		}
		data = append(data, convertFields(m))
		//lib.response.Data = append(lib.response.Data, convertFields(m))
	}
	return data
	//lib.response.Data = data
}

func (lib *Lib) numberOfItems() {
	var (
		stmt  strings.Builder
		total int
	)

	switch len(lib.Request.itemIDs) {
	case 0:
		stmt.WriteString("SELECT COUNT(*) FROM ")
		stmt.WriteString(lib.Request.cat)
		row := lib.db.QueryRowx(stmt.String())
		row.Scan(&total)
		lib.response.numberOfItems = total
	default:
		lib.response.numberOfItems = len(lib.Request.itemIDs)
	}
}

func (lib *Lib) renderSqlTmpl(name string) string {
	var buf bytes.Buffer
	err := lib.bookTmpl.ExecuteTemplate(&buf, name, lib)
	if err != nil {
		log.Println("executing template:", err)
	}

	return buf.String()
}

func (lib *Lib) queryStmt() (string, []interface{}) {
	var args []interface{}
	if lib.Request.bookQuery {
		return lib.bookStmt()
	} else {
		switch table := lib.Request.cat; table {
		case "preferences":
			return prefSql, args
		case "customColumns":
			return customColumnsSql, args
		default:
			return lib.relationStmt(table)
		}
	}
}

func (lib *Lib) booksInCatStmt(table string, id string) string {
	var (
		ids   string
		value string
	)

	lib.Request.PathID = id
	lib.Request.CatLabel = table

	row := lib.db.QueryRowx(lib.renderSqlTmpl("booksInCategory"))
	row.Scan(&value, &ids)

	lib.response.booksInCat = value

	return ids
}

func (lib *Lib) bookStmt() (string, []interface{}) {
	if !lib.Request.HasFields {
		lib.Request.Fields = lib.AllFields()
	}

	return lib.filterQuery(lib.renderSqlTmpl("book"))
}

func (lib *Lib) custStmt() {
	println(lib.renderSqlTmpl("custCol"))
}

// Build association Queries
func (lib *Lib) relationStmt(table string) (string, []interface{}) {
	lib.Request.isSorted = true
	//lib.Request.sort = lib.GetField(table).Column
	field := GetTableColumns(table, lib.Name)
	lib.Request.sort = field["value"]
	if table == "formats" {
		lib.Request.sort = "format"
	}

	return lib.filterQuery(lib.renderSqlTmpl("category"))
}

func (lib *Lib) filterQuery(q string) (string, []interface{}) {
	var (
		stmt strings.Builder
	)

	stmt.WriteString(q)
	stmt.WriteString("\n")

	if len(lib.Request.itemIDs) > 0 {
		if lib.Request.bookQuery {
			stmt.WriteString(" WHERE books.id IN (?) ")
		} else {
			stmt.WriteString(" WHERE id IN (?) ")
		}
		stmt.WriteString("\n")
	}

	stmt.WriteString(" ORDER BY ")
	if lib.Request.isSorted {
		if lib.Request.bookQuery {
			stmt.WriteString(BookSortField(lib.Request.sort))
		} else {
			stmt.WriteString(lib.Request.sort)
		}
	} else if !lib.Request.isSorted {
		if lib.Request.isCustom {
			if lib.Request.bookQuery {
				stmt.WriteString("timestamp ")
			} else {
				stmt.WriteString("value ")
			}
		} else {
			stmt.WriteString("timestamp ")
		}
	}

	stmt.WriteString("\n")

	if lib.Request.desc {
		stmt.WriteString(" DESC ")
		stmt.WriteString("\n")
	}

	if lib.Request.itemsPerPage != 0 {
		stmt.WriteString(" LIMIT ")
		stmt.WriteString(strconv.Itoa(lib.Request.itemsPerPage))
		stmt.WriteString("\n")
	}

	if offset := lib.Request.calculateOffset(); offset != 0 {
		stmt.WriteString(" OFFSET ")
		stmt.WriteString(strconv.Itoa(offset))
		stmt.WriteString("\n")
	}

	stmt.WriteString(" ;")

	var (
		query string
		args  []interface{}
	)
	switch len(lib.Request.itemIDs) > 0 {
	case true:
		var err error
		query, args, err = sqlx.In(stmt.String(), lib.Request.itemIDs)
		if err != nil {
			log.Fatal(err)
		}
	case false:
		query = stmt.String()
	}

	//fmt.Println(query)
	return query, args
}

func BookSortField(f string) string {
	var bookSortField = map[string]string{
		"authorSort":  "author_sort",
		"sortAs":      "sort",
		"added":       "timestamp",
		"languages":   "lang_code",
		"title":       "sort",
		"identifiers": "val",
	}

	switch name := bookSortField[f]; name {
	case "":
		return f
	default:
		return name
	}
}

func GetCalibreField(f string) string {
	var jsonFieldToCalibre = map[string]string{
		"authorSort":  "author_sort",
		"rating":      "rating",
		"description": "comments",
		"modified":    "last_modified",
		"published":   "pubdate",
		"publishers":  "publisher",
		"sortAs":      "sort",
		"added":       "timestamp",
		"position":    "series_index",
	}
	switch name := jsonFieldToCalibre[f]; name {
	case "":
		return f
	default:
		return name
	}
}

func GetJsonField(f string) string {
	var calibreFieldToJson = map[string]string{
		"authors":       "authors",
		"author_sort":   "authorSort",
		"rating":        "rating",
		"publisher":     "publisher",
		"comments":      "description",
		"last_modified": "modified",
		"pubdate":       "published",
		"sort":          "sortAs",
		"timestamp":     "added",
		"series_index":  "position",
	}
	switch name := calibreFieldToJson[f]; name {
	case "":
		return f
	default:
		return name
	}
}

//Fields

const customColumnsSql = `
SELECT 
	IFNULL(JSON_QUOTE(id) , "") id,  
	IFNULL(JSON_QUOTE(label) , "") label,
	IFNULL(JSON_QUOTE(name) , "") name,
	IFNULL(JSON_QUOTE(editable) , "") editable,
	IFNULL(JSON_QUOTE(is_multiple) , "") is_multiple,
	IFNULL(JSON_QUOTE(JSON_EXTRACT(display, "$.is_names")), 0) is_names,
	IFNULL(JSON_QUOTE(JSON_EXTRACT(display, "$.description")), "") description
FROM custom_columns; `
