// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package power_schedule

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/hpe_morpheus_power_schedule/example.tf example.tf.tmpl Name "Business Hours" Description "Power on during business hours" ScheduleType "power" ScheduleTimezone "America/New_York" Enabled "true" MondayOnTime "08:00" MondayOffTime "18:00" TuesdayOnTime "08:00" TuesdayOffTime "18:00"

func RenderPowerScheduleConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":             "Business Hours",
		"Description":      "Power on during business hours",
		"ScheduleType":     "power",
		"ScheduleTimezone": "America/New_York",
		"Enabled":          "true",
		"MondayOnTime":     "08:00",
		"MondayOffTime":    "18:00",
		"TuesdayOnTime":    "08:00",
		"TuesdayOffTime":   "18:00",
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
