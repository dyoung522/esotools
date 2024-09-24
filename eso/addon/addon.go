package addon

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/dyoung522/esotools/ostools"
	"github.com/gertd/go-pluralize"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var AppFs = afero.NewReadOnlyFs(afero.NewOsFs())

func Run() (AddOns, []error) {
	var verbosity = viper.GetInt("verbosity")

	if verbosity >= 1 {
		fmt.Println("Building a list of addons and their dependencies... please wait...")
	}

	return Get(AppFs)
}

func AddOnsPath() string {
	return filepath.Join(filepath.Clean(ESOHome()), "live", "AddOns")
}

func SavedVariablesPath() string {
	return filepath.Join(filepath.Clean(ESOHome()), "live", "SavedVariables")
}

func Pluralize(s string, c int) string {
	var pluralize = pluralize.NewClient()

	if c == 1 {
		return s
	}

	return pluralize.Plural(s)
}

func ValidateESOHOME() error {
	verbosity := viper.GetInt("verbosity")
	esoHome := ESOHome()

	if !checkESODir(esoHome) {
		if esoHome != "" {
			fmt.Println(fmt.Errorf("%q does not appear to be a valid ESO directory, attempting auto-detect", esoHome))
		}

		documentsDir, err := ostools.DocumentsDir()
		if err != nil {
			return err
		}

		esoHome = esoDir(documentsDir)

		for !checkESODir(esoHome) {
			fmt.Println(fmt.Errorf("%q is not a valid ESO directory\n", esoHome))

			esoHome, err = pterm.DefaultInteractiveTextInput.WithDefaultValue(documentsDir).Show(`Enter the directory where your "Elder Scrools Online" documents folder lives [CTRL+C to exit]`)
			if err != nil {
				return err
			}

			esoHome = esoDir(esoHome)
		}

		if verbosity >= 2 {
			fmt.Printf("ESO_HOME set to %q\n", esoHome)
		}

		viper.Set("eso_home", string(esoHome))
	}

	return nil
}

func ESOHome() string {
	return esoDir(viper.GetString("eso_home"))
}

// Helper function to convert a string to a key.
// It replaces any spaces in the input string with hyphens.
func ToKey(input string) string {
	return strings.ReplaceAll(strings.TrimSpace(input), " ", "-")
}

func StripESOColorCodes(input string) string {
	if !strings.Contains(input, `|c`) {
		return input
	}

	colorRE := regexp.MustCompile(`\|c[[:xdigit:]]{6}`)
	re := regexp.MustCompile(`\|c[[:xdigit:]]{6}(.*?)(?:\|r)`)

	cleanString := re.ReplaceAllStringFunc(input, func(match string) string {
		parts := re.FindStringSubmatch(match)
		return parts[1]
	})

	// Strip out any remaining color codes from the clean title.
	return colorRE.ReplaceAllString(cleanString, "")
}

// Removes version dependencies and returns the plain dependency name
func DependencyName(input string) []string {
	return strings.Split(strings.TrimRight(input, "\r\n"), ">=")
}

func Find(AppFs afero.Fs) ([]AddOnFile, error) {
	var err error
	var addons []AddOnFile

	verbosity := viper.GetInt("verbosity")
	addonsPath := AddOnsPath()

	if verbosity >= 2 {
		fmt.Println("Searching", addonsPath)
	}

	err = afero.Walk(AppFs, addonsPath, func(path string, info fs.FileInfo, err error) error { return List(path, &addons, err) })
	if err != nil {
		return nil, fmt.Errorf("error occurred while walking %q: %w", addonsPath, err)
	}

	if verbosity >= 2 {
		fmt.Println("Found", len(addons), "AddOn directories")
	}

	return addons, err
}

func List(path string, addons *[]AddOnFile, err error) error {
	var verbosity = viper.GetInt("verbosity")

	if err != nil {
		return err
	}

	if verbosity >= 5 {
		fmt.Println("Searching", path)
	}

	md := AddOnFile{
		Name: filepath.Base(path),
		Dir:  strings.TrimPrefix(filepath.Dir(path), AddOnsPath()),
	}

	if filepath.Ext(md.Name) == ".txt" && ToKey(filepath.Base(md.Dir)) == md.Key() {
		if verbosity >= 3 {
			fmt.Println("Found", md.Name)
		}

		*addons = append(*addons, md)
	}

	return nil
}

// Get returns a list of AddOns and any errors encountered
func Get(AppFs afero.Fs) (AddOns, []error) {
	var errs = []error{}
	var addons = AddOns{}
	var verbosity = viper.GetInt("verbosity")

	addonlist, err := Find(AppFs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var re = regexp.MustCompile(`##\s+(?P<Type>\w+):\s(?P<Data>.*)\s*$`)

	for _, addonFile := range addonlist {
		file, err := AppFs.Open(filepath.Join(AddOnsPath(), addonFile.Path()))
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

		addon, err := New(addonFile.Key())
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
		} else {
			for _, err := range addon.Errors() {
				errs = append(errs, fmt.Errorf("addon %s: %w", addon.Key(), err))
			}
		}
	}

	markDependencies(&addons)

	return addons, errs
}

/*
 * Private Functions
 */

func cleanString(input string) string {
	output := strings.TrimFunc(input, func(r rune) bool {
		return !unicode.IsPrint(r)
	})
	return strings.TrimSpace(output)
}

func markDependencies(addons *AddOns) {
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
			if dependencyName == "" || ToKey(dependencyName) == key {
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

func esoDir(dir string) string {
	dir = filepath.Clean(strings.TrimSpace(dir))

	if dir == "" || dir == "." || dir == "/" || dir == `\\` || dir == `C:\` {
		return ""
	}

	// If the user entered a SavedVariables or AddOns directory, strip it off
	for strings.Contains(dir, "live") {
		dir = filepath.Dir(dir)
	}

	// Add the "Elder Scrolls Online" directory if it's not there
	if !strings.Contains(dir, "Elder Scrolls Online") {
		dir = filepath.Join(dir, "Elder Scrolls Online")
	}

	return dir
}

func checkESODir(dir string) bool {
	dir = esoDir(dir)

	ok, err := afero.DirExists(AppFs, filepath.Join(dir, "live"))
	if err != nil {
		return false
	}

	return ok
}
