
terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 1.6.0"
    }
  }
}

provider "hpe" {
  opsramp {
    client_id = "abcdefghijklmnopqrstuvwxyz123456"
    client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
    endpoint      = "tenant.api.pov.opsramp.com"
    tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
  }
}

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

resource "hpe_opsramp_integration_event" "custom_event" {
  integration_id = hpe_opsramp_integration.custom_bidirectional.id
  name = "Custom Event"
  entity = "DEFAULT_RESOURCE"
  event_type = "CREATE"

  resource_group_allowed = false
  use_base_notifier = true

  endpoint_uri = "https://external-system.example.com/event-webhook"
  third_party_event_type = "POST"
  event_payload = "test $default_resource.aliasName"

  headers = {
    "Content-Type" = "application/json"
    "Accept" = "application/json"
  }

  response_headers = {
    "Status Message" = "$status",
    "extTicketURL" = "$url",
    "extTicketId" = "$id",
    "Error Message" = "$error",
  }

}

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
    #process_definition_ids = ["PROCESS_843e00d7-24d0-4bf1-a705-9d955930f09e"]
  }
}


data "hpe_opsramp_custom_event_alert_source" "custom_source" {
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
        #process_definition_ids = ["PROCESS_843e00d7-24d0-4bf1-a705-9d955930f09e"]
    }
}

output "client_event_integration_id" {
    value = hpe_opsramp_integration.client_event_integration.id
}
output "client_event_integration_token" {
    value = hpe_opsramp_integration.client_event_integration.inbound.token
    sensitive = true
}

resource "hpe_opsramp_integration_app" "kubernetes_2_0_integration_app" {
  application  = "Kubernetes-2.0"
  version = "2.3.0"
  bypass_resource_reconciliation = true
}

resource "hpe_opsramp_integration_config" "kubernetes_2_0_integration_config" {
  integration_id = hpe_opsramp_integration_app.kubernetes_2_0_integration_app.id
  
  name = "test"
  config = jsonencode({"Etcd":true,"coreDNS":true,"KubeProxy":true,"enableLog":true,"enableEBPF":true,"kubeEvents":true,"enableTrace":true,"KubeletStats":true,"KubeAPIServer":true,"KubeScheduler":true,"ebpfFeatureFlag":false,"KubeClusterReceiver":true,"KubeControllerManager":true,"clientLevelLogPermission":false,"clientLevelTracePermission":false,"replicaCount":3,"DistributionType":"K8S","eBPFContainerEngine":"cri-o"})
  all_resources = false
}

resource "hpe_opsramp_management_profile" "test_profile" {
  name = "Test Profile"
  description = "Profile for testing integrations"
}

resource "hpe_opsramp_integration" "snmp_integration" {
  application  = "SNMP"
  display_name = "SNMP Integration"
  profile_id = hpe_opsramp_management_profile.test_profile.uuid
}

resource "hpe_opsramp_credential_set" "snmp_credential_set" {
  name        = "SNMP Credential Set"
  description = "Credential set for VMware monitoring"

  credential_type = "SNMP"
  community        = "public"
  port            = 161

  secure              = true
  timeout_ms          = 15000
  security_level      = "NOAUTHNOPRIV"
  snmp_version        = "V2"
  ssh_credential_type = "PASSWORD"
  transport_type      = "HTTP"
}

resource "hpe_opsramp_integration_config" "snmp_integration_config" {
  integration_id = hpe_opsramp_integration.snmp_integration.id
  
  name = "SNMP Config"
  config = jsonencode({"nmapResult":true,"packetCount":"default","deviceType":"SNMP Network Device","networkDepth":"1","discoveryType":"Iprange","ipRange":"10.0.0.1","credentials":[hpe_opsramp_credential_set.snmp_credential_set.id],"cdp":true,"snmpEnabledDevice":true,"ospf":true})
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
  config = jsonencode({"nmapResult":true,"packetCount":"default","deviceType":"SNMP Network Device","networkDepth":"1","discoveryType":"Iprange","ipRange":"10.0.1.1","credentials":[hpe_opsramp_credential_set.snmp_credential_set.id],"cdp":true,"snmpEnabledDevice":true,"ospf":true})
  all_resources = true
}