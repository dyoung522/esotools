package addonTools

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/dyoung522/esotools/pkg/esoAddOnFiles"
	"github.com/dyoung522/esotools/pkg/esoAddOns"
	"github.com/spf13/afero"
)

func FindAddOns(AppFs afero.Fs) ([]esoAddOnFiles.AddOnDefinition, error) {
	var err error
	var addons []esoAddOnFiles.AddOnDefinition

	if ok, err := afero.DirExists(AppFs, AddOnsPath); !ok || err != nil {
		return nil, fmt.Errorf("%+q is not a valid ESO HOME directory", ESOHOME)
	}

	if verbosity >= 1 {
		fmt.Println("Searching", AddOnsPath)
	}

	err = afero.Walk(AppFs, AddOnsPath, func(path string, info fs.FileInfo, err error) error { return getAddOnList(path, &addons, err) })
	if err != nil {
		return nil, fmt.Errorf("error occurred while walking %q: %w", AddOnsPath, err)
	}

	if verbosity >= 1 {
		fmt.Println("Found", len(addons), "AddOn directories")
	}

	return addons, err
}

func getAddOnList(path string, addons *[]esoAddOnFiles.AddOnDefinition, err error) error {
	if err != nil {
		return err
	}

	if verbosity >= 5 {
		fmt.Println("Searching", path)
	}

	md := esoAddOnFiles.AddOnDefinition{
		Name: filepath.Base(path),
		Dir:  strings.TrimPrefix(filepath.Dir(path), AddOnsPath),
	}

	if filepath.Ext(md.Name) == ".txt" && esoAddOns.ToKey(filepath.Base(md.Dir)) == md.Key() {
		if verbosity >= 3 {
			fmt.Println("Found", md.Name)
		}

		*addons = append(*addons, md)
	}

	return nil
}
