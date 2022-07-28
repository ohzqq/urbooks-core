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
	url *url.URL
	//url         *queryUrl
	countryCode string
	suffix      string
	api         *ApiRequest
	scraper     *WebScraper
	isApi       bool
	isWeb       bool
	Authors     string
	Keywords    string
	Narrators   string
	Title       string
	NoCovers    bool
}

type queryUrl struct {
	*url.URL
	countryCode string
	suffix      string
	values      url.Values
	asin        string
}

func (q *queryUrl) string() string {
	if !strings.HasSuffix(q.Host, q.suffix) {
		q.Host = q.Host + q.suffix
	}

	q.RawQuery = q.values.Encode()

	return q.String()
}

func NewAudibleQuery() *AudibleQuery {
	return &AudibleQuery{
		url: &url.URL{
			Scheme: "https",
		},
	}
}

func (q *AudibleQuery) Get(u string) *book.Book {
	var b *book.Book

	req := q.parseUrl(u)
	if q.isApi {
		b = q.api.getBook(req)
	}

	if q.isWeb {
		b = q.scraper.getBook(u)
	}

	return b
}

func (q *AudibleQuery) Search() *AudibleQuery {
	if q.isApi {
		results := q.api.searchResults(q.buildUrl())
		fmt.Printf("%+v\n", results)
	}
	return q
}

func (q *AudibleQuery) parseUrl(u string) *queryUrl {
	aUrl, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	query := &queryUrl{
		URL:    aUrl,
		values: url.Values{},
	}

	paths := strings.Split(query.Path, "/")
	query.asin = paths[len(paths)-1]

	host := strings.Split(aUrl.Host, ".")
	query.countryCode = host[len(host)-1]
	query.suffix = countrySuffix(query.countryCode)

	return query
}

func (q *AudibleQuery) buildUrl() string {
	if !strings.HasSuffix(q.url.Host, q.suffix) {
		q.url.Host = q.url.Host + q.suffix
	}

	q.url.RawQuery = q.buildQuery().Encode()

	return q.url.String()
}

func (q *AudibleQuery) buildQuery() url.Values {
	var query = url.Values{}

	if a := q.Authors; a != "" {
		if q.isApi {
			query.Set("author", a)
		}
		if q.isWeb {
			query.Set("searchAuthor", a)
		}
	}

	if n := q.Narrators; n != "" {
		if q.isApi {
			query.Set("narrator", n)
		}
		if q.isWeb {
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

func (a *AudibleQuery) SetKeywords(words string) *AudibleQuery {
	a.Keywords = words
	return a
}

func (a *AudibleQuery) SetAuthors(words string) *AudibleQuery {
	a.Authors = words
	return a
}

func (a *AudibleQuery) SetNarrators(words string) *AudibleQuery {
	a.Narrators = words
	return a
}

func (a *AudibleQuery) SetTitle(words string) *AudibleQuery {
	a.Title = words
	return a
}

func (a *AudibleQuery) SetNoCovers() *AudibleQuery {
	a.NoCovers = true
	return a
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
