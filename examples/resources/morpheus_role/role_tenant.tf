resource "hpe_morpheus_role" "example" {
  name        = "ExampleTenantRole"
  description = "An example tenant role"
  role_type   = "tenant"
}
