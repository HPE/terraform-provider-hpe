// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optionsourcevalues_test

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// These tests read option sources on the shared acceptance appliance; they
// provision nothing and run in seconds. The provisioning cases target its HVM
// cluster "Duck" (cloud 1, resource pool 1, layout 1193), as the cluster and
// resource-pool tests do.

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

// listLen reads the length of a list attribute from state.
func listLen(s *terraform.State, address, attr string) (int, error) {
	rs, ok := s.RootModule().Resources[address]
	if !ok {
		return 0, fmt.Errorf("%s not in state", address)
	}

	raw, ok := rs.Primary.Attributes[attr+".#"]
	if !ok {
		return 0, fmt.Errorf("%s: %s.# not in state", address, attr)
	}

	return strconv.Atoi(raw)
}

// checkListNotEmpty asserts a list attribute has at least one element.
func checkListNotEmpty(address, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		n, err := listLen(s, address, attr)
		if err != nil {
			return err
		}

		if n == 0 {
			return fmt.Errorf("%s: %s is empty", address, attr)
		}

		return nil
	}
}

// checkEveryElem applies want to every element of a list attribute. want gets
// the element's attributes with the "values.N." prefix stripped, so keys read
// "value" or "attributes.providerType".
func checkEveryElem(address, attr string, want func(elem map[string]string) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		n, err := listLen(s, address, attr)
		if err != nil {
			return err
		}

		rs := s.RootModule().Resources[address]

		for i := 0; i < n; i++ {
			prefix := fmt.Sprintf("%s.%d.", attr, i)
			elem := map[string]string{}

			for k, v := range rs.Primary.Attributes {
				if len(k) > len(prefix) && k[:len(prefix)] == prefix {
					elem[k[len(prefix):]] = v
				}
			}

			if err := want(elem); err != nil {
				return fmt.Errorf("%s: %s[%d]: %w (element: %v)", address, attr, i, err, elem)
			}
		}

		return nil
	}
}

// matchAll asserts each key of elem matches its pattern.
func matchAll(elem map[string]string, keyPatterns ...string) error {
	for i := 0; i+1 < len(keyPatterns); i += 2 {
		key, pattern := keyPatterns[i], keyPatterns[i+1]

		v, ok := elem[key]
		if !ok {
			return fmt.Errorf("missing %q", key)
		}

		if !regexp.MustCompile(pattern).MatchString(v) {
			return fmt.Errorf("%s = %q does not match %s", key, v, pattern)
		}
	}

	return nil
}

// A source that needs no parameters returns rows.
func TestAccMorpheusOptionSourceValuesGroups(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	const ds = "data.hpe_morpheus_option_source_values.groups"

	resource.Test(t, singleStep(t, `
data "hpe_morpheus_option_source_values" "groups" {
  source = "groups"
}`, resource.ComposeAggregateTestCheckFunc(
		checkListNotEmpty(ds, "values"),
		resource.TestCheckResourceAttrSet(ds, "values.0.name"),
		resource.TestCheckResourceAttrSet(ds, "values.0.value"),
	)))
}

// The wrong-bucket regression: with cloud_id the networks list is narrowed to
// the cloud. If the provider ever stopped sending the top-level key, the server
// would silently return every network on the appliance instead.
func TestAccMorpheusOptionSourceValuesNetworksFilteredByZone(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	resource.Test(t, singleStep(t, `
data "hpe_morpheus_option_source_values" "in_cloud" {
  source  = "networks"
  cloud_id = 1
}

data "hpe_morpheus_option_source_values" "all" {
  source = "networks"
}`, func(s *terraform.State) error {
		inCloud, err := listLen(s, "data.hpe_morpheus_option_source_values.in_cloud", "values")
		if err != nil {
			return err
		}

		all, err := listLen(s, "data.hpe_morpheus_option_source_values.all", "values")
		if err != nil {
			return err
		}

		if inCloud == 0 || inCloud >= all {
			return fmt.Errorf("networks in cloud 1 = %d, all = %d; want 0 < in_cloud < all", inCloud, all)
		}

		return nil
	}))
}

