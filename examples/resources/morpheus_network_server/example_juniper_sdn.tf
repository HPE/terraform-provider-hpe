# HPE Morpheus SDN (Juniper CN2) Network Server
#
# The SDN controller of an HVM cluster is reached from the cluster's hosts, so
# this type takes no service credentials: only the CN2 controller address and,
# optionally, the introspect and sub-tenant toggles. The controller address
# cannot be changed after creation (changing it replaces the network server).
resource "hpe_morpheus_network_server" "juniper_sdn" {
  name      = "HPE Morpheus SDN"
  type_code = "morpheus-junipersdn-network"
  cloud_id  = 1

  config_juniper_sdn = {
    controller_address = "10.0.0.10"
    show_introspect    = false
    enable_sub_tenants = false
  }
}
