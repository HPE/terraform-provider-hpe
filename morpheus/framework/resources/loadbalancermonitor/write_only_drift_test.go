// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusLoadBalancerMonitorResourceWriteOnlyDriftDetected: the monitor
// password would be rotated out of band (a direct API PUT, as the Morpheus UI
// would do) to assert the appliance-side drift path. On Morpheus 9.1.x this
// cannot be staged -- see the skip below -- so the test is skipped. Gated like
// the other NSX-T monitor tests.
func TestAccMorpheusLoadBalancerMonitorResourceWriteOnlyDriftDetected(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)
	testhelpers.SkipUnlessAcc(t)

	// Morpheus 9.1.x accepts monitorPassword on update but does not change the
	// stored password (monitorPasswordHash unchanged), so an out-of-band
	// rotation cannot be staged; tracked as a Morpheus defect. Set
	// TF_ACC_LB_MONITOR_ROTATABLE=1 to run against an appliance that stores it;
	// the rotate step then verifies the hash changed.
	if os.Getenv("TF_ACC_LB_MONITOR_ROTATABLE") == "" {
		t.Skip("Morpheus 9.1.x does not store monitorPassword on update (monitorPasswordHash unchanged); " +
			"set TF_ACC_LB_MONITOR_ROTATABLE=1 on an appliance that does")
	}

	name := acctest.RandomWithPrefix(t.Name())
	lbName := name[0:min(32, len(name))]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{"Name": lbName})
	if err != nil {
		t.Fatal(err)
	}

	addr := "hpe_morpheus_load_balancer_monitor.wo"

	config := testhelpers.ProviderBlock() + lbConfig + `
resource "hpe_morpheus_load_balancer_monitor" "wo" {
  load_balancer_id    = hpe_morpheus_load_balancer.lb.id
  name                = "` + name + `"
  monitor_type        = "http"
  monitor_interval    = 5
  monitor_timeout     = 3
  monitor_destination = "/health"
  monitor_username    = "probe"
  monitor_password_wo = "Secret123!"
}
`

	// Captured from state by the create step, used by the next step's
	// PreConfig to rotate the secret out of band.
	var lbID, id int64

	capture := func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("%s not in state", addr)
		}

		var err error
		if lbID, err = strconv.ParseInt(rs.Primary.Attributes["load_balancer_id"], 10, 64); err != nil {
			return err
		}

		id, err = strconv.ParseInt(rs.Primary.ID, 10, 64)

		return err
	}

	rotate := func() {
		ctx := context.Background()

		client, err := testhelpers.NewClientForServer(ctx, "")
		if err != nil {
			t.Fatalf("client: %s", err)
		}

		rotated := "Rotated456!"
		body := sdk.UpdateLoadBalancerMonitorRequest{
			LoadBalancerMonitor: &sdk.UpdateLoadBalancerMonitorRequestLoadBalancerMonitor{
				MonitorPassword: *sdk.NewNullableString(&rotated),
			},
		}

		hash := func() string {
			resp, _, err := client.LoadBalancersAPI.GetLoadBalancerMonitor(ctx, lbID, id).Execute()
			if err != nil || resp == nil || resp.LoadBalancerMonitor == nil {
				return ""
			}
			if h := resp.LoadBalancerMonitor.MonitorPasswordHash.Get(); h != nil {
				return *h
			}

			return ""
		}
		before := hash()

		if _, httpResp, err := client.LoadBalancersAPI.
			UpdateLoadBalancerMonitor(ctx, lbID, id).
			UpdateLoadBalancerMonitorRequest(body).Execute(); err != nil {
			t.Fatalf("out-of-band rotate failed: %s (status %v)", err, httpResp)
		}

		if after := hash(); after == before {
			t.Fatalf("rotate PUT did not change monitorPasswordHash (%q); the appliance did not store the "+
				"new password, so drift cannot be detected", before)
		}
	}

	// If Morpheus started honouring the monitorPassword update, the plan after
	// the rotation would re-send the configured value (ResourceActionUpdate).
	// The test is skipped above because 9.1.x does not, so this is never
	// reached; it is kept so the body stays compilable.
	update := resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
	}}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{Config: config, Check: capture},
			{
				PreConfig:        rotate,
				Config:           config,
				ConfigPlanChecks: update,
			},
			{Config: config, PlanOnly: true},
		},
	})
}
