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

func TestAccServicemapResource(t *testing.T) {
	// CLIENT-only resource: runs in Run 2 (MSP + target_client) and Run 3 (CLIENT creds).
	// Skipped in Run 1 (MSP creds without target_client).
	acctest.SkipIfNotClient(t)
	clientOverride := acctest.RequireClientScope(t)

	t.Run("create_and_import", func(t *testing.T) {
		smName := acctest.RandomName("servicemap")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckServicemapDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccServicemapConfig(smName, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureServicemapExists(t, "hpe_opsramp_servicemap.test_sm"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_servicemap.test_sm", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_servicemap.test_sm", "name", smName),
						resource.TestCheckResourceAttr("hpe_opsramp_servicemap.test_sm", "type", "Service"),
					),
				},
				// ImportState testing
				{
					ResourceName:      "hpe_opsramp_servicemap.test_sm",
					ImportState:       true,
					ImportStateIdFunc: testAccServicemapImportStateIdFunc("hpe_opsramp_servicemap.test_sm", clientOverride),
					ImportStateVerify: true,
				},
			},
		})
	})
}

func TestAccServicemapWithChildResource(t *testing.T) {
	acctest.SkipIfNotClient(t)
	clientOverride := acctest.RequireClientScope(t)

	t.Run("parent_and_child", func(t *testing.T) {
		rootName := acctest.RandomName("sm-root")
		childName := acctest.RandomName("sm-child")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckServicemapDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccServicemapWithChildConfig(rootName, childName, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureServicemapExists(t, "hpe_opsramp_servicemap.test_sm_root"),
						testAccEnsureServicemapExists(t, "hpe_opsramp_servicemap.test_sm_child"),
						resource.TestCheckResourceAttr("hpe_opsramp_servicemap.test_sm_root", "name", rootName),
						resource.TestCheckResourceAttr("hpe_opsramp_servicemap.test_sm_child", "name", childName),
					),
				},
			},
		})
	})
}

func testAccServicemapConfig(name string, clientOverride string) string {
	clientAttr := ""
	if clientOverride != "" {
		clientAttr = fmt.Sprintf(`client = "%s"`, clientOverride)
	}

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_servicemap" "test_sm" {
	name = "%s"
	type = "Service"
	%s
}
`, acctest.ProviderConfigHCL(), name, clientAttr)
}

func testAccServicemapWithChildConfig(rootName string, childName string, clientOverride string) string {
	clientAttr := ""
	if clientOverride != "" {
		clientAttr = fmt.Sprintf(`client = "%s"`, clientOverride)
	}

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_servicemap" "test_sm_root" {
	name = "%s"
	type = "Service"
	%s
}

resource "hpe_opsramp_servicemap" "test_sm_child" {
	name   = "%s"
	type   = "Service"
	parent = hpe_opsramp_servicemap.test_sm_root.id
	%s
}
`, acctest.ProviderConfigHCL(), rootName, clientAttr, childName, clientAttr)
}

func testAccEnsureServicemapExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetServicemap(tenantID, id)
		if err != nil {
			return fmt.Errorf("servicemap %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckServicemapDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_servicemap" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetServicemap(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("servicemap still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no service group exists") {
				return fmt.Errorf("unexpected error checking deleted servicemap %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}

func testAccServicemapImportStateIdFunc(resourceName string, clientOverride string) resource.ImportStateIdFunc {
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
