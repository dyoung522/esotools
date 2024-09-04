package addonTools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dyoung522/esotools/pkg/esoAddOns"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func Run() (esoAddOns.AddOns, []error) {
	var AppFs = afero.NewReadOnlyFs(afero.NewOsFs())
	var verbosity = viper.GetInt("verbosity")
	var esohome = string(viper.GetString("eso_home"))

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
	var esohome = string(viper.GetString("eso_home"))

	return filepath.Join(filepath.Clean(string(esohome)), "live", "AddOns")
}

func SavedVariablesPath() string {
	var esohome = string(viper.GetString("eso_home"))

	return filepath.Join(filepath.Clean(string(esohome)), "live", "SavedVariables")
}

func Pluralize(s string, c int) string {
	if c == 1 {
		return s
	}

	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}

	return s + "s"
}
