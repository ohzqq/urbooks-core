package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	//"strings"
	//"encoding/json"
	//"io"

	//"github.com/ohzqq/urbooks/web"
	//"github.com/ohzqq/urbooks/bubbles"
	//"github.com/ohzqq/urbooks/ui"
	"github.com/ohzqq/urbooks/ui/utils"
	//"github.com/ohzqq/urbooks/calibredb"
	"github.com/ohzqq/urbooks/urbooks"

	//tea "github.com/charmbracelet/bubbletea"
	"github.com/integrii/flaggy"
	"github.com/spf13/viper"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Config()
	urbooks.InitConfig(viper.GetStringMapString("library_options"))

	utils.Config(viper.Sub("tui"))
	//bubbles.Config(viper.Sub("tui"))
	//tquery()
	//pref()

	start := flaggy.NewSubcommand("start")
	flaggy.AttachSubcommand(start, 1)

	serve := flaggy.NewSubcommand("serve")
	flaggy.AttachSubcommand(serve, 1)

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

	//if serve.Used {
	//  //urbooks.InitLibraries(viper.Sub("libraries"), viper.GetStringMapString("libraries"), true)
	//  urbooks.Cfg().Opts["serve"] = "true"
	//  urbooks.InitLibraries(viper.Sub("libraries"), viper.GetStringMapString("libraries"), false)
	//  web.InitServer(viper.Sub("website"))
	//  web.Serve()
	//}

	if start.Used {
		urbooks.Cfg().Opts["serve"] = "false"
		urbooks.InitLibraries(viper.Sub("libraries"), viper.GetStringMapString("libraries"), false)
		//fmt.Printf("%V\n", urbooks.NewURL().String())
		//fmt.Println(ui.TermHeight())
		//tui()
		someBooks()
		//rss()
		//feeds()
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

//func rss() {
//  resp, err := urbooks.Get(url)
//  if err != nil {
//    log.Fatal(err)
//  }
//  books := urbooks.ParseBooks(resp)
//  feed := web.NewRSSFeed()
//  feed.NewChannel(books)
//  fmt.Printf("%v\n", string(feed.ToXML()))
//}

func someBooks() {
	//u := "http://localhost:9932/tags/11?currentPage=2&itemsPerPage=50&library=audiobooks&order=desc&sort=added"
	req := urbooks.NewRequest("audiobooks")
	req.From("books")
	req.ID("98")
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

	book := urbooks.ParseBooks(resp).Books()[0]
	//fmt.Printf("%+V\n", book.Get("formats").String())
	fmt.Printf("%+V\n", book.Get("authors").String())
	//fmt.Printf("%+V\n", book.Get("seriesAndTitle").String())
	//fmt.Printf("%+V\n", book.Get("tags").String())
	//fmt.Printf("%+V\n", book.Get("authors").String())

	//for _, b := range book.Get("tags").Data() {
	//fmt.Printf("%+V\n",b.URL() )
	//}
	//return urbooks.ParseBooks(resp).Data[0]
	//fmt.Printf("%V\n", string(resp))

}

func pref() {
	//fmt.Printf("%v", db.LibCfg.Cur().EditableFields())
}

func tquery() {
	//db.LibCfg.Cur().GetBooks()
}

func tui() {
	//var prompt = map[string]string{"one": "one", "two": "two"}
	//bubbles.NewPrompt("test", prompt).Choose()
	//b := someBooks()
	//bubbles.RenderMarkdown(b.ToPlain())
	//bubbles.Simple()
	//fmt.Printf("%+V\n", ui.Cfg)
	//  ui.StylesConfig()
	//p := tea.NewProgram(bubbles.NewMetaViewer().SetHeight(0).SetWidth(0).Book(b).Model())
	//  //p := tea.NewProgram(ui.NewCategoryBrowser())
	//p := tea.NewProgram(bubbles.NewPrompt(prompt))
	//p.EnableMouseCellMotion()

	//if err := p.Start(); err != nil {
	//log.Fatal(err)
	//}
}

func printResults[T any](r []T) {
	fmt.Printf("Results: %v\n", r)
}

func Config() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(filepath.Join(home, ".config/urbooks"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
}

func TUIConfig() *viper.Viper {
	//theme := viper.GetStringMapString("tui")
	return viper.Sub("tui")
}
