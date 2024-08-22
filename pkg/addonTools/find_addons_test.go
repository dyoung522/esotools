package addonTools

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func init() {
	AppFs = afero.NewMemMapFs()
	esoHome := filepath.Join("/tmp", "eso")
	addonPath := filepath.Join(esoHome, "live", "AddOns", "TestAddOn")

	AppFs.MkdirAll(addonPath, 0755)
	afero.WriteFile(AppFs, filepath.Join(addonPath, "TestAddOn.txt"), []byte("## Title: Test AddOn\n"), 0644)
}

func TestFindAddOns(t *testing.T) {
	t.Run("ENV-SET", func(t *testing.T) {
		esoHome := filepath.Join("/tmp", "eso")
		afero.WriteFile(AppFs, ".env", []byte(fmt.Sprintf("ESO_HOME=%q", esoHome)), 0644)
		t.Setenv("ESO_HOME", esoHome)

		addons, err := FindAddOns()

		assert.Nil(t, err, "Expected no error")
		assert.Equal(t, 1, len(addons), "Expected 1 AddOn to be found")
	})
	t.Run("ENV-NOT-SET", func(t *testing.T) {
		_, err := FindAddOns()

		assert.NotNil(t, err, "Expected an error")
		assert.Equal(t, "please set the ESO_HOME environment variable and try again", err.Error(), "Expected error message")
	})
	t.Run("INVALID-ESO-HOME", func(t *testing.T) {
		esoHome := filepath.Join("/tmp", "eso", "invalid")
		afero.WriteFile(AppFs, ".env", []byte(fmt.Sprintf("ESO_HOME=%q", esoHome)), 0644)
		t.Setenv("ESO_HOME", esoHome)

		_, err := FindAddOns()

		assert.NotNil(t, err, "Expected an error")
		assert.Equal(t, "\"/tmp/eso/invalid\" is not a valid ESO HOME directory", err.Error(), "Expected error message")
	})
}
