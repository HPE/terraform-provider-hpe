// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package compute_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dscompute "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/compute"
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

func TestAccMorpheusDataSourceResourcePoolExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.ResourcePool)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dependenciesConfig := testhelpers.WhoamiBlock()

	if currentDependency, err := compute.RenderResourcePoolGroupConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	datasourceConfig, err := dscompute.RenderResourcePoolConfig(t, map[string]string{
		"Name": "resource.hpe_morpheus_resource_pool_group.example.name",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_resource_pool.example",
			"name",
			name,
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

// The three tests below target the HVM cluster "Duck" (id 1) in cloud 1 on the
// shared acceptance appliance, as the hpe_morpheus_cluster tests do. None
// provisions anything; each is a single read.

// A lookup by the cluster's name must fail — the cloud listing does not include
// cluster pools — and the error must hand the user the pool id via the cluster
// hint, which is the situation that motivated the change.
func TestAccMorpheusDataSourceResourcePoolByNameMissesClusterPoolWithHint(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	datasourceConfig, err := dscompute.RenderResourcePoolConfig(t, map[string]string{
		"Name":    `"Duck"`,
		"CloudId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + datasourceConfig,
				ExpectError: regexp.MustCompile(
					`(?s)No resource pool named "Duck" was found in cloud 1\..*` +
						`returns cloud-level pools only.*` +
						`A cluster named "Duck" exists in cloud 1; its provisioning pool has id 1\.`,
				),
			},
		},
	})
}

// The same pool is reachable by id, and reports the type Morpheus gives a
// cluster's pool.
func TestAccMorpheusDataSourceResourcePoolByIDFindsClusterPool(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	datasourceConfig, err := dscompute.RenderResourcePoolByIDConfig(t, map[string]string{
		"Id":      "1",
		"CloudId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + datasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hpe_morpheus_resource_pool.by_id", "id", "1"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_resource_pool.by_id", "name", "Duck"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_resource_pool.by_id", "type", "namespace"),
				),
			},
		},
	})
}

// An id that does not exist is reported as such, not as a raw HTTP error.
func TestAccMorpheusDataSourceResourcePoolByIDNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	datasourceConfig, err := dscompute.RenderResourcePoolByIDConfig(t, map[string]string{
		"Id":      "999999",
		"CloudId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      testhelpers.ProviderBlock() + datasourceConfig,
				ExpectError: regexp.MustCompile(`Resource pool id 999999 was not found in cloud 1\.`),
			},
		},
	})
}
