package addonTools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dyoung522/esotools/pkg/esoAddOns"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var AppFs afero.Fs
var ESOHOME string
var AddOnsPath string

func Run() (esoAddOns.AddOns, []error) {
	ESOHOME = string(viper.GetString("eso_home"))

	if ESOHOME == "" {
		fmt.Println(errors.New("please set the ESO_HOME environment variable and try again"))
		os.Exit(1)
	}

	AddOnsPath = filepath.Join(filepath.Clean(string(ESOHOME)), "live", "AddOns")

	if viper.GetInt("verbosity") >= 1 {
		fmt.Printf("Searching for AddOns in %q\n\n", AddOnsPath)
	}

	return GetAddOns()
}

func init() {
	AppFs = afero.NewReadOnlyFs(afero.NewOsFs())
}
