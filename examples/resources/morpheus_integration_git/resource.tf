resource "hpe_morpheus_integration_git" "tf_example_git_integration" {
  name               = "tfexample git integration"
  enabled            = true
  url                = "https://github.com/gomorpheus/morpheus-automation-examples.git"
  default_branch     = "main"
  enable_git_caching = true
}
