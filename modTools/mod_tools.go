package modTools

import (
	"github.com/spf13/afero"
)

var AppFs afero.Fs

func init() {
	AppFs = afero.NewReadOnlyFs(afero.NewOsFs())
}
