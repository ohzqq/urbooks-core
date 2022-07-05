package cmd

import (
	"fmt"

	"github.com/ohzqq/urbooks-core/urbooks"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		someBooks()
		//newBook()
	},
}

func newBook() {
	//b := urbooks.NewBook()
	//b.Set("title")
}

func someBooks() {
	//u := "http://localhost:9932/tags/11?currentPage=2&itemsPerPage=50&library=audiobooks&order=desc&sort=added"
	req := urbooks.NewRequest("test-library")
	req.From("books")
	req.ID("216")
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
	fmt.Printf("%+V\n", book)
	fmt.Printf("%+V\n", book.Get("series").String())
	//fmt.Printf("%+V\n", book.Get("description").String())
	//fmt.Printf("%+V\n", book.Get("authors").String())

	//for _, b := range book.Get("tags").Data() {
	//fmt.Printf("%+V\n",b.URL() )
	//}
	//return urbooks.ParseBooks(resp).Data[0]
	//fmt.Printf("%V\n", string(resp))

}
func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
