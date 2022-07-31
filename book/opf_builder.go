package book

import (
	"bytes"
	"encoding/xml"
	"log"

	"golang.org/x/exp/slices"
)

type OPFpackage struct {
	XMLName  xml.Name     `xml:"http://www.idpf.org/2007/opf package"`
	Version  string       `xml:"version,attr"`
	Metadata *OPFmetadata `xml:"metadata"`
}

type OPFmetadata struct {
	DC          string           `xml:"xmlns:dc,attr"`
	OPF         string           `xml:"xmlns:opf,attr"`
	Creator     []OPFcreator     `xml:"dc:creator,omitempty"`
	Description string           `xml:"dc:description,omitempty"`
	Identifier  []OPFIdentifier  `xml:"dc:identifier,omitempty"`
	Language    []string         `xml:"dc:languages,omitempty"`
	Date        string           `xml:"dc:date,omitempty"`
	Publisher   string           `xml:"dc:publisher,omitempty"`
	Subject     []string         `xml:"dc:subject,omitempty"`
	Title       string           `xml:"dc:title"`
	Meta        []OPFcalibreMeta `xml:"meta"`
}

type OPFcreator struct {
	Creator string `xml:",chardata"`
	Role    string `xml:"opf:role,attr"`
}

type OPFIdentifier struct {
	Id     string `xml:",chardata"`
	Scheme string `xml:"opf:scheme,attr"`
}

type OPFcalibreMeta struct {
	Name    string `xml:"name,attr"`
	Content any    `xml:"content,attr"`
}

func NewOpfMetadata() *OPFmetadata {
	return &OPFmetadata{
		DC:  "http://purl.org/dc/terms/",
		OPF: "http://www.idpf.org/2007/opf",
	}
}

func opfFields(b *Book) []*Field {
	opfFields := []string{
		"authors",
		"tags",
		"languages",
		"identifiers",
		"title",
		"published",
		"description",
		"series",
		"position",
	}

	var fields []*Field
	for _, f := range b.EachField() {
		if slices.Contains(opfFields, f.JsonLabel) {
			fields = append(fields, f)
		}
	}

	return fields
}

func (b *Book) ConvertToOPF() *OPFmetadata {
	return buildOPF(b)
}

func buildOPF(b *Book) *OPFmetadata {
	opf := NewOpfMetadata()
	for _, field := range opfFields(b) {
		if !field.IsNull() {
			switch {
			case field.IsCollection():
				for _, item := range field.Collection().EachItem() {
					switch field.JsonLabel {
					case "authors":
						opf.AddAuthor(item.String(field))
					case "tags":
						opf.AddSubject(item.String(field))
					case "languages":
						opf.AddLanguage(item.String(field))
					case "identifiers":
						opf.AddIdentifier(item.Get("value"), item.Get("type"))
					}
				}
			default:
				switch field.JsonLabel {
				case "series":
					opf.AddMeta("series", field.String())
				case "position":
					opf.AddMeta("series_index", field.String())
				case "title":
					opf.SetTitle(field.String())
				case "rating":
					opf.AddMeta("rating", field.String())
				case "published":
					opf.SetDate(field.String())
				case "description":
					opf.SetDescription(field.String())
				}
			}
		}
	}
	return opf
}

func (opf *OPFmetadata) Marshal() *bytes.Buffer {
	pkg := bytes.NewBufferString(xml.Header)
	enc := xml.NewEncoder(pkg)
	enc.Indent("", "  ")
	err := enc.Encode(opf.BuildOPFpackage())
	if err != nil {
		log.Fatal(err)
	}
	return pkg
}

func (m *OPFmetadata) SetTitle(title string) *OPFmetadata {
	m.Title = title
	return m
}

func (m *OPFmetadata) SetPublisher(publisher string) *OPFmetadata {
	m.Publisher = publisher
	return m
}

func (m *OPFmetadata) AddMeta(name string, content any) *OPFmetadata {
	m.Meta = append(m.Meta, OPFcalibreMeta{Name: "calibre:" + name, Content: content})
	return m
}

func (m *OPFmetadata) SetSeries(name string) *OPFmetadata {
	m.AddMeta("calibre:series", name)
	return m
}

func (m *OPFmetadata) SetRating(rating string) *OPFmetadata {
	m.AddMeta("calibre:rating", rating)
	return m
}

func (m *OPFmetadata) SetSeriesIndex(pos string) *OPFmetadata {
	m.AddMeta("calibre:series_index", pos)
	return m
}

func (m *OPFmetadata) AddCustomColumn(name string, val any) *OPFmetadata {
	m.AddMeta("user_metadata:"+name, val)
	return m
}

func (m *OPFmetadata) AddLanguage(lang string) *OPFmetadata {
	m.Language = append(m.Language, lang)
	return m
}

func (m *OPFmetadata) AddAuthor(author string) *OPFmetadata {
	m.Creator = append(m.Creator, OPFcreator{Creator: author, Role: "aut"})
	return m
}

func (m *OPFmetadata) AddSubject(subject string) *OPFmetadata {
	m.Subject = append(m.Subject, subject)
	return m
}

func (m *OPFmetadata) AddIdentifier(id, scheme string) *OPFmetadata {
	m.Identifier = append(m.Identifier, OPFIdentifier{Id: id, Scheme: scheme})
	return m
}

func (m *OPFmetadata) SetDate(published string) *OPFmetadata {
	m.Date = published
	return m
}

func (m *OPFmetadata) SetDescription(summary string) *OPFmetadata {
	m.Description = summary
	return m
}

func (m *OPFmetadata) BuildOPFpackage() OPFpackage {
	return OPFpackage{
		Version:  "2.0",
		Metadata: m,
	}
}
