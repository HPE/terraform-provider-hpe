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

func TestAccClientResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		acctest.SkipIfNotMSP(t)
		clientName := acctest.RandomName("client")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckClientDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccClientConfig(clientName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureClientExists(t, "hpe_opsramp_client.test_client"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_client.test_client", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_client.test_client", "name", clientName),
						resource.TestCheckResourceAttr("hpe_opsramp_client.test_client", "country", "Spain"),
						resource.TestCheckResourceAttr("hpe_opsramp_client.test_client", "time_zone", "Europe/Paris"),
					),
				},
			},
		})
	})
}

func testAccClientConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_client" "test_client" {
	name      = "%s"
	address   = "Valencia, Spain"
	country   = "Spain"
	time_zone = "Europe/Paris"

	packages = [
		"Hybrid Discovery and Monitoring",
		"Event and Incident Management",
		"Remediation and Automation"
	]
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccEnsureClientExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetClient(id)
		if err != nil {
			return fmt.Errorf("client %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckClientDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_client" {
				continue
			}

			_, err := apiClient.GetClient(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("client still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "404") && !strings.Contains(errText, "not found") {
				return fmt.Errorf("unexpected error checking deleted client %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
