package urbooks

import (
	"encoding/xml"
	"time"

	"github.com/lann/builder"
)

type Package struct {
	XMLName  xml.Name `xml:"http://www.idpf.org/2007/opf package"`
	Version  string   `attr,`
	Metadata Metadata
}

type Metadata struct {
	DC          string `xml:"xmlns:dc,attr"`
	OPF         string `xml:"xmlns:opf,attr"`
	Author      []Author
	Description string `xml:"description"`
	Identifier  []*Identifier
	Language    []string `xml:"languages"`
	Date        string   `xml:"date"`
	Publisher   string
	Subject     []string `xml:"subject"`
	Title       string
	Meta        []*CalibreMeta
}

type Author struct {
	Creator string `xml:"creator"`
	Role    string `xml:"opf:role,attr"`
}

type Identifier struct {
	XMLName xml.Name `xml:"identifier"`
	Id      string   `xml:"id,attr"`
	Scheme  string   `xml:"opf:scheme,attr"`
}

type CalibreMeta struct {
	XMLName xml.Name `xml:"meta"`
	Name    string   `xml:"name"`
	Content string   `xml:"content"`
}

type opfBuilder builder.Builder

func NewOpfPackage() opfBuilder {
	return OpfBuilder
}

func (e opfBuilder) Title(title string) opfBuilder {
	return builder.Set(e, "Title", title).(opfBuilder)
}

func (e opfBuilder) Series(name string) opfBuilder {
	return builder.Append(e, "Meta", CalibreMeta{Name: "calibre:series", Content: name}).(opfBuilder)
}

func (e opfBuilder) SeriesIndex(pos string) opfBuilder {
	return builder.Append(e, "Meta", CalibreMeta{Name: "calibre:series_index", Content: pos}).(opfBuilder)
}

func (e opfBuilder) AddCustomColumn(name, val string) opfBuilder {
	return builder.Append(e, "Meta", CalibreMeta{Name: "calibre:#" + name, Content: val}).(opfBuilder)
}

func (e opfBuilder) AddLanguage(lang string) opfBuilder {
	return builder.Append(e, "Language", lang).(opfBuilder)
}

func (e opfBuilder) AddAuthor(author string) opfBuilder {
	return builder.Append(e, "Author", Author{Creator: author, Role: "aut"}).(opfBuilder)
}

func (e opfBuilder) AddSubject(subject string) opfBuilder {
	return builder.Append(e, "Subject", subject).(opfBuilder)
}

func (e opfBuilder) AddIdentifier(id, scheme string) opfBuilder {
	return builder.Append(e, "Identifier", Identifier{Id: id, Scheme: scheme}).(opfBuilder)
}

func (e opfBuilder) Date(published time.Time) opfBuilder {
	return builder.Set(e, "Date", Time(published)).(opfBuilder)
}

func (e opfBuilder) Description(summary string) opfBuilder {
	return builder.Set(e, "Description", summary).(opfBuilder)
}

func (e opfBuilder) BuildPackage() Package {
	return Package{
		Version:  "2.0",
		Metadata: builder.GetStruct(e).(Metadata),
	}
}

// OpfBuilder is a fluent immutable builder to build OPDS entries
var OpfBuilder = builder.Register(opfBuilder{}, Metadata{}).(opfBuilder)

type TimeStr string

func Time(t time.Time) TimeStr {
	return TimeStr(t.Format("Monday, 02 January 2006 15:04:05 MST"))
}
