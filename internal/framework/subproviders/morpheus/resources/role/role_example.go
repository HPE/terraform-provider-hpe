package role

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../../cmd/render example.tf.tmpl Name "ExampleRole" Multitenant "false" Description "An example role" RoleType "user"

//go:generate go run ../../../../../../cmd/render example-using-legacy-provider.tf.tmpl TaskDataSourceName "example_legacy_task" TaskName "example_task" ResourceName "example_with_legacy_provider" Name "ExampleRoleWithLegacyProvider" Description "An example role using legacy provider" RoleType "user" Task0Access "full"

//go:embed example.tf.tmpl
var templateExample string

func RenderRoleUserConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        "ExampleRole",
		"Multitenant": "false",
		"Description": "An example user role",
		"RoleType":    "user",
	}

	// Apply overrides to defaults
	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"example.tf.tmpl",
		"Name",
		defaults["Name"],
		"Multitenant",
		defaults["Multitenant"],
		"Description",
		defaults["Description"],
		"RoleType",
		defaults["RoleType"],
	)
}

func RenderRoleTenantConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        "ExampleRole",
		"Description": "An example tenant role",
		"RoleType":    "tenant",
	}

	// Apply overrides to defaults
	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "role_tenant.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name",
		defaults["Name"],
		"Description",
		defaults["Description"],
		"RoleType",
		defaults["RoleType"],
	)
}
