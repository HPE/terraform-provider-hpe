# Import by user ID. The write-only passwords (password_wo, linux_password_wo,
# windows_password_wo) are not stored in state, so the import cannot tell
# whether the values in your configuration match the ones on the appliance:
# they are taken as already applied, the plan after import is empty, and a
# change to them is only detected once the user is next updated for another
# reason. Terraform warns about this on every plan until then.

terraform import hpe_morpheus_user.example 123

# To have changes detected from the start, add the sha256 digest of each
# configured write-only value to the import ID, as <attribute>=<hex>:

terraform import hpe_morpheus_user.example \
  "123#password_wo=$(printf '%s' 'Secret123!' | sha256sum | cut -d' ' -f1)"

# Several values are separated with '&'. The same ID works in an import block:
#
#   import {
#     to = hpe_morpheus_user.example
#     id = "123#password_wo=${sha256("Secret123!")}&linux_password_wo=${sha256("Linux123!")}"
#   }
