package esoAddOnFiles

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dyoung522/esotools/pkg/esoAddOns"
)

// AddOnDefinition represents a single ESO add-on definition.
type AddOnDefinition struct {
	Name string // Name of the add-on file (without the .txt extension).
	Dir  string // Directory where the add-on file is located.
}

// String returns a string representation of the AddOnDefinition.
// The format is "[ID: <key>, name: <Name>, dir: <Dir>]".
func (AD AddOnDefinition) String() string {
	return fmt.Sprintf("[ID: %s, name: %s, dir: %s]", AD.Key(), AD.Name, AD.Dir)
}

// Path returns the full path to the add-on file.
// It joins the directory (Dir) and the name (with .txt extension) of the add-on file.
func (AD AddOnDefinition) Path() string {
	return filepath.Join(AD.Dir, AD.Name)
}

// Key returns the unique identifier of the add-on.
// It trims the .txt extension from the name and converts it to a key using the ToKey function from the esoAddOns package.
func (AD AddOnDefinition) Key() string {
	return esoAddOns.ToKey(strings.TrimSuffix(AD.Name, ".txt"))
}
