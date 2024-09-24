package eso_test

import (
	"path/filepath"
	"testing"

	"github.com/dyoung522/esotools/eso"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	// Arrange
	viper.Set("eso_home", "/home/user/eso")
	expected := filepath.Clean("/home/user/eso/Elder Scrolls Online")

	// Assert
	assert.Equal(t, expected, eso.Home())
}

func TestPath_WithEmptyESOHome(t *testing.T) {
	// Arrange
	viper.Set("eso_home", "")

	// Assert
	assert.Equal(t, "", eso.Home())
}

func TestHome_WithSpaces(t *testing.T) {
	// Arrange
	viper.Set("eso_home", "/home/user/path with spaces")
	expected := filepath.Clean("/home/user/path with spaces/Elder Scrolls Online")

	// Assert
	assert.Equal(t, expected, eso.Home())
}

func TestHome_WithSpecialCharacters(t *testing.T) {
	// Arrange
	path := filepath.Clean("/home/user/path-with-sp3c14l-char$")
	viper.Set("eso_home", path)
	expected := filepath.Join(path, "Elder Scrolls Online")

	// Assert
	assert.Equal(t, expected, eso.Home())
}

func TestPluralize_WordEndingInS(t *testing.T) {
	assert := assert.New(t)

	word := "bus"
	count := 2
	expected := "buses"

	actual := eso.Pluralize(word, count)

	assert.Equal(expected, actual, "Expected %q, but got %q", expected, actual)
}

func TestPluralize_WordEndingInX(t *testing.T) {
	assert := assert.New(t)

	word := "box"
	count := 2
	expected := "boxes"

	actual := eso.Pluralize(word, count)

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
			actual := eso.Pluralize(tt.word, tt.count)
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
			actual := eso.StripESOColorCodes(tt.input)
			assert.Equal(tt.expect, actual, "Expected %q, but got %q", tt.expect, actual)
		})
	}
}
