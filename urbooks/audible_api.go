package urbooks

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/ohzqq/urbooks-core/book"
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

func NewAudibleQuery() *AudibleQuery {
	return &AudibleQuery{
		url: &url.URL{
			Scheme: "https",
			Path:   audibleApi,
		},
		client: audibleClient,
	}
}

func (q *AudibleQuery) Get() map[string]json.RawMessage {
	println(q.url.String())
	resp, err := audibleClient.Get(q.url.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, errr := io.ReadAll(resp.Body)
	if errr != nil {
		log.Fatal(errr)
	}

	var result map[string]json.RawMessage
	uerr := json.Unmarshal(body, &result)
	if uerr != nil {
		log.Fatalf("failed to unmarshal audible api search %v\n", uerr)
	}

	return result
}

func (q *AudibleQuery) Product(asin string) *book.Book {
	query := url.Values{}
	query.Set("response_groups", "media,product_desc,contributors,series,product_extended_attrs,product_attrs")
	q.url.RawQuery = query.Encode()
	q.url.Path = path.Join(audibleApi, asin)

	result := q.Get()
	b := book.UnmarshalAudibleApiProduct(result["product"])

	return b
}

func (q *AudibleQuery) Search() []string {
	query := url.Values{}

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

	result := q.Get()

	var total int
	err := json.Unmarshal(result["total_results"], &total)
	if err != nil {
		log.Fatalf("failed to unmarshal total results %v\n", err)
	}

	var products []map[string]string
	err = json.Unmarshal(result["products"], &products)
	if err != nil {
		log.Fatalf("failed to unmarshal products %v\n", err)
	}

	var asin []string
	for _, p := range products {
		asin = append(asin, p["asin"])
	}

	return asin
}
