package calibredb

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/slices"
)

type request struct {
	URL          *url.URL
	query        url.Values
	prefix       string
	itemIDs      []int
	pathParams   []string
	library      string
	path         string
	cat          string
	CatLabel     string
	ids          string
	sort         string
	pathID       string
	PathID       string
	queryIDs     string
	Fields       []string
	itemsPerPage int
	currentPage  int
	allItems     bool
	isCustom     bool
	HasFields    bool
	isSorted     bool
	desc         bool
	collection   bool
	bookQuery    bool
	mtx          sync.Mutex
}

func (lib *Lib) newRequest(u string) (*request, error) {
	var req request

	req.mtx.Lock()
	defer req.mtx.Unlock()

	uri, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case strings.HasPrefix(uri.Path, "/opds"):
		req.prefix = "/opds"
		req.path = strings.TrimPrefix(uri.Path, "/opds/")
	case strings.HasPrefix(uri.Path, "/rss"):
		req.prefix = "/rss"
		req.path = strings.TrimPrefix(uri.Path, "/rss/")
	case strings.HasPrefix(uri.Path, "/api"):
		req.prefix = "/api"
		req.path = strings.TrimPrefix(uri.Path, "/api/")
	default:
		req.path = uri.Path
	}

	req.URL = uri
	req.query = req.URL.Query()
	req.library = req.query.Get("library")

	routeRegex := regexp.MustCompile("^/?([a-zA-Z]+)/?([0-9]+)?/?$")
	matches := routeRegex.FindStringSubmatch(req.path)
	if len(matches) == 0 {
		return &req, fmt.Errorf("400 Bad Request:'%v' is not a valid URL", u)
	}
	req.pathParams = matches[1:]

	if matches[1] != "" {
		req.cat = matches[1]
		if !lib.validEndpoint(req.cat) {
			return &req, fmt.Errorf("400 Bad Request:'%v' is not a valid enpoint", req.cat)
		}
		if req.cat == "narrators" {
			req.isCustom = true
		}
	}

	if req.query.Has("ids") {
		req.queryIDs = req.query.Get("ids")
		req.ids = req.query.Get("ids")
	}

	if matches[2] != "" {
		req.ids = matches[2]
		req.pathID = matches[2]
		req.PathID = matches[2]
		if req.cat == "narrators" {
			req.isCustom = true
			req.bookQuery = true
		}
	}

	if req.queryIDs == "all" {
		req.allItems = true
	}

	req.CatLabel = req.cat

	switch req.cat {
	case "customColumns":
		req.collection = false
	case "preferences":
		req.collection = false
	case "books":
		req.collection = true
		req.bookQuery = true
	default:
		if req.pathID != "" {
			req.ids = lib.booksInCatStmt(req.cat, req.pathID)
			req.bookQuery = true
		}
		req.collection = true
	}

	if req.query.Has("sort") {
		req.isSorted = true
		req.sort = req.query.Get("sort")
	}

	if req.query.Has("order") {
		if req.query.Get("order") == "desc" {
			req.desc = true
		}
	}

	if req.query.Has("fields") {
		req.HasFields = true
		req.Fields = strings.Split(req.query.Get("fields"), ",")

		if req.cat != "books" {
			if len(req.Fields) == 1 {
				if slices.Contains(req.Fields, "books") {
					req.ids = lib.booksInCatStmt(req.cat, req.ids)
					req.collection = true
					req.bookQuery = true
				}
			}
		}
	}

	switch ids := req.ids; {
	case ids != "":
		for _, id := range strings.Split(req.ids, ",") {
			newID, err := strconv.Atoi(id)
			if err != nil {
				log.Fatal(err)
			}
			req.itemIDs = append(req.itemIDs, newID)
		}
	}

	if len(req.itemIDs) > 1 {
		req.collection = true
	}

	if req.collection {
		if req.bookQuery {
			if !req.allItems {
				var err error
				req.itemsPerPage, err = strconv.Atoi(req.query.Get("itemsPerPage"))
				if err != nil {
					req.itemsPerPage = 50
					req.query.Set("itemsPerPage", "50")
				}

				req.currentPage, err = strconv.Atoi(req.query.Get("currentPage"))
				if err != nil {
					req.query.Set("currentPage", "1")
					req.currentPage = 1
				}
			}
		}
	}

	//fmt.Printf("request params: %+v\n", req)

	return &req, nil
}

func (r *request) calculateOffset() int {
	return r.itemsPerPage*r.currentPage - r.itemsPerPage
}
