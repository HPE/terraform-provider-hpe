resource "hpe_opsramp_site" "site_valencia" {
  parent    = hpe_opsramp_site.vmware_site_spain.id
  name      = "Valencia Data Center"
  address   = "Av. del General Avilés, 35-37, Benicalap"
  country   = "Spain"
  zip       = "46035"
  state     = "Comunitat Valenciana"
  city      = "València"
  resources = [hpe_opsramp_resource.resource1.uuid]
}