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

const audible = "audible.ca"

type AudibleScraper struct {
	Scraper     *geziyor.Geziyor
	ScraperOpts *geziyor.Options
	URLs        map[string]string
	Books       []*book.Book
	URL         *url.URL
	searchQuery url.Values
	AudibleURL  string
	Suffix      string
	IsList      bool
	IsSearch    bool
	IsSingle    bool
	NoCovers    bool
	Keywords    string
	Authors     string
	Narrators   string
	Title       string
}

func NewAudibleSearch() *AudibleScraper {
	return &AudibleScraper{searchQuery: make(url.Values)}
}

func (s *AudibleScraper) SetKeywords(words string) *AudibleScraper {
	s.Keywords = words
	return s
}

func (s *AudibleScraper) SetAuthors(words string) *AudibleScraper {
	s.Authors = words
	return s
}

func (s *AudibleScraper) SetNarrators(words string) *AudibleScraper {
	s.Narrators = words
	return s
}

func (s *AudibleScraper) SetTitle(words string) *AudibleScraper {
	s.Title = words
	return s
}

func (s *AudibleScraper) SetNoCovers() *AudibleScraper {
	s.NoCovers = true
	return s
}

func (s *AudibleScraper) String() string {
	url := &url.URL{
		Scheme: "https",
		Host:   audible,
		Path:   "/search",
	}

	if s.Keywords != "" {
		s.searchQuery.Set("keywords", s.Keywords)
	}

	if s.Authors != "" {
		s.searchQuery.Set("searchAuthor", s.Authors)
	}

	if s.Narrators != "" {
		s.searchQuery.Set("searchNarrator", s.Narrators)
	}

	if s.Title != "" {
		s.searchQuery.Set("title", s.Title)
	}

	url.RawQuery = s.searchQuery.Encode()
	return url.String()
}

func (s *AudibleScraper) Search() *AudibleScraper {
	a := NewAudibleScraper().Get(s.String())
	if s.NoCovers {
		a.SetNoCovers()
	}
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
		URLs:        make(map[string]string),
		searchQuery: make(url.Values),
	}
}

func (a *AudibleScraper) Get(audible string) *AudibleScraper {
	a.AudibleURL = audible
	a.IsSingle = true
	a.ParseURL()
	return a
}

func (a *AudibleScraper) List(audible string) *AudibleScraper {
	a.AudibleURL = audible
	a.IsList = true
	a.ParseURL()
	return a
}

func (a *AudibleScraper) ParseURL() {
	aUrl, err := url.Parse(a.AudibleURL)
	if err != nil {
		log.Fatal(err)
	}
	a.URL = aUrl
}

func (a *AudibleScraper) Scrape() []*book.Book {
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

func (a *AudibleScraper) getListURLs(aUrl string) map[string]string {
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
					Scheme: a.URL.Scheme,
					Host:   a.URL.Host,
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

func (a *AudibleScraper) scrapeBook() func(g *geziyor.Geziyor, r *client.Response) {
	return func(g *geziyor.Geziyor, r *client.Response) {
		b := book.NewBook()

		var title string
		if f := b.GetField("title"); f.IsNull() {
			title = strings.TrimSpace(r.HTMLDoc.Find("li.bc-list-item h1.bc-heading").Text())
			f.SetData(title)
		}

		coverURL, _ := r.HTMLDoc.Find(".hero-content img.bc-pub-block").Attr("src")
		if !a.NoCovers {
			DownloadCover(title, coverURL)
		}

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
