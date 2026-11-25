# Allocate a floating IP from a pool on an SDN network server.
#
# The floating IP pool and the network server are not returned by the API, so
# they are kept in Terraform state; changing either replaces the resource.
# Floating IP pools are defined on the SDN router of a cluster resource pool
# and are discoverable through the option source, filtered by that pool; the
# option source returns ids as strings.
data "hpe_morpheus_network_server" "sdn" {
  name = "HPE Morpheus SDN"
}

data "hpe_morpheus_option_source_values" "floating_ip_pools" {
  source     = "juniperFloatingIpPools"
  parameters = { juniperResourcePool = "1" }
}

resource "hpe_morpheus_network_floating_ip" "example" {
  network_server_id   = data.hpe_morpheus_network_server.sdn.id
  floating_ip_pool_id = tonumber(data.hpe_morpheus_option_source_values.floating_ip_pools.values[0].value)
}
