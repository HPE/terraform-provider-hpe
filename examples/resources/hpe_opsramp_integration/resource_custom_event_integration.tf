resource "hpe_opsramp_custom_event_alert_source" "custom_source" {
  name      = "Custom"
}

resource "hpe_opsramp_integration" "client_event_integration" {
    display_name = "Custom Integration Example"
    application  = "CUSTOM-EVENT"
    description = "Custom event integration description"
   
    alert_source_id = data.hpe_opsramp_custom_event_alert_source.custom_source.id

    inbound = {

        auth_type = "WEBHOOK"

        map_attributes = [
          {
              third_party_attribute = "alert_time"
              opsramp_attribute = "alert.alertTime"
          },
          {
              third_party_attribute = "alert_id"
              opsramp_attribute = "alert.extAlertId"
              parsing_operators = [
                  {
                      operator = "BEFORE"
                      end_word = " is the alert ID"
                  }
              ]
              default_parsing_value = "default_ext_alert_id"
          },
          {
              third_party_attribute = "alert_component2"
              opsramp_attribute = "alert.component"
              parsing_operators = [
                  {
                      operator = "AFTER"
                      start_word = "Component:"
                  }
              ]
              default_parsing_value = "default_component"
          },
          {
              third_party_attribute = "alert_resource_ip_address"
              opsramp_attribute = "alert.resource.ipAddress"
              parsing_operators = [
                  {
                      operator = "BETWEEN"
                      start_word = "IPs:"
                      end_word = ","
                  }
              ]
              default_parsing_value = "default_ip_address"
          },
          {
              third_party_attribute = "alert_device_name"
              opsramp_attribute = "alert.deviceName"
              parsing_operators = [
                  {
                      operator = "MATCHES"
                      regex_str = "Resource Name: (.+)"
                  }
              ]
              default_parsing_value = "default_device_name"
          }
        ]

        enable_drop_alerts = false
        process_definition_ids = ["PROCESS_843e00d7-24d0-4bf1-a705-9d955930f09e"]
    }
}