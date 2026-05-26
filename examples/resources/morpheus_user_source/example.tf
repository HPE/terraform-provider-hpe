resource "hpe_morpheus_user_source" "example" {
  name       = "Example User Source"
  type       = "ldap"
  account_id = 1
}
