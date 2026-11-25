# HPE Morpheus SDN logical router (juniperSdnRouter)
#
# Joins SDN virtual networks on an HVM cluster resource pool. The router type
# is selected by the presence of the config_juniper_sdn block; the SDN network
# server is given through network_integration_id.
resource "hpe_morpheus_network_router" "example" {
  name                   = "TestRouter"
  group_id               = 1
  network_integration_id = 1

  config_juniper_sdn = {
    resource_pool_id     = 1
    virtual_network_uids = ["0be6b79b-757c-4d92-affc-8084e328cb4c"]
  }
}
