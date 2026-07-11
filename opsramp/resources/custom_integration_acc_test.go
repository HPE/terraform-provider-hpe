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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCustomIntegrationResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		displayName := acctest.RandomName("integration")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckCustomIntegrationDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccCustomIntegrationConfig(displayName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureCustomIntegrationExists(t, "opsramp_custom_integration.test_integration"),
						resource.TestCheckResourceAttrSet("opsramp_custom_integration.test_integration", "id"),
						resource.TestCheckResourceAttr("opsramp_custom_integration.test_integration", "display_name", displayName),
					),
				},
			},
		})
	})
}

func testAccCustomIntegrationConfig(displayName string) string {
	return fmt.Sprintf(`
%s

resource "opsramp_permission_set" "test_role_perms" {
	name        = "Test permissions for %s"
	description = "Permissions for role test"

	permissions = [
		{
			name = "Alerts"
			type = "View Alerts "
		}
	]
}

resource "opsramp_role" "test_role" {
	name        = "Test role for %s"
	description = "Acceptance test role"

	permissions = [
		opsramp_permission_set.test_role_perms.unique_id
	]
}
resource "opsramp_custom_integration" "test_integration" {
	display_name = "%s"
	role_name    = opsramp_role.test_role.name
}
`, acctest.ProviderConfigHCL(), displayName, displayName, displayName)
}

func testAccEnsureCustomIntegrationExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetCustomIntegration(tenantID, id)
		if err != nil {
			return fmt.Errorf("custom integration %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckCustomIntegrationDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_custom_integration" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetCustomIntegration(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("custom integration still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no installed integration found with id") {
				return fmt.Errorf("unexpected error checking deleted custom integration %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
