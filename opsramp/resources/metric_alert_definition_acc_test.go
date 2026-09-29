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

func TestAccMetricAlertDefinitionResource(t *testing.T) {
	acctest.SkipIfNotClient(t)

	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("create_and_update_dynamic_change", func(t *testing.T) {
		name := acctest.RandomName("metric-alert")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckMetricAlertDefinitionDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccMetricAlertDefinitionDynamicChangeConfig(name, "description one", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("hpe_opsramp_metric_alert_definition.test", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "name", name),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_metric_alert_definition.test", "alert_threshold_type", "DYNAMIC_CHANGE_DETECTION"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "entity_type", "RESOURCE"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "component", "$$__name__"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "attributes.0.name", "host"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "attributes.0.value", "$$__name__"),
					),
				},
				{
					Config: testAccMetricAlertDefinitionDynamicChangeConfig(name, "description two", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("hpe_opsramp_metric_alert_definition.test", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "name", name),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "description", "description two"),
					),
				},
			},
		})
	})

	t.Run("create_and_update_static_threshold", func(t *testing.T) {
		name := acctest.RandomName("metric-alert-static")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckMetricAlertDefinitionDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccMetricAlertDefinitionStaticThresholdConfig(name, "50000", "30000", "description one", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("hpe_opsramp_metric_alert_definition.test", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "name", name),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "alert_threshold_type", "STATIC_THRESHOLD"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "no_data_condition", "NO_DATA_ALERT"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "alert_threshold_data.warning_condition", "50000"),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_metric_alert_definition.test", "alert_threshold_data.critical_condition", "30000"),
					),
				},
				{
					Config: testAccMetricAlertDefinitionStaticThresholdConfig(name, "40000", "20000", "description two", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("hpe_opsramp_metric_alert_definition.test", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "name", name),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "description", "description two"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "alert_threshold_data.warning_condition", "40000"),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_metric_alert_definition.test", "alert_threshold_data.critical_condition", "20000"),
					),
				},
			},
		})
	})

	t.Run("create_and_update_dynamic_threshold", func(t *testing.T) {
		name := acctest.RandomName("metric-alert-dyn")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckMetricAlertDefinitionDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccMetricAlertDefinitionDynamicThresholdConfig(name, 2, "WARNING_ALERT", "description one", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("hpe_opsramp_metric_alert_definition.test", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "name", name),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "alert_threshold_type", "DYNAMIC_THRESHOLD"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "alert_threshold_data.limit", "2"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "no_data_condition", "WARNING_ALERT"),
					),
				},
				{
					Config: testAccMetricAlertDefinitionDynamicThresholdConfig(name, 3, "CRITICAL_ALERT", "description two", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("hpe_opsramp_metric_alert_definition.test", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "name", name),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "description", "description two"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "alert_threshold_data.limit", "3"),
						resource.TestCheckResourceAttr("hpe_opsramp_metric_alert_definition.test", "no_data_condition", "CRITICAL_ALERT"),
					),
				},
			},
		})
	})
}

func testAccMetricAlertDefinitionDynamicChangeConfig(name string, description string, clientOverride string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_metric_alert_definition" "test" {
  name       = %q
  %s

  alert_type           = "METRICS"
  query                = "metrics_samples_count"
  alert_threshold_type = "DYNAMIC_CHANGE_DETECTION"

  alert_threshold_data = {
    direction          = "increaseordecrease"
    learning_period    = "4h"
    standard_deviation = 2
  }

  alert_trigger_duration = "5m"

  subject     = "$$__name__ alert for $$resource.name$$ - $$component.name$$ - $$metric.value$$ ($$threshold)"
  description = %q

  entity_type = "RESOURCE"
  component   = "$$__name__"
  status      = true

  labels = [
    {
      name  = "environment"
      value = "test"
    }
  ]

  attributes = [
    {
      name  = "host"
      value = "$$__name__"
    }
  ]
}
`, acctest.ProviderConfigHCL(), name, acctest.ClientAttrHCL(clientOverride), description)
}

func testAccMetricAlertDefinitionStaticThresholdConfig(
	name string,
	warningCondition string,
	criticalCondition string,
	description string,
	clientOverride string,
) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_metric_alert_definition" "test" {
	name       = %q
	%s

	alert_type           = "METRICS"
	query                = "metrics_samples_count"
	alert_threshold_type = "STATIC_THRESHOLD"

	alert_threshold_data = {
		warning_condition  = %q
		critical_condition = %q
	}

	alert_trigger_duration = "1m"
	no_data_condition      = "NO_DATA_ALERT"

	subject     = "$$__name__ alert for $$resource.name$$ - $$component.name$$ - $$metric.value$$ ($$threshold)"
	description = %q

	entity_type = "RESOURCE"
	component   = "$$__name__"
	status      = true

	labels = [
		{
			name  = "environment"
			value = "test"
		}
	]

	attributes = [
		{
			name  = "host"
			value = "$$__name__"
		}
	]
}
`, acctest.ProviderConfigHCL(), name, acctest.ClientAttrHCL(clientOverride), warningCondition, criticalCondition, description)
}

func testAccMetricAlertDefinitionDynamicThresholdConfig(
	name string,
	limit int,
	noDataCondition string,
	description string,
	clientOverride string,
) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_metric_alert_definition" "test" {
	name       = %q
	%s

	alert_type           = "METRICS"
	query                = "metrics_samples_count"
	alert_threshold_type = "DYNAMIC_THRESHOLD"

	alert_threshold_data = {
		limit = %d
	}

	alert_trigger_duration = "1m"
	no_data_condition      = %q

	subject     = "$$__name__ alert for $$resource.name$$ - $$component.name$$ - $$metric.value$$ ($$threshold)"
	description = %q

	entity_type = "RESOURCE"
	component   = "$$__name__"
	status      = true

	labels = [
		{
			name  = "environment"
			value = "test"
		}
	]

	attributes = [
		{
			name  = "host"
			value = "$$__name__"
		}
	]
}
`, acctest.ProviderConfigHCL(), name, acctest.ClientAttrHCL(clientOverride), limit, noDataCondition, description)
}

func testAccCheckMetricAlertDefinitionDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_metric_alert_definition" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			err := apiClient.DeleteMetricAlertDefinition(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("metric alert definition still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "alert definition is not available") &&
				!strings.Contains(errText, "not found") &&
				!strings.Contains(errText, "404") {
				return fmt.Errorf("unexpected error checking deleted metric alert definition %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
