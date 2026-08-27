resource "hpe_morpheus_vdi_app" "example" {
  name          = "Example"
  description   = "An example description"
  launch_prefix = "||example-launch-prefix"
}
