package esoAddOnFiles

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dyoung522/esotools/pkg/esoAddOns"
)

type AddOnDefinition struct {
	Name string
	Dir  string
}

func (AD AddOnDefinition) String() string {
	return fmt.Sprintf("[ID: %s, name: %s, dir: %s]", AD.Key(), AD.Name, AD.Dir)
}

func (AD AddOnDefinition) Path() string {
	return filepath.Join(AD.Dir, AD.Name)
}

func (AD AddOnDefinition) Key() string {
	return esoAddOns.ToKey(strings.TrimSuffix(AD.Name, ".txt"))
}
