// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network_pool_server

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_network_pool_server/example_infoblox.tf example_infoblox.tf.tmpl Name "Infoblox IPAM" TypeId "1" Enabled "true" ServiceUrl "https://infoblox.example.com/wapi/v2.12" ServiceUsername "admin" ServicePassword "changeme" ServicePasswordVersion "1" IgnoreSsl "true" NetworkFilter "10.0.0.0/8" ZoneFilter "example.com" TenantMatch ".*" ServiceMode "static" ServiceThrottleRate "0"
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_pool_server/example_bluecat.tf example_bluecat.tf.tmpl Name "Bluecat IPAM" TypeId "2" Enabled "true" ServiceUrl "https://bluecat.example.com/api" ServiceUsername "admin" ServicePassword "changeme" ServicePasswordVersion "1" IgnoreSsl "false" NetworkFilter "192.168.0.0/16" ServiceThrottleRate "50"
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_pool_server/example_phpipam.tf example_phpipam.tf.tmpl Name "phpIPAM" TypeId "3" Enabled "true" ServiceUrl "https://phpipam.example.com/api/app" ServiceUsername "admin" ServicePassword "changeme" ServicePasswordVersion "1" IgnoreSsl "false" NetworkFilter "172.16.0.0/12" ServiceThrottleRate "0"
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_pool_server/example_solarwinds.tf example_solarwinds.tf.tmpl Name "SolarWinds IPAM" TypeId "4" Enabled "true" ServiceUrl "https://solarwinds.example.com:17778/SolarWinds/InformationService/v3/Json" ServiceUsername "admin" ServicePassword "changeme" ServicePasswordVersion "1" IgnoreSsl "true" ServiceThrottleRate "100"
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_pool_server/example_credential.tf example_credential.tf.tmpl Name "Infoblox with Credential" TypeId "1" Enabled "true" ServiceUrl "https://infoblox.example.com/wapi/v2.12" CredentialId "42" IgnoreSsl "true"

func RenderNetworkPoolServerInfobloxConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                   "Infoblox IPAM",
		"TypeId":                 "1",
		"Enabled":                "true",
		"ServiceUrl":             "https://infoblox.example.com/wapi/v2.12",
		"ServiceUsername":        "admin",
		"ServicePassword":        "changeme",
		"ServicePasswordVersion": "1",
		"IgnoreSsl":              "true",
		"NetworkFilter":          "10.0.0.0/8",
		"ZoneFilter":             "example.com",
		"TenantMatch":            ".*",
		"ServiceMode":            "static",
		"ServiceThrottleRate":    "0",
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
	templatePath := filepath.Join(dir, "example_infoblox.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderNetworkPoolServerBluecatConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                   "Bluecat IPAM",
		"TypeId":                 "2",
		"Enabled":                "true",
		"ServiceUrl":             "https://bluecat.example.com/api",
		"ServiceUsername":        "admin",
		"ServicePassword":        "changeme",
		"ServicePasswordVersion": "1",
		"IgnoreSsl":              "false",
		"NetworkFilter":          "192.168.0.0/16",
		"ServiceThrottleRate":    "50",
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
	templatePath := filepath.Join(dir, "example_bluecat.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderNetworkPoolServerPhpIPAMConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                   "phpIPAM",
		"TypeId":                 "3",
		"Enabled":                "true",
		"ServiceUrl":             "https://phpipam.example.com/api/app",
		"ServiceUsername":        "admin",
		"ServicePassword":        "changeme",
		"ServicePasswordVersion": "1",
		"IgnoreSsl":              "false",
		"NetworkFilter":          "172.16.0.0/12",
		"ServiceThrottleRate":    "0",
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
	templatePath := filepath.Join(dir, "example_phpipam.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderNetworkPoolServerSolarWindsConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                   "SolarWinds IPAM",
		"TypeId":                 "4",
		"Enabled":                "true",
		"ServiceUrl":             "https://solarwinds.example.com:17778/SolarWinds/InformationService/v3/Json",
		"ServiceUsername":        "admin",
		"ServicePassword":        "changeme",
		"ServicePasswordVersion": "1",
		"IgnoreSsl":              "true",
		"ServiceThrottleRate":    "100",
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
	templatePath := filepath.Join(dir, "example_solarwinds.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}

func RenderNetworkPoolServerCredentialConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "Infoblox with Credential",
		"TypeId":       "1",
		"Enabled":      "true",
		"ServiceUrl":   "https://infoblox.example.com/wapi/v2.12",
		"CredentialId": "42",
		"IgnoreSsl":    "true",
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
	templatePath := filepath.Join(dir, "example_credential.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}
