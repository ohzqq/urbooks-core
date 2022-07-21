package cmd

import (
	"fmt"

	"github.com/ohzqq/urbooks-core/book"
	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lib := urbooks.Lib(lib)
		println(lib.Name)
		req = urbooks.NewRequest(lib.Name).From("books").Limit("1")
		//somecat(req)
		j := makeReq(req)
		parsed := book.ParseBooks(j)[0]
		fmt.Printf("%+v\n", parsed.GetField("title").String())
		parsed.ConvertTo("ffmeta").Print()
		b := book.NewBook()
		b.GetField("authors").SetMeta("bleep & blorp")
		b.GetField("titleAndSeries").SetMeta("poot")
		b.ConvertTo("ffmeta").Print()
		//fmt.Printf("%+v\n", b.ConvertTo("markdown").Print())
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

}
