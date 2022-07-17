package cmd

import (
	"fmt"
	"log"

	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

var cmdLib *urbooks.Library
var req *urbooks.Request

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmdLib = urbooks.Lib(lib)
		req = buildRequest(args)
		//fmt.Printf("%+v\n", req.String())
		//fmt.Printf("%v\n", string(req.Response()))
		somebooks()
	},
}

func buildRequest(args []string) *urbooks.Request {
	return cmdLib.DefaultRequest
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
		fmt.Println(book.Get("title").String())
	}
}

func init() {
	rootCmd.AddCommand(testCmd)

}
