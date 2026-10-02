resource "hpe_opsramp_template" "example-snmp-simple" {
  name        = "tf-snmp-template-example-simple"
  description = "Terraform-managed SNMP template simple."

  template_type = "SNMP"

  #status = "INACTIVE"

  metric_groups = [
    {
      name      = "snmp metric group"
      frequency = 5

      adapter_type = "SNMP"

      # single instance
      is_multi_instance = false

      metrics = [
        {
          name           = "account.percent.used.capacity"
          display_name   = "Percent used capacity"
          description    = "Percent used capacity"
          category       = "METRIC"
          threshold_type = "STATIC"
          availability   = true

          # Inline OID definitions feeding the expression below.
          metric_values = [
            {
              name           = "hrStorageUsed"
              oid            = "1.3.6.1.2.1.25.2.3.1.6"
              oid_type       = "SCALAROID"
              data_type      = "COUNTER64"
              poll_frequency = "NORMAL"
            },
            {
              name           = "hrStorageSize"
              oid            = "1.3.6.1.2.1.25.2.3.1.5"
              poll_frequency = "NORMAL"
              oid_type       = "SCALAROID"
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