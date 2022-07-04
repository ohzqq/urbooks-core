package urbooks

import (
	//"strings"
	"log"
	//"strconv"
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	//"github.com/ohzqq/urbooks-core/calibredb"

	"golang.org/x/exp/slices"
)

var _ = fmt.Sprintf("%v", "")

func Get(u string) ([]byte, error) {
	url, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	var lib *Library
	if url.Query().Has("library") {
		l := url.Query().Get("library")
		if slices.Contains(Libraries(), l) {
			lib = Lib(l)
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

type Request struct {
	path       []string
	endpoint   string
	prefix     string
	resourceID string
	file       string
	query      url.Values
	lib        *Library
}

func NewRequest(lib string) *Request {
	req := &Request{}
	req.lib = Lib(lib)
	req.query = url.Values{}
	req.query.Add("library", lib)
	//fmt.Printf("%+v\n\n", req)
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
	return r.lib.DB.Get(r.String())
}

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

type ResponseLinks map[string]string

type ResponseLink struct {
	Rel  string
	Href string
}

type ResponseMeta map[string]string

type ResponseErrors []map[string]string

type Response struct {
	Links  ResponseLinks  `json:"links"`
	Meta   ResponseMeta   `json:"meta"`
	Errors ResponseErrors `json:"errors"`
}

func ParseResponse(resp map[string]json.RawMessage) Response {
	var (
		response Response
		err      error
	)

	err = json.Unmarshal(resp["meta"], &response.Meta)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = json.Unmarshal(resp["errors"], &response.Errors)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = json.Unmarshal(resp["links"], &response.Links)
	if err != nil {
		fmt.Println("error:", err)
	}

	return response
}
