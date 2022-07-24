package urbooks

import (
	"encoding/json"
	"log"
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

func GetResponse(r *Request) Response {
	resp := Response{}
	err := json.Unmarshal(r.library.DB.Get(r.String()), &resp)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}
