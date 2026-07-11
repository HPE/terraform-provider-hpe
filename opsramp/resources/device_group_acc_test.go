// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-opsramp/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDeviceGroupResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		groupNameOne := acctest.RandomName("device-group")
		groupNameTwo := acctest.RandomName("device-group")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckDeviceGroupDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccDeviceGroupConfig(groupNameOne, "resourceType = \"Server\""),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"opsramp_device_group.test_group",
							tfjsonpath.New("name"),
							knownvalue.StringExact(groupNameOne),
						),
						statecheck.ExpectKnownValue(
							"opsramp_device_group.test_group",
							tfjsonpath.New("entity_type"),
							knownvalue.StringExact("DEVICE_GROUP"),
						),
					},
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureDeviceGroupExists(t, "opsramp_device_group.test_group"),
						resource.TestCheckResourceAttrSet("opsramp_device_group.test_group", "id"),
						resource.TestCheckResourceAttr("opsramp_device_group.test_group", "search_query", "resourceType = \"Server\""),
					),
				},
				{
					Config: testAccDeviceGroupConfig(groupNameTwo, "name CONTAINS \"updated\""),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"opsramp_device_group.test_group",
							tfjsonpath.New("name"),
							knownvalue.StringExact(groupNameTwo),
						),
						statecheck.ExpectKnownValue(
							"opsramp_device_group.test_group",
							tfjsonpath.New("entity_type"),
							knownvalue.StringExact("DEVICE_GROUP"),
						),
					},
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureDeviceGroupExists(t, "opsramp_device_group.test_group"),
						resource.TestCheckResourceAttr("opsramp_device_group.test_group", "search_query", "name CONTAINS \"updated\""),
						resource.TestCheckResourceAttr("opsramp_device_group.test_group", "name", groupNameTwo),
					),
				},
			},
		})
	})
}

func testAccDeviceGroupConfig(name string, searchQuery string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_device_group" "test_group" {
	name         = "%s"
	search_query = %q
}
`, acctest.ProviderConfigHCL(), name, searchQuery)
}

func testAccEnsureDeviceGroupExists(t *testing.T, resourceName string) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		id := strings.TrimSpace(rs.Primary.ID)
		if id == "" {
			return fmt.Errorf("resource id is empty in state for %s", resourceName)
		}

		tenantID := os.Getenv("OPSRAMP_TENANT")
		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetDeviceGroup(tenantID, id)
		if err != nil {
			return fmt.Errorf("device group %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckDeviceGroupDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_device_group" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetDeviceGroup(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("device group still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no device group exists") {
				return fmt.Errorf("unexpected error checking deleted device group %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
