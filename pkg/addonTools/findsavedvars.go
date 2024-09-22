package addOnTools

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type SavedVars struct {
	fs.FileInfo
}

func (sv SavedVars) FullPath() string {
	return filepath.Join(SavedVariablesPath(), sv.Name())
}

func FindSavedVars(AppFs afero.Fs) ([]SavedVars, error) {
	var err error
	var savedVars []SavedVars

	verbosity := viper.GetInt("verbosity")
	savedVarsPath := SavedVariablesPath()

	if ok, err := afero.DirExists(AppFs, savedVarsPath); !ok || err != nil {
		return nil, fmt.Errorf("could not find any SavedVariables in %+q", savedVarsPath)
	}

	if verbosity >= 2 {
		fmt.Println("Searching", savedVarsPath)
	}

	files, err := afero.ReadDir(AppFs, savedVarsPath)
	if err != nil {
		return nil, fmt.Errorf("error occurred while reading %q: %w", savedVarsPath, err)
	}

	if verbosity >= 2 {
		fmt.Println("Found", len(files), "savedVariable files")
	}

	for _, file := range files {
		savedVars = append(savedVars, SavedVars{file})
	}

	return savedVars, err
}
