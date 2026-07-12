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

func TestAccAlertPredictionPolicyResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		policyName := acctest.RandomName("pred-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertPredictionPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertPredictionPolicyConfig(policyName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertPredictionPolicyExists(t, "hpe_opsramp_alert_prediction_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_prediction_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_prediction_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_prediction_policy.test_policy", "seasonality_time_frame", "7D"),
					),
				},
			},
		})
	})
}

func testAccAlertPredictionPolicyConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_alert_prediction_policy" "test_policy" {
	name = "%s"

	enabled_mode = "OFF"
	filter_query = ""

	seasonality_time_frame    = "7D"
	generate_prediction_alert = true
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccEnsureAlertPredictionPolicyExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetAlertPredictionPolicy(tenantID, id)
		if err != nil {
			return fmt.Errorf("alert prediction policy %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckAlertPredictionPolicyDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_alert_prediction_policy" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetAlertPredictionPolicy(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("alert prediction policy still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no alert prediction policy exists") {
				return fmt.Errorf("unexpected error checking deleted alert prediction policy %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
