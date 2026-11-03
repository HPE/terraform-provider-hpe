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

func TestAccTemplateResourceRSE(t *testing.T) {
	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("create_and_update", func(t *testing.T) {
		name := acctest.RandomName("tf-rse-template")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckTemplateDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccTemplateRSEConfig(name, "Terraform-managed G2 Remote Script Executor template.", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureTemplateExists(t, "hpe_opsramp_template.test_rse"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_template.test_rse", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "name", name),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "template_type", "Remote Script Executor-Agent"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "application_type", "RSE"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "metric_groups.#", "2"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "metric_groups.0.name", "Terraform RSE Metric Group"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "metric_groups.0.frequency", "5"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "metric_groups.1.name", "Terraform RSE Metric Group 2"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "metric_groups.1.frequency", "2"),
					),
				},
				{
					Config: testAccTemplateRSEConfig(name, "Terraform-managed G2 Remote Script Executor template, updated.", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureTemplateExists(t, "hpe_opsramp_template.test_rse"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "name", name),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_rse", "description", "Terraform-managed G2 Remote Script Executor template, updated.",
						),
						// The API returns metric groups in a non-deterministic order; re-applying
						// the same metric group data must not report the two groups as swapped.
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "metric_groups.0.name", "Terraform RSE Metric Group"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_rse", "metric_groups.1.name", "Terraform RSE Metric Group 2"),
					),
				},
			},
		})
	})
}

func TestAccTemplateResourcePrometheus(t *testing.T) {
	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("create_and_update", func(t *testing.T) {
		name := acctest.RandomName("tf-prometheus-template")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckTemplateDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccTemplatePrometheusConfig(name, "Terraform-managed OTEL-Based Templates-Agent template.", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureTemplateExists(t, "hpe_opsramp_template.test_prometheus"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_template.test_prometheus", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_prometheus", "name", name),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_prometheus", "template_type", "OTEL Based Templates-Agent",
						),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_prometheus", "application_type", "OTEL_METRICS"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_prometheus", "metric_groups.#", "2"),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_prometheus", "metric_groups.0.name", "Terraform Otel Metric Group",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_prometheus", "metric_groups.0.otel_metrics_config.config_type", "nativeOtel",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_prometheus", "metric_groups.1.name", "Terraform Otel Metric Group 2",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_prometheus", "metric_groups.1.otel_metrics_config.config_type", "prometheus",
						),
					),
				},
				{
					Config: testAccTemplatePrometheusConfig(
						name, "Terraform-managed OTEL-Based Templates-Agent template, updated.", clientOverride,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureTemplateExists(t, "hpe_opsramp_template.test_prometheus"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_prometheus", "name", name),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_prometheus", "description",
							"Terraform-managed OTEL-Based Templates-Agent template, updated.",
						),
						// The API returns metric groups in a non-deterministic order; re-applying
						// the same metric group data must not report the two groups as swapped.
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_prometheus", "metric_groups.0.name", "Terraform Otel Metric Group",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_prometheus", "metric_groups.1.name", "Terraform Otel Metric Group 2",
						),
					),
				},
			},
		})
	})
}

