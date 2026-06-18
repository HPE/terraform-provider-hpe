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

// ProviderBlock returns a provider block that can be used for acceptance testing
func ProviderBlock() string {
	return providerConfig
}

// ProviderBlockLegacy returns a provider block for the legacy morpheus provider
func ProviderBlockLegacy() string {
	return providerConfigLegacy
}

// ProviderBlockMixed returns a provider block for mixed usage of the new and old providers
func ProviderBlockMixed() string {
	return providerConfig + providerConfigLegacyProviderBlockOnly
}

//nolint:lll
const providerConfigUnitTest = `
provider "hpe" {
  morpheus {
    url          = "http://localhost"
    access_token = "unit-test"
  }
}
`

// ProviderBlockUnitTest returns a provider block with a placeholder URL and
// token so that plan-time validation tests (CustomizeDiff / schema) can run as
// unit tests (IsUnitTest, no TF_ACC) without real credentials. The Morpheus
// client is created lazily, so no connection is made before the validation
// under test fires.
func ProviderBlockUnitTest() string {
	return providerConfigUnitTest
}
