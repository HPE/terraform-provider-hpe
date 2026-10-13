# Import by network domain ID. The write-only value domain_password_wo is not
# stored in state, so the import cannot tell whether the value in your
# configuration matches the one on the appliance: it is taken as already
# applied, the plan after import is empty, and a change to it is only detected
# once the domain is next updated for another reason. Terraform warns about
# this on every plan until then.

terraform import hpe_morpheus_network_domain.example 38785

# To have changes detected from the start, add the sha256 digest of the
# configured write-only value to the import ID, as <attribute>=<hex>:

terraform import hpe_morpheus_network_domain.example \
  "38785#domain_password_wo=$(printf '%s' 'Secret123!' | sha256sum | cut -d' ' -f1)"

# The same ID works in an import block, where Terraform's sha256() function
# produces the digest:
#
#   import {
#     to = hpe_morpheus_network_domain.example
#     id = "38785#domain_password_wo=${sha256("Secret123!")}"
#   }