func testAccTemplateRSEConfig(name string, description string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
data "hpe_opsramp_metric" "test_rse" {
  metric_name = "haproxy.denied_resp_rate"
  scope       = "GLOBAL"
}

data "hpe_opsramp_adapter_type" "test_rse" {
  name = "RSE"
}

resource "hpe_opsramp_template" "test_rse" {
  name        = "%s"
  description = "%s"
  %s

  template_type    = "Remote Script Executor-Agent"
  application_type = "RSE"

  metric_groups = [
    {
      name      = "Terraform RSE Metric Group"
      frequency = 5

      adapter_type     = "APPLICATION"
      application_type = "RSE"

      script = format("sleep 1; echo '{\"%%s\":{\"components\":{\"default\":1}}}'", data.hpe_opsramp_metric.test_rse.metric_name)

      configuration_parameters = [
        {
          vistara_custom_param_id = 738
          name                    = "collector.application.type"
          display_name            = "Application Type"
          default_value           = data.hpe_opsramp_adapter_type.test_rse.name
          description             = "Application Type"
          configurable            = true
          input_type              = "text"
        },
        {
          vistara_custom_param_id = 1492
          name                    = "remote.server.operatingsystem"
          display_name            = "Remote Server Operating System"
          default_value           = "linux"
          description             = "Select Remote Server Operating System either windows or Linux Based"
          configurable            = true
          input_type              = "radio"
          input_option_values     = "[{\"Windows\":\"windows\"}, {\"Linux or Unix\":\"linux\"}]"
        }
      ]

      metrics = [
        {
          name                        = data.hpe_opsramp_metric.test_rse.metric_name
          threshold_type              = "STATIC"
          availability                = true
          alert_on_missing_components = true
          alert_subject               = "test"
          alert_content               = "test"

          thresholds = [
            {
              critical_repeat_count = 1
              warning_repeat_count  = 1
              critical_operator     = "GREATER_THAN_EQUAL"
              warning_operator      = "GREATER_THAN"
              raise_alert           = true
              level_type            = "TEMPLATE"
              warning_threshold     = "10"
              critical_threshold    = "15"
            }
          ]
        }
      ]
    },
    {
      name      = "Terraform RSE Metric Group 2"
      frequency = 2

      adapter_type     = "APPLICATION"
      application_type = "RSE"

      script = format("sleep 1; echo '{\"%%s\":{\"components\":{\"default\":1}}}'", data.hpe_opsramp_metric.test_rse.metric_name)

      configuration_parameters = [
        {
          vistara_custom_param_id = 738
          name                    = "collector.application.type"
          display_name            = "Application Type"
          default_value           = data.hpe_opsramp_adapter_type.test_rse.name
          description             = "Application Type"
          configurable            = true
          input_type              = "text"
        }
      ]

      metrics = [
        {
          name                        = "haproxy.denied_req_rate"
          threshold_type              = "STATIC"
          availability                = true
          alert_on_missing_components = true
          alert_subject               = "test 2"
          alert_content               = "test 2"

          thresholds = [
            {
              critical_repeat_count = 2
              warning_repeat_count  = 1
              critical_operator     = "GREATER_THAN_EQUAL"
              warning_operator      = "GREATER_THAN"
              raise_alert           = true
              level_type            = "TEMPLATE"
              warning_threshold     = "12"
              critical_threshold    = "13"
            }
          ]
        }
      ]
    }
  ]
}
`, acctest.ProviderConfigHCL(), name, description, clientAttr)
}

func testAccTemplatePrometheusConfig(name string, description string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
data "hpe_opsramp_metric" "test_prometheus" {
  metric_name = "haproxy.denied_resp_rate"
  scope       = "GLOBAL"
}

resource "hpe_opsramp_template" "test_prometheus" {
  name        = "%s"
  description = "%s"
  %s

  template_type    = "OTEL Based Templates-Agent"
  application_type = "OTEL_METRICS"

  metric_groups = [
    {
      name      = "Terraform Otel Metric Group"
      frequency = 5

      adapter_type     = "APPLICATION"
      application_type = "OTEL_METRICS"

      script = ""

      configuration_parameters = [
        {
          vistara_custom_param_id = 1492
          name                    = "remote.server.operatingsystem"
          display_name            = "Remote Server Operating System"
          default_value           = "linux"
          description             = "Select Remote Server Operating System either windows or Linux Based"
          configurable            = true
          input_type              = "radio"
          input_option_values     = "[{\"Windows\":\"windows\"}, {\"Linux or Unix\":\"linux\"}]"
        }
      ]

      otel_metrics_config = {
        config_type            = "nativeOtel"
        receiver_config_yaml   = "options: true"
        alert_definitions_yaml = "options: true"
        external_exporter      = false
      }

      metrics = [
        {
          name                        = data.hpe_opsramp_metric.test_prometheus.metric_name
          threshold_type              = "STATIC"
          availability                = true
          alert_on_missing_components = true
          alert_subject               = "test"
          alert_content               = "test"

          thresholds = [
            {
              critical_repeat_count = 1
              warning_repeat_count  = 1
              critical_operator     = "GREATER_THAN_EQUAL"
              warning_operator      = "GREATER_THAN_EQUAL"
              raise_alert           = true
              level_type            = "TEMPLATE"
              warning_threshold     = "10"
              critical_threshold    = "15"
            }
          ]
        }
      ]
    },
    {
      name      = "Terraform Otel Metric Group 2"
      frequency = 2

      adapter_type     = "APPLICATION"
      application_type = "OTEL_METRICS"

      configuration_parameters = [
        {
          vistara_custom_param_id = 1492
          name                    = "remote.server.operatingsystem"
          display_name            = "Remote Server Operating System"
          default_value           = "linux"
          description             = "Select Remote Server Operating System either windows or Linux Based"
          configurable            = true
          input_type              = "radio"
          input_option_values     = "[{\"Windows\":\"windows\"}, {\"Linux or Unix\":\"linux\"}]"
        }
      ]

      otel_metrics_config = {
        config_type            = "prometheus"
        receiver_config_yaml   = "options: true"
        alert_definitions_yaml = "options: true"
        external_exporter      = false
      }

      metrics = [
        {
          name                        = "haproxy.denied_req_rate"
          threshold_type              = "STATIC"
          availability                = true
          alert_on_missing_components = true
          alert_subject               = "test 3"
          alert_content               = "test 3"

          thresholds = [
            {
              critical_repeat_count = 1
              warning_repeat_count  = 1
              critical_operator     = "GREATER_THAN_EQUAL"
              warning_operator      = "GREATER_THAN"
              raise_alert           = true
              level_type            = "TEMPLATE"
              warning_threshold     = "12"
              critical_threshold    = "14"
            }
          ]
        }
      ]
    }
  ]
}
`, acctest.ProviderConfigHCL(), name, description, clientAttr)
}

