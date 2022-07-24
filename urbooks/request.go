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

type Request struct {
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

func NewRequest(l string) *Request {
	req := &Request{}
	req.query = url.Values{}
	req.query.Add("library", l)
	req.lib = l
	req.library = GetLib(l)
	//fmt.Printf("%+V\n\n", req.query)
	return req
}

func (r *Request) String() string {
	u := url.URL{}
	u.RawQuery = r.query.Encode()
	u.Path = path.Join("/", r.prefix, r.endpoint, r.resourceID, r.file)
	return u.String()
}

func (r *Request) JoinPath() string {
	return path.Join(r.prefix, r.endpoint, r.resourceID, r.file)
}

func (r *Request) Query() url.Values {
	return r.query
}

func (r *Request) Set(key, value string) *Request {
	switch key {
	case "endpoint":
		r.endpoint = value
	case "prefix":
		r.prefix = value
	default:
		r.query.Set(key, value)
	}
	return r
}

func (r *Request) Get(key string) string {
	switch key {
	case "endpoint":
		return r.endpoint
	case "prefix":
		return r.prefix
	default:
		return r.query.Get(key)
	}
}

func (r *Request) Response() []byte {
	lib := GetLib(r.lib)
	//fmt.Printf("%+v\n", lib)
	return lib.DB.Get(r.String())
}

func (r *Request) GetResponse() Response {
	//var resp response
	resp := Response{}
	err := json.Unmarshal(r.library.DB.Get(r.String()), &resp)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

//func (r *Request) UnmarshalBooks() BookResponse {
//  lib := GetLib(r.lib)
//  resp := Response{}
//  err := json.Unmarshal(lib.DB.Get(r.String()), &resp)
//  if err != nil {
//    log.Fatal(err)
//  }
//  return resp
//}

func (r *Request) From(table string) *Request {
	r.endpoint = table
	return r
}

func (r *Request) Find(ids string) *Request {
	r.query.Add("ids", ids)
	return r
}

func (r *Request) ID(id string) *Request {
	r.resourceID = id
	return r
}

func (r *Request) Fields(fields string) *Request {
	r.query.Add("fields", fields)
	return r
}

func (r *Request) Limit(limit string) *Request {
	r.query.Add("itemsPerPage", limit)
	return r
}

func (r *Request) Page(page string) *Request {
	r.query.Add("currentPage", page)
	return r
}

func (r *Request) Sort(sort string) *Request {
	r.query.Add("sort", sort)
	return r
}

func (r *Request) Order(order string) *Request {
	r.query.Add("order", order)
	return r
}

func (r *Request) Desc() *Request {
	r.query.Add("order", "desc")
	return r
}
