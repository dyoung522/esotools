package addonTools

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/dyoung522/esotools/pkg/esoAddOnFiles"
	"github.com/dyoung522/esotools/pkg/esoAddOns"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var verbose = viper.GetBool("verbose")

func FindAddOns() ([]esoAddOnFiles.AddOnDefinition, error) {
	var err error
	var addons []esoAddOnFiles.AddOnDefinition

	if ok, err := afero.DirExists(AppFs, AddOnsPath); !ok || err != nil {
		return nil, fmt.Errorf("%+q is not a valid ESO HOME directory", ESOHOME)
	}

	if verbose {
		fmt.Println("Searching", AddOnsPath)
	}

	err = afero.Walk(AppFs, AddOnsPath, func(path string, info fs.FileInfo, err error) error { return getAddOnList(path, &addons, err) })
	if err != nil {
		return nil, fmt.Errorf("error occurred while walking %q: %w", AddOnsPath, err)
	}

	if verbose {
		fmt.Println("Found", len(addons), "AddOn directories")
	}

	return addons, err
}

func getAddOnList(path string, addons *[]esoAddOnFiles.AddOnDefinition, err error) error {
	if err != nil {
		return err
	}

	md := esoAddOnFiles.AddOnDefinition{
		Name: filepath.Base(path),
		Dir:  strings.TrimPrefix(strings.TrimPrefix(filepath.Dir(path), AddOnsPath), "/"),
	}

	if filepath.Ext(md.Name) == ".txt" && esoAddOns.ToKey(filepath.Base(md.Dir)) == md.Key() {
		*addons = append(*addons, md)
	}

	return nil
}
