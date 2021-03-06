package audible

import (
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/ohzqq/urbooks-core/book"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

const audibleHost = "www.audible"

type WebScraper struct {
	Scraper     *geziyor.Geziyor
	ScraperOpts *geziyor.Options
	url         *url.URL
	URLs        map[string]string
	Books       []*book.Book
	searchQuery url.Values

	AudibleURL string
	Suffix     string
	IsList     bool
	IsSearch   bool
	IsSingle   bool
}

func newScraper() *WebScraper {
	return &WebScraper{
		ScraperOpts: &geziyor.Options{
			ConcurrentRequests: 1,
			LogDisabled:        true,
		},
		URLs:        make(map[string]string),
		searchQuery: make(url.Values),
	}
}

func (a *WebScraper) getBook(u string) *book.Book {
	books := a.scrapeUrls(u)
	if len(books) > 0 {
		return books[0]
	}
	return nil
}

func (a *WebScraper) scrapeUrls(urls ...string) []*book.Book {
	for _, u := range urls {
		a.ScraperOpts.StartURLs = []string{u}
		a.ScraperOpts.ParseFunc = a.scrapeBook()
		geziyor.NewGeziyor(a.ScraperOpts).Start()
	}
	return a.Books
}

func (a *WebScraper) getListURLs(aUrl string) []string {
	var urls []string
	//urls := make(map[string]string)
	a.ScraperOpts.StartURLs = []string{aUrl}
	a.ScraperOpts.ParseFunc = func(g *geziyor.Geziyor, r *client.Response) {
		metaList := r.HTMLDoc.Find("li.productListItem")
		metaList.Each(func(_ int, s *goquery.Selection) {
			link := s.Find("li.bc-list-item h3.bc-heading a")
			//var authors []string
			//s.Find(".authorLabel a").Each(func(_ int, a *goquery.Selection) {
			//  authors = append(authors, a.Text())
			//})
			href, _ := link.Attr("href")
			if href != "" {
				pd, err := url.Parse(href)
				if err != nil {
					log.Fatal(err)
				}
				urls = append(urls, pd.Path)

				//text := fmt.Sprintf("%s by %s", link.Text(), strings.Join(authors, ", "))
				//urls[text] = pd.Path
			}
		})
	}
	geziyor.NewGeziyor(a.ScraperOpts).Start()

	return urls
}

func (a *WebScraper) scrapeBook() func(g *geziyor.Geziyor, r *client.Response) {
	return func(g *geziyor.Geziyor, r *client.Response) {
		b := book.NewBook()

		title := b.GetField("title")
		t := strings.TrimSpace(r.HTMLDoc.Find("li.bc-list-item h1.bc-heading").Text())
		title.SetMeta(t)

		coverURL, _ := r.HTMLDoc.Find(".hero-content img.bc-pub-block").Attr("src")
		b.GetField("cover").Item().Set("url", coverURL)

		if f := b.GetField("authors"); f.IsNull() {
			authors := f.Collection()
			r.HTMLDoc.Find(".authorLabel a").Each(func(_ int, s *goquery.Selection) {
				if text := s.Text(); text != "" {
					authors.AddItem().Set("value", text)
				}
			})
		}

		b.AddField(book.NewCollection("#narrators"))
		narrators := b.GetField("#narrators").SetIsNames().SetIsMultiple().Collection()
		r.HTMLDoc.Find(".narratorLabel a").Each(func(_ int, s *goquery.Selection) {
			if text := s.Text(); text != "" {
				narrators.AddItem().Set("value", text)
			}
		})

		seriesHtml := strings.TrimPrefix(strings.TrimSpace(r.HTMLDoc.Find(".seriesLabel").Text()), "Series:")
		allSeries := regexp.MustCompile(`(\w+\s?){1,}, (Book \d+)`).FindAllString(seriesHtml, -1)
		if len(allSeries) > 0 {
			split := strings.Split(allSeries[0], ", Book ")

			series := b.GetField("series").Item()
			series.Set("value", split[0]).Set("position", split[1])

			position := b.GetField("position")
			position.SetMeta(split[1])
		}

		tags := b.GetField("tags").Collection()
		r.HTMLDoc.Find(".bc-chip-text").Each(func(_ int, s *goquery.Selection) {
			tags.AddItem().Set("value", strings.TrimSpace(s.Text()))
		})

		description := b.GetField("description")
		desc, err := r.HTMLDoc.Find(".productPublisherSummary span.bc-text").Html()
		if err != nil {
			log.Fatal(err)
		}
		description.SetMeta(desc)

		a.Books = append(a.Books, b)
	}
}
