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

func TestAccAlertPredictionPolicyResource(t *testing.T) {
	clientOverride := acctest.OptionalClientOverride(t)
	isMSPCreds := acctest.EffectiveScope(t) == "MSP"

	t.Run("create", func(t *testing.T) {
		policyName := acctest.RandomName("pred-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertPredictionPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertPredictionPolicyConfig(policyName, clientOverride, isMSPCreds, ""),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertPredictionPolicyExists(t, "hpe_opsramp_alert_prediction_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_prediction_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_prediction_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_prediction_policy.test_policy",
							"seasonality_time_frame",
							"7D",
						),
					),
				},
			},
		})
	})
}

func TestAccAlertPredictionPolicyResourceWithInclude(t *testing.T) {
	acctest.SkipIfNotMSPWithClientOverride(t)
	includedClient := acctest.OptionalClientOverride(t)

	t.Run("create_with_include", func(t *testing.T) {
		policyName := acctest.RandomName("pred-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertPredictionPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertPredictionPolicyConfig(policyName, "", true, includedClient),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertPredictionPolicyExists(t, "hpe_opsramp_alert_prediction_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_prediction_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_prediction_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_prediction_policy.test_policy",
							"organization_matching_type",
							"INCLUDE",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_prediction_policy.test_policy",
							"included_clients.#",
							"1",
						),
					),
				},
			},
		})
	})
}

func testAccAlertPredictionPolicyConfig(name string, clientOverride string, isMSPCreds bool, targetClient string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)
	orgMatchingType := ""
	if isMSPCreds {
		if targetClient != "" {
			orgMatchingType = fmt.Sprintf("organization_matching_type = \"INCLUDE\"\n\tincluded_clients = [\"%s\"]", targetClient)
		} else {
			orgMatchingType = `organization_matching_type = "ALL"`
		}
	}

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_alert_prediction_policy" "test_policy" {
	name = "%s"
	%s
	%s

	enabled_mode = "OFF"
	filter_query = ""

	seasonality_time_frame    = "7D"
	generate_prediction_alert = true
}
`, acctest.ProviderConfigHCL(), name, clientAttr, orgMatchingType)
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
