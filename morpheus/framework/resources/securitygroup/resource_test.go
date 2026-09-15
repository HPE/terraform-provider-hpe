package securitygroup_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/securitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusSecurityGroupResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := securitygroup.RenderSecurityGroupConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("hpe_morpheus_security_group.example", "name", name),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_security_group.example",
			"description",
			"Security group for web servers",
		),
		resource.TestCheckResourceAttr("hpe_morpheus_security_group.example", "active", "true"),
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
				ResourceName:      "hpe_morpheus_security_group.example",
			},
		},
	})
}

func TestAccMorpheusSecurityGroupResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedDescription := "Updated security group for web servers"

	createConfig, err := securitygroup.RenderSecurityGroupConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := securitygroup.RenderSecurityGroupConfig(t, map[string]string{
		"Name":        name,
		"Description": updatedDescription,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_security_group.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", "Security group for web servers"),
		resource.TestCheckResourceAttr(resourceName, "active", "true"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
		resource.TestCheckResourceAttr(resourceName, "active", "true"),
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
				Config: providerConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:           providerConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

// permissionGroupID is the group the restricted-access tests grant access to.
// Group 1 is the primary group present on every appliance, and is what the
// other acceptance tests in this repository use for a group id.
const permissionGroupID = "1"

// permissionConfig renders a security group whose group access is restricted to
// specific groups: resource_permission_groups_all = false together with a list
// of group ids -- the documented (and shipped-example) way to express "only
// these groups", and the configuration MORPH-16355 found the resource rejected.
func permissionConfig(name string, groupsAll bool) string {
	all := "true"
	ids := ""
	if !groupsAll {
		all = "false"
		ids = "\n  resource_permission_group_ids  = [" + permissionGroupID + "]"
	}

	return `
resource "hpe_morpheus_security_group" "example" {
  name                           = "` + name + `"
  description                    = "Security group with restricted group access"
  visibility                     = "private"
  resource_permission_groups_all = ` + all + ids + `
}
`
}

// TestAccMorpheusSecurityGroupResourcePermissionsOnCreate is the MORPH-16355
// regression guard for two defects that hid behind one another:
//
//  1. A ConflictsWith validator forbade setting resource_permission_group_ids
//     together with resource_permission_groups_all, so the documented
//     "specific groups" configuration failed at plan time.
//  2. Once that was reachable, the create endpoint turned out to ignore
//     resourcePermissions, so a group created with all = false came back as
//     all = true -- an "inconsistent result after apply". The resource now
//     applies the permissions with a follow-up update after the create.
//
// The checks assert the values READ BACK from the API, not the plan, so the
// test fails if either the validator returns or the create path stops
// applying permissions. The steady-state plan and the import both complete the
// round trip.
func TestAccMorpheusSecurityGroupResourcePermissionsOnCreate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_security_group.example"

	config := providerConfig + permissionConfig(name, false)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "resource_permission_groups_all", "false"),
		resource.TestCheckResourceAttr(resourceName, "resource_permission_group_ids.#", "1"),
		resource.TestCheckTypeSetElemAttr(resourceName, "resource_permission_group_ids.*", permissionGroupID),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
			{
				Config:             config,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

// TestAccMorpheusSecurityGroupResourcePermissionsUpdate covers the update
// direction: a group created open to all groups is then restricted to one
// group, and the restriction is read back from the API. It also exercises the
// two attributes changing together in a single in-place update, which the
// removed validator made impossible to even plan.
func TestAccMorpheusSecurityGroupResourcePermissionsUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_security_group.example"

	openConfig := providerConfig + permissionConfig(name, true)
	restrictedConfig := providerConfig + permissionConfig(name, false)

	openChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "resource_permission_groups_all", "true"),
		resource.TestCheckNoResourceAttr(resourceName, "resource_permission_group_ids.#"),
	)
	restrictedChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "resource_permission_groups_all", "false"),
		resource.TestCheckResourceAttr(resourceName, "resource_permission_group_ids.#", "1"),
		resource.TestCheckTypeSetElemAttr(resourceName, "resource_permission_group_ids.*", permissionGroupID),
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
				Config: openConfig,
				Check:  openChecks,
			},
			{
				Config:           restrictedConfig,
				Check:            restrictedChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             restrictedConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
