package urbooks

import (
	"strings"

	"github.com/ohzqq/urbooks-core/calibredb"
)

type Category struct {
	Field *calibredb.Field
	items []*Item
	value string
	item  Item
}

func NewCategory() *Category {
	cat := &Category{}
	return cat
}

func (c *Category) AddItem() *Item {
	item := NewCategoryItem()
	c.items = append(c.items, item)
	return item
}

func (c Category) String() string {
	return c.Join("value")
}

const (
	nameSep = " & "
	itemSep = ", "
)

func (c Category) Join(v string) string {
	var meta []string
	for _, field := range c.items {
		meta = append(meta, field.Get(v))
	}
	switch c.Field.IsNames {
	case true:
		return strings.Join(meta, nameSep)
	default:
		return strings.Join(meta, itemSep)
	}
}

func (c *Category) Split() {
	sep := itemSep
	if c.Field.IsNames {
		sep = nameSep
	}
	for _, val := range strings.Split(c.value, sep) {
		c.AddItem().SetValue(val)
	}
}

func (c Category) IsNull() bool {
	return len(c.Items()) == 0
}

func (c Category) FieldMeta() *calibredb.Field {
	return c.Field
}

func (c Category) Value() string {
	return c.Join("value")
}

func (c Category) Items() []*Item {
	return c.items
}

func (c Category) URL() string {
	return c.Field.Label + "/"
}

func (c *Category) SetFieldMeta(k, v string) *Category {
	switch k {
	case "isNames":
		if v == "true" {
			c.Field.IsNames = true
		}
	}
	return c
}
