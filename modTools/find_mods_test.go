package modTools

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func init() {
	AppFs = afero.NewMemMapFs()
	modPath := filepath.Join("/", "tmp", "eso", "live", "AddOns", "TestMod")

	AppFs.MkdirAll(modPath, 0755)
	afero.WriteFile(AppFs, filepath.Join(modPath, "TestMod.txt"), []byte("## Name: TestMNod\n"), 0644)
}

func TestFindMods(t *testing.T) {
	t.Run("ENV-SET", func(t *testing.T) {
		esoHome := filepath.Join("/tmp", "eso")
		t.Setenv("ESO_HOME", esoHome)

		mods, err := FindMods()
		if err != nil {
			t.Errorf("did not expect error message; got %v", err)
			return
		}

		if len(mods) != 1 {
			t.Error("Did not find any mods")
		}

	})
	t.Run("ENV-NOT-SET", func(t *testing.T) {
		t.Setenv("ESO_HOME", "")
		_, err := FindMods()

		if err == nil {
			t.Error("expected error message, not nil")
		}
	})
}
