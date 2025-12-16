// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_kubernetes/resource_git.tf morpheus_spec_template_kubernetes_resource_git.tf.tmpl Name tf-kubernetes-spec-example-git SourceType repository RepositoryId 2 VersionRef main SpecPath ./spec.yaml

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_kubernetes/resource_local.tf morpheus_spec_template_kubernetes_resource_local.tf.tmpl Name tf-terraform-spec-example-local SourceType local SpecContent "---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx-deployment\n  labels:\n    app: nginx\nspec:\n  replicas: 3\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      containers:\n      - name: nginx\n        image: nginx:1.14.2\n        ports:\n        - containerPort: 80\n"

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_kubernetes/resource_url.tf morpheus_spec_template_kubernetes_resource_url.tf.tmpl Name tf-kubernetes-spec-example-url SourceType url SpecPath http://example.com/spec.yaml

func RenderHpeMorpheusFileTemplateConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":            name,
		"Labels":          `["demo", "template", "terraform"]`,
		"FileName":        "tfcustom.cnf",
		"FilePath":        "/etc/my.cnf.d",
		"Phase":           "preProvision",
		"FileContent":     `"# Test MySQL Configuration\n[mysqld]\ninnodb_buffer_pool_size = 128M"`,
		"FileOwner":       "root",
		"SettingName":     "myCnf",
		"SettingCategory": "master",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_file_template_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Labels", defaults["Labels"],
		"FileName", defaults["FileName"],
		"FilePath", defaults["FilePath"],
		"Phase", defaults["Phase"],
		"FileContent", defaults["FileContent"],
		"FileOwner", defaults["FileOwner"],
		"SettingName", defaults["SettingName"],
		"SettingCategory", defaults["SettingCategory"],
	)
}

// RenderClusterPackageConfig generates a test configuration for cluster package resource.
// It accepts a name and a map of field overrides to customize the default values.
func RenderClusterPackageConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":            name,
		"Code":            "tf-example-cluster-package",
		"Description":     "Terraform example cluster package",
		"PackageVersion":  "1.2.3",
		"Type":            "apps",
		"PackageType":     "example",
		"Enabled":         "true",
		"RepeatInstall":   "true",
		"SpecTemplateIds": "[1,2]",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"cluster_package_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Description", defaults["Description"],
		"PackageVersion", defaults["PackageVersion"],
		"Type", defaults["Type"],
		"PackageType", defaults["PackageType"],
		"Enabled", defaults["Enabled"],
		"RepeatInstall", defaults["RepeatInstall"],
		"SpecTemplateIds", defaults["SpecTemplateIds"],
	)
	if err != nil {
		return "", err
	}

	return resourceConfig, nil
}

// RenderSpecTemplateKubernetesGitConfig renders the configuration for the Git-based
// Kubernetes spec template resource. Pass overrides as a map to customize field values.
func RenderSpecTemplateKubernetesGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         name,
		"SourceType":   "repository",
		"RepositoryId": "2",
		"VersionRef":   "main",
		"SpecPath":     "./spec.yaml",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"morpheus_spec_template_kubernetes_resource_git.tf.tmpl",
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
		"RepositoryId", defaults["RepositoryId"],
		"VersionRef", defaults["VersionRef"],
		"SpecPath", defaults["SpecPath"],
	)
	if err != nil {
		return "", err
	}

	return resourceConfig, nil
}

func RenderMorpheusSpecTemplateArmUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       name,
		"SourceType": "url",
		"SpecPath":   "http://example.com/spec.json",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_arm_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderSpecTemplateTerraformGitConfig renders the Terraform config for
// spec_template_terraform_resource_git tests
func RenderSpecTemplateTerraformGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         name,
		"SourceType":   "repository",
		"RepositoryId": "2",
		"VersionRef":   "main",
		"SpecPath":     "Instance Types/Terraform/CloudResource/aws/vpc.tf",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_terraform_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateCloudFormationUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 name,
		"SourceType":           "url",
		"SpecPath":             "http://example.com/spec.yaml",
		"CapabilityIam":        "true",
		"CapabilityNamedIam":   "true",
		"CapabilityAutoExpand": "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
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

// RenderSpecTemplateKubernetesUrlConfig renders the configuration for the URL-based
// Kubernetes spec template resource. Pass overrides as a map to customize field values.
func RenderSpecTemplateKubernetesUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       name,
		"SourceType": "url",
		"SpecPath":   "http://example.com/spec.yaml",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"morpheus_spec_template_kubernetes_resource_url.tf.tmpl",
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
		"SpecPath", defaults["SpecPath"],
	)
	if err != nil {
		return "", err
	}

	return resourceConfig, nil
}

