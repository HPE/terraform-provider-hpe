// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hosts_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const dataSourceName = "data.hpe_morpheus_hosts.example"

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// checkHostsNotEmpty asserts the data source returned at least one host.
func checkHostsNotEmpty() resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(dataSourceName, "hosts.#",
		func(value string) error {
			count, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("hosts.# is not a number: %q", value)
			}
			if count == 0 {
				return fmt.Errorf("expected at least one host, got none")
			}

			return nil
		})
}

// TestAccMorpheusHostsByCloudID lists the HVM cluster hosts of cloud 17, which
// are hypervisors, and asserts at least one is returned with the cloud filter.
func TestAccMorpheusHostsByCloudID(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	cloudID := testhelpers.ComputeServerCloudID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
data "hpe_morpheus_hosts" "example" {
  cloud_id = %s
}
`, cloudID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "hosts.#"),
					checkHostsNotEmpty(),
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceName, "hosts.*", map[string]string{
							"cloud_id": cloudID,
						}),
				),
			},
		},
	})
}

// TestAccMorpheusHostsBareMetalFalse asserts the bare_metal = false filter
// returns at least one host on the HVM cloud (HVM hypervisors are not bare
// metal).
func TestAccMorpheusHostsBareMetalFalse(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	cloudID := testhelpers.ComputeServerCloudID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
data "hpe_morpheus_hosts" "example" {
  cloud_id   = %s
  bare_metal = false
}
`, cloudID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkHostsNotEmpty(),
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceName, "hosts.*", map[string]string{
							"bare_metal": "false",
						}),
				),
			},
		},
	})
}

// TestAccMorpheusHostsBareMetalTrue asserts the bare_metal = true filter applies
// without error on the BMaaS (Omega, 15174) cloud. The result may legitimately
// be empty, so only that hosts.# is set (a valid number, including 0) is
// asserted.
func TestAccMorpheusHostsBareMetalTrue(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.BareMetalStub)

	cloudID := os.Getenv("TF_VAR_testacc_morpheus_bmaas_cloud_id")
	if cloudID == "" {
		t.Skip("TF_VAR_testacc_morpheus_bmaas_cloud_id not set; skipping bare_metal=true test")
	}

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
data "hpe_morpheus_hosts" "example" {
  cloud_id   = %s
  bare_metal = true
}
`, cloudID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "hosts.#"),
					// Every host the filter returns is bare metal (the Omega
					// cloud always has at least its seeded node).
					checkAllHosts(dataSourceName, func(h map[string]string) error {
						if h["bare_metal"] != "true" {
							return fmt.Errorf("bare_metal=true filter returned host %s with bare_metal=%q", h["id"], h["bare_metal"])
						}

						return nil
					}),
				),
			},
		},
	})
}

// TestAccMorpheusHostsComputeServerTypeFilter checks that compute_server_type
// separates an HVM cloud's hypervisor hosts (type mvm) from the guest VMs that
// also live in /api/servers (type morpheusVm): each filter returns only its own
// type, and the hypervisors have no parent host while every guest has one.
func TestAccMorpheusHostsComputeServerTypeFilter(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	cloudID := testhelpers.ComputeServerCloudID(t)

	// A guest VM known to live on the cloud, so the guest filter is proven on
	// a real member rather than on whatever happens to be there (an HVM cloud
	// with only hypervisors would otherwise fail the all-hosts check).
	guestID := os.Getenv("TF_VAR_testacc_morpheus_hvm_guest_server_id")
	if guestID == "" {
		t.Skip("TF_VAR_testacc_morpheus_hvm_guest_server_id must be set")
	}

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
data "hpe_morpheus_hosts" "hypervisors" {
  cloud_id            = %[1]s
  compute_server_type = "mvm"
}

data "hpe_morpheus_hosts" "guests" {
  cloud_id            = %[1]s
  compute_server_type = "morpheusVm"
}
`, cloudID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAllHosts("data.hpe_morpheus_hosts.hypervisors", func(h map[string]string) error {
						if h["compute_server_type"] != "mvm" {
							return fmt.Errorf("hypervisor filter returned type %q", h["compute_server_type"])
						}
						if h["parent_host_id"] != "" {
							return fmt.Errorf("hypervisor %s has a parent host %s", h["id"], h["parent_host_id"])
						}

						return nil
					}),
					checkHasHost("data.hpe_morpheus_hosts.guests", guestID),
					checkAllHosts("data.hpe_morpheus_hosts.guests", func(h map[string]string) error {
						if h["compute_server_type"] != "morpheusVm" {
							return fmt.Errorf("guest filter returned type %q", h["compute_server_type"])
						}
						if h["parent_host_id"] == "" {
							return fmt.Errorf("guest %s has no parent host", h["id"])
						}

						return nil
					}),
				),
			},
		},
	})
}

// checkAllHosts runs check over every element of the data source's hosts set
// (attributes flattened per element) and requires at least one element.
func checkAllHosts(name string, check func(map[string]string) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("%s not in state", name)
		}
		n, err := strconv.Atoi(rs.Primary.Attributes["hosts.#"])
		if err != nil || n == 0 {
			return fmt.Errorf("%s: hosts.# = %q, want > 0", name, rs.Primary.Attributes["hosts.#"])
		}
		for i := range n {
			prefix := fmt.Sprintf("hosts.%d.", i)
			h := map[string]string{}
			for k, v := range rs.Primary.Attributes {
				if strings.HasPrefix(k, prefix) {
					h[strings.TrimPrefix(k, prefix)] = v
				}
			}
			if err := check(h); err != nil {
				return fmt.Errorf("%s element %d: %w", name, i, err)
			}
		}

		return nil
	}
}

// checkHasHost asserts that the hosts set of the named data source contains a
// host with the given id.
func checkHasHost(name, id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("%s not found in state", name)
		}

		for k, v := range rs.Primary.Attributes {
			if strings.HasPrefix(k, "hosts.") && strings.HasSuffix(k, ".id") && v == id {
				return nil
			}
		}

		return fmt.Errorf("%s: no host with id %s in hosts", name, id)
	}
}
