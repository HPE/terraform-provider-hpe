// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/clusteraffinitygroup"
	agresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/clusteraffinitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

func TestAccMorpheusFindClusterAffinityGroupByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	clusterID := testhelpers.AffinityClusterID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := agresource.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig := `
data "hpe_morpheus_cluster_affinity_group" "example" {
  name       = "` + name + `"
  cluster_id = ` + clusterID + `
  depends_on = [hpe_morpheus_cluster_affinity_group.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(affinityGroupChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterAffinityGroupById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	clusterID := testhelpers.AffinityClusterID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := agresource.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := clusteraffinitygroup.RenderAffinityGroupByIdConfig(t, map[string]string{
		"Id":        "hpe_morpheus_cluster_affinity_group.example.id",
		"ClusterId": clusterID,
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(affinityGroupChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterAffinityGroupNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	clusterID := testhelpers.AffinityClusterID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_cluster_affinity_group" "example" {
  name       = "nonexistent-affinity-group-name-that-should-not-exist"
  cluster_id = ` + clusterID + `
}
`

	expected := regexp.MustCompile(`no cluster affinity group found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindClusterAffinityGroupNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	t.Parallel()

	// cluster_id is required by the schema but never used: the data source
	// rejects the config for having no id or name before it looks a cluster up,
	// so this test needs no real affinity-group-capable cluster and
	// deliberately does not consult TF_VAR_testacc_morpheus_affinity_cluster_id.
	config := providerConfigOffline + `
      data "hpe_morpheus_cluster_affinity_group" "test" {
        cluster_id = 1
      }`

	expected := clusteraffinitygroup.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func affinityGroupChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_cluster_affinity_group.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "cluster_id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "active"),
	}
}
