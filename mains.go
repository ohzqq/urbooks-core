package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ohzqq/urbooks-core/urbooks"

	"github.com/integrii/flaggy"
	"github.com/spf13/viper"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Config()
	urbooks.InitConfig(viper.GetStringMapString("library_options"))

	start := flaggy.NewSubcommand("start")
	flaggy.AttachSubcommand(start, 1)

	scrape := flaggy.NewSubcommand("scrape")
	var (
		audibleURL string
		keywords   string
		authors    string
		narrators  string
		title      string
	)
	scrape.String(&audibleURL, "u", "url", "audible url")
	scrape.String(&authors, "a", "authors", "search authors on audible")
	scrape.String(&narrators, "n", "narrators", "search narrators on audible")
	scrape.String(&title, "t", "title", "search titles on audible")
	scrape.AddPositionalValue(&keywords, "Keywords", 1, false, "keywords for searching audible")
	flaggy.AttachSubcommand(scrape, 1)

	flaggy.Parse()

	if scrape.Used {
		var scraper *urbooks.AudibleScraper
		switch audibleURL {
		case "":
			search := urbooks.NewAudibleSearch()
			if keywords != "" {
				search.Keywords(keywords)
			}
			if authors != "" {
				search.Authors(authors)
			}
			if narrators != "" {
				search.Narrators(narrators)
			}
			if title != "" {
				search.Title(title)
			}
			scraper = search.Search()
		default:
			scraper = urbooks.NewAudibleScraper().Get(audibleURL)
		}
		books := scraper.Scrape()
		if books.Books == nil {
			fmt.Println("no results")
		}
		fmt.Printf("%+V\n", books)
	}

	if start.Used {
		urbooks.Cfg().Opts["serve"] = "false"
		urbooks.InitLibraries(viper.Sub("libraries"), viper.GetStringMapString("libraries"), false)
		//fmt.Printf("%+v\n", urbooks.Lib("test-library"))
		//fmt.Printf("%V\n", urbooks.NewURL().String())
		//fmt.Println(ui.TermHeight())
		someBooks()
		//url := "preferences?library=audiobooks"
		//books := resp.Books()
		//req := urbooks.ResponseToRequest(&books.Params).Page(10)
		//urbooks.ParseURL(resp.Request.URL().String())
		//fmt.Printf("%+v\n", books)
		//fmt.Printf("%V\n", resp.Request.URL().PathPrefix("opds").String())
		//feed := urbooks.NewNavFeed(urbooks.Lib("audiobooks"))
		//feed.ToXML()
		//fmt.Printf("%V", feed)

		//u := urbooks.NewURL().SetParams(&books.Params)
		//u.PathPrefix("poot")
		//fmt.Printf("%V\n", u.String())
	}
}

var url = "books?library=audiobooks&order=desc&sort=added"

func someBooks() {
	//u := "http://localhost:9932/tags/11?currentPage=2&itemsPerPage=50&library=audiobooks&order=desc&sort=added"
	req := urbooks.NewRequest("test-library")
	req.From("books")
	//req.ID("98")
	//req.Sort("formats")
	//req.Desc()
	//req.Limit("1")
	//req.Page("6")
	//req.Fields("added")
	//req.Find("all")
	//req.Response()
	resp := req.Response()
	//fmt.Printf("%V\n\n", req.String())

	//resp, err := urbooks.Get(u)

	//if err != nil {
	//log.Fatal(err)
	//}
	//cat := urbooks.ParseCat(resp)
	//for _, c := range cat.Items() {
	//  fmt.Printf("%V\n", c.Get("books"))
	//}
	//url := authors
	//if err != nil {
	//log.Fatal(err)
	//}

	//book := urbooks.ParseBooks(resp).Books()[0]
	//fmt.Printf("%+V\n", book.Get("formats").String())
	//fmt.Printf("%+V\n", book.Get("authors").String())
	//fmt.Printf("%+V\n", book.Get("seriesAndTitle").String())
	//fmt.Printf("%+V\n", book.Get("tags").String())
	//fmt.Printf("%+V\n", book.Get("authors").String())

	//for _, b := range book.Get("tags").Data() {
	//fmt.Printf("%+V\n",b.URL() )
	//}
	//return urbooks.ParseBooks(resp).Data[0]
	fmt.Printf("%V\n", string(resp))

}

func pref() {
	//fmt.Printf("%v", db.LibCfg.Cur().EditableFields())
}

func Config() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	//viper.AddConfigPath(filepath.Join(home, ".config/urbooks"))
	viper.AddConfigPath(filepath.Join(home, "Code/urbooks-core/tmp/"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
}
