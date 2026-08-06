// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

func TestAccCustomIntegrationResource(t *testing.T) {
	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("create_and_import", func(t *testing.T) {
		displayName := acctest.RandomName("integration")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckCustomIntegrationDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccCustomIntegrationConfig(displayName, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureCustomIntegrationExists(t, "hpe_opsramp_custom_integration.test_integration"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_custom_integration.test_integration", "id"),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_custom_integration.test_integration",
							"display_name",
							displayName,
						),
					),
				},
				// ImportState testing
				{
					ResourceName: "hpe_opsramp_custom_integration.test_integration",
					ImportState:  true,
					ImportStateIdFunc: testAccCustomIntegrationImportStateIdFunc(
						"hpe_opsramp_custom_integration.test_integration",
						clientOverride,
					),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"api_client_id", "api_client_secret", "role_name"},
				},
			},
		})
	})
}

func testAccCustomIntegrationConfig(displayName string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s

resource "hpe_opsramp_permission_set" "test_role_perms" {
	name        = "Test permissions for %s"
	description = "Permissions for role test"
	%s

	permissions = [
		{
			name = "Alerts"
			type = "View Alerts "
		}
	]
}

resource "hpe_opsramp_role" "test_role" {
	name        = "Test role for %s"
	description = "Acceptance test role"
	%s

	permissions = [
		hpe_opsramp_permission_set.test_role_perms.unique_id
	]
}
resource "hpe_opsramp_custom_integration" "test_integration" {
	display_name = "%s"
	role_name    = hpe_opsramp_role.test_role.name
	%s
}
`, acctest.ProviderConfigHCL(), displayName, clientAttr, displayName, clientAttr, displayName, clientAttr)
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

		tenantID, _ := acctest.LookupProviderEnv("tenant")

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
			if rs.Type != "hpe_opsramp_custom_integration" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

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

func testAccCustomIntegrationImportStateIdFunc(resourceName string, clientOverride string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		id := rs.Primary.ID
		if clientOverride != "" {
			return clientOverride + ":" + id, nil
		}

		return id, nil
	}
}
