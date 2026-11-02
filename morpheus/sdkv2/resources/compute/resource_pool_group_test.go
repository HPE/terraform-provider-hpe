// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package compute_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/compute"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusResourcePoolGroupExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.ResourcePool)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	dependenciesConfig := testhelpers.WhoamiBlock()
	visibility := testhelpers.TenantVisibility(t)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := compute.RenderResourcePoolGroupConfig(t, map[string]string{
		"Name":       name,
		"Visibility": visibility,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"description",
			"TFExample Resource Pool Group",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"mode",
			"roundrobin",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"resource_pool_ids.0",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"resource_pool_ids.1",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"resource_pool_ids.2",
			"3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"all_group_access",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"group_access.0.group_id",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"group_access.0.default",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"visibility",
			visibility,
		),
		// MORPH-8207 regression: tenants must round-trip (apply then plan-empty).
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_resource_pool_group.example", "tenant_ids.0",
			"data.hpe_morpheus_whoami.current", "tenant_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"tenant_ids.1",
			"2",
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

// resourcePoolGroupTenantConfig renders a minimal resource pool group with
// tenant_ids rendered verbatim from tenantIDs, or omitted when it is empty.
func resourcePoolGroupTenantConfig(name, tenantIDs string) string {
	tenantLine := ""
	if tenantIDs != "" {
		tenantLine = "\n  tenant_ids        = " + tenantIDs
	}

	return `
resource "hpe_morpheus_resource_pool_group" "test" {
  name              = "` + name + `"
  mode              = "roundrobin"
  resource_pool_ids = [1]` + tenantLine + `
}
`
}

// TestAccMorpheusResourcePoolGroupTenantIdsOk reproduces MORPH-8207: the API
// always reports the owning tenant in tenants, so a config that does not set
// tenant_ids used to plan "- <owner>" forever after apply. The owner must be
// absent from state unless listed explicitly, in which case it is kept, and
// import must agree with the applied state.
func TestAccMorpheusResourcePoolGroupTenantIdsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.ResourcePool)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	dependenciesConfig := testhelpers.WhoamiBlock()
	name := acctest.RandomWithPrefix(t.Name())
	const addr = "hpe_morpheus_resource_pool_group.test"

	unsetConfig := providerConfig + dependenciesConfig + resourcePoolGroupTenantConfig(name, "")
	ownerConfig := providerConfig + dependenciesConfig +
		resourcePoolGroupTenantConfig(name, "["+testhelpers.WhoamiTenantIDRef+"]")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// tenant_ids unset: the injected owner must not reach state.
			{
				Config: unsetConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "name", name),
					resource.TestCheckResourceAttr(addr, "mode", "roundrobin"),
					resource.TestCheckResourceAttr(addr, "tenant_ids.#", "0"),
				),
			},
			{
				Config:   unsetConfig,
				PlanOnly: true,
			},
			{
				ResourceName:      addr,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Listing the owning tenant explicitly keeps it in state.
			{
				Config: ownerConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "tenant_ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						addr, "tenant_ids.0",
						"data.hpe_morpheus_whoami.current", "tenant_id",
					),
				),
			},
			{
				Config:   ownerConfig,
				PlanOnly: true,
			},
			// Unsetting it again clears the state without re-planning the owner.
			{
				Config: unsetConfig,
				Check:  resource.TestCheckResourceAttr(addr, "tenant_ids.#", "0"),
			},
			{
				Config:   unsetConfig,
				PlanOnly: true,
			},
		},
	})
}
