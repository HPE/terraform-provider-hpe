resource "hpe_morpheus_app_blueprint_cloud_formation" "tf_example_cloud_formation_app_blueprint_yaml" {
  name                   = "example_cloud_formation_app_blueprint_yaml"
  description            = "Example cloud formation app blueprint"
  category               = "cloudformation"
  install_agent          = true
  cloud_init_enabled     = true
  capability_iam         = true
  capability_named_iam   = true
  capability_auto_expand = true
  source_type            = "yaml"
  blueprint_content      = <<TFEOF
...
TFEOF
}