func TestAccTemplateResourceSNMPSimple(t *testing.T) {
	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("create_and_update", func(t *testing.T) {
		name := acctest.RandomName("tf-snmp-template-simple")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckTemplateDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccTemplateSNMPSimpleConfig(name, "Terraform-managed SNMP template simple.", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureTemplateExists(t, "hpe_opsramp_template.test_snmp_simple"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_template.test_snmp_simple", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_snmp_simple", "name", name),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_snmp_simple", "template_type", "SNMP"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_snmp_simple", "metric_groups.#", "1"),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_simple", "metric_groups.0.name", "snmp metric group",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_simple", "metric_groups.0.is_multi_instance", "false",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_simple", "metric_groups.0.metrics.0.name", "account.percent.used.capacity",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_simple", "metric_groups.0.metrics.0.category", "METRIC",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_simple", "metric_groups.0.metrics.0.metric_values.#", "2",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_simple", "metric_groups.0.metrics.0.expression",
							"(hrStorageUsed / hrStorageSize) * 100",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_simple",
							"metric_groups.0.metrics.0.format_metric_value_details.units", "MBytes",
						),
					),
				},
				{
					Config: testAccTemplateSNMPSimpleConfig(name, "Terraform-managed SNMP template simple, updated.", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureTemplateExists(t, "hpe_opsramp_template.test_snmp_simple"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_snmp_simple", "name", name),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_simple", "description",
							"Terraform-managed SNMP template simple, updated.",
						),
					),
				},
			},
		})
	})
}

func TestAccTemplateResourceSNMPMultiple(t *testing.T) {
	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("create_and_update", func(t *testing.T) {
		name := acctest.RandomName("tf-snmp-template-multiple")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckTemplateDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccTemplateSNMPMultipleConfig(name, "Terraform-managed SNMP template multiple.", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureTemplateExists(t, "hpe_opsramp_template.test_snmp_multiple"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_template.test_snmp_multiple", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_snmp_multiple", "name", name),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_snmp_multiple", "template_type", "SNMP"),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_multiple", "metric_groups.0.is_multi_instance", "true",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_multiple", "metric_groups.0.instance_filter.instance_type", "ALL",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_multiple", "metric_groups.0.component_properties.name", "Name",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_multiple", "metric_groups.0.component_properties.mode", "INDEX_LOOKUP",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_multiple", "metric_groups.0.additional_properties.#", "2",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_multiple", "metric_groups.0.additional_properties.0.name", "ifAdminStatus",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_multiple", "metric_groups.0.additional_properties.1.name", "ifOperStatus",
						),
					),
				},
				{
					Config: testAccTemplateSNMPMultipleConfig(name, "Terraform-managed SNMP template multiple, updated.", clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureTemplateExists(t, "hpe_opsramp_template.test_snmp_multiple"),
						resource.TestCheckResourceAttr("hpe_opsramp_template.test_snmp_multiple", "name", name),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_multiple", "description",
							"Terraform-managed SNMP template multiple, updated.",
						),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_template.test_snmp_multiple", "metric_groups.0.additional_properties.#", "2",
						),
					),
				},
			},
		})
	})
}

