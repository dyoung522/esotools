package modTools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dyoung522/esotools/pkg/esoMods"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var AppFs afero.Fs
var ESOHOME string
var AddOnsPath string

func Run() (esoMods.Mods, []error) {
	ESOHOME = string(viper.GetString("eso_home"))
	if ESOHOME == "" {
		fmt.Println(errors.New("please set the ESO_HOME environment variable and try again"))
		os.Exit(1)
	}

	AddOnsPath = filepath.Join(filepath.Clean(string(ESOHOME)), "live", "AddOns")

	verbose := viper.GetBool("verbose")

	if verbose {
		fmt.Printf("Getting list of mods\n\n")
	}

	return GetMods()
}

func init() {
	AppFs = afero.NewReadOnlyFs(afero.NewOsFs())
}
