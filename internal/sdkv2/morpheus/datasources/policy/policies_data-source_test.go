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

	t.Skip("Skipping due to bug in terraform code")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	// create a role as a dependency to not affect any existing resources
	if currentDependency, err := role.RenderRoleUserConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	// create a policy as a dependency purely for testing this
	dependenciesConfig += `
	resource "hpe_morpheus_policy" "example" {
		name = "` + name + `"
		description = "Example role-scoped policy"
		associated_resource_type = "Role"
		associated_resource_id = resource.hpe_morpheus_role.example.id
		enabled = false
		policy_type = {
			code = "workflow"
		}
	}
	`
	datasourceConfig, err := dspolicy.RenderPoliciesConfig(t, map[string]string{
		"Name":          name,
		"Filter1Values": "[\"Role\"]",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_policies.example",
			"ids[0]",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"filter.0.name",
			"\"name\"",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"filter.1.name",
			"\"type\"",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"filter.0.values",
			"[\".*\"]",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"filter.1.values",
			"[\"Role\"]",
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
