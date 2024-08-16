package main

import (
	"fmt"
	"os"

	"github.com/dyoung522/esotools/modTools"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.BindEnv("ESO_HOME")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	ESOHOME := string(viper.GetString("ESO_HOME"))
	if ESOHOME == "" {
		fmt.Println(fmt.Errorf("please set the ESO_HOME environment variable and try again"))
		os.Exit(1)
	}
}

func main() {
	modlist, err := modTools.FindMods()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mods, errs := modTools.ReadMods(&modlist)
	_ = mods

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		os.Exit(1)
	}

	fmt.Println(mods.Print())
}
