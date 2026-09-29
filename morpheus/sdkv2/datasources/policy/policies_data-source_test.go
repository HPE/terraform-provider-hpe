// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package policy_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dspolicy "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/policy"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusDataSourcePoliciesExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

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

	// create a policy as a dependency purely for testing this. Use a
	// self-contained policy type (Max VMs) so the create needs no external
	// dependency; a workflow policy would require a workflowId.
	dependenciesConfig += `
	resource "hpe_morpheus_policy" "example" {
		name = "` + name + `"
		description = "Example role-scoped policy"
		associated_resource_type = "Role"
		associated_resource_id = resource.hpe_morpheus_role.example.id
		enabled = false
		policy_type = {
			code = "maxVms"
		}
		config_max_vms = {
			max_vms = "5"
		}
	}
	`
	datasourceConfig, err := dspolicy.RenderPoliciesConfig(t, map[string]string{
		"Name":          name,
		"Filter1Values": "[\"Max VMs\"]",
	})
	if err != nil {
		t.Fatal(err)
	}

	// The data source has no attribute reference to the policy, so without an
	// explicit dependency Terraform reads it during plan (before the policy
	// exists) and the filter matches nothing. Force it to be read after apply.
	datasourceConfig = strings.Replace(
		datasourceConfig,
		`data "hpe_morpheus_policies" "example" {`,
		`data "hpe_morpheus_policies" "example" {`+"\n  depends_on = [hpe_morpheus_policy.example]",
		1,
	)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_policies.example",
			"ids.0",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"filter.#",
			"2",
		),

		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_policies.example",
			"filter.*",
			map[string]string{
				"name":     "name",
				"values.#": "1",
			},
		),

		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_policies.example",
			"filter.*",
			map[string]string{
				"name":     "type",
				"values.#": "1",
			},
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_policies.example",
			"sort_ascending",
			"true",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
