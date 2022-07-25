package urbooks

import (
	//"strings"
	"encoding/json"
	"log"

	//"strconv"

	"fmt"
	"net/url"
	"path"

	//"github.com/ohzqq/urbooks-core/calibredb"

	"golang.org/x/exp/slices"
)

type request struct {
	path       []string
	endpoint   string
	prefix     string
	resourceID string
	file       string
	query      url.Values
	lib        string
	library    *Library
}

func Get(u string) ([]byte, error) {
	url, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	var lib *Library
	if url.Query().Has("library") {
		l := url.Query().Get("library")
		if slices.Contains(Libraries(), l) {
			lib = GetLib(l)
		} else {
			return []byte{}, fmt.Errorf("She %v doesn't even go here!", l)
		}
	} else {
		lib = DefaultLib()
		url.Query().Set("library", lib.Name)
	}
	//fmt.Printf("%v\n", url)
	return lib.DB.Get(url.String()), nil
}

func Query(u string) ([]byte, error) {
	url, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	var lib *Library
	if url.Query().Has("library") {
		l := url.Query().Get("library")
		if slices.Contains(Libraries(), l) {
			lib = GetLib(l)
		} else {
			return []byte{}, fmt.Errorf("She %v doesn't even go here!", l)
		}
	} else {
		lib = DefaultLib()
		url.Query().Set("library", lib.Name)
	}
	//fmt.Printf("%v\n", url)
	return lib.DB.Get(url.String()), nil
}

func NewRequest(l string) *request {
	req := &request{
		query:   url.Values{},
		lib:     l,
		library: GetLib(l),
	}
	req.query.Add("library", l)
	//fmt.Printf("%+V\n\n", req.query)
	return req
}

func (r *request) String() string {
	u := url.URL{}
	u.RawQuery = r.query.Encode()
	u.Path = path.Join("/", r.prefix, r.endpoint, r.resourceID, r.file)
	return u.String()
}

func (r *request) JoinPath() string {
	return path.Join(r.prefix, r.endpoint, r.resourceID, r.file)
}

func (r *request) Query() *url.Values {
	return &r.query
}

func (r *request) Set(key, value string) *request {
	switch key {
	case "endpoint":
		r.endpoint = value
	case "prefix":
		r.prefix = value
	default:
		r.Query().Set(key, value)
	}
	return r
}

func (r *request) Get(key string) string {
	switch key {
	case "endpoint":
		return r.endpoint
	case "prefix":
		return r.prefix
	default:
		return r.Query().Get(key)
	}
}

func (r *request) Response() []byte {
	lib := GetLib(r.lib)
	//fmt.Printf("%+v\n", lib)
	return lib.DB.Get(r.String())
}

func (r *request) GetResponse() Response {
	//var resp response
	resp := Response{}
	err := json.Unmarshal(r.library.DB.Get(r.String()), &resp)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

//func (r *request) UnmarshalBooks() BookResponse {
//  lib := GetLib(r.lib)
//  resp := Response{}
//  err := json.Unmarshal(lib.DB.Get(r.String()), &resp)
//  if err != nil {
//    log.Fatal(err)
//  }
//  return resp
//}

func (r *request) From(table string) *request {
	r.endpoint = table
	return r
}

func (r *request) Find(ids string) *request {
	r.query.Add("ids", ids)
	return r
}

func (r *request) ID(id string) *request {
	r.resourceID = id
	return r
}

func (r *request) Fields(fields string) *request {
	r.query.Add("fields", fields)
	return r
}

func (r *request) Limit(limit string) *request {
	r.query.Add("itemsPerPage", limit)
	return r
}

func (r *request) Page(page string) *request {
	r.query.Add("currentPage", page)
	return r
}

func (r *request) Sort(sort string) *request {
	r.query.Add("sort", sort)
	return r
}

func (r *request) Order(order string) *request {
	r.query.Add("order", order)
	return r
}

func (r *request) Desc() *request {
	r.query.Add("order", "desc")
	return r
}
