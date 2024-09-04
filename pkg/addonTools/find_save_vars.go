package addonTools

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
	return filepath.Join(SavedVarsPath, sv.Name())
}

func FindSavedVars(AppFs afero.Fs) ([]SavedVars, error) {
	var err error
	var verbosity = viper.GetInt("verbosity")
	var savedVars []SavedVars

	if ok, err := afero.DirExists(AppFs, SavedVarsPath); !ok || err != nil {
		return nil, fmt.Errorf("could not find any SavedVariables in %+q", SavedVarsPath)
	}

	if verbosity >= 2 {
		fmt.Println("Searching", SavedVarsPath)
	}

	files, err := afero.ReadDir(AppFs, SavedVarsPath)
	if err != nil {
		return nil, fmt.Errorf("error occurred while reading %q: %w", SavedVarsPath, err)
	}

	if verbosity >= 2 {
		fmt.Println("Found", len(files), "SavedVariable files")
	}

	for _, file := range files {
		savedVars = append(savedVars, SavedVars{file})
	}

	return savedVars, err
}
