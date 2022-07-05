package urbooks

type Column struct {
	Field
	meta string
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

func (c Column) FieldMeta() Field {
	return c.Field
}

func (c *Column) Set(v string) {
	c.meta = v
}

func (c Column) IsNull() bool {
	return c.meta == ""
}
