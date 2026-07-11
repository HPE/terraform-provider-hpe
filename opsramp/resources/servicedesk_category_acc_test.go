// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccServicedeskCategoryResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		catName := acctest.RandomName("sd-cat")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckServicedeskCategoryDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccServicedeskCategoryConfig(catName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureServicedeskCategoryExists(t, "opsramp_servicedesk_category.test_category"),
						resource.TestCheckResourceAttrSet("opsramp_servicedesk_category.test_category", "id"),
						resource.TestCheckResourceAttr("opsramp_servicedesk_category.test_category", "name", catName),
						resource.TestCheckResourceAttr("opsramp_servicedesk_category.test_category", "ticket_type", "serviceRequests"),
					),
				},
			},
		})
	})
}

func testAccServicedeskCategoryConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_servicedesk_category" "test_category" {
	name        = "%s"
	description = "Acceptance test category"
	ticket_type = "serviceRequests"
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccEnsureServicedeskCategoryExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetServiceDeskCategory(id)
		if err != nil {
			return fmt.Errorf("servicedesk category %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckServicedeskCategoryDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_servicedesk_category" {
				continue
			}

			category, err := apiClient.GetServiceDeskCategory(rs.Primary.ID)
			if category != nil {
				return fmt.Errorf("servicedesk category still exists: %s, object: %+v", rs.Primary.ID, category)
			}

			if err != nil {
				errText := strings.ToLower(err.Error())
				if !strings.Contains(errText, "404") && !strings.Contains(errText, "not found") {
					return fmt.Errorf("unexpected error checking deleted servicedesk category %s: %w", rs.Primary.ID, err)
				}
			}
		}

		return nil
	}
}
