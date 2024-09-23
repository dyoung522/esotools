package eso_test

import (
	"fmt"
	"testing"

	"github.com/dyoung522/esotools/lib/eso"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.Set("noColor", true)
}

func TestPrint_NoAddOns(t *testing.T) {
	addons := eso.AddOns{}
	output := addons.Print("markdown")
	expected := "Total: 0 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_MarkdownFormat(t *testing.T) {
	addons := eso.AddOns{
		"addon1": eso.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": eso.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
	}
	output := addons.Print("markdown")
	expected := "## Addon One (v1.0) by Author One\n\n## Addon Two (v2.0) by Author Two\n\nTotal: 2 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_HeaderFormat(t *testing.T) {
	addons := eso.AddOns{
		"addon1": eso.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": eso.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
	}
	output := addons.Print("header")
	expected := "## Title: Addon One\n## Description: \n## Author: Author One\n## Version: 1.0\n## AddOnVersion: \n## APIVersion: \n## SavedVariables: \n## DependsOn: \n## OptionalDependsOn: \n## IsDependency: false\n## IsLibrary: false\n\n## Title: Addon Two\n## Description: \n## Author: Author Two\n## Version: 2.0\n## AddOnVersion: \n## APIVersion: \n## SavedVariables: \n## DependsOn: \n## OptionalDependsOn: \n## IsDependency: false\n## IsLibrary: false\n\nTotal: 2 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_JsonFormat(t *testing.T) {
	var addon1, addon2 eso.AddOn

	addon1, _ = eso.NewAddOn("addon1")
	addon2, _ = eso.NewAddOn("addon2")

	addon1.Title = "Addon One"
	addon1.Author = "Author One"
	addon1.Version = "1.0"

	addon2.Title = "Addon Two"
	addon2.Author = "Author Two"
	addon2.Version = "2.0"

	addons := eso.AddOns{"addon1": addon1, "addon2": addon2}
	output := addons.Print("json")
	expected := `[{"Title":"Addon One","Author":"Author One","Contributors":"","Version":"1.0","Description":"","AddOnVersion":"","APIVersion":"","SavedVariables":null,"DependsOn":null,"OptionalDependsOn":null},{"Title":"Addon Two","Author":"Author Two","Contributors":"","Version":"2.0","Description":"","AddOnVersion":"","APIVersion":"","SavedVariables":null,"DependsOn":null,"OptionalDependsOn":null}]`
	assert.JSONEq(t, expected, output)
}

func TestPrint_OnelineMarkdownFormat(t *testing.T) {
	addons := eso.AddOns{
		"addon1": eso.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": eso.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
	}
	output := addons.Print("oneline")
	expected := "- Addon One (v1.0) by Author One\n- Addon Two (v2.0) by Author Two\nTotal: 2 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_WithDependenciesAndLibraries(t *testing.T) {
	addons := eso.AddOns{
		"addon1": eso.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"},
		"addon2": eso.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"},
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
	var addon1, addon2 eso.AddOn

	addon1 = eso.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"}
	addon2 = eso.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"}

	addon1.SetDependency(true)
	addon2.SetLibrary(true)

	addons := eso.AddOns{"addon1": addon1, "addon2": addon2}

	viper.Set("noDeps", true)
	viper.Set("noLibs", true)

	output := addons.Print("markdown")
	expected := "Total: 0 AddOns"
	assert.Equal(t, expected, output)
}

func TestPrint_WithSubmodules(t *testing.T) {
	var addon1, addon2 eso.AddOn

	addon1 = eso.AddOn{Title: "Addon One", Author: "Author One", Version: "1.0"}
	addon2 = eso.AddOn{Title: "Addon Two", Author: "Author Two", Version: "2.0"}

	addon1.SetDir("path/to/submodule")

	addons := eso.AddOns{"addon1": addon1, "addon2": addon2}

	output := addons.Print("markdown")
	expected := "## Addon Two (v2.0) by Author Two\n\nTotal: 1 AddOns"
	assert.Equal(t, expected, output)
}

func TestValidate_WithValidAddon(t *testing.T) {
	addon1, _ := eso.NewAddOn("ValidAddon")

	assert.True(t, addon1.Validate())
}

func TestValidateTitle_MissingTitle(t *testing.T) {
	addon1, _ := eso.NewAddOn("InvalidTitle")
	addon1.Author = "Author One"
	addon1.Version = "1.0"

	valid := addon1.Validate()
	assert.True(t, valid)
	assert.Equal(t, addon1.Title, "Invalid Title")

}

func TestValidateTitle_MissingAuthor(t *testing.T) {
	addon1, _ := eso.NewAddOn("InvalidAuthor")
	addon1.Title = "Invalid Author"
	addon1.Version = "1.0"

	valid := addon1.Validate()
	assert.True(t, valid)
	assert.Equal(t, addon1.Author, "Unknown")
}

func TestValidateTitle_MissingVersion(t *testing.T) {
	addon1, _ := eso.NewAddOn("InvalidVersion")
	addon1.Title = "Invalid Version"
	addon1.Author = "Author One"

	valid := addon1.Validate()
	assert.True(t, valid)
	assert.Equal(t, addon1.Version, "0")
}

func TestValidateTitle_InvalidAddon(t *testing.T) {
	var errs []error

	addon1, _ := eso.NewAddOn("")

	expectedErrors := append(errs, fmt.Errorf("'Title' is required"))

	valid := addon1.Validate()
	assert.False(t, valid)
	assert.Equal(t, addon1.Errors(), expectedErrors)
}

func TestStripESOColorCodes(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "no color codes",
			input:  "Test String",
			expect: "Test String",
		},
		{
			name:   "single color code",
			input:  "|c121212Test|r String",
			expect: "Test String",
		},
		{
			name:   "multiple color codes",
			input:  "|c121212Test|r String |c323232Multiple|r Codes",
			expect: "Test String Multiple Codes",
		},
		{
			name:   "color codes at end of string",
			input:  "|c121212Test|r String |c323232Multiple|r",
			expect: "Test String Multiple",
		},
		{
			name:   "invalid color codes",
			input:  "|c00FF00Test |cFFFF00String|r",
			expect: "Test String",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			actual := eso.StripESOColorCodes(tt.input)
			assert.Equal(tt.expect, actual, "Expected %q, but got %q", tt.expect, actual)
		})
	}
}
