// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate go run ../../../../cmd/render -out examples/resources/morpheus_task/example_conditional_workflow.tf example_conditional_workflow.tf.tmpl Name "Example Conditional Workflow Task" IfOperationalWorkflowId "4090" IfOperationalWorkflowName "Example If Workflow" ElseOperationalWorkflowId "4091" ElseOperationalWorkflowName "Example Else Workflow"
//go:generate go run ../../../../cmd/render -out examples/resources/morpheus_task/example_conditional_workflow_null_else.tf example_conditional_workflow_null_else.tf.tmpl Name "Example Conditional Workflow Task" IfOperationalWorkflowId "4090" IfOperationalWorkflowName "Example If Workflow"
//go:generate go run ../../../../cmd/render -out examples/resources/morpheus_task/example_generic_config.tf example_generic_config.tf.tmpl Name "Example Generic Task" OperationalWorkflowId "4090" OperationalWorkflowName "Example Workflow"

func RenderTaskGenericConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                    "Example Generic Task",
		"OperationalWorkflowId":   "4090",
		"OperationalWorkflowName": "Example Workflow",
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
	templatePath := filepath.Join(dir, "example_generic_config.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderTaskConditionalWorkflowConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                        "Example Conditional Workflow Task",
		"IfOperationalWorkflowId":     "4090",
		"IfOperationalWorkflowName":   "Example If Workflow",
		"ElseOperationalWorkflowId":   "4091",
		"ElseOperationalWorkflowName": "Example Else Workflow",
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
	templatePath := filepath.Join(dir, "example_conditional_workflow.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
