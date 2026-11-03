// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/skip"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusClusterResourceWriteOnlyOk: the nested, create-only
// server.ssh_password_wo is tracked by digest (MORPH-18005). Morpheus cannot
// change it after creation, so a change to it -- or to its deprecated marker
// -- plans a replacement of the cluster; an unrelated update leaves it alone;
// the plan is empty otherwise. Skipped by default like the other cluster
// tests: it needs the HVM cloud, group, key pair and hosts they use.
func TestAccMorpheusClusterResourceWriteOnlyOk(t *testing.T) {
	capabilities.MustHaveOrSkip(t, capabilities.All)

	if skip.SkipByDefault(t) {
		t.Log("Skipping test not opted in via RUN_SKIPPED_BY_DEFAULT")

		return
	}

	defer testhelpers.RecordResult(t)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())
	addr := "hpe_morpheus_cluster.wo"

	config := func(password, extra, description string) string {
		return testhelpers.ProviderBlock() + `
data "hpe_morpheus_cloud" "test" {
	name = "hvm"
}

data "hpe_morpheus_group" "test" {
	name = "HVM"
}

data "hpe_morpheus_service_plan" "test" {
	name = "Default Manual"
	provision_type_code = "manual"
}

resource "hpe_morpheus_cluster" "wo" {
	name        = "` + name + `"
	description = "` + description + `"
	cloud_id    = data.hpe_morpheus_cloud.test.id
	group_id    = data.hpe_morpheus_group.test.id
	layout_id   = 449

	server = {
		service_plan_id = data.hpe_morpheus_service_plan.test.id
		ssh_username    = "hvmserviceuser"
		ssh_password_wo = "` + password + `"` + extra + `

		ssh_hosts = [
			{
				name = "` + name + `-worker-1"
				ip   = "10.118.67.20"
			},
		]
	}
}
`
	}

	action := func(a plancheck.ResourceActionType) resource.ConfigPlanChecks {
		return resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(addr, a)}}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: config("Secret123!", "", "initial"),
				Check:  resource.TestCheckResourceAttrSet(addr, "write_only_revision"),
			},
			{Config: config("Secret123!", "", "initial"), PlanOnly: true},
			{
				// An unrelated change is an in-place update.
				Config:           config("Secret123!", "", "changed"),
				ConfigPlanChecks: action(plancheck.ResourceActionUpdate),
			},
			{
				// The password cannot be changed in place: replacement.
				Config:           config("Secret456!", "", "changed"),
				ConfigPlanChecks: action(plancheck.ResourceActionReplace),
			},
			{Config: config("Secret456!", "", "changed"), PlanOnly: true},
			{
				// Adding the deprecated marker re-sends -- which here means
				// replacement too.
				Config:           config("Secret456!", "\n\t\tssh_password_wo_version = 1", "changed"),
				ConfigPlanChecks: action(plancheck.ResourceActionReplace),
			},
			{
				// Removing it is one in-place update (digests recorded, nothing
				// sent), then clean.
				Config:           config("Secret456!", "", "changed"),
				ConfigPlanChecks: action(plancheck.ResourceActionUpdate),
			},
			{Config: config("Secret456!", "", "changed"), PlanOnly: true},
		},
	})
}
