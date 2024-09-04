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

func FindAddOns(AppFs afero.Fs) ([]esoAddOnFiles.AddOnDefinition, error) {
	var err error
	var addons []esoAddOnFiles.AddOnDefinition

	verbosity := viper.GetInt("verbosity")
	esohome := string(viper.GetString("eso_home"))
	addonsPath := AddOnsPath()

	if ok, err := afero.DirExists(AppFs, addonsPath); !ok || err != nil {
		return nil, fmt.Errorf("%+q is not a valid ESO HOME directory", esohome)
	}

	if verbosity >= 2 {
		fmt.Println("Searching", addonsPath)
	}

	err = afero.Walk(AppFs, addonsPath, func(path string, info fs.FileInfo, err error) error { return getAddOnList(path, &addons, err) })
	if err != nil {
		return nil, fmt.Errorf("error occurred while walking %q: %w", addonsPath, err)
	}

	if verbosity >= 2 {
		fmt.Println("Found", len(addons), "AddOn directories")
	}

	return addons, err
}

func getAddOnList(path string, addons *[]esoAddOnFiles.AddOnDefinition, err error) error {
	var verbosity = viper.GetInt("verbosity")

	if err != nil {
		return err
	}

	if verbosity >= 5 {
		fmt.Println("Searching", path)
	}

	md := esoAddOnFiles.AddOnDefinition{
		Name: filepath.Base(path),
		Dir:  strings.TrimPrefix(filepath.Dir(path), AddOnsPath()),
	}

	if filepath.Ext(md.Name) == ".txt" && esoAddOns.ToKey(filepath.Base(md.Dir)) == md.Key() {
		if verbosity >= 3 {
			fmt.Println("Found", md.Name)
		}

		*addons = append(*addons, md)
	}

	return nil
}
