// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_network_router/example_generic.tf example_generic.tf.tmpl Name "TestRouter" TypeId "1" GroupId "1" NetworkIntegrationId "1"
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_router/example_nsxt_gateway_tier0.tf example_nsxt_gateway_tier0.tf.tmpl Name "TestRouter" TypeId "1" GroupId "1" NetworkIntegrationId "1" "IpServerId" "1"
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_router/example_nsxt_gateway_tier1.tf example_nsxt_gateway_tier1.tf.tmpl Name "TestRouter" TypeId "1" GroupId "1" NetworkIntegrationId "1"

func RenderNetworkRouterGenericConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "TestRouter",
		"TypeId":               "1",
		"GroupId":              "1",
		"NetworkIntegrationId": "1",
		"IpServerId":           "1",
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
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderNetworkRouterNSXGatewayTier0Config(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "TestRouter",
		"TypeId":               "1",
		"GroupId":              "1",
		"NetworkIntegrationId": "1",
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
	templatePath := filepath.Join(dir, "example_nsxt_gateway_tier0.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderNetworkRouterNSXGatewayTier1Config(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "TestRouter",
		"TypeId":               "1",
		"GroupId":              "1",
		"NetworkIntegrationId": "1",
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
	templatePath := filepath.Join(dir, "example_nsxt_gateway_tier1.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
