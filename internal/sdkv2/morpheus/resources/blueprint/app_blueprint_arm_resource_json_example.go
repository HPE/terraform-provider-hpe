// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c " ../../../../../bin/render -out examples/resources/morpheus_app_blueprint_arm/resource_json.tf app_blueprint_arm_resource_json.tf.tmpl BlueprintContent '{\n \"type\": \"Microsoft.Storage/storageAccounts\",\n \"apiVersion\": \"2019-04-01\",\n \"name\": \"string\",\n \"location\": \"string\",\n \"tags\": {\n \"tagName1\": \"tagValue1\",\n \"tagName2\": \"tagValue2\"\n },\n \"sku\": {\n \"name\": \"string\",\n \"restrictions\": [\n {\n \"reasonCode\": \"string\"\n }\n ]\n },\n \"kind\": \"string\",\n \"identity\": {\n \"type\": \"SystemAssigned\"\n },\n \"properties\": {\n \"accessTier\": \"string\",\n \"allowBlobPublicAccess\": \"bool\",\n \"allowSharedKeyAccess\": \"bool\",\n \"azureFilesIdentityBasedAuthentication\": {\n \"activeDirectoryProperties\": {\n \"azureStorageSid\": \"string\",\n \"domainGuid\": \"string\",\n \"domainName\": \"string\",\n \"domainSid\": \"string\",\n \"forestName\": \"string\",\n \"netBiosDomainName\": \"string\"\n },\n \"directoryServiceOptions\": \"string\"\n },\n \"customDomain\": {\n \"name\": \"string\",\n \"useSubDomainName\": \"bool\"\n },\n \"encryption\": {\n \"keySource\": \"string\",\n \"keyvaultproperties\": {\n \"keyname\": \"string\",\n \"keyvaulturi\": \"string\",\n \"keyversion\": \"string\"\n },\n \"services\": {\n \"blob\": {\n \"enabled\": \"bool\"\n },\n \"file\": {\n \"enabled\": \"bool\"\n }\n }\n },\n \"isHnsEnabled\": \"bool\",\n \"largeFileSharesState\": \"string\",\n \"minimumTlsVersion\": \"string\",\n \"networkAcls\": {\n \"bypass\": \"string\",\n \"defaultAction\": \"string\",\n \"ipRules\": [\n {\n \"action\": \"Allow\",\n \"value\": \"string\"\n }\n ],\n \"virtualNetworkRules\": [\n {\n \"action\": \"Allow\",\n \"id\": \"string\",\n \"state\": \"string\"\n }\n ]\n },\n \"supportsHttpsTrafficOnly\": \"bool\"\n }\n}' Category 'armtemplates' CloudInitEnabled 'true' Description 'example arm app blueprint' InstallAgent 'true' Name 'example_app_arm_blueprint_json' OsType 'linux' SourceType 'json'"

// RenderAppBlueprintArmJSONConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderAppBlueprintArmJSONConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		//nolint:lll
		"BlueprintContent": "{\n \"type\": \"Microsoft.Storage/storageAccounts\",\n \"apiVersion\": \"2019-04-01\",\n \"name\": \"string\",\n \"location\": \"string\",\n \"tags\": {\n \"tagName1\": \"tagValue1\",\n \"tagName2\": \"tagValue2\"\n },\n \"sku\": {\n \"name\": \"string\",\n \"restrictions\": [\n {\n \"reasonCode\": \"string\"\n }\n ]\n },\n \"kind\": \"string\",\n \"identity\": {\n \"type\": \"SystemAssigned\"\n },\n \"properties\": {\n \"accessTier\": \"string\",\n \"allowBlobPublicAccess\": \"bool\",\n \"allowSharedKeyAccess\": \"bool\",\n \"azureFilesIdentityBasedAuthentication\": {\n \"activeDirectoryProperties\": {\n \"azureStorageSid\": \"string\",\n \"domainGuid\": \"string\",\n \"domainName\": \"string\",\n \"domainSid\": \"string\",\n \"forestName\": \"string\",\n \"netBiosDomainName\": \"string\"\n },\n \"directoryServiceOptions\": \"string\"\n },\n \"customDomain\": {\n \"name\": \"string\",\n \"useSubDomainName\": \"bool\"\n },\n \"encryption\": {\n \"keySource\": \"string\",\n \"keyvaultproperties\": {\n \"keyname\": \"string\",\n \"keyvaulturi\": \"string\",\n \"keyversion\": \"string\"\n },\n \"services\": {\n \"blob\": {\n \"enabled\": \"bool\"\n },\n \"file\": {\n \"enabled\": \"bool\"\n }\n }\n },\n \"isHnsEnabled\": \"bool\",\n \"largeFileSharesState\": \"string\",\n \"minimumTlsVersion\": \"string\",\n \"networkAcls\": {\n \"bypass\": \"string\",\n \"defaultAction\": \"string\",\n \"ipRules\": [\n {\n \"action\": \"Allow\",\n \"value\": \"string\"\n }\n ],\n \"virtualNetworkRules\": [\n {\n \"action\": \"Allow\",\n \"id\": \"string\",\n \"state\": \"string\"\n }\n ]\n },\n \"supportsHttpsTrafficOnly\": \"bool\"\n }\n}",
		"Category":         "armtemplates",
		"CloudInitEnabled": "true",
		"Description":      "example arm app blueprint",
		"InstallAgent":     "true",
		"Name":             "example_app_arm_blueprint_json",
		"OsType":           "linux",
		"SourceType":       "json",
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
	templatePath := filepath.Join(dir, "app_blueprint_arm_resource_json.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
