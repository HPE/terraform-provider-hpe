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

func TestAccPermissionSetResource(t *testing.T) {
	// Any-level resource: runs in all three runs.
	//   Run 1 (MSP, no target_client): created at MSP level, no client attr.
	//   Run 2 (MSP + target_client):   created under target client, client attr set.
	//   Run 3 (CLIENT creds):          created directly, no client attr.

	t.Run("create", func(t *testing.T) {
		clientOverride := acctest.OptionalClientOverride(t)
		permName := acctest.RandomName("permset")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckPermissionSetDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccPermissionSetConfig(permName, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsurePermissionSetExists(t, "hpe_opsramp_permission_set.test_perms"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_permission_set.test_perms", "unique_id"),
						resource.TestCheckResourceAttr("hpe_opsramp_permission_set.test_perms", "name", permName),
					),
				},
			},
		})
	})
}

func testAccPermissionSetConfig(name string, clientOverride string) string {
	clientAttr := ""
	if clientOverride != "" {
		clientAttr = fmt.Sprintf(`client = "%s"`, clientOverride)
	}

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_permission_set" "test_perms" {
	name        = "%s"
	description = "Acceptance test permission set"
	%s

	permissions = [
		{
			name = "Alerts"
			type = "View Alerts "
		},
		{
			name = "Devices"
			type = "View Device "
		}
	]
}
`, acctest.ProviderConfigHCL(), name, clientAttr)
}

func testAccEnsurePermissionSetExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetPermissionSet(tenantID, id)
		if err != nil {
			return fmt.Errorf("permission set %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckPermissionSetDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_permission_set" {
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

			_, err := apiClient.GetPermissionSet(tenantID, id)
			if err == nil {
				return fmt.Errorf("permission set still exists: %s", id)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no permissionset found") {
				return fmt.Errorf("unexpected error checking deleted permission set %s: %w", id, err)
			}
		}

		return nil
	}
}
