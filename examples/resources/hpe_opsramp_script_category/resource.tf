resource "hpe_opsramp_script_category" "network" {
  name      = "Network"
  parent_id = hpe_opsramp_script_category.automation.uuid
}