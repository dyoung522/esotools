//go:build windows

package ostools

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

	docs, _, err := k.GetStringValue(`Personal`)
	if err != nil {
		return "", err
	}

	if docs == "" {
		return "", nil
	}

	return filepath.Clean(docs), nil
}
