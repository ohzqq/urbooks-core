package book

import (
	"encoding/json"
	"log"
)

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
