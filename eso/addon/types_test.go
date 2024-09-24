package addon_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/dyoung522/esotools/eso/addon"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.Set("noColor", true)
}

func TestPrint_NoAddOns(t *testing.T) {
	addons := addon.AddOns{}
	output := addons.Print("markdown")
	expected := "Total: 0 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_MarkdownFormat(t *testing.T) {
	addons := addon.AddOns{
		"addon1": addon.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": addon.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
	}
	output := addons.Print("markdown")
	expected := "## Addon One (v1.0) by Author One\n\n## Addon Two (v2.0) by Author Two\n\nTotal: 2 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_HeaderFormat(t *testing.T) {
	addons := addon.AddOns{
		"addon1": addon.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": addon.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
	}
	output := addons.Print("header")
	expected := "## Title: Addon One\n## Description: \n## Author: Author One\n## Version: 1.0\n## AddOnVersion: \n## APIVersion: \n## SavedVariables: \n## DependsOn: \n## OptionalDependsOn: \n## IsDependency: false\n## IsLibrary: false\n\n## Title: Addon Two\n## Description: \n## Author: Author Two\n## Version: 2.0\n## AddOnVersion: \n## APIVersion: \n## SavedVariables: \n## DependsOn: \n## OptionalDependsOn: \n## IsDependency: false\n## IsLibrary: false\n\nTotal: 2 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_JsonFormat(t *testing.T) {
	var addon1, addon2 addon.AddOn

	addon1, _ = addon.New("addon1")
	addon2, _ = addon.New("addon2")

	addon1.Title = "Addon One"
	addon1.Author = "Author One"
	addon1.Version = "1.0"

	addon2.Title = "Addon Two"
	addon2.Author = "Author Two"
	addon2.Version = "2.0"

	addons := addon.AddOns{"addon1": addon1, "addon2": addon2}
	output := addons.Print("json")
	expected := `[{"Title":"Addon One","Author":"Author One","Contributors":"","Version":"1.0","Description":"","AddOnVersion":"","APIVersion":"","SavedVariables":null,"DependsOn":null,"OptionalDependsOn":null},{"Title":"Addon Two","Author":"Author Two","Contributors":"","Version":"2.0","Description":"","AddOnVersion":"","APIVersion":"","SavedVariables":null,"DependsOn":null,"OptionalDependsOn":null}]`
	assert.JSONEq(t, expected, output)
}

func TestPrint_OnelineMarkdownFormat(t *testing.T) {
	addons := addon.AddOns{
		"addon1": addon.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": addon.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
	}
	output := addons.Print("oneline")
	expected := "- Addon One (v1.0) by Author One\n- Addon Two (v2.0) by Author Two\nTotal: 2 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_WithDependenciesAndLibraries(t *testing.T) {
	addons := addon.AddOns{
		"addon1": addon.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": addon.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
	}
	addon1 := addons["addon1"]
	addon2 := addons["addon2"]
	addon1.SetDependency(true)
	addon2.SetLibrary(true)

	viper.Set("noDeps", false)
	viper.Set("noLibs", false)
	output := addons.Print("markdown")
	expected := "## Addon One (v1.0) by Author One\n\n## Addon Two (v2.0) by Author Two\n\nTotal: 2 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_WithoutDependenciesAndLibraries(t *testing.T) {
	var addon1, addon2 addon.AddOn

	addon1 = addon.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"}
	addon2 = addon.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"}

	addon1.SetDependency(true)
	addon2.SetLibrary(true)

	addons := addon.AddOns{"addon1": addon1, "addon2": addon2}

	viper.Set("noDeps", true)
	viper.Set("noLibs", true)

	output := addons.Print("markdown")
	expected := "Total: 0 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_WithSubmodules(t *testing.T) {
	var addon1, addon2 addon.AddOn

	addon1 = addon.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"}
	addon2 = addon.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"}

	addon1.SetDir("path/to/submodule")

	addons := addon.AddOns{"addon1": addon1, "addon2": addon2}

	output := addons.Print("markdown")
	expected := "## Addon Two (v2.0) by Author Two\n\nTotal: 1 AddOns"
	assert.Equal(t, expected, output)
}

func TestValidate_WithValidAddon(t *testing.T) {
	addon1, _ := addon.New("ValidAddon")

	assert.True(t, addon1.Validate())
}

func TestValidateTitle_MissingTitle(t *testing.T) {
	addon1, _ := addon.New("InvalidTitle")
	addon1.Author = "Author One"
	addon1.Version = "1.0"

	valid := addon1.Validate()
	assert.True(t, valid)
	assert.Equal(t, addon1.Title, "Invalid Title")

}

func TestValidateTitle_MissingAuthor(t *testing.T) {
	addon1, _ := addon.New("InvalidAuthor")
	addon1.Title = "Invalid Author"
	addon1.Version = "1.0"

	valid := addon1.Validate()
	assert.True(t, valid)
	assert.Equal(t, addon1.Author, "Unknown")
}

