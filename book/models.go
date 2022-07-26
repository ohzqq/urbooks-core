package book

import (
	"encoding/json"
	"log"
)

type BookSchema struct {
	Abridged      string
	Author        []PersonSchema
	BookFormat    string
	Name          string `json:"name"`
	DatePublished string `json:"datePublished"`
	Description   string `json:"description"`
	Duration      string
	Identifier    string
	Image         string
	InLanguage    string `json:"inLanguage"`
	IsPartOf      string `json:"isPartOf"`
	Keywords      []string
	Position      string
	Publisher     string
	ReadBy        []PersonSchema
	Url           string
}

type PersonSchema struct {
	Type string `json:"@type"`
	Name string
}

func UnmarshalBookSchema(d []byte) *Book {
	var data []map[string]interface{}
	err := json.Unmarshal(d, &data)
	if err != nil {
		log.Fatalf("issue unmarshalling schema.org book %v\n", err)
	}

	book := NewBook()
	for f, d := range data[0] {
		switch f {
		case "name":
			book.GetField("title").SetData(d)
		case "datePublished":
			book.GetField("published").SetData(d)
		case "image":
			book.GetField("cover").Item().Set("url", d.(string))
		case "inLanguage":
			book.GetField("languages").Collection().AddItem().Set("value", d.(string))
		case "IsPartOf":
			book.GetField("series").SetData(d)
		case "url":
			book.GetField("uri").SetData(d)
		case "author":
			authors := book.GetField("authors").Collection()
			for _, author := range d.([]interface{}) {
				authors.AddItem().Set("value", author.(map[string]interface{})["name"].(string))
			}
		case "readBy":
			narrators := book.AddField(NewCollection("#narrators")).
				SetIsNames().
				SetIsMultiple().
				Collection()
			for _, narrator := range d.([]interface{}) {
				narrators.AddItem().Set("value", narrator.(map[string]interface{})["name"].(string))
			}
		case "duration":
			book.AddField(NewColumn("#duration")).SetData(d)
		case "keywords":
			tags := book.GetField("tags").Collection()
			for _, tag := range d.([]interface{}) {
				tags.AddItem().Set("value", tag.(string))
			}
		case "publisher":
			book.GetField("publisher").Item().Set("value", d.(string))
		case "identifier":
			book.GetField("identifiers").Collection().AddItem().Set("value", d.(string))
		}
	}
	return book
}
