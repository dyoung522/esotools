package addonTools

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/dyoung522/esotools/pkg/esoAddOns"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func GetAddOns(AppFs afero.Fs) (esoAddOns.AddOns, []error) {
	var errs = []error{}
	var addons = esoAddOns.AddOns{}
	var verbosity = viper.GetInt("verbosity")

	addonlist, err := FindAddOns(AppFs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var re = regexp.MustCompile(`##\s+(?P<Type>\w+):\s(?P<Data>.*)\s*$`)

	for _, addonFile := range addonlist {
		file, err := AppFs.Open(filepath.Join(AddOnsPath, addonFile.Path()))
		if err != nil {
			errs = append(errs, fmt.Errorf("error opening file: %w", err))
			continue
		}
		defer file.Close()

		data, err := afero.ReadAll(file)
		if err != nil {
			errs = append(errs, fmt.Errorf("error reading file: %w", err))
			continue
		}

		addon, err := esoAddOns.NewAddOn(addonFile.Key())
		if err != nil {
			errs = append(errs, fmt.Errorf("could not create addon: %w", err))
			continue
		}

		if verbosity >= 3 {
			fmt.Printf("Parsing %s\n", addonFile.Path())
		}

		addon.SetDir(addonFile.Dir)

		// Create a reader from the byte slice
		reader := bufio.NewReader(bytes.NewReader(data))

		// Read lines until EOF
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				break // EOF or error
			}

			// Remove the trailing newline character
			line = bytes.TrimSuffix(line, []byte("\n"))

			matches := re.FindStringSubmatch(string(line))
			if len(matches) > 1 {
				typeIndex := re.SubexpIndex("Type")
				dataIndex := re.SubexpIndex("Data")

				rawString := matches[dataIndex]
				cleanedString := cleanString(rawString)

				switch matches[typeIndex] {
				case "Title":
					addon.Title = cleanedString
				case "Description":
					addon.Description = rawString
				case "Author":
					addon.Author = rawString
				case "Contributors":
					addon.Contributors = rawString
				case "Version":
					addon.Version = strings.TrimPrefix(cleanedString, "v")
				case "AddOnVersion", "AddonVersion":
					addon.AddOnVersion = cleanedString
				case "APIVersion":
					addon.APIVersion = cleanedString
				case "SavedVariables":
					addon.SavedVariables = strings.Split(cleanedString, " ")
				case "DependsOn":
					addon.DependsOn = strings.Split(cleanedString, " ")
				case "OptionalDependsOn":
					addon.OptionalDependsOn = strings.Split(cleanedString, " ")
				case "IsLibrary":
					addon.SetLibrary(cleanedString == "true")
				default:
					if verbosity >= 3 {
						fmt.Println(fmt.Errorf("unknown type: %s with value: %s", matches[typeIndex], matches[dataIndex]))
					}
				}
			}
		}

		// Don't add submodules to the list (for now)
		if dup, exists := addons.Find(addon.Key()); exists {
			if !addon.IsSubmodule() {
				if dup.IsSubmodule() {
					addons.Update(addon)
				} else {
					fmt.Println(fmt.Errorf("duplicate addons found for %s\n%v\n%v", addon.Key(), addon, dup))
				}
			}

			continue
		}

		if addon.Validate() {
			addons.Add(addon)
		}
	}

	markDependencies(&addons)

	return addons, errs
}

// Cleans up a string by removing any non-graphic characters and extraneous whitespace
func cleanString(input string) string {
	output := strings.TrimFunc(input, func(r rune) bool {
		return !unicode.IsPrint(r)
	})
	return strings.TrimSpace(output)
}

// Removes version dependencies and returns the plain dependency name
func DependencyName(input string) []string {
	return strings.Split(strings.TrimRight(input, "\r\n"), ">=")
}

func markDependencies(addons *esoAddOns.AddOns) {
	for key, addon := range *addons {
		if len(addon.DependsOn) == 0 {
			continue
		}

		// Mark submodules as dependencies (of their parent)
		if addon.IsSubmodule() {
			addon.SetDependency(true)
			addons.Update(addon)
		}

		for _, dependency := range addon.DependsOn {
			dependencyName := DependencyName(dependency)[0]

			// Skip self-references
			if dependencyName == "" || esoAddOns.ToKey(dependencyName) == key {
				continue
			}

			if depaddon, exists := addons.Find(dependencyName); exists {
				depaddon.SetDependency(true)
				addons.Update(depaddon)
			} else {
				fmt.Println(fmt.Errorf("missing Dependency: %s", dependencyName))
			}
		}
	}
}
