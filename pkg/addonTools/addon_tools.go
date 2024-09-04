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

var ESOHOME string
var AddOnsPath string
var SavedVarsPath string

func Run() (esoAddOns.AddOns, []error) {
	var AppFs = afero.NewReadOnlyFs(afero.NewOsFs())
	var verbosity = viper.GetInt("verbosity")

	ESOHOME = string(viper.GetString("eso_home"))

	if ESOHOME == "" {
		fmt.Println(errors.New("please set the ESO_HOME environment variable and try again"))
		os.Exit(1)
	}

	AddOnsPath = filepath.Join(filepath.Clean(string(ESOHOME)), "live", "AddOns")
	SavedVarsPath = filepath.Join(filepath.Clean(string(ESOHOME)), "live", "SavedVariables")

	if verbosity >= 2 {
		fmt.Printf("Searching for AddOns in %q\n", AddOnsPath)
		fmt.Printf("Searching for SavedVariables in %q\n", SavedVarsPath)
	}

	return GetAddOns(AppFs)
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
