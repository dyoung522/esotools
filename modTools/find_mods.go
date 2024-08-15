package modTools

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func FindMods() ([]string, error) {
	var err error
	var mods []string

	esoHome, ok := os.LookupEnv("ESO_HOME")
	if !ok {
		err = fmt.Errorf("please set the ESO_HOME environment variable and try again")
		return nil, err
	}

	var addonsPath = filepath.Join(filepath.Clean(esoHome), "live", "AddOns")

	if fileInfo, err := os.Stat(addonsPath); err != nil || !fileInfo.IsDir() {
		return nil, fmt.Errorf("%+q is not a valid ESO HOME directory", esoHome)
	}

	fmt.Println("Searching", addonsPath)

	err = filepath.WalkDir(addonsPath, func(path string, d fs.DirEntry, err error) error { return getModList(path, &mods, err) })
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
