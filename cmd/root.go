/*
Copyright © 2024 Donovan C. Young <dyoung522@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	sub1 "github.com/dyoung522/esotools/cmd/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "esotools",
	Version: "0.1.0",
	Short:   "toosl used to list, install, and validate ESO mods",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {

	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var Verbose bool

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.esotools.yaml)")

	RootCmd.PersistentFlags().StringP("esohome", "H", "", "The full installation path of your ESO game files (where the `live` folder lives).")
	viper.BindPFlag("eso_home", RootCmd.PersistentFlags().Lookup("esohome"))

	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Add subcommands
	RootCmd.AddCommand(sub1.ListCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".esotools" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".esotools")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

}
