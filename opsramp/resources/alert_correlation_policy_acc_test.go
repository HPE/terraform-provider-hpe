// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

func TestAccAlertCorrelationPolicyResource(t *testing.T) {
	clientOverride := acctest.OptionalClientOverride(t)
	isMSPCreds := acctest.EffectiveScope(t) == "MSP"

	t.Run("similarity", func(t *testing.T) {
		policyName := acctest.RandomName("esc-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertCorrelationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertCorrelationPolicySimilarityConfig(policyName, clientOverride, isMSPCreds, ""),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertCorrelationPolicyExists(t, "hpe_opsramp_alert_correlation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_correlation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"enabled_mode",
							"OBSERVED",
						),
					),
				},
				{
					ResourceName: "hpe_opsramp_alert_correlation_policy.test_policy",
					ImportState:  true,
					ImportStateIdFunc: testAccAlertCorrelationPolicyImportStateIdFunc(
						"hpe_opsramp_alert_correlation_policy.test_policy",
						clientOverride,
					),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"organization_matching_type"},
				},
			},
		})
	})

	t.Run("topology", func(t *testing.T) {
		policyName := acctest.RandomName("esc-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertCorrelationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertCorrelationPolicyTopologyConfig(policyName, clientOverride, isMSPCreds),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertCorrelationPolicyExists(t, "hpe_opsramp_alert_correlation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_correlation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"enabled_mode",
							"OBSERVED",
						),
					),
				},

				{
					ResourceName: "hpe_opsramp_alert_correlation_policy.test_policy",
					ImportState:  true,
					ImportStateIdFunc: testAccAlertCorrelationPolicyImportStateIdFunc(
						"hpe_opsramp_alert_correlation_policy.test_policy",
						clientOverride,
					),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"organization_matching_type"},
				},
			},
		})
	})

	t.Run("algorithm", func(t *testing.T) {
		policyName := acctest.RandomName("esc-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertCorrelationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					// alerts_time_window set as a number.
					Config: testAccAlertCorrelationPolicyAlgorithmConfig(policyName, clientOverride, isMSPCreds, "15"),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertCorrelationPolicyExists(t, "hpe_opsramp_alert_correlation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_correlation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"enabled_mode",
							"ON",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"type",
							"ALGORITHM",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"algorithm_correlation.alerts_time_window",
							"15",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"algorithm_correlation.matching_conditions.#",
							"1",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"algorithm_correlation.matching_conditions.0.property",
							"resource_name",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"algorithm_correlation.matching_conditions.0.match_type",
							"Identical",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"inference_subject",
							"Multiple business alerts on $resource.name",
						),
					),
				},
				{
					// alerts_time_window set as a quoted string holding the same number:
					// Terraform converts it to the schema's number type, so switching between
					// the two forms for an unchanged value plans no update.
					Config: testAccAlertCorrelationPolicyAlgorithmConfig(policyName, clientOverride, isMSPCreds, `"15"`),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(
								"hpe_opsramp_alert_correlation_policy.test_policy",
								plancheck.ResourceActionNoop,
							),
						},
					},
					Check: resource.TestCheckResourceAttr(
						"hpe_opsramp_alert_correlation_policy.test_policy",
						"algorithm_correlation.alerts_time_window",
						"15",
					),
				},
				{
					ResourceName: "hpe_opsramp_alert_correlation_policy.test_policy",
					ImportState:  true,
					ImportStateIdFunc: testAccAlertCorrelationPolicyImportStateIdFunc(
						"hpe_opsramp_alert_correlation_policy.test_policy",
						clientOverride,
					),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"organization_matching_type"},
				},
			},
		})
	})
}

func TestAccAlertCorrelationPolicyResourceWithInclude(t *testing.T) {
	acctest.SkipIfNotMSPWithClientOverride(t)
	includedClient := acctest.OptionalClientOverride(t)

	t.Run("similarity_with_include", func(t *testing.T) {
		policyName := acctest.RandomName("esc-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertCorrelationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertCorrelationPolicySimilarityConfig(policyName, "", true, includedClient),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertCorrelationPolicyExists(t, "hpe_opsramp_alert_correlation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_correlation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"organization_matching_type",
							"INCLUDE",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_correlation_policy.test_policy",
							"included_clients.#",
							"1",
						),
					),
				},
			},
		})
	})
}

