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

func TestAccFirstResponsePolicyResource(t *testing.T) {
	acctest.SkipIfNotClient(t)

	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("create", func(t *testing.T) {
		policyName := acctest.RandomName("frp")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckFirstResponsePolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccFirstResponsePolicyConfig(policyName, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureFirstResponsePolicyExists(t, "hpe_opsramp_first_response_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_first_response_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_first_response_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr("hpe_opsramp_first_response_policy.test_policy", "enabled_mode", "OBSERVED"),
					),
				},
			},
		})
	})
}

func testAccFirstResponsePolicyConfig(name string, clientOverride string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_first_response_policy" "test_policy" {
	name = "%s"
	%s

	enabled_mode = "OBSERVED"
	filter_query = ""

	pattern_actions = {
		seasonality_time_frame = "7D"
		suppress = {
			seasonal_alerts = true
		}
	}
}
`, acctest.ProviderConfigHCL(), name, acctest.ClientAttrHCL(clientOverride))
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

		tenantID, _ := acctest.LookupProviderEnv("tenant")

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
			if rs.Type != "hpe_opsramp_first_response_policy" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

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
