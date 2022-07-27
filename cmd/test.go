package cmd

import (
	"fmt"
	"strings"

	"github.com/ohzqq/urbooks-core/urbooks"
	"github.com/spf13/cobra"
)

var (
	api = urbooks.NewAudibleQuery()
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		api.Keywords = strings.Join(args, " ")
		apicall()
		//cmdLib = urbooks.GetLib(lib)
		//somebooks()
		//getPreferences(cmdLib)
		//resp := lib.GetResponse()
		//req = urbooks.NewRequest(lib.Name).From("authors").Limit("1")
		//req = buildRequest(args)
		//somecat(req)
		//j := makeReq(req)
		//fmt.Println(string(j))
		//cat, err := urbooks.ParseCatResponse(j)
		//if err != nil {
		//log.Fatal(err)
		//}
		//fmt.Printf("%+V\n", lib.Query().Get("library"))
		//fmt.Printf("%+V\n", parsed.Books()[0].GetField("cover"))

		//if err != nil {
		//log.Fatal(err)
		//}
		//bb := parsed[0]

	},
}

func apicall() {
	resp := api.Search()
	//for _, asin := range resp {
	//println(asin)
	fmt.Printf("%+V\n", api.Product(resp[0]))
	//}
}

func getPreferences(lib *urbooks.Library) {
	//val := url.Values{}
	//val.Set("fields", "#duration")
	//val.Set("library", lib.Name)
	//u := url.URL{Path: "preferences", RawQuery: val.Encode()}
	//resp := lib.DB.Get(u.String())
	//println(string(resp))
	lib.GetDBCustomColumns()
	fmt.Printf("duration %+v\n", lib.CustomColumns["duration"])
	//field, err := book.UnmarshalField([]byte(c["meta"]))
	//if err != nil {
	//log.Fatalf("cust col fail %v\n", err)
	//}
	//l.CustomColumns[c["label"]] = field
}

func init() {
	rootCmd.AddCommand(testCmd)
	//testCmd.Flags().BoolVar(&api.NoCovers, "nc", false, "don't download covers")

	testCmd.Flags().StringVarP(&uri, "url", "u", "", "audible url")
	testCmd.Flags().StringVarP(&batch, "batch", "b", "", "batch scrape from audible search list")
	testCmd.MarkFlagsMutuallyExclusive("url", "batch")

	testCmd.Flags().StringVarP(&api.Authors, "authors", "a", "", "book authors")
	testCmd.MarkFlagsMutuallyExclusive("authors", "url")
	testCmd.MarkFlagsMutuallyExclusive("authors", "batch")

	testCmd.Flags().StringVarP(&api.Narrators, "narrators", "n", "", "book narrators")
	testCmd.MarkFlagsMutuallyExclusive("narrators", "url")
	testCmd.MarkFlagsMutuallyExclusive("narrators", "batch")

	testCmd.Flags().StringVarP(&api.Title, "title", "t", "", "book title")
	testCmd.MarkFlagsMutuallyExclusive("title", "url")
	testCmd.MarkFlagsMutuallyExclusive("title", "batch")

}
