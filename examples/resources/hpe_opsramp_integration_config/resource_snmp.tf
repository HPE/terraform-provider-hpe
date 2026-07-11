resource "hpe_opsramp_integration" "snmp_integration" {
  application  = "SNMP"
  display_name = "SNMP Integration"
  profile_id = "7f087cbf-2ce5-4e01-b6e8-a9d9d2f81728"
}

resource "hpe_opsramp_integration_config" "snmp_integration_config" {
  integration_id = hpe_opsramp_integration.snmp_integration.id
  
  name = "SNMP Config"
  config = jsonencode({"nmapResult":true,"packetCount":"default","deviceType":"SNMP Network Device","networkDepth":"1","discoveryType":"Iprange","ipRange":"10.0.0.1","credentials":["NpjSqRafZPMnjbSYHzmFXmCr"],"cdp":true,"snmpEnabledDevice":true,"ospf":true})
  schedule = {
    pattern_type = "DAILY"
    pattern = 3
    start_time = "01"
  }
  all_resources = true
}

resource "hpe_opsramp_integration_config" "snmp_integration_config2" {
  integration_id = hpe_opsramp_integration.snmp_integration.id
  
  name = "SNMP Config 2"
  config = jsonencode({"nmapResult":true,"packetCount":"default","deviceType":"SNMP Network Device","networkDepth":"1","discoveryType":"Iprange","ipRange":"10.0.1.1","credentials":["NpjSqRafZPMnjbSYHzmFXmCr"],"cdp":true,"snmpEnabledDevice":true,"ospf":true})
  all_resources = true
}