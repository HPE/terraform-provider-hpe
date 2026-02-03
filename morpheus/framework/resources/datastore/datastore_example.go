// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_datastore/example_alletramp_hvm.tf example_alletramp_hvm.tf.tmpl Name "TestAlletraDatastore" AssociatedResourceID 1 StorageServerID 1 GroupID 1 TenantID 1

func RenderDatastoreAlletraMPHVMConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "TestAlletraDatastore",
		"AssociatedResourceID": "1",
		"StorageServerID":      "1",
		"GroupID":              "1",
		"TenantID":             "1",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "role_user.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
