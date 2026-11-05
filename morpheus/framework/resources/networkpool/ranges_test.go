// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const rangesAddr = "hpe_morpheus_network_pool.ranges"

type rng struct{ start, end string }

// rangesConfig renders a morpheus-type pool with the given ip_address_ranges.
func rangesConfig(name string, ranges ...rng) string {
	var b strings.Builder
	for _, r := range ranges {
		fmt.Fprintf(&b, "    { starting_address = %q, ending_address = %q },\n", r.start, r.end)
	}

	return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_pool" "ranges" {
  name           = "` + name + `"
  type_code      = "morpheus"
  subnet_address = "10.9.0.0"
  netmask        = "255.255.255.0"
  gateway        = "10.9.0.1"

  ip_address_ranges = [
` + b.String() + `  ]
}
`
}

// singleRangeConfig renders the same pool through the deprecated ip_ranges.
func singleRangeConfig(name string, r rng) string {
	return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_pool" "ranges" {
  name           = "` + name + `"
  type_code      = "morpheus"
  subnet_address = "10.9.0.0"
  netmask        = "255.255.255.0"
  gateway        = "10.9.0.1"

  ip_ranges = {
    starting_address = "` + r.start + `"
    ending_address   = "` + r.end + `"
  }
}
`
}

// rangeIDRecorder remembers a range's id across steps so edits can be shown
// to happen in place (same id) rather than as delete-and-create.
type rangeIDRecorder struct{ ids map[int]string }

func (r *rangeIDRecorder) record(index int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rangesAddr]
		if !ok {
			return fmt.Errorf("%s not in state", rangesAddr)
		}
		id := rs.Primary.Attributes[fmt.Sprintf("ip_address_ranges.%d.id", index)]
		if id == "" {
			return fmt.Errorf("ip_address_ranges.%d.id is empty", index)
		}
		if r.ids == nil {
			r.ids = map[int]string{}
		}
		r.ids[index] = id

		return nil
	}
}

func (r *rangeIDRecorder) expectSame(index int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[rangesAddr]
		got := rs.Primary.Attributes[fmt.Sprintf("ip_address_ranges.%d.id", index)]
		if want := r.ids[index]; got != want {
			return fmt.Errorf("ip_address_ranges.%d.id changed: %s -> %s (the range was recreated, not edited)", index, want, got)
		}

		return nil
	}
}

func inPlace() resource.ConfigPlanChecks {
	return resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(rangesAddr, plancheck.ResourceActionUpdate),
	}}
}

