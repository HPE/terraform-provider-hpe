package testhelpers

//nolint:lll
const providerConfig = `
variable "testacc_morpheus_url" {
  default = null
}
variable "testacc_morpheus_username" {
  default = null
}
variable "testacc_morpheus_password" {
  default = null
}
variable "testacc_morpheus_access_token" {
  default = null
}
variable "testacc_morpheus_insecure" {
  default = false
}

provider "hpe" {
        morpheus {
                url = var.testacc_morpheus_url
                access_token    = var.testacc_morpheus_access_token
                username = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_username : null
                password = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_password : null
                insecure = var.testacc_morpheus_insecure
        }
}
`

// Returns a provider block that can be used for acceptance testing
func ProviderBlock() string {
	return providerConfig
}
