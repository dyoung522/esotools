package eso

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/* Example
## Title: LibMapData
## Author: Sharlikran
## Version: 1.14
## AddOnVersion: 114
## IsLibrary: true
## APIVersion: 101041 101042
## SavedVariables: LibMapData_SavedVariables
## DependsOn: LibGPS>=71
## OptionalDependsOn: LibDebugLogger>=263 DebugLogViewer>=558
*/

type addonMeta struct {
	key        string
	dir        string
	dependency bool
	library    bool
	errs       []error
}

func (AM addonMeta) String() string {
	return fmt.Sprintf("[key: %s, dir: %s, dependency: %v, library: %v]", AM.key, AM.dir, AM.dependency, AM.library)
}

type AddOn struct {
	Title             string
	Author            string
	Contributors      string
	Version           string
	Description       string
	AddOnVersion      string
	APIVersion        string
	SavedVariables    []string
	DependsOn         []string
	OptionalDependsOn []string
	meta              addonMeta
}

// NewAddOn creates a new instance of AddOn with the provided key.
// The key is converted to a standardized format using ToKey function.
// If the key is empty, an error is returned.
func NewAddOn(key string) (AddOn, error) {
	key = ToKey(key)

	if key == "" {
		return AddOn{}, fmt.Errorf("key is required")
	}

	return AddOn{meta: addonMeta{key: key}}, nil
}

// String returns a string representation of the AddOn.
// It includes all fields of the AddOn, separated by commas.
func (A AddOn) String() string {
	return fmt.Sprintf(
		"Title: %s, Description: %s, Author: %v, Version: %s, AddOnVersion: %s, APIVersion: %s, SavedVariables: %s, DependsOn: %s, OptionalDependsOn: %s, IsDependency: %v, IsLibrary: %v",
		A.CleanTitle(),
		A.CleanDescription(),
		A.CleanAuthor(),
		A.Version,
		A.AddOnVersion,
		A.APIVersion,
		strings.Join(A.SavedVariables, " "),
		strings.Join(A.DependsOn, " "),
		strings.Join(A.OptionalDependsOn, " "),
		A.meta.dependency,
		A.meta.library,
	)
}

// ToOnelineMarkdown returns a string representation of the AddOn in Markdown format,
// with a hyphen (-) before the title.
func (A AddOn) ToOnelineMarkdown() string {
	return fmt.Sprint("- ", A.TitleString())
}

// ToHeader returns a string representation of the AddOn in Markdown format,
// with a header (##) before each field.
func (A AddOn) ToHeader() string {
	return fmt.Sprintf(
		"## Title: %s\n## Description: %s\n## Author: %s\n## Version: %s\n## AddOnVersion: %s\n## APIVersion: %s\n## SavedVariables: %s\n## DependsOn: %s\n## OptionalDependsOn: %s\n## IsDependency: %v\n## IsLibrary: %v\n",
		A.Title,
		A.Description,
		A.Author,
		A.Version,
		A.AddOnVersion,
		A.APIVersion,
		strings.Join(A.SavedVariables, " "),
		strings.Join(A.DependsOn, " "),
		strings.Join(A.OptionalDependsOn, " "),
		A.meta.dependency,
		A.meta.library,
	)
}

// ToMarkdown returns a string representation of the AddOn in Markdown format,
// with a header (##) before the title and a newline before the description.
func (A AddOn) ToMarkdown() string {
	var description string

	if A.Description != "" {
		description = fmt.Sprintf("\n%s\n", A.CleanDescription())
	}

	return fmt.Sprintf("## %s\n%s", A.TitleString(), description)
}

// ToJson returns a byte slice and an error.
// It marshals the AddOn into JSON format and returns the byte slice.
// If an error occurs during the marshalling process, it returns the error.
func (A AddOn) ToJson() ([]byte, error) {
	output, err := json.Marshal(A)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling JSON: %w", err)
	}
	return output, nil
}

// TitleString returns a string representation of the AddOn's title, version, and author.
// It includes the title, version, and author, separated by parentheses and the word "by".
func (A AddOn) TitleString() string {
	var (
		cyan = pterm.NewStyle(pterm.Bold, pterm.FgCyan)
		blue = pterm.NewStyle(pterm.Bold, pterm.FgBlue)
	)

	var (
		title   = cyan.Sprint(A.CleanTitle())
		version = blue.Sprintf("v%s", A.Version)
		author  = A.CleanAuthor()
	)

	if viper.GetBool("noColor") {
		pterm.DisableColor()
	}

	return fmt.Sprintf("%s (%s) by %v", title, version, author)
}

// SetDir sets the directory of the AddOn.
func (A *AddOn) SetDir(dir string) {
	A.meta.dir = dir
}

// Dir returns the directory of the AddOn.
func (A *AddOn) Dir() string {
	return A.meta.dir
}

// SetDependency sets the dependency status of the AddOn.
func (A *AddOn) SetDependency(value bool) {
	A.meta.dependency = value
}

// IsDependency returns the dependency status of the AddOn.
func (A *AddOn) IsDependency() bool {
	return A.meta.dependency
}

// SetLibrary sets the library status of the AddOn.
func (A *AddOn) SetLibrary(value bool) {
	A.meta.library = value
}

// IsLibrary returns the library status of the AddOn.
func (A *AddOn) IsLibrary() bool {
	return A.meta.library
}

