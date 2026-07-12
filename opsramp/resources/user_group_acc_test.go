// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
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
						testAccEnsureUserGroupExists(t, "hpe_opsramp_user_group.test_group"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_user_group.test_group", "unique_id"),
						resource.TestCheckResourceAttr("hpe_opsramp_user_group.test_group", "name", groupName),
					),
				},
			},
		})
	})
}

func testAccUserGroupConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_user_group" "test_group" {
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

		tenantID, _ := acctest.LookupProviderEnv("tenant")

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
			if rs.Type != "hpe_opsramp_user_group" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

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
