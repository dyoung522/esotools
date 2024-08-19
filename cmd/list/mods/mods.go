package cmd

import (
	"fmt"
	"os"

	"github.com/dyoung522/esotools/modTools"
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

// modsCmd represents the mods command
var ListModsCmd = &cobra.Command{
	Use:   "mods",
	Short: "Lists installed ESO mods",
	Long: `Lists mods installed in the ESO mods directory.

By default, this will print out a simple list with only one mod per line. However, other formats may be specified via the flags.
`,

	Run: func(cmd *cobra.Command, args []string) {
		mods, errs := modTools.Run()

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
			fmt.Println(mods.Print("markdown", !noDependency))
		case ofRaw:
			fmt.Println(mods.Print("header", !noDependency))
		default:
			fmt.Println(mods.Print("simple", !noDependency))
		}
	},
}

func init() {
	ListModsCmd.Flags().BoolVarP(&ofJSON, "json", "j", false, "Print out the list in JSON format")
	ListModsCmd.Flags().BoolVarP(&ofMarkdown, "markdown", "m", false, "Print out the list in markdown format")
	ListModsCmd.Flags().BoolVarP(&ofRaw, "raw", "r", false, "Print out the list in the RAW ESO mod header format (most verbose)")
	ListModsCmd.Flags().BoolVarP(&ofSimple, "simple", "s", false, "Prints the mod listing in simple plain text")
	ListModsCmd.MarkFlagsMutuallyExclusive("json", "markdown", "raw", "simple")

	ListModsCmd.Flags().BoolVarP(&noDependency, "no-deps", "D", false, "Suppresses printing of mods that are dependencies of other mods")
	viper.BindPFlag("no-deps", ListModsCmd.Flags().Lookup("no-deps"))
}
