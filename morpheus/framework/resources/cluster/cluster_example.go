// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_cluster/example_hvm.tf example_hvm.tf.tmpl Name "TestCluster" Description "A test HVM cluster" CloudId "1" GroupId "1" LayoutId "2" CreateUser "false" DynamicPlacement "false" CpuArch "x86_64" CpuModel "host-model" PowerPolicy "default" ServicePlanId "1" SshPort "22" SshUsername "user" SshKeyPairId "1" SshHost1Name "host1" SshHost1Ip "10.0.0.1" SshHost2Name "host2" SshHost2Ip "10.0.0.2" SshHost3Name "host3" SshHost3Ip "10.0.0.3" ManagementNetInterface "eth0" Visibility "private" Tag1Name "source" Tag1Value "terraform" Tag2Name "environment" Tag2Value "example" Label1 "terraform" Label2 "example"

//go:generate ../../../../bin/render -out examples/resources/morpheus_cluster/example_generic_hvm.tf example_generic_hvm.tf.tmpl Name "TestCluster" Description "A HVM cluster created with a dynamic config" CloudId "1" GroupId "1" LayoutId "2" ClusterTypeCode "mvm-cluster" ServicePlanId "1" SshPort "22" SshUsername "user" SshKeyPairId "1" SshHost1Name "host1" SshHost1Ip "10.0.0.1" SshHost2Name "host2" SshHost2Ip "10.0.0.2" SshHost3Name "host3" SshHost3Ip "10.0.0.3" ManagementNetInterface "eth0" Visibility "private" Tag1Name "source" Tag1Value "terraform" Tag2Name "environment" Tag2Value "example" Label1 "terraform" Label2 "example"

func RenderClusterHvmConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                   "TestCluster",
		"Description":            "A test HVM cluster",
		"CloudId":                "1",
		"GroupId":                "1",
		"LayoutId":               "2",
		"CreateUser":             "false",
		"DynamicPlacement":       "false",
		"CpuArch":                "x86_64",
		"CpuModel":               "host-model",
		"PowerPolicy":            "default",
		"ServicePlanId":          "1",
		"SshPort":                "22",
		"SshUsername":            "user",
		"SshKeyPairId":           "1",
		"SshHost1Name":           "host1",
		"SshHost1Ip":             "10.0.0.1",
		"SshHost2Name":           "host2",
		"SshHost2Ip":             "10.0.0.2",
		"SshHost3Name":           "host3",
		"SshHost3Ip":             "10.0.0.3",
		"ManagementNetInterface": "eth0",
		"Visibility":             "private",
		"Tag1Name":               "source",
		"Tag1Value":              "terraform",
		"Tag2Name":               "environment",
		"Tag2Value":              "example",
		"Label1":                 "terraform",
		"Label2":                 "example",
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
	templatePath := filepath.Join(dir, "example_hvm.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderClusterGenericConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                   "TestCluster",
		"Description":            "A HVM cluster created with a dynamic config",
		"CloudId":                "1",
		"GroupId":                "1",
		"LayoutId":               "2",
		"ClusterTypeCode":        "mvm-cluster",
		"ServicePlanId":          "1",
		"SshPort":                "22",
		"SshUsername":            "user",
		"SshKeyPairId":           "1",
		"ManagementNetInterface": "eth0",
		"SshHost1Name":           "host1",
		"SshHost1Ip":             "10.0.0.1",
		"SshHost2Name":           "host2",
		"SshHost2Ip":             "10.0.0.2",
		"SshHost3Name":           "host3",
		"SshHost3Ip":             "10.0.0.3",
		"Visibility":             "private",
		"Tag1Name":               "source",
		"Tag1Value":              "terraform",
		"Tag2Name":               "environment",
		"Tag2Value":              "example",
		"Label1":                 "terraform",
		"Label2":                 "example",
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
	templatePath := filepath.Join(dir, "example_generic_hvm.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
