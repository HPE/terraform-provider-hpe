// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusLoadBalancerResourceWriteOnlyOk: the write-only group_id is
// taken by Morpheus only at creation, so a change plans a replacement
// (MORPH-18005); the deprecated group_id_version marker still replaces when
// bumped and does nothing when removed; both import forms leave a clean plan.
// Uses an HAProxy load balancer (never NSX-T). Needs two groups the caller may
// create load balancers in: TF_ACC_LB_GROUP_ID and TF_ACC_LB_GROUP_ID_2, plus
// the HAProxy cloud/plan the other tests use.
func TestAccMorpheusLoadBalancerResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkLoadBalancer, capabilities.NetworkLoadBalancerHAProxy)
	testhelpers.SkipUnlessAcc(t)

	g1, g2 := os.Getenv("TF_ACC_LB_GROUP_ID"), os.Getenv("TF_ACC_LB_GROUP_ID_2")
	if g1 == "" || g2 == "" {
		t.Skip("TF_ACC_LB_GROUP_ID and TF_ACC_LB_GROUP_ID_2 must be set")
	}

	name := acctest.RandomWithPrefix("tfwo-lb")
	name = name[0:min(32, len(name))]

	c := testhelpers.WriteOnlyCase{ //nolint:gosec // attribute names, not credentials
		Address:   "hpe_morpheus_load_balancer.wo",
		Attribute: "group_id",
		Marker:    "group_id_version",
		APIPath:   "/api/load-balancers",
		APIKey:    "loadBalancer.groupId",
		Config: func(value, extra string, unrelated bool) string {
			desc := "write-only matrix"
			if unrelated {
				desc = "write-only matrix, changed"
			}

			return testhelpers.ProviderBlock() + `
data "hpe_morpheus_cloud" "example" {
  name = "hvm"
}

resource "hpe_morpheus_load_balancer" "wo" {
  name        = "` + name + `"
  description = "` + desc + `"
  cloud_id    = data.hpe_morpheus_cloud.example.id
  group_id    = ` + value + extra + `
  visibility  = "private"

  config_haproxy = {
    plan_id = 8
    pool    = "pool-1"
  }

  permissions = {
    all = true
  }
}
`
		},
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyReplaceSteps(c, rec, g1, g2)
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, g2, "config", "config_haproxy", "cloud_id")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
