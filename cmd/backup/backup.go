package cmd

import (
	sub1 "github.com/dyoung522/esotools/cmd/backup/saved_vars"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Various backup commands",
}

func init() {
	BackupCmd.AddCommand(sub1.BackupSavedVarsCmd)
}
