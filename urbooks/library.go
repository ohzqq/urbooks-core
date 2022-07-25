package urbooks

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/ohzqq/urbooks-core/book"
	"github.com/ohzqq/urbooks-core/calibredb"
	"golang.org/x/exp/slices"
)

func Libraries() []string {
	return cfg.list
}

func Lib(l string) *Library {
	return cfg.libs[l]
}

func GetLib(l string) *Library {
	return cfg.libs[l]
}

func DefaultLib() *Library {
	var lib = Libraries()[0]
	if l := Cfg().Opts["default"]; l != "" {
		if slices.Contains(Libraries(), l) {
			lib = l
		}
	}
	return Lib(lib)
}

type Library struct {
	Cfg            *libCfg
	Name           string
	Path           string
	DefaultRequest *request
	DB             *calibredb.Lib
	pref           dbPreferences
	Pref           dbPreferences
	RawPref        map[string]json.RawMessage
	CustomColumns  Fields
	Books          book.Books
	Category       *book.Category
	*request
	Response
}

func NewLibrary(name, path string) *Library {
	return &Library{
		Cfg:  Cfg().libCfg[name],
		Name: name,
		Path: path,
	}
}

func (l *Library) ConnectDB() *Library {
	l.DB = calibredb.NewLib(l.Path)
	//l.pref = l.DB.Preferences
	return l
}

func (l *Library) IsAudiobooks() bool {
	return l.Cfg.Audiobooks
}

type dbPreferences struct {
	HiddenCategories []string                   `json:"tag_browser_hidden_categories"`
	DisplayFields    json.RawMessage            `json:"book_display_fields"`
	SavedSearches    map[string]string          `json:"saved_searches"`
	FieldMeta        map[string]map[string]bool `json:"field_metadata"`
}

func (l *Library) GetDBPreferences() {
	dbPref := NewRequest(l.Name).From("preferences").Response()
	data := make(map[string]json.RawMessage)
	err := json.Unmarshal(dbPref, &data)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data["data"], &l.Pref)
	if err != nil {
		log.Fatal(err)
	}
}

func (l *Library) GetDBCustomColumns() {
	cols := make(Fields)
	for _, c := range l.DB.CustCols {
		data := []byte(c["meta"])
		field := Field{}
		err := json.Unmarshal(data, &field)
		if err != nil {
			log.Fatalf("cust col fail %v\n", err)
		}
		cols[c["label"]] = &field
	}
	l.CustomColumns = cols
}

func (l *Library) GetBooks() *Library {
	l.request = l.NewRequest().From("books")
	return l
}

func (l *Library) NewRequest() *request {
	l.request = &request{query: url.Values{}, library: l}
	l.Query().Set("library", l.Name)
	return l.request
}

func (l *Library) GetResponse() *Library {
	l.Response = GetResponse(l.request)
	switch l.endpoint {
	case "books":
		err := json.Unmarshal(l.Data, &l.Books)
		if err != nil {
			log.Fatal(err)
		}
	case "field_meta":
		err := json.Unmarshal(l.Data, &l.CustomColumns)
		if err != nil {
			log.Fatal(err)
		}
	default:
		err := json.Unmarshal(l.Data, &l.Category)
		if err != nil {
			log.Fatal(err)
		}
	}
	return l
}

func (l *Library) From(table string) *Library {
	l.NewRequest().Set("endpoint", table)
	return l
}

func (l *Library) Find(ids string) *Library {
	l.Query().Set("ids", ids)
	return l
}

func (l *Library) ID(id string) *Library {
	l.resourceID = id
	return l
}

func (l *Library) Fields(fields string) *Library {
	l.Query().Set("fields", fields)
	return l
}

func (l *Library) Limit(limit string) *Library {
	l.Query().Set("itemsPerPage", limit)
	return l
}

func (l *Library) Page(page string) *Library {
	l.Query().Set("currentPage", page)
	return l
}

func (l *Library) Sort(sort string) *Library {
	l.Query().Set("sort", sort)
	return l
}

func (l *Library) Order(order string) *Library {
	l.Query().Set("order", order)
	return l
}

func (l *Library) Desc() *Library {
	l.Query().Set("order", "desc")
	return l
}
