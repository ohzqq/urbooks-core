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
	scraper   *urbooks.AudibleScraper
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch url {
		case "":
			search := urbooks.NewAudibleSearch()
			if keywords := strings.Join(args, " "); keywords != "" {
				fmt.Println(keywords)
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
			scraper = urbooks.NewAudibleScraper().Get(url)
		}
		books := scraper.Scrape()
		if books.Books == nil {
			fmt.Println("no results")
		}
		fmt.Printf("%+V\n", books)
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
	scrapeCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "audible url")
	scrapeCmd.PersistentFlags().StringVarP(&authors, "authors", "a", "", "book authors")
	scrapeCmd.PersistentFlags().StringVarP(&narrators, "narrators", "n", "", "book narrators")
	scrapeCmd.PersistentFlags().StringVarP(&title, "title", "t", "", "book title")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scrapeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
