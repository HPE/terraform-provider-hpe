// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// The shape of the tree tenantTreeConfig builds: every node has
// tenantTreeFanout children and the tree is tenantTreeDepth levels deep, so
// with a fanout of 2 that is 2 + 4 + 8 = 14 tenants. The fanout is substituted
// into the configuration; the depth is fixed by its three explicit levels, so
// tenantTreeDepth must stay at 3 unless the configuration grows a level.
const (
	tenantTreeFanout = 2
	tenantTreeDepth  = 3

	// tenantTreeRoleAddress is the single base role every tenant in the tree
	// uses. A tenant's own base role is always assignable to its children, so
	// one role serves the whole tree.
	tenantTreeRoleAddress = "hpe_morpheus_role.example"
)

// tenantTreeConfig returns the configuration for a tree of tenants rooted at
// the master tenant. Level 1 omits parent_id, so those tenants are created
// under the caller's own tenant -- the master tenant, the only tenant permitted
// to nominate a parent. Levels 2 and 3 nominate their parent with parent_id,
// referencing the level above by for_each key. Every tenant uses the base role
// at tenantTreeRoleAddress.
//
// Both the name and the subdomain of every tenant are "<name>-<key>", where
// key is the node's for_each key ("1", "1-2", "1-2-1", ...). tenantTreeNode
// mirrors the keys and descriptions so the checks agree with the
// configuration.
func tenantTreeConfig(name, description, accountName string) string {
	return fmt.Sprintf(`
locals {
  tenant_tree_fanout = %[1]d

  tenant_tree_level_1 = [
    for i in range(local.tenant_tree_fanout) : {
      key    = tostring(i + 1)
      parent = null
    }
  ]

  tenant_tree_level_2 = flatten([
    for node in local.tenant_tree_level_1 : [
      for i in range(local.tenant_tree_fanout) : {
        key    = "${node.key}-${i + 1}"
        parent = node.key
      }
    ]
  ])

  tenant_tree_level_3 = flatten([
    for node in local.tenant_tree_level_2 : [
      for i in range(local.tenant_tree_fanout) : {
        key    = "${node.key}-${i + 1}"
        parent = node.key
      }
    ]
  ])
}

# Level 1: children of the master tenant (no parent_id).
resource "hpe_morpheus_tenant" "tenant_tree_l1" {
  for_each = { for node in local.tenant_tree_level_1 : node.key => node }

  name            = "%[2]s-${each.key}"
  description     = "%[3]s level 1 node ${each.key}, child of the master tenant"
  enabled         = true
  subdomain       = "%[2]s-${each.key}"
  base_role_id    = %[4]s.id
  currency        = "USD"
  account_number  = "1337"
  account_name    = "%[5]s"
  customer_number = "1337"
}

# Level 2: children of the level 1 tenants.
resource "hpe_morpheus_tenant" "tenant_tree_l2" {
  for_each = { for node in local.tenant_tree_level_2 : node.key => node }

  name            = "%[2]s-${each.key}"
  description     = "%[3]s level 2 node ${each.key}, child of ${each.value.parent}"
  enabled         = true
  subdomain       = "%[2]s-${each.key}"
  base_role_id    = %[4]s.id
  parent_id       = hpe_morpheus_tenant.tenant_tree_l1[each.value.parent].id
  currency        = "USD"
  account_number  = "1337"
  account_name    = "%[5]s"
  customer_number = "1337"
}

# Level 3: children of the level 2 tenants.
resource "hpe_morpheus_tenant" "tenant_tree_l3" {
  for_each = { for node in local.tenant_tree_level_3 : node.key => node }

  name            = "%[2]s-${each.key}"
  description     = "%[3]s level 3 node ${each.key}, child of ${each.value.parent}"
  enabled         = true
  subdomain       = "%[2]s-${each.key}"
  base_role_id    = %[4]s.id
  parent_id       = hpe_morpheus_tenant.tenant_tree_l2[each.value.parent].id
  currency        = "USD"
  account_number  = "1337"
  account_name    = "%[5]s"
  customer_number = "1337"
}
`, tenantTreeFanout, name, description, tenantTreeRoleAddress, accountName)
}

// tenantTreeNode is one tenant in the tree, enumerated the same way the
// template's locals enumerate it so the checks and the configuration agree on
// for_each keys and therefore on resource addresses.
type tenantTreeNode struct {
	// depth is the tree level, 1 for the children of the master tenant.
	depth int
	// key is the for_each key: "1", "1-2", "1-2-1", ...
	key string
	// parent is the node's parent, or nil at level 1 (the master tenant).
	parent *tenantTreeNode
}

