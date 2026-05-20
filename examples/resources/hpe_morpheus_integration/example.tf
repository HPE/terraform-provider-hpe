resource "hpe_morpheus_integration" "example" {
  name     = "Git Repository"
  type     = "git"
  url      = "https://github.com/example/repo.git"
  username = "deploy-user"
  enabled  = true
}
