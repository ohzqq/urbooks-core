package book_builder

import (
	"encoding/xml"
	"time"
)

type Book struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Namespaces
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Link    []Link   `xml:"link"`
	Updated TimeStr  `xml:"updated"`
	Author  *Person  `xml:"author"`
	Entry   []*Entry `xml:"entry"`
}

type Namespaces struct {
	DCTerms    string `xml:"xmlns:dcterms,attr"`
	OPF        string `xml:"xmlns:opf,attr"`
	OPDS       string `xml:"xmlns:opds,attr"`
	Schema     string `xml:"xmlns:schema,attr"`
	OpenSearch string `xml:"xmlns:opensearch,attr"`
}

type Entry struct {
	Title       string      `xml:"title"`
	ID          string      `xml:"id"`
	Published   TimeStr     `xml:"published"`
	Updated     TimeStr     `xml:"updated"`
	AltTitle    string      `xml:"schema:alternativeHeadline,omitempty"`
	Series      Series      `xml:"schema:Series,omitempty"`
	Author      []*Person   `xml:"author"`
	Contributor []*Person   `xml:"contributor"`
	Language    []string    `xml:"dcterms:language,omitempty"`
	Category    []*Category `xml:"category,omitempty"`
	Summary     *Text       `xml:"summary"`
	Content     *Text       `xml:"content"`
	Link        []Link      `xml:"link"`
}

type Category struct {
	Term   string `xml:"term,attr"`
	Scheme string `xml:"scheme,attr,omitempty"`
	Label  string `xml:"label,attr"`
}

type Series struct {
	Name     string `xml:"name,attr,omitempty"`
	Position string `xml:"position,attr,omitempty"`
}

type Link struct {
	Rel         string `xml:"rel,attr,omitempty"`
	Href        string `xml:"href,attr"`
	Type        string `xml:"type,attr,omitempty"`
	HrefLang    string `xml:"hreflang,attr,omitempty"`
	Title       string `xml:"title,attr,omitempty"`
	Length      uint   `xml:"length,attr,omitempty"`
	FacetGroup  string `xml:"opds:facetGroup,attr,omitempty"`
	ActiveFacet bool   `xml:"opds:activeFacet,attr,omitempty"`
	Count       uint   `xml:"thr:count,attr,omitempty"`
}

type Person struct {
	Name     string `xml:"name"`
	URI      string `xml:"uri,omitempty"`
	Email    string `xml:"email,omitempty"`
	InnerXML string `xml:",innerxml"`
}

type Text struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

type TimeStr string

func Time(t time.Time) TimeStr {
	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}
