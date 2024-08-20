package esoModFiles

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dyoung522/esotools/pkg/esoMods"
)

type ModDefinition struct {
	Name string
	Dir  string
}

func (MD ModDefinition) String() string {
	return fmt.Sprintf("[ID: %s, name: %s, dir: %s]", MD.Key(), MD.Name, MD.Dir)
}

func (MD ModDefinition) Path() string {
	return filepath.Join(MD.Dir, MD.Name)
}

func (MD ModDefinition) Key() string {
	return esoMods.ToKey(strings.TrimSuffix(MD.Name, ".txt"))
}
