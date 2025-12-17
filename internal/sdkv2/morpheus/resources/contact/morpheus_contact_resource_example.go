// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package contact

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_contact/resource.tf contact_resource.tf.tmpl Name tfcontactdemo EmailAddress tfcontact@demo.com MobileNumber 123-456-7890

// RenderContactConfig renders the contact resource configuration with
// optional field overrides
func RenderContactConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         name,
		"EmailAddress": "tfcontact@demo.com",
		"MobileNumber": "123-456-7890",
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
	templatePath := filepath.Join(dir, "contact_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"EmailAddress", defaults["EmailAddress"],
		"MobileNumber", defaults["MobileNumber"],
	)
}
