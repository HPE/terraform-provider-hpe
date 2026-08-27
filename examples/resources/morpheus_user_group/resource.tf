resource "hpe_morpheus_user_group" "example" {
  name         = "tftest"
  description  = "terraform"
  sudo_access  = true
  server_group = "test"
  user_ids     = [19, 10]
}
