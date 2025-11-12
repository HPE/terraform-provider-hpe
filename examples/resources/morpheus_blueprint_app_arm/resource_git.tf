resource "hpe_morpheus_blueprint_app_arm" "tf_example_app_arm_blueprint_git" {
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
