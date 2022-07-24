package calibredb

import (
	"encoding/json"
	"log"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
)

type response struct {
	Links         map[string]string `json:"links"`
	Data          any               `json:"data"`
	Errors        []responseErr     `json:"errors"`
	Meta          map[string]string `json:"meta"`
	booksInCat    string
	mtx           sync.Mutex
	numberOfItems int
}

func newResponse() *response {
	var resp response
	resp.Links = make(map[string]string)
	resp.Meta = make(map[string]string)
	return &resp
}

func (lib *Lib) setResponseURL() {
	if lib.response.numberOfItems > 1 {
		if !lib.Request.query.Has("currentPage") {
			lib.Request.query.Set("currentPage", "1")
		}
	}
	url := url.URL{
		Path:     path.Join(lib.Request.path),
		RawQuery: lib.Request.query.Encode(),
	}
	lib.response.addLink("self", url.String())
}

func (lib *Lib) setResponseMeta() {
	lib.response.addMeta("library", lib.Name)
	lib.response.addMeta("numberOfItems", strconv.Itoa(lib.response.numberOfItems))
	lib.response.addMeta("currentPage", lib.Request.query.Get("currentPage"))
	lib.response.addMeta("itemsPerPage", lib.Request.query.Get("itemsPerPage"))
	lib.response.addMeta("endpoint", lib.Request.cat)
	lib.response.addMeta("categoryLabel", lib.Request.cat)

	if lib.response.booksInCat != "" {
		lib.response.addMeta("categoryLabel", lib.response.booksInCat)
	}
}

func (lib *Lib) setResponseData(data any) *Lib {
	lib.response.Data = data
	return lib
}

func (lib *Lib) calculatePagination() {
	var (
		prev  int
		next  int
		last  int
		first = 1
	)

	if lib.Request.currentPage <= first {
		prev = first
	} else {
		prev = lib.Request.currentPage - 1
	}

	last = lib.response.numberOfItems / lib.Request.itemsPerPage
	if r := lib.response.numberOfItems % lib.Request.itemsPerPage; r == 0 {
		last = last
	} else {
		last = last + 1
	}

	if lib.Request.currentPage >= last {
		next = last
	} else {
		next = lib.Request.currentPage + 1
	}

	rPath := path.Join(lib.Request.path)
	lib.Request.query.Set("currentPage", strconv.Itoa(first))
	firstPage := url.URL{Path: rPath, RawQuery: lib.Request.query.Encode()}
	lib.response.addLink("first", firstPage.String())

	lib.Request.query.Set("currentPage", strconv.Itoa(next))
	nextPage := url.URL{Path: rPath, RawQuery: lib.Request.query.Encode()}
	lib.response.addLink("next", nextPage.String())

	lib.Request.query.Set("currentPage", strconv.Itoa(prev))
	prevPage := url.URL{Path: rPath, RawQuery: lib.Request.query.Encode()}
	lib.response.addLink("prev", prevPage.String())

	lib.Request.query.Set("currentPage", strconv.Itoa(last))
	lastPage := url.URL{Path: rPath, RawQuery: lib.Request.query.Encode()}
	lib.response.addLink("last", lastPage.String())

	lib.Request.query.Set("currentPage", strconv.Itoa(lib.Request.currentPage))
}

type responseErr struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}

func (r *response) addErr(e error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	errMsg := strings.Split(e.Error(), ":")
	respErr := responseErr{Status: errMsg[0], Detail: errMsg[1]}
	r.Errors = append(r.Errors, respErr)
}

func (r *response) addLink(title, url string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.Links[title] = url
}

func (r *response) addMeta(key string, value string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.Meta[key] = value
}

func (r *response) json() []byte {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	result, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

type field string

func (f field) MarshalJSON() ([]byte, error) {
	return []byte(f), nil
}

func convertFields(book map[string]interface{}) map[string]field {
	meta := make(map[string]field)
	for key, val := range book {
		meta[key] = field(val.(string))
	}
	return meta
}
