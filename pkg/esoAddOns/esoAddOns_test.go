package esoAddOns_test

import (
	"fmt"
	"testing"

	"github.com/dyoung522/esotools/pkg/esoAddOns"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.Set("noColor", true)
}

func TestPrint_NoAddOns(t *testing.T) {
	addons := esoAddOns.AddOns{}
	output := addons.Print("markdown")
	expected := "Total: 0 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_MarkdownFormat(t *testing.T) {
	addons := esoAddOns.AddOns{
		"addon1": esoAddOns.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": esoAddOns.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
	}
	output := addons.Print("markdown")
	expected := "## Addon One (v1.0) by Author One\n\n## Addon Two (v2.0) by Author Two\n\nTotal: 2 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_HeaderFormat(t *testing.T) {
	addons := esoAddOns.AddOns{
		"addon1": esoAddOns.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": esoAddOns.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
	}
	output := addons.Print("header")
	expected := "## Title: Addon One\n## Description: \n## Author: Author One\n## Version: 1.0\n## AddOnVersion: \n## APIVersion: \n## SavedVariables: \n## DependsOn: \n## OptionalDependsOn: \n## IsDependency: false\n## IsLibrary: false\n\n## Title: Addon Two\n## Description: \n## Author: Author Two\n## Version: 2.0\n## AddOnVersion: \n## APIVersion: \n## SavedVariables: \n## DependsOn: \n## OptionalDependsOn: \n## IsDependency: false\n## IsLibrary: false\n\nTotal: 2 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_JsonFormat(t *testing.T) {
	var addon1, addon2 esoAddOns.AddOn

	addon1, _ = esoAddOns.NewAddOn("addon1")
	addon2, _ = esoAddOns.NewAddOn("addon2")

	addon1.Title = "Addon One"
	addon1.Author = "Author One"
	addon1.Version = "1.0"

	addon2.Title = "Addon Two"
	addon2.Author = "Author Two"
	addon2.Version = "2.0"

	addons := esoAddOns.AddOns{"addon1": addon1, "addon2": addon2}
	output := addons.Print("json")
	expected := `[{"Title":"Addon One","Author":"Author One","Contributors":"","Version":"1.0","Description":"","AddOnVersion":"","APIVersion":"","SavedVariables":null,"DependsOn":null,"OptionalDependsOn":null},{"Title":"Addon Two","Author":"Author Two","Contributors":"","Version":"2.0","Description":"","AddOnVersion":"","APIVersion":"","SavedVariables":null,"DependsOn":null,"OptionalDependsOn":null}]`
	assert.JSONEq(t, expected, output)
}

func TestPrint_OnelineMarkdownFormat(t *testing.T) {
	addons := esoAddOns.AddOns{
		"addon1": esoAddOns.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": esoAddOns.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
	}
	output := addons.Print("oneline")
	expected := "- Addon One (v1.0) by Author One\n- Addon Two (v2.0) by Author Two\nTotal: 2 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_WithDependenciesAndLibraries(t *testing.T) {
	addons := esoAddOns.AddOns{
		"addon1": esoAddOns.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": esoAddOns.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
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
	var addon1, addon2 esoAddOns.AddOn

	addon1 = esoAddOns.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"}
	addon2 = esoAddOns.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"}

	addon1.SetDependency(true)
	addon2.SetLibrary(true)

	addons := esoAddOns.AddOns{"addon1": addon1, "addon2": addon2}

	viper.Set("noDeps", true)
	viper.Set("noLibs", true)

	output := addons.Print("markdown")
	expected := "Total: 0 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_WithSubmodules(t *testing.T) {
	var addon1, addon2 esoAddOns.AddOn

	addon1 = esoAddOns.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"}
	addon2 = esoAddOns.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"}

	addon1.SetDir("path/to/submodule")

	addons := esoAddOns.AddOns{"addon1": addon1, "addon2": addon2}

	output := addons.Print("markdown")
	expected := "## Addon Two (v2.0) by Author Two\n\nTotal: 1 AddOns"
	assert.Equal(t, expected, output)
}

func TestValidate_WithValidAddon(t *testing.T) {
	addon1, _ := esoAddOns.NewAddOn("ValidAddon")

	assert.True(t, addon1.Validate())
}

func TestValidateTitle_MissingTitle(t *testing.T) {
	addon1, _ := esoAddOns.NewAddOn("InvalidTitle")
	addon1.Author = "Author One"
	addon1.Version = "1.0"

	valid := addon1.Validate()
	assert.True(t, valid)
	assert.Equal(t, addon1.Title, "Invalid Title")

}

func TestValidateTitle_MissingAuthor(t *testing.T) {
	addon1, _ := esoAddOns.NewAddOn("InvalidAuthor")
	addon1.Title = "Invalid Author"
	addon1.Version = "1.0"

	valid := addon1.Validate()
	assert.True(t, valid)
	assert.Equal(t, addon1.Author, "Unknown")
}

func TestValidateTitle_MissingVersion(t *testing.T) {
	addon1, _ := esoAddOns.NewAddOn("InvalidVersion")
	addon1.Title = "Invalid Version"
	addon1.Author = "Author One"

	valid := addon1.Validate()
	assert.True(t, valid)
	assert.Equal(t, addon1.Version, "0")
}

func TestValidateTitle_InvalidAddon(t *testing.T) {
	var errs []error

	addon1, _ := esoAddOns.NewAddOn("")

	expectedErrors := append(errs, fmt.Errorf("'Title' is required"))

	valid := addon1.Validate()
	assert.False(t, valid)
	assert.Equal(t, addon1.Errors(), expectedErrors)
}

func TestCleanTitle_WithColorCode(t *testing.T) {
	addon1 := esoAddOns.AddOn{Title: "|c121212Addon|r One", Author: "Author One", Version: "1.0"}

	output := addon1.CleanTitle()
	expected := "Addon One"

	assert.Equal(t, expected, output)
}

func TestCleanTitle_WithMultipleColorCode(t *testing.T) {
	addon1 := esoAddOns.AddOn{Title: "|c121212Addon|r One |c323232Forever|r Clean", Author: "Author One", Version: "1.0"}

	output := addon1.CleanTitle()
	expected := "Addon One Forever Clean"

	assert.Equal(t, expected, output)
}

func TestCleanTitle_WithMultipleColorCodeAtEnd(t *testing.T) {
	addon1 := esoAddOns.AddOn{Title: "|c121212Addon|r One |c323232Forever|r", Author: "Author One", Version: "1.0"}

	output := addon1.CleanTitle()
	expected := "Addon One Forever"

	assert.Equal(t, expected, output)
}

func TestCleanTitle_WithInvalidColorCodes(t *testing.T) {
	addon1 := esoAddOns.AddOn{Title: "|c00FF00Addon |cFFFF00One|r", Author: "Author One", Version: "1.0"}

	output := addon1.CleanTitle()
	expected := "Addon One"

	assert.Equal(t, expected, output)
}

func TestCleanTitle_WithoutColorCode(t *testing.T) {
	addon1 := esoAddOns.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"}

	output := addon1.CleanTitle()
	expected := "Addon One"

	assert.Equal(t, expected, output)
}
