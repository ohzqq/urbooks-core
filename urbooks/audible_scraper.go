package urbooks

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	//"github.com/ohzqq/urbooks/ui"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	//tea "github.com/charmbracelet/bubbletea"
)

var _ = fmt.Sprintf("%v", "")

const audible = "audible.ca"

type AudibleScraper struct {
	Scraper     *geziyor.Geziyor
	ScraperOpts *geziyor.Options
	URLs        map[string]string
	Books       []BookMeta
	URL         *url.URL
	AudibleURL  string
	Suffix      string
	IsList      bool
	IsSearch    bool
	IsSeries    bool
}

type AudibleSearch struct {
	url   *url.URL
	query url.Values
}

func NewAudibleSearch() *AudibleSearch {
	return &AudibleSearch{
		url: &url.URL{
			Scheme: "https",
			Host:   audible,
			Path:   "/search",
		},
		query: make(url.Values),
	}
}

func (s *AudibleSearch) Keywords(words string) *AudibleSearch {
	s.query.Set("keywords", words)
	return s
}

func (s *AudibleSearch) Authors(words string) *AudibleSearch {
	s.query.Set("searchAuthor", words)
	return s
}

func (s *AudibleSearch) Narrators(words string) *AudibleSearch {
	s.query.Set("searchNarrator", words)
	return s
}

func (s *AudibleSearch) Title(words string) *AudibleSearch {
	s.query.Set("title", words)
	return s
}

func (s *AudibleSearch) String() string {
	s.url.RawQuery = s.query.Encode()
	return s.url.String()
}

func (s *AudibleSearch) Search() *AudibleScraper {
	a := NewAudibleScraper().Get(s.String())
	a.IsSearch = true
	a.URLs = a.getListURLs(a.AudibleURL)
	return a
}

func NewAudibleScraper() *AudibleScraper {
	return &AudibleScraper{
		ScraperOpts: &geziyor.Options{
			ConcurrentRequests: 1,
			LogDisabled:        true,
		},
		URLs: make(map[string]string),
	}
}

func (a *AudibleScraper) Get(audible string) *AudibleScraper {
	a.AudibleURL = audible

	var err error
	a.URL, err = url.Parse(a.AudibleURL)
	if err != nil {
		log.Fatal(err)
	}

	return a
}

func (a *AudibleScraper) Scrape() []BookMeta {
	var urls map[string]string
	audible := a.AudibleURL

	switch strings.Contains(a.URL.Path, "/pd") {
	case true:
		urls = map[string]string{"self": audible}
		a.URLs["self"] = a.AudibleURL
	case false:
		a.IsList = true
		urls = a.getListURLs(audible)
	}

	if a.IsSearch {
		switch len(a.URLs) {
		case 0:
			fmt.Println("No results")
		case 1:
			break
		default:
			//choice := ui.NewPrompt(a.URLs).Choose()
			//a.IsSearch = false
			//urls = map[string]string{"self": choice}
		}
	}

	for _, u := range urls {
		a.ScraperOpts.StartURLs = []string{u}
		a.ScraperOpts.ParseFunc = a.scrapeBook()
		geziyor.NewGeziyor(a.ScraperOpts).Start()
	}
	return a.Books
}

func (a *AudibleScraper) getListURLs(aUrl string) map[string]string {
	urls := make(map[string]string)
	a.ScraperOpts.StartURLs = []string{aUrl}
	a.ScraperOpts.ParseFunc = func(g *geziyor.Geziyor, r *client.Response) {
		metaList := r.HTMLDoc.Find("li.productListItem")
		metaList.Each(func(_ int, s *goquery.Selection) {
			link := s.Find("li.bc-list-item h3.bc-heading a")
			href, _ := link.Attr("href")
			if href != "" {
				pd, err := url.Parse(href)
				if err != nil {
					log.Fatal(err)
				}
				linkURL := url.URL{
					Scheme: a.URL.Scheme,
					Host:   a.URL.Host,
					Path:   pd.Path,
				}
				//a.URLs[link.Text()] = linkURL.String()
				urls[link.Text()] = linkURL.String()
			}
		})
	}
	geziyor.NewGeziyor(a.ScraperOpts).Start()
	return urls
}

func (a *AudibleScraper) scrapeBook() func(g *geziyor.Geziyor, r *client.Response) {
	return func(g *geziyor.Geziyor, r *client.Response) {
		//book := Book{Audiobook: true}
		book := NewBook()

		book.Set("title", NewColumn(strings.TrimSpace(r.HTMLDoc.Find("li.bc-list-item h1.bc-heading").Text())))

		coverURL, _ := r.HTMLDoc.Find(".hero-content img.bc-pub-block").Attr("src")
		book.Set("cover", NewCategoryItem().Set("url", coverURL))

		authors := NewCategory()
		authors.SetField("isNames", "true")
		r.HTMLDoc.Find(".authorLabel a").Each(func(_ int, s *goquery.Selection) {
			authors.AddItem(NewCategoryItem().Set("value", s.Text()))
		})
		book.Set("authors", authors)

		narrators := NewCategory()
		narrators.SetField("isNames", "true")
		r.HTMLDoc.Find(".narratorLabel a").Each(func(_ int, s *goquery.Selection) {
			narrators.AddItem(NewCategoryItem().Set("value", s.Text()))
		})
		book.Set("narrators", narrators)

		series := r.HTMLDoc.Find(".seriesLabel").Text()
		splitSeries := strings.Split(strings.TrimPrefix(strings.TrimSpace(series), "Series:"), ",")
		n := 0
		p := 1
		for i := 0; i < len(splitSeries)/2; i++ {
			s := NewCategoryItem()
			s.Set("name", strings.TrimPrefix(strings.TrimSpace(splitSeries[n]), "Book "))
			s.Set("position", strings.TrimPrefix(strings.TrimSpace(splitSeries[p]), "Book "))
			n = n + 2
			p = p + 2
			book.Set("series", s)
		}

		tags := NewCategory()
		r.HTMLDoc.Find(".bc-chip-text").Each(func(_ int, s *goquery.Selection) {
			tags.AddItem(NewCategoryItem().Set("value", strings.TrimSpace(s.Text())))
		})
		book.Set("tags", tags)

		desc, err := r.HTMLDoc.Find(".productPublisherSummary span.bc-text").Html()
		if err != nil {
			log.Fatal(err)
		}
		book.Set("description", NewColumn(desc))

		fmt.Printf("%+v\n", book)
		a.Books = append(a.Books, book)
	}
}
