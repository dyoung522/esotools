package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ofSimple bool
var ofMarkdown bool
var ofJSON bool
var ofRaw bool

// modsCmd represents the mods command
var ListModsCmd = &cobra.Command{
	Use:   "mods",
	Short: "Lists installed ESO mods",
	Long: `Lists out all the mods that are installed in the ESO mods directory.

By default, this will print out a simple list with only one mod per line. However, other formats may be specified via the flags.
`,

	Run: func(cmd *cobra.Command, args []string) {
		verbose := viper.GetBool("verbose")
		fmt.Println("verbose:", verbose)

		for _, arg := range args {
			fmt.Println(arg)
		}

		switch {
		case ofJSON:
			fmt.Println("using JSON format")
		case ofMarkdown:
			fmt.Println("using Markdown format")
		case ofRaw:
			fmt.Println("using RAW format")
		default:
			fmt.Println("using Simple format")
		}

		fmt.Println("list mods called")
	},
}

func init() {
	ListModsCmd.Flags().BoolVarP(&ofJSON, "json", "j", false, "Print out the list in JSON format")
	ListModsCmd.Flags().BoolVarP(&ofMarkdown, "markdown", "m", false, "Print out the list in markdown format")
	ListModsCmd.Flags().BoolVarP(&ofRaw, "raw", "r", false, "Print out the list in the RAW ESO mod header format (most verbose)")
	ListModsCmd.Flags().BoolVarP(&ofSimple, "simple", "s", false, "Prints the mod listing in simple plain text (default)")
	ListModsCmd.MarkFlagsMutuallyExclusive("json", "markdown", "raw", "simple")
}