// RenderSpecTemplateTerraformUrlConfig renders the Terraform config for
// spec_template_terraform_resource_url tests
func RenderSpecTemplateTerraformUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       name,
		"SourceType": "url",
		"SpecPath":   "http://example.com/spec.tf",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_terraform_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateCloudFormationGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 name,
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

	args := []string{}
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

func RenderMorpheusSpecTemplateArmGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         name,
		"SourceType":   "repository",
		"RepositoryId": "2",
		"VersionRef":   "main",
		"SpecPath":     "./test.json",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_arm_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateHelmGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         name,
		"RepositoryId": "2",
		"SourceType":   "repository",
		"SpecPath":     "./spec.yaml",
		"VersionRef":   "main",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_helm_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"RepositoryId", defaults["RepositoryId"],
		"SourceType", defaults["SourceType"],
		"SpecPath", defaults["SpecPath"],
		"VersionRef", defaults["VersionRef"],
	)
}

func RenderSpecTemplateCloudFormationLocalConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       name,
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

	args := []string{}
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

// RenderScriptTemplateConfig renders the template with provided overrides
func RenderScriptTemplateConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":          name,
		"Labels":        "[\"demo\", \"template\", \"terraform\"]",
		"ScriptType":    "bash",
		"ScriptPhase":   "provision",
		"ScriptContent": "echo \"testing\"",
		"RunAsUser":     "root",
		"Sudo":          "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_script_template_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Labels", defaults["Labels"],
		"ScriptType", defaults["ScriptType"],
		"ScriptPhase", defaults["ScriptPhase"],
		"ScriptContent", defaults["ScriptContent"],
		"RunAsUser", defaults["RunAsUser"],
		"Sudo", defaults["Sudo"],
	)
}

// RenderSpecTemplateTerraformLocalConfig renders the Terraform config for
// spec_template_terraform_resource_local tests
func RenderSpecTemplateTerraformLocalConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       name,
		"SourceType": "local",
		"SpecContent": `resource "aws_instance" "instance_1" {
  ami           = "ami-0b91a410940e82c54"
  instance_type = "t2.micro"
}`,
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_terraform_resource_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateHelmLocalConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaultSpecContent := "apiVersion: v1\nkind: Service\nmetadata:\nname: {{ template \"fullname\" . }}\n" +
		"labels:\n    chart: \"{{ .Chart.Name }}-{{ .Chart.Version | replace \"+\" \"_\" }}\"\nspec:\n" +
		"type: {{ .Values.service.type }}\nports:\n- port: {{ .Values.service.externalPort }}\n" +
		"    targetPort: {{ .Values.service.internalPort }}\n    protocol: TCP\n" +
		"    name: {{ .Values.service.name }}\nselector:\n    app: {{ template \"fullname\" . }}"

	defaults := map[string]string{
		"Name":        name,
		"SourceType":  "local",
		"SpecContent": defaultSpecContent,
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_helm_resource_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
		"SpecContent", defaults["SpecContent"],
	)
}

func RenderSecurityPackageConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        name,
		"Description": "Terraform security package example",
		"Labels":      "[\"demo\", \"terraform\"]",
		"Enabled":     "true",
		"Url": "https://github.com/ComplianceAsCode/content/releases/download/" +
			"v0.1.59/scap-security-guide-0.1.59.zip",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_security_package_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderMorpheusSpecTemplateArmLocalConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       name,
		"SourceType": "local",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_arm_resource_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateHelmUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       name,
		"SourceType": "url",
		"SpecPath":   "http://example.com/chart.yaml",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_helm_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
		"SpecPath", defaults["SpecPath"],
	)
}

// RenderSpecTemplateKubernetesLocalConfig renders the configuration for the local
// Kubernetes spec template resource. Pass overrides as a map to customize field values.
func RenderSpecTemplateKubernetesLocalConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaultSpecContent := `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80`

	defaults := map[string]string{
		"Name":        name,
		"SourceType":  "local",
		"SpecContent": defaultSpecContent,
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"morpheus_spec_template_kubernetes_resource_local.tf.tmpl",
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
		"SpecContent", defaults["SpecContent"],
	)
	if err != nil {
		return "", err
	}

	return resourceConfig, nil
}

