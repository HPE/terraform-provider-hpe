resource "hpe_morpheus_vdi_pool" "example" {
  name              = "Developer Desktops"
  description       = "VDI pool for development team"
  max_pool_size     = 10
  min_idle          = 2
  max_idle          = 5
  initial_pool_size = 3
  enabled           = true
  persistent_user   = true
  idle_timeout      = 30
}
