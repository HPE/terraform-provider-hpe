// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

func TestAccResource(t *testing.T) {
	acctest.SkipIfNotClient(t)
	clientOverride := acctest.RequireClientScope(t)

	// t.Run("happy path", func(t *testing.T) {
	resourceName := acctest.RandomName("resource")
	hostname := acctest.RandomName("host")

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceConfig(resourceName, "one", hostname, clientOverride),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hpe_opsramp_resource.test_resource",
						tfjsonpath.New("resource_name"),
						knownvalue.StringExact(resourceName),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureResourceExists(t, "hpe_opsramp_resource.test_resource"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_resource.test_resource", "uuid"),
					resource.TestCheckResourceAttr("hpe_opsramp_resource.test_resource", "alias_name", "one"),
					resource.TestCheckResourceAttr("hpe_opsramp_resource.test_resource", "resource_name", resourceName),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "hpe_opsramp_resource.test_resource",
				ImportState:                          true,
				ImportStateIdFunc:                    testAccResourceImportStateIdFunc("hpe_opsramp_resource.test_resource", clientOverride),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
			// Update and Read testing
			{
				Config: testAccResourceConfig(resourceName, "onetwo", hostname, clientOverride),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"hpe_opsramp_resource.test_resource",
						tfjsonpath.New("resource_name"),
						knownvalue.StringExact(resourceName),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureResourceExists(t, "hpe_opsramp_resource.test_resource"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_resource.test_resource", "uuid"),
					resource.TestCheckResourceAttr("hpe_opsramp_resource.test_resource", "alias_name", "onetwo"),
					resource.TestCheckResourceAttr("hpe_opsramp_resource.test_resource", "resource_name", resourceName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
	//})
}

func testAccResourceConfig(name string, alias string, hostname string, clientOverride string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_resource" "test_resource" {
	resource_name = "%s"
	alias_name = "%s"
	hostname = "%s"
	resource_type = "Linux"
	%s
}`, acctest.ProviderConfigHCL(), name, alias, hostname, acctest.ClientAttrHCL(clientOverride))
}

func testAccEnsureResourceExists(t *testing.T, resourceName string) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		resourceUuid := strings.TrimSpace(rs.Primary.Attributes["uuid"])
		if resourceUuid == "" {
			return fmt.Errorf("resource uuid is empty in state for %s", resourceName)
		}

		tenantID, _ := acctest.LookupProviderEnv("tenant")

		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetResource(tenantID, resourceUuid)
		if err != nil {
			return fmt.Errorf("resource %s (%s) was not found in opsramp api: %w", resourceName, resourceUuid, err)
		}

		return nil
	}
}

func testAccResourceImportStateIdFunc(resourceName string, clientOverride string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		id := rs.Primary.Attributes["uuid"]
		if clientOverride != "" {
			return clientOverride + ":" + id, nil
		}

		return id, nil
	}
}
