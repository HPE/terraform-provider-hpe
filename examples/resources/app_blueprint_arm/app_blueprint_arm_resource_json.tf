resource "hpe_morpheus_app_blueprint_arm" "tf_example_app_arm_blueprint_json" {
  name               = "example_app_blueprint_json"
  description        = "example arm app blueprint"
  category           = "armtemplates"
  source_type        = "json"
  install_agent      = true
  cloud_init_enabled = true
  os_type            = "linux"
  blueprint_content  = <<EOF
...
EOF
}
