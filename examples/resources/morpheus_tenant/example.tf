# Look up the base role by name using the role data source.
data "hpe_morpheus_role" "example" {
  name = "Tenant Admin"
}

resource "hpe_morpheus_tenant" "example" {
  name            = "Example Tenant"
  description     = "Terraform example tenant"
  enabled         = true
  subdomain       = "tfexample"
  base_role_id    = data.hpe_morpheus_role.example.id
  currency        = "USD"
  account_number  = "12345"
  account_name    = "tenant 12345"
  customer_number = "12345"
}
