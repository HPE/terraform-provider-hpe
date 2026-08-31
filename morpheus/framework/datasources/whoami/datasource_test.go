// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package whoami_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusWhoamiDataSource(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
data "hpe_morpheus_whoami" "current" {}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_whoami.current", "id"),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_whoami.current", "tenant_id"),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_whoami.current", "username"),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_whoami.current", "is_master_account"),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_whoami.current", "permissions.#"),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_whoami.current", "appliance_build_version"),
		resource.TestMatchResourceAttr(
			"data.hpe_morpheus_whoami.current", "tenant_id",
			regexp.MustCompile(`^[1-9][0-9]*$`)),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checkFn,
			},
		},
	})
}
