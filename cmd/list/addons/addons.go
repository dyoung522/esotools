package cmd

import (
	"fmt"
	"os"

	"github.com/dyoung522/esotools/pkg/addonTools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ofSimple     bool = true
	ofMarkdown   bool
	ofJSON       bool
	ofRaw        bool
	noDependency bool
)

// ListAddOnsCmd represents the addons command
var ListAddOnsCmd = &cobra.Command{
	Use:   "addons",
	Short: "Lists installed ESO AddOns",
	Long: `Lists AddOns installed in the ESO AddOns directory.

By default, this will print out a simple list with only one AddOn per line. However, other formats may be specified via the flags.
`,

	Run: func(cmd *cobra.Command, args []string) {
		addons, errs := addonTools.Run()

		if len(errs) > 0 {
			for _, e := range errs {
				fmt.Println(e)
			}
			os.Exit(1)
		}

		switch {
		case ofJSON:
			fmt.Println("using JSON format")
		case ofMarkdown:
			fmt.Println(addons.Print("markdown", !noDependency))
		case ofRaw:
			fmt.Println(addons.Print("header", !noDependency))
		default:
			fmt.Println(addons.Print("simple", !noDependency))
		}
	},
}

func init() {
	ListAddOnsCmd.Flags().BoolVarP(&ofJSON, "json", "j", false, "Print out the list in JSON format")
	ListAddOnsCmd.Flags().BoolVarP(&ofMarkdown, "markdown", "m", false, "Print out the list in markdown format")
	ListAddOnsCmd.Flags().BoolVarP(&ofRaw, "raw", "r", false, "Print out the list in the RAW ESO AddOn header format (most verbose)")
	ListAddOnsCmd.Flags().BoolVarP(&ofSimple, "simple", "s", false, "Prints the AddOn listing in simple plain text")
	ListAddOnsCmd.MarkFlagsMutuallyExclusive("json", "markdown", "raw", "simple")

	ListAddOnsCmd.Flags().BoolVarP(&noDependency, "no-deps", "D", false, "Suppresses printing of AddOns that are dependencies of other AddOns")
	viper.BindPFlag("no-deps", ListAddOnsCmd.Flags().Lookup("no-deps"))
}
