// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package usergroup_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/usergroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusUserGroupExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to API error")
	// t.Skip("Skipping due to missing infrastructure in test environment")
	// t.Skip("Skipping due to missing resource implementation")
	// t.Skip("Skipping due to mismatch between Morpheus API and Terraform schema")

	providerConfig := testhelpers.ProviderBlock()

	dependenciesConfig := testhelpers.WhoamiBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := usergroup.RenderUserGroupConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user_group.example",
			"description",
			"terraform",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_user_group.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_user_group.example",
			"server_group",
			"test",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_user_group.example",
			"sudo_access",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_user_group.example",
			"user_ids.#",
			"1",
		),

		// MORPH-16495: the fixture must reference the caller's own user id
		// (via the whoami data source), not the hardcoded master user id.
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_user_group.example", "user_ids.0",
			"data.hpe_morpheus_whoami.current", "id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependenciesConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependenciesConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
