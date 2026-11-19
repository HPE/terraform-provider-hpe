// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const dataSourceName = "data.hpe_morpheus_host.example"

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// hostID returns an existing host id to look up, skipping when unset. Feature's
// HVM cluster hosts on cloud 17 are hypervisors; any of their ids works.
func hostID(t *testing.T) string {
	t.Helper()

	v := os.Getenv("TF_VAR_testacc_morpheus_host_id")
	if v == "" {
		t.Skip("TF_VAR_testacc_morpheus_host_id not set; skipping host-by-id test")
	}

	return v
}

func hostName(t *testing.T) string {
	t.Helper()

	v := os.Getenv("TF_VAR_testacc_morpheus_host_name")
	if v == "" {
		t.Skip("TF_VAR_testacc_morpheus_host_name not set; skipping host-by-name test")
	}

	return v
}

func TestAccMorpheusHostByID(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	id := hostID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
data "hpe_morpheus_host" "example" {
  id = %s
}
`, id)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", id),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "compute_server_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cloud_id"),
				),
			},
		},
	})
}

func TestAccMorpheusHostByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	name := hostName(t)
	cloudID := testhelpers.ComputeServerCloudID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
data "hpe_morpheus_host" "example" {
  name     = %q
  cloud_id = %s
}
`, name, cloudID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "cloud_id", cloudID),
				),
			},
		},
	})
}
