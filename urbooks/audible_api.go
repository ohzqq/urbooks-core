package urbooks

import (
	"net/http"
	"net/url"
	"path"
)

const audibleApi = `api.audible.com/1.0/catalog/products`

type AudibleQuery struct {
	client    *http.Client
	url       *url.URL
	Authors   string
	Keywords  string
	Narrators string
	Title     string
}

var audibleClient = &http.Client{}

func NewAudibleClient() *http.Client {
	return &http.Client{}
}

func NewQuery() *AudibleQuery {
	return &AudibleQuery{
		url: &url.URL{
			Scheme: "https",
			Path:   audibleApi,
		},
		client: audibleClient,
	}
}

func (q *AudibleQuery) searchQuery(asin string) *AudibleQuery {
	query := url.Values{}
	query.Set("response_groups", "product_desc,contributors,series,product_extended_attrs,product_attrs")

	if a := q.Authors; a != "" {
		query.Set("author", a)
	}

	if n := q.Narrators; n != "" {
		query.Set("narrator", n)
	}

	if k := q.Keywords; k != "" {
		query.Set("keywords", k)
	}

	if t := q.Title; t != "" {
		query.Set("title", t)
	}

	q.url.RawQuery = query.Encode()
	q.url.Path = path.Join(q.url.Path, asin)

	return q
}
