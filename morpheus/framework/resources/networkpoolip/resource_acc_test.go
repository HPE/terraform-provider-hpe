// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

const resourceName = "hpe_morpheus_network_pool_ip.test"

// requirePoolID skips the test unless TF_ACC_NETWORK_POOL_ID names a pool the
// appliance under test exposes (the test appliance has pool id 1 covering
// 192.168.100.10-199).
func requirePoolID(t *testing.T) string {
	t.Helper()

	poolID := os.Getenv("TF_ACC_NETWORK_POOL_ID")
	if poolID == "" {
		t.Skip("TF_ACC_NETWORK_POOL_ID not set; skipping network pool IP test")
	}

	return poolID
}

// TestAccMorpheusNetworkPoolIpAutoAllocateOk exercises the main feature: a
// reservation created with a hostname only is allocated the pool's next free
// address. It then adds and removes a ref_type/ref_id association, and imports.
func TestAccMorpheusNetworkPoolIpAutoAllocateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	poolID := requirePoolID(t)
	hostname := acctest.RandomWithPrefix(t.Name())
	hostname = hostname[0:min(60, len(hostname))]

	// The association target. Morpheus stores ref_type/ref_id as an opaque
	// pair and does not check that the referenced object exists, so by default
	// the test uses a type nothing resolves, which attaches the reservation to
	// nothing on the appliance. Set TF_ACC_NETWORK_POOL_IP_REF_ID to the id of a
	// real server to exercise a ComputeServer association instead.
	refType, refID := "TerraformAcceptanceTest", "424242"
	if id := os.Getenv("TF_ACC_NETWORK_POOL_IP_REF_ID"); id != "" {
		refType, refID = "ComputeServer", id
	}

	config := func(extra string) string {
		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_ip" "test" {
  network_pool_id = %s
  hostname        = %q
%s
}
`, poolID, hostname, extra)
	}

	// Only assert the allocated address matches the known CIDR prefix when the
	// test is told what that prefix is; otherwise just confirm an address was
	// set, so the test is portable across appliances.
	ipChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
	}
	if prefix := os.Getenv("TF_ACC_NETWORK_POOL_CIDR_PREFIX"); prefix != "" {
		ipChecks = append(ipChecks, resource.TestMatchResourceAttr(
			resourceName, "ip_address", regexp.MustCompile(prefix),
		))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(""),
				Check:  resource.ComposeAggregateTestCheckFunc(ipChecks...),
			},
			{
				Config: config(fmt.Sprintf("  ref_type = %q\n  ref_id   = %s", refType, refID)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ref_type", refType),
					resource.TestCheckResourceAttr(resourceName, "ref_id", refID),
				),
			},
			{
				Config: config(""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "ref_type"),
					resource.TestCheckNoResourceAttr(resourceName, "ref_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportStateVerifyIgnore: []string{"create_dns", "domain_id"},
			},
		},
	})
}

// TestAccMorpheusNetworkPoolIpPinnedOk reserves a specific in-range address.
func TestAccMorpheusNetworkPoolIpPinnedOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	poolID := requirePoolID(t)

	ipAddress := os.Getenv("TF_ACC_NETWORK_POOL_IP")
	if ipAddress == "" {
		t.Skip("TF_ACC_NETWORK_POOL_IP not set; skipping pinned address test")
	}

	hostname := acctest.RandomWithPrefix(t.Name())
	hostname = hostname[0:min(60, len(hostname))]

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_ip" "test" {
  network_pool_id = %s
  hostname        = %q
  ip_address      = %q
}
`, poolID, hostname, ipAddress)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
			},
		},
	})
}

// TestAccMorpheusNetworkPoolIpOutOfRangeErr confirms an out-of-range address
// is rejected with Morpheus's own message. Terraform re-wraps the diagnostic
// at roughly 80 columns, so the pattern tolerates the inserted whitespace.
func TestAccMorpheusNetworkPoolIpOutOfRangeErr(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	poolID := requirePoolID(t)
	hostname := acctest.RandomWithPrefix(t.Name())
	hostname = hostname[0:min(60, len(hostname))]

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_ip" "test" {
  network_pool_id = %s
  hostname        = %q
  ip_address      = "203.0.113.5"
}
`, poolID, hostname)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`not\s+within\s+the\s+pool`),
			},
		},
	})
}

// importStateIDFunc builds the "<network_pool_id>:<id>" import ID from state.
func importStateIDFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", name)
		}

		return fmt.Sprintf("%s:%s",
			rs.Primary.Attributes["network_pool_id"],
			rs.Primary.Attributes["id"],
		), nil
	}
}
