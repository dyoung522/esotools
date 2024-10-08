package eso_test

import (
	"path/filepath"
	"testing"

	"github.com/dyoung522/esotools/lib/eso"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var AppFs afero.Fs

func init() {
	AppFs = afero.NewMemMapFs()
	_ = AppFs.MkdirAll("/tmp/eso/Elder Scrolls Online/live/AddOns", 0755)

	viper.Set("eso_home", "/tmp/eso")
}

func TestGetAddOns_EmptyAddonList(t *testing.T) {
	// Arrange
	expected := eso.AddOns{}
	expectedErrs := []error{}

	// Act
	actual, actualErrs := eso.GetAddOns(AppFs)

	// Assert
	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedErrs, actualErrs)
}

func TestGetAddOns_MissingRequiredTitle(t *testing.T) {
	// Arrange
	addonName := "MissingRequiredTitle"
	data := []byte(`## Description: My Addon
## Author: Author Name
## Contributors: Contributor 1, Contributor 2
## AddOnVersion: 1.0.0
## APIVersion: 10001
## SavedVariables: MyAddon_SV
## DependsOn:
## OptionalDependsOn:
## IsLibrary: false`)

	err := afero.WriteFile(AppFs, filepath.Join(eso.AddOnsPath(), addonName, addonName+".txt"), data, 0644)
	if err != nil {
		t.Fatal(err)
	}

	expectedTitle := "Missing Required Title"

	// Act
	addons, actualErrs := eso.GetAddOns(AppFs)
	require.Empty(t, actualErrs)
	require.Len(t, addons, 1)

	// Assert
	addon, found := addons.Find(addonName)
	require.True(t, found)

	assert.Equal(t, expectedTitle, addon.Title)
}
