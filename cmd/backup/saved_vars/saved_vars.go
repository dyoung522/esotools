package cmd

import (
	"archive/zip"
	"fmt"
	"time"

	esoSavedVar "github.com/dyoung522/esotools/eso/savedvar"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var BackupSavedVarsCmd = &cobra.Command{
	Use:   "savedvars",
	Short: "Create a ZIP backup file of all SavedVariables",
	Long:  `Creates a ZIP backup file of all SavedVariables in the current directory.`,
	Run:   execute,
}

func execute(cmd *cobra.Command, args []string) {
	AppFs := afero.NewOsFs()
	if err := BackupSavedVars(AppFs); err != nil {
		cmd.Println(err)
	}
}

func BackupSavedVars(AppFs afero.Fs) error {
	var err error
	verbosity := viper.GetInt("verbosity")

	t := time.Now()
	archiveTime := fmt.Sprintf("%d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	archiveFileName := fmt.Sprintf("saved_variables_%s.zip", archiveTime)

	saveVarFiles, err := esoSavedVar.Find(AppFs)
	if err != nil {
		fmt.Println(err)
		return err
	}

	archiveFile, err := AppFs.Create(archiveFileName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer archiveFile.Close()

	if verbosity >= 1 {
		fmt.Printf("Backing up SavedVariables to %s\n", archiveFileName)
	}

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, sv := range saveVarFiles {
		if verbosity >= 2 {
			fmt.Printf("Adding %s to %s\n", sv.Name(), archiveFileName)
		}

		fileData, err := afero.ReadFile(AppFs, sv.Path())
		if err != nil {
			fmt.Println(err)
			return err
		}

		zipFile, err := zipWriter.Create(sv.Name())
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, err = zipFile.Write(fileData)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
