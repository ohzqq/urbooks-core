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

func UnmarshalAudibleApiProduct(d []byte) *Book {
	var data map[string]json.RawMessage
	err := json.Unmarshal(d, &data)
	if err != nil {
		log.Fatalf("issue unmarshalling audible api book %v\n", err)
	}

	book := NewBook()
	for f, dd := range data {
		switch f {
		case "narrators", "authors", "series":
			var c []map[string]string
			err := json.Unmarshal(dd, &c)
			if err != nil {
				log.Fatal(err)
			}

			if f == "series" {
				if len(c) > 0 {
					series := c[0]
					book.GetField("series").Item().
						Set("value", series["title"]).
						Set("position", series["sequence"])
					book.GetField("position").SetMeta(series["sequence"])
				}
				break
			}

			var contributors *Field
			switch f {
			case "narrators":
				contributors = book.AddField(NewCollection("#narrators")).SetIsNames().SetIsEditable().SetIsCustom()
			case "authors":
				contributors = book.GetField(f)
			}
			var cc []string
			for _, contributor := range c {
				cc = append(cc, contributor["name"])
			}
			contributors.SetMeta(cc)
		case "title", "release_date", "publisher_summary", "language", "publisher_name":
			var val string
			err := json.Unmarshal(dd, &val)
			if err != nil {
				log.Fatal(err)
			}
			switch f {
			case "title":
				book.GetField("title").SetMeta(val)
			case "release_date":
				book.GetField("published").SetMeta(val)
			case "publisher_summary":
				book.GetField("description").SetMeta(val)
			case "language":
				book.GetField("languages").SetMeta(val)
			case "publisher_name":
				book.GetField("publisher").SetMeta(val)
			}
		case "product_images":
			var val = make(map[string]string)
			err := json.Unmarshal(dd, &val)
			if err != nil {
				log.Fatal(err)
			}
			book.GetField("cover").Item().Set("url", val["500"])
		case "runtime_length_min":
		}
	}
	return book
}

//func UnmarshalBookSchema(d []byte) *Book {
//  var data []map[string]interface{}
//  err := json.Unmarshal(d, &data)
//  if err != nil {
//    log.Fatalf("issue unmarshalling schema.org book %v\n", err)
//  }

//  book := NewBook()
//  for f, d := range data[0] {
//    switch f {
//    case "name":
//      book.GetField("title").SetData(d)
//    case "datePublished":
//      book.GetField("published").SetData(d)
//    case "image":
//      book.GetField("cover").Item().Set("url", d.(string))
//    case "inLanguage":
//      book.GetField("languages").Collection().AddItem().Set("value", d.(string))
//    case "IsPartOf":
//      book.GetField("series").SetData(d)
//    case "url":
//      book.GetField("uri").SetData(d)
//    case "author":
//      authors := book.GetField("authors").Collection()
//      for _, author := range d.([]interface{}) {
//        authors.AddItem().Set("value", author.(map[string]interface{})["name"].(string))
//      }
//    case "readBy":
//      narrators := book.AddField(NewCollection("#narrators")).
//        SetIsNames().
//        SetIsMultiple().
//        Collection()
//      for _, narrator := range d.([]interface{}) {
//        narrators.AddItem().Set("value", narrator.(map[string]interface{})["name"].(string))
//      }
//    case "duration":
//      book.AddField(NewColumn("#duration")).SetData(d)
//    case "keywords":
//      tags := book.GetField("tags").Collection()
//      for _, tag := range d.([]interface{}) {
//        tags.AddItem().Set("value", tag.(string))
//      }
//    case "publisher":
//      book.GetField("publisher").Item().Set("value", d.(string))
//    case "identifier":
//      book.GetField("identifiers").Collection().AddItem().Set("value", d.(string))
//    }
//  }
//  return book
//}
