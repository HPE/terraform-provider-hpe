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

func TestAccCredentialSetResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		credentialSetName := acctest.RandomName("credentialSet")
		description := acctest.RandomName("description")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckCredentialSetDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccCredentialSetConfig(credentialSetName, description),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureCredentialSetExists(t, "hpe_opsramp_credential_set.test_credential_set"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_credential_set.test_credential_set", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_credential_set.test_credential_set", "name", credentialSetName),
						resource.TestCheckResourceAttr("hpe_opsramp_credential_set.test_credential_set", "description", description),
					),
				},
			},
		})
	})
}

func testAccCredentialSetConfig(name string, description string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_credential_set" "test_credential_set" {
  name                = "%s"
  description         = "%s"

  credential_type     = "VMWARE"
  user_name           = "administrator"
  password            = "**********"
  port                = 443

  timeout_ms		  = 15000
  security_level      = "NOAUTHNOPRIV"
  snmp_version        = "V2"
  ssh_credential_type = "PASSWORD"
  transport_type      = "HTTP"
}
`, acctest.ProviderConfigHCL(), name, description)
}

func testAccEnsureCredentialSetExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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
			return fmt.Errorf("failed to initialize opsramp api credentialSet: %w", err)
		}

		_, err = apiClient.GetCredentialSet(tenantID, id)
		if err != nil {
			return fmt.Errorf("credentialSet %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckCredentialSetDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api credentialSet: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_credential_set" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetCredentialSet(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("credentialSet still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no credential set found") {
				return fmt.Errorf("unexpected error checking deleted credentialSet %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
