package modTools

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	esoMods "github.com/dyoung522/esotools/esoMods"
	"github.com/spf13/afero"
)

func GetMods() (esoMods.Mods, []error) {
	var mods = esoMods.Mods{}
	var errs = []error{}

	modlist, err := FindMods()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var re = regexp.MustCompile(`##\s+(?P<Type>\w+):\s(?P<Data>.*)\s*$`)

	for _, modfile := range modlist {
		file, err := AppFs.Open(modfile)
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

		basename := filepath.Base(modfile)
		relativePath, _ := filepath.Rel(AddOnsPath, modfile)
		key := cleanString(strings.ToLower(strings.TrimSuffix(basename, filepath.Ext(basename))))

		mod := esoMods.NewMod(key)
		mod.SetDir(filepath.Dir(relativePath))

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
					mod.Title = cleanedString
				case "Author":
					mod.Author = rawString
				case "Version":
					mod.Version = cleanedString
				case "AddOnVersion":
					mod.AddOnVersion = cleanedString
				case "APIVersion":
					mod.APIVersion = strings.Split(cleanedString, " ")
				case "SavedVariables":
					mod.SavedVariables = strings.Split(cleanedString, " ")
				case "DependsOn":
					mod.DependsOn = strings.Split(cleanedString, " ")
				}
			}
		}

		// Don't add submodules to the list (for now)
		if dup, exists := mods.Find(mod.ID()); exists {
			if !mod.IsSubmodule() {
				if dup.IsSubmodule() {
					mods.Update(mod)
				} else {
					fmt.Println(fmt.Errorf("duplicate mods found for %s\n%v\n%v", mod.ID(), mod, dup))
				}
			}

			continue
		}

		if mod.Valid() {
			mods.Add(mod)
		}
	}

	markDependencies(&mods)

	return mods, errs
}

// Cleans up a string by removing any non-graphic characters and extraneous whitespace
func cleanString(input string) string {
	output := strings.TrimFunc(input, func(r rune) bool {
		return !unicode.IsPrint(r)
	})
	return strings.TrimSpace(output)
}

// Removes version dependencies and returns the plain dependency name
func dependencyName(input string) string {
	return strings.Split(input, ">")[0]
}

func markDependencies(mods *esoMods.Mods) {
	for key, mod := range *mods {
		if len(mod.DependsOn) == 0 {
			continue
		}

		// Mark submodules as dependencies (of their parent)
		if mod.IsSubmodule() {
			mod.SetDependency()
			mods.Update(mod)
		}

		for _, dependency := range mod.DependsOn {
			dependencyName := dependencyName(dependency)

			// Skip self-references
			if dependencyName == "" || dependencyName == key {
				continue
			}

			if depmod, exists := mods.Find(dependencyName); exists {
				depmod.SetDependency()
				mods.Update(depmod)
			} else {
				fmt.Println(fmt.Errorf("missing Dependency: %s", dependencyName))
			}
		}
	}
}
