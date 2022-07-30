package urbooks

import (
	"encoding/json"
	"log"

	"github.com/ohzqq/urbooks-core/book"
)

type ResponseLinks map[string]string

type ResponseMeta map[string]string

type ResponseErrors []map[string]string

type Response struct {
	ResponseLinks  ResponseLinks   `json:"links"`
	ResponseMeta   ResponseMeta    `json:"meta"`
	ResponseErrors ResponseErrors  `json:"errors"`
	Data           json.RawMessage `json:"data"`
}

type BookResponse struct {
	Response
	Books book.Books `json:"data"`
}

type CategoryResponse struct {
	Response
	Category book.Category `json:"data"`
}

func (r Response) ParseBooks() BookResponse {
	books := BookResponse{Response: r}
	err := json.Unmarshal(r.Data, &books.Books)
	if err != nil {
		log.Fatal(err)
	}
	return books
}

func (r *Response) GetResponseMeta(key string) string {
	return r.ResponseMeta[key]
}

func (r *Response) GetResponseLink(key string) string {
	return r.ResponseLinks[key]
}

func GetResponse(r *request) Response {
	resp := Response{}
	err := json.Unmarshal(r.library.DB.Get(r.String()), &resp)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}
