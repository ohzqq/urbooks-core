package urbooks

import (
	"strconv"

	"github.com/ohzqq/urbooks-core/calibredb"
)

type Item struct {
	Field *calibredb.Field
	meta  map[string]string
}

func NewCategoryItem() *Item {
	return &Item{meta: make(map[string]string)}
}

func (i Item) Value() string {
	return i.Get("value")
}

func (i Item) String() string {
	return i.Get("value")
}

func (i Item) FieldMeta() *calibredb.Field {
	return i.Field
}

func (i Item) IsNull() bool {
	return len(i.meta) == 0
}

func (i Item) ID() string {
	return i.Get("id")
}

func (i Item) TotalBooks() int {
	if t := i.Get("books"); t != "" {
		b, err := strconv.Atoi(t)
		if err != nil {
			return 0
		}
		return b
	}
	return 0
}

func (i Item) URL() string {
	return i.Get("url")
}

func (i Item) Get(v string) string {
	return i.meta[v]
}

func (i *Item) Set(k, v string) *Item {
	i.meta[k] = v
	return i
}

func (i *Item) SetValue(v string) *Item {
	i.meta["value"] = v
	return i
}
