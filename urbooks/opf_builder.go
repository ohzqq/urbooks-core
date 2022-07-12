package urbooks

import (
	"encoding/xml"
)

type Package struct {
	XMLName  xml.Name  `xml:"http://www.idpf.org/2007/opf package"`
	Version  string    `xml:"version,attr"`
	Metadata *Metadata `xml:"metadata"`
}

type Metadata struct {
	DC          string        `xml:"xmlns:dc,attr"`
	OPF         string        `xml:"xmlns:opf,attr"`
	Creator     []Creator     `xml:"dc:creator"`
	Description string        `xml:"dc:description,omitempty"`
	Identifier  []Identifier  `xml:"dc:identifier"`
	Language    []string      `xml:"dc:languages"`
	Date        string        `xml:"dc:date"`
	Publisher   string        `xml:"dc:publisher,omitempty"`
	Subject     []string      `xml:"dc:subject,omitempty"`
	Title       string        `xml:"dc:title"`
	Meta        []CalibreMeta `xml:"meta"`
}

type Creator struct {
	Creator string `xml:",chardata"`
	Role    string `xml:"opf:role,attr"`
}

type Identifier struct {
	Id     string `xml:",chardata"`
	Scheme string `xml:"opf:scheme,attr"`
}

type CalibreMeta struct {
	Name    string `xml:"name,attr"`
	Content string `xml:"content,attr"`
}

func NewOpfMetadata() *Metadata {
	return &Metadata{
		DC:  "http://purl.org/dc/terms/",
		OPF: "http://www.idpf.org/2007/opf",
	}
}

func (m *Metadata) SetTitle(title string) *Metadata {
	m.Title = title
	return m
}

func (m *Metadata) SetPublisher(publisher string) *Metadata {
	m.Publisher = publisher
	return m
}

func (m *Metadata) AddMeta(name, content string) *Metadata {
	m.Meta = append(m.Meta, CalibreMeta{Name: name, Content: content})
	return m
}

func (m *Metadata) SetSeries(name string) *Metadata {
	m.AddMeta("calibre:series", name)
	return m
}

func (m *Metadata) SetSeriesIndex(pos string) *Metadata {
	m.AddMeta("calibre:series_index", pos)
	return m
}

func (m *Metadata) AddCustomColumn(name, val string) *Metadata {
	m.AddMeta("calibre:#"+name, val)
	return m
}

func (m *Metadata) AddLanguage(lang string) *Metadata {
	m.Language = append(m.Language, lang)
	return m
}

func (m *Metadata) AddAuthor(author string) *Metadata {
	m.Creator = append(m.Creator, Creator{Creator: author, Role: "aut"})
	return m
}

func (m *Metadata) AddSubject(subject string) *Metadata {
	m.Subject = append(m.Subject, subject)
	return m
}

func (m *Metadata) AddIdentifier(id, scheme string) *Metadata {
	m.Identifier = append(m.Identifier, Identifier{Id: id, Scheme: scheme})
	return m
}

func (m *Metadata) SetDate(published string) *Metadata {
	m.Date = published
	return m
}

func (m *Metadata) SetDescription(summary string) *Metadata {
	m.Description = summary
	return m
}

func (m *Metadata) BuildPackage() Package {
	return Package{
		Version:  "2.0",
		Metadata: m,
	}
}
