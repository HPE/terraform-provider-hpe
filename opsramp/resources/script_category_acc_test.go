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

func TestAccScriptCategoryResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		catName := acctest.RandomName("script-cat")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckScriptCategoryDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccScriptCategoryConfig(catName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureScriptCategoryExists(t, "hpe_opsramp_script_category.test_category"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_script_category.test_category", "uuid"),
						resource.TestCheckResourceAttr("hpe_opsramp_script_category.test_category", "name", catName),
					),
				},
			},
		})
	})
}

func TestAccScriptCategoryWithParentResource(t *testing.T) {
	t.Run("parent and child", func(t *testing.T) {
		parentName := acctest.RandomName("script-parent")
		childName := acctest.RandomName("script-child")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckScriptCategoryDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccScriptCategoryWithParentConfig(parentName, childName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureScriptCategoryExists(t, "hpe_opsramp_script_category.test_parent"),
						testAccEnsureScriptCategoryExists(t, "hpe_opsramp_script_category.test_child"),
						resource.TestCheckResourceAttr("hpe_opsramp_script_category.test_parent", "name", parentName),
						resource.TestCheckResourceAttr("hpe_opsramp_script_category.test_child", "name", childName),
					),
				},
			},
		})
	})
}

func testAccScriptCategoryConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_script_category" "test_category" {
	name = "%s"
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccScriptCategoryWithParentConfig(parentName string, childName string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_script_category" "test_parent" {
	name = "%s"
}

resource "hpe_opsramp_script_category" "test_child" {
	name      = "%s"
	parent_id = hpe_opsramp_script_category.test_parent.uuid
}
`, acctest.ProviderConfigHCL(), parentName, childName)
}

func testAccEnsureScriptCategoryExists(t *testing.T, resourceName string) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		id := strings.TrimSpace(rs.Primary.Attributes["uuid"])
		if id == "" {
			return fmt.Errorf("resource uuid is empty in state for %s", resourceName)
		}

		// The tenant lookup result is not checked here: PreCheck already
		// validates the required environment, and an explicit client attribute
		// on the resource takes precedence below.
		tenantID, _ := acctest.LookupProviderEnv("tenant")
		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetScriptCategory(tenantID, id)
		if err != nil {
			return fmt.Errorf("script category %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckScriptCategoryDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_script_category" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			id := strings.TrimSpace(rs.Primary.Attributes["uuid"])
			if id == "" {
				continue
			}

			_, err := apiClient.GetScriptCategory(tenantID, id)
			if err == nil {
				return fmt.Errorf("script category still exists: %s", id)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no task category found with id") {
				return fmt.Errorf("unexpected error checking deleted script category %s: %w", id, err)
			}
		}

		return nil
	}
}
