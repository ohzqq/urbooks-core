package urbooks

import (
	"net/url"

	"github.com/ohzqq/urbooks-core/calibredb"
)

type Column struct {
	Field *calibredb.Field
	query url.Values
	meta  string
}

func NewColumn() *Column {
	return &Column{}
}

func (c Column) URL() string {
	return ""
}

func (c Column) Value() string {
	return c.meta
}

func (c Column) String() string {
	return c.meta
}

func (c Column) FieldMeta() *calibredb.Field {
	return c.Field
}

func (c *Column) SetValue(v string) *Column {
	c.meta = v
	return c
}

func (c Column) IsNull() bool {
	return c.meta == ""
}
