resource "hpe_morpheus_library_layout" "example" {
  instance_type_id    = 1
  name                = "Single Node"
  instance_version    = "1.0"
  provision_type_code = "vmware"
  description         = "Single node layout for VMware"
  creatable           = true
}
