package cmd

import (
	"fmt"

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
		req = urbooks.NewRequest(lib.Name).From("authors").Limit("1")
		//req = buildRequest(args)
		//somebooks(req)
		//somecat(req)
		j := makeReq(req)
		fmt.Println(string(j))

		//parsed, err := book.ParseBooks(j)
		//if err != nil {
		//log.Fatal(err)
		//}
		//bb := parsed[0]

		//fmt.Printf("%+v\n", bb.ListFields())
		//fmt.Printf("%+v\n", parsed.GetSeriesString())
		//parsed.ConvertTo("ffmeta").Print()
		//b := book.NewBook()
		//b.NewField("narrators").SetKind("collection").SetMeta("forever")
		//fmt.Printf("%+v\n", b.GetField("narrators"))
		//fmt.Printf("%+v\n", b.ListFields())
		//b.GetField("duration").SetMeta("bleep & blorp")
		//b.GetField("titleAndSeries").SetMeta("poot")
		//b.ConvertTo("ffmeta").Print()
		//fmt.Printf("%+v\n", b.GetField("titleAndSeries").String())
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

}
