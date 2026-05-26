resource "hpe_morpheus_vdi_app" "example" {
  name          = "Chrome Browser"
  description   = "Google Chrome virtual application"
  launch_prefix = "/usr/bin/google-chrome"
}
