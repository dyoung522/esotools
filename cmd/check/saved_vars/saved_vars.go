package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dyoung522/esotools/pkg/addonTools"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	clean  bool
	dryRun bool
}

var (
	red    = pterm.NewStyle(pterm.FgRed)
	yellow = pterm.NewStyle(pterm.FgYellow)
	// green  = pterm.NewStyle(pterm.FgGreen)
	cyan = pterm.NewStyle(pterm.FgCyan)
	// blue   = pterm.NewStyle(pterm.FgBlue)
)

// ListAddOnsCmd represents the addons command
var CheckSavedVarsCmd = &cobra.Command{
	Use:   "savedvars",
	Short: "Checks validity of ESO SavedVariables files",
	Long:  "Specifically, it reports on extraneous SavedVariable files that do not correspond to any known AddOn.\nOptionally, you can auto-remove them with the --clean flag.",

	Run: func(cmd *cobra.Command, args []string) {
		var verbosity = viper.GetInt("verbosity")
		var AppFs = afero.NewReadOnlyFs(afero.NewOsFs())
		var extraneousSavedVars []addonTools.SavedVars

		if verbosity >= 1 {
			fmt.Println("Building a list of addons and their dependencies... please wait...")
		}

		addons, errs := addonTools.Run()

		if len(errs) > 0 {
			for _, e := range errs {
				fmt.Println(e)
			}
			os.Exit(2)
		}

		if verbosity >= 2 {
			cmd.Println("Checking SavedVariables")
		}

		savedVarFiles, err := addonTools.FindSavedVars(AppFs)
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
			yellow.Printf("Found %d extraneous SavedVariable %s\n", numberOfExtraneousSavedVars, addonTools.Pluralize("file", numberOfExtraneousSavedVars))

			if verbosity >= 1 || flags.clean {
				for _, savedVar := range extraneousSavedVars {
					fmt.Printf("- %s\n", cyan.Sprint(savedVar.FileInfo.Name()))
				}
			}

			if flags.clean {
				var prompt = "Remove above files?"

				if flags.dryRun {
					prompt += " [dry-run enabled]"
				}

				if result, _ := pterm.DefaultInteractiveConfirm.Show(prompt); result {
					for _, savedVar := range extraneousSavedVars {
						if flags.dryRun {
							yellow.Printf("Would have removed: %q\n", savedVar.FullPath())
						} else {
							if err := AppFs.Remove(savedVar.FullPath()); err != nil {
								red.Printf("Error removing %s: %s\n", savedVar.FileInfo.Name(), err)
							}
						}
					}
				}
			}
		}
	},
}

func init() {
	CheckSavedVarsCmd.Flags().BoolVarP(&flags.clean, "clean", "", false, "Removes extranious SavedVariable files")
	CheckSavedVarsCmd.Flags().BoolVarP(&flags.dryRun, "dry-run", "", false, "Shows what changes would be made without actually making them. Use this to double-check before using --clean")
}
