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
		cmdLib = urbooks.GetLib(lib)
		somebooks()
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

}
