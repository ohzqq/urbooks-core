package urbooks

import (
	"encoding/json"
	"fmt"

	"github.com/ohzqq/urbooks-core/book"
)

type ResponseLinks map[string]string

type ResponseMeta map[string]string

type ResponseErrors []map[string]string

type Response struct {
	Links  ResponseLinks  `json:"links"`
	Meta   ResponseMeta   `json:"meta"`
	Errors ResponseErrors `json:"errors"`
}

func ParseResponse(d []byte) (Response, error) {
	resp := Response{}
	err := json.Unmarshal(d, &resp)
	if err != nil {
		return resp, fmt.Errorf("Response failed to unmarshal")
	}
	return resp, nil
}

type BookResponse struct {
	Response Response
	books    book.Books
}

func ParseBookResponse(d []byte) (BookResponse, error) {
	response := BookResponse{}
	err := json.Unmarshal(d, &response)
	if err != nil {
		return response, fmt.Errorf("%v\n", err)
	}
	return response, nil
}

func (b BookResponse) Books() []*book.Book {
	return b.books.EachBook()
}

func (b BookResponse) GetMeta(k string) string {
	if m := b.Response.Meta[k]; m != "" {
		return m
	}
	return ""
}

func (b BookResponse) GetLink(k string) string {
	if m := b.Response.Links[k]; m != "" {
		return m
	}
	return ""
}

func (b *BookResponse) UnmarshalJSON(d []byte) error {
	var err error

	b.Response, err = ParseResponse(d)
	if err != nil {
		return err
	}

	b.books, err = book.ParseBooks(d)
	if err != nil {
		return err
	}

	return nil
}

type CatResponse struct {
	Response
	data *book.Category
}

func ParseCatResponse(d []byte) (CatResponse, error) {
	response := CatResponse{}
	err := json.Unmarshal(d, &response)
	if err != nil {
		return response, fmt.Errorf("%v\n", err)
	}
	return response, nil
}

func (c CatResponse) Items() []*book.Item {
	return c.data.EachItem()
}

func (c *CatResponse) UnmarshalJSON(d []byte) error {
	var err error

	c.Response, err = ParseResponse(d)
	if err != nil {
		return err
	}

	c.data, err = book.ParseCategory(d)
	if err != nil {
		return err
	}

	return nil
}

//func ParseCategory(r []byte) *CatResponse {
//  var (
//    response CatResponse
//    err      error
//  )

//  var resp map[string]json.RawMessage
//  err = json.Unmarshal(r, &resp)
//  if err != nil {
//    fmt.Println("error:", err)
//  }

//  response.Response = ParseResponse(resp)

//  //response.data.Field.query = url.Values{}
//  //response.data.Field.query.Set("library", response.Meta["library"])

//  lib := Lib(response.Meta["library"])

//  cats := Category{
//    Field: lib.DB.GetField(response.Meta["endpoint"]),
//  }
//  err = json.Unmarshal(resp["data"], &cats.items)
//  if err != nil {
//    fmt.Println("error:", err)
//  }

//  return &response
//}

//func (c CatResponse) Label() string {
//  return c.data.Field.Label
//}