// TestAccMorpheusNetworkPoolResourceRangesLifecycleOk covers a multi-range
// pool end to end: create with two ranges, add a third, edit one range's end
// address in place (its id does not change), remove one, and import.
func TestAccMorpheusNetworkPoolResourceRangesLifecycleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)
	testhelpers.SkipUnlessAcc(t)

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	r1, r2, r3 := rng{"10.9.0.10", "10.9.0.20"}, rng{"10.9.0.30", "10.9.0.40"}, rng{"10.9.0.50", "10.9.0.60"}
	r2edited := rng{"10.9.0.30", "10.9.0.45"}
	ids := &rangeIDRecorder{}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: rangesConfig(name, r1, r2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.#", "2"),
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.0.starting_address", r1.start),
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.1.ending_address", r2.end),
					resource.TestCheckNoResourceAttr(rangesAddr, "ip_ranges.starting_address"),
					ids.record(0), ids.record(1),
				),
			},
			{Config: rangesConfig(name, r1, r2), PlanOnly: true},
			{
				// Add a third range: the first two keep their ids.
				Config:           rangesConfig(name, r1, r2, r3),
				ConfigPlanChecks: inPlace(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.#", "3"),
					ids.expectSame(0), ids.expectSame(1), ids.record(2),
				),
			},
			{
				// Edit the second range's end address: same id, new address.
				Config:           rangesConfig(name, r1, r2edited, r3),
				ConfigPlanChecks: inPlace(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.1.ending_address", r2edited.end),
					ids.expectSame(0), ids.expectSame(1), ids.expectSame(2),
				),
			},
			{
				// Remove the first range: the others keep their ids and move up.
				Config:           rangesConfig(name, r2edited, r3),
				ConfigPlanChecks: inPlace(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.#", "2"),
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.0.starting_address", r2edited.start),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[rangesAddr]
						if got := rs.Primary.Attributes["ip_address_ranges.0.id"]; got != ids.ids[1] {
							return fmt.Errorf("range %s should keep id %s after the first range was removed, got %s", r2edited.start, ids.ids[1], got)
						}

						return nil
					},
				),
			},
			{Config: rangesConfig(name, r2edited, r3), PlanOnly: true},
			{
				// Edit and remove in one apply: positional ids are not trusted
				// when the list shape changes, so the edited range is recreated
				// (new id) while the untouched one keeps its id. Documented
				// trade-off; a reserved range would be kept and reported instead.
				Config:           rangesConfig(name, rng{"10.9.0.50", "10.9.0.65"}),
				ConfigPlanChecks: inPlace(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.#", "1"),
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.0.ending_address", "10.9.0.65"),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[rangesAddr]
						if got := rs.Primary.Attributes["ip_address_ranges.0.id"]; got == ids.ids[2] {
							return fmt.Errorf("edit+remove must recreate the edited range (expected a new id, got the old %s)", got)
						}

						return nil
					},
				),
			},
			{Config: rangesConfig(name, rng{"10.9.0.50", "10.9.0.65"}), PlanOnly: true},
			{
				ResourceName:      rangesAddr,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccMorpheusNetworkPoolResourceRangesSwitchFromDeprecatedOk proves the
// migration story: a pool created with the deprecated ip_ranges, then the same
// values moved to ip_address_ranges, plans no change.
func TestAccMorpheusNetworkPoolResourceRangesSwitchFromDeprecatedOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)
	testhelpers.SkipUnlessAcc(t)

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	r := rng{"10.9.0.100", "10.9.0.120"}
	r2 := rng{"10.9.0.130", "10.9.0.140"}
	ids := &rangeIDRecorder{}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: singleRangeConfig(name, r),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rangesAddr, "ip_ranges.starting_address", r.start),
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.#", "1"),
					resource.TestCheckResourceAttrSet(rangesAddr, "ip_address_ranges.0.id"),
					ids.record(0),
				),
			},
			{
				// Same values under the new attribute: empty plan. The deprecated
				// attribute stays mirrored in state (it is computed) until its
				// removal in 2.6; only the configuration has moved.
				Config: rangesConfig(name, r),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				}},
				Check: resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.0.starting_address", r.start),
			},
			{Config: rangesConfig(name, r), PlanOnly: true},
			{
				// Add a range after the move. ip_ranges is still mirrored in
				// state, so Update must take the path the CONFIGURATION names:
				// both ranges present afterwards, the original keeping its id.
				Config: rangesConfig(name, r, r2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.#", "2"),
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.0.starting_address", r.start),
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.1.starting_address", r2.start),
					ids.expectSame(0),
				),
			},
			{Config: rangesConfig(name, r, r2), PlanOnly: true},
		},
	})
}

// TestAccMorpheusNetworkPoolResourceRangesValidation: exactly one of the two
// attributes; the deprecated one refuses a pool that has several ranges.
func TestAccMorpheusNetworkPoolResourceRangesValidation(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)
	testhelpers.SkipUnlessAcc(t)

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	r1, r2 := rng{"10.9.0.130", "10.9.0.140"}, rng{"10.9.0.150", "10.9.0.160"}

	both := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_pool" "ranges" {
  name      = "` + name + `"
  type_code = "morpheus"
  ip_ranges = { starting_address = "` + r1.start + `", ending_address = "` + r1.end + `" }
  ip_address_ranges = [{ starting_address = "` + r1.start + `", ending_address = "` + r1.end + `" }]
}
`
	neither := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_pool" "ranges" {
  name      = "` + name + `"
  type_code = "morpheus"
}
`

	missingEnd := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_pool" "ranges" {
  name      = "` + name + `"
  type_code = "morpheus"
  ip_address_ranges = [{ starting_address = "` + r1.start + `" }]
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{Config: missingEnd, ExpectError: regexp.MustCompile(`Missing range address`)},
			{
				Config: both,
				ExpectError: regexp.MustCompile(
					`(?s)Invalid Attribute Combination.*(ip_address_ranges.*ip_ranges|ip_ranges.*ip_address_ranges)`),
			},
			{
				Config: neither,
				ExpectError: regexp.MustCompile(
					`(?s)Invalid Attribute Combination.*No attribute specified when one \(and only one\) of`),
			},
			{
				Config: rangesConfig(name, r1, r2),
				Check:  resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.#", "2"),
			},
			{
				// Switching a two-range pool to the deprecated single-range
				// attribute would delete a range: refused before any call.
				Config:      singleRangeConfig(name, r1),
				ExpectError: regexp.MustCompile(`ip_ranges cannot manage a pool with several ranges`),
			},
			{
				// Back to the list: unchanged on the appliance.
				Config: rangesConfig(name, r1, r2),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				}},
			},
		},
	})
}

