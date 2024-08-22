package esoAddOns

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
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

func (M addonMeta) String() string {
	return fmt.Sprintf("[dir: %v, dependency: %v]", M.dir, M.dependency)
}

type AddOn struct {
	Title          string
	Author         string
	Version        string
	AddOnVersion   string
	APIVersion     []string
	SavedVariables []string
	DependsOn      []string
	Errs           []error
	meta           addonMeta
}

func NewAddOn(key string) AddOn {
	return AddOn{meta: addonMeta{key: ToKey(key)}}
}

func (M AddOn) String() string {
	return M.Header()
}

func (M AddOn) Header() string {
	return fmt.Sprintf(
		"## Title: %s\n"+
			"## Author: %v\n"+
			"## Version: %s\n"+
			"## AddOnVersion: %s\n"+
			"## APIVersion: %v\n"+
			"## SavedVariables: %v\n"+
			"## DependsOn: %v\n",
		M.Title,
		M.Author,
		M.Version,
		M.AddOnVersion,
		strings.Join(M.APIVersion, " "),
		strings.Join(M.SavedVariables, " "),
		strings.Join(M.DependsOn, " "),
	)
}

func (M AddOn) Simple() string {
	return fmt.Sprintf("- %s (v%s) by %v", M.Title, M.Version, M.Author)
}

func (M AddOn) Markdown() string {
	return fmt.Sprintf("## %s (v%s)\nby %s\n", M.Title, M.Version, M.Author)
}

func (M *AddOn) SetDir(dir string) {
	M.meta.dir = dir
}

func (M *AddOn) GetDir() string {
	return M.meta.dir
}

func (M *AddOn) SetDependency() {
	M.meta.dependency = true
}

func (M *AddOn) ClearDependency() {
	M.meta.dependency = false
}

func (M *AddOn) IsDependency() bool {
	return M.meta.dependency
}

func (M AddOn) IsSubmodule() bool {
	return len(strings.Split(M.meta.dir, "/")) > 1
}

func (M AddOn) Key() string {
	return M.meta.key
}

func (M *AddOn) Valididate() bool {
	if M.Title == "" {
		caser := cases.Title(language.English)
		M.Title = caser.String(strcase.ToDelimited(filepath.Base(M.meta.dir), ' '))
		if M.Title == "" {
			M.Errs = append(M.Errs, fmt.Errorf("'Title' is required"))
		}
	}

	if M.Author == "" {
		M.Author = "Unknown"
	}

	if M.Version == "" {
		M.Version = "0"
	}

	return len(M.Errs) == 0
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
func (A AddOns) Print(format string, deps bool) string {
	count := 0
	output := []string{}

	for _, key := range A.keys() {
		addon := A[key]

		// Don't print out submodules
		if addon.IsSubmodule() {
			continue
		}

		if !addon.meta.dependency || deps {
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
