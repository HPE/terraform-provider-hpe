// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ---------------------------------------------------------------------------
// hpe_opsramp_integration_config – SNMP integration with schedule
// ---------------------------------------------------------------------------

func TestAccIntegrationConfigResource_WithSchedule(t *testing.T) {
	configName := acctest.RandomName("intg-cfg")
	configNameUpdated := configName + "-upd"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckIntegrationConfigDestroy(t),
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccIntegrationConfigWithScheduleConfig(configName, "DAILY", 1, "01"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureIntegrationConfigExists(t, "hpe_opsramp_integration_config.test"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_integration_config.test", "id"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_config.test", "name", configName),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_config.test", "schedule.pattern_type", "DAILY"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_config.test", "schedule.pattern", "1"),
				),
			},
			// Update – rename and change schedule
			{
				Config: testAccIntegrationConfigWithScheduleConfig(configNameUpdated, "HOURLY", 2, "3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_opsramp_integration_config.test", "name", configNameUpdated),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_config.test", "schedule.pattern_type", "HOURLY"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_config.test", "schedule.pattern", "2"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// hpe_opsramp_integration_config – config without schedule (all_resources)
// ---------------------------------------------------------------------------

func TestAccIntegrationConfigResource_NoSchedule(t *testing.T) {
	configName := acctest.RandomName("intg-cfg-nosched")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acctest.PreCheck(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckIntegrationConfigDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConfigNoScheduleConfig(configName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureIntegrationConfigExists(t, "hpe_opsramp_integration_config.test_nosched"),
					resource.TestCheckResourceAttrSet("hpe_opsramp_integration_config.test_nosched", "id"),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_config.test_nosched", "name", configName),
					resource.TestCheckResourceAttr("hpe_opsramp_integration_config.test_nosched", "all_resources", "true"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Config helpers
// ---------------------------------------------------------------------------

func testAccIntegrationConfigWithScheduleConfig(name string, patternType string, pattern int, startTime string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_management_profile" "test_management_profile" {
  name = "Test Management Profile - 1"
  description = "Management profile for integration config acceptance tests"
}
resource "hpe_opsramp_integration" "config_parent" {
  application  = "SNMP"
  display_name = "tf-acc-snmp-config-parent - 1"
  profile_id = hpe_opsramp_management_profile.test_management_profile.uuid
}
resource "hpe_opsramp_credential_set" "snmp_credential_set" {
  name        = "SNMP Credential Set - 1"
  description = "Credential set for SNMP tests"

  credential_type = "SNMP"
  password        = "**********"
  port            = 161
  
  security_level       = "NOAUTHNOPRIV"
  transport_type       = "HTTP"

  secure              = true
  timeout_ms          = 15000
  snmp_version        = "V2"
  ssh_credential_type = "PASSWORD"
  community           = "public"
}

resource "hpe_opsramp_integration_config" "test" {
  integration_id = hpe_opsramp_integration.config_parent.id
  name           = %q
  config         = jsonencode({
    nmapResult      = true
    deviceType      = "SNMP Network Device"
    discoveryType   = "Iprange"
    ipRange         = "10.0.0.1"
    networkDepth    = "1"
	credentials     = [hpe_opsramp_credential_set.snmp_credential_set.id]
	packetCount     = "default"
  })

  all_resources = false

  schedule = {
    pattern_type = %q
    pattern      = %d
	start_time  = %q
  }
}
`, acctest.ProviderConfigHCL(), name, patternType, pattern, startTime)
}

func testAccIntegrationConfigNoScheduleConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_management_profile" "test_management_profile" {
  name = "Test Management Profile - 2"
  description = "Management profile for integration config acceptance tests"
}
resource "hpe_opsramp_integration" "config_parent" {
  application  = "SNMP"
  display_name = "tf-acc-snmp-config-parent 2"
  profile_id = hpe_opsramp_management_profile.test_management_profile.uuid
}
resource "hpe_opsramp_credential_set" "snmp_credential_set" {
  name        = "SNMP Credential Set - 2"
  description = "Credential set for SNMP tests"

  credential_type = "SNMP"
  password        = "**********"
  port            = 161

  security_level       = "NOAUTHNOPRIV"
  transport_type       = "HTTP"

  secure              = true
  timeout_ms          = 15000
  snmp_version        = "V2"
  ssh_credential_type = "PASSWORD"
  community           = "public"
}
resource "hpe_opsramp_integration_config" "test_nosched" {
  integration_id = hpe_opsramp_integration.config_parent.id
  name           = %q
  config         = jsonencode({
    nmapResult      = true
    deviceType      = "SNMP Network Device"
    discoveryType   = "Iprange"
    ipRange         = "10.0.1.1"
    networkDepth    = "1"
    credentials     = [hpe_opsramp_credential_set.snmp_credential_set.id]
	packetCount     = "default"
  })
  all_resources = true
}
`, acctest.ProviderConfigHCL(), name)
}

// ---------------------------------------------------------------------------
// Shared check helpers
// ---------------------------------------------------------------------------

func testAccEnsureIntegrationConfigExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		integrationID := rs.Primary.Attributes["integration_id"]

		tenantID, _ := acctest.LookupProviderEnv("tenant")

		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetIntegrationConfig(tenantID, integrationID, id)
		if err != nil {
			return fmt.Errorf("integration config %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckIntegrationConfigDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_integration_config" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			integrationID := rs.Primary.Attributes["integration_id"]

			_, err := apiClient.GetIntegrationConfig(tenantID, integrationID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("integration config still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no installed integration found with") {
				return fmt.Errorf("unexpected error checking deleted integration config %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
