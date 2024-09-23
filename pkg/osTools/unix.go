//go:build !windows

package osTools

import (
	"os"
	"path/filepath"
)

func DocumentsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if home == "" {
		return "", nil
	}

	return filepath.Join(home, "Documents"), nil
}
