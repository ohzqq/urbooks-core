package urbooks

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type BookResponse struct {
	Response
	books Books
}

func (b BookResponse) Books() []*Book {
	return b.books.Books
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

func ParseBooks(r []byte) *BookResponse {
	var (
		response BookResponse
		err      error
	)

	var resp map[string]json.RawMessage
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println("error:", err)
	}

	response.Response = ParseResponse(resp)

	var books []map[string]json.RawMessage
	err = json.Unmarshal(resp["data"], &books)
	if err != nil {
		fmt.Println("error:", err)
	}

	response.books.query = url.Values{}
	response.books.query.Set("library", response.GetMeta("library"))
	lib := Lib(response.GetMeta("library"))

	for _, book := range books {
		bb := NewBook(lib.Name)
		formats := bb.NewCategory("formats")
		formats.query = response.books.query
		for key, val := range book {
			var err error
			switch key {
			case "cover":
				item := formats.AddItem()
				item.query = response.books.query
				err = json.Unmarshal(val, &item)
			case "formats":
				err = json.Unmarshal(val, &formats.items)
			case "series", "publishers":
				item := bb.NewItem(key)
				item.query = response.books.query
				err = json.Unmarshal(val, &item)
				u := &url.URL{Path: item.Get("uri"), RawQuery: response.books.query.Encode()}
				item.Set("url", u.String())
			case "authors", "narrators", "identifiers", "languages", "tags":
				cat := bb.NewCategory(key)
				cat.query = response.books.query

				err = json.Unmarshal(val, &cat.items)

				for _, item := range cat.items {
					u := &url.URL{Path: item.Get("uri"), RawQuery: response.books.query.Encode()}
					item.Set("url", u.String())
				}
			default:
				col := bb.NewColumn(key)
				col.query = response.books.query
				err = json.Unmarshal(val, &col.meta)
			}
			if err != nil {
				fmt.Printf("%v: %v\n", key, err)
			}
		}
		response.books.Add(bb)
	}
	return &response
}

func (i *Item) UnmarshalJSON(b []byte) error {
	i.meta = make(map[string]string)
	if err := json.Unmarshal(b, &i.meta); err != nil {
		return err
	}
	return nil
}

type CatResponse struct {
	Response
	data Category
}

func ParseCategory(r []byte) *CatResponse {
	var (
		response CatResponse
		err      error
	)

	var resp map[string]json.RawMessage
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println("error:", err)
	}

	response.Response = ParseResponse(resp)

	//response.data.Field.query = url.Values{}
	//response.data.Field.query.Set("library", response.Meta["library"])

	lib := Lib(response.Meta["library"])

	cats := Category{
		Field: lib.DB.GetField(response.Meta["endpoint"]),
	}
	err = json.Unmarshal(resp["data"], &cats.items)
	if err != nil {
		fmt.Println("error:", err)
	}

	return &response
}

func (c CatResponse) Items() []*Item {
	return c.data.Items()
}

func (c CatResponse) Label() string {
	return c.data.Field.Label
}
