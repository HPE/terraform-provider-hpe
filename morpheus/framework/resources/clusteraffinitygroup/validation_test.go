// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusClusterAffinityGroupResourceValidationEmptyName verifies that
// an empty name is rejected at plan time rather than being sent to the API,
// which otherwise returns a misleading 403 (MORPH-16730).
func TestAccMorpheusClusterAffinityGroupResourceValidationEmptyName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
resource "hpe_morpheus_cluster_affinity_group" "test" {
  cluster_id    = 1
  name          = ""
  affinity_type = "KEEP_TOGETHER"
  active        = true
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("at least 1"),
			},
		},
	})
}
