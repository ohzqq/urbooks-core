package urbooks

func NewBookMeta() BookMeta {
	return make(BookMeta)
}

func (meta *BookMeta) NewColumn(k string) *Column {
	col := NewColumn()
	m := *meta
	if lib := m["library"].Value(); lib != "" {
		col.Field = Lib(lib).DB.GetField(k)
	}
	m[k] = col
	return col
}

func (meta *BookMeta) NewItem(k string) *Item {
	item := NewCategoryItem()
	m := *meta
	if lib := m["library"].Value(); lib != "" {
		item.Field = Lib(lib).DB.GetField(k)
	}
	m[k] = item
	return item
}

func (meta *BookMeta) NewCategory(k string) *Category {
	cat := NewCategory(k)
	m := *meta
	if lib := m["library"].Value(); lib != "" {
		cat.Field = Lib(lib).DB.GetField(k)
	}
	m[k] = cat
	return cat
}

func (c *Category) AddItem() *Item {
	item := NewCategoryItem()
	c.items = append(c.items, item)
	return item
}
