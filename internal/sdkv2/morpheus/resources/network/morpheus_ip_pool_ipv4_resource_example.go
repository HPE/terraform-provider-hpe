// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_ip_pool_ipv4/resource.tf morpheus_ip_pool_ipv4_resource.tf.tmpl Name "\"Terraform Example IPv4 IP pool\"" StartingAddress1 "\"192.168.1.1\"" EndingAddress1 "\"192.168.1.10\"" StartingAddress2 "\"10.0.0.1\"" EndingAddress2 "\"10.0.0.10\""

// RenderIPPoolIPv4Config generates a Terraform configuration for hpe_morpheus_ip_pool_ipv4 resource.
// It accepts an optional map of field overrides. If nil or empty, default values are used.
func RenderIPPoolIPv4Config(t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":             name,
		"StartingAddress1": "\"192.168.1.1\"",
		"EndingAddress1":   "\"192.168.1.10\"",
		"StartingAddress2": "\"10.0.0.1\"",
		"EndingAddress2":   "\"10.0.0.10\"",
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
	templatePath := filepath.Join(dir, "morpheus_ip_pool_ipv4_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name",
		defaults["Name"],
		"StartingAddress1",
		defaults["StartingAddress1"],
		"EndingAddress1",
		defaults["EndingAddress1"],
		"StartingAddress2",
		defaults["StartingAddress2"],
		"EndingAddress2",
		defaults["EndingAddress2"],
	)
}
