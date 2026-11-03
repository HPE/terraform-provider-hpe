resource "hpe_opsramp_template" "example-prometheus" {
  name        = "tf-prometheus-template-example"
  description = "Terraform-managed OTEL-Based Templates-Agent template."

  template_type    = "OTEL Based Templates-Agent"
  application_type = "OTEL_METRICS"

  #status = "INACTIVE"

  metric_groups = [
    {
      name      = "Terraform Otel Metric Group prometheus receiver"
      frequency = 5

      adapter_type     = "APPLICATION"
      application_type = "OTEL_METRICS"

      # script = ""

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
        },
        {
          vistara_custom_param_id = 1497
          name                    = "remote.server.scripttype"
          display_name            = "Remote Server Script Type"
          default_value           = "bash"
          description             = "Select the Remote Server Script Type"
          configurable            = true
          input_type              = "radio"
          input_option_values     = "[{\"Bash\":\"bash\"},{\"PowerShell\":\"powershell\"},{\"Python\":\"python\"},{\"Custom\":\"custom\"}]"
        }
      ]

      otel_metrics_config = {
        config_type            = "prometheus"
        receiver_config_yaml   = <<EOF
        receivers:
          prometheus:
            config:
              scrape_configs:
                - job_name: 'node'
                  static_configs:
                    - targets: ['localhost:9100']
        EOF
        alert_definitions_yaml = ""
        external_exporter      = false
      }

      metrics = [
        {
          name                        = "haproxy.denied_req_rate"
          threshold_type              = "STATIC"
          availability                = true
          alert_on_missing_components = false
          # alert_subject               = ""
          # alert_content               = ""

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
      name      = "Terraform Otel Metric Group nativeOtel receiver"
      frequency = 10

      adapter_type     = "APPLICATION"
      application_type = "OTEL_METRICS"

      # script                = ""

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
        },
        {
          vistara_custom_param_id = 1497
          name                    = "remote.server.scripttype"
          display_name            = "Remote Server Script Type"
          default_value           = "bash"
          description             = "Select the Remote Server Script Type"
          configurable            = true
          input_type              = "radio"
          input_option_values     = "[{\"Bash\":\"bash\"},{\"PowerShell\":\"powershell\"},{\"Python\":\"python\"},{\"Custom\":\"custom\"}]"
      }]

      otel_metrics_config = {
        config_type            = "nativeOtel"
        receiver_config_yaml   = <<EOF
        receivers:
          redis:
            endpoint: "localhost:6379"
            password: ""
            collection_interval: 10s
        EOF
        alert_definitions_yaml = ""
        external_exporter      = false
      }

      metrics = [
        {
          name                        = "haproxy.denied_req_rate"
          threshold_type              = "STATIC"
          availability                = false
          alert_on_missing_components = false
          # alert_subject               = ""
          # alert_content               = ""

          thresholds = [
            {
              raise_alert = false
            }
          ]
        }
      ]
    }
  ]
}