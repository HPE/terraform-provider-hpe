// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_host_profile/example.tf example.tf.tmpl Name "Example Host Profile" Description "An example host profile" IscsiTransport "tcp" IscsiChapEnabled "false" MultipathFindMultipaths "true" MultipathUserFriendlyNames "true" CpuPowerManagement "performance" LoggingLevel "info" AdditionalPackages "openvswitch=2.17.0"

func RenderHostProfileConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                       "Example Host Profile",
		"Description":                "An example host profile",
		"IscsiTransport":             "tcp",
		"IscsiChapEnabled":           "false",
		"MultipathFindMultipaths":    "true",
		"MultipathUserFriendlyNames": "true",
		"CpuPowerManagement":         "performance",
		"LoggingLevel":               "info",
		"AdditionalPackages":         "openvswitch=2.17.0",
	}

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
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
