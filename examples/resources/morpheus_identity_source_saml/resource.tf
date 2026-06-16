resource "hpe_morpheus_identity_source_saml" "samldemo" {
  tenant_id                      = 1
  name                           = "samldemo"
  description                    = "TF example SAML identity source"
  login_redirect_url             = "https://tfexamplesaml.test.local:8443/realms/master/protocol/saml"
  logout_redirect_url            = "https://tfexamplesaml.test.local:8443/realms/master/protocol/saml"
  include_saml_request_parameter = true
  saml_request                   = "SelfSigned"
  validate_assertion_signature   = false
  given_name_attribute           = "givenName"
  surname_attribute              = "surname"
  email_attribute                = "email"
  default_account_role_id        = 4
  role_attribute_name            = "memberOf"
  required_role_attribute_value  = "test"
  role_mapping {
    role_id             = 4
    role_name           = "Demo"
    assertion_attribute = "developers"
  }

  role_mapping {
    role_id             = 5
    role_name           = "tf-example-user-role"
    assertion_attribute = "developers"
  }
  enable_role_mapping_permission = false
}

# Morpheus computes the SAML Service Provider metadata. Expose it as outputs so
# the SP can be registered in the IdP (e.g. Zitadel) within the same apply.
output "saml_entity_id" {
  value = hpe_morpheus_identity_source_saml.samldemo.entity_id
}

output "saml_acs_url" {
  value = hpe_morpheus_identity_source_saml.samldemo.acs_url
}

output "saml_sp_metadata" {
  value = hpe_morpheus_identity_source_saml.samldemo.sp_metadata
}

