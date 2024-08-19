package modTools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var AppFs afero.Fs
var ESOHOME string
var AddOnsPath string

func init() {
	AppFs = afero.NewReadOnlyFs(afero.NewOsFs())

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.BindEnv("ESO_HOME")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(fmt.Errorf("error reading config file, %w", err))
	}

	ESOHOME = string(viper.GetString("ESO_HOME"))
	if ESOHOME == "" {
		fmt.Println(errors.New("please set the ESO_HOME environment variable and try again"))
		os.Exit(1)
	}

	AddOnsPath = filepath.Join(filepath.Clean(string(ESOHOME)), "live", "AddOns")
}
