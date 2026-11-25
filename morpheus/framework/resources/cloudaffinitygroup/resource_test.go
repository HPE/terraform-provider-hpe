// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroup_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cloudaffinitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusCloudAffinityGroupResourceExampleOk tests create, read, and import.
func TestAccMorpheusCloudAffinityGroupResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware, capabilities.AffinityGroup)

	cloudID := testhelpers.AffinityCloudID(t)
	serverID := testhelpers.AffinityMemberServerID(t)
	poolID := testhelpers.AffinityPoolID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	// Serialised across packages: every test seeds its group with the same
	// fixture VM, and on a VMware cloud Morpheus pushes each group to vCenter
	// as a DRS rule -- vCenter refuses the same VM in several KEEP_TOGETHER
	// rules at once, so the groups must be created one at a time.
	testhelpers.LockFixture(t, "affinity-member-server")

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := cloudaffinitygroup.RenderCloudAffinityGroupConfig(t, map[string]string{
		"CloudId": cloudID,
		"Name":    name,
		"PoolId":  poolID,
		"Servers": "[" + serverID + "]",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_cloud_affinity_group.example"

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "cloud_id", cloudID),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "affinity_type", "KEEP_TOGETHER"),
		// Morpheus rejects the create without a pool, so it is always configured.
		resource.TestCheckResourceAttr(resourceName, "pool_id", poolID),
		// CRITICAL BEHAVIOUR 1: active must default to true even when omitted.
		resource.TestCheckResourceAttr(resourceName, "active", "true"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			// CRITICAL BEHAVIOUR 6: Import with composite ID.
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["cloud_id"] + "." +
						rs.Primary.Attributes["id"], nil
				},
			},
		},
	})
}

// TestAccMorpheusCloudAffinityGroupResourceUpdateOk tests update (name change, in-place).
func TestAccMorpheusCloudAffinityGroupResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware, capabilities.AffinityGroup)

	cloudID := testhelpers.AffinityCloudID(t)
	serverID := testhelpers.AffinityMemberServerID(t)
	poolID := testhelpers.AffinityPoolID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	// Serialised across packages: every test seeds its group with the same
	// fixture VM, and on a VMware cloud Morpheus pushes each group to vCenter
	// as a DRS rule -- vCenter refuses the same VM in several KEEP_TOGETHER
	// rules at once, so the groups must be created one at a time.
	testhelpers.LockFixture(t, "affinity-member-server")

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	// The in-place change is `active`, not `name`: on Morpheus 9.1 VMware a renamed
	// cloud group's DRS rule is leaked in vCenter when the group is destroyed
	// (MORPH-20098), which poisons the fixture VM for every later test.

	createConfig, err := cloudaffinitygroup.RenderCloudAffinityGroupConfig(t, map[string]string{
		"CloudId": cloudID,
		"Name":    name,
		"PoolId":  poolID,
		"Servers": "[" + serverID + "]",
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := cloudaffinitygroup.RenderCloudAffinityGroupConfig(t, map[string]string{
		"CloudId": cloudID,
		"Name":    name,
		"Active":  "false",
		"PoolId":  poolID,
		"Servers": "[" + serverID + "]",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_cloud_affinity_group.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "cloud_id", cloudID),
		resource.TestCheckResourceAttr(resourceName, "name", name),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "cloud_id", cloudID),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "active", "false"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{Config: providerConfig + createConfig, Check: createChecks},
			{
				Config:           providerConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{Config: providerConfig + updateConfig, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}

// TestAccMorpheusCloudAffinityGroupResourceRequiresReplace verifies affinity_type change
// triggers a replacement (RequiresReplace plan modifier).
func TestAccMorpheusCloudAffinityGroupResourceRequiresReplace(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware, capabilities.AffinityGroup)

	cloudID := testhelpers.AffinityCloudID(t)
	serverID := testhelpers.AffinityMemberServerID(t)
	poolID := testhelpers.AffinityPoolID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	// Serialised across packages: every test seeds its group with the same
	// fixture VM, and on a VMware cloud Morpheus pushes each group to vCenter
	// as a DRS rule -- vCenter refuses the same VM in several KEEP_TOGETHER
	// rules at once, so the groups must be created one at a time.
	testhelpers.LockFixture(t, "affinity-member-server")

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	createConfig, err := cloudaffinitygroup.RenderCloudAffinityGroupConfig(t, map[string]string{
		"CloudId": cloudID,
		"Name":    name,
		"PoolId":  poolID,
		"Servers": "[" + serverID + "]",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Build a config that changes affinity_type to KEEP_SEPARATE.
	// We can't use RenderCloudAffinityGroupConfig because the template hardcodes
	// KEEP_TOGETHER, so inline the HCL directly. pool is repeated verbatim
	// because Morpheus rejects a cloud affinity group created without one.
	replaceConfig := fmt.Sprintf(`
resource "hpe_morpheus_cloud_affinity_group" "example" {
  cloud_id        = %s
  name            = "%s"
  affinity_type   = "KEEP_SEPARATE"
  pool_id         = %s
  servers         = [%s]
}
`, cloudID, name, poolID, serverID)

	resourceName := "hpe_morpheus_cloud_affinity_group.example"

	checkReplace := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				resourceName, plancheck.ResourceActionDestroyBeforeCreate,
			),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.TestCheckResourceAttr(
					resourceName, "affinity_type", "KEEP_TOGETHER",
				),
			},
			{
				Config:           providerConfig + replaceConfig,
				ConfigPlanChecks: checkReplace,
				Check: resource.TestCheckResourceAttr(
					resourceName, "affinity_type", "KEEP_SEPARATE",
				),
			},
		},
	})
}
