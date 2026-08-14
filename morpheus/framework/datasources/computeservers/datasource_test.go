// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package computeservers_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const dataSourceName = "data.hpe_morpheus_compute_servers.example"

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// checkServersNotEmpty asserts the data source returned at least one server. An
// empty result would satisfy every other check in these tests while never
// exercising how set elements are built.
func checkServersNotEmpty() resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(dataSourceName, "servers.#",
		func(value string) error {
			count, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("servers.# is not a number: %q", value)
			}
			if count == 0 {
				return fmt.Errorf("expected at least one server, got none")
			}

			return nil
		})
}

// checkSomeServerHasParentHost asserts that at least one returned server reports
// the hypervisor host it runs on.
//
// This exists because API contract validation flags parentServer as possibly
// absent from the implementation, which would make parent_host_id and
// parent_host_name silently null. It is not absent: it is populated for guests
// and null for hypervisor hosts, which have no parent. A sample made up of host
// records therefore looks like a missing field.
//
// Asserting "at least one" rather than naming a server keeps the check honest
// without tying it to a particular guest that someone may later delete.
func checkSomeServerHasParentHost() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("data source %s not found in state", dataSourceName)
		}

		count, err := strconv.Atoi(rs.Primary.Attributes["servers.#"])
		if err != nil {
			return fmt.Errorf("servers.# is not a number: %w", err)
		}

		hosts := 0

		for i := range count {
			name := rs.Primary.Attributes[fmt.Sprintf("servers.%d.parent_host_name", i)]
			id := rs.Primary.Attributes[fmt.Sprintf("servers.%d.parent_host_id", i)]

			if name != "" && id != "" && id != "0" {
				hosts++
			}
		}

		if hosts == 0 {
			return fmt.Errorf(
				"no server reported parent_host_name/parent_host_id out of %d returned; "+
					"either the API stopped returning parentServer, or the result "+
					"contained only hypervisor hosts", count)
		}

		return nil
	}
}

func TestAccMorpheusListComputeServers(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_compute_servers" "example" {
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "servers.#"),
					checkSomeServerHasParentHost(),
				),
			},
		},
	})
}

func TestAccMorpheusListComputeServersWithCloudId(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	// No fallback: cloud IDs are allocated per appliance, and a guessed "1" is
	// as likely to name a cloud with no hosts as one that does not exist.
	cloudID := testhelpers.ComputeServerCloudID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_compute_servers" "example" {
  cloud_id = ` + cloudID + `
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "servers.#"),
					checkServersNotEmpty(),
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceName, "servers.*", map[string]string{
							"cloud_id": cloudID,
						}),
				),
			},
		},
	})
}

func TestAccMorpheusListComputeServersWithInstanceId(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	// The instance filter needs an instance that actually owns a host. Most
	// hosts on an appliance belong to no instance, so this cannot be derived
	// from EnvComputeServerID and gets its own variable.
	instanceID := testhelpers.ComputeServerInstanceID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_compute_servers" "example" {
  instance_id = ` + instanceID + `
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "servers.#"),
					checkServersNotEmpty(),
					resource.TestCheckTypeSetElemNestedAttrs(
						dataSourceName, "servers.*", map[string]string{
							"instance_id": instanceID,
						}),
				),
			},
		},
	})
}

func TestAccMorpheusListComputeServersWithFilter(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_compute_servers" "example" {
  filter {
    name   = "status"
    values = ["provisioned"]
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "servers.#"),
				),
			},
		},
	})
}
