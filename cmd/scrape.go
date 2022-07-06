package cmd

import (
	"fmt"
	"strings"

	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

var (
	url       string
	authors   string
	narrators string
	title     string
	batch     string
	scraper   *urbooks.AudibleScraper
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "scrape audiobook metadata from audible",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		keywords := strings.Join(args, " ")
		switch {
		case authors != "" || narrators != "" || title != "" || keywords != "":
			search := urbooks.NewAudibleSearch()
			var terms []string
			if keywords != "" {
				terms = append(terms, "keywords: "+keywords)
				search.Keywords(keywords)
			}
			if authors != "" {
				terms = append(terms, "authors: "+authors)
				search.Authors(authors)
			}
			if narrators != "" {
				terms = append(terms, "narrators: "+narrators)
				search.Narrators(narrators)
			}
			if title != "" {
				terms = append(terms, "title: "+title)
				search.Title(title)
			}
			fmt.Println(strings.Join(terms, "\n"))
			scraper = search.Search()
		case batch != "":
			scraper = urbooks.NewAudibleScraper().List(batch)
		case url != "":
			scraper = urbooks.NewAudibleScraper().Get(url)
		}
		books := scraper.Scrape()
		if books == nil {
			fmt.Println("no results")
		}
		fmt.Printf("%+V\n", books)
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "audible url")
	scrapeCmd.PersistentFlags().StringVarP(&batch, "batch", "b", "", "batch scrape from audible search list")
	scrapeCmd.MarkFlagsMutuallyExclusive("url", "batch")

	scrapeCmd.PersistentFlags().StringVarP(&authors, "authors", "a", "", "book authors")
	scrapeCmd.MarkFlagsMutuallyExclusive("authors", "url")
	scrapeCmd.MarkFlagsMutuallyExclusive("authors", "batch")

	scrapeCmd.PersistentFlags().StringVarP(&narrators, "narrators", "n", "", "book narrators")
	scrapeCmd.MarkFlagsMutuallyExclusive("narrators", "url")
	scrapeCmd.MarkFlagsMutuallyExclusive("narrators", "batch")

	scrapeCmd.PersistentFlags().StringVarP(&title, "title", "t", "", "book title")
	scrapeCmd.MarkFlagsMutuallyExclusive("title", "url")
	scrapeCmd.MarkFlagsMutuallyExclusive("title", "batch")
}
