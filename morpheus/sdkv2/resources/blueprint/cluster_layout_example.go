// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate../../../../bin/render -out examples/resources/morpheus_cluster_layout/resource.tf cluster_layout_resource.tf.tmpl Name "tfexample cluster layout" Description "Terraform example cluster layout" Version "1.0" Creatable "false" MinimumMemory "4294967296" ClusterTypeId "1" ProvisionTypeId "3" EnableScaling "false" EvarName "application" EvarValue "first" EvarExport "true" MasterNodePoolCount "1" MasterNodePoolNodeTypeId "3" MasterNodePoolPriorityOrder "0" WorkerNodePool1Count "4" WorkerNodePool1NodeTypeId "4" WorkerNodePool1PriorityOrder "1" WorkerNodePool2Count "4" WorkerNodePool2NodeTypeId "3" WorkerNodePool2PriorityOrder "2"

// RenderClusterLayoutConfig generates a Terraform configuration for the cluster layout resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderClusterLayoutConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                         "tfexample cluster layout",
		"Description":                  "Terraform example cluster layout",
		"Version":                      "1.0",
		"Creatable":                    "false",
		"MinimumMemory":                "4294967296",
		"ClusterTypeId":                "1",
		"ProvisionTypeId":              "3",
		"EnableScaling":                "false",
		"EvarName":                     "application",
		"EvarValue":                    "first",
		"EvarExport":                   "true",
		"MasterNodePoolCount":          "1",
		"MasterNodePoolNodeTypeId":     "3",
		"MasterNodePoolPriorityOrder":  "0",
		"WorkerNodePool1Count":         "4",
		"WorkerNodePool1NodeTypeId":    "4",
		"WorkerNodePool1PriorityOrder": "1",
		"WorkerNodePool2Count":         "4",
		"WorkerNodePool2NodeTypeId":    "3",
		"WorkerNodePool2PriorityOrder": "2",
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
	templatePath := filepath.Join(dir, "cluster_layout_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