func testAccTemplateSNMPSimpleConfig(name string, description string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_template" "test_snmp_simple" {
  name        = "%s"
  description = "%s"
  %s

  template_type = "SNMP"

  metric_groups = [
    {
      name      = "snmp metric group"
      frequency = 5

      adapter_type = "SNMP"

      is_multi_instance = false

      metrics = [
        {
          name           = "account.percent.used.capacity"
          display_name   = "Percent used capacity"
          description    = "Percent used capacity"
          category       = "METRIC"
          threshold_type = "STATIC"
          availability   = true

          metric_values = [
            {
              name      = "hrStorageUsed"
              oid       = "1.3.6.1.2.1.25.2.3.1.6"
              oid_type  = "SCALAROID"
              data_type = "COUNTER64"
            },
            {
              name = "hrStorageSize"
              oid  = "1.3.6.1.2.1.25.2.3.1.5"
              oid_type  = "SCALAROID"
            }
          ]
          expression = "(hrStorageUsed / hrStorageSize) * 100"

          format_metric_value_details = {
            units = "MBytes"
            graph = true
            format_metric_value = {
              used_by = "NONE"
              alerts  = false
              graph   = false
            }
          }

          thresholds = [
            {
              level_type            = "TEMPLATE"
              raise_alert           = true
              warning_operator      = "NONE"
              critical_operator     = "NONE"
              warning_repeat_count  = 1
              critical_repeat_count = 1
            }
          ]
        }
      ]
    }
  ]
}
`, acctest.ProviderConfigHCL(), name, description, clientAttr)
}

func testAccTemplateSNMPMultipleConfig(name string, description string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_template" "test_snmp_multiple" {
  name        = "%s"
  description = "%s"
  %s

  template_type = "SNMP"

  metric_groups = [
    {
      name      = "snmp metric group"
      frequency = 5

      adapter_type = "SNMP"

      is_multi_instance = true

      instance_filter = {
        instance_type = "ALL"
      }

      component_properties = {
        name       = "Name"
        mode       = "INDEX_LOOKUP"
        oid        = "1.3.6.1.4.1.9.9.109.1.1.1.1.3"
        lookup_oid = "1.3.6.1.2.1.47.1.1.1.1.4"
        format     = "COUNTER64"
      }

      additional_properties = [
        {
          name           = "ifAdminStatus"
          oid_type       = "SCALAROID"
          oid            = "1.3.6.1.2.1.2.2.1.7"
          enumerated_map = []
          format_value   = "NONE"
        },
        {
          name           = "ifOperStatus"
          oid_type       = "SCALAROID"
          oid            = "1.3.6.1.2.1.2.2.1.8"
          enumerated_map = []
          format_value   = "NONE"
        }
      ]

      metrics = [
        {
          name           = "account.percent.used.capacity"
          display_name   = "Percent used capacity"
          description    = "Percent used capacity"
          category       = "METRIC"
          threshold_type = "STATIC"
          availability   = true

          metric_values = [
            {
              name      = "hrStorageUsed"
              oid       = "1.3.6.1.2.1.25.2.3.1.6"
              oid_type  = "SCALAROID"
              data_type = "COUNTER64"
            },
            {
              name     = "hrStorageSize"
              oid      = "1.3.6.1.2.1.25.2.3.1.5"
              oid_type = "SCALAROID"
            }
          ]
          expression = "(hrStorageUsed / hrStorageSize) * 100"

          format_metric_value_details = {
            units = "MBytes"
            graph = true
            format_metric_value = {
              used_by = "NONE"
              alerts  = false
              graph   = false
            }
          }

          thresholds = [
            {
              level_type            = "TEMPLATE"
              raise_alert           = true
              warning_operator      = "NONE"
              critical_operator     = "NONE"
              warning_repeat_count  = 1
              critical_repeat_count = 1
            }
          ]
        }
      ]
    }
  ]
}
`, acctest.ProviderConfigHCL(), name, description, clientAttr)
}

func testAccEnsureTemplateExists(t *testing.T, resourceName string) resource.TestCheckFunc {
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

		_, err = apiClient.GetTemplate(tenantID, id)
		if err != nil {
			return fmt.Errorf("template %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckTemplateDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_opsramp_template" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			id := strings.TrimSpace(rs.Primary.ID)
			if id == "" {
				continue
			}

			_, err := apiClient.GetTemplate(tenantID, id)
			if err == nil {
				return fmt.Errorf("template still exists: %s", id)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "not found") && !strings.Contains(errText, "404") {
				return fmt.Errorf("unexpected error checking deleted template %s: %w", id, err)
			}
		}

		return nil
	}
}
