// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package policy_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/role"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	dspolicy "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/datasources/policy"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusDataSourcePoliciesExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to API error")
	// t.Skip("Skipping due to missing infrastructure in test environment")
	// t.Skip("Skipping due to missing resource implementation")
	// t.Skip("Skipping due to bug in terraform code")
	// t.Skip("Skipping due to mismatch between Morpheus API and Terraform schema")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	if currentDependency, err := role.RenderRoleUserConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	dependenciesConfig += `
	resource "hpe_morpheus_policy" "example" {
	
	}	
	`

	datasourceConfig, err := dspolicy.RenderPoliciesConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(providerConfig + dependenciesConfig + datasourceConfig)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"filter_name",
			"\"name\"",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"filter_name2",
			"\"type\"",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"filter_values",
			"[\"Test*\"]",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"filter_values2",
			"[\"Max VMs\", \"Workflow\"]",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"sort_ascending",
			"true",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
