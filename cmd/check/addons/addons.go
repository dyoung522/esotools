package cmd

import (
	"fmt"
	"os"

	"github.com/dyoung522/esotools/pkg/addonTools"
	"github.com/spf13/cobra"
)

// ListAddOnsCmd represents the addons command
var CheckAddOnsCmd = &cobra.Command{
	Use:   "addons",
	Short: "Checks dependencies for ESO AddOns",
	Long: `Checks AddOns installed in the ESO AddOns directory, and reports any errors`,

	Run: func(cmd *cobra.Command, args []string) {
		addons, errs := addonTools.Run()

		if len(errs) > 1 {
			for _, e := range errs {
				fmt.Println(e)
			}
			os.Exit(2)
		}

		// Check addons for errors
		// If there are errors, print them and exit with a non-zero status
		fmt.Println("This is where we would check for errors")
		_ = addons
	},
}
