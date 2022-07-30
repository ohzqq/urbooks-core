package book

import (
	"bytes"
	"encoding/xml"
	"log"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	ITunes  string   `xml:"xmlns:itunes,attr"`
	Version string   `xml:"version,attr"`
	*Channel
}

func NewFeed() *RSS {
	return &RSS{
		ITunes:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
		Version: "2.0",
		Channel: NewChannel(),
	}
}

func (rss *RSS) Marshal() []byte {
	pkg := bytes.NewBufferString(xml.Header)
	enc := xml.NewEncoder(pkg)
	enc.Indent("", "  ")
	err := enc.Encode(rss)
	if err != nil {
		log.Fatal(err)
	}
	return pkg.Bytes()
}

func (r *RSS) SetChannel() *Channel {
	return r.Channel
}

type SharedRss struct {
	Description *Description
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate,omitempty"`
	Image       Image    `xml:"itunes:image"`
	Category    []string `xml:"category,omitempty"`
}

type Channel struct {
	XMLName  xml.Name `xml:"channel"`
	Language string   `xml:"language,omitempty"`
	*SharedRss
	Item []*RssItem
}

func BookToRssChannel(b *Book) *RSS {
	rss := NewFeed()

	channel := rss.SetChannel()
	channel.SetShared(sharedRss(b))
	channel.SetLanguage(b.GetField("languages").String())

	channel.AddItem(BookToRssItem(b))
	return rss
}

func BookToRssItem(b *Book) *RssItem {
	item := NewRssItem()
	item.SetShared(sharedRss(b))
	item.SetGuid(b.GetMeta("uri"))
	item.SetDuration(b.GetMeta("duration"))
	item.SetAuthor(b.GetMeta("authors"))
	item.SetEnclosure(b.GetFile("audio"))
	return item
}

func sharedRss(b *Book) *SharedRss {
	rss := NewRssObject()
	rss.SetTitle(b.GetMeta("title"))
	rss.SetLink(b.GetMeta("uri"))
	rss.SetPubdate(b.GetMeta("published"))
	rss.SetImage(b.GetFile("cover").Get("url"))
	rss.SetDescription(b.GetMeta("description"))

	for _, item := range b.GetField("tags").Collection().EachItem() {
		rss.AddCategory(item.Get("value"))
	}

	return rss
}

type RssItem struct {
	XMLName   xml.Name   `xml:"item"`
	Guid      string     `xml:"guid,omitempty"`
	Author    string     `xml:"author"`
	Duration  string     `xml:"itunes:duration,omitempty"`
	Enclosure *Enclosure `xml:"enclosure"`
	*SharedRss
}

type Description struct {
	XMLName xml.Name `xml:"description,omitempty"`
	Body    string   `xml:",cdata"`
}

func NewChannel() *Channel {
	return &Channel{SharedRss: NewRssObject()}
}

func NewRssItem() *RssItem {
	return &RssItem{SharedRss: NewRssObject()}
}

func NewRssObject() *SharedRss {
	return &SharedRss{}
}

func (rss *SharedRss) SetTitle(title string) *SharedRss {
	rss.Title = title
	return rss
}

func (rss *SharedRss) SetLink(link string) *SharedRss {
	rss.Link = link
	return rss
}

func (rss *SharedRss) SetDescription(desc string) *SharedRss {
	rss.Description = &Description{Body: desc}
	return rss
}

func (rss *SharedRss) AddCategory(category string) *SharedRss {
	rss.Category = append(rss.Category, category)
	return rss
}

func (rss *SharedRss) SetPubdate(pubdate string) *SharedRss {
	rss.PubDate = pubdate
	//rss.PubDate = Time(pubdate)
	return rss
}

func (rss *SharedRss) SetImage(uri string) *SharedRss {
	rss.Image = Image{Href: uri}
	return rss
}

func (rss *SharedRss) Channel() *Channel {
	return &Channel{SharedRss: rss}
}

func (rss *SharedRss) Item() *RssItem {
	return &RssItem{SharedRss: rss}
}

func (c *Channel) SetShared(rss *SharedRss) *Channel {
	c.SharedRss = rss
	return c
}

func (c *Channel) SetLanguage(lang string) *Channel {
	c.Language = lang
	return c
}

func (c *Channel) AddItem(item *RssItem) *Channel {
	c.Item = append(c.Item, item)
	return c
}

type Image struct {
	XMLName xml.Name `xml:"itunes:image"`
	Href    string   `xml:"href,attr"`
}

func (i *RssItem) SetShared(rss *SharedRss) *RssItem {
	i.SharedRss = rss
	return i
}

func (i *RssItem) SetGuid(id string) *RssItem {
	i.Guid = id
	return i
}

func (i *RssItem) SetAuthor(author string) *RssItem {
	i.Author = author
	return i
}

func (i *RssItem) SetDuration(dur string) *RssItem {
	i.Duration = dur
	return i
}

func (i *RssItem) SetEnclosure(file *Item) *RssItem {
	i.Enclosure = &Enclosure{
		Url:    file.Get("url"),
		Length: file.Get("size"),
		Type:   AudioMimeType(file.Get("extension")),
	}
	return i
}

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

func NewEnclosure() *Enclosure {
	return &Enclosure{}
}

func (e *Enclosure) SetUrl(url string) *Enclosure {
	e.Url = url
	return e
}

func (e *Enclosure) SetLength(length string) *Enclosure {
	e.Length = length
	return e
}

func (e *Enclosure) SetType(t string) *Enclosure {
	e.Type = t
	return e
}

type TimeStr string

func Time(t time.Time) TimeStr {
	return TimeStr(t.Format("Monday, 02 January 2006 15:04:05 MST"))
}
