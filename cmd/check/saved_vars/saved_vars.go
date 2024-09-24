package cmd

import (
	"fmt"
	"os"
	"strings"

	backupCmd "github.com/dyoung522/esotools/cmd/backup/saved_vars"
	"github.com/dyoung522/esotools/eso"
	"github.com/dyoung522/esotools/eso/addon"
	"github.com/dyoung522/esotools/eso/savedvar"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	backup bool
	clean  bool
	dryRun bool
}

var (
	red     = pterm.NewStyle(pterm.FgRed)
	caution = pterm.NewStyle(pterm.BgRed, pterm.FgYellow, pterm.Bold)
	yellow  = pterm.NewStyle(pterm.FgYellow)
	green   = pterm.NewStyle(pterm.FgGreen)
	cyan    = pterm.NewStyle(pterm.FgCyan)
	// blue   = pterm.NewStyle(pterm.FgBlue)
)

// ListAddOnsCmd represents the addons command
var CheckSavedVarsCmd = &cobra.Command{
	Use:   "savedvars",
	Short: "Checks validity of ESO SavedVariables files",
	Long:  "Specifically, it reports on extraneous SavedVariable files that do not correspond to any known AddOn.\nOptionally, you can auto-remove them with the --clean flag.",
	Run:   execute,
}

func execute(cmd *cobra.Command, args []string) {
	var verbosity = viper.GetInt("verbosity")
	var AppFs = afero.NewOsFs()
	var extraneousSavedVars []savedvar.SavedVar

	addons, errs := addon.Run()
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		os.Exit(2)
	}

	if verbosity >= 2 {
		cmd.Println("Checking SavedVariables")
	}

	savedVarFiles, err := savedvar.Find(AppFs)
	if err != nil {
		cmd.Println(err)
		os.Exit(2)
	}

	for _, savedVar := range savedVarFiles {
		savedVarKey := strings.TrimSuffix(savedVar.FileInfo.Name(), ".lua")

		// Skip Zenemax Online files
		if strings.HasPrefix(savedVarKey, "ZO_") {
			continue
		}

		if _, err := addons.Find(savedVarKey); !err {
			extraneousSavedVars = append(extraneousSavedVars, savedVar)
		}
	}

	var numberOfExtraneousSavedVars = len(extraneousSavedVars)

	if numberOfExtraneousSavedVars > 0 {
		yellow.Printf("Found %d extraneous SavedVariable %s\n", numberOfExtraneousSavedVars, eso.Pluralize("file", numberOfExtraneousSavedVars))

		if verbosity >= 1 || flags.clean {
			for _, savedVar := range extraneousSavedVars {
				fmt.Printf("- %s\n", cyan.Sprint(savedVar.FileInfo.Name()))
			}
		}

		if flags.clean {
			var removePrompt = "Remove above files?"

			if flags.dryRun {
				removePrompt += " [dry-run enabled, no destructive actions will be taken]"
			}

			if result, _ := pterm.DefaultInteractiveConfirm.Show(removePrompt); result {
				if !flags.backup && !flags.dryRun {
					savePrompt := caution.Sprint("This opperation is destructive, do you want to make a backup first?")

					if result, _ := pterm.DefaultInteractiveConfirm.Show(savePrompt); result {
						flags.backup = true
					}
				}

				if flags.backup {
					if err = backupCmd.BackupSavedVars(AppFs); err != nil {
						fmt.Println(err)
						os.Exit(2)
					}
				}

				for _, savedVar := range extraneousSavedVars {
					if flags.dryRun {
						yellow.Printf("Would have removed: %q\n", savedVar.Path())
					} else {
						fmt.Println("Removing:", savedVar.Path())
						if err := AppFs.Remove(savedVar.Path()); err != nil {
							red.Printf("Error removing %s: %s\n", savedVar.FileInfo.Name(), err)
						}
					}
				}
			}
		}
	} else {
		green.Println("No extraneous SavedVariables found")
	}
}

func init() {
	CheckSavedVarsCmd.Flags().BoolVarP(&flags.backup, "backup", "", false, "Performs a backup prior to any destructive actions")
	CheckSavedVarsCmd.Flags().BoolVarP(&flags.clean, "clean", "", false, "Removes extranious SavedVariable files")
	CheckSavedVarsCmd.Flags().BoolVarP(&flags.dryRun, "dry-run", "", false, "Shows what changes would be made without actually making them. Use this to double-check before using --clean")
}
