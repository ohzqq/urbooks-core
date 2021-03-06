package cmd

import (
	"fmt"

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
	cmdLib *urbooks.Library
	//req          *urbooks.request
	searchFields = make([]string, 11)
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list books in library",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmdLib = urbooks.Lib(lib)
		//req = buildRequest(args)
		//somebooks()
		fmt.Printf("%+v\n", searchFields[8])
	},
}

//func buildRequest(args []string) *urbooks.Request {
//return cmdLib.Books().Limit("1")
//return cmdLib.DefaultRequest
//if len(args) != 0 {
//req = buildRequest(args)
//}
//return req
//}

func somebooks() {
	req := cmdLib.GetBooks().Limit("1").GetResponse()
	books := req.ParseBooks()
	for _, b := range books.Books {
		println(b.GetMeta("title"))
		//for label, f := range b.EachField() {
		//  fmt.Printf("field: %v\n", label)
		//  fmt.Printf("data %+V\n", f.String())
		//}
	}
}

//func somecat(r *urbooks.request) {
//  resp := makeReq(r)
//  fmt.Println(string(resp))
//}

//func makeReq(r *urbooks.request) []byte {
//  resp, err := urbooks.Get(r.String())
//  if err != nil {
//    log.Fatal(err)
//  }
//  return resp
//}

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
