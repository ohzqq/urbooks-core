package cmd

import (
	"fmt"
	"log"

	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

const (
	authors   = iota
	added     = iota
	limit     = iota
	narrators = iota
	order     = iota
	published = iota
	publisher = iota
	rating    = iota
	series    = iota
	sort      = iota
	tags      = iota
)

var (
	cmdLib       *urbooks.Library
	req          *urbooks.Request
	searchFields = make([]string, 11)
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmdLib = urbooks.Lib(lib)
		req = buildRequest(args)
		somebooks()
		fmt.Printf("%+v\n", searchFields[8])
	},
}

func buildRequest(args []string) *urbooks.Request {
	return urbooks.NewRequest(lib).From("books").Limit("1")
	//return cmdLib.DefaultRequest
	//if len(args) != 0 {
	//req = buildRequest(args)
	//}
	//return req
}

func somebooks() {
	resp, err := urbooks.Get(req.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp))
	bookResp := urbooks.ParseBooks(resp)
	for _, book := range bookResp.Books() {
		fmt.Println(book.Get("rating").String())
	}
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().StringVarP(&searchFields[authors], "authors", "a", "", "author field")
	//lsCmd.PersistentFlags().StringVarP(&authorsFlag, "authors", "a", "", "author field")
	lsCmd.PersistentFlags().StringVarP(&searchFields[added], "added", "d", "", "date added field")
	//  lsCmd.PersistentFlags().StringVarP(&limitFlag, "limit", "l", "", "limit results")
	//  lsCmd.PersistentFlags().StringVarP(&narratorsFlag, "narrators", "n", "", "narrator field")
	//  lsCmd.PersistentFlags().StringVarP(&orderFlag, "order", "O", "", "order of results (asc or desc)")
	//  lsCmd.PersistentFlags().StringVarP(&publishedFlag, "published", "p", "", "date published")
	//  lsCmd.PersistentFlags().StringVarP(&publisherFlag, "publisher", "P", "", "publisher field")
	//  lsCmd.PersistentFlags().StringVarP(&ratingFlag, "rating", "r", "", "rating field")
	//  lsCmd.PersistentFlags().StringVarP(&seriesFlag, "series", "s", "", "series field")
	//  lsCmd.PersistentFlags().StringVarP(&sortFlag, "sort", "S", "", "sort results by...")
	//  lsCmd.PersistentFlags().StringVarP(&tagsFlag, "tags", "t", "", "tags field")
	//  lsCmd.PersistentFlags().StringVarP(&titleFlag, "title", "T", "", "title field")
}
