resource "hpe_morpheus_integration_ansible" "tf_example_ansible_integration" {
  name                          = "tfexample ansible integration"
  enabled                       = true
  url                           = "https://github.com/gomorpheus/morpheus-automation-examples.git"
  default_branch                = "main"
  playbooks_path                = "/"
  roles_path                    = "/roles"
  group_variables_path          = "/vars"
  host_variables_path           = "/vars"
  enable_ansible_galaxy_install = true
  enable_verbose_logging        = true
  enable_agent_command_bus      = true
  enable_git_caching            = true
}
