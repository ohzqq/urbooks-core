package audible

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/ohzqq/urbooks-core/book"
	"github.com/ohzqq/urbooks-core/bubbles"
)

type query struct {
	*url.URL
	suffix      string
	countryCode string
	values      url.Values
	asin        string
	IsApi       bool
	IsWeb       bool
}

func newQuery() *query {
	return &query{
		URL: &url.URL{
			Scheme: "https",
		},
		values: url.Values{},
		suffix: ".com",
	}
}

func newApiQuery() *query {
	query := newQuery()
	query.Host = apiHost
	query.Path = apiPath
	return query
}

func newScraperQuery() *query {
	query := newQuery()
	query.Host = audibleHost
	query.Path = "/search"
	query.suffix = ".ca"
	return query
}

func (q *query) string() string {
	if !strings.HasSuffix(q.Host, q.suffix) {
		q.Host = q.Host + q.suffix
	}

	if q.asin != "" {
		q.values.Set("response_groups", responseGroups)
		q.Path = path.Join(apiPath, q.asin)
	}

	q.RawQuery = q.values.Encode()

	return q.String()
}

type AudibleQuery struct {
	cliArgs
	query   *query
	api     *ApiRequest
	scraper *WebScraper
	IsApi   bool
	IsWeb   bool
}

type cliArgs struct {
	Url       string
	Authors   string
	Keywords  string
	Narrators string
	Title     string
	IsBatch   bool
}

func NewQuery() *AudibleQuery {
	audible := &AudibleQuery{
		IsApi:   true,
		api:     NewApiRequest(),
		scraper: newScraper(),
		query:   newApiQuery(),
	}
	return audible
}

func (q *AudibleQuery) GetBook() *book.Book {
	q.parseCliUrl()

	var b *book.Book

	if q.IsApi {
		b = q.api.getBook(q.query.string())
	}

	if q.IsWeb {
		b = q.scraper.getBook(q.Url)
	}

	return b
}

func (q *AudibleQuery) GetBookBatch() []*book.Book {
	q.parseCliUrl()

	urls := q.scraper.getListURLs(q.Url)

	var b []*book.Book

	if q.IsApi {
		for _, u := range urls {
			q.query.asin = getAsin(u)
			b = append(b, q.api.getBook(q.query.string()))
		}
	}

	if q.IsWeb {
		var webUrls []string
		for _, u := range urls {
			q.query = newScraperQuery()
			q.query.Path = u
			webUrls = append(webUrls, q.query.string())
		}
		return q.scraper.scrapeUrls(webUrls...)
	}
	return b
}

func (q *AudibleQuery) Search() []*book.Book {
	var b []*book.Book

	switch {
	case q.IsWeb:
		q.query = newScraperQuery()
		q.query.values = q.parseCliSearch()
		var urls []string
		for _, u := range q.scraper.getListURLs(q.query.string()) {
			q.query.values = url.Values{}
			q.query.Path = u
			println(q.query.string())
			urls = append(urls, q.query.string())
		}
		b = q.scraper.scrapeUrls(urls...)
	case q.IsApi:
		q.query.values = q.parseCliSearch()
		results := q.api.searchResults(q.query.string())
		for _, result := range results {
			q.query.asin = result
			b = append(b, q.api.getBook(q.query.string()))
		}
	}

	if len(b) > 1 {
		b = q.selectResults(b)
	}

	return b
}

func (q *AudibleQuery) selectResults(books []*book.Book) []*book.Book {
	var selected []*book.Book
	choices := make(map[string]string)
	idx := make(map[string]*book.Book)

	for i, c := range books {
		title := c.GetField("title").String()
		authors := c.GetField("authors").String()
		text := fmt.Sprintf("%s by %s", title, authors)
		num := strconv.Itoa(i)
		choices[text] = num
		idx[num] = c
	}

	choice := bubbles.NewPrompt("search results: pick one", choices).Choose()

	selected = append(selected, idx[choice])
	return selected
}

func (q *AudibleQuery) parseCliSearch() url.Values {
	var query = url.Values{}

	if a := q.Authors; a != "" {
		if q.IsApi {
			query.Set("author", a)
		}
		if q.IsWeb {
			query.Set("searchAuthor", a)
		}
	}

	if n := q.Narrators; n != "" {
		if q.IsApi {
			query.Set("narrator", n)
		}
		if q.IsWeb {
			query.Set("searchNarrator", n)
		}
	}

	if k := q.Keywords; k != "" {
		query.Set("keywords", k)
	}

	if t := q.Title; t != "" {
		query.Set("title", t)
	}

	return query
}

func (q *AudibleQuery) parseCliUrl() *AudibleQuery {
	aURL, err := url.Parse(q.Url)
	if err != nil {
		log.Fatal(err)
	}

	if !q.IsBatch {
		q.query.asin = getAsin(aURL.Path)
	}

	host := strings.Split(aURL.Host, ".")
	q.query.countryCode = host[len(host)-1]
	q.query.suffix = countrySuffix(q.query.countryCode)

	return q
}

func (args *cliArgs) SetKeywords(words []string) *cliArgs {
	args.Keywords = strings.Join(words, " ")
	return args
}

func (args *cliArgs) SetUrl(u string) *cliArgs {
	args.Url = u
	return args
}

func (args *cliArgs) SetAuthors(words string) *cliArgs {
	args.Authors = words
	return args
}

func (args *cliArgs) SetNarrators(words string) *cliArgs {
	args.Narrators = words
	return args
}

func (args *cliArgs) SetTitle(words string) *cliArgs {
	args.Title = words
	return args
}
