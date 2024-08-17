package modTools

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/spf13/afero"
)

func ReadMods(modList *[]string) (Mods, []error) {
	var mods = Mods{}
	var errors []error
	var re = regexp.MustCompile(`##\s+(?P<Type>\w+):\s(?P<Data>.*)\s*$`)

	for _, modfile := range *modList {
		file, err := AppFs.Open(modfile)
		if err != nil {
			fmt.Println("Error opening file", modfile, err)
			errors = append(errors, err)
			continue
		}

		data, err := afero.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading file", modfile, err)
			errors = append(errors, err)
			continue
		}

		basename := filepath.Base(modfile)
		relativePath, _ := filepath.Rel(AddOnsPath, modfile)
		key := cleanString(strings.ToLower(strings.TrimSuffix(basename, filepath.Ext(basename))))

		mod := NewMod(key)
		mod.meta.dir = filepath.Dir(relativePath)

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
		file.Close()

		// Don't add submodules to the list (for now)
		if dup, exists := mods.Find(mod.key); exists {
			if !mod.IsSubmodule() {
				if dup.IsSubmodule() {
					mods.Replace(mod)
				} else {
					errors = append(errors, fmt.Errorf("duplicate mods found for %s\n%v\n%v", mod.key, mod, dup))
				}
			}

			continue
		}

		if mod.Valid() {
			mods.Add(mod)
		}
	}

	markDependencies(&mods)

	return mods, errors
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

func markDependencies(mods *Mods) {
	for key, mod := range *mods {
		if len(mod.DependsOn) == 0 {
			continue
		}

		// Mark submodules as dependencies (of their parent)
		if mod.IsSubmodule() {
			mod.meta.dependency = true
			mods.Replace(mod)
		}

		for _, dependency := range mod.DependsOn {
			dependencyName := dependencyName(dependency)

			// Skip self-references
			if dependencyName == "" || dependencyName == key {
				continue
			}

			if depmod, exists := mods.Find(dependencyName); exists {
				depmod.meta.dependency = true
				mods.Replace(depmod)
			} else {
				fmt.Println("Missing Dependency:", dependencyName)
			}
		}
	}
}
