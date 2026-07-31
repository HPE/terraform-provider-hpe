resource "hpe_opsramp_integration_event" "custom_event" {
  integration_id         = hpe_opsramp_integration.custom_bidirectional.id
  name                   = "Custom Event"
  entity                 = "DEFAULT_RESOURCE"
  event_type             = "CREATE"
  resource_group_allowed = false
  use_base_notifier      = true
  endpoint_uri           = "https://external-system.example.com/event-webhook"
  third_party_event_type = "POST"
  event_payload          = "test $default_resource.aliasName"
  headers = {
    "Content-Type" = "application/json"
    "Accept"       = "application/json"
  }
  response_headers = {
    "Status Message" = "$status",
    "extTicketURL"   = "$url",
    "extTicketId"    = "$id",
    "Error Message"  = "$error",
  }
}