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

//nolint:lll
const providerConfigLegacy = `
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

provider "morpheus" {
  url          = var.testacc_morpheus_url
  access_token = var.testacc_morpheus_access_token
  username     = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_username : null
  password     = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_password : null
}
`

//nolint:lll
const providerConfigLegacyProviderBlockOnly = `
provider "morpheus" {
  url          = var.testacc_morpheus_url
  access_token = var.testacc_morpheus_access_token
  username     = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_username : null
  password     = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_password : null
}
`

// Returns a provider block that can be used for acceptance testing
func ProviderBlock() string {
	return providerConfig
}

// Returns a provider block for the legacy morpheus provider that can be used for acceptance testing
func ProviderBlockLegacy() string {
	return providerConfigLegacy
}

// Returns a provider block for mixed usage of the new and old providers in accedptance testing
func ProviderBlockMixed() string {
	return providerConfig + providerConfigLegacyProviderBlockOnly
}
