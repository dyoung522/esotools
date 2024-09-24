package addon_test

import (
	"path/filepath"
	"testing"

	"github.com/dyoung522/esotools/eso/addon"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var AppFs afero.Fs

func TestAddOnsPath_WithValidESOHome(t *testing.T) {
	// Arrange
	expected := filepath.Clean("/home/user/eso/Elder Scrolls Online/live/AddOns")
	viper.Set("eso_home", "/home/user/eso")

	// Act
	actual := addon.AddOnsPath()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestAddOnsPath_WithEmptyESOHome(t *testing.T) {
	// Arrange
	expected := filepath.Clean("./live/AddOns")
	viper.Set("eso_home", "")

	// Act
	actual := addon.AddOnsPath()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestAddOnsPath_WithESOHomeWithSpaces(t *testing.T) {
	// Arrange
	expected := filepath.Clean("/home/user/path with spaces/Elder Scrolls Online/live/AddOns")
	viper.Set("eso_home", "/home/user/path with spaces")

	// Act
	actual := addon.AddOnsPath()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestAddOnsPath_WithESOHomeWithSpecialCharacters(t *testing.T) {
	// Arrange
	path := filepath.Clean("/home/user/path-with-sp3c14l-char$")
	viper.Set("eso_home", path)

	expected := filepath.Join(path, "/Elder Scrolls Online/live/AddOns")

	// Act
	actual := addon.AddOnsPath()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestSavedVariablesPath_WithValidESOHome(t *testing.T) {
	// Arrange
	expected := filepath.Clean("/home/user/eso/Elder Scrolls Online/live/SavedVariables")
	viper.Set("eso_home", "/home/user/eso")

	// Act
	actual := addon.SavedVariablesPath()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestSavedVariablesPath_WithEmptyESOHome(t *testing.T) {
	// Arrange
	expected := filepath.Clean("./live/SavedVariables")
	viper.Set("eso_home", "")

	// Act
	actual := addon.SavedVariablesPath()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestPluralize_WordEndingInS(t *testing.T) {
	assert := assert.New(t)

	word := "bus"
	count := 2
	expected := "buses"

	actual := addon.Pluralize(word, count)

	assert.Equal(expected, actual, "Expected %q, but got %q", expected, actual)
}

func TestPluralize_WordEndingInX(t *testing.T) {
	assert := assert.New(t)

	word := "box"
	count := 2
	expected := "boxes"

	actual := addon.Pluralize(word, count)

	assert.Equal(expected, actual, "Expected %q, but got %q", expected, actual)
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		name   string
		word   string
		count  int
		expect string
	}{
		{
			name:   "singular word",
			word:   "box",
			count:  1,
			expect: "box",
		},
		{
			name:   "plural word ending in 's'",
			word:   "bus",
			count:  2,
			expect: "buses",
		},
		{
			name:   "plural word ending in 'y' proceeded by a consonant",
			word:   "fly",
			count:  2,
			expect: "flies",
		},
		{
			name:   "plural word ending in 'y' and proceeded by a vowel",
			word:   "flay",
			count:  2,
			expect: "flays",
		},
		{
			name:   "plural word ending in 'ch'",
			word:   "match",
			count:  2,
			expect: "matches",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			actual := addon.Pluralize(tt.word, tt.count)
			assert.Equal(tt.expect, actual, "Expected %q, but got %q", tt.expect, actual)
		})
	}
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
			actual := addon.StripESOColorCodes(tt.input)
			assert.Equal(tt.expect, actual, "Expected %q, but got %q", tt.expect, actual)
		})
	}
}

func TestFind(t *testing.T) {
	// Create a new in-memory file system
	var fs = afero.NewMemMapFs()
	viper.Set("eso_home", "/tmp/eso/Elder Scrolls Online")

	// Populate the file system with some addon files
	_ = fs.MkdirAll("/tmp/eso/Elder Scrolls Online/live/AddOns/MyAddon1", 0755)
	_ = fs.MkdirAll("/tmp/eso/Elder Scrolls Online/live/AddOns/MyAddon2", 0755)
	_ = afero.WriteFile(fs, "/tmp/eso/Elder Scrolls Online/live/AddOns/MyAddon1/MyAddon1.txt", []byte("## Title: MyAddon1"), 0644)
	_ = afero.WriteFile(fs, "/tmp/eso/Elder Scrolls Online/live/AddOns/MyAddon2/MyAddon2.txt", []byte("## Title: MyAddon2"), 0644)

	// Call the function we're testing
	addonList, err := addon.Find(fs)

	// Check that the function didn't return an error
	assert.Nil(t, err, "expected no error")

	// Check that the function returned the expected number of addons
	assert.Len(t, addonList, 2, "expected 2 addons")

	// Check that the function returned the expected paths for the addons
	assert.Contains(t, addonList, addon.AddOnFile{Name: "MyAddon1.txt", Dir: filepath.Clean("/MyAddon1")}, "expected MyAddon1")
	assert.Contains(t, addonList, addon.AddOnFile{Name: "MyAddon2.txt", Dir: filepath.Clean("/MyAddon2")}, "expected MyAddon2")
}

func TestFindError(t *testing.T) {
	var addonList []addon.AddOnFile

	g := NewWithT(t)

	// Create a new in-memory file system
	var fs = afero.NewMemMapFs()
	viper.Set("eso_home", "/tmp/eso")

	// Populate the file system with some addon files
	_ = fs.MkdirAll("/tmp/eso/live/AddOns/MyAddon1", 0755)
	_ = fs.MkdirAll("/tmp/eso/live/AddOns/MyAddon2", 0755)
	_ = afero.WriteFile(fs, "/tmp/eso/live/addons/MyAddon1/MyAddon1.txt", []byte("## Title: MyAddon1"), 0644)
	_ = afero.WriteFile(fs, "/tmp/eso/live/addons/MyAddon2/MyAddon2.txt", []byte("## Title: MyAddon2"), 0644)

	// Remove the read permission from the addons directory
	_ = fs.Chmod("/tmp/eso/live/AddOns", 0555)

	// Check that the function panics
	g.Expect(func() { addonList, _ = addon.Find(fs) }).NotTo(Panic())

	// Check that the function returned the expected message
	assert.Len(t, addonList, 0, "expected 0 addons")
}

func init() {
	AppFs = afero.NewMemMapFs()
	_ = AppFs.MkdirAll("/tmp/eso/Elder Scrolls Online/live/AddOns", 0755)

	viper.Set("eso_home", "/tmp/eso")
}

func TestGetAddOns_EmptyAddonList(t *testing.T) {
	// Arrange
	expected := addon.AddOns{}
	expectedErrs := []error{}

	// Act
	actual, actualErrs := addon.Get(AppFs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedErrs, actualErrs)
}

func TestGetAddOns_MissingRequiredTitle(t *testing.T) {
	// Arrange
	addonName := "MissingRequiredTitle"
	data := []byte(`## Description: My Addon
## Author: Author Name
## Contributors: Contributor 1, Contributor 2
## AddOnVersion: 1.0.0
## APIVersion: 10001
## SavedVariables: MyAddon_SV
## DependsOn:
## OptionalDependsOn:
## IsLibrary: false`)

	err := afero.WriteFile(AppFs, filepath.Join(addon.AddOnsPath(), addonName, addonName+".txt"), data, 0644)
	if err != nil {
		t.Fatal(err)
	}

	expectedTitle := "Missing Required Title"

	// Act
	addons, actualErrs := addon.Get(AppFs)
	require.Empty(t, actualErrs)
	require.Len(t, addons, 1)

	// Assert
	addon, found := addons.Find(addonName)
	require.True(t, found)

	assert.Equal(t, expectedTitle, addon.Title)
}
