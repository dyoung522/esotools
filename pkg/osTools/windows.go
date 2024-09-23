//go:build windows

package osTools

import (
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func DocumentsDir() (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	winDocsDir, _, err := k.GetStringValue("Personal")
	if err != nil {
		return "", err
	}

	if winDocsDir == "" {
		return "", nil
	}

	return filepath.Clean(winDocsDir), nil
}
