// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package job

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_job_task/resource_date_and_time.tf job_task_date_and_time.tf.tmpl Name "TF Example Job Task Date and Time" Enabled true Labels "[\"aws\", \"demo\"]" TaskId 1 ScheduleMode date_and_time ScheduledDateAndTime 2022-12-30T06:00:00Z ContextType instance InstanceIds "[1, 2]"
//go:generate ../../../../../bin/render -out examples/resources/morpheus_job_task/resource_schedule.tf job_task_schedule.tf.tmpl Name "TF Example Job Task Schedule" Enabled true Labels "[\"aws\", \"demo\"]" TaskId 1 ScheduleMode scheduled ExecutionScheduleId 1 ContextType instance InstanceIds "[91]" CustomConfig "{\\\"test\\\":\\\"new\\\"}"
//go:generate ../../../../../bin/render -out examples/resources/morpheus_job_task/resource_manual.tf job_task_manual.tf.tmpl Name "TF Example Job Task Manual" Enabled true Labels "[\"aws\", \"demo\"]" TaskId 1 ScheduleMode manual ContextType instance-label InstanceLabel demo

// RenderJobTaskDateAndTimeConfig generates a Terraform configuration for the job task resource
// with date and time scheduling. It accepts optional overrides for field values.
// Default values are used if not overridden.
func RenderJobTaskDateAndTimeConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "Example Job Task Date and Time",
		"Enabled":              "true",
		"Labels":               "[\"aws\", \"demo\"]",
		"TaskId":               "1",
		"ScheduleMode":         "date_and_time",
		"ScheduledDateAndTime": "2022-12-30T06:00:00Z",
		"ContextType":          "instance",
		"InstanceIds":          "[1, 2]",
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
	templatePath := filepath.Join(dir, "job_task_date_and_time.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderJobTaskScheduleConfig generates a Terraform configuration for the job task resource
// with execution schedule. It accepts optional overrides for field values.
// Default values are used if not overridden.
func RenderJobTaskScheduleConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                "Example Job Task Schedule",
		"Enabled":             "true",
		"Labels":              "[\"aws\", \"demo\"]",
		"TaskId":              "1",
		"ScheduleMode":        "scheduled",
		"ExecutionScheduleId": "1",
		"ContextType":         "instance",
		"InstanceIds":         "[91]",
		"CustomConfig":        "{\\\"test\\\":\\\"new\\\"}",
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
	templatePath := filepath.Join(dir, "job_task_schedule.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderJobTaskManualConfig generates a Terraform configuration for the job task resource
// with manual scheduling. It accepts optional overrides for field values.
// Default values are used if not overridden.
func RenderJobTaskManualConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":          "Example Job Task Manual",
		"Enabled":       "true",
		"Labels":        "[\"aws\", \"demo\"]",
		"TaskId":        "1",
		"ScheduleMode":  "manual",
		"ContextType":   "instance-label",
		"InstanceLabel": "demo",
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
	templatePath := filepath.Join(dir, "job_task_manual.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
