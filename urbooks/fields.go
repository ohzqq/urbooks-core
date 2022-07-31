package urbooks

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/ohzqq/urbooks-core/book"
)

func (l *Library) ConvertBook(b *book.Book, fmt string) book.Fmt {
	return b.ConvertTo(fmt)
}

func (l *Library) ToOPF(b *book.Book) *bytes.Buffer {
	opf := b.ConvertToOPF()
	for name, field := range l.CustomColumns {
		println(name)
		field.Value = b.GetField(name).String()
		val := field.Marshal()
		opf.AddCustomColumn(name, string(val))
	}
	return opf.Marshal()
}

type Fields map[string]*Field

type Field struct {
	CategorySort string   `json:"category_sort"`
	Colnum       int      `json:"colnum"`
	Column       string   `json:"column"`
	Datatype     string   `json:"datatype"`
	Display      Display  `json:"display"`
	IsCategory   bool     `json:"is_category"`
	IsCustom     bool     `json:"is_custom"`
	IsCsp        bool     `json:"is_csp"`
	IsEditable   bool     `json:"is_editable"`
	Multiple     Multiple `json:"is_multiple"`
	Kind         string   `json:"kind"`
	CalibreLabel string   `json:"label"`
	LinkColumn   string   `json:"link_column"`
	Name         string   `json:"name"`
	RecIndex     int      `json:"rec_index"`
	SearchTerms  []string `json:"search_terms"`
	Table        string   `json:"table"`
	Value        any      `json:"#value#"`
	Extra        any      `json:"#extra#"`
}

type Display struct {
	Description     string `json:"description"`
	HeadingPosition string `json:"heading_position"`
	InterpretAs     string `json:"long-text"`
	IsNames         bool   `json:"is_names"`
}

type Multiple struct {
	CacheToList string `json:"cache_to_list"`
	ListToUi    string `json:"list_to_ui"`
	UiToList    string `json:"ui_to_list"`
}

func (f *Field) Marshal() []byte {
	val, err := json.Marshal(f)
	if err != nil {
		log.Fatalf("failed to marshal field %v\n", err)
	}
	return val
}
