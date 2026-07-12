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

func TestAccSiteResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		siteNameOne := acctest.RandomName("site")
		siteNameTwo := acctest.RandomName("site")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckSiteDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccSiteConfig(
						siteNameOne,
						"Initial site description",
						"Av. del General Avilés, 35-37, Benicalap,  València, Valencia",
						"Valencia",
						"València",
						"Spain",
						"46035",
						"902027020",
						"34",
						"name CONTAINS \"init\"",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureSiteExists(t, "hpe_opsramp_site.test_site"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_site.test_site", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "name", siteNameOne),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "state", "València"),
					),
				},
				{
					Config: testAccSiteConfig(
						siteNameTwo,
						"Updated site description",
						"Vicente Aleixandre, 1",
						"Las Rozas de Madrid",
						"Madrid",
						"Spain",
						"28232",
						"911237104",
						"34",
						"name CONTAINS \"updated\"",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureSiteExists(t, "hpe_opsramp_site.test_site"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "name", siteNameTwo),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "description", "Updated site description"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "address", "Vicente Aleixandre, 1"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "city", "Las Rozas de Madrid"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "state", "Madrid"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "country", "Spain"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "zip", "28232"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "phone_number", "911237104"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "phone_extension", "34"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.test_site", "search_query", "name CONTAINS \"updated\""),
					),
				},
				// ImportState testing
				{
					ResourceName:      "hpe_opsramp_site.test_site",
					ImportState:       true,
					ImportStateIdFunc: testAccSiteImportStateIdFunc("hpe_opsramp_site.test_site"),
					ImportStateVerify: true,
				},
			},
		})
	})
}

func testAccSiteConfig(name string, description string, address string, city string, state string, country string, zip string, phoneNumber string, phoneExtension string, searchQuery string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_site" "test_site" {
	name            = "%s"
	description     = "%s"
	address         = "%s"
	city            = "%s"
	state           = "%s"
	country         = "%s"
	zip             = "%s"
	phone_number    = "%s"
	phone_extension = "%s"
	search_query    = %q
}
`, acctest.ProviderConfigHCL(), name, description, address, city, state, country, zip, phoneNumber, phoneExtension, searchQuery)
}

func testAccEnsureSiteExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetSite(tenantID, id)
		if err != nil {
			return fmt.Errorf("site %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckSiteDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_site" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetSite(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("site still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no site found") {
				return fmt.Errorf("unexpected error checking deleted site %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}

func testAccSiteImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
