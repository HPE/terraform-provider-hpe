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

// ---------------------------------------------------------------------------
// hpe_opsramp_integration – CUSTOM-EVENT (inbound-only with mapping attributes)
// ---------------------------------------------------------------------------

func TestAccIntegrationResource_CustomEvent(t *testing.T) {
	displayName := acctest.RandomName("intg-custom-event")
	displayNameUpdated := displayName + "-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckIntegrationDestroy(t),
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccIntegrationCustomEventConfig(displayName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureIntegrationExists(t, "hpe_opsramp_integration.test"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_integration.test", "id"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration.test", "display_name", displayName),
					resource.TestCheckResourceAttr("hpe_opsramp_integration.test", "application", "CUSTOM-EVENT"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration.test", "category", "Monitoring"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration.test", "inbound.enable_drop_alerts", "true"),
				),
			},
			// Update – change display_name and toggle enable_drop_alerts
			{
				Config: testAccIntegrationCustomEventConfigUpdated(displayNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_opsramp_integration.test", "display_name", displayNameUpdated),
					resource.TestCheckResourceAttr("hpe_opsramp_integration.test", "inbound.enable_drop_alerts", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "hpe_opsramp_integration.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccIntegrationImportStateIdFunc("hpe_opsramp_integration.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"application", "inbound"},
			},
		},
	})
}

// ---------------------------------------------------------------------------
// hpe_opsramp_integration – CUSTOM (inbound OAUTH2 + outbound REST_API)
// ---------------------------------------------------------------------------

func TestAccIntegrationResource_Custom(t *testing.T) {
	displayName := acctest.RandomName("intg-custom")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckIntegrationDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationCustomConfig(displayName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureIntegrationExists(t, "hpe_opsramp_integration.test_custom"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_integration.test_custom", "id"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration.test_custom", "application", "CUSTOM"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration.test_custom", "inbound.auth_type", "OAUTH2"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration.test_custom", "outbound.auth_type", "BASIC"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "hpe_opsramp_integration.test_custom",
				ImportState:             true,
				ImportStateIdFunc:       testAccIntegrationImportStateIdFunc("hpe_opsramp_integration.test_custom"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"application", "inbound", "outbound"},
			},
		},
	})
}

// ---------------------------------------------------------------------------
// hpe_opsramp_integration – NEWRELIC (pre-configured, inbound auto-provisioned)
// ---------------------------------------------------------------------------

func TestAccIntegrationResource_NewRelic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckIntegrationDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationNewRelicConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureIntegrationExists(t, "hpe_opsramp_integration.test_newrelic"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_integration.test_newrelic", "id"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration.test_newrelic", "application", "NEWRELIC"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_integration.test_newrelic", "inbound.token"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_integration.test_newrelic", "inbound.webhook_url"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Config helpers
// ---------------------------------------------------------------------------

func testAccIntegrationCustomEventConfig(displayName string) string {
	return fmt.Sprintf(`
%s

data "hpe_opsramp_custom_event_alert_source" "custom_source" {
  name      = "Custom"
}

resource "hpe_opsramp_integration" "test" {
  display_name    = %q
  application     = "CUSTOM-EVENT"
  alert_source_id = data.hpe_opsramp_custom_event_alert_source.custom_source.id

  inbound = {
    auth_type          = "WEBHOOK"
    enable_drop_alerts = true

    map_attributes = [
      {
        third_party_attribute = "alert_time"
        opsramp_attribute     = "alert.alertTime"
      },
    ]
  }
}
`, acctest.ProviderConfigHCL(), displayName)
}

func testAccIntegrationCustomEventConfigUpdated(displayName string) string {
	return fmt.Sprintf(`
%s
data "hpe_opsramp_custom_event_alert_source" "custom_source" {
  name      = "Custom"
}
  
resource "hpe_opsramp_integration" "test" {
  display_name    = %q
  application     = "CUSTOM-EVENT"
  alert_source_id = data.hpe_opsramp_custom_event_alert_source.custom_source.id

  inbound = {
    auth_type          = "WEBHOOK"
    enable_drop_alerts = false

    map_attributes = [
      {
        third_party_attribute = "alert_time"
        opsramp_attribute     = "alert.alertTime"
      },
    ]
  }
}
`, acctest.ProviderConfigHCL(), displayName)
}

func testAccIntegrationCustomConfig(displayName string) string {
	return fmt.Sprintf(`
%s

resource "hpe_opsramp_integration" "test_custom" {
  display_name = %q
  application  = "CUSTOM"
  category     = "Custom"

  inbound = {
    auth_type = "OAUTH2"
  }

  outbound = {
    base_uri  = "https://httpbin.org/post"
    auth_type = "BASIC"
    username  = "api-user"
    password  = "secret"
  }
}
`, acctest.ProviderConfigHCL(), displayName)
}

func testAccIntegrationNewRelicConfig() string {
	return fmt.Sprintf(`
%s

resource "hpe_opsramp_integration" "test_newrelic" {
  application = "NEWRELIC"

  inbound = {
    auth_type          = "WEBHOOK"
    enable_drop_alerts = false
  }
}
`, acctest.ProviderConfigHCL())
}

// ---------------------------------------------------------------------------
// Shared check helpers
// ---------------------------------------------------------------------------

func testAccEnsureIntegrationExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetIntegration(tenantID, id)
		if err != nil {
			return fmt.Errorf("integration %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckIntegrationDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_integration" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetIntegration(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("integration still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no installed integration found with id") {
				return fmt.Errorf("unexpected error checking deleted integration %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}

func testAccIntegrationImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
