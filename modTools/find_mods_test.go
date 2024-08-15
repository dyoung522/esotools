package modTools

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func init() {
	AppFs = afero.NewMemMapFs()
	modPath := filepath.Join("/", "tmp", "eso", "live", "AddOns", "TestMod")

	AppFs.MkdirAll(modPath, 0755)
	afero.WriteFile(AppFs, filepath.Join(modPath, "TestMod.txt"), []byte("## Title: Test Mod\n"), 0644)
}

func TestFindMods(t *testing.T) {
	t.Run("ENV-SET", func(t *testing.T) {
		esoHome := filepath.Join("/tmp", "eso")
		t.Setenv("ESO_HOME", esoHome)

		mods, err := FindMods()

		assert.Nil(t, err, "Expected no error")
		assert.Equal(t, 1, len(mods), "Expected 1 mod to be found")
	})
	t.Run("ENV-NOT-SET", func(t *testing.T) {
		_, err := FindMods()

		assert.NotNil(t, err, "Expected an error")
		assert.Equal(t, "please set the ESO_HOME environment variable and try again", err.Error(), "Expected error message")
	})
	t.Run("INVALID-ESO-HOME", func(t *testing.T) {
		esoHome := filepath.Join("/tmp", "eso", "invalid")
		t.Setenv("ESO_HOME", esoHome)

		_, err := FindMods()

		assert.NotNil(t, err, "Expected an error")
		assert.Equal(t, "\"/tmp/eso/invalid\" is not a valid ESO HOME directory", err.Error(), "Expected error message")
	})
}
