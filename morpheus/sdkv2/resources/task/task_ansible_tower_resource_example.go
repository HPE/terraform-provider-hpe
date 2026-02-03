// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c " ../../../../../bin/render -out examples/resources/morpheus_task_ansible_tower/resource.tf task_ansible_tower_resource.tf.tmpl AllowCustomConfig 'true' AnsibleTowerIntegrationId '1' AnsibleTowerInventoryId '5' Code 'tfexample-ansible-tower-task' ExecuteMode 'executeAll' ExecuteTarget 'local' Group 'demo' JobTemplateId '3' Labels '[\"demo\", \"terraform\"]' Name 'tfexample_task_ansible_tower' RetryCount '5' RetryDelaySeconds '10' Retryable 'true' ScmOverride 'main' Visibility 'public'"

// RenderTaskAnsibleTowerConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderTaskAnsibleTowerConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"AllowCustomConfig":         "true",
		"AnsibleTowerIntegrationId": "1",
		"AnsibleTowerInventoryId":   "5",
		"Code":                      "tfexample-ansible-tower-task",
		"ExecuteMode":               "executeAll",
		"ExecuteTarget":             "local",
		"Group":                     "demo",
		"JobTemplateId":             "3",
		"Labels":                    "[\"demo\", \"terraform\"]",
		"Name":                      "tfexample_task_ansible_tower",
		"RetryCount":                "5",
		"RetryDelaySeconds":         "10",
		"Retryable":                 "true",
		"ScmOverride":               "main",
		"Visibility":                "public",
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
	templatePath := filepath.Join(dir, "task_ansible_tower_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
