package savedvar

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/dyoung522/esotools/eso"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type SavedVar struct {
	fs.FileInfo
}

// New creates a new SavedVar from a fs.FileInfo
func New(f fs.FileInfo) SavedVar {
	return SavedVar{f}
}

// Path returns the full path to the SavedVar file
func (sv SavedVar) Path() string {
	return filepath.Join(Path(), sv.Name())
}

// Find returns a list of SavedVar files in the SavedVariables directory
func Find(AppFs afero.Fs) ([]SavedVar, error) {
	var err error
	var savedVars []SavedVar

	verbosity := viper.GetInt("verbosity")
	savedVarsPath := Path()

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
		savedVars = append(savedVars, SavedVar{file})
	}

	return savedVars, err
}

func Path() string {
	return filepath.Join(filepath.Clean(eso.Home()), "live", "SavedVariables")
}
