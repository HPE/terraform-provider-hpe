resource "hpe_opsramp_site" "vmware_site_madrid" {
  parent       = hpe_opsramp_site.vmware_site_spain.id
  name         = "CTC Madrid"
  address      = "Calle Vicente Aleixandre, 1"
  country      = "Spain"
  zip          = "28232"
  state        = "Madrid"
  city         = "Las Rozas de Madrid"
  search_query = "name CONTAINS \"MAD\""
}