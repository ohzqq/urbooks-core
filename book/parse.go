package book

import (
	"encoding/json"
	"fmt"
	"log"
)

func ParseBooks(r []byte) Books {
	var (
		err error
	)

	var resp map[string]json.RawMessage
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Printf("response err: %v\n", err)
	}

	var rmeta map[string]string
	err = json.Unmarshal(resp["meta"], &rmeta)
	if err != nil {
		fmt.Printf("response err: %v\n", err)
	}
	lib := rmeta["library"]

	var rawbooks []map[string]json.RawMessage
	err = json.Unmarshal(resp["data"], &rawbooks)
	if err != nil {
		fmt.Printf("book parsing error: %v\n", err)
	}

	var books Books
	for _, b := range rawbooks {
		book := NewBook()
		book.lib = lib
		for key, value := range b {
			field := book.GetField(key)
			field.Library = lib

			if key != "customColumns" {
				field.Meta.UnmarshalJSON(value)
			}

			field.Data = value

			if key != field.JsonLabel {
				log.Fatalf("json: %v\n field meta: %v\n", key, field.JsonLabel)
			}

			if key == "customColumns" {
				var custom = make(map[string]map[string]json.RawMessage)
				err = json.Unmarshal(value, &custom)
				if err != nil {
					fmt.Printf("custom column parsing error: %v\n", err)
				}
				for name, cdata := range custom {
					col := &Field{
						IsCustom:     true,
						Data:         cdata["data"],
						CalibreLabel: name,
						JsonLabel:    name,
						IsEditable:   true,
					}

					meta := make(map[string]string)
					err := json.Unmarshal(cdata["meta"], &meta)
					if err != nil {
						fmt.Printf("custom column parsing error: %v\n", err)
					}

					switch meta["is_multiple"] {
					case "true":
						col.IsMultiple = true
						col.IsCollection = true
						col.Meta = &Collection{}
					case "false":
						col.IsColumn = true
						col.Meta = NewColumn()
					}
					col.Meta.UnmarshalJSON(col.Data)

					if meta["is_names"] == "true" {
						col.IsNames = true
					}

					book.AddField(col)
				}
			}
		}
		books = append(books, book)
	}
	return books
}

//func (i *Item) UnmarshalJSON(b []byte) error {
//  i.meta = make(map[string]string)
//  if err := json.Unmarshal(b, &i.meta); err != nil {
//    return err
//  }
//  return nil
//}
