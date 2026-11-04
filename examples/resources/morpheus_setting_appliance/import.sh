# The appliance settings are a singleton: the import ID identifies nothing, and
# any value works. The write-only credentials (smtp_password, proxy_password,
# currency_provider_api_key, twilio_auth_token) are not stored in state, so the
# import cannot tell whether the values in your configuration match the
# appliance's: they are taken as already applied, the plan after import is
# empty, and a change to one is only detected once the settings are next
# updated for another reason. Terraform warns about this on every plan until
# then.

terraform import hpe_morpheus_setting_appliance.example 1

# To have changes detected from the start, add the sha256 digest of each
# configured credential to the import ID, as <attribute>=<hex>, joined by &:

terraform import hpe_morpheus_setting_appliance.example \
  "1#smtp_password=$(printf '%s' 'Password12' | sha256sum | cut -d' ' -f1)&proxy_password=$(printf '%s' 'Password123456' | sha256sum | cut -d' ' -f1)"

# The same ID works in an import block, where Terraform's sha256() function
# produces the digests:
#
#   import {
#     to = hpe_morpheus_setting_appliance.example
#     id = "1#smtp_password=${sha256("Password12")}&proxy_password=${sha256("Password123456")}"
#   }
