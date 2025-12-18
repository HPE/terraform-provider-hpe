// (C) Copyright 2025 Hewlett Packard Enterprise Development LP
package user

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../../cmd/render -out examples/resources/morpheus_user/example.tf example.tf.tmpl TenantId 1 Username "example-user" RoleIds 1 LinuxKeyPairId 100

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
		"TenantId", defaults["TenantId"],
		"Username", defaults["Username"],
		"RoleIds", defaults["RoleIds"],
		"LinuxKeyPairId", defaults["LinuxKeyPairId"],
	)
}
