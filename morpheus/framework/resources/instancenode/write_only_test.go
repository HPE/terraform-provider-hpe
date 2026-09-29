// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusInstanceNodeWriteOnlyOk: the write-only ssh_password of a
// pre-provisioned node is sent at creation (it used to be read from the plan,
// where a write-only value is always null, so it was never sent) and a change
// to it plans a replacement (MORPH-18005). Needs a pre-provisioned server:
// TF_VAR_metal_instance_id, TF_VAR_node_ssh_host, TF_VAR_node_ssh_username.
func TestAccMorpheusInstanceNodeWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Metal)
	testhelpers.SkipUnlessAcc(t)

	instanceID := os.Getenv("TF_VAR_metal_instance_id")
	host, user := os.Getenv("TF_VAR_node_ssh_host"), os.Getenv("TF_VAR_node_ssh_username")
	if instanceID == "" || host == "" || user == "" {
		t.Skip("TF_VAR_metal_instance_id, TF_VAR_node_ssh_host and TF_VAR_node_ssh_username not set")
	}

	addr := "hpe_morpheus_instance_node.wo"
	config := func(password string) string {
		return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_instance_node" "wo" {
  instance_id  = ` + instanceID + `
  ssh_host     = "` + host + `"
  ssh_username = "` + user + `"
  ssh_password = "` + password + `"
}
`
	}

	rec := testhelpers.RecordRequests(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				PreConfig: rec.Reset,
				Config:    config("Secret123!"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(addr, "write_only_revision"),
					// The password must reach the add-node action.
					rec.ExpectSent("/api/instances", "action.sshPassword", true),
				),
			},
			{Config: config("Secret123!"), PlanOnly: true},
			{
				Config: config("Secret456!"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(addr, plancheck.ResourceActionReplace),
				}},
			},
			{Config: config("Secret456!"), PlanOnly: true},
		},
	})
}
