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
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_router/example_nsxt_gateway_tier0.tf example_nsxt_gateway_tier0.tf.tmpl Name "TestRouter" GroupId "1" NetworkIntegrationId "1"
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_router/example_nsxt_gateway_tier1.tf example_nsxt_gateway_tier1.tf.tmpl Name "TestRouter" GroupId "1" NetworkIntegrationId "1"
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_router/example_juniper_sdn.tf example_juniper_sdn.tf.tmpl Name "TestRouter" GroupId "1" NetworkIntegrationId "1" ResourcePoolId "1" VirtualNetworkUids "\"0be6b79b-757c-4d92-affc-8084e328cb4c\""
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_router/example_juniper_bgp.tf example_juniper_bgp.tf.tmpl Name "TestRouter" GroupId "1" NetworkIntegrationId "1" Vendor "juniper" Asn "65000" PeerAddress1 "10.0.0.1" AddressFamilies "\"inet\"" HoldTime "90"

// SdnTestNamePrefix is the name prefix of the SDN routers the acceptance
// tests create. The SDN plugin caps a router name at 63 characters, which
// the suite-wide "TestAccMorpheus<Test>" name exceeds, so the tests use this
// short prefix and the sweeper matches it as well as the suite-wide one. The
// sweeper (build-tagged, so it cannot import this package) carries its own
// copy; TestUnitSweepPrefixMatches pins the two equal.
const SdnTestNamePrefix = "tf-sdn-rtr"

func RenderNetworkRouterGenericConfig(t *testing.T, overrides map[string]string) (string, error) {
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

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}

	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "example_generic.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderNetworkRouterNSXTGatewayTier0Config(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "TestRouter",
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

func RenderNetworkRouterNSXTGatewayTier1Config(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "TestRouter",
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

func RenderNetworkRouterJuniperSdnConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "TestRouter",
		"GroupId":              "1",
		"NetworkIntegrationId": "1",
		"ResourcePoolId":       "1",
		"VirtualNetworkUids":   `"0be6b79b-757c-4d92-affc-8084e328cb4c"`,
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
	templatePath := filepath.Join(dir, "example_juniper_sdn.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderNetworkRouterJuniperBgpConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "TestRouter",
		"GroupId":              "1",
		"NetworkIntegrationId": "1",
		"Vendor":               "juniper",
		"Asn":                  "65000",
		"PeerAddress1":         "10.0.0.1",
		"AddressFamilies":      `"inet"`,
		"HoldTime":             "90",
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
	templatePath := filepath.Join(dir, "example_juniper_bgp.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
