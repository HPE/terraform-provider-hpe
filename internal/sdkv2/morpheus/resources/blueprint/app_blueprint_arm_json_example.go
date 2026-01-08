// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/app_blueprint_arm/app_blueprint_arm_resource_json.tf app_blueprint_arm_resource_json.tf.tmpl Category "armtemplates" CloudInitEnabled "true" Description "example arm app blueprint" InstallAgent "true" Name "example_app_blueprint_json" OsType "linux" SourceType "json" BlueprintContent "..."

// RenderAppBlueprintArmJSONConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderAppBlueprintArmJSONConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	//nolint:lll
	defaults := map[string]string{
		"BlueprintContent": `{
  "type": "Microsoft.Storage/storageAccounts",
  "apiVersion": "2019-04-01",
  "name": "string",
  "location": "string",
  "tags": {
    "tagName1": "tagValue1",
    "tagName2": "tagValue2"
  },
  "sku": {
    "name": "string",
    "restrictions": [
      {
        "reasonCode": "string"
      }
    ]
  },
  "kind": "string",
  "identity": {
    "type": "SystemAssigned"
  },
  "properties": {
    "accessTier": "string",
    "allowBlobPublicAccess": "bool",
    "allowSharedKeyAccess": "bool",
    "azureFilesIdentityBasedAuthentication": {
      "activeDirectoryProperties": {
        "azureStorageSid": "string",
        "domainGuid": "string",
        "domainName": "string",
        "domainSid": "string",
        "forestName": "string",
        "netBiosDomainName": "string"
      },
      "directoryServiceOptions": "string"
    },
    "customDomain": {
      "name": "string",
      "useSubDomainName": "bool"
    },
    "encryption": {
      "keySource": "string",
      "keyvaultproperties": {
        "keyname": "string",
        "keyvaulturi": "string",
        "keyversion": "string"
      },
      "services": {
        "blob": {
          "enabled": "bool"
        },
        "file": {
          "enabled": "bool"
        }
      }
    },
    "isHnsEnabled": "bool",
    "largeFileSharesState": "string",
    "minimumTlsVersion": "string",
    "networkAcls": {
      "bypass": "string",
      "defaultAction": "string",
      "ipRules": [
        {
          "action": "Allow",
          "value": "string"
        }
      ],
      "virtualNetworkRules": [
        {
          "action": "Allow",
          "id": "string",
          "state": "string"
        }
      ]
    },
    "supportsHttpsTrafficOnly": "bool"
  }
}`,
		"Category":         "armtemplates",
		"CloudInitEnabled": "true",
		"Description":      "example arm app blueprint",
		"InstallAgent":     "true",
		"Name":             "example_app_blueprint_json",
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
