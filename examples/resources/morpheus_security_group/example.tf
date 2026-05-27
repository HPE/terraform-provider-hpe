resource "hpe_morpheus_security_group" "example" {
  name        = "web-servers"
  description = "Security group for web servers"
  active      = true
}
