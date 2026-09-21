# Discovering what a source offers: every library share on an SCVMM server,
# through a stored credential — and the same list driving for_each.

resource "hpe_morpheus_credential" "scvmm" {
  name     = "scvmm-prod"
  type     = "username-password"
  username = var.scvmm_username
  password = var.scvmm_password
}

data "hpe_morpheus_option_source_values" "library_shares" {
  source        = "scvmmLibraryShares"
  service_type  = "scvmm"
  credential_id = hpe_morpheus_credential.scvmm.id
  parameters    = { host = var.scvmm_host }
}

output "scvmm_library_shares" {
  value = data.hpe_morpheus_option_source_values.library_shares.values[*].value
}

# One cloud per share.
resource "hpe_morpheus_cloud" "per_share" {
  for_each = { for v in data.hpe_morpheus_option_source_values.library_shares.values : v.value => v }

  name            = "scvmm-${replace(each.key, "/[^A-Za-z0-9]+/", "-")}"
  tenant_id       = 1
  cloud_type_code = "scvmm"

  config = {
    host         = var.scvmm_host
    username     = var.scvmm_username
    password     = var.scvmm_password
    workingPath  = "c:\\Temp"
    diskPath     = "c:\\VirtualDisks"
    libraryShare = each.key
  }
}
