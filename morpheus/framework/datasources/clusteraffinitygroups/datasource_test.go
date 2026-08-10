// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroups_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	agresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/clusteraffinitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const dataSourceName = "data.hpe_morpheus_cluster_affinity_groups.example"

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusListClusterAffinityGroups(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	clusterID := testhelpers.AffinityClusterID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	// Create a group so the list is guaranteed to return at least one element.
	// Without a fixture the data source can return an empty set, which passes
	// every check while never exercising how set elements are built.
	resourceConfig, err := agresource.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig := `
data "hpe_morpheus_cluster_affinity_groups" "example" {
  cluster_id = ` + clusterID + `
  depends_on = [hpe_morpheus_cluster_affinity_group.example]
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "affinity_groups.#"),
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceName, "affinity_groups.*", map[string]string{
							"name":          name,
							"affinity_type": "KEEP_TOGETHER",
						}),
				),
			},
		},
	})
}

func TestAccMorpheusListClusterAffinityGroupsWithFilter(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	clusterID := testhelpers.AffinityClusterID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	// The rendered resource is a KEEP_TOGETHER group, so it must survive the
	// filter below and give the filtered set at least one element.
	resourceConfig, err := agresource.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig := `
data "hpe_morpheus_cluster_affinity_groups" "example" {
  cluster_id = ` + clusterID + `

  filter {
    name   = "affinity_type"
    values = ["KEEP_TOGETHER"]
  }

  depends_on = [hpe_morpheus_cluster_affinity_group.example]
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "affinity_groups.#"),
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceName, "affinity_groups.*", map[string]string{
							"name":          name,
							"affinity_type": "KEEP_TOGETHER",
						}),
				),
			},
		},
	})
}
