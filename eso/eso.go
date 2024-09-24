package eso

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dyoung522/esotools/ostools"
	"github.com/gertd/go-pluralize"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func Home() string {
	return esoDir(viper.GetString("eso_home"))
}

func ValidateESOHOME() error {
	verbosity := viper.GetInt("verbosity")
	esoHome := Home()

	if !checkESODir(esoHome) {
		if esoHome != "" {
			fmt.Println(fmt.Errorf("%q does not appear to be a valid ESO directory, attempting auto-detect", esoHome))
		}

		documentsDir, err := ostools.DocumentsDir()
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

// StripESOColorCodes removes ESO color codes from a string.
func StripESOColorCodes(input string) string {
	if !strings.Contains(input, `|c`) {
		return input
	}

	colorRE := regexp.MustCompile(`\|c[[:xdigit:]]{6}`)
	re := regexp.MustCompile(`\|c[[:xdigit:]]{6}(.*?)(?:\|r)`)

	cleanString := re.ReplaceAllStringFunc(input, func(match string) string {
		parts := re.FindStringSubmatch(match)
		return parts[1]
	})

	// Strip out any remaining color codes from the clean title.
	return colorRE.ReplaceAllString(cleanString, "")
}

// Removes version dependencies and returns the plain dependency name
func DependencyName(input string) []string {
	return strings.Split(strings.TrimRight(input, "\r\n"), ">=")
}

func Pluralize(s string, c int) string {
	var pluralize = pluralize.NewClient()

	if c == 1 {
		return s
	}

	return pluralize.Plural(s)
}

/*
 * Private Functions
 */

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
	AppFs := afero.NewReadOnlyFs(afero.NewOsFs())
	dir = esoDir(dir)

	ok, err := afero.DirExists(AppFs, filepath.Join(dir, "live"))
	if err != nil {
		return false
	}

	return ok
}
