package cmd

import (
	sub1 "github.com/dyoung522/esotools/cmd/list/addons"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Various listing commands",
}

func init() {
	ListCmd.AddCommand(sub1.ListAddOnsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
