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

func TestAccScriptResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		catName := acctest.RandomName("script-cat")
		scriptName := acctest.RandomName("script")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckScriptDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccScriptConfig(catName, scriptName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureScriptExists(t, "hpe_opsramp_script.test_script"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_script.test_script", "uuid"),
						resource.TestCheckResourceAttr("hpe_opsramp_script.test_script", "name", scriptName),
						resource.TestCheckResourceAttr("hpe_opsramp_script.test_script", "execution_type", "SHELL"),
					),
				},
				// ImportState testing — import ID is <category_id>:<script_uuid>
				{
					ResourceName:                         "hpe_opsramp_script.test_script",
					ImportState:                          true,
					ImportStateIdFunc:                    testAccScriptImportStateIdFunc("hpe_opsramp_script.test_script"),
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "uuid",
					ImportStateVerifyIgnore:              []string{"attachment"},
				},
			},
		})
	})
}

func testAccScriptConfig(catName string, scriptName string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_script_category" "script_test_cat" {
	name = "%s"
}

resource "hpe_opsramp_script" "test_script" {
	category_id     = hpe_opsramp_script_category.script_test_cat.uuid
	name            = "%s"
	description     = "Acceptance test script"
	platforms       = ["LINUX"]
	execution_type  = "SHELL"
	install_timeout = 120

	attachment = {
		name = "test_script.sh"
		file = "#!/bin/bash\necho hello"
	}

	parameters = [
		{
			name          = "param1"
			description   = "Test parameter"
			default_value = "default"
			type          = "REQUIRED"
			data_type     = "STRING"
		}
	]
}
`, acctest.ProviderConfigHCL(), catName, scriptName)
}

func testAccEnsureScriptExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		tenantID, _ := acctest.LookupProviderEnv("tenant")

		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		categoryID := rs.Primary.Attributes["category_id"]

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetScript(tenantID, categoryID, id)
		if err != nil {
			return fmt.Errorf("script %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckScriptDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_script" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			categoryID := rs.Primary.Attributes["category_id"]

			id := strings.TrimSpace(rs.Primary.Attributes["uuid"])
			if id == "" {
				continue
			}

			_, err := apiClient.GetScript(tenantID, categoryID, id)
			if err == nil {
				return fmt.Errorf("script still exists: %s", id)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no task found with") {
				return fmt.Errorf("unexpected error checking deleted script %s: %w", id, err)
			}
		}

		return nil
	}
}

// testAccScriptImportStateIdFunc builds the composite import ID: <category_id>:<script_uuid>
func testAccScriptImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		categoryID := rs.Primary.Attributes["category_id"]
		scriptUUID := rs.Primary.Attributes["uuid"]

		return categoryID + ":" + scriptUUID, nil
	}
}