// IsSubmodule returns true if the AddOn is a submodule, and false otherwise.
// It checks if the directory of the AddOn contains more than one path element.
func (A AddOn) IsSubmodule() bool {
	files, _ := filepath.Split(A.meta.dir)
	return len(files) > 1
}

// AddError appends an error to the list of errors of the AddOn.
func (A *AddOn) AddError(err error) {
	A.meta.errs = append(A.meta.errs, err)
}

// Errors returns the list of errors of the AddOn.
func (A AddOn) Errors() []error {
	return A.meta.errs
}

// Key returns the key of the AddOn.
func (A AddOn) Key() string {
	return A.meta.key
}

// Validate validates the AddOn.
// It checks if the title, author, and version fields are not empty.
// If any of these fields are empty, it sets an error and returns false.
// Otherwise, it returns true.
func (A *AddOn) Validate() bool {
	if A.Title == "" {
		caser := cases.Title(language.English)
		A.Title = caser.String(strcase.ToDelimited(A.meta.key, ' '))
		if A.Title == "" {
			A.AddError(fmt.Errorf("'Title' is required"))
		}
	}

	if A.Author == "" {
		A.Author = "Unknown"
	}

	if A.Version == "" {
		A.Version = "0"
	}

	return len(A.Errors()) == 0
}

// CleanAuthor returns the Author without any ESO color codes
func (A AddOn) CleanAuthor() string {
	return StripESOColorCodes(A.Author)
}

// CleanTitle returns the Title without any ESO color codes
func (A AddOn) CleanTitle() string {
	return StripESOColorCodes(A.Title)
}

// CleanDescription returns the Author without any ESO color codes
func (A AddOn) CleanDescription() string {
	return StripESOColorCodes(A.Description)
}

/*
 * ADDONS
 */

// Global AddOn Container to store all AddOns.
type AddOns map[string]AddOn

// Add adds an AddOn to the global AddOns map.
// It checks if the AddOn with the same key already exists in the map.
// If it does, it panics with an error message.
// Otherwise, it adds the AddOn to the map.
func (A *AddOns) Add(addon AddOn) {
	if _, exists := (*A)[addon.meta.key]; exists && (*A)[addon.meta.key].meta.key != "" {
		panic(fmt.Errorf("attempt to duplicate an existing AddOn: %v\nPerhaps you meant to use Update()?", (*A)[addon.meta.key]))
	}

	(*A)[addon.meta.key] = addon
}

// Update updates an AddOn in the global AddOns map.
// It checks if the AddOn with the same key exists in the map.
// If it doesn't, it panics with an error message.
// Otherwise, it updates the AddOn in the map.
func (A *AddOns) Update(addon AddOn) {
	if _, exists := (*A)[addon.meta.key]; !exists {
		panic(fmt.Errorf("attempt to add a new AddOn via Update\nPerhaps you meant to use Add()?"))
	}

	(*A)[addon.meta.key] = addon
}

// Get returns an AddOn from the global AddOns map based on the given key.
func (A *AddOns) Get(key string) AddOn {
	return (*A)[ToKey(key)]
}

// Find returns an AddOn and a boolean from the global AddOns map based on the given key.
// The boolean is true if the AddOn is found, and false otherwise.
func (A *AddOns) Find(key string) (AddOn, bool) {
	addon, exists := (*A)[ToKey(key)]
	return addon, exists
}

// Strings returns a string representation of the global AddOns map.
// It includes all AddOns in the map, separated by newlines.
func (A AddOns) Strings() string {
	output := []string{}

	for _, key := range A.Keys() {
		addon := A[key]
		output = append(output, fmt.Sprintf("%s:\n%v", addon.meta.key, addon))
	}

	return strings.Join(output, "\n")
}

// Keys returns a slice of keys from the global AddOns map.
// It sorts the keys in alphabetical order.
func (A AddOns) Keys() []string {
	keys := make([]string, 0, len(A))
	for key := range A {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

// Print prints the global AddOns map in the specified format.
// It checks the format and calls the corresponding function to print the AddOns.
// If the format is "json", it marshals the AddOns into JSON format and prints the byte slice.
// Otherwise, it prints the AddOns in Markdown format.
func (A AddOns) Print(format string) string {
	var (
		deps   bool = !viper.GetBool("noDeps")
		libs   bool = !viper.GetBool("noLibs")
		count       = 0
		output      = []string{}
		addons      = []AddOn{}
	)

	if viper.GetBool("noColor") {
		pterm.DisableColor()
	}

	for _, key := range A.Keys() {
		addon := A[key]

		// Don't print out submodules
		if addon.IsSubmodule() {
			continue
		}

		if (!addon.meta.dependency || deps) && (!addon.meta.library || libs) {
			switch format {
			case "json":
				addons = append(addons, addon)
			case "header":
				output = append(output, fmt.Sprintln(addon.ToHeader()))
			case "markdown":
				output = append(output, fmt.Sprintln(addon.ToMarkdown()))
			default:
				output = append(output, fmt.Sprintln(addon.ToOnelineMarkdown()))
			}
			count++
		}
	}

	if format == "json" {
		jout, _ := json.Marshal(addons)
		return string(jout)
	}

	blue := pterm.NewStyle(pterm.FgBlue)
	return strings.Join(append(output, blue.Sprintf("Total: %d AddOns", count)), "")
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
