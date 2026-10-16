// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusInstanceHVMUserDataOk (MORPH-19115): config_hvm.user_data
// is sent at create, stored by Morpheus and read back on import (it forces
// replacement, so import must recover it or the next plan would recreate the
// instance), a re-apply plans nothing, and a change plans a replacement.
// Whether cloud-init ran the content is not observable from the API.
func TestAccMorpheusInstanceHVMUserDataOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	// No `packages`: the test network may have no route to a package mirror,
	// and a package install that cannot reach one blocks cloud-init -- and so
	// Morpheus's own runcmd and IP detection -- until the provision times out.
	userData := "      #cloud-config\n      write_files:\n        - path: /etc/tf-user-data\n          content: \"" + name + "\\n\""
	changed := userData + "\n      runcmd:\n        - true"

	render := func(ud string) string {
		cfg, err := instance.RenderInstanceHVMConfig(t, map[string]string{"Name": name, "UserData": ud})
		if err != nil {
			t.Fatal(err)
		}

		return testhelpers.ProviderBlock() + cfg
	}
	addr := "hpe_morpheus_instance.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: render(userData),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(addr, "config_hvm.user_data"),
					resource.TestMatchResourceAttr(addr, "config_hvm.user_data", regexp.MustCompile(`^#cloud-config\n`)),
				),
			},
			{Config: render(userData), PlanOnly: true},
			{
				ResourceName:      addr,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"network_interfaces",
					"volumes.0.datastore_auto_selection",
				},
			},
			{
				// A changed value cannot reach the running VM, so it replaces.
				Config: render(changed),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(addr, plancheck.ResourceActionReplace),
				}},
				Check: resource.TestMatchResourceAttr(addr, "config_hvm.user_data", regexp.MustCompile(`runcmd:`)),
			},
			{Config: render(changed), PlanOnly: true},
		},
	})
}