func testAccAlertCorrelationPolicySimilarityConfig(policyName, clientOverride string, isMSPCreds bool, targetClient string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)
	orgMatchingType := ""
	if isMSPCreds {
		if targetClient != "" {
			orgMatchingType = fmt.Sprintf("organization_matching_type = \"INCLUDE\"\n  included_clients = [\"%s\"]", targetClient)
		} else {
			orgMatchingType = `organization_matching_type = "ALL"`
		}
	}

	return fmt.Sprintf(`
%s

resource "hpe_opsramp_alert_correlation_policy" "test_policy" {
  %s
  %s
  name = "%s"
  enabled_mode    = "OBSERVED"
  filter_query    = ""
  inference_query = ""
  type            = "CO_OCCURRENCE"
  machine_learning = {
    continuous_learning = true
    topology            = false
    matching_conditions = [
      {
        property   = "service_group"
        match_type = "Identical"
      }
    ]
  }

  inference_subject = ""

  
}
`, acctest.ProviderConfigHCL(), clientAttr, orgMatchingType, policyName)
}

func testAccAlertCorrelationPolicyTopologyConfig(policyName, clientOverride string, isMSPCreds bool) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)
	orgMatchingType := ""
	if isMSPCreds {
		orgMatchingType = `  organization_matching_type = "ALL"`
	}

	return fmt.Sprintf(`
%s

resource "hpe_opsramp_alert_correlation_policy" "test_policy" {
  %s
  %s
  name = "%s"
  enabled_mode    = "OBSERVED"
  filter_query    = ""
  inference_query = ""
  type            = "CO_OCCURRENCE"
  machine_learning = {
    continuous_learning = true
    topology            = true
    topology_depth      = 3
  }

  inference_subject = ""
}
`, acctest.ProviderConfigHCL(), clientAttr, orgMatchingType, policyName)
}

func testAccAlertCorrelationPolicyAlgorithmConfig(policyName, clientOverride string, isMSPCreds bool, alertsTimeWindow string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)
	orgMatchingType := ""
	if isMSPCreds {
		orgMatchingType = `  organization_matching_type = "ALL"`
	}

	return fmt.Sprintf(`
%s

resource "hpe_opsramp_alert_correlation_policy" "test_policy" {
  %s
  %s
  name         = "%s"
  enabled_mode = "ON"
  type         = "ALGORITHM"
  algorithm_correlation = {
    alerts_time_window = %s
    matching_conditions = [
      {
        property   = "resource_name"
        match_type = "Identical"
      }
    ]
  }
  filter_query      = "metric STARTS WITH \"demo\" OR metric CONTAINS \"trace\" OR metric CONTAINS \"Rate\" OR problemArea CONTAINS \"Log_Alert\""
  inference_query   = ""
  inference_subject = "Multiple business alerts on $resource.name"
}
`, acctest.ProviderConfigHCL(), clientAttr, orgMatchingType, policyName, alertsTimeWindow)
}

func testAccEnsureAlertCorrelationPolicyExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		tenantID := apiClient.TenantId
		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		_, err = apiClient.GetAlertCorrelationPolicy(tenantID, id)
		if err != nil {
			return fmt.Errorf("alert correlation policy %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckAlertCorrelationPolicyDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_alert_correlation_policy" {
				continue
			}

			tenantID := apiClient.TenantId
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetAlertCorrelationPolicy(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("alert correlation policy still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no alert correlation policy exists") {
				return fmt.Errorf("unexpected error checking deleted alert correlation policy %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}

func testAccAlertCorrelationPolicyImportStateIdFunc(resourceName, clientPrefix string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		id := rs.Primary.ID
		if clientPrefix != "" {
			return clientPrefix + ":" + id, nil
		}

		return id, nil
	}
}
