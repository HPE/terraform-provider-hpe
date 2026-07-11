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

func TestAccKBCategoryResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		catName := acctest.RandomName("kb-cat")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckKBCategoryDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccKBCategoryConfig(catName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureKBCategoryExists(t, "opsramp_kb_category.test_category"),
						resource.TestCheckResourceAttrSet("opsramp_kb_category.test_category", "id"),
						resource.TestCheckResourceAttr("opsramp_kb_category.test_category", "name", catName),
					),
				},
			},
		})
	})
}

func testAccKBCategoryConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_kb_category" "test_category" {
	name        = "%s"
	description = "Acceptance test KB category"
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccEnsureKBCategoryExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetKBCategory(tenantID, id)
		if err != nil {
			return fmt.Errorf("KB category %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckKBCategoryDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_kb_category" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			category, err := apiClient.GetKBCategory(tenantID, rs.Primary.ID)
			if category != nil && category.State != "TRASH" {
				return fmt.Errorf("KB category still exists: %s (%s), object %+v", rs.Primary.ID, rs.Primary.Attributes["name"], category)
			}

			if err != nil {
				errText := strings.ToLower(err.Error())
				if !strings.Contains(errText, "no category found with id") {
					return fmt.Errorf("unexpected error checking deleted KB category %s: %w", rs.Primary.ID, err)
				}
			}
		}

		return nil
	}
}
