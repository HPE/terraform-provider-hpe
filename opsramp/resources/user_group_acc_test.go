// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccUserGroupResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		groupName := acctest.RandomName("usergroup")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckUserGroupDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccUserGroupConfig(groupName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureUserGroupExists(t, "opsramp_user_group.test_group"),
						resource.TestCheckResourceAttrSet("opsramp_user_group.test_group", "unique_id"),
						resource.TestCheckResourceAttr("opsramp_user_group.test_group", "name", groupName),
					),
				},
			},
		})
	})
}

func testAccUserGroupConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_user_group" "test_group" {
	name        = "%s"
	description = "Acceptance test user group"
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccEnsureUserGroupExists(t *testing.T, resourceName string) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		id := strings.TrimSpace(rs.Primary.Attributes["unique_id"])
		if id == "" {
			return fmt.Errorf("resource unique_id is empty in state for %s", resourceName)
		}

		tenantID := os.Getenv("OPSRAMP_TENANT")
		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetUserGroup(tenantID, id)
		if err != nil {
			return fmt.Errorf("user group %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckUserGroupDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_user_group" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			id := strings.TrimSpace(rs.Primary.Attributes["unique_id"])
			if id == "" {
				continue
			}

			_, err := apiClient.GetUserGroup(tenantID, id)
			if err == nil {
				return fmt.Errorf("user group still exists: %s", id)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no usergroup found") {
				return fmt.Errorf("unexpected error checking deleted user group %s: %w", id, err)
			}
		}

		return nil
	}
}
