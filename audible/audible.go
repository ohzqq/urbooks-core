package audible

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gosimple/slug"
	"github.com/ohzqq/urbooks-core/book"
)

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

type query struct {
	*url.URL
	suffix      string
	countryCode string
	values      url.Values
	asin        string
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

func NewQuery() *AudibleQuery {
	audible := &AudibleQuery{
		IsApi:   true,
		api:     NewApiRequest(),
		scraper: newScraper(),
		query:   newApiQuery(),
	}
	if audible.IsWeb {
		audible.query = newScraperQuery()
	}
	return audible
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
	return query
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

func (q *query) setValues(val url.Values) *query {
	q.values = val
	return q
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
	var b []*book.Book
	urls := q.scraper.getListURLs(q.Url)
	if q.IsApi {
		for _, u := range urls {
			q.query.asin = getAsin(u)
			b = append(b, q.api.getBook(q.query.string()))
		}
	}
	return b
}

func (q *AudibleQuery) Search() *AudibleQuery {
	q.query.setValues(q.parseCliSearch())
	if q.IsApi {
		results := q.api.searchResults(q.buildUrl())
		fmt.Printf("%+v\n", results)
	}
	return q
}

func (q *AudibleQuery) buildQuery() url.Values {
	var query = url.Values{}

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

	return query
}

func getAsin(path string) string {
	paths := strings.Split(path, "/")
	return paths[len(paths)-1]
}

func (q *AudibleQuery) buildUrl() string {
	if !strings.HasSuffix(q.query.Host, q.query.suffix) {
		q.query.Host = q.query.Host + q.query.suffix
	}

	q.query.RawQuery = q.buildQuery().Encode()

	return q.query.String()
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
