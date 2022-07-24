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
		lib := urbooks.GetLib(lib)
		getPreferences(lib)
		lib.From("books").Find("333").GetResponse()
		//resp := lib.GetResponse()
		//req = urbooks.NewRequest(lib.Name).From("books").Find("333")
		//req = urbooks.NewRequest(lib.Name).From("authors").Limit("1")
		//req = buildRequest(args)
		//somebooks(req)
		//somecat(req)
		//j := makeReq(req)
		//fmt.Println(string(j))
		//parsed, err := urbooks.ParseBookResponse(j)
		//cat, err := urbooks.ParseCatResponse(j)
		//if err != nil {
		//log.Fatal(err)
		//}
		//fmt.Printf("category %+V\n", lib)
		for _, book := range lib.Books {
			//println("book")
			//fmt.Printf("book %+V\n", book)
			fmt.Printf("%+V\n", book.GetField("added").String())
		}
		//fmt.Printf("%+V\n", parsed.Books()[0].GetField("cover"))

		//if err != nil {
		//log.Fatal(err)
		//}
		//bb := parsed[0]

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

func getPreferences(lib *urbooks.Library) {
	//resp := lib.DB.Get("/preferences?library=test-library")
	//println(string(resp))
	lib.GetDBPreferences()
	fmt.Printf("%+v\n", lib.Pref)
}

func init() {
	rootCmd.AddCommand(testCmd)

}
