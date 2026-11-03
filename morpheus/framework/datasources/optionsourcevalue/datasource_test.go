// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optionsourcevalue_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// These tests read option sources on the shared acceptance appliance; they
// provision nothing and run in seconds. They target its HVM cluster "Duck"
// (cloud 1, resource pool 1, layout 1193), as the cluster tests do.

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func accTest(t *testing.T) {
	t.Helper()

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()
}

func singleStep(t *testing.T, config string, check resource.TestCheckFunc) resource.TestCase {
	t.Helper()

	return resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{{
			Config: testhelpers.ProviderBlock() + config,
			Check:  check,
		}},
	}
}

func expectError(t *testing.T, config string, pattern string) resource.TestCase {
	t.Helper()

	return resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{{
			Config:      testhelpers.ProviderBlock() + config,
			ExpectError: regexp.MustCompile(pattern),
		}},
	}
}

// The second motivating case: the HVM cluster's resource pool, by cluster name,
// in the form config.resourcePoolId takes, with poolProviderType alongside.
func TestAccMorpheusOptionSourceValueClusterPoolByName(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	const ds = "data.hpe_morpheus_option_source_value.cluster_pool"

	resource.Test(t, singleStep(t, `
data "hpe_morpheus_option_source_value" "cluster_pool" {
  source    = "zonePools"
  cloud_id   = 1
  layout_id = 1193
  name      = "Duck"
}`, resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(ds, "name", "Duck"),
		resource.TestCheckResourceAttr(ds, "value", "pool-1"),
		resource.TestCheckResourceAttr(ds, "attributes.providerType", "mvm"),
		resource.TestCheckResourceAttr(ds, "attributes.id", "1"),
	)))
}

// Selecting by value works too, and the name comes back.
func TestAccMorpheusOptionSourceValueClusterPoolByValue(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	const ds = "data.hpe_morpheus_option_source_value.by_value"

	resource.Test(t, singleStep(t, `
data "hpe_morpheus_option_source_value" "by_value" {
  source    = "zonePools"
  cloud_id   = 1
  layout_id = 1193
  value     = "pool-1"
}`, resource.TestCheckResourceAttr(ds, "name", "Duck")))
}

// Several matches is an error that lists the candidates and both remedies.
func TestAccMorpheusOptionSourceValueMultipleMatches(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	resource.Test(t, expectError(t, `
data "hpe_morpheus_option_source_value" "x" {
  source    = "zonePools"
  cloud_id   = 1
  layout_id = 1193
}`, `(?s)multiple\s+matches\s+for\s+option\s+source\s+"zonePools".*`+
		`Expected\s+exactly\s+one\s+value,\s+found\s+\d+:.*Duck\s+\(pool-1\).*`+
		`Narrow\s+the\s+result\s+with\s+name\s+or\s+value,\s+or\s+use\s+hpe_morpheus_option_source_values`))
}

// A name nothing has is an error that lists what the source did return.
func TestAccMorpheusOptionSourceValueNoMatchListsValues(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	resource.Test(t, expectError(t, `
data "hpe_morpheus_option_source_value" "x" {
  source    = "zonePools"
  cloud_id   = 1
  layout_id = 1193
  name      = "no-such-cluster"
}`, `(?s)no\s+match\s+for\s+option\s+source\s+"zonePools".*`+
		`None\s+of\s+the\s+\d+\s+values\s+the\s+source\s+returned\s+has\s+name\s+"no-such-cluster".*`+
		`Duck\s+\(pool-1\)`))
}

// cloud_id alone: the source returns nothing, and the error names the missing
// provisioning context rather than just "no match".
func TestAccMorpheusOptionSourceValueMissingContextHint(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	resource.Test(t, expectError(t, `
data "hpe_morpheus_option_source_value" "x" {
  source  = "zonePools"
  cloud_id = 1
  name    = "Duck"
}`, `(?s)no\s+match\s+for\s+option\s+source\s+"zonePools".*`+
		`The\s+source\s+returned\s+no\s+values\.\s+zonePools\s+requires\s+a\s+provisioning\s+context:\s+`+
		`set\s+layout_id\s+or\s+provision_type_code`))
}

// The object-shaped source with a collection resolves a single network.
func TestAccMorpheusOptionSourceValueNetworkFromCollection(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	const ds = "data.hpe_morpheus_option_source_value.mgmt"

	resource.Test(t, singleStep(t, `
data "hpe_morpheus_option_source_value" "mgmt" {
  source     = "zoneNetworkOptions"
  cloud_id    = 1
  layout_id  = 1193
  pool_id    = "pool-1"
  collection = "networks"
  value      = "network-1"
}`, resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(ds, "value", "network-1"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttr(ds, "attributes.id", "network-1"),
	)))
}
