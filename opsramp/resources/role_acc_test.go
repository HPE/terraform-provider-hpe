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

func TestAccRoleResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		permName := acctest.RandomName("role-perm")
		roleName := acctest.RandomName("role")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckRoleDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccRoleConfig(permName, roleName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureRoleExists(t, "opsramp_role.test_role"),
						resource.TestCheckResourceAttrSet("opsramp_role.test_role", "id"),
						resource.TestCheckResourceAttr("opsramp_role.test_role", "name", roleName),
					),
				},
			},
		})
	})
}

func testAccRoleConfig(permName string, roleName string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_permission_set" "test_role_perms" {
	name        = "%s"
	description = "Permissions for role test"

	permissions = [
		{
			name = "Alerts"
			type = "View Alerts "
		}
	]
}

resource "opsramp_role" "test_role" {
	name        = "%s"
	description = "Acceptance test role"

	permissions = [
		opsramp_permission_set.test_role_perms.unique_id
	]
}
`, acctest.ProviderConfigHCL(), permName, roleName)
}

func testAccEnsureRoleExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetRole(tenantID, id)
		if err != nil {
			return fmt.Errorf("role %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckRoleDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_role" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetRole(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("role still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no role found with id") {
				return fmt.Errorf("unexpected error checking deleted role %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
