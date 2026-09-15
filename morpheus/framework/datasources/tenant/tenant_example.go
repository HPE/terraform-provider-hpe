// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_tenant/data-source.tf example-name.tf.tmpl Name "tenant name"

func RenderExample(t *testing.T, tmpl string, overrides map[string]string) (string, error) {
	t.Helper()

	var args []string
	for key, value := range overrides {
		args = append(args, key, value)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}

	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, tmpl)

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
