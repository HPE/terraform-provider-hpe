
terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 2.1.0"
    }
  }
}

provider "hpe" {
  opsramp {
    client_id     = "abcdefghijklmnopqrstuvwxyz123456"
    client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
    endpoint      = "tenant.api.pov.opsramp.com"
    tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
  }
}

data "hpe_opsramp_metric" "example" {
  metric_name = "haproxy.denied_resp_rate"
  scope       = "GLOBAL"
}

data "hpe_opsramp_adapter_type" "rse" {
  name = "RSE"
}

data "hpe_opsramp_template" "example" {
  name  = "Active Directory Performance Counters DotNet v4"
  scope = "GLOBAL"
}

resource "hpe_opsramp_template" "example-snmp-multiple" {
  name        = "tf-snmp-template-example-multiple"
  description = "Terraform-managed SNMP template multiple."

  template_type = "SNMP"

  #status = "INACTIVE"

  metric_groups = [
    {
      name      = "snmp metric group"
      frequency = 5

      adapter_type = "SNMP"

      # multiple instance
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

      # Extra OIDs collected alongside the metrics below, e.g. for use in alert content.
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

          # Inline OID definitions feeding the expression below.
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

resource "hpe_opsramp_device_management_policy" "example" {
  name = "example-policy"
  # parent_policy = 573963
  search_query = "name CONTAINS \"test\""
  actions = [
    {
      action                   = "ASSIGN MONITORINGTEMPLATE"
      force_assign_or_unassign = true
      items = [
        {
          item       = hpe_opsramp_template.example-snmp-multiple.id
          parameters = {}
        },
        {
          item = data.hpe_opsramp_template.example.id
          parameters = {
            "connection timeout in ms"       = "5000",
            "Remote Server Script Type"      = "bash",
            "Remote Server Operating System" = "linux",
            "Application Type"               = "RSE"
          }
        }
      ]
    }
  ]
}
