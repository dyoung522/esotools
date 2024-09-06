package addonTools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dyoung522/esotools/pkg/esoAddOns"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func Run() (esoAddOns.AddOns, []error) {
	var AppFs = afero.NewReadOnlyFs(afero.NewOsFs())
	var verbosity = viper.GetInt("verbosity")
	var esohome = string(viper.GetString("esohome"))

	if esohome == "" {
		fmt.Println(errors.New("please set the ESO_HOME environment variable and try again"))
		os.Exit(1)
	}

	if verbosity >= 2 {
		fmt.Printf("Searching for AddOns in %q\n", AddOnsPath())
		fmt.Printf("Searching for SavedVariables in %q\n", SavedVariablesPath())
	}

	if verbosity >= 1 {
		fmt.Println("Building a list of addons and their dependencies... please wait...")
	}

	return GetAddOns(AppFs)
}

func AddOnsPath() string {
	var esohome = string(viper.GetString("esohome"))

	if esohome == "" {
		panic("ESO_HOME environment variable not set")
	}

	return filepath.Join(filepath.Clean(string(esohome)), "live", "AddOns")
}

func SavedVariablesPath() string {
	var esohome = string(viper.GetString("esohome"))

	if esohome == "" {
		panic("ESO_HOME environment variable not set")
	}

	return filepath.Join(filepath.Clean(string(esohome)), "live", "SavedVariables")
}

func Pluralize(s string, c int) string {
	var pluralize = pluralize.NewClient()

	if c == 1 {
		return s
	}

	return pluralize.Plural(s)
}
