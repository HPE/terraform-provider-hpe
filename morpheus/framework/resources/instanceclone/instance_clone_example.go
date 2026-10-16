// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_instance_clone/example.tf example.tf.tmpl Name "my-clone" SourceInstanceName "my-instance" SourceInstanceId "data.hpe_morpheus_instance.source.id" ClusterName "MyCluster" NetworkName "MyNetwork" NetworkId "data.hpe_morpheus_network.clone.id"

func RenderInstanceCloneConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	// The defaults render the published example: a source instance and a
	// network looked up by name. The tests override SourceInstanceId and
	// NetworkId with references to resources and data sources of their own
	// and leave the name fields empty, so no lookup block is rendered.
	defaults := map[string]string{
		"Name":               "my-clone",
		"SourceInstanceName": "my-instance",
		"SourceInstanceId":   "data.hpe_morpheus_instance.source.id",
		"ClusterName":        "MyCluster",
		"NetworkName":        "MyNetwork",
		"NetworkId":          "data.hpe_morpheus_network.clone.id",
		"DatastoreId":        "",
		"UserGroup":          "",
		"StorageProfile":     "",
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