// zonePools with a layout returns the cluster pools, with value in the form
// the instance API takes and providerType in attributes; header rows are gone;
// an attributes filter narrows the list; a non-matching filter empties it.
func TestAccMorpheusOptionSourceValuesZonePools(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	resource.Test(t, singleStep(t, `
data "hpe_morpheus_option_source_values" "pools" {
  source    = "zonePools"
  cloud_id   = 1
  layout_id = 1193
}

data "hpe_morpheus_option_source_values" "mvm_pools" {
  source    = "zonePools"
  cloud_id   = 1
  layout_id = 1193

  filter {
    name   = "attributes.providerType"
    values = ["^mvm$"]
  }
}

data "hpe_morpheus_option_source_values" "no_such_pools" {
  source    = "zonePools"
  cloud_id   = 1
  layout_id = 1193

  filter {
    name   = "name"
    values = ["^no-such-cluster$"]
  }
}`, resource.ComposeAggregateTestCheckFunc(
		checkListNotEmpty("data.hpe_morpheus_option_source_values.pools", "values"),
		checkEveryElem("data.hpe_morpheus_option_source_values.pools", "values", func(elem map[string]string) error {
			return matchAll(elem, "value", `^pool-\d+$`, "attributes.providerType", `^mvm$`, "name", `.+`)
		}),
		checkListNotEmpty("data.hpe_morpheus_option_source_values.mvm_pools", "values"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_option_source_values.no_such_pools", "values.#", "0"),
	)))
}

// zoneNetworkOptions is object-shaped: collection selects the list, and each
// value is the id in the prefixed form the instance API takes.
func TestAccMorpheusOptionSourceValuesZoneNetworkOptions(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	const ds = "data.hpe_morpheus_option_source_values.networks"

	resource.Test(t, singleStep(t, `
data "hpe_morpheus_option_source_values" "networks" {
  source     = "zoneNetworkOptions"
  cloud_id    = 1
  layout_id  = 1193
  pool_id    = "pool-1"
  collection = "networks"
}`, resource.ComposeAggregateTestCheckFunc(
		checkListNotEmpty(ds, "values"),
		checkEveryElem(ds, "values", func(elem map[string]string) error {
			return matchAll(elem, "value", `^network-\d+$`, "attributes.id", `^network-\d+$`)
		}),
	)))
}

// Without collection an object-shaped source is an error that names the lists.
func TestAccMorpheusOptionSourceValuesObjectNeedsCollection(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	resource.Test(t, expectError(t, `
data "hpe_morpheus_option_source_values" "x" {
  source    = "zoneNetworkOptions"
  cloud_id   = 1
  layout_id = 1193
}`, `(?s)returns\s+an\s+object,\s+not\s+a\s+list.*set\s+collection\s+to\s+one\s+of:\s+`+
		`networkGroups,\s+networkSubnets,\s+networkTypes,\s+networks`))
}

// zonePools with cloud_id alone is legitimately empty for the plural, with the
// missing-context hint raised as a warning. The testing framework does not
// surface warnings, so this asserts an empty list and a successful plan; the
// hint text is pinned by the singular test and the unit tests.
func TestAccMorpheusOptionSourceValuesMissingContextIsEmptyNotError(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	resource.Test(t, singleStep(t, `
data "hpe_morpheus_option_source_values" "x" {
  source  = "zonePools"
  cloud_id = 1
}`, resource.TestCheckResourceAttr("data.hpe_morpheus_option_source_values.x", "values.#", "0")))
}

// An unknown source fails on the server with the generic error; ours must say
// the name is the likely cause and show the request.
func TestAccMorpheusOptionSourceValuesUnknownSource(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	resource.Test(t, expectError(t, `
data "hpe_morpheus_option_source_values" "x" {
  source  = "noSuchOptionSource"
  cloud_id = 1
}`, `(?s)returned HTTP 500.*not\s+in\s+the\s+provider's\s+catalogue.*`+
		`name\s+is\s+the\s+most\s+likely\s+cause.*GET\s+/api/options/noSuchOptionSource`))
}

// A source defined only under a service type, called without it, fails on the
// server; ours must point at service_type.
func TestAccMorpheusOptionSourceValuesNeedsServiceType(t *testing.T) {
	defer testhelpers.RecordResult(t)
	accTest(t)

	resource.Test(t, expectError(t, `
data "hpe_morpheus_option_source_values" "x" {
  source = "scvmmLibraryShares"
}`, `(?s)returned HTTP 500.*advice\s+on\s+service_type.*GET\s+/api/options/scvmmLibraryShares`))
}
