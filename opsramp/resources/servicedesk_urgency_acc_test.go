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

func TestAccServicedeskUrgencyResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		urgencyName := acctest.RandomName("sd-urgency")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckServicedeskUrgencyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccServicedeskUrgencyConfig(urgencyName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureServicedeskUrgencyExists(t, "hpe_opsramp_servicedesk_urgency.test_urgency"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_servicedesk_urgency.test_urgency", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_servicedesk_urgency.test_urgency", "name", urgencyName),
					),
				},
			},
		})
	})
}

func testAccServicedeskUrgencyConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_servicedesk_urgency" "test_urgency" {
	name        = "%s"
	description = "Acceptance test urgency"
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccEnsureServicedeskUrgencyExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetServiceDeskUrgency(id)
		if err != nil {
			return fmt.Errorf("servicedesk urgency %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckServicedeskUrgencyDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_servicedesk_urgency" {
				continue
			}

			urgency, err := apiClient.GetServiceDeskUrgency(rs.Primary.ID)
			if urgency != nil {
				return fmt.Errorf("servicedesk urgency still exists: %s, object: %+v", rs.Primary.ID, urgency)
			}

			if err != nil {
				errText := strings.ToLower(err.Error())
				if !strings.Contains(errText, "404") && !strings.Contains(errText, "not found") {
					return fmt.Errorf("unexpected error checking deleted servicedesk urgency %s: %w", rs.Primary.ID, err)
				}
			}
		}

		return nil
	}
}
