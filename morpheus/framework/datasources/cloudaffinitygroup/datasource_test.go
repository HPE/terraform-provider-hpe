// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroup_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cloudaffinitygroup"
	agresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cloudaffinitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusFindCloudAffinityGroupByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware, capabilities.AffinityGroup)

	cloudID := testhelpers.AffinityCloudID(t)
	poolID := testhelpers.AffinityPoolID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	// The data source needs a group that is known to exist, so the test creates
	// one rather than assuming the appliance already has one with a given name.
	resourceConfig, err := agresource.RenderCloudAffinityGroupConfig(t, map[string]string{
		"CloudId": cloudID,
		"Name":    name,
		"PoolId":  poolID,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig := `
data "hpe_morpheus_cloud_affinity_group" "example" {
  name       = "` + name + `"
  cloud_id   = ` + cloudID + `
  depends_on = [hpe_morpheus_cloud_affinity_group.example]
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

func TestAccMorpheusFindCloudAffinityGroupById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware, capabilities.AffinityGroup)

	cloudID := testhelpers.AffinityCloudID(t)
	poolID := testhelpers.AffinityPoolID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	// The data source needs a group that is known to exist, so the test creates
	// one and looks it up by the id the resource reports rather than guessing.
	resourceConfig, err := agresource.RenderCloudAffinityGroupConfig(t, map[string]string{
		"CloudId": cloudID,
		"Name":    name,
		"PoolId":  poolID,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := cloudaffinitygroup.RenderAffinityGroupByIdConfig(t, map[string]string{
		"Id":      "hpe_morpheus_cloud_affinity_group.example.id",
		"CloudId": cloudID,
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

func TestAccMorpheusFindCloudAffinityGroupNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware, capabilities.AffinityGroup)

	cloudID := testhelpers.AffinityCloudID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_cloud_affinity_group" "example" {
  name     = "nonexistent-affinity-group-name-that-should-not-exist"
  cloud_id = ` + cloudID + `
}
`

	expected := regexp.MustCompile(`no cloud affinity group found`)

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

func TestAccMorpheusFindCloudAffinityGroupNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware, capabilities.AffinityGroup)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// cloud_id is required by the schema but never used: the data source
	// rejects the config for having no id or name before it looks a cloud up,
	// so this test needs no real affinity-group-capable cloud and deliberately
	// does not consult TF_VAR_testacc_morpheus_affinity_cloud_id.
	config := providerConfig + `
data "hpe_morpheus_cloud_affinity_group" "test" {
  cloud_id = 1
}
`

	expected := cloudaffinitygroup.ErrorNoValidSearchTerms

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
	ds := "data.hpe_morpheus_cloud_affinity_group.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "cloud_id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "active"),
	}
}
