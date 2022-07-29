package audible

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"github.com/ohzqq/urbooks-core/book"
)

type AudibleQuery struct {
	cliArgs
	url     *url.URL
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
	NoCovers  bool
	IsBatch   bool
}

type query struct {
	*url.URL
	suffix      string
	countryCode string
	values      url.Values
	asin        string
}

func (q *query) string() string {
	if !strings.HasSuffix(q.Host, q.suffix) {
		q.Host = q.Host + q.suffix
	}

	q.RawQuery = q.values.Encode()

	return q.String()
}

func NewQuery() *AudibleQuery {
	return &AudibleQuery{
		IsApi:   true,
		api:     NewApiRequest(),
		scraper: newScraper(),
		url: &url.URL{
			Scheme: "https",
		},
		query: &query{
			values: url.Values{},
			suffix: ".com",
		},
	}
}

func (q *AudibleQuery) GetBookMeta() *book.Book {
	q.ParseArgs()

	var b *book.Book

	if q.IsApi {
		b = q.api.getBook(q.query)
	}

	if q.IsWeb {
		b = q.scraper.getBook(q.query.String())
	}

	return b
}

func (q *AudibleQuery) GetBookBatch() []*book.Book {
	var b []*book.Book
	if q.IsWeb {
		urls := q.scraper.getListURLs(q.Url)
		for t, u := range urls {
			fmt.Printf("title: %v, url: %v\n", t, u)
		}
	}
	return b
}

func (q *AudibleQuery) Search() *AudibleQuery {
	if q.IsApi {
		results := q.api.searchResults(q.buildUrl())
		fmt.Printf("%+v\n", results)
	}
	return q
}

func (q *AudibleQuery) ParseArgs() *AudibleQuery {
	if q.Url != "" {
		q.query = q.parseUrl(q.Url)
		return q
	}
	return q
}

func (q *AudibleQuery) buildQuery() url.Values {
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

func (q *AudibleQuery) parseUrl(u string) *query {
	aUrl, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	q.query.URL = aUrl

	if !q.IsBatch {
		paths := strings.Split(q.query.Path, "/")
		q.query.asin = paths[len(paths)-1]
	}

	host := strings.Split(aUrl.Host, ".")
	q.query.countryCode = host[len(host)-1]
	q.query.suffix = countrySuffix(q.query.countryCode)

	return q.query
}

func (q *AudibleQuery) buildSearchUrl() string {
	if !strings.HasSuffix(q.url.Host, q.query.suffix) {
		q.url.Host = q.url.Host + q.query.suffix
	}

	if a := q.Authors; a != "" {
		if q.IsApi {
			q.query.values.Set("author", a)
		}
		if q.IsWeb {
			q.query.values.Set("searchAuthor", a)
		}
	}

	if n := q.Narrators; n != "" {
		if q.IsApi {
			q.query.values.Set("narrator", n)
		}
		if q.IsWeb {
			q.query.values.Set("searchNarrator", n)
		}
	}

	if k := q.Keywords; k != "" {
		q.query.values.Set("keywords", k)
	}

	if t := q.Title; t != "" {
		q.query.values.Set("title", t)
	}

	q.query.URL.RawQuery = q.query.values.Encode()

	return q.url.String()
}

func (q *AudibleQuery) buildUrl() string {
	if !strings.HasSuffix(q.url.Host, q.query.suffix) {
		q.url.Host = q.url.Host + q.query.suffix
	}

	q.url.RawQuery = q.buildQuery().Encode()

	return q.url.String()
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

func (args *cliArgs) SetNoCovers() *cliArgs {
	args.NoCovers = true
	return args
}

func DownloadCover(name, u string) {
	response, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatal(response.StatusCode)
	}

	file, err := os.Create(slug.Make(name) + ".jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}

var countryCodes = map[string]string{
	"us": ".com",
	"ca": ".ca",
	"uk": ".co.uk",
	"au": ".co.uk",
	"fr": "fr",
}

func countrySuffix(code string) string {
	switch code {
	case "uk", "au", "jp", "in":
		return ".co." + code
	default:
		return "." + code
	}
}

func countryCode(suffix string) string {
	switch suffix {
	case "com":
		return "us"
	default:
		return suffix
	}
}
