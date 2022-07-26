package cmd

import (
	"fmt"
	"strings"

	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

var (
	uri     string
	batch   string
	scraper = urbooks.NewAudibleScraper()
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "scrape audiobook metadata from audible",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		scraper.Keywords = strings.Join(args, " ")
		switch {
		case batch != "":
			scraper = scraper.List(batch)
		case uri != "":
			scraper = scraper.Get(uri)
		case scraper.Authors != "" || scraper.Narrators != "" || scraper.Title != "" || scraper.Keywords != "":
			scraper = scraper.Search()
		}

		books := scraper.Scrape()

		if books == nil {
			fmt.Println("cli no results")
		}

		for _, book := range books {
			book.ConvertTo("ffmeta").Print()
		}
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().BoolVar(&scraper.NoCovers, "nc", false, "don't download covers")

	scrapeCmd.Flags().StringVarP(&uri, "url", "u", "", "audible url")
	scrapeCmd.Flags().StringVarP(&batch, "batch", "b", "", "batch scrape from audible search list")
	scrapeCmd.MarkFlagsMutuallyExclusive("url", "batch")

	scrapeCmd.Flags().StringVarP(&scraper.Authors, "authors", "a", "", "book authors")
	scrapeCmd.MarkFlagsMutuallyExclusive("authors", "url")
	scrapeCmd.MarkFlagsMutuallyExclusive("authors", "batch")

	scrapeCmd.Flags().StringVarP(&scraper.Narrators, "narrators", "n", "", "book narrators")
	scrapeCmd.MarkFlagsMutuallyExclusive("narrators", "url")
	scrapeCmd.MarkFlagsMutuallyExclusive("narrators", "batch")

	scrapeCmd.Flags().StringVarP(&scraper.Title, "title", "t", "", "book title")
	scrapeCmd.MarkFlagsMutuallyExclusive("title", "url")
	scrapeCmd.MarkFlagsMutuallyExclusive("title", "batch")

}
