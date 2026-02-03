// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

// Disabling for now
//
//nolint:lll
//xgo:generate sh -c " ../../../../../bin/render -out examples/resources/morpheus_app_blueprint_cloud_formation/resource_yaml.tf app_blueprint_cloud_formation_resource_yaml.tf.tmpl BlueprintContent '---\nAWSTemplateFormatVersion: '2010-09-09'\nDescription: 'AWS CloudFormation Sample Template S3_Website_Bucket_With_Retain_On_Delete:\n Sample template showing how to create a publicly accessible S3 bucket configured\n for website access with a deletion policy of retain on delete. **WARNING** This\n template creates an S3 bucket that will NOT be deleted when the stack is deleted.\n You will be billed for the AWS resources used if you create a stack from this template.'\nResources:\n S3Bucket:\n Type: AWS::S3::Bucket\n Properties:\n AccessControl: PublicRead\n WebsiteConfiguration:\n IndexDocument: index.html\n ErrorDocument: error.html\n DeletionPolicy: Retain\nOutputs:\n WebsiteURL:\n Value:\n Fn::GetAtt:\n - S3Bucket\n - WebsiteURL\n Description: URL for website hosted on S3\n S3BucketSecureURL:\n Value:\n Fn::Join:\n - ''\n - - https://\n - Fn::GetAtt:\n - S3Bucket\n - DomainName\n Description: Name of S3 bucket to hold website content' CapabilityAutoExpand 'true' CapabilityIam 'true' CapabilityNamedIam 'true' Category 'cloudformation' CloudInitEnabled 'true' Description 'Example cloud formation app blueprint' InstallAgent 'true' Name 'example_cloud_formation_app_blueprint_yaml' SourceType 'yaml'"

// RenderAppBlueprintCloudFormationYamlConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.

func RenderAppBlueprintCloudFormationYamlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		//nolint:lll
		"BlueprintContent":     "---\nAWSTemplateFormatVersion: '2010-09-09'\nDescription: 'AWS CloudFormation Sample Template S3_Website_Bucket_With_Retain_On_Delete:\n Sample template showing how to create a publicly accessible S3 bucket configured\n for website access with a deletion policy of retain on delete. **WARNING** This\n template creates an S3 bucket that will NOT be deleted when the stack is deleted.\n You will be billed for the AWS resources used if you create a stack from this template.'\nResources:\n S3Bucket:\n Type: AWS::S3::Bucket\n Properties:\n AccessControl: PublicRead\n WebsiteConfiguration:\n IndexDocument: index.html\n ErrorDocument: error.html\n DeletionPolicy: Retain\nOutputs:\n WebsiteURL:\n Value:\n Fn::GetAtt:\n - S3Bucket\n - WebsiteURL\n Description: URL for website hosted on S3\n S3BucketSecureURL:\n Value:\n Fn::Join:\n - ''\n - - https://\n - Fn::GetAtt:\n - S3Bucket\n - DomainName\n Description: Name of S3 bucket to hold website content",
		"CapabilityAutoExpand": "true",
		"CapabilityIam":        "true",
		"CapabilityNamedIam":   "true",
		"Category":             "cloudformation",
		"CloudInitEnabled":     "true",
		"Description":          "Example cloud formation app blueprint",
		"InstallAgent":         "true",
		"Name":                 "example_cloud_formation_app_blueprint_yaml",
		"SourceType":           "yaml",
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
	templatePath := filepath.Join(dir, "app_blueprint_cloud_formation_resource_yaml.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
