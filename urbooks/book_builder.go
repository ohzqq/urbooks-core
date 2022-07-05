package urbooks

func NewBookMeta() BookMeta {
	return make(BookMeta)
}

func (meta *BookMeta) NewColumn(k string) *Column {
	col := NewColumn()
	m := *meta
	if lib := m["library"]; !lib.IsNull() {
		col.Field = Lib(lib.Value()).DB.GetField(k)
	}
	m[k] = col
	return col
}

func (meta *BookMeta) NewItem(k string) *Item {
	item := NewCategoryItem()
	m := *meta
	if lib := m["library"]; !lib.IsNull() {
		item.Field = Lib(lib.Value()).DB.GetField(k)
	}
	m[k] = item
	return item
}

func (meta *BookMeta) NewCategory(k string) *Category {
	cat := NewCategory(k)
	m := *meta
	if lib := m["library"]; !lib.IsNull() {
		cat.Field = Lib(lib.Value()).DB.GetField(k)
	}
	m[k] = cat
	return cat
}

func (c *Category) AddItem() *Item {
	item := NewCategoryItem()
	c.items = append(c.items, item)
	return item
}
