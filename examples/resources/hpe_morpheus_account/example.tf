resource "hpe_morpheus_account" "example" {
  name        = "Acme Corp"
  description = "Acme Corporation tenant"
  subdomain   = "acme"
  role_id     = 1
  currency    = "USD"
  active      = true
}
