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
}
