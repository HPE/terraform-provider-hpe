// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup_test

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/clusteraffinitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusClusterAffinityGroupResourceExampleOk tests create, read, and import.
func TestAccMorpheusClusterAffinityGroupResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	clusterID := testhelpers.AffinityClusterID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := clusteraffinitygroup.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_cluster_affinity_group.example"

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "cluster_id", clusterID),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "affinity_type", "KEEP_TOGETHER"),
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
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["cluster_id"] + "." +
						rs.Primary.Attributes["id"], nil
				},
			},
		},
	})
}

// TestAccMorpheusClusterAffinityGroupResourceUpdateOk tests update (name change, in-place).
func TestAccMorpheusClusterAffinityGroupResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	clusterID := testhelpers.AffinityClusterID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig, err := clusteraffinitygroup.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := clusteraffinitygroup.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      updatedName,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_cluster_affinity_group.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "cluster_id", clusterID),
		resource.TestCheckResourceAttr(resourceName, "name", name),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "cluster_id", clusterID),
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
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

// TestAccMorpheusClusterAffinityGroupResourceRequiresReplace verifies affinity_type change
// triggers a replacement (RequiresReplace plan modifier).
func TestAccMorpheusClusterAffinityGroupResourceRequiresReplace(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	clusterID := testhelpers.AffinityClusterID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	createConfig, err := clusteraffinitygroup.RenderClusterAffinityGroupConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Build a config that changes affinity_type to KEEP_SEPARATE.
	replaceConfig := fmt.Sprintf(`
resource "hpe_morpheus_cluster_affinity_group" "example" {
  cluster_id    = %s
  name          = "%s"
  affinity_type = "KEEP_SEPARATE"
}
`, clusterID, name)

	resourceName := "hpe_morpheus_cluster_affinity_group.example"

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

// TestAccMorpheusClusterAffinityGroupResourceMembersOk creates a group that
// owns its membership. An HVM instance is provisioned first and the group is
// created with the instance's compute server as its member -- the only order
// that works on Morpheus 9.1, which refuses a group with no members. The group
// is then renamed with servers still set, which is the update path that sends
// the configured set; imported; and finally left without servers, which hands
// the membership back without evicting anyone: the member is still there and
// the plan is clean.
//
// The instance's configuration resolves the HVM fixture by name, including the
// cluster the group is placed on, so no id-bearing variable is needed.
func TestAccMorpheusClusterAffinityGroupResourceMembersOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	// A short instance name: Morpheus derives the guest's hostname from it,
	// and a test name with a random suffix is longer than a hostname label.
	instanceConfig, err := instance.RenderInstanceHVMConfig(t, map[string]string{
		"Name": acctest.RandomWithPrefix("TestAccMorpheusAGMember"),
	})
	if err != nil {
		t.Fatal(err)
	}

	groupConfig := func(groupName string, withServers bool) string {
		overrides := map[string]string{
			"ClusterId": "data.hpe_morpheus_cluster.hvm.id",
			"Name":      groupName,
		}
		if withServers {
			overrides["Servers"] = "hpe_morpheus_instance.example.compute_servers"
		}

		cfg, err := clusteraffinitygroup.RenderClusterAffinityGroupMembersConfig(t, overrides)
		if err != nil {
			t.Fatal(err)
		}

		return providerConfig + instanceConfig + "\n" + cfg
	}

	resourceName := "hpe_morpheus_cluster_affinity_group.members"
	instanceName := "hpe_morpheus_instance.example"

	// The group's members are exactly the instance's compute servers.
	memberChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrPair(resourceName, "cluster_id", "data.hpe_morpheus_cluster.hvm", "id"),
		resource.TestCheckResourceAttr(resourceName, "affinity_type", "KEEP_TOGETHER"),
		resource.TestCheckResourceAttr(resourceName, "servers.#", "1"),
		checkSameSet(resourceName, "servers", instanceName, "compute_servers"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: groupConfig(name, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					memberChecks,
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{Config: groupConfig(name, true), ExpectNonEmptyPlan: false, PlanOnly: true},
			{
				Config:           groupConfig(updatedName, true),
				ConfigPlanChecks: checkInPlaceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					memberChecks,
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["cluster_id"] + "." +
						rs.Primary.Attributes["id"], nil
				},
			},
			// servers dropped from the configuration: membership is no longer
			// Terraform's to assert, and the member it has is kept.
			{
				Config: groupConfig(updatedName, false),
				Check:  memberChecks,
			},
			{Config: groupConfig(updatedName, false), ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}

// checkSameSet asserts that two set attributes in state hold the same values,
// whatever order or index the state gives them.
func checkSameSet(firstName, firstAttr, secondName, secondAttr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		first, err := setValues(s, firstName, firstAttr)
		if err != nil {
			return err
		}

		second, err := setValues(s, secondName, secondAttr)
		if err != nil {
			return err
		}

		if len(first) == 0 {
			return fmt.Errorf("%s.%s is empty in state", firstName, firstAttr)
		}

		if !maps.Equal(first, second) {
			return fmt.Errorf("%s.%s %v and %s.%s %v differ",
				firstName, firstAttr, slices.Sorted(maps.Keys(first)),
				secondName, secondAttr, slices.Sorted(maps.Keys(second)))
		}

		return nil
	}
}

// setValues returns the values of a set attribute in state, keyed by value.
// A set is stored as <attr>.# plus one <attr>.<index> entry per element.
func setValues(s *terraform.State, name, attr string) (map[string]struct{}, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return nil, fmt.Errorf("%s not found in state", name)
	}

	values := make(map[string]struct{})
	for k, v := range rs.Primary.Attributes {
		if k == attr+".#" || !strings.HasPrefix(k, attr+".") {
			continue
		}

		values[v] = struct{}{}
	}

	return values, nil
}
