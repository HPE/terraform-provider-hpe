// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package job

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_job_workflow/resource_date_and_time.tf job_workflow_date_and_time.tf.tmpl Name "TF Example Workflow Job Date and Time" Enabled true Labels "[\"aws\", \"demo\"]" WorkflowId 1 ScheduleMode date_and_time ScheduledDateAndTime 2022-12-30T06:00:00Z ContextType instance InstanceIds "[1, 2]" CustomOptions "{\"demo\" = \"testing\"}"

//go:generate ../../../../../bin/render -out examples/resources/morpheus_job_workflow/resource_schedule.tf job_workflow_schedule.tf.tmpl Name "TF Example Workflow Job Schedule" Enabled true Labels "[\"aws\", \"demo\"]" WorkflowId 1 ScheduleMode scheduled ExecutionScheduleId 1 ContextType instance InstanceIds "[91]" CustomOptions "{\"demo\" = \"testing\"}"

//go:generate ../../../../../bin/render -out examples/resources/morpheus_job_workflow/resource_manual.tf job_workflow_manual.tf.tmpl Name "TF Example Workflow Job Manual" Enabled true Labels "[\"aws\", \"demo\"]" WorkflowId 1 ScheduleMode manual ContextType instance-label InstanceLabel demo CustomOptions "{\"demo\" = \"testing\"}"

// RenderJobWorkflowDateAndTimeConfig generates a Terraform configuration for the job workflow resource
// with date and time scheduling. It accepts optional overrides for field values.
// Default values are used if not overridden.
func RenderJobWorkflowDateAndTimeConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "TF Example Workflow Job Date and Time",
		"Enabled":              "true",
		"Labels":               "[\"aws\", \"demo\"]",
		"WorkflowId":           "1",
		"ScheduleMode":         "date_and_time",
		"ScheduledDateAndTime": "2022-12-30T06:00:00Z",
		"ContextType":          "instance",
		"InstanceIds":          "[1, 2]",
		"CustomOptions":        "{\"demo\" = \"testing\"}",
	}

	// Apply overrides to defaults
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
	templatePath := filepath.Join(dir, "job_workflow_date_and_time.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderJobWorkflowScheduleConfig generates a Terraform configuration for the job workflow resource
// with execution schedule. It accepts optional overrides for field values.
// Default values are used if not overridden.
func RenderJobWorkflowScheduleConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                "TF Example Workflow Job Schedule",
		"Enabled":             "true",
		"Labels":              "[\"aws\", \"demo\"]",
		"WorkflowId":          "1",
		"ScheduleMode":        "scheduled",
		"ExecutionScheduleId": "1",
		"ContextType":         "instance",
		"InstanceIds":         "[91]",
		"CustomOptions":       "{\"demo\" = \"testing\"}",
	}

	// Apply overrides to defaults
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
	templatePath := filepath.Join(dir, "job_workflow_schedule.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderJobWorkflowManualConfig generates a Terraform configuration for the job workflow resource
// with manual execution. It accepts optional overrides for field values.
// Default values are used if not overridden.
func RenderJobWorkflowManualConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":          "TF Example Workflow Job Manual",
		"Enabled":       "true",
		"Labels":        "[\"aws\", \"demo\"]",
		"WorkflowId":    "1",
		"ScheduleMode":  "manual",
		"ContextType":   "instance-label",
		"InstanceLabel": "demo",
		"CustomOptions": "{\"demo\" = \"testing\"}",
	}

	// Apply overrides to defaults
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
	templatePath := filepath.Join(dir, "job_workflow_manual.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
