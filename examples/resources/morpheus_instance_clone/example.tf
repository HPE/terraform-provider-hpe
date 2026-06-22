resource "hpe_morpheus_instance_clone" "example" {
  source_instance_id = 1
  name               = "my-clone"

  volumes {
    name        = "root"
    size        = 20
    root_volume = true
  }

  network_interfaces {
    network_id = 1
    ip_mode    = ""
  }
}
