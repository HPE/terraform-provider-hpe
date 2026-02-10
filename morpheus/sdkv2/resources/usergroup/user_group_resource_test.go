// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package usergroup_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/usergroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusUserGroupExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to API error")
	// t.Skip("Skipping due to missing infrastructure in test environment")
	// t.Skip("Skipping due to missing resource implementation")
	// t.Skip("Skipping due to mismatch between Morpheus API and Terraform schema")

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

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
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
