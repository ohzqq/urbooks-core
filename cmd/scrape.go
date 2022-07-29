package cmd

import (
	"fmt"

	"github.com/ohzqq/urbooks-core/audible"
	"github.com/ohzqq/urbooks-core/book"
	"github.com/spf13/cobra"
)

var (
	audibleUrl string
	batchUrl   string
	query      = audible.NewQuery()
	scraper    = audible.NewWebScraper()
	api        = audible.NewRequest()
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "scrape audiobook metadata from audible",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		query.SetKeywords(args)
		//switch {
		//case batch != "":
		//  scraper = scraper.List(batch)
		//case uri != "":
		//  scraper = scraper.Get(uri)
		//case scraper.Authors != "" || scraper.Narrators != "" || scraper.Title != "" || scraper.Keywords != "":
		//  scraper = scraper.Search()
		//}

		//books := scraper.Scrape()

		//if books == nil {
		//fmt.Println("cli no results")
		//}

		//for _, book := range books {
		//book.ConvertTo("ffmeta").Write()
		//}
	},
}

func apicall() {
	var books []*book.Book
	if audibleUrl != "" {
		query.SetUrl(audibleUrl)
		query.IsWeb = true
		books = append(books, query.GetBookMeta())
	}

	if batchUrl != "" {
		println("get batch!")
		query.IsBatch = true
		query.IsWeb = true
		query.SetUrl(batchUrl)
		books = query.GetBookBatch()
	}

	for _, b := range books {
		fmt.Printf("%+V\n", b.GetField("#narrators").String())
	}
	//resp := api.Search()
	//for _, asin := range resp {
	//println(asin)
	//b := api.Product(resp[0])
	//b.ConvertTo("ffmeta").Write()
	//urbooks.DownloadCover(b.GetField("title").String(), b.GetField("cover").Item().Get("url"))
	//fmt.Printf("%+V\n", b.GetField("cover").Item().Get("url"))
	//}
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().BoolVar(&query.NoCovers, "nc", false, "don't download covers")

	scrapeCmd.Flags().StringVarP(&audibleUrl, "url", "u", "", "audible url")
	scrapeCmd.Flags().StringVarP(&batchUrl, "batch", "b", "", "batch scrape from audible search list")
	scrapeCmd.MarkFlagsMutuallyExclusive("url", "batch")

	scrapeCmd.Flags().StringVarP(&query.Authors, "authors", "a", "", "book authors")
	scrapeCmd.MarkFlagsMutuallyExclusive("authors", "url")
	scrapeCmd.MarkFlagsMutuallyExclusive("authors", "batch")

	scrapeCmd.Flags().StringVarP(&query.Narrators, "narrators", "n", "", "book narrators")
	scrapeCmd.MarkFlagsMutuallyExclusive("narrators", "url")
	scrapeCmd.MarkFlagsMutuallyExclusive("narrators", "batch")

	scrapeCmd.Flags().StringVarP(&query.Title, "title", "t", "", "book title")
	scrapeCmd.MarkFlagsMutuallyExclusive("title", "url")
	scrapeCmd.MarkFlagsMutuallyExclusive("title", "batch")

}
