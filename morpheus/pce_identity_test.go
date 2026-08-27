// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheus_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Test end to end the auth flow to PCE.
func TestAccMorpheusPCEIdentityAuthFlowOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.PCE)

	fixture := testhelpers.RequirePceIdentityFixture(t)

	dsName := "data.hpe_morpheus_cloud.test"
	cloudName := fixture.CloudName

	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlockPceIdentity()

	// We'll test we can find something using pce identity creds.
	datasourceConfig := `
data "hpe_morpheus_cloud" "test" {
  name = "` + cloudName + `"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + datasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "name", cloudName),
					resource.TestCheckResourceAttrSet(dsName, "id"),
				),
			},
		},
	})
}
