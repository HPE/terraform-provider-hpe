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

func TestAccKBCategoryResource(t *testing.T) {
	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("create_and_import", func(t *testing.T) {
		catName := acctest.RandomName("kb-cat")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckKBCategoryDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccKBCategoryConfig(catName, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureKBCategoryExists(t, "hpe_opsramp_kb_category.test_category"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_kb_category.test_category", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_kb_category.test_category", "name", catName),
					),
				},
				// ImportState testing
				{
					ResourceName:      "hpe_opsramp_kb_category.test_category",
					ImportState:       true,
					ImportStateIdFunc: testAccKBCategoryImportStateIdFunc("hpe_opsramp_kb_category.test_category", clientOverride),
					ImportStateVerify: true,
				},
			},
		})
	})
}

func testAccKBCategoryConfig(name string, clientOverride string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_kb_category" "test_category" {
	name        = "%s"
	description = "Acceptance test KB category"
	%s
}
`, acctest.ProviderConfigHCL(), name, acctest.ClientAttrHCL(clientOverride))
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

		tenantID, _ := acctest.LookupProviderEnv("tenant")

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
			if rs.Type != "hpe_opsramp_kb_category" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			category, err := apiClient.GetKBCategory(tenantID, rs.Primary.ID)
			if category != nil && category.State != "TRASH" {
				return fmt.Errorf(
					"KB category still exists: %s (%s), object %+v",
					rs.Primary.ID,
					rs.Primary.Attributes["name"],
					category,
				)
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

func testAccKBCategoryImportStateIdFunc(resourceName string, clientOverride string) resource.ImportStateIdFunc {
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
