// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

func TestAccScheduledMaintenanceResource(t *testing.T) {
	acctest.SkipIfNotClient(t)
	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("one-time schedule", func(t *testing.T) {
		nameOne := acctest.RandomName("sm")
		nameTwo := acctest.RandomName("sm")

		// Use future dates to ensure the maintenance window is in Pending state
		startTime := time.Now().Add(48 * time.Hour).Format("2006-01-02T15:04")
		endTime := time.Now().Add(49 * time.Hour).Format("2006-01-02T15:04")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckScheduledMaintenanceDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccScheduledMaintenanceOneTimeConfig(nameOne, "Initial maintenance window", startTime, endTime, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureScheduledMaintenanceExists(t, "hpe_opsramp_scheduled_maintenance.test"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_scheduled_maintenance.test", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_scheduled_maintenance.test", "name", nameOne),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test",
							"description",
							"Initial maintenance window",
						),
						resource.TestCheckResourceAttr("hpe_opsramp_scheduled_maintenance.test", "schedule.type", "one-time"),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test",
							"schedule.timezone",
							"Europe/Paris",
						),
						resource.TestCheckResourceAttr("hpe_opsramp_scheduled_maintenance.test", "correlate_alerts", "true"),
						resource.TestCheckResourceAttr("hpe_opsramp_scheduled_maintenance.test", "run_rba", "false"),
						resource.TestCheckResourceAttr("hpe_opsramp_scheduled_maintenance.test", "install_patch", "false"),
					),
				},
				{
					Config: testAccScheduledMaintenanceOneTimeConfig(nameTwo, "Updated maintenance window", startTime, endTime, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureScheduledMaintenanceExists(t, "hpe_opsramp_scheduled_maintenance.test"),
						resource.TestCheckResourceAttr("hpe_opsramp_scheduled_maintenance.test", "name", nameTwo),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test",
							"description",
							"Updated maintenance window",
						),
					),
				},
				// ImportState testing
				{
					ResourceName: "hpe_opsramp_scheduled_maintenance.test",
					ImportState:  true,
					ImportStateIdFunc: testAccScheduledMaintenanceImportStateIdFunc(
						"hpe_opsramp_scheduled_maintenance.test", clientOverride,
					),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"schedule.start_time", "schedule.end_time"},
				},
			},
		})
	})

	t.Run("recurring weekly schedule", func(t *testing.T) {
		name := acctest.RandomName("sm-recurring")

		startTime := time.Now().Add(72 * time.Hour).Format("2006-01-02T15:04")
		endTime := time.Now().Add(73 * time.Hour).Format("2006-01-02T15:04")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckScheduledMaintenanceDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccScheduledMaintenanceRecurringWeeklyConfig(name, startTime, endTime, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureScheduledMaintenanceExists(t, "hpe_opsramp_scheduled_maintenance.test_recurring"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_scheduled_maintenance.test_recurring", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_scheduled_maintenance.test_recurring", "name", name),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test_recurring",
							"schedule.type",
							"recurring",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test_recurring",
							"schedule.timezone",
							"Europe/Paris",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test_recurring",
							"schedule.pattern.type",
							"weekly",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test_recurring",
							"schedule.pattern.week_days",
							"Monday,Wednesday,Friday",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test_recurring",
							"schedule.end_by",
							"Never",
						),
					),
				},
			},
		})
	})

	t.Run("with alert conditions", func(t *testing.T) {
		name := acctest.RandomName("sm-alerts")

		startTime := time.Now().Add(96 * time.Hour).Format("2006-01-02T15:04")
		endTime := time.Now().Add(97 * time.Hour).Format("2006-01-02T15:04")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckScheduledMaintenanceDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccScheduledMaintenanceWithAlertConditionsConfig(name, startTime, endTime, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureScheduledMaintenanceExists(t, "hpe_opsramp_scheduled_maintenance.test_alerts"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_scheduled_maintenance.test_alerts", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_scheduled_maintenance.test_alerts", "name", name),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test_alerts",
							"alert_conditions.matching_type",
							"ANY",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test_alerts",
							"alert_conditions.rules.0.key",
							"subject",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test_alerts",
							"alert_conditions.rules.0.operator",
							"contains",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_scheduled_maintenance.test_alerts",
							"alert_conditions.rules.0.value",
							"maintenance",
						),
					),
				},
			},
		})
	})

	t.Run("with device group", func(t *testing.T) {
		name := acctest.RandomName("sm-dg")
		groupName := acctest.RandomName("sm-group")

		startTime := time.Now().Add(120 * time.Hour).Format("2006-01-02T15:04")
		endTime := time.Now().Add(121 * time.Hour).Format("2006-01-02T15:04")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckScheduledMaintenanceDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccScheduledMaintenanceWithDeviceGroupConfig(name, groupName, startTime, endTime, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureScheduledMaintenanceExists(t, "hpe_opsramp_scheduled_maintenance.test_dg"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_scheduled_maintenance.test_dg", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_scheduled_maintenance.test_dg", "name", name),
					),
				},
			},
		})
	})
}

