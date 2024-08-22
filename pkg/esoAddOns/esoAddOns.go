package esoAddOns

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

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
}

func (A addonMeta) String() string {
	return fmt.Sprintf("[dir: %v, dependency: %v]", A.dir, A.dependency)
}

type AddOn struct {
	Title             string
	Author            string
	Contributors      string
	Version           string
	Description       string
	AddOnVersion      string
	APIVersion        []string
	SavedVariables    []string
	DependsOn         []string
	OptionalDependsOn []string
	IsLibrary         bool
	Errs              []error
	meta              addonMeta
}

func NewAddOn(key string) AddOn {
	return AddOn{meta: addonMeta{key: ToKey(key)}}
}

func (A AddOn) String() string {
	return A.Header()
}

func (A AddOn) Header() string {
	return fmt.Sprintf(
		"## Title: %s\n"+
			"## Description: %s\n"+
			"## Author: %v\n"+
			"## Version: %s\n"+
			"## AddOnVersion: %s\n"+
			"## APIVersion: %v\n"+
			"## SavedVariables: %v\n"+
			"## DependsOn: %v\n"+
			"## OptionalDependsOn: %v\n"+
			"## IsLibrary: %v\n",
		A.Title,
		A.Description,
		A.Author,
		A.Version,
		A.AddOnVersion,
		strings.Join(A.APIVersion, " "),
		strings.Join(A.SavedVariables, " "),
		strings.Join(A.DependsOn, " "),
		strings.Join(A.OptionalDependsOn, " "),
		A.IsLibrary,
	)
}

func (A AddOn) Simple() string {
	return fmt.Sprintf("- %s (v%s) by %v", A.Title, A.Version, A.Author)
}

func (A AddOn) Markdown() string {
	return fmt.Sprintf("## %s (v%s)\nby %s\n", A.Title, A.Version, A.Author)
}

func (A *AddOn) SetDir(dir string) {
	A.meta.dir = dir
}

func (A *AddOn) GetDir() string {
	return A.meta.dir
}

func (A *AddOn) SetDependency() {
	A.meta.dependency = true
}

func (A *AddOn) ClearDependency() {
	A.meta.dependency = false
}

func (A *AddOn) IsDependency() bool {
	return A.meta.dependency
}

func (A AddOn) IsSubmodule() bool {
	return len(strings.Split(A.meta.dir, "/")) > 1
}

func (A AddOn) Key() string {
	return A.meta.key
}

func (A *AddOn) Valididate() bool {
	if A.Title == "" {
		caser := cases.Title(language.English)
		A.Title = caser.String(strcase.ToDelimited(filepath.Base(A.meta.dir), ' '))
		if A.Title == "" {
			A.Errs = append(A.Errs, fmt.Errorf("'Title' is required"))
		}
	}

	if A.Author == "" {
		A.Author = "Unknown"
	}

	if A.Version == "" {
		A.Version = "0"
	}

	return len(A.Errs) == 0
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
	)

	if viper.GetInt("verbosity") >= 2 {
		fmt.Println("Print Dependencies?", deps)
		fmt.Println("Print Libraries?", libs)
	}

	_ = libs // TODO: Implement this

	for _, key := range A.keys() {
		addon := A[key]

		// Don't print out submodules
		if addon.IsSubmodule() {
			continue
		}

		if (!addon.meta.dependency || deps) && (!addon.IsLibrary || libs) {
			switch format {
			case "json":
				// output = append(output, fmt.Sprintln(addon))
			case "header":
				output = append(output, fmt.Sprintln(addon.Header()))
			case "markdown":
				output = append(output, fmt.Sprintln(addon.Markdown()))
			default:
				output = append(output, fmt.Sprintln(addon.Simple()))
			}
			count++
		}
	}

	return strings.Join(append(output, fmt.Sprintln("\nTotal:", count, "AddOns")), "")
}

// Helper Functions
func ToKey(input string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(input), " ", "-"))
}
