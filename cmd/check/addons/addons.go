package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/dyoung522/esotools/pkg/addonTools"
	"github.com/dyoung522/esotools/pkg/esoAddOns"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	red    = pterm.NewStyle(pterm.FgRed)
	yellow = pterm.NewStyle(pterm.FgYellow)
	green  = pterm.NewStyle(pterm.FgGreen)
	cyan   = pterm.NewStyle(pterm.FgCyan)
	blue   = pterm.NewStyle(pterm.FgBlue)
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
		var errors, warnings map[string][]string
		var missingDependencies []string
		var dependencyArray = [2][]string{}
		var verbosity = viper.GetInt("verbosity")

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

		errors = make(map[string][]string)
		warnings = make(map[string][]string)

		if viper.GetBool("noColor") {
			pterm.DisableColor()
		}

		// Check each addon for dependencies, we use Keys() because it's sorted
		for _, key := range addons.Keys() {
			addon := addons[key]

			if verbosity >= 2 {
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

				if verbosity >= 2 {
					var descriptor = addonTools.Pluralize("dependency", numberOfDependencies)

					if first {
						blue.Printf("\tchecking %2d required %-15s ", numberOfDependencies, descriptor)
					} else {
						blue.Printf("\tchecking %2d optional %-15s ", numberOfDependencies, descriptor)
					}
				}

				missingDependencies = checkDependencies(&addons, dependencies)

				if len(missingDependencies) > 0 {
					for _, missingDependency := range missingDependencies {
						if missingDependency == "" {
							continue
						}

						if first {
							errors[missingDependency] = append(errors[missingDependency], key)
						} else {
							warnings[missingDependency] = append(warnings[missingDependency], key)
						}
					}

					if verbosity >= 2 {
						if first {
							red.Println("X")
						} else {
							yellow.Println("X")
						}
					}
				} else {
					if verbosity >= 2 {
						green.Println("√")
					}
				}

				first = false
			}
		}

		if verbosity >= 1 {
			fmt.Println()
		}

		if len(warnings) > 0 {
			printErrors(&warnings, "optional")
		}

		if len(errors) > 0 {
			printErrors(&errors, "required")
			os.Exit(1)
		}

		green.Printf("\nAll %d Addons Ok\n", len(addons))
	},
}

func printErrors(errors *map[string][]string, dependencyType string) {
	var color = pterm.NewStyle(pterm.FgRed)
	var keys = []string{}

	if dependencyType == "optional" {
		color = pterm.NewStyle(pterm.FgYellow)
	}

	for k := range *errors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		numberOfDependencies := len((*errors)[key])
		descriptor := addonTools.Pluralize("AddOn", numberOfDependencies)

		// fmt.Printf("%s is an  %d %s %s: %s\n", key, len((*errors)[key]), dependencyType, descriptor, color.Sprint(strings.Join((*errors)[key], ", ")))
		fmt.Printf(
			"%s is a missing %s dependency for %s -> %s\n",
			color.Add(*pterm.Bold.ToStyle()).Sprintf("%-30s", key),
			color.Sprint(dependencyType),
			cyan.Add(*pterm.Bold.ToStyle()).Sprintf("%d %-6s", numberOfDependencies, descriptor),
			blue.Sprint(pterm.DefaultParagraph.WithMaxWidth(80).Sprint(strings.Join((*errors)[key], ", "))),
		)
	}
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
