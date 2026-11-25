package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestUnitSHA512CryptTerraform(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{{
			Config: `
terraform {
  required_providers {
    hpe = { source = "hashicorp/hpe" }
  }
}
resource "terraform_data" "password" {
  input = "password"
}
output "hash" {
  value = provider::hpe::morpheus_sha512_crypt(terraform_data.password.output, "saltsalt", 5000)
}
output "sensitivity" {
  value = issensitive(provider::hpe::morpheus_sha512_crypt(sensitive("password"), "saltsalt", 5000))
}
output "explicit_sensitivity" {
  value = issensitive(sensitive(provider::hpe::morpheus_sha512_crypt("password", "saltsalt", 5000)))
}
`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckOutput("hash", "$6$saltsalt$qFmFH.bQmmtXzyBY0s9v7Oicd2z4XSIecDzlB5KiA2/jctKu9YterLp8wwnSq.qc.eoxqOmSuNp2xS0ktL3nh/"),
				resource.TestCheckOutput("sensitivity", "true"),
				resource.TestCheckOutput("explicit_sensitivity", "true"),
			),
		}},
	})
}
