package audible

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/ohzqq/urbooks-core/book"
	"github.com/ohzqq/urbooks-core/bubbles"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

const audible = "www.audible"

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

func NewWebScraper() *AudibleQuery {
	return &AudibleQuery{
		url: &url.URL{
			Scheme: "https",
			Host:   audible,
			Path:   "/search",
		},
		IsWeb:   true,
		scraper: newScraper(),
	}
}

//func (s *WebScraper) Search() *WebScraper {
//  a := NewWebScraper().Get(s.buildUrl())
//  if s.NoCovers {
//    a.SetNoCovers()
//  }
//  a.IsSearch = true
//  a.URLs = a.getListURLs(a.AudibleURL)
//  return a
//}

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

func (a *WebScraper) Get(audible string) *WebScraper {
	a.AudibleURL = audible
	a.IsSingle = true
	//a.ParseURL()
	return a
}

func (a *WebScraper) List(audible string) *WebScraper {
	a.AudibleURL = audible
	a.IsList = true
	//a.ParseURL()
	return a
}

//func (a *WebScraper) ParseURL() {
//  aUrl, err := url.Parse(a.AudibleURL)
//  if err != nil {
//    log.Fatal(err)
//  }
//  a.url = aUrl
//}

func (a *WebScraper) getBook(u string) *book.Book {
	s := a.Get(u)
	books := s.Scrape()
	fmt.Printf("%+v\n", books)
	//books := a.scrapeUrls([]string{u})
	if len(books) > 0 {
		return books[0]
	}
	return nil
}

func (a *WebScraper) scrapeUrls(urls []string) []*book.Book {
	for _, u := range urls {
		a.ScraperOpts.StartURLs = []string{u}
		a.ScraperOpts.ParseFunc = a.scrapeBook()
		geziyor.NewGeziyor(a.ScraperOpts).Start()
	}
	return a.Books
}

func (a *WebScraper) Scrape() []*book.Book {
	var urls map[string]string
	switch {
	case a.IsSingle:
		urls = map[string]string{"self": a.AudibleURL}
	case a.IsSearch:
		switch len(a.URLs) {
		case 0:
			fmt.Println("No search results")
		case 1:
			for _, u := range a.URLs {
				urls = map[string]string{"self": u}
			}
			break
		default:
			choice := bubbles.NewPrompt("search results: pick one", a.URLs).Choose()
			a.IsSearch = false
			urls = map[string]string{"self": choice}
		}
	case a.IsList:
		urls = a.getListURLs(a.AudibleURL)
	}

	for _, u := range urls {
		a.ScraperOpts.StartURLs = []string{u}
		a.ScraperOpts.ParseFunc = a.scrapeBook()
		geziyor.NewGeziyor(a.ScraperOpts).Start()
	}
	return a.Books
}

func (a *WebScraper) getListURLs(aUrl string) map[string]string {
	println("get list urls!")
	urls := make(map[string]string)
	a.ScraperOpts.StartURLs = []string{aUrl}
	a.ScraperOpts.ParseFunc = func(g *geziyor.Geziyor, r *client.Response) {
		metaList := r.HTMLDoc.Find("li.productListItem")
		metaList.Each(func(_ int, s *goquery.Selection) {
			link := s.Find("li.bc-list-item h3.bc-heading a")
			var authors []string
			s.Find(".authorLabel a").Each(func(_ int, a *goquery.Selection) {
				authors = append(authors, a.Text())
			})
			href, _ := link.Attr("href")
			if href != "" {
				pd, err := url.Parse(href)
				if err != nil {
					log.Fatal(err)
				}
				linkURL := url.URL{
					Scheme: "https",
					Host:   a.url.Host,
					Path:   pd.Path,
				}
				text := fmt.Sprintf("%s by %s", link.Text(), strings.Join(authors, ", "))
				urls[text] = linkURL.String()
			}
		})
	}
	geziyor.NewGeziyor(a.ScraperOpts).Start()
	return urls
}

func (a *WebScraper) scrapeBook() func(g *geziyor.Geziyor, r *client.Response) {
	return func(g *geziyor.Geziyor, r *client.Response) {
		b := book.NewBook()

		var title string
		if f := b.GetField("title"); f.IsNull() {
			title = strings.TrimSpace(r.HTMLDoc.Find("li.bc-list-item h1.bc-heading").Text())
			f.SetData(title)
		}

		coverURL, _ := r.HTMLDoc.Find(".hero-content img.bc-pub-block").Attr("src")
		//if !a.NoCovers {
		DownloadCover(title, coverURL)
		//}

		if f := b.GetField("authors"); f.IsNull() {
			authors := f.Collection()
			r.HTMLDoc.Find(".authorLabel a").Each(func(_ int, s *goquery.Selection) {
				if text := s.Text(); text != "" {
					authors.AddItem().Set("value", text)
				}
			})
		}

		if f := b.GetField("#narrators"); f.IsNull() {
			b.AddField(book.NewCollection("#narrators"))
			narrators := b.GetField("#narrators").SetIsNames().SetIsMultiple().Collection()
			r.HTMLDoc.Find(".narratorLabel a").Each(func(_ int, s *goquery.Selection) {
				if text := s.Text(); text != "" {
					narrators.AddItem().Set("value", text)
				}
			})
		}

		seriesHtml := strings.TrimPrefix(strings.TrimSpace(r.HTMLDoc.Find(".seriesLabel").Text()), "Series:")
		allSeries := regexp.MustCompile(`(\w+\s?){1,}, (Book \d+)`).FindAllString(seriesHtml, -1)
		if len(allSeries) > 0 {
			split := strings.Split(allSeries[0], ", Book ")

			series := b.GetField("series").Item()
			series.Set("value", split[0]).Set("position", split[1])

			position := b.GetField("position")
			position.SetData(split[1])
		}

		if f := b.GetField("tags"); f.IsNull() {
			tags := f.Collection()
			r.HTMLDoc.Find(".bc-chip-text").Each(func(_ int, s *goquery.Selection) {
				tags.AddItem().Set("value", strings.TrimSpace(s.Text()))
			})
		}

		if f := b.GetField("description"); f.IsNull() {
			desc, err := r.HTMLDoc.Find(".productPublisherSummary span.bc-text").Html()
			if err != nil {
				log.Fatal(err)
			}
			f.SetData(desc)
		}

		a.Books = append(a.Books, b)
	}
}
