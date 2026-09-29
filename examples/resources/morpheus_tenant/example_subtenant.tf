# Look up the base role by name using the role data source. The base role must
# be one the parent tenant can assign; a tenant's own base role always is.
data "hpe_morpheus_role" "example" {
  name = "Tenant Admin"
}

# Look up the parent tenant by name using the tenant data source.
data "hpe_morpheus_tenant" "parent" {
  name = "Parent Tenant"
}

# Create a subtenant under that parent. Nominating a parent is only permitted
# for the master tenant, and changing parent_id forces a new resource.
resource "hpe_morpheus_tenant" "subtenant" {
  name         = "Example Subtenant"
  description  = "Terraform example subtenant"
  base_role_id = data.hpe_morpheus_role.example.id
  parent_id    = data.hpe_morpheus_tenant.parent.id
  subdomain    = "subtenant"
  currency     = "USD"
}
