# Custom integration with both inbound and outbound
resource "hpe_opsramp_integration" "custom_bidirectional" {
  display_name = "Custom Bidirectional"
  application  = "CUSTOM"
  category     = "Custom"

  inbound = {
    auth_type = "OAUTH2"

    map_attributes = [
      {
        entity_type = "INCIDENT"
        third_party_attribute = "description"
        opsramp_attribute = "incident.impact"
        attribute_values = {"attr1": "value1","attr2": "value2"}
      }
    ]

    enable_drop_alerts = true
  }

  outbound = {
    base_uri  = "https://external-system.example.com/webhook"
    auth_type = "OAUTH2"
    username  = "api-user"
    password  = "secret"
    api_key   = "my-api-key"
    api_secret = "my-api-secret"
    grant_type = "PASSWORD"
    access_token_url = "https://external-system.example.com/oauth/token"
    scope     = "alerts"

    additional_properties = {
      custom_key = "custom_value"
    }

    map_attributes = [
      {
        entity_type = "INCIDENT"
        third_party_attribute = "priority.name"
        opsramp_attribute = "incident.priority.name"
        attribute_values = {"Very Low": "vl","Low": "l"}
      }
    ]
  }
}