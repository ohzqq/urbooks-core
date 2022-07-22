package book

import (
	"encoding/json"
	"fmt"
	"log"
)

func ParseBooks(r []byte) Books {
	var books Books
	err := json.Unmarshal(r, &books)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", books)
	return books
}

//func (i *Item) UnmarshalJSON(b []byte) error {
//  i.meta = make(map[string]string)
//  if err := json.Unmarshal(b, &i.meta); err != nil {
//    return err
//  }
//  return nil
//}