func testAccScheduledMaintenanceOneTimeConfig(name string, description string,
	startTime string, endTime string, clientOverride string,
) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_scheduled_maintenance" "test" {
	name            = "%s"
	description     = "%s"
	correlate_alerts = true
	%s

	run_escalate_action = true

	schedule = {
		type       = "one-time"
		start_time = "%s:00+0200"
		end_time   = "%s:00+0200"
		timezone   = "Europe/Paris"
	}
}
`, acctest.ProviderConfigHCL(), name, description, acctest.ClientAttrHCL(clientOverride), startTime, endTime)
}

func testAccScheduledMaintenanceRecurringWeeklyConfig(name string, startTime string, endTime string, clientOverride string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_scheduled_maintenance" "test_recurring" {
	name            = "%s"
	description     = "Recurring weekly maintenance"
	correlate_alerts = true
	%s

	run_escalate_action = true
	schedule = {
		type       = "recurring"
		start_time = "%s:00+0200"
		end_time   = "%s:00+0200"
		timezone   = "Europe/Paris"
		end_by     = "Never"

		pattern = {
			type      = "weekly"
			week_days = "Monday,Wednesday,Friday"
		}
	}
}
`, acctest.ProviderConfigHCL(), name, acctest.ClientAttrHCL(clientOverride), startTime, endTime)
}

func testAccScheduledMaintenanceWithAlertConditionsConfig(
	name string,
	startTime string,
	endTime string,
	clientOverride string,
) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_scheduled_maintenance" "test_alerts" {
	name            = "%s"
	description     = "Maintenance with alert conditions"
	correlate_alerts = true
	%s

	run_escalate_action = true
	schedule = {
		type       = "one-time"
		start_time = "%s:00+0200"
		end_time   = "%s:00+0200"
		timezone   = "Europe/Paris"
	}

	alert_conditions = {
		matching_type = "ANY"
		rules = [{
			key      = "subject"
			operator = "contains"
			value    = "maintenance"
		}]
	}
}
`, acctest.ProviderConfigHCL(), name, acctest.ClientAttrHCL(clientOverride), startTime, endTime)
}

func testAccScheduledMaintenanceWithDeviceGroupConfig(name string,
	groupName string, startTime string, endTime string, clientOverride string,
) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_device_group" "sm_test_group" {
	name         = "%s"
	search_query = "resourceType = \"Server\""
	%s
}

resource "hpe_opsramp_scheduled_maintenance" "test_dg" {
	name            = "%s"
	description     = "Maintenance with device group"
	correlate_alerts = true
	run_escalate_action = true
	%s

	schedule = {
		type       = "one-time"
		start_time = "%s:00+0200"
		end_time   = "%s:00+0200"
		timezone   = "Europe/Paris"
	}

	device_group_ids = [hpe_opsramp_device_group.sm_test_group.id]
}
`, acctest.ProviderConfigHCL(), groupName, clientAttr, name, clientAttr, startTime, endTime)
}

func testAccEnsureScheduledMaintenanceExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetScheduledMaintenance(tenantID, id)
		if err != nil {
			return fmt.Errorf("scheduled maintenance %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckScheduledMaintenanceDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_scheduled_maintenance" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetScheduledMaintenance(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("scheduled maintenance still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "no schedule maintenance exists with id") {
				return fmt.Errorf("unexpected error checking deleted scheduled maintenance %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}

func testAccScheduledMaintenanceImportStateIdFunc(resourceName string, clientOverride string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		id := rs.Primary.ID
		if clientOverride != "" {
			return clientOverride + ":" + id, nil
		}

		return id, nil
	}
}