// address is the node's resource address as it appears in Terraform's JSON
// state and plan, which is how statecheck and plancheck locate instances of a
// for_each resource.
func (n tenantTreeNode) address() string {
	return fmt.Sprintf(`hpe_morpheus_tenant.tenant_tree_l%d["%s"]`, n.depth, n.key)
}

// description mirrors the description tenantTreeConfig gives the node.
func (n tenantTreeNode) description(prefix string) string {
	if n.parent == nil {
		return fmt.Sprintf("%s level 1 node %s, child of the master tenant", prefix, n.key)
	}

	return fmt.Sprintf("%s level %d node %s, child of %s", prefix, n.depth, n.key, n.parent.key)
}

// tenantTreeNodes enumerates the tree breadth-first, level 1 first.
func tenantTreeNodes() []*tenantTreeNode {
	var nodes, level []*tenantTreeNode

	for i := 1; i <= tenantTreeFanout; i++ {
		level = append(level, &tenantTreeNode{depth: 1, key: strconv.Itoa(i)})
	}
	nodes = append(nodes, level...)

	for depth := 2; depth <= tenantTreeDepth; depth++ {
		var next []*tenantTreeNode
		for _, parent := range level {
			for i := 1; i <= tenantTreeFanout; i++ {
				next = append(next, &tenantTreeNode{
					depth:  depth,
					key:    fmt.Sprintf("%s-%d", parent.key, i),
					parent: parent,
				})
			}
		}
		nodes = append(nodes, next...)
		level = next
	}

	return nodes
}

// tenantTreePlanChecks expects every tenant in the tree to be planned with the
// given action, and the shared base role to be left alone.
func tenantTreePlanChecks(
	nodes []*tenantTreeNode,
	action plancheck.ResourceActionType,
) resource.ConfigPlanChecks {
	checks := make([]plancheck.PlanCheck, 0, len(nodes)+1)
	for _, n := range nodes {
		checks = append(checks, plancheck.ExpectResourceAction(n.address(), action))
	}

	if action != plancheck.ResourceActionCreate {
		checks = append(checks,
			plancheck.ExpectResourceAction(tenantTreeRoleAddress, plancheck.ResourceActionNoop))
	}

	return resource.ConfigPlanChecks{PreApply: checks}
}

// tenantTreeStateChecks verifies every tenant in the tree after apply: its own
// attributes, that it uses the shared base role, and that it hangs off the
// right parent -- parent_id, parent_name and parent_subdomain all match the
// parent tenant's id, name and subdomain. Level 1 tenants omit parent_id, so
// it is computed to the caller's (master) tenant: it is checked to be set and
// the same for every level 1 tenant.
//
// These are statecheck checks rather than resource.TestCheck* functions
// because the latter run on the legacy flatmap state shim, which rejects
// for_each instances; see TestAccMorpheusTenantTreeOk.
func tenantTreeStateChecks(
	nodes []*tenantTreeNode,
	name, description, accountName string,
) []statecheck.StateCheck {
	var checks []statecheck.StateCheck
	var firstLevel1 *tenantTreeNode

	for _, n := range nodes {
		addr := n.address()

		checks = append(checks,
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"),
				knownvalue.StringExact(name+"-"+n.key)),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("subdomain"),
				knownvalue.StringExact(name+"-"+n.key)),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("description"),
				knownvalue.StringExact(n.description(description))),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("account_name"),
				knownvalue.StringExact(accountName)),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("master"), knownvalue.Bool(false)),
			statecheck.CompareValuePairs(
				addr, tfjsonpath.New("base_role_id"),
				tenantTreeRoleAddress, tfjsonpath.New("id"),
				compare.ValuesSame()),
		)

		if n.parent == nil {
			checks = append(checks,
				statecheck.ExpectKnownValue(addr, tfjsonpath.New("parent_id"), knownvalue.NotNull()),
				statecheck.ExpectKnownValue(addr, tfjsonpath.New("parent_name"), knownvalue.NotNull()),
			)
			if firstLevel1 == nil {
				firstLevel1 = n
			} else {
				checks = append(checks, statecheck.CompareValuePairs(
					addr, tfjsonpath.New("parent_id"),
					firstLevel1.address(), tfjsonpath.New("parent_id"),
					compare.ValuesSame()))
			}

			continue
		}

		parentAddr := n.parent.address()
		checks = append(checks,
			statecheck.CompareValuePairs(
				addr, tfjsonpath.New("parent_id"),
				parentAddr, tfjsonpath.New("id"),
				compare.ValuesSame()),
			statecheck.CompareValuePairs(
				addr, tfjsonpath.New("parent_name"),
				parentAddr, tfjsonpath.New("name"),
				compare.ValuesSame()),
			statecheck.CompareValuePairs(
				addr, tfjsonpath.New("parent_subdomain"),
				parentAddr, tfjsonpath.New("subdomain"),
				compare.ValuesSame()),
		)
	}

	return checks
}

