package testhelpers

import (
	"bytes"
	"text/template"
)

const providerVariables = `
variable "testacc_morpheus_{{.}}_url" {
  default = null
}
variable "testacc_morpheus_{{.}}_username" {
  default = null
}
variable "testacc_morpheus_{{.}}_password" {
  default = null
}
variable "testacc_morpheus_{{.}}_access_token" {
  default = null
}
variable "testacc_morpheus_{{.}}_insecure" {
  default = false
}
`

//nolint:lll
const providerConfig = providerVariables + `

provider "hpe" {
    morpheus {
        url = var.testacc_morpheus_{{.}}_url
        access_token    = var.testacc_morpheus_{{.}}_access_token
        username = var.testacc_morpheus_{{.}}_access_token == null ? var.testacc_morpheus_{{.}}_username : null
        password = var.testacc_morpheus_{{.}}_access_token == null ? var.testacc_morpheus_{{.}}_password : null
        insecure = var.testacc_morpheus_{{.}}_insecure
    }
}
`

//nolint:lll
const providerConfigLegacy = providerVariables + `

provider "morpheus" {
  url          = var.testacc_morpheus_{{.}}_url
  access_token = var.testacc_morpheus_{{.}}_access_token
  username     = var.testacc_morpheus_{{.}}_access_token == null ? var.testacc_morpheus_{{.}}_username : null
  password     = var.testacc_morpheus_{{.}}_access_token == null ? var.testacc_morpheus_{{.}}_password : null
}
`

//nolint:lll
const providerConfigLegacyProviderBlockOnly = `
provider "morpheus" {
  url          = var.testacc_morpheus_{{.}}_url
  access_token = var.testacc_morpheus_{{.}}_access_token
  username     = var.testacc_morpheus_{{.}}_access_token == null ? var.testacc_morpheus_{{.}}_username : null
  password     = var.testacc_morpheus_{{.}}_access_token == null ? var.testacc_morpheus_{{.}}_password : null
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

// Returns a provider block for mixed usage of the new and old providers in acceptance testing
func ProviderBlockMixed() string {
	return providerConfig + providerConfigLegacyProviderBlockOnly
}

// Returns a provider block for tests, using the preferred system as a parameter for configuration
func ProviderBlockForServer(preferredSystem string) string {
	tmpl, err := template.New("provider-block").Parse(providerConfig)
	if err != nil {
		panic("could not parse template" + err.Error())
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, preferredSystem); err != nil {
		panic("could not execute template" + err.Error())
	}

	return out.String()
}

// Returns a provider block for legacy provider tests, using the preferred system for configuration
func ProviderBlockLegacyForServer(preferredSystem string) string {
	tmpl, err := template.New("provider-block").Parse(providerConfigLegacy)
	if err != nil {
		panic("could not parse template" + err.Error())
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, preferredSystem); err != nil {
		panic("could not execute template" + err.Error())
	}

	return out.String()
}

// Returns a provider block for mixed new and old providers in tests, using the preferred system as a parameter
func ProviderBlockMixedForServer(preferredSystem string) string {
	tmpl, err := template.New("provider-block").Parse(providerConfig + providerConfigLegacyProviderBlockOnly)
	if err != nil {
		panic("could not parse template" + err.Error())
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, preferredSystem); err != nil {
		panic("could not execute template" + err.Error())
	}

	return out.String()
}
