package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dyoung522/esotools/pkg/addonTools"
	"github.com/dyoung522/esotools/pkg/esoAddOns"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	green  = color.New(color.FgGreen)
	cyan   = color.New(color.FgCyan)
	blue   = color.New(color.FgBlue)
)

var flags struct {
	optional bool
}

// ListAddOnsCmd represents the addons command
var CheckAddOnsCmd = &cobra.Command{
	Use:   "addons",
	Short: "Checks dependencies for ESO AddOns",
	Long:  `Checks AddOns installed in the ESO AddOns directory, and reports any errors`,

	Run: func(cmd *cobra.Command, args []string) {
		var errCount, totalErrs, totalWarns int
		var verbosity = viper.GetInt("verbosity")
		var missingDependencies []string
		var dependencyArray = [2][]string{}

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

		color.NoColor = viper.GetBool("noColor")

		// Check each addon for dependencies, we use Keys() because it's sorted
		for _, key := range addons.Keys() {
			addon := addons[key]

			if verbosity >= 1 {
				cyan.Printf("Checking %s\n", key)
			}

			dependencyArray[0] = addon.DependsOn
			if flags.optional {
				dependencyArray[1] = addon.OptionalDependsOn
			}

			var first = true // Used to determine if we're checking optional dependencies
			for _, dependencies := range dependencyArray {
				numberOfDependencies := len(dependencies)

				if numberOfDependencies == 0 {
					first = false // first element may be empty
					continue
				}

				errCount = 0

				if verbosity >= 1 {
					var descriptor string

					if numberOfDependencies == 1 {
						descriptor = "dependency"
					} else {
						descriptor = "dependencies"
					}

					if first {
						blue.Printf("\tchecking %2d required %-15s ", numberOfDependencies, descriptor)
					} else {
						blue.Printf("\tchecking %2d optional %-15s ", numberOfDependencies, descriptor)
					}
				}

				missingDependencies = checkDependencies(&addons, dependencies)
				errCount += len(missingDependencies)

				if len(missingDependencies) > 0 {
					var errString = fmt.Sprintf("Missing %s\n", strings.Join(missingDependencies, ", "))

					if first {
						red.Print(errString)
						totalErrs += errCount
					} else {
						yellow.Print(errString)
						totalWarns += errCount
					}
				}

				if errCount == 0 && verbosity >= 1 {
					green.Println("√")
				}

				first = false
			}
		}

		if totalErrs > 0 {
			red.Printf("\n%d missing dependencies\n", totalErrs)
		}

		if totalWarns > 0 {
			yellow.Printf("\n%d missing optional dependencies\n", totalWarns)
		}

		if totalErrs > 0 || totalWarns > 0 {
			os.Exit(1)
		} else {
			green.Printf("\nAll %d Addons Ok\n", len(addons))
		}
	},
}

func checkDependencies(addons *esoAddOns.AddOns, dependencies []string) []string {
	var missingDependencies = []string{}

	for _, dependency := range dependencies {
		if len(addonTools.DependencyName(dependency)) == 0 {
			continue
		}

		dependencyName := addonTools.DependencyName(dependency)[0]
		if dependencyName == "" {
			continue
		}

		// Check if the dependency exists
		if _, exists := addons.Find(dependencyName); !exists {
			missingDependencies = append(missingDependencies, dependencyName)
		}
	}

	return missingDependencies
}

func init() {
	CheckAddOnsCmd.Flags().BoolVarP(&flags.optional, "optional", "o", false, "Warn if optional dependencies aren't installed as well")
}
