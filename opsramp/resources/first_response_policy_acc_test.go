// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFirstResponsePolicyResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		policyName := acctest.RandomName("frp")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckFirstResponsePolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccFirstResponsePolicyConfig(policyName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureFirstResponsePolicyExists(t, "opsramp_first_response_policy.test_policy"),
						resource.TestCheckResourceAttrSet("opsramp_first_response_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("opsramp_first_response_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr("opsramp_first_response_policy.test_policy", "enabled_mode", "OBSERVED"),
					),
				},
			},
		})
	})
}

func testAccFirstResponsePolicyConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_first_response_policy" "test_policy" {
	name = "%s"

	enabled_mode = "OBSERVED"
	filter_query = ""

	pattern_actions = {
		seasonality_time_frame = "7D"
		suppress = {
			seasonal_alerts = true
		}
	}
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccEnsureFirstResponsePolicyExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetFirstResponsePolicy(tenantID, id)
		if err != nil {
			return fmt.Errorf("first response policy %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckFirstResponsePolicyDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_first_response_policy" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetFirstResponsePolicy(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("first response policy still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no first response policy exists") {
				return fmt.Errorf("unexpected error checking deleted first response policy %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
