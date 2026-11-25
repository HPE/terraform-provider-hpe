// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagebucket

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_storage_bucket/example.tf example.tf.tmpl Name "Example Storage Bucket" ProviderType "s3" BucketName "example-bucket"
//go:generate ../../../../bin/render -out examples/resources/morpheus_storage_bucket/example_local.tf example_local.tf.tmpl ResourceName "images" Name "images" BasePath "/var/opt/morpheus/morpheus-ui/images" DefaultVirtualImageTarget "true"

func RenderStorageBucketConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	// bucket_name is optional in the schema (it only applies to Amazon, Azure,
	// CIFS, NFSv3, Openstack Swift and Rackspace CDN) but the API rejects an
	// s3 bucket without one:
	//
	//   400 {"errors":{"bucketName":"bucketName is required"}}
	//
	// so the example has to supply it.
	defaults := map[string]string{
		"Name":         "Example Storage Bucket",
		"ProviderType": "s3",
		"BucketName":   "example-bucket",
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

// RenderStorageBucketLocalConfig renders the local-storage example: a bucket
// on the appliance's disk, identified by base_path. default_virtual_image_target
// is rendered only when given, so a test can show what leaving it out does; a
// test that sets it true takes the default image store from whatever bucket the
// appliance under test already uses, and has to give it back.
func RenderStorageBucketLocalConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"ResourceName": "images",
		"Name":         "images",
		"BasePath":     "/var/opt/morpheus/morpheus-ui/images",
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
	templatePath := filepath.Join(dir, "example_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
