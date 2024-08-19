package modTools

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
)

func FindMods() ([]string, error) {
	var err error
	var mods []string

	if ok, err := afero.DirExists(AppFs, AddOnsPath); !ok || err != nil {
		return nil, fmt.Errorf("%+q is not a valid ESO HOME directory", ESOHOME)
	}

	fmt.Println("Searching", AddOnsPath)

	err = afero.Walk(AppFs, AddOnsPath, func(path string, info fs.FileInfo, err error) error { return getModList(path, &mods, err) })
	if err != nil {
		return nil, fmt.Errorf("error occurred while walking %q: %w", AddOnsPath, err)
	}

	fmt.Println("Found", len(mods), "mod directories")

	return mods, err
}

func getModList(path string, a *[]string, err error) error {
	if err != nil {
		return err
	}

	filename := filepath.Base(path)
	directoryName := filepath.Base(filepath.Dir(path))

	// We're filtering this list to only capture mods directly under the base AddOns directory, elimating subdirectories
	// We may want to include subdirectories in the future, but for now we're only interested in the mods themselves
	// directoryBase := filepath.Base(filepath.Dir(filepath.Dir(path)))
	// if directoryBase != "AddOns" {
	// 	return nil
	// }

	if filepath.Ext(filename) == ".txt" && directoryName+".txt" == filename {
		*a = append(*a, path)
	}

	return nil
}
