// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_host/example-id.tf example-id.tf.tmpl Id 1
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_host/example-name.tf example-name.tf.tmpl Name "Example host" CloudId 17

func RenderHostByIdConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderHost(t, "example-id.tf.tmpl", map[string]string{
		"Id": "1",
	}, overrides)
}

func RenderHostByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderHost(t, "example-name.tf.tmpl", map[string]string{
		"Name":    "Example host",
		"CloudId": "17",
	}, overrides)
}

func renderHost(
	t *testing.T,
	tmpl string,
	defaults, overrides map[string]string,
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
	templatePath := filepath.Join(dir, tmpl)

	return testhelpers.RenderExample(t, templatePath, args...)
}
