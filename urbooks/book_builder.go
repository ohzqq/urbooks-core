package urbooks

func NewBook() BookMeta {
	return make(BookMeta)
}

func (meta *BookMeta) NewColumn(k string) *Column {
	col := NewColumn()
	m := *meta
	m[k] = col
	return col
}

func (meta *BookMeta) NewItem(k string) *Item {
	item := NewCategoryItem()
	m := *meta
	m[k] = item
	return item
}

func (meta *BookMeta) NewCategory(k string) *Category {
	cat := NewCategory(k)
	m := *meta
	m[k] = cat
	return cat
}

func (c *Category) AddItem() *Item {
	item := NewCategoryItem()
	c.items = append(c.items, item)
	return item
}
