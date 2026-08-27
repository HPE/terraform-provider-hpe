// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_load_balancer_profile/example-id.tf example-id.tf.tmpl LoadBalancerId 1 Id 99
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_load_balancer_profile/example-name.tf example-name.tf.tmpl LoadBalancerId 1 Name "My Profile"

func RenderLoadBalancerProfileDataSourceByIDConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId": "1",
		"Id":             "99",
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
	templatePath := filepath.Join(dir, "example-id.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderLoadBalancerProfileDataSourceByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"LoadBalancerId": "1",
		"Name":           "My Profile",
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
	templatePath := filepath.Join(dir, "example-name.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