// TestAccMorpheusTenantTreeOk provisions a three-level tree of tenants under
// the master tenant -- two children per node, 14 tenants -- with for_each and
// cross-level parent_id references, then verifies the hierarchy, that a plan
// after apply is empty, that an update to every tenant is applied in place
// (parent_id is bound on create and must not force replacement on update),
// and finally that the tree destroys cleanly, children before parents.
//
// The tree is destroyed by the test's own final step, not by the framework's
// post-test destroy. terraform-plugin-testing (through v1.16.0) converts the
// state through a legacy flatmap shim in two places: under every
// resource.TestCheck* function, and in the deferred cleanup that decides
// whether to run the post-test destroy. That shim rejects string-indexed
// instances ("for_each is not supported"), so a test whose final state still
// holds the tree fails in that cleanup -- after every step has passed -- and,
// because it fails there, skips the destroy and leaks all 14 tenants. The
// final step therefore removes the tree from the configuration, leaving only
// the role in state for the framework to shim and destroy, and the checks
// throughout use statecheck/plancheck, which read Terraform's JSON state and
// plan directly. Should a step fail before the tree is removed, that same
// shim failure prevents cleanup and the tree is left for the tenant sweeper,
// which removes trees children-first.
//
// Nominating a parent is only permitted for the master tenant, so this test
// requires master-tenant credentials; gate it on the master-tenant capability
// once MORPH-16401 lands. The tenant hierarchy only exists on Morpheus 8.1.0
// and later (constants.TenantParentMinVersion), so it is skipped on older
// appliances.
func TestAccMorpheusTenantTreeOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testhelpers.SkipUnlessApplianceVersionAtLeast(
		context.Background(), t, constants.TenantParentMinVersion,
	)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	// Both the tenant names and subdomains are "<name>-<key>"; name is unique
	// per run and, as it starts with the test name, sweepable.
	name := acctest.RandomWithPrefix(t.Name())
	nodes := tenantTreeNodes()

	// The single tenant-scoped base role shared by every tenant in the tree.
	roleConfig, err := role.RenderRoleTenantConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	const (
		description        = "tenant tree"
		updatedDescription = "updated tenant tree"
		accountName        = "tenant 1337"
		updatedAccountName = "tenant 99999"
	)

	treeConfig := tenantTreeConfig(name, description, accountName)
	updatedTreeConfig := tenantTreeConfig(name, updatedDescription, updatedAccountName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			// Apply: every tenant is created and the hierarchy is as configured.
			{
				Config:            providerConfig + roleConfig + treeConfig,
				ConfigPlanChecks:  tenantTreePlanChecks(nodes, plancheck.ResourceActionCreate),
				ConfigStateChecks: tenantTreeStateChecks(nodes, name, description, accountName),
			},
			// Plan after apply: no drift anywhere in the tree.
			{
				Config:   providerConfig + roleConfig + treeConfig,
				PlanOnly: true,
			},
			// Update every tenant in place. parent_id is unchanged, so no
			// tenant may be replaced.
			{
				Config:            providerConfig + roleConfig + updatedTreeConfig,
				ConfigPlanChecks:  tenantTreePlanChecks(nodes, plancheck.ResourceActionUpdate),
				ConfigStateChecks: tenantTreeStateChecks(nodes, name, updatedDescription, updatedAccountName),
			},
			// Remove the tree, keeping the role: every tenant is destroyed,
			// children before parents (each level depends on the one above
			// through parent_id, and Delete waits for the asynchronous
			// deletion to finish), and the role is untouched. This leaves
			// only the role in state for the framework's post-test destroy;
			// see the test comment.
			{
				Config:           providerConfig + roleConfig,
				ConfigPlanChecks: tenantTreePlanChecks(nodes, plancheck.ResourceActionDestroy),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						tenantTreeRoleAddress, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
		},
	})
}
