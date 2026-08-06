// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

func TestAccClientResource(t *testing.T) {
	// MSP-only resource: only runs in Run 1 (MSP creds, no target_client).
	// Skipped in Run 2 (MSP + target_client) and Run 3 (CLIENT creds).
	acctest.SkipIfNotMSP(t)

	t.Run("create", func(t *testing.T) {
		clientName := acctest.RandomName("client")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			// TODO: find a way to check destroy for clients, as they are deleted asynch-
			// ronously and may take a long time to be fully removed from the system.
			// CheckDestroy:             testAccCheckClientDestroy(t),
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
		// Allow time for the client to be fully provisioned.
		time.Sleep(10 * time.Second)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		id := strings.TrimSpace(rs.Primary.Attributes["unique_id"])
		if id == "" {
			return fmt.Errorf("resource unique_id is empty in state for %s", resourceName)
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
