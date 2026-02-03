// (C) Copyright 2025 Hewlett Packard Enterprise Development LP
package user

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../../bin/render -out examples/resources/morpheus_user/example.tf example.tf.tmpl TenantId 1 Username "example-user" RoleIds 1 LinuxKeyPairId 100

func RenderUserConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"TenantId":       "1",
		"Username":       "example-user",
		"RoleIds":        "1",
		"LinuxKeyPairId": "100",
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
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
