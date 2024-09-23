package addOnTools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dyoung522/esotools/lib/esoAddOns"
	"github.com/dyoung522/esotools/pkg/osTools"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var AppFs = afero.NewReadOnlyFs(afero.NewOsFs())

func Run() (esoAddOns.AddOns, []error) {
	var verbosity = viper.GetInt("verbosity")
	err := validateESOHOME()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if verbosity >= 2 {
		fmt.Printf("Using ESO_HOME: %q\n", filepath.Clean(viper.GetString("eso_home")))
		fmt.Printf("Searching for AddOns in %q\n", AddOnsPath())
		fmt.Printf("Searching for SavedVariables in %q\n", SavedVariablesPath())
	}

	if verbosity >= 1 {
		fmt.Println("Building a list of addons and their dependencies... please wait...")
	}

	return GetAddOns(AppFs)
}

func AddOnsPath() string {
	return filepath.Join(filepath.Clean(esoHome()), "live", "AddOns")
}

func SavedVariablesPath() string {
	return filepath.Join(filepath.Clean(esoHome()), "live", "SavedVariables")
}

func Pluralize(s string, c int) string {
	var pluralize = pluralize.NewClient()

	if c == 1 {
		return s
	}

	return pluralize.Plural(s)
}

func validateESOHOME() error {
	var err error
	verbosity := viper.GetInt("verbosity")
	eso_home := viper.GetString("eso_home")

	if eso_home == "" {
		documentsDir, err := osTools.DocumentsDir()

		if err != nil {
			return err
		}

		if documentsDir == "" {
			return errors.New("Cannot determine the location of the ESO_HOME directory, please set it manually")
		}

		eso_home = filepath.Join(documentsDir, "Elder Scrolls Online")

		if verbosity > 1 {
			fmt.Printf("Autodiscovered ESO_HOME at %q\n", eso_home)
		}
	}

	ok, err := afero.DirExists(AppFs, filepath.Join(eso_home, "live"))
	if !ok || err != nil {
		return fmt.Errorf("%q does not appear to be a valid ESO_HOME; please set the ESO_HOME environment variable and try again", eso_home)
	}

	viper.Set("eso_home", string(eso_home))

	return nil
}

func esoHome() string {
	var eso_home string = viper.GetString("eso_home")

	if eso_home == "" {
		eso_home = "."
	}

	return eso_home
}
