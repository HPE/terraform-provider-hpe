package testhelpers_test

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestBuildProviderBlockWithAccessToken(t *testing.T) {
	t.Setenv("TF_VAR_testacc_morpheus_url", "https://test.morpheus.com")
	t.Setenv("TF_VAR_testacc_morpheus_access_token", "abcdefg")
	t.Setenv("TF_VAR_testacc_morpheus_username", "")
	t.Setenv("TF_VAR_testacc_morpheus_password", "")

	providerConfig, err := testhelpers.BuildProviderBlock()
	assert.NoError(t, err)

	expected := `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_insecure" {
  default = false
}
variable "testacc_morpheus_access_token" {}

provider "hpe" {
  morpheus {
    url          = var.testacc_morpheus_url
    insecure     = var.testacc_morpheus_insecure
    access_token = var.testacc_morpheus_access_token
  }
}
`
	assert.Equal(t, expected, providerConfig)
}

func TestBuildProviderBlockWithCredentials(t *testing.T) {
	t.Setenv("TF_VAR_testacc_morpheus_url", "https://test.morpheus.com")
	t.Setenv("TF_VAR_testacc_morpheus_username", "foo@test.com")
	t.Setenv("TF_VAR_testacc_morpheus_password", "testpass")
	t.Setenv("TF_VAR_testacc_morpheus_access_token", "")

	providerConfig, err := testhelpers.BuildProviderBlock()
	assert.NoError(t, err)

	expected := `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_insecure" {
  default = false
}
variable "testacc_morpheus_username" {}
variable "testacc_morpheus_password" {}

provider "hpe" {
  morpheus {
    url          = var.testacc_morpheus_url
    insecure     = var.testacc_morpheus_insecure
    username     = var.testacc_morpheus_username
    password     = var.testacc_morpheus_password
  }
}
`

	assert.Equal(t, expected, providerConfig)
}

func TestBuildProviderBlockNoneSet(t *testing.T) {
	// can't unset these for t without a cleanup func,
	// but setting to "" is good enough
	t.Setenv("TF_VAR_testacc_morpheus_url", "")
	t.Setenv("TF_VAR_testacc_morpheus_username", "")
	t.Setenv("TF_VAR_testacc_morpheus_password", "")
	t.Setenv("TF_VAR_testacc_morpheus_access_token", "")

	providerConfig, err := testhelpers.BuildProviderBlock()

	assert.Equal(t, "", providerConfig)

	//nolint:lll
	assert.EqualError(t, err, "One or more environment variables were not set: TF_VAR_testacc_morpheus_url, TF_VAR_testacc_morpheus_username, TF_VAR_testacc_morpheus_password, TF_VAR_testacc_morpheus_access_token")
}
