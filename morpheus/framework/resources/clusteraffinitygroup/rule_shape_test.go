// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/clusteraffinitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusClusterAffinityGroupRuleShapeGate (MORPH-20083): on Morpheus
// 9.1+ a group with no initial_servers is refused at plan with the provider's
// own diagnostic, not the API's; on older appliances the empty group is
// created. The test picks its expectation from the appliance version.
func TestAccMorpheusClusterAffinityGroupRuleShapeGate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)
	clusterID := testhelpers.AffinityClusterID(t)
	testhelpers.SkipUnlessAcc(t)

	// Serialised with the cloud tests: all of them seed their group with the
	// same fixture VM. HVM does not push DRS rules, so sharing it is probably
	// harmless here, but one lock keeps the question from ever arising.
	testhelpers.LockFixture(t, "affinity-member-server")

	enforced := applianceEnforcesRuleShape(t)
	name := acctest.RandomWithPrefix(t.Name())
	config, err := clusteraffinitygroup.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID, "Name": name, // no Servers: an empty group
	})
	if err != nil {
		t.Fatal(err)
	}

	step := resource.TestStep{Config: testhelpers.ProviderBlock() + config}
	if enforced {
		step.PlanOnly = true
		step.ExpectError = regexp.MustCompile(`Affinity groups need a member on this appliance`)
	} else {
		step.Check = resource.TestCheckResourceAttr("hpe_morpheus_cluster_affinity_group.example", "servers.#", "0")
	}

	steps := []resource.TestStep{step}
	if enforced {
		// A replacement creates a new group too: a seeded group whose
		// affinity_type changes while initial_servers is dropped must be
		// refused at plan, not fail after the destroy.
		serverID := testhelpers.AffinityMemberServerID(t)
		seeded, err := clusteraffinitygroup.RenderClusterAffinityGroupConfig(t, map[string]string{
			"ClusterId": clusterID, "Name": name, "Servers": "[" + serverID + "]",
		})
		if err != nil {
			t.Fatal(err)
		}
		replaceUnseeded := fmt.Sprintf(`
resource "hpe_morpheus_cluster_affinity_group" "example" {
  cluster_id    = %s
  name          = "%s"
  affinity_type = "KEEP_SEPARATE"
}
`, clusterID, name)
		// An in-place update to servers = [] would empty the group, which
		// 9.1 refuses too; the plan must refuse it before the apply fails.
		emptied := fmt.Sprintf(`
resource "hpe_morpheus_cluster_affinity_group" "example" {
  cluster_id    = %s
  name          = "%s"
  affinity_type = "KEEP_TOGETHER"
  servers       = []
}
`, clusterID, name)
		steps = append(steps,
			resource.TestStep{Config: testhelpers.ProviderBlock() + seeded},
			resource.TestStep{
				Config:      testhelpers.ProviderBlock() + emptied,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Affinity groups need a member on this appliance`),
			},
			resource.TestStep{
				Config:      testhelpers.ProviderBlock() + replaceUnseeded,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Affinity groups need a member on this appliance`),
			},
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}

func applianceEnforcesRuleShape(t *testing.T) bool {
	t.Helper()

	client, err := testhelpers.NewClientForServer(t.Context(), "")
	if err != nil {
		t.Fatal(err)
	}
	v, err := versioncheck.Appliance(t.Context(), client)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := versioncheck.Satisfies(v, constants.AffinityRuleShapeMinVersion)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("appliance %s enforces affinity rule shape: %v", v.Original(), ok)

	return ok
}
