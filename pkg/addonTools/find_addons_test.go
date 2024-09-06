package addonTools_test

import (
	"testing"

	"github.com/dyoung522/esotools/pkg/addonTools"
	"github.com/dyoung522/esotools/pkg/esoAddOnFiles"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFindAddOns(t *testing.T) {
	// Create a new in-memory file system
	var fs = afero.NewMemMapFs()
	viper.Set("esohome", "/tmp/eso")

	// Populate the file system with some addon files
	_ = fs.MkdirAll("/tmp/eso/live/AddOns/MyAddon1", 0755)
	_ = fs.MkdirAll("/tmp/eso/live/AddOns/MyAddon2", 0755)
	_ = afero.WriteFile(fs, "/tmp/eso/live/AddOns/MyAddon1/MyAddon1.txt", []byte("## Title: MyAddon1"), 0644)
	_ = afero.WriteFile(fs, "/tmp/eso/live/AddOns/MyAddon2/MyAddon2.txt", []byte("## Title: MyAddon2"), 0644)

	// Call the function we're testing
	addonList, err := addonTools.FindAddOns(fs)

	// Check that the function didn't return an error
	assert.Nil(t, err, "expected no error")

	// Check that the function returned the expected number of addons
	assert.Len(t, addonList, 2, "expected 2 addons")

	// Check that the function returned the expected paths for the addons
	assert.Contains(t, addonList, esoAddOnFiles.AddOnDefinition{Name: "MyAddon1.txt", Dir: "/MyAddon1"}, "expected MyAddon1")
	assert.Contains(t, addonList, esoAddOnFiles.AddOnDefinition{Name: "MyAddon2.txt", Dir: "/MyAddon2"}, "expected MyAddon2")
}

func TestFindAddOnsError(t *testing.T) {
	var addonList []esoAddOnFiles.AddOnDefinition

	g := NewWithT(t)

	// Create a new in-memory file system
	var fs = afero.NewMemMapFs()
	viper.Set("esohome", "/tmp/eso")

	// Populate the file system with some addon files
	_ = fs.MkdirAll("/tmp/eso/live/AddOns/MyAddon1", 0755)
	_ = fs.MkdirAll("/tmp/eso/live/AddOns/MyAddon2", 0755)
	_ = afero.WriteFile(fs, "/tmp/eso/live/addons/MyAddon1/MyAddon1.txt", []byte("## Title: MyAddon1"), 0644)
	_ = afero.WriteFile(fs, "/tmp/eso/live/addons/MyAddon2/MyAddon2.txt", []byte("## Title: MyAddon2"), 0644)

	// Remove the read permission from the addons directory
	_ = fs.Chmod("/tmp/eso/live/AddOns", 0555)

	// Check that the function panics
	g.Expect(func() { addonList, _ = addonTools.FindAddOns(fs) }).NotTo(Panic())

	// Check that the function returned the expected message
	assert.Len(t, addonList, 0, "expected 0 addons")
}
