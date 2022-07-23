package book

import (
	"encoding/json"
	"fmt"
)

type Category struct {
	*Field
	lib string
}

func NewCategory(name string) *Category {
	return &Category{
		Field: NewCollection(name),
	}
}

func ParseCategory(r []byte) (*Category, error) {
	var err error

	var resp map[string]json.RawMessage
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response error: %v\n", err)
	}

	var rmeta map[string]string
	err = json.Unmarshal(resp["meta"], &rmeta)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response meta error: %v\n", err)
	}

	cat := NewCategory(rmeta["categoryLabel"])
	cat.lib = rmeta["library"]

	err = json.Unmarshal(resp["data"], cat)
	if err != nil {
		return nil, fmt.Errorf("unmarshal cat data error: %v\n", err)
	}

	return cat, nil
}

func (c *Category) GetMeta() *Collection {
	return c.Meta.(*Collection)
}

func (c *Category) EachItem() []*Item {
	return c.GetMeta().EachItem()
}
