//go:build !windows

package osTools

import (
	"os"
	"path/filepath"
)

func DocumentsDir() (string, error) {
	eso_home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if eso_home == "" {
		return "", nil
	}

	return filepath.Join(eso_home, "Documents"), nil
}
