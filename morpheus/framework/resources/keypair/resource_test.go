// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package keypair_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusKeyPairBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_key_pair" "test" {
  name       = "` + name + `"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJ+5Di5CbCrDWUBicVSsWM5baTLxTXA88ZF+EzbZuUg4 test@example.com"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_key_pair.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_key_pair.test", "name", name),
				),
			},
			{
				ResourceName:            "hpe_morpheus_key_pair.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key", "passphrase"},
			},
		},
	})
}
