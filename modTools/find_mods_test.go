package modTools

import (
	"testing"
)

func TestFindMods(t *testing.T) {
	// t.Run("ENV-SET", func(t *testing.T) {
	// 	t.Setenv("ESO_HOME", ".")
	// 	_, err := FindMods()
	// 	if err != nil {
	// 		t.Errorf("did not expect error message; got %v", err)
	// 	}
	// })
	t.Run("ENV-NOT-SET", func(t *testing.T) {
		t.Setenv("ESO_HOME", "")
		_, err := FindMods()

		if err == nil {
			t.Error("expected error message, not nil")
		}
	})
}
