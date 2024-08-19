package esomods

import (
	"errors"
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
	key        string
	dir        string
	dependency bool
}

func (M modMeta) String() string {
	return fmt.Sprintf("[dir: %v, dependency: %v]", M.dir, M.dependency)
}

type Mod struct {
	Title          string
	Author         string
	Version        string
	AddOnVersion   string
	APIVersion     []string
	SavedVariables []string
	DependsOn      []string
	Errs           []error
	meta           modMeta
}

func NewMod(key string) Mod {
	return Mod{meta: modMeta{key: ToKey(key)}}
}

func (M Mod) String() string {
	return M.Header()
}

func (M Mod) Markdown() string {
	return fmt.Sprintf("- %s (v%s) by %v", M.Title, M.Version, M.Author)
}

func (M Mod) Header() string {
	return fmt.Sprintf(
		"## Title: %s\n"+
			"## Author: %v\n"+
			"## Version: %s\n"+
			"## AddOnVersion: %s\n"+
			"## APIVersion: %v\n"+
			"## SavedVariables: %v\n"+
			"## DependsOn: %v\n"+
			"## Meta: %v\n"+
			"Errors: %v\n",
		M.Title,
		M.Author,
		M.Version,
		M.AddOnVersion,
		strings.Join(M.APIVersion, " "),
		strings.Join(M.SavedVariables, " "),
		strings.Join(M.DependsOn, " "),
		M.meta,
		M.Errs,
	)
}

func (M *Mod) SetDir(dir string) {
	M.meta.dir = dir
}

func (M *Mod) GetDir() string {
	return M.meta.dir
}

func (M *Mod) SetDependency() {
	M.meta.dependency = true
}

func (M *Mod) ClearDependency() {
	M.meta.dependency = false
}

func (M *Mod) IsDependency() bool {
	return M.meta.dependency
}

func (M Mod) IsSubmodule() bool {
	return len(strings.Split(M.meta.dir, "/")) > 1
}

func (M Mod) ID() string {
	return M.meta.key
}

func (M *Mod) Valid() bool {
	M.Errs = []error{}

	switch {
	case M.Title == "":
		M.Errs = append(M.Errs, errors.New("'Title' is missing"))
	case M.Author == "":
		M.Errs = append(M.Errs, errors.New("'Author' is missing"))
		// case M.Version == "":
		// 	M.errs = append(M.errs, fmt.Errorf("'Version' is missing"))
	}

	return len(M.Errs) == 0
}

/**
 ** Global Mods Container
 **/
type Mods map[string]Mod

// These methods may seem a big draconian, but they're intended to prevent bugs
func (M *Mods) Add(mod Mod) {
	if _, exists := (*M)[mod.meta.key]; exists && (*M)[mod.meta.key].meta.key != "" {
		panic(fmt.Errorf("attempt to duplicate an existing Mod: %v\nPerhaps you meant to use Update()?", (*M)[mod.meta.key]))
	}

	(*M)[mod.meta.key] = mod
}

// These methods may seem a big draconian, but they're intended to prevent bugs
func (M *Mods) Update(mod Mod) {
	if _, exists := (*M)[mod.meta.key]; !exists {
		panic(fmt.Errorf("attempt to add a new Mod via Update\nPerhaps you meant to use Add()?"))
	}

	(*M)[mod.meta.key] = mod
}

func (M *Mods) Get(key string) Mod {
	return (*M)[ToKey(key)]
}

func (M *Mods) Find(key string) (Mod, bool) {
	mod, exists := (*M)[ToKey(key)]
	return mod, exists
}

func (M Mods) Strings() string {
	output := []string{}

	for _, key := range M.keys() {
		mod := M[key]
		output = append(output, fmt.Sprintf("%s:\n%v", mod.meta.key, mod))
	}

	return strings.Join(output, "\n")
}

// Prints installed mods (excluding dependencies)
func (M Mods) Print() string {
	return M.print(false)
}

// Prints all mods (including dependencies)
func (M Mods) PrintAll() string {
	return M.print(true)
}

func (M Mods) keys() []string {
	keys := make([]string, 0, len(M))
	for key := range M {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

// Helper function to print mods
func (M Mods) print(all bool) string {
	count := 0
	output := []string{}

	for _, key := range M.keys() {
		mod := M[key]
		if all || !mod.meta.dependency {
			output = append(output, fmt.Sprintln(mod.Markdown()))
			count++
		}
	}

	return strings.Join(append(output, fmt.Sprintln("\nTotal:", count, "mods installed")), "")
}

// Helper Functions
func ToKey(input string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(input), " ", "-"))
}
