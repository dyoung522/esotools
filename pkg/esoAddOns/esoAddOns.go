package esoAddOns

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/* Example
## Title: Sous Chef
## Author: Wobin, CrazyDutchGuy, KatKat42 & Baertram
## Version: v2.31
## AddOnVersion: 231
## APIVersion: 1010036 101037
## SavedVariables: SousChef_Settings
## DependsOn: LibSort LibAddonMenu-2.0>=3
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

func NewAddOn(key string) (AddOn, error) {
	key = ToKey(key)

	if key == "" {
		return AddOn{}, fmt.Errorf("key is required")
	}

	return AddOn{meta: addonMeta{key: key}}, nil
}

func (A AddOn) String() string {
	return fmt.Sprintf(
		"Title: %s, "+
			"Description: %s, "+
			"Author: %v, "+
			"Version: %s, "+
			"AddOnVersion: %s, "+
			"APIVersion: %s, "+
			"SavedVariables: %s, "+
			"DependsOn: %s, "+
			"OptionalDependsOn: %s, "+
			"IsDependency: %v, "+
			"IsLibrary: %v",
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

func (A AddOn) ToOnelineMarkdown() string {
	return fmt.Sprint("- ", A.TitleString())
}

func (A AddOn) ToHeader() string {
	return fmt.Sprintf(
		"## Title: %s\n"+
			"## Description: %s\n"+
			"## Author: %s\n"+
			"## Version: %s\n"+
			"## AddOnVersion: %s\n"+
			"## APIVersion: %s\n"+
			"## SavedVariables: %s\n"+
			"## DependsOn: %s\n"+
			"## OptionalDependsOn: %s\n"+
			"## IsDependency: %v\n"+
			"## IsLibrary: %v\n",
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

func (A AddOn) ToMarkdown() string {
	var description string

	if A.Description != "" {
		description = fmt.Sprintf("\n%s\n", A.Description)
	}

	return fmt.Sprintf("## %s\n%s", A.TitleString(), description)
}

func (A AddOn) ToJson() ([]byte, error) {
	output, err := json.Marshal(A)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling JSON: %w", err)
	}
	return output, nil
}

func (A AddOn) TitleString() string {
	var (
		cyan  = color.New(color.Bold, color.FgCyan).SprintfFunc()
		blue  = color.New(color.Bold, color.FgBlue).SprintfFunc()
		white = color.New(color.FgHiWhite).SprintfFunc()
	)

	var (
		title   = cyan(A.Title)
		version = blue("v%s", A.Version)
		author  = white(A.Author)
	)

	color.NoColor = viper.GetBool("noColor")

	return fmt.Sprintf("%s (%s) by %v", title, version, author)
}

func (A *AddOn) SetDir(dir string) {
	A.meta.dir = dir
}

func (A *AddOn) Dir() string {
	return A.meta.dir
}

func (A *AddOn) SetDependency(value bool) {
	A.meta.dependency = value
}

func (A *AddOn) IsDependency() bool {
	return A.meta.dependency
}

func (A *AddOn) SetLibrary(value bool) {
	A.meta.library = value
}

func (A AddOn) IsLibrary() bool {
	return A.meta.library
}

func (A AddOn) IsSubmodule() bool {
	files, _ := filepath.Split(A.meta.dir)
	return len(files) > 1
}

func (A *AddOn) AddError(err error) {
	A.meta.errs = append(A.meta.errs, err)
}

func (A AddOn) Errors() []error {
	return A.meta.errs
}

func (A AddOn) Key() string {
	return A.meta.key
}

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

/**
 ** Global AddOn Container
 **/
type AddOns map[string]AddOn

// These methods may seem a big draconian, but they're intended to prevent bugs
func (A *AddOns) Add(addon AddOn) {
	if _, exists := (*A)[addon.meta.key]; exists && (*A)[addon.meta.key].meta.key != "" {
		panic(fmt.Errorf("attempt to duplicate an existing AddOn: %v\nPerhaps you meant to use Update()?", (*A)[addon.meta.key]))
	}

	(*A)[addon.meta.key] = addon
}

// These methods may seem a big draconian, but they're intended to prevent bugs
func (A *AddOns) Update(addon AddOn) {
	if _, exists := (*A)[addon.meta.key]; !exists {
		panic(fmt.Errorf("attempt to add a new AddOn via Update\nPerhaps you meant to use Add()?"))
	}

	(*A)[addon.meta.key] = addon
}

func (A *AddOns) Get(key string) AddOn {
	return (*A)[ToKey(key)]
}

func (A *AddOns) Find(key string) (AddOn, bool) {
	addon, exists := (*A)[ToKey(key)]
	return addon, exists
}

func (A AddOns) Strings() string {
	output := []string{}

	for _, key := range A.keys() {
		addon := A[key]
		output = append(output, fmt.Sprintf("%s:\n%v", addon.meta.key, addon))
	}

	return strings.Join(output, "\n")
}

func (A AddOns) keys() []string {
	keys := make([]string, 0, len(A))
	for key := range A {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

// Helper function to print addons
func (A AddOns) Print(format string) string {
	var (
		deps   bool = !viper.GetBool("noDeps")
		libs   bool = !viper.GetBool("noLibs")
		count       = 0
		output      = []string{}
		addons      = []AddOn{}
	)

	color.NoColor = viper.GetBool("noColor")

	for _, key := range A.keys() {
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

	toBlue := color.New(color.FgBlue).SprintfFunc()
	return strings.Join(append(output, fmt.Sprintln(toBlue("\nTotal: %d AddOns", count))), "")
}

// Helper Functions
func ToKey(input string) string {
	return strings.ReplaceAll(strings.TrimSpace(input), " ", "-")
}
