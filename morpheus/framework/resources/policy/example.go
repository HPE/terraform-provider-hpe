// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package policy

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render example.tf.tmpl ResourceName "group_policy" Name "TestMaxMemoryGroupPolicy" Description "Example group-scoped policy" AssociatedResourceType "Group" AssociatedResourceID "1" PolicyTypeCode "maxMemory" ConfigKey "maxMemory" ConfigValue "8"
//go:generate ../../../../bin/render example.tf.tmpl ResourceName "cloud_policy" Name "TestMaxMemoryCloudPolicy" Description "Example cloud-scoped policy" AssociatedResourceType "Cloud" AssociatedResourceID "1" PolicyTypeCode "maxMemory" ConfigKey "maxMemory" ConfigValue "8"
//go:generate ../../../../bin/render example.tf.tmpl ResourceName "user_policy" Name "TestMaxMemoryUserPolicy" Description "Example user-scoped policy" AssociatedResourceType "User" AssociatedResourceID "1" PolicyTypeCode "maxMemory" ConfigKey "maxMemory" ConfigValue "8"
//go:generate ../../../../bin/render example.tf.tmpl ResourceName "role_policy" Name "TestMaxMemoryRolePolicy" Description "Example role-scoped policy" AssociatedResourceType "Role" AssociatedResourceID "1" PolicyTypeCode "maxMemory" ConfigKey "maxMemory" ConfigValue "8"

func RenderPolicyConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name": "Example Policy",
		"Code": "maxContainers",
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
