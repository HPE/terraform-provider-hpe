# Pre-configured integration (e.g. NewRelic) - auth is auto-provisioned on install
resource "hpe_opsramp_integration" "newrelic" {
  application  = "NEWRELIC"

  inbound = {
    auth_type = "WEBHOOK"

    map_attributes = [
        {
            third_party_attribute = "device.ip"
	        opsramp_attribute = "alert.alertTime"
            attribute_values = {
                "attr1": "value1",
                "attr2": "value2"
            }

	        default_parsing_value = "test"
        	parsing_operators = [
                {
                    operator = "BEFORE"
                    end_word = "IP:"
                }
            ]
        }
    ]

    enable_drop_alerts = true
    process_definition_ids = ["PROCESS_843e00d7-24d0-4bf1-a705-9d955930f09e"]
  }
}