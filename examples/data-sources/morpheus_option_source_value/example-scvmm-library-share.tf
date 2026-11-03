# Declaring an SCVMM cloud: resolve the required libraryShare from the VMM
# server itself, through a stored credential, before the cloud exists.

# The SCVMM credentials, stored once in Morpheus. The secret lives here and is
# referenced by id below; it does not pass through the data source, so nothing
# sensitive is recorded in its state.
resource "hpe_morpheus_credential" "scvmm" {
  name     = "scvmm-prod"
  type     = "username-password"
  username = var.scvmm_username
  password = var.scvmm_password
}

# The library share. The source is defined on the SCVMM service, so service_type
# selects it, and it reads the VMM server address from the host parameter. When
# the server has several shares, narrow with name — they are UNC paths — or the
# plan fails listing them.
data "hpe_morpheus_option_source_value" "library_share" {
  source        = "scvmmLibraryShares"
  service_type  = "scvmm"
  credential_id = hpe_morpheus_credential.scvmm.id
  parameters    = { host = var.scvmm_host }
}

resource "hpe_morpheus_cloud" "scvmm" {
  name            = "scvmm-prod"
  tenant_id       = 1
  cloud_type_code = "scvmm"

  config = {
    host         = var.scvmm_host
    username     = var.scvmm_username
    password     = var.scvmm_password
    workingPath  = "c:\\Temp"
    diskPath     = "c:\\VirtualDisks"
    libraryShare = data.hpe_morpheus_option_source_value.library_share.value
  }
}
