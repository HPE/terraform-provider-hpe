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

func TestAccServicedeskBusinessImpactResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		impactName := acctest.RandomName("sd-impact")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckServicedeskBusinessImpactDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccServicedeskBusinessImpactConfig(impactName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureServicedeskBusinessImpactExists(t, "hpe_opsramp_servicedesk_business_impact.test_impact"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_servicedesk_business_impact.test_impact", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_servicedesk_business_impact.test_impact", "name", impactName),
					),
				},
			},
		})
	})
}
func testAccServicedeskBusinessImpactConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_servicedesk_business_impact" "test_impact" {
	name        = "%s"
	description = "Acceptance test business impact"
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccEnsureServicedeskBusinessImpactExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetServiceDeskBusinessImpact(id)
		if err != nil {
			return fmt.Errorf("servicedesk business impact %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckServicedeskBusinessImpactDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_servicedesk_business_impact" {
				continue
			}

			businessImpact, err := apiClient.GetServiceDeskBusinessImpact(rs.Primary.ID)
			if businessImpact != nil {
				return fmt.Errorf("servicedesk business impact still exists: %s, object: %+v", rs.Primary.ID, businessImpact)
			}

			if err != nil {
				errText := strings.ToLower(err.Error())
				if !strings.Contains(errText, "404") && !strings.Contains(errText, "not found") {
					return fmt.Errorf("unexpected error checking deleted servicedesk business impact %s: %w", rs.Primary.ID, err)
				}
			}
		}

		return nil
	}
}
