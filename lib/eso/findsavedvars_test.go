package eso_test

import (
	"testing"

	"github.com/dyoung522/esotools/lib/eso"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFindSavedVars(t *testing.T) {
	// Create a new in-memory file system
	var fs = afero.NewMemMapFs()
	viper.Set("eso_home", "/tmp/eso/Elder Scrolls Online")

	// Populate the file system with some addon files
	_ = fs.MkdirAll("/tmp/eso/Elder Scrolls Online/live/SavedVariables", 0755)
	_ = afero.WriteFile(fs, "/tmp/eso/Elder Scrolls Online/live/SavedVariables/MyAddon1.lua", []byte("MyAddon1"), 0644)
	_ = afero.WriteFile(fs, "/tmp/eso/Elder Scrolls Online/live/SavedVariables/MyAddon2.lua", []byte("MyAddon2"), 0644)

	// Call the function we're testing
	savedVarsList, err := eso.FindSavedVars(fs)

	// Check that the function didn't return an error
	assert.Nil(t, err, "expected no error")

	// Check that the function returned the expected number of addons
	assert.Len(t, savedVarsList, 2, "expected 2 SavedVariable files")
}

func TestFindSavedVarsError(t *testing.T) {
	var savedVarsList []eso.SavedVars

	g := NewWithT(t)

	// Create a new in-memory file system
	var fs = afero.NewMemMapFs()
	viper.Set("eso_home", "/tmp/eso")

	// Populate the file system with some addon files
	_ = fs.MkdirAll("/tmp/eso/live/SavedVariables", 0755)

	// Check that the function panics
	g.Expect(func() { savedVarsList, _ = eso.FindSavedVars(fs) }).NotTo(Panic())

	// Check that the function returned the expected message
	assert.Len(t, savedVarsList, 0, "expected 0 SavedVariable files")
}
