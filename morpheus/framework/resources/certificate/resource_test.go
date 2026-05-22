// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package certificate_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

const fakeCertPEM = `-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJALHHzKLBcg+uMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3RjYTAeFw0yNDAxMDEwMDAwMDBaFw0yNTAxMDEwMDAwMDBaMBExDzANBgNVBAMM
BnRlc3RjYTBcMA0GCSqGSIb3DQEBAQUAAwsAMEgCQQC7o4r6MGKR7ByMMv+2xXNy
QUOTWUQ6QUOTQUOTQUOTQUOTQUOTQUOTQUOT
-----END CERTIFICATE-----`

//nolint:gosec // test fixture
const fakeKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIBogIBAAJBALujivoYYpHsHIwy/7bFc3IZRM8Y0iOtxEE7N2jn7iyQRoW1y1lZ
QUOTQUOTQUOTQUOTQUOTQUOTQUOTQUOTQUOT
-----END RSA PRIVATE KEY-----`

func TestAccMorpheusCertificateBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	createConfig := `
resource "hpe_morpheus_certificate" "test" {
  name      = "` + name + `"
  cert_file = <<-EOT
` + fakeCertPEM + `
EOT
  key_file  = <<-EOT
` + fakeKeyPEM + `
EOT
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_certificate.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_certificate.test", "name", name),
				),
			},
			{
				ResourceName:            "hpe_morpheus_certificate.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cert_file", "key_file"},
			},
		},
	})
}

func TestAccMorpheusCertificateUpdate(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig := `
resource "hpe_morpheus_certificate" "test" {
  name      = "` + name + `"
  cert_file = <<-EOT
` + fakeCertPEM + `
EOT
  key_file  = <<-EOT
` + fakeKeyPEM + `
EOT
}
`

	updateConfig := `
resource "hpe_morpheus_certificate" "test" {
  name      = "` + updatedName + `"
  cert_file = <<-EOT
` + fakeCertPEM + `
EOT
  key_file  = <<-EOT
` + fakeKeyPEM + `
EOT
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_certificate.test", "name", name),
				),
			},
			{
				Config: providerConfig + updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_certificate.test", "name", updatedName),
				),
			},
		},
	})
}
