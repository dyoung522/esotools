/*
Copyright © 2024 Donovan C. Young <dyoung522@gmail.com>
*/
package cmd

import (
	"os"

	sub1 "github.com/dyoung522/esotools/cmd/list"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "esotools",
	Version: "0.1.0",
	Short:   "tools used to list, install, and validate ESO AddOns",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbosity, _ := cmd.Flags().GetCount("verbose")
		viper.Set("verbosity", verbosity)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	cc.Init(&cc.Config{
		RootCmd:  RootCmd,
		Headings: cc.HiCyan + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})

	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// var Verbose int

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.esotools.yaml)")
	RootCmd.PersistentFlags().StringP("esohome", "H", "", "The full installation path of your ESO game files (where the `live` folder lives).")
	RootCmd.PersistentFlags().CountP("verbose", "v", "counted verbosity")

	viper.BindPFlag("eso_home", RootCmd.PersistentFlags().Lookup("esohome"))

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
	viper.ReadInConfig()
}
