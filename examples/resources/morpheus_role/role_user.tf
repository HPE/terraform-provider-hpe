resource "hpe_morpheus_role" "example" {
  name        = "ExampleUserRole"
  multitenant = false
  description = "An example user role"
  role_type   = "user"
}
