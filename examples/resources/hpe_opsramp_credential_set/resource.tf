resource "hpe_opsramp_credential_set" "vmware_credential_set" {
  name        = "VMWare Credential Set"
  description = "Credential set for VMware monitoring"

  credential_type = "VMWARE"
  user_name       = "administrator"
  password        = "**********"
  port            = 443

  secure              = true
  timeout_ms          = 15000
  security_level      = "NOAUTHNOPRIV"
  snmp_version        = "V2"
  ssh_credential_type = "PASSWORD"
  transport_type      = "HTTP"
}