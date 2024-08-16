package modTools

import (
	"fmt"
	"sort"
	"strings"
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

type modMeta struct {
	dir        string
	dependency bool
}

func (M *modMeta) String() string {
	return fmt.Sprintf("[dir: %v, dependency: %v]", M.dir, M.dependency)
}

type Mod struct {
	key            string
	errors         []error
	meta           modMeta
	Title          string
	Author         string
	Version        string
	AddOnVersion   string
	APIVersion     []string
	SavedVariables []string
	DependsOn      []string
}

func NewMod(key string) Mod {
	return Mod{key: ToKey(key)}
}

func (M *Mod) String() string {
	return M.Header()
}

func (M *Mod) Markdown() string {
	return fmt.Sprintf("- %s (v%s) by %v", M.Title, M.Version, M.Author)
}

func (M *Mod) Header() string {
	return fmt.Sprintf(
		"## Title: %s\n"+
			"## Author: %v\n"+
			"## Version: %s\n"+
			"## AddOnVersion: %s\n"+
			"## APIVersion: %v\n"+
			"## SavedVariables: %v\n"+
			"## DependsOn: %v\n"+
			"## Meta: %v\n"+
			"errors: %v\n",
		M.Title,
		M.Author,
		M.Version,
		M.AddOnVersion,
		strings.Join(M.APIVersion, " "),
		strings.Join(M.SavedVariables, " "),
		strings.Join(M.DependsOn, " "),
		M.meta,
		M.errors,
	)
}

func (M *Mod) IsSubmodule() bool {
	return len(strings.Split(M.meta.dir, "/")) > 1
}

func (M *Mod) Valid() bool {
	M.errors = []error{}

	switch {
	case M.Title == "":
		M.errors = append(M.errors, fmt.Errorf("Title is missing"))
	case M.Author == "":
		M.errors = append(M.errors, fmt.Errorf("Author is missing"))
		// case M.Version == "":
		// 	M.errors = append(M.errors, fmt.Errorf("Version is missing"))
	}

	return len(M.errors) == 0
}

/**
 ** Global Mods Container
 **/
type Mods map[string]Mod

func (M *Mods) Add(mod Mod) {
	if _, exists := (*M)[mod.key]; exists && (*M)[mod.key].key != "" {
		panic(fmt.Sprintf("Attempt to duplicate an existing Mod: %v\nPerhaps you meant to use Replace()?\n", (*M)[mod.key]))
	}

	(*M)[mod.key] = mod
}

func (M *Mods) Replace(mod Mod) {
	if _, exists := (*M)[mod.key]; !exists {
		panic(fmt.Sprintln("Attempt to add a new Mod via Replace\nPerhaps you meant to use Add()?"))
	}

	(*M)[mod.key] = mod
}

func (M *Mods) Get(key string) Mod {
	return (*M)[ToKey(key)]
}

func (M *Mods) Find(key string) (Mod, bool) {
	mod, exists := (*M)[ToKey(key)]
	return mod, exists
}

func (M *Mods) Strings() string {
	output := []string{}

	for _, key := range M.keys() {
		mod := (*M)[key]
		output = append(output, fmt.Sprintf("%s:\n%v", mod.key, mod))
	}

	return strings.Join(output, "\n")
}

// Prints installed mods (excluding dependencies)
func (M *Mods) Print() string {
	return M.print(false)
}

// Prints all mods (including dependencies)
func (M *Mods) PrintAll() string {
	return M.print(true)
}

func (M *Mods) keys() []string {
	keys := make([]string, 0, len(*M))
	for key := range *M {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

// Helper function to print mods
func (M *Mods) print(all bool) string {
	count := 0
	output := []string{}

	for _, key := range M.keys() {
		mod := (*M)[key]
		if all || !mod.meta.dependency {
			output = append(output, fmt.Sprintln(mod.Markdown()))
			count++
		}
	}

	return strings.Join(append(output, fmt.Sprintln("Total:", count)), "")
}

// Helper Functions
func ToKey(input string) string {
	return strings.ToLower(strings.ReplaceAll(input, " ", "-"))
}
