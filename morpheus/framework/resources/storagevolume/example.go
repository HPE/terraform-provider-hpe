// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_storage_volume/example.tf example.tf.tmpl Name "Example Storage Volume" TypeCode hpealletraMPLUN StorageServerId 1
//go:generate ../../../../bin/render -out examples/resources/morpheus_storage_volume/example_alletramp_bmaas.tf example_alletramp_bmaas.tf.tmpl Name "Example Alletra MP BMaaS Volume" TypeCode hpealletraMPLUN WoVersion 1
//go:generate ../../../../bin/render -out examples/resources/morpheus_storage_volume/example_config.tf example_config.tf.tmpl Name "Example Storage Volume" TypeCode hpealletraMPLUN
//go:generate ../../../../bin/render -out examples/resources/morpheus_storage_volume/example_complete.tf example_complete.tf.tmpl Name "Example Storage Volume" TypeCode hpealletraMPLUN StorageServerId 1 MaxStorage 10

func RenderStorageVolumeConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderStorageVolumeExample(t, "example.tf.tmpl", map[string]string{
		"Name":            "Example Storage Volume",
		"TypeCode":        "hpealletraMPLUN",
		"StorageServerId": "1",
	}, overrides)
}

// RenderStorageVolumeCompleteConfig renders a storage volume with a size and a
// storage server, exercising the max_storage (GiB) path.
func RenderStorageVolumeCompleteConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderStorageVolumeExample(t, "example_complete.tf.tmpl", map[string]string{
		"Name":            "Example Storage Volume",
		"TypeCode":        "hpealletraMPLUN",
		"StorageServerId": "1",
		"MaxStorage":      "10",
	}, overrides)
}

// RenderStorageVolumeAlletraMPBMaaSConfig renders a storage volume with the
// write-only config_alletramp_bmaas block, exercising buildCreateConfig and the
// config_alletramp_bmaas_wo_version replace trigger. Override WoVersion to plan a
// replacement.
func RenderStorageVolumeAlletraMPBMaaSConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderStorageVolumeExample(t, "example_alletramp_bmaas.tf.tmpl", map[string]string{
		"Name":      "Example Alletra MP BMaaS Volume",
		"TypeCode":  "hpealletraMPLUN",
		"WoVersion": "1",
	}, overrides)
}

func renderStorageVolumeExample(
	t *testing.T,
	templateName string,
	defaults map[string]string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, templateName)

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
