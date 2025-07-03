// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

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
  default_group_access = "full"
}
`
	expectJSONWithWhitespace := `
{
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
  "globalSiteAccess": "full"
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
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role_permissions.testacc_permissions_json_ok",
			"default_group_access",
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
