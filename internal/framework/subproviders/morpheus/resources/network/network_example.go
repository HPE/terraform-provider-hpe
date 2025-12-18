// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../../cmd/render -out examples/resources/morpheus_network/example_aws.tf example_aws.tf.tmpl Name "example-terraform-aws" Description "AWS subnet" CloudName "AWS Cloud" PoolId "1" GroupName "Example Group" TypeId "36" AssignPublicIP "true" AvailabilityZone "us-west-1a" Active "true" DHCPServer "true" ApplianceURLProxyBypass "true" TenantName "Master Tenant" Visibility "private" CIDR "10.200.99.0/24" ZonePoolId "12329" Labels "\"terraform\", \"example\""

//go:generate go run ../../../../../../cmd/render -out examples/resources/morpheus_network/example_azure.tf example_azure.tf.tmpl Name "example-terraform-azure" Description "Azure network" CloudName "Azure Cloud" PoolId "1" GroupName "Example Group" TypeId "35" CIDR "10.100.0.0/16" Visibility "public" Active "true" DHCPServer "true" ApplianceURLProxyBypass "false" ResourceGroupId "all-attrs-resource-group" SubnetName "all-attrs-subnet" SubnetCIDR "10.100.1.0/24" Location "eastus" TenantName "Master Tenant" Labels "\"terraform\", \"example\""

//go:generate go run ../../../../../../cmd/render -out examples/resources/morpheus_network/example_gcp.tf example_gcp.tf.tmpl Name "example-terraform-gcp" Description "GCP network" CloudName "Google Cloud" PoolId "1" GroupName "Examle Group" TypeId "38" MTU "1460" AutoCreate "true" Active "true" DHCPServer "false" ApplianceURLProxyBypass "true" TenantName "Master Tenant" Visibility "private" CIDR "10.0.0.0/8" ZonePoolId "85990" Labels "\"terraform\", \"example\""

//go:generate go run ../../../../../../cmd/render -out examples/resources/morpheus_network/example_host.tf example_host.tf.tmpl Name "example-terraform-host" Description "A host network" CloudName "Standard Cloud" PoolId "1" GroupName "Example Group" TypeId "1" Active "true" DHCPServer "false" ApplianceURLProxyBypass "true" TenantName "Master Tenant" Visibility "private" CIDR "10.0.0.0/8" Labels "terraform, example"

//go:generate go run ../../../../../../cmd/render -out examples/resources/morpheus_network/example_ovs_port_group.tf example_ovs_port_group.tf.tmpl Name "Terraform OVS Port Group" Description "OVS Port Group network" CloudName "Morpheus Standard Cloud" PoolId "3251" GroupName "ExampleGroup" TypeId "63" SwitchId "Compute" Active "true" DHCPServer "false" ApplianceURLProxyBypass "true" TenantName "Master Tenant" Visibility "public" CIDR "10.32.148.0/22" ZonePoolId "62299" VLANId "43" Labels "\"terraform\", \"example\""

func RenderNetworkHostConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"CloudName":               "Standard Cloud",
		"GroupName":               "Example Group",
		"TenantName":              "Master Tenant",
		"Name":                    "example-terraform-host",
		"Description":             "A host network",
		"PoolId":                  "1",
		"TypeId":                  "1",
		"Active":                  "true",
		"DHCPServer":              "false",
		"ApplianceURLProxyBypass": "true",
		"Visibility":              "private",
		"CIDR":                    "10.0.0.0/8",
		"Labels":                  "terraform, example",
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
	templatePath := filepath.Join(dir, "example_host.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderNetworkAWSConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"CloudName":               "AWS Cloud",
		"GroupName":               "Example Group",
		"TenantName":              "Master Tenant",
		"Name":                    "example-terraform-aws",
		"Description":             "AWS subnet",
		"PoolId":                  "1",
		"TypeId":                  "36",
		"AssignPublicIP":          "true",
		"AvailabilityZone":        "us-west-1a",
		"Active":                  "true",
		"DHCPServer":              "true",
		"ApplianceURLProxyBypass": "true",
		"Visibility":              "private",
		"CIDR":                    "10.200.99.0/24",
		"ZonePoolId":              "12329",
		"Labels":                  "\"terraform\", \"example\"",
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
	templatePath := filepath.Join(dir, "example_aws.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderNetworkAzureConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"CloudName":               "Azure Cloud",
		"GroupName":               "Example Group",
		"TenantName":              "Master Tenant",
		"Name":                    "example-terraform-azure",
		"Description":             "Azure network",
		"PoolId":                  "1",
		"TypeId":                  "35",
		"CIDR":                    "10.100.0.0/16",
		"Visibility":              "public",
		"Active":                  "true",
		"DHCPServer":              "true",
		"ApplianceURLProxyBypass": "false",
		"ResourceGroupId":         "all-attrs-resource-group",
		"SubnetName":              "all-attrs-subnet",
		"SubnetCIDR":              "10.100.1.0/24",
		"Location":                "eastus",
		"Labels":                  "\"terraform\", \"example\"",
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
	templatePath := filepath.Join(dir, "example_azure.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderNetworkGCPConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"CloudName":               "Google Cloud",
		"GroupName":               "Example Group",
		"TenantName":              "Master Tenant",
		"Name":                    "example-terraform-gcp",
		"Description":             "GCP network",
		"PoolId":                  "1",
		"TypeId":                  "38",
		"MTU":                     "1460",
		"AutoCreate":              "true",
		"Active":                  "true",
		"DHCPServer":              "false",
		"ApplianceURLProxyBypass": "true",
		"Visibility":              "private",
		"CIDR":                    "10.0.0.0/8",
		"ZonePoolId":              "85990",
		"Labels":                  "\"terraform\", \"example\"",
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
	templatePath := filepath.Join(dir, "example_gcp.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderNetworkOVSPortGroupConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"CloudName":               "Morpheus Standard Cloud",
		"GroupName":               "ExampleGroup",
		"TenantName":              "Master Tenant",
		"Name":                    "Terraform OVS Port Group",
		"Description":             "OVS Port Group network",
		"PoolId":                  "3251",
		"TypeId":                  "63",
		"SwitchId":                "Compute",
		"Active":                  "true",
		"DHCPServer":              "false",
		"ApplianceURLProxyBypass": "true",
		"Visibility":              "public",
		"CIDR":                    "10.32.148.0/22",
		"ZonePoolId":              "62299",
		"VLANId":                  "43",
		"Labels":                  "\"terraform\", \"example\"",
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
	templatePath := filepath.Join(dir, "example_ovs_port_group.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
