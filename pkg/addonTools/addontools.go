package addOnTools

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dyoung522/esotools/lib/esoAddOns"
	"github.com/dyoung522/esotools/pkg/osTools"
	"github.com/gertd/go-pluralize"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var AppFs = afero.NewReadOnlyFs(afero.NewOsFs())

func Run() (esoAddOns.AddOns, []error) {
	var verbosity = viper.GetInt("verbosity")

	if verbosity >= 1 {
		fmt.Println("Building a list of addons and their dependencies... please wait...")
	}

	return GetAddOns(AppFs)
}

func AddOnsPath() string {
	return filepath.Join(filepath.Clean(ESOHome()), "live", "AddOns")
}

func SavedVariablesPath() string {
	return filepath.Join(filepath.Clean(ESOHome()), "live", "SavedVariables")
}

func Pluralize(s string, c int) string {
	var pluralize = pluralize.NewClient()

	if c == 1 {
		return s
	}

	return pluralize.Plural(s)
}

func ValidateESOHOME() error {
	verbosity := viper.GetInt("verbosity")
	esoHome := ESOHome()

	if !checkESODir(esoHome) {
		if esoHome != "" {
			fmt.Println(fmt.Errorf("%q does not appear to be a valid ESO directory, attempting auto-detect", esoHome))
		}

		documentsDir, err := osTools.DocumentsDir()
		if err != nil {
			return err
		}

		esoHome = esoDir(documentsDir)

		for !checkESODir(esoHome) {
			fmt.Println(fmt.Errorf("%q is not a valid ESO directory\n", esoHome))

			esoHome, err = pterm.DefaultInteractiveTextInput.WithDefaultValue(documentsDir).Show(`Enter the directory where your "Elder Scrools Online" documents folder lives [CTRL+C to exit]`)
			if err != nil {
				return err
			}

			esoHome = esoDir(esoHome)
		}

		if verbosity >= 2 {
			fmt.Printf("ESO_HOME set to %q\n", esoHome)
		}

		viper.Set("eso_home", string(esoHome))
	}

	return nil
}

func ESOHome() string {
	return esoDir(viper.GetString("eso_home"))
}

func esoDir(dir string) string {
	dir = filepath.Clean(strings.TrimSpace(dir))

	if dir == "" || dir == "." || dir == "/" || dir == `\\` || dir == `C:\` {
		return ""
	}

	// If the user entered a SavedVariables or AddOns directory, strip it off
	for strings.Contains(dir, "live") {
		dir = filepath.Dir(dir)
	}

	// Add the "Elder Scrolls Online" directory if it's not there
	if !strings.Contains(dir, "Elder Scrolls Online") {
		dir = filepath.Join(dir, "Elder Scrolls Online")
	}

	return dir
}

func checkESODir(dir string) bool {
	dir = esoDir(dir)

	ok, err := afero.DirExists(AppFs, filepath.Join(dir, "live"))
	if err != nil {
		return false
	}

	return ok
}
