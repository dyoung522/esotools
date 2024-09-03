package cmd

import (
	sub1 "github.com/dyoung522/esotools/cmd/check/addons"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Various check commands",
}

func init() {
	CheckCmd.AddCommand(sub1.CheckAddOnsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
