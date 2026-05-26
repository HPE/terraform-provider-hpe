resource "hpe_morpheus_user_source" "example" {
  name                    = "Corporate LDAP"
  type                    = "ldap"
  account_id              = 1
  description             = "Corporate Active Directory integration"
  default_account_role_id = 1
}
