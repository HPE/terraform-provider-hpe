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

func TestAccAlertCorrelationPolicyResource(t *testing.T) {
	// CLIENT scope — direct CLIENT credentials, no client override
	t.Run("client_similarity", func(t *testing.T) {
		acctest.SkipIfNotClient(t)
		policyName := acctest.RandomName("esc-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertCorrelationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertCorrelationPolicySimilarityConfig(policyName, ""),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertCorrelationPolicyExists(t, "hpe_opsramp_alert_correlation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_correlation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "enabled_mode", "OBSERVED"),
					),
				},
				// Import — CLIENT scope: import ID is just the policy ID
				{
					ResourceName:      "hpe_opsramp_alert_correlation_policy.test_policy",
					ImportState:       true,
					ImportStateIdFunc: testAccAlertCorrelationPolicyImportStateIdFunc("hpe_opsramp_alert_correlation_policy.test_policy", ""),
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("client_topology", func(t *testing.T) {
		acctest.SkipIfNotClient(t)
		policyName := acctest.RandomName("esc-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertCorrelationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertCorrelationPolicyTopologyConfig(policyName, ""),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertCorrelationPolicyExists(t, "hpe_opsramp_alert_correlation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_correlation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "enabled_mode", "OBSERVED"),
					),
				},
				// Import — CLIENT scope: import ID is just the policy ID
				{
					ResourceName:      "hpe_opsramp_alert_correlation_policy.test_policy",
					ImportState:       true,
					ImportStateIdFunc: testAccAlertCorrelationPolicyImportStateIdFunc("hpe_opsramp_alert_correlation_policy.test_policy", ""),
					ImportStateVerify: true,
				},
			},
		})
	})

	// MSP scope — MSP credentials alone, no client override (policy at MSP level)
	t.Run("msp_similarity", func(t *testing.T) {
		acctest.SkipIfNotMSP(t)
		policyName := acctest.RandomName("esc-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertCorrelationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertCorrelationPolicyMSPSimilarityConfig(policyName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertCorrelationPolicyExists(t, "hpe_opsramp_alert_correlation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_correlation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "enabled_mode", "OBSERVED"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "organization_matching_type", "ALL"),
					),
				},
				// Import — MSP scope without client override: import ID is just the policy ID
				{
					ResourceName:      "hpe_opsramp_alert_correlation_policy.test_policy",
					ImportState:       true,
					ImportStateIdFunc: testAccAlertCorrelationPolicyImportStateIdFunc("hpe_opsramp_alert_correlation_policy.test_policy", ""),
					ImportStateVerify: true,
				},
			},
		})
	})

	// MSP scope + target client override (policy pushed to a specific client)
	t.Run("msp_target_similarity", func(t *testing.T) {
		acctest.SkipIfNotMSP(t)
		targetClient := acctest.TargetClientID(t)
		policyName := acctest.RandomName("esc-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertCorrelationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertCorrelationPolicyMSPTargetSimilarityConfig(policyName, targetClient),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertCorrelationPolicyExists(t, "hpe_opsramp_alert_correlation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_correlation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "enabled_mode", "OBSERVED"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "organization_matching_type", "ALL"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "client", targetClient),
					),
				},
				// Import — MSP scope with client prefix: import ID is <client_id>:<policy_id>
				{
					ResourceName:      "hpe_opsramp_alert_correlation_policy.test_policy",
					ImportState:       true,
					ImportStateIdFunc: testAccAlertCorrelationPolicyImportStateIdFunc("hpe_opsramp_alert_correlation_policy.test_policy", targetClient),
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("msp_target_topology", func(t *testing.T) {
		acctest.SkipIfNotMSP(t)
		targetClient := acctest.TargetClientID(t)
		policyName := acctest.RandomName("esc-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertCorrelationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertCorrelationPolicyMSPTargetTopologyConfig(policyName, targetClient),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertCorrelationPolicyExists(t, "hpe_opsramp_alert_correlation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_correlation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "enabled_mode", "OBSERVED"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "organization_matching_type", "ALL"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.test_policy", "client", targetClient),
					),
				},
				// Import — MSP scope with client prefix: import ID is <client_id>:<policy_id>
				{
					ResourceName:      "hpe_opsramp_alert_correlation_policy.test_policy",
					ImportState:       true,
					ImportStateIdFunc: testAccAlertCorrelationPolicyImportStateIdFunc("hpe_opsramp_alert_correlation_policy.test_policy", targetClient),
					ImportStateVerify: true,
				},
			},
		})
	})
}

func testAccAlertCorrelationPolicySimilarityConfig(policyName, clientOverride string) string {
	clientAttr := ""
	if clientOverride != "" {
		clientAttr = fmt.Sprintf(`  client = "%s"`, clientOverride)
	}

	return fmt.Sprintf(`
%s

resource "hpe_opsramp_alert_correlation_policy" "test_policy" {
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
`, acctest.ProviderConfigHCL(), clientAttr, policyName)
}

func testAccAlertCorrelationPolicyTopologyConfig(policyName, clientOverride string) string {
	clientAttr := ""
	if clientOverride != "" {
		clientAttr = fmt.Sprintf(`  client = "%s"`, clientOverride)
	}

	return fmt.Sprintf(`
%s

resource "hpe_opsramp_alert_correlation_policy" "test_policy" {
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
`, acctest.ProviderConfigHCL(), clientAttr, policyName)
}

func testAccAlertCorrelationPolicyMSPSimilarityConfig(policyName string) string {
	return fmt.Sprintf(`
%s

resource "hpe_opsramp_alert_correlation_policy" "test_policy" {
  organization_matching_type = "ALL"
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
`, acctest.ProviderConfigHCL(), policyName)
}

func testAccAlertCorrelationPolicyMSPTargetSimilarityConfig(policyName, targetClient string) string {
	return fmt.Sprintf(`
%s

resource "hpe_opsramp_alert_correlation_policy" "test_policy" {
  client                     = "%s"
  organization_matching_type = "ALL"
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
`, acctest.ProviderConfigHCL(), targetClient, policyName)
}

func testAccAlertCorrelationPolicyMSPTargetTopologyConfig(policyName, targetClient string) string {
	return fmt.Sprintf(`
%s

resource "hpe_opsramp_alert_correlation_policy" "test_policy" {
  client                     = "%s"
  organization_matching_type = "ALL"
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
`, acctest.ProviderConfigHCL(), targetClient, policyName)
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
