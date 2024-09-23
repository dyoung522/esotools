package eso_test

import (
	"path/filepath"
	"testing"

	"github.com/dyoung522/esotools/lib/eso"
)

func TestAddOnDefinition_Key_TrimSuffix(t *testing.T) {
	// Arrange
	addOnDef := eso.AddOnDefinition{
		Name: "MyAddon.txt",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	expectedKey := "MyAddon"
	if key != expectedKey {
		t.Errorf("expected key '%s' for name '%s', but got: '%s'", expectedKey, addOnDef.Name, key)
	}
}

func TestAddOnDefinition_Key_TrimSuffix_NotEndingInTxt(t *testing.T) {
	// Arrange
	addOnDef := eso.AddOnDefinition{
		Name: "MyAddon",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	expectedKey := "MyAddon"
	if key != expectedKey {
		t.Errorf("expected key '%s' for name '%s', but got: '%s'", expectedKey, addOnDef.Name, key)
	}
}

func TestAddOnDefinition_Path(t *testing.T) {
	// Arrange
	addOnDef := eso.AddOnDefinition{
		Name: "MyAddon.txt",
		Dir:  "path/to/addon",
	}

	// Act
	path := addOnDef.Path()

	// Assert
	expectedPath := filepath.Clean("path/to/addon/MyAddon.txt")
	if path != expectedPath {
		t.Errorf("expected path '%s' for name '%s', but got: '%s'", expectedPath, addOnDef.Name, path)
	}
}

func TestAddOnDefinition_Path_NameNotEndingInTxt(t *testing.T) {
	// Arrange
	addOnDef := eso.AddOnDefinition{
		Name: "MyAddon.txt",
		Dir:  "path/to/addon",
	}

	// Act
	path := addOnDef.Path()

	// Assert
	expectedPath := filepath.Clean("path/to/addon/MyAddon.txt")
	if path != expectedPath {
		t.Errorf("expected path '%s' for name '%s', but got: '%s'", expectedPath, addOnDef.Name, path)
	}
}

func TestAddOnDefinition_Key_EmptyName(t *testing.T) {
	// Arrange
	addOnDef := eso.AddOnDefinition{
		Name: "",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	if key != "" {
		t.Errorf("expected an empty key for an empty name, but got: %s", key)
	}
}

func TestAddOnDefinition_Path_LeadingAndTrailingSlashes(t *testing.T) {
	// Arrange
	addOnDef := eso.AddOnDefinition{
		Name: "MyAddon.txt",
		Dir:  "/path/to/addon/",
	}

	// Act
	path := addOnDef.Path()

	// Assert
	expectedPath := filepath.Clean("/path/to/addon/MyAddon.txt")
	if path != expectedPath {
		t.Errorf("expected path '%s' for name '%s' and dir '%s', but got: '%s'", expectedPath, addOnDef.Name, addOnDef.Dir, path)
	}
}

func TestAddOnDefinition_Path_DirectoryWithSpaces(t *testing.T) {
	// Arrange
	addOnDef := eso.AddOnDefinition{
		Name: "MyAddon.txt",
		Dir:  "/path with spaces/to/addon",
	}

	// Act
	path := addOnDef.Path()

	// Assert
	expectedPath := filepath.Clean("/path with spaces/to/addon/MyAddon.txt")
	if path != expectedPath {
		t.Errorf("expected path '%s' for name '%s' and dir '%s', but got: '%s'", expectedPath, addOnDef.Name, addOnDef.Dir, path)
	}
}

func TestAddOnDefinition_Key_UppercaseLetters(t *testing.T) {
	// Arrange
	addOnDef := eso.AddOnDefinition{
		Name: "MyADDON.txt",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	expectedKey := "MyADDON"
	if key != expectedKey {
		t.Errorf("expected key '%s' for name '%s', but got: '%s'", expectedKey, addOnDef.Name, key)
	}
}

func TestAddOnDefinition_Key_SpecialCharacters(t *testing.T) {
	// Arrange
	addOnDef := eso.AddOnDefinition{
		Name: "My@Addon#1.txt",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	expectedKey := "My@Addon#1"
	if key != expectedKey {
		t.Errorf("expected key '%s' for name '%s', but got: '%s'", expectedKey, addOnDef.Name, key)
	}
}

func TestAddOnDefinition_Key_NonASCIICharacters(t *testing.T) {
	// Arrange
	addOnDef := eso.AddOnDefinition{
		Name: "MyÁddón.txt",
		Dir:  "path/to/addon",
	}

	// Act
	key := addOnDef.Key()

	// Assert
	expectedKey := "MyÁddón"
	if key != expectedKey {
		t.Errorf("expected key '%s' for name '%s', but got: '%s'", expectedKey, addOnDef.Name, key)
	}
}
