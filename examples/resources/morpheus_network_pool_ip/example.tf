# Auto-allocate: omit ip_address and Morpheus reserves the pool's next free
# address, returned as the computed ip_address.
resource "hpe_morpheus_network_pool_ip" "auto" {
  network_pool_id = 1
  hostname        = "host01"
}

# Pinned address: the given ip_address must lie within the pool's IP ranges.
resource "hpe_morpheus_network_pool_ip" "pinned" {
  network_pool_id = 1
  hostname        = "host01-pinned"
  ip_address      = "10.0.0.50"
}

# Attached to a server: associate the reservation with a resource by type/id.
resource "hpe_morpheus_network_pool_ip" "attached" {
  network_pool_id = 1
  hostname        = "host01-attached"
  ref_type        = "ComputeServer"
  ref_id          = 1
}
