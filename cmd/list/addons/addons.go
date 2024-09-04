package cmd

import (
	"fmt"
	"os"

	"github.com/dyoung522/esotools/pkg/addonTools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	simple   bool
	markdown bool
	json     bool
	raw      bool
	noDeps   bool
	noLibs   bool
}

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
		case flags.json:
			fmt.Println(addons.Print("json"))
		case flags.markdown:
			fmt.Println(addons.Print("markdown"))
		case flags.raw:
			fmt.Println(addons.Print("header"))
		default:
			fmt.Println(addons.Print("simple"))
		}
	},
}

func init() {
	var err error

	ListAddOnsCmd.Flags().BoolVarP(&flags.json, "json", "j", false, "Print out the list in JSON format")
	ListAddOnsCmd.Flags().BoolVarP(&flags.markdown, "markdown", "m", false, "Print out the list in markdown format")
	ListAddOnsCmd.Flags().BoolVarP(&flags.raw, "raw", "r", false, "Print out the list in the RAW ESO AddOn header format (most verbose)")
	ListAddOnsCmd.Flags().BoolVarP(&flags.simple, "simple", "s", false, "Prints the AddOn listing in simple plain text")
	ListAddOnsCmd.MarkFlagsMutuallyExclusive("json", "markdown", "raw", "simple")

	ListAddOnsCmd.Flags().BoolVarP(&flags.noLibs, "no-libs", "L", false, "Suppresses printing of AddOns that are considered Libraries")
	err = viper.BindPFlag("noLibs", ListAddOnsCmd.Flags().Lookup("no-libs"))
	if err != nil {
		panic(err)
	}

	ListAddOnsCmd.Flags().BoolVarP(&flags.noDeps, "no-deps", "D", false, "Suppresses printing of AddOns that are dependencies of other AddOns")
	err = viper.BindPFlag("noDeps", ListAddOnsCmd.Flags().Lookup("no-deps"))
	if err != nil {
		panic(err)
	}
}
