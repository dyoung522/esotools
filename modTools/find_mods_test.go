package modTools

import (
	"os"
	"testing"
)

func TestFindMods(t *testing.T) {
	os.Setenv("ESO_HOME", "")

	_, err := FindMods()

	if err == nil {
		t.Error("expected error message, not nil")
	}
}
