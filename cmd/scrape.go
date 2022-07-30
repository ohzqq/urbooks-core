package cmd

import (
	"github.com/ohzqq/urbooks-core/audible"
	"github.com/ohzqq/urbooks-core/book"
	"github.com/spf13/cobra"
)

var (
	audibleUrl string
	batchUrl   string
	noCovers   bool
	query      = audible.NewQuery()
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "scrape audiobook metadata from audible",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		query.SetKeywords(args)
		apicall()
	},
}

func apicall() {
	var books []*book.Book

	switch {
	case audibleUrl != "":
		query.SetUrl(audibleUrl)
		books = append(books, query.GetBook())
	case batchUrl != "":
		query.IsBatch = true
		query.SetUrl(batchUrl)
		books = query.GetBookBatch()
	case query.Keywords != "":
		books = query.Search()
	}

	for _, b := range books {
		b.ConvertTo("ini").Write()
		if !noCovers {
			audible.DownloadCover(b.GetField("title").String(), b.GetFile("cover").Get("url"))
		}
	}
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().BoolVar(&noCovers, "nc", false, "don't download covers")

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
