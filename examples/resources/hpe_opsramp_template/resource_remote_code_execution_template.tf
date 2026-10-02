resource "hpe_opsramp_template" "example" {
  name        = "tf-rse-template-example"
  description = "Terraform-managed G2 Remote Script Executor template."

  template_type    = "Remote Script Executor-Agent"
  application_type = "RSE"

  #status = "INACTIVE"

  metric_groups = [
    {
      name      = "Terraform RSE Metric Group"
      frequency = 5

      adapter_type     = "APPLICATION"
      application_type = "RSE"

      script = "echo '{\"haproxy.denied_resp_rate\":{\"components\":{\"default\":1}}}'"

      configuration_parameters = [
        {
          vistara_custom_param_id = 738
          name                    = "collector.application.type"
          display_name            = "Application Type"
          default_value           = "RSE"
          description             = "Application Type"
          configurable            = true
          input_type              = "text"
        },
        {
          vistara_custom_param_id = 752
          name                    = "connection.timeout.ms"
          display_name            = "connection timeout in ms"
          default_value           = "15000"
          description             = "connection timeout in ms"
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

      metrics = [
        {
          name                        = "haproxy.denied_resp_rate"
          threshold_type              = "STATIC"
          availability                = true
          alert_on_missing_components = true
          # alert_subject               = ""
          # alert_content               = ""

          thresholds = [
            {
              critical_repeat_count = 1
              warning_repeat_count  = 3
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
    }
  ]
}