// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-opsramp/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccManagementProfile(t *testing.T) {
	//t.Run("happy path", func(t *testing.T) {
	managementProfileName := acctest.RandomName("managementProfile")
	description := acctest.RandomName("description")

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckManagementProfileDestroy(t),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccManagementProfileConfig(managementProfileName, description),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"opsramp_management_profile.test_management_profile",
						tfjsonpath.New("name"),
						knownvalue.StringExact(managementProfileName),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureManagementProfileExists(t, "opsramp_management_profile.test_management_profile"),
					resource.TestCheckResourceAttrSet("opsramp_management_profile.test_management_profile", "uuid"),
					resource.TestCheckResourceAttrSet("opsramp_management_profile.test_management_profile", "id"),
					resource.TestCheckResourceAttr("opsramp_management_profile.test_management_profile", "name", managementProfileName),
				),
			},
			// Update and Read testing
			{
				Config: testAccManagementProfileConfig(managementProfileName, description+" updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"opsramp_management_profile.test_management_profile",
						tfjsonpath.New("name"),
						knownvalue.StringExact(managementProfileName),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureManagementProfileExists(t, "opsramp_management_profile.test_management_profile"),
					resource.TestCheckResourceAttrSet("opsramp_management_profile.test_management_profile", "uuid"),
					resource.TestCheckResourceAttrSet("opsramp_management_profile.test_management_profile", "id"),
					resource.TestCheckResourceAttr("opsramp_management_profile.test_management_profile", "name", managementProfileName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
	//})
}

func testAccManagementProfileConfig(name string, description string) string {
	return fmt.Sprintf(`
%s

resource "opsramp_management_profile" "test_management_profile" {
  name = "%s"
  description = "%s"
}`, acctest.ProviderConfigHCL(), name, description)
}

func testAccEnsureManagementProfileExists(t *testing.T, managementProfileName string) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[managementProfileName]
		if !ok {
			return fmt.Errorf("managementProfile not found in state: %s", managementProfileName)
		}

		managementProfileId, err := strconv.Atoi(rs.Primary.Attributes["id"])
		if err != nil {
			return fmt.Errorf("invalid management profile id in state for %s: %w", rs.Primary.Attributes["id"], err)
		}

		tenantID := os.Getenv("OPSRAMP_TENANT")
		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetManagementProfile(tenantID, managementProfileId)
		if err != nil {
			return fmt.Errorf("managementProfile %s (%d) was not found in opsramp api: %w", managementProfileName, managementProfileId, err)
		}

		return nil
	}
}

func testAccCheckManagementProfileDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_management_profile" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			managementProfileId, err := strconv.Atoi(rs.Primary.Attributes["id"])
			if err != nil {
				return fmt.Errorf("invalid management profile id in state for %s: %w", rs.Primary.Attributes["id"], err)
			}

			_, err = apiClient.GetManagementProfile(tenantID, managementProfileId)
			if err == nil {
				return fmt.Errorf("management profile still exists: %d", managementProfileId)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no management profile found with id") {
				return fmt.Errorf("unexpected error checking deleted management profile %d: %w", managementProfileId, err)
			}
		}

		return nil
	}
}
