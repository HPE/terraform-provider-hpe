resource "hpe_morpheus_app_blueprint_arm" "example" {
  name               = "example_app_arm_blueprint_git"
  description        = "example arm app blueprint"
  category           = "armtemplates"
  source_type        = "repository"
  install_agent      = true
  cloud_init_enabled = true
  os_type            = "linux"
  working_path       = "./test"
  integration_id     = 3
  repository_id      = 1
  version_ref        = "main"
}
