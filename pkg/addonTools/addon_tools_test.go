package addonTools_test

import (
	"testing"

	"github.com/dyoung522/esotools/pkg/addonTools"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAddOnsPath_WithValidESOHome(t *testing.T) {
	// Arrange
	expected := "/home/user/eso/live/AddOns"
	viper.Set("eso_home", "/home/user/eso")

	// Act
	actual := addonTools.AddOnsPath()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestAddOnsPath_WithEmptyESOHome(t *testing.T) {
	g := NewWithT(t)

	// Arrange
	viper.Set("eso_home", "")

	// Assert
	g.Expect(func() { addonTools.AddOnsPath() }).To(PanicWith("ESO_HOME environment variable not set"))
}

func TestAddOnsPath_WithESOHomeWithSpaces(t *testing.T) {
	// Arrange
	expected := "/home/user/path with spaces/live/AddOns"
	viper.Set("eso_home", "/home/user/path with spaces")

	// Act
	actual := addonTools.AddOnsPath()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestAddOnsPath_WithESOHomeWithSpecialCharacters(t *testing.T) {
	// Arrange
	path := "/home/user/path-with-sp3c14l-char$"
	viper.Set("eso_home", path)

	expected := path + "/live/AddOns"

	// Act
	actual := addonTools.AddOnsPath()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestSavedVariablesPath_WithValidESOHome(t *testing.T) {
	// Arrange
	expected := "/home/user/eso/live/SavedVariables"
	viper.Set("eso_home", "/home/user/eso")

	// Act
	actual := addonTools.SavedVariablesPath()

	// Assert
	assert.Equal(t, expected, actual)
}

func TestSavedVariablesPath_WithEmptyESOHome(t *testing.T) {
	g := NewWithT(t)

	// Arrange
	viper.Set("eso_home", "")

	// Assert
	g.Expect(func() { addonTools.SavedVariablesPath() }).To(PanicWith("ESO_HOME environment variable not set"))
}

func TestPluralize_WordEndingInS(t *testing.T) {
	assert := assert.New(t)

	word := "bus"
	count := 2
	expected := "buses"

	actual := addonTools.Pluralize(word, count)

	assert.Equal(expected, actual, "Expected %q, but got %q", expected, actual)
}

func TestPluralize_WordEndingInX(t *testing.T) {
	assert := assert.New(t)

	word := "box"
	count := 2
	expected := "boxes"

	actual := addonTools.Pluralize(word, count)

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
			actual := addonTools.Pluralize(tt.word, tt.count)
			assert.Equal(tt.expect, actual, "Expected %q, but got %q", tt.expect, actual)
		})
	}
}
