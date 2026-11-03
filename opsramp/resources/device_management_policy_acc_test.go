// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

func TestAccDeviceManagementPolicyResource(t *testing.T) {
	clientOverride := acctest.RequireClientScope(t)

	t.Run("create_and_update", func(t *testing.T) {
		name := acctest.RandomName("tf-device-mgmt-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckDeviceManagementPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccDeviceManagementPolicyConfig(name, false, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureDeviceManagementPolicyExists(t, "hpe_opsramp_device_management_policy.test"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_device_management_policy.test", "id"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_device_management_policy.test", "uid"),
						resource.TestCheckResourceAttr("hpe_opsramp_device_management_policy.test", "name", name),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_device_management_policy.test", "actions.0.action", "ASSIGN MONITORINGTEMPLATE",
						),
						resource.TestCheckResourceAttrPair(
							"hpe_opsramp_device_management_policy.test", "actions.0.items.0.item",
							"hpe_opsramp_template.test_policy_snmp", "id",
						),
					),
				},
				{
					Config: testAccDeviceManagementPolicyConfig(name, true, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureDeviceManagementPolicyExists(t, "hpe_opsramp_device_management_policy.test"),
						resource.TestCheckResourceAttr("hpe_opsramp_device_management_policy.test", "name", name),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_device_management_policy.test", "actions.0.force_assign_or_unassign", "true",
						),
					),
				},
			},
		})
	})

	t.Run("requires_item_or_item_uid", func(t *testing.T) {
		name := acctest.RandomName("tf-device-mgmt-policy-err")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      testAccDeviceManagementPolicyMissingItemConfig(name, clientOverride),
					ExpectError: regexp.MustCompile(`(?s)[Aa]t least one attribute out of.*item_uid`),
				},
			},
		})
	})

	t.Run("requires_item_type_for_assign_jobs", func(t *testing.T) {
		name := acctest.RandomName("tf-device-mgmt-policy-jobs-err")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      testAccDeviceManagementPolicyAssignJobsMissingItemTypeConfig(name, clientOverride),
					ExpectError: regexp.MustCompile(`item_type is required for each item when action is 'Assign Jobs'`),
				},
			},
		})
	})
}

func testAccDeviceManagementPolicyConfig(name string, forceAssign bool, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_template" "test_policy_snmp" {
  name        = "%s-template"
  description = "Terraform-managed SNMP template for device management policy tests."
  %s

  template_type = "SNMP"

  metric_groups = [
    {
      name      = "snmp metric group"
      frequency = 5

      adapter_type = "SNMP"

      is_multi_instance = false

      metrics = [
        {
          name           = "account.percent.used.capacity"
          display_name   = "Percent used capacity"
          description    = "Percent used capacity"
          category       = "METRIC"
          threshold_type = "STATIC"
          availability   = true

		  format_metric_value_details = {
            units = ""
            graph = true
            format_metric_value = {
              used_by = "NONE"
              alerts  = false
              graph   = false
            }
          }

          metric_values = [
            {
              name      = "hrStorageUsed"
              oid       = "1.3.6.1.2.1.25.2.3.1.6"
              oid_type  = "SCALAROID"
              data_type = "COUNTER64"
            },
            {
              name     = "hrStorageSize"
              oid      = "1.3.6.1.2.1.25.2.3.1.5"
              oid_type = "SCALAROID"
            }
          ]
          expression = "(hrStorageUsed / hrStorageSize) * 100"

          thresholds = [
            {
              level_type            = "TEMPLATE"
              raise_alert           = true
              warning_operator      = "NONE"
              critical_operator     = "NONE"
              warning_repeat_count  = 1
              critical_repeat_count = 1
            }
          ]
        }
      ]
    }
  ]
}

resource "hpe_opsramp_device_management_policy" "test" {
  name = "%s"
  %s

  search_query = "name = \"testing-iat-resource\""

  actions = [
    {
      action                   = "ASSIGN MONITORINGTEMPLATE"
      force_assign_or_unassign = %t

      items = [
        {
          item = tostring(hpe_opsramp_template.test_policy_snmp.id)
        }
      ]
    }
  ]
}
`, acctest.ProviderConfigHCL(), name, clientAttr, name, clientAttr, forceAssign)
}

func testAccDeviceManagementPolicyMissingItemConfig(name string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_device_management_policy" "test_err" {
  name = "%s"
  %s

  search_query = "name = \"testing-iat-resource\""

  actions = [
    {
      action = "Assign Availability Rule"
      items = [
        {}
      ]
    }
  ]
}
`, acctest.ProviderConfigHCL(), name, clientAttr)
}

func testAccDeviceManagementPolicyAssignJobsMissingItemTypeConfig(name string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_device_management_policy" "test_jobs_err" {
  name = "%s"
  %s

  search_query = "name = \"testing-iat-resource\""

  actions = [
    {
      action = "Assign Jobs"
      items = [
        {
          item     = "assetinfo"
          item_uid = "assetinfo"
        }
      ]
    }
  ]
}
`, acctest.ProviderConfigHCL(), name, clientAttr)
}

func testAccEnsureDeviceManagementPolicyExists(t *testing.T, resourceName string) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		uid := strings.TrimSpace(rs.Primary.Attributes["uid"])
		if uid == "" {
			return fmt.Errorf("resource uid is empty in state for %s", resourceName)
		}

		tenantID, _ := acctest.LookupProviderEnv("tenant")
		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetDeviceManagementPolicy(tenantID, uid)
		if err != nil {
			return fmt.Errorf("device management policy %s (%s) was not found in opsramp api: %w", resourceName, uid, err)
		}

		return nil
	}
}

func testAccCheckDeviceManagementPolicyDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_device_management_policy" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			uid := strings.TrimSpace(rs.Primary.Attributes["uid"])
			if uid == "" {
				continue
			}

			_, err := apiClient.GetDeviceManagementPolicy(tenantID, uid)
			if err == nil {
				return fmt.Errorf("device management policy still exists: %s", uid)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no policy exists with id") {
				return fmt.Errorf("unexpected error checking deleted device management policy %s: %w", uid, err)
			}
		}

		return nil
	}
}
