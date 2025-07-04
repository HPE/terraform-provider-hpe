// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package rolepermissions_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"

	"github.com/stretchr/testify/assert"
)

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hpe": newProviderWithError,
}

// tests that the JSON body from setting a permissions config is as expected
func TestAccMorpheusDataSourceRolePermissionsJsonOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	// we're testing the construction of a JSON body from a config,
	// so we'll include both user AND tenant permissions
	dataSourceConfig := `
data "hpe_morpheus_role_permissions" "testacc_permissions_json_ok" {
  feature_permissions = [
    {
      "code"   = "activity"
      "access" = "full"
    },
    {
      "code"   = "admin-accounts"
      "access" = "full"
    }
  ]
  cloud_permissions = [
    {
      "id"   = 1
      "access" = "full"
    }
  ]
  group_permissions = [
    {
      "id"   = 1
      "access" = "full"
    }
  ]
  blueprint_permissions = [
    {
      "id"   = 1
      "access" = "full"
    }
  ]
  instance_type_permissions = [
    {
      "id"   = 1
      "access" = "full"
    }
  ]
  persona_permissions = [
    {
      "code"   = "standard"
      "access" = "full"
    }
  ]
  report_type_permissions = [
    {
      "code"   = "appCost"
      "access" = "full"
    }
  ]
  task_permissions = [
    {
      "id"   = 1
      "access" = "full"
    }
  ]
  workflow_permissions = [
    {
      "id"   = 1
      "access" = "full"
    }
  ]
  vdi_pool_permissions = [
    {
      "id"   = 1
      "access" = "full"
    }
  ]
  default_group_access = "full"
  default_cloud_access = "full"
  default_blueprint_access = "full"
  default_catalog_item_type_access = "full"
  default_instance_type_access = "full"
  default_persona_access = "full"
  default_report_type_access = "full"
  default_task_access = "full"
  default_workflow_access = "full"
  default_vdi_pool_access = "full"
}
`
	// the json Marshal performed will sort the keys lexicographically
	expectJSONWithWhitespace := `
{
  "appTemplatePermissions": [
    {
      "access": "full",
      "id": 1
    }
  ],
  "featurePermissions": [
      {
        "access": "full",
        "code": "activity"
      },
      {
        "access": "full",
        "code" : "admin-accounts"
      }
    ],
  "globalAppTemplateAccess": "full",
  "globalCatalogItemTypeAccess": "full",
  "globalInstanceTypeAccess": "full",
  "globalPersonaAccess": "full",
  "globalReportTypeAccess": "full",
  "globalSiteAccess": "full",
  "globalTaskAccess": "full",
  "globalTaskSetAccess": "full",
  "globalVdiPoolAccess": "full",
  "globalZoneAccess": "full",
  "instanceTypePermissions": [
    {
      "access": "full",
      "id": 1
    }
  ],
  "personaPermissions": [
    {
      "access": "full",
      "code": "standard"
    }
  ],
  "reportTypePermissions": [
    {
      "access": "full",
      "code": "appCost"
    }
  ],
  "sites": [
    {
      "access": "full",
      "id": 1
    }
  ],
  "taskPermissions": [
    {
      "access": "full",
      "id": 1
    }
  ],
  "taskSetPermissions": [
    {
      "access": "full",
      "id": 1
    }
  ],
  "vdiPoolPermissions": [
    {
      "access": "full",
      "id": 1
    }
  ],
  "zones": [
    {
      "access": "full",
      "id": 1
    }
  ]
}
`
	var bufJSONCompact bytes.Buffer
	err := json.Compact(&bufJSONCompact, []byte(expectJSONWithWhitespace))
	assert.NoError(t, err)

	expectJSON := bufJSONCompact.String()

	checks := []resource.TestCheckFunc{
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"feature_permissions.*",
			map[string]string{
				"code":   "activity",
				"access": "full",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"feature_permissions.*",
			map[string]string{
				"code":   "admin-accounts",
				"access": "full",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"cloud_permissions.*",
			map[string]string{
				"id":     "1",
				"access": "full",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"group_permissions.*",
			map[string]string{
				"id":     "1",
				"access": "full",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"blueprint_permissions.*",
			map[string]string{
				"id":     "1",
				"access": "full",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"instance_type_permissions.*",
			map[string]string{
				"id":     "1",
				"access": "full",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"persona_permissions.*",
			map[string]string{
				"code":   "standard",
				"access": "full",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"report_type_permissions.*",
			map[string]string{
				"code":   "appCost",
				"access": "full",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"task_permissions.*",
			map[string]string{
				"id":     "1",
				"access": "full",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"workflow_permissions.*",
			map[string]string{
				"id":     "1",
				"access": "full",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"vdi_pool_permissions.*",
			map[string]string{
				"id":     "1",
				"access": "full",
			},
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_group_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_cloud_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_blueprint_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_catalog_item_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_instance_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_persona_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_report_type_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_task_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_workflow_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_vdi_pool_access",
			"full",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"json",
			expectJSON,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigOffline + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

// tests that when none of the optional properties are set that we get an empty JSON string
func TestAccMorpheusDataSourceRolePermissionsNoneSetOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	dataSourceConfig := `
data "hpe_morpheus_role_permissions" "testacc_permissions_none_set_ok" {}
`

	// a json.Marshal on an empty struct produces "{}"
	expectJson := "{}"

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"feature_permissions",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"cloud_permissions",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"group_permissions",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"blueprint_permissions",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"instance_type_permissions",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"persona_permissions",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"report_type_permissions",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"task_permissions",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"workflow_permissions",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"vdi_pool_permissions",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"default_group_access",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"default_cloud_access",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"default_blueprint_access",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"default_catalog_item_type_access",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"default_instance_type_access",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"default_persona_access",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"default_report_type_access",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"default_task_access",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"default_workflow_access",
		),
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"default_vdi_pool_access",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_none_set_ok",
			"json",
			expectJson,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigOffline + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

// test that we can use the permissions data source to create a user role
func TestAccMorpheusDataSourceRolePermissionsUserRoleOk(t *testing.T) {
	//

}
