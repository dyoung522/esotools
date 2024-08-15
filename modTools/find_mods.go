package modTools

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

var AppFs = afero.NewReadOnlyFs(afero.NewOsFs())

func FindMods() ([]string, error) {
	var err error
	var mods []string

	esoHome, ok := os.LookupEnv("ESO_HOME")
	if !ok {
		err = fmt.Errorf("please set the ESO_HOME environment variable and try again")
		return nil, err
	}

	var addonsPath = filepath.Join(filepath.Clean(esoHome), "live", "AddOns")

	if ok, err := afero.DirExists(AppFs, addonsPath); !ok || err != nil {
		return nil, fmt.Errorf("%+q is not a valid ESO HOME directory", esoHome)
	}

	fmt.Println("Searching", addonsPath)

	err = afero.Walk(AppFs, addonsPath, func(path string, info fs.FileInfo, err error) error { return getModList(path, &mods, err) })
	if err != nil {
		return mods, err
	}

	fmt.Println("Found", len(mods), "mods")

	return mods, err
}

func getModList(path string, a *[]string, err error) error {
	if err != nil {
		return err
	}

	filename := filepath.Base(path)
	directory := filepath.Base(filepath.Dir(path))

	if filepath.Ext(filename) == ".txt" && directory+".txt" == filename {
		*a = append(*a, path)
		// fmt.Println(path)
	}

	return nil
}
