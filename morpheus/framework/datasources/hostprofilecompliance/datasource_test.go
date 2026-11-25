// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofilecompliance_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusHostProfileComplianceDataSourceExactlyOne verifies the
// ExactlyOneOf config validator rejects configs that set neither cluster_id
// nor server_id.
func TestAccMorpheusHostProfileComplianceDataSourceExactlyOne(t *testing.T) {
	defer testhelpers.RecordResult(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_host_profile_compliance" "example" {
}
`

	expected := regexp.MustCompile(`Invalid Attribute Combination|Missing Attribute Configuration|Exactly one of these attributes`)

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

// TestAccMorpheusHostProfileComplianceDataSourceCluster exercises the cluster_id
// Read branch against a live HVM cluster, confirming the compliance report is
// fetched and mapped into state.
func TestAccMorpheusHostProfileComplianceDataSourceCluster(t *testing.T) {
	defer testhelpers.RecordResult(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	capabilities.MustHaveOrSkip(t, capabilities.HostProfile)

	clusterID := testhelpers.HVMClusterID(t)

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := fmt.Sprintf(`
data "hpe_morpheus_host_profile_compliance" "cluster" {
  cluster_id = %s
}
`, clusterID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_host_profile_compliance.cluster",
						"cluster_id", clusterID,
					),
				),
			},
		},
	})
}

// TestAccMorpheusHostProfileComplianceDataSourceServer exercises the server_id
// Read branch against a live compute server, confirming the per-server
// compliance report is fetched and mapped into state.
func TestAccMorpheusHostProfileComplianceDataSourceServer(t *testing.T) {
	defer testhelpers.RecordResult(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	capabilities.MustHaveOrSkip(t, capabilities.HostProfile)

	serverID := testhelpers.ComputeServerID(t)

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := fmt.Sprintf(`
data "hpe_morpheus_host_profile_compliance" "server" {
  server_id = %s
}
`, serverID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_host_profile_compliance.server",
						"server_id", serverID,
					),
				),
			},
		},
	})
}
