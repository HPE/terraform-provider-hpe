// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/app_blueprint_cloud_formation/app_blueprint_cloud_formation_resource_yaml.tf app_blueprint_cloud_formation_resource_yaml.tf.tmpl CapabilityAutoExpand "true" CapabilityIam "true" CapabilityNamedIam "true" Category "cloudformation" CloudInitEnabled "true" Description "Example cloud formation app blueprint" InstallAgent "true" Name "example_cloud_formation_app_blueprint_yaml" SourceType "yaml" BlueprintContent "..."

// RenderAppBlueprintCloudFormationYamlConfig generates a Terraform configuration
// for the app blueprint cloud formation yaml resource.
// It accepts optional overrides for field values. Default values are used if not overridden.

func RenderAppBlueprintCloudFormationYamlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"BlueprintContent": `---
AWSTemplateFormatVersion: '2010-09-09'
Description: 'AWS CloudFormation Sample Template S3_Website_Bucket_With_Retain_On_Delete:
  Sample template showing how to create a publicly accessible S3 bucket configured
  for website access with a deletion policy of retain on delete. **WARNING** This
  template creates an S3 bucket that will NOT be deleted when the stack is deleted.
  You will be billed for the AWS resources used if you create a stack from this template.'
Resources:
  S3Bucket:
    Type: AWS::S3::Bucket
    Properties:
      AccessControl: PublicRead
      WebsiteConfiguration:
        IndexDocument: index.html
        ErrorDocument: error.html
    DeletionPolicy: Retain
Outputs:
  WebsiteURL:
    Value:
      Fn::GetAtt:
      - S3Bucket
      - WebsiteURL
    Description: URL for website hosted on S3
  S3BucketSecureURL:
    Value:
      Fn::Join:
      - ''
      - - https://
        - Fn::GetAtt:
          - S3Bucket
          - DomainName
    Description: Name of S3 bucket to hold website content`,
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
