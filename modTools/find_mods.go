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
		fmt.Println("Please set the ESO_HOME environment variable and try again.")
		os.Exit(1)
	}

	addonsPath := filepath.Join(filepath.Clean(esoHome), "live", "AddOns")

	fmt.Println("Searching", addonsPath)

	absPath, err := filepath.Abs(addonsPath)
	if err != nil {
		return mods, err
	}

	err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error { return getModList(path, &mods, err) })
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
