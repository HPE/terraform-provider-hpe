// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_cloud_formation/resource_git.tf morpheus_spec_template_cloud_formation_resource_git.tf.tmpl Name tf-cloud-formation-spec-example-git SourceType repository RepositoryId 2 VersionRef main SpecPath ./spec.yaml CapabilityIam true CapabilityNamedIam true CapabilityAutoExpand true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_cloud_formation/resource_local.tf morpheus_spec_template_cloud_formation_resource_local.tf.tmpl Name tf_cloud_formation_spec_example_local SourceType local SpecContent "{\n  \"AWSTemplateFormatVersion\" : \"2010-09-09\",\n  \"Description\" : \"AWS CloudFormation Sample Template S3_Website_Bucket_With_Retain_On_Delete: Sample template showing how to create a publicly accessible S3 bucket configured for website access with a deletion policy of retain on delete. **WARNING** This template creates an S3 bucket that will NOT be deleted when the stack is deleted. You will be billed for the AWS resources used if you create a stack from this template.\",\n  \"Resources\" : {\n    \"S3Bucket\" : {\n      \"Type\" : \"AWS::S3::Bucket\",\n      \"Properties\" : {\n        \"AccessControl\" : \"PublicRead\",\n        \"WebsiteConfiguration\" : {\n          \"IndexDocument\" : \"index.html\",\n          \"ErrorDocument\" : \"error.html\"\n         }\n      },\n      \"DeletionPolicy\" : \"Retain\"\n    }\n  },\n\n  \"Outputs\" : {\n    \"WebsiteURL\" : {\n      \"Value\" : { \"Fn::GetAtt\" : [ \"S3Bucket\", \"WebsiteURL\" ] },\n      \"Description\" : \"URL for website hosted on S3\"\n    },\n    \"S3BucketSecureURL\" : {\n      \"Value\" : { \"Fn::Join\" : [ \"\", [ \"https://\", { \"Fn::GetAtt\" : [ \"S3Bucket\", \"DomainName\" ] } ] ] },\n      \"Description\" : \"Name of S3 bucket to hold website content\"\n    }\n  }\n}" CapabilityIam true CapabilityNamedIam true CapabilityAutoExpand true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_cloud_formation/resource_url.tf morpheus_spec_template_cloud_formation_resource_url.tf.tmpl Name tf_cloud_formation_spec_example_url SourceType url SpecPath http://example.com/spec.yaml CapabilityIam true CapabilityNamedIam true CapabilityAutoExpand true

func RenderSpecTemplateCloudFormationLocalConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "Example",
		"SourceType": "local",
		"SpecContent": `{
  "AWSTemplateFormatVersion" : "2010-09-09",
  "Description" : "AWS CloudFormation Sample Template S3_Website_Bucket_With_Retain_On_Delete: ` +
			`Sample template showing how to create a publicly accessible S3 bucket configured for ` +
			`website access with a deletion policy of retain on delete. **WARNING** This template ` +
			`creates an S3 bucket that will NOT be deleted when the stack is deleted. You will be ` +
			`billed for the AWS resources used if you create a stack from this template.",
  "Resources" : {
    "S3Bucket" : {
      "Type" : "AWS::S3::Bucket",
      "Properties" : {
        "AccessControl" : "PublicRead",
        "WebsiteConfiguration" : {
          "IndexDocument" : "index.html",
          "ErrorDocument" : "error.html"
         }
      },
      "DeletionPolicy" : "Retain"
    }
  },

  "Outputs" : {
    "WebsiteURL" : {
      "Value" : { "Fn::GetAtt" : [ "S3Bucket", "WebsiteURL" ] },
      "Description" : "URL for website hosted on S3"
    },
    "S3BucketSecureURL" : {
      "Value" : { "Fn::Join" : [ "", [ "https://", { "Fn::GetAtt" : [ "S3Bucket", "DomainName" ] } ] ] },
      "Description" : "Name of S3 bucket to hold website content"
    }
  }
}`,
		"CapabilityIam":        "true",
		"CapabilityNamedIam":   "true",
		"CapabilityAutoExpand": "true",
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_cloud_formation_resource_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateCloudFormationUrlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "Example",
		"SourceType":           "url",
		"SpecPath":             "http://example.com/spec.yaml",
		"CapabilityIam":        "true",
		"CapabilityNamedIam":   "true",
		"CapabilityAutoExpand": "true",
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_cloud_formation_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateCloudFormationGitConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 "Example",
		"SourceType":           "repository",
		"RepositoryId":         "2",
		"VersionRef":           "main",
		"SpecPath":             "./spec.yaml",
		"CapabilityIam":        "true",
		"CapabilityNamedIam":   "true",
		"CapabilityAutoExpand": "true",
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_cloud_formation_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
