package audible

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/ohzqq/urbooks-core/book"
)

const (
	apiHost = `api.audible`
	apiPath = `/1.0/catalog/products`
)

type ApiRequest struct {
	client *http.Client
	url    *url.URL
	asin   []string
}

var audibleClient = &http.Client{}

func NewRequest() *AudibleQuery {
	return &AudibleQuery{
		url: &url.URL{
			Scheme: "https",
			Host:   apiHost,
			Path:   apiPath,
		},
		query: &query{
			suffix: ".com",
		},
		IsApi: true,
		api: &ApiRequest{
			client: audibleClient,
		},
	}
}

func NewApiRequest() *ApiRequest {
	return &ApiRequest{
		client: audibleClient,
		url: &url.URL{
			Scheme: "https",
			Host:   apiHost,
			Path:   apiPath,
		},
	}
}

func newRequestUrl() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   apiHost,
		Path:   apiPath,
	}
}

func (a *ApiRequest) makeRequest(u string) map[string]json.RawMessage {
	resp, err := audibleClient.Get(u)
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

func (a *ApiRequest) getBook(req *query) *book.Book {
	req.values.Set("response_groups", "media,product_desc,contributors,series,product_extended_attrs,product_attrs")

	req.Path = path.Join(apiPath, req.asin)
	req.Host = apiHost

	result := a.makeRequest(req.string())
	return book.UnmarshalAudibleApiProduct(result["product"])
}

//func (a *ApiRequest) search(req string) []*book.Book {
//}

func (a *ApiRequest) searchResults(req string) []string {
	result := a.makeRequest(req)

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
