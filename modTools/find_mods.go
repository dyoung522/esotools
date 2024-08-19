package modTools

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	esoMods "github.com/dyoung522/esotools/esoMods"
	esoModFiles "github.com/dyoung522/esotools/eso_mod_files"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var verbose = viper.GetBool("verbose")

func FindMods() ([]esoModFiles.ModDefinition, error) {
	var err error
	var mods []esoModFiles.ModDefinition

	if ok, err := afero.DirExists(AppFs, AddOnsPath); !ok || err != nil {
		return nil, fmt.Errorf("%+q is not a valid ESO HOME directory", ESOHOME)
	}

	if verbose {
		fmt.Println("Searching", AddOnsPath)
	}

	err = afero.Walk(AppFs, AddOnsPath, func(path string, info fs.FileInfo, err error) error { return getModList(path, &mods, err) })
	if err != nil {
		return nil, fmt.Errorf("error occurred while walking %q: %w", AddOnsPath, err)
	}

	if verbose {
		fmt.Println("Found", len(mods), "mod directories")
	}

	return mods, err
}

func getModList(path string, mods *[]esoModFiles.ModDefinition, err error) error {
	if err != nil {
		return err
	}

	md := esoModFiles.ModDefinition{
		Name: filepath.Base(path),
		Dir:  strings.TrimPrefix(strings.TrimPrefix(filepath.Dir(path), AddOnsPath), "/"),
	}

	if filepath.Ext(md.Name) == ".txt" && esoMods.ToKey(filepath.Base(md.Dir)) == md.Key() {
		*mods = append(*mods, md)
	}

	return nil
}
