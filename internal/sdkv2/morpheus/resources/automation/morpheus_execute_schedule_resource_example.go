// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package automation

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_execute_schedule/resource.tf hpe_morpheus_execute_schedule_resource.tf.tmpl Name "Run daily at 7 AM" Description "This schedule runs daily at 7 AM Mountain Time" Enabled false TimeZone "America/Denver" Schedule "7 0 * * *"

func RenderExecuteScheduleConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        "Example",
		"Description": "This schedule runs daily at 7 AM Mountain Time",
		"Enabled":     "false",
		"TimeZone":    "America/Denver",
		"Schedule":    "7 0 * * *",
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
	templatePath := filepath.Join(dir, "hpe_morpheus_execute_schedule_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
