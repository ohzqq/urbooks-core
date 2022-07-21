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
		req = urbooks.NewRequest(lib.Name).From("books").Find("333")
		//somecat(req)
		j := makeReq(req)
		parsed := book.ParseBooks(j)[0]
		fmt.Printf("%+v\n", parsed.ListFields())
		//fmt.Printf("%+v\n", parsed.GetSeriesString())
		//parsed.ConvertTo("ffmeta").Print()
		b := book.NewBook()
		b.NewField("narrators").SetKind("collection").SetMeta("forever")
		fmt.Printf("%+v\n", b.GetField("narrators"))
		fmt.Printf("%+v\n", b.ListFields())
		//b.GetField("duration").SetMeta("bleep & blorp")
		b.GetField("titleAndSeries").SetMeta("poot")
		b.ConvertTo("ffmeta").Print()
		//fmt.Printf("%+v\n", b.GetField("titleAndSeries").String())
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

}
