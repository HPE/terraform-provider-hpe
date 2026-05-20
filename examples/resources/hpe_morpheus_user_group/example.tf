resource "hpe_morpheus_user_group" "example" {
  name         = "developers"
  description  = "Development team user group"
  sudo_user    = true
  server_group = "developers"
}