func TestValidateTitle_MissingVersion(t *testing.T) {
	addon1, _ := addon.New("InvalidVersion")
	addon1.Title = "Invalid Version"
	addon1.Author = "Author One"

	valid := addon1.Validate()
	assert.True(t, valid)
	assert.Equal(t, addon1.Version, "0")
}

func TestValidateTitle_InvalidAddon(t *testing.T) {
	var errs []error

	addon1, _ := addon.New("")

	expectedErrors := append(errs, fmt.Errorf("'Title' is required"))

	valid := addon1.Validate()
	assert.False(t, valid)
	assert.Equal(t, addon1.Errors(), expectedErrors)
}

func TestAddOnFile_Key_TrimSuffix(t *testing.T) {
	// Arrange
	addOnDef := addon.AddOnFile{
		Name: "MyAddon.txt",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	expectedKey := "MyAddon"
	if key != expectedKey {
		t.Errorf("expected key '%s' for name '%s', but got: '%s'", expectedKey, addOnDef.Name, key)
	}
}

func TestAddOnFile_Key_TrimSuffix_NotEndingInTxt(t *testing.T) {
	// Arrange
	addOnDef := addon.AddOnFile{
		Name: "MyAddon",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	expectedKey := "MyAddon"
	if key != expectedKey {
		t.Errorf("expected key '%s' for name '%s', but got: '%s'", expectedKey, addOnDef.Name, key)
	}
}

func TestAddOnFile_Path(t *testing.T) {
	// Arrange
	addOnDef := addon.AddOnFile{
		Name: "MyAddon.txt",
		Dir:  "path/to/addon",
	}

	// Act
	path := addOnDef.Path()

	// Assert
	expectedPath := filepath.Clean("path/to/addon/MyAddon.txt")
	if path != expectedPath {
		t.Errorf("expected path '%s' for name '%s', but got: '%s'", expectedPath, addOnDef.Name, path)
	}
}

func TestAddOnFile_Path_NameNotEndingInTxt(t *testing.T) {
	// Arrange
	addOnDef := addon.AddOnFile{
		Name: "MyAddon.txt",
		Dir:  "path/to/addon",
	}

	// Act
	path := addOnDef.Path()

	// Assert
	expectedPath := filepath.Clean("path/to/addon/MyAddon.txt")
	if path != expectedPath {
		t.Errorf("expected path '%s' for name '%s', but got: '%s'", expectedPath, addOnDef.Name, path)
	}
}

func TestAddOnFile_Key_EmptyName(t *testing.T) {
	// Arrange
	addOnDef := addon.AddOnFile{
		Name: "",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	if key != "" {
		t.Errorf("expected an empty key for an empty name, but got: %s", key)
	}
}

func TestAddOnFile_Path_LeadingAndTrailingSlashes(t *testing.T) {
	// Arrange
	addOnDef := addon.AddOnFile{
		Name: "MyAddon.txt",
		Dir:  "/path/to/addon/",
	}

	// Act
	path := addOnDef.Path()

	// Assert
	expectedPath := filepath.Clean("/path/to/addon/MyAddon.txt")
	if path != expectedPath {
		t.Errorf("expected path '%s' for name '%s' and dir '%s', but got: '%s'", expectedPath, addOnDef.Name, addOnDef.Dir, path)
	}
}

func TestAddOnFile_Path_DirectoryWithSpaces(t *testing.T) {
	// Arrange
	addOnDef := addon.AddOnFile{
		Name: "MyAddon.txt",
		Dir:  "/path with spaces/to/addon",
	}

	// Act
	path := addOnDef.Path()

	// Assert
	expectedPath := filepath.Clean("/path with spaces/to/addon/MyAddon.txt")
	if path != expectedPath {
		t.Errorf("expected path '%s' for name '%s' and dir '%s', but got: '%s'", expectedPath, addOnDef.Name, addOnDef.Dir, path)
	}
}

func TestAddOnFile_Key_UppercaseLetters(t *testing.T) {
	// Arrange
	addOnDef := addon.AddOnFile{
		Name: "MyADDON.txt",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	expectedKey := "MyADDON"
	if key != expectedKey {
		t.Errorf("expected key '%s' for name '%s', but got: '%s'", expectedKey, addOnDef.Name, key)
	}
}

func TestAddOnFile_Key_SpecialCharacters(t *testing.T) {
	// Arrange
	addOnDef := addon.AddOnFile{
		Name: "My@Addon#1.txt",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	expectedKey := "My@Addon#1"
	if key != expectedKey {
		t.Errorf("expected key '%s' for name '%s', but got: '%s'", expectedKey, addOnDef.Name, key)
	}
}

func TestAddOnFile_Key_NonASCIICharacters(t *testing.T) {
	// Arrange
	addOnDef := addon.AddOnFile{
		Name: "MyÁddón.txt",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	expectedKey := "MyÁddón"
	if key != expectedKey {
		t.Errorf("expected key '%s' for name '%s', but got: '%s'", expectedKey, addOnDef.Name, key)
	}
}