// TestAccMorpheusNetworkPoolRangesDataSourceOk: the data source lists the
// pool's ranges.
func TestAccMorpheusNetworkPoolRangesDataSourceOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)
	testhelpers.SkipUnlessAcc(t)

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	r1, r2 := rng{"10.9.0.170", "10.9.0.180"}, rng{"10.9.0.190", "10.9.0.200"}
	config := rangesConfig(name, r1, r2) + `
data "hpe_morpheus_network_pool" "ranges" {
  id = hpe_morpheus_network_pool.ranges.id
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.hpe_morpheus_network_pool.ranges", "ip_address_ranges.#", "2"),
				resource.TestCheckTypeSetElemNestedAttrs("data.hpe_morpheus_network_pool.ranges", "ip_address_ranges.*", map[string]string{
					"starting_address": r1.start, "ending_address": r1.end,
				}),
				resource.TestCheckTypeSetElemNestedAttrs("data.hpe_morpheus_network_pool.ranges", "ip_address_ranges.*", map[string]string{
					"starting_address": r2.start, "ending_address": r2.end,
				}),
			),
		}},
	})
}

// TestAccMorpheusNetworkPoolRangesWithReservedIPOk exercises the interplay
// with hpe_morpheus_network_pool_ip: a reservation pinned in the second range
// survives an in-place edit of that range; removing the range while the
// reservation exists is refused with a clear error; removing both works.
func TestAccMorpheusNetworkPoolRangesWithReservedIPOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)
	testhelpers.SkipUnlessAcc(t)

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	host := strings.ToLower(acctest.RandomWithPrefix("tfacc-pool-ip"))
	r1, r2 := rng{"10.9.0.210", "10.9.0.220"}, rng{"10.9.0.230", "10.9.0.240"}
	r2edited := rng{"10.9.0.230", "10.9.0.245"}
	ipInR2 := "10.9.0.235"
	ids := &rangeIDRecorder{}

	withIP := func(pool string) string {
		return pool + `
resource "hpe_morpheus_network_pool_ip" "pinned" {
  network_pool_id = hpe_morpheus_network_pool.ranges.id
  hostname        = "` + host + `"
  ip_address      = "` + ipInR2 + `"
}
`
	}
	ipAddr := "hpe_morpheus_network_pool_ip.pinned"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: withIP(rangesConfig(name, r1, r2)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ipAddr, "ip_address", ipInR2),
					ids.record(0), ids.record(1),
				),
			},
			{
				// Grow the range holding the reservation: in place, IP intact.
				Config:           withIP(rangesConfig(name, r1, r2edited)),
				ConfigPlanChecks: inPlace(),
				Check: resource.ComposeAggregateTestCheckFunc(
					ids.expectSame(1),
					resource.TestCheckResourceAttr(ipAddr, "ip_address", ipInR2),
				),
			},
			{
				// Drop the range while the reservation still exists: refused.
				Config:      withIP(rangesConfig(name, r1)),
				ExpectError: regexp.MustCompile(`still has reserved addresses`),
			},
			{
				// The removal is still planned, not silently forgotten: the same
				// configuration fails the same way again.
				Config:      withIP(rangesConfig(name, r1)),
				ExpectError: regexp.MustCompile(`still has reserved addresses`),
			},
			{
				// After the refusals the kept range is in state with its id
				// (Check does not run on an ExpectError step, hence this one):
				// putting it back in the configuration plans nothing.
				Config: withIP(rangesConfig(name, r1, r2edited)),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.#", "2"),
					ids.expectSame(1),
				),
			},
			{
				// Drop the reservation and the range together: the IP is
				// released first (it depends on the pool), then the range goes.
				Config: rangesConfig(name, r1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.#", "1"),
					resource.TestCheckResourceAttr(rangesAddr, "ip_address_ranges.0.starting_address", r1.start),
					ids.expectSame(0),
				),
			},
			{Config: rangesConfig(name, r1), PlanOnly: true},
		},
	})
}
