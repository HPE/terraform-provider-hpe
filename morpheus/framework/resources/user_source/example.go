// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package user_source

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/hpe_morpheus_user_source/example.tf example.tf.tmpl Name "Corporate LDAP" Type "ldap" AccountId "1" Description "Corporate Active Directory integration" DefaultAccountRoleId "1"

func RenderUserSourceConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  "Corporate LDAP",
		"Type":                  "ldap",
		"AccountId":            "1",
		"Description":           "Corporate Active Directory integration",
		"DefaultAccountRoleId": "1",
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
