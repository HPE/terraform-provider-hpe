// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// --- Removing an optional attribute clears it (MORPH-17990) ------------------

// clearableHostConfig is a host network with the clearable attributes the
// caller names set; the rest are left out. display_name is never set: the
// server fills it in from the name when an update sends one, and that value
// must survive the clears — a value the user never configured is not the
// provider's to touch.
func clearableHostConfig(name string, with ...string) string {
	set := map[string]string{
		"description":    `description    = "to be cleared"`,
		"gateway":        `gateway        = "10.212.0.1"`,
		"dns_primary":    `dns_primary    = "10.212.0.53"`,
		"dns_secondary":  `dns_secondary  = "10.212.0.54"`,
		"search_domains": `search_domains = "example.test"`,
		"pool_id":        `pool_id        = 207`,
		"labels":         `labels         = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]`,
		"tenant_ids":     `tenant_ids     = [` + testhelpers.WhoamiTenantIDRef + `, 2]`,
		"cidr":           `cidr           = "10.212.0.0/24"`,
		// Associations and the VLAN: the host type keeps none of them as its
		// cluster, so all three can go.
		"network_domain_id": `network_domain_id = 57`,
		"zone_pool_id":      `zone_pool_id      = 153047`,
		"vlan_id":           `vlan_id           = 2212`,
	}

	var lines []string
	for _, attr := range with {
		lines = append(lines, "  "+set[attr])
	}

	return testhelpers.ProviderBlock() + testhelpers.WhoamiBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network" "host" {
  name       = %q
  cloud_id   = 17
  group_id   = 1
  type_code  = "host"
  visibility = "private"
%s
}
`, name, strings.Join(lines, "\n"))
}

// The host type edits every one of these in place, so each removal is an
// in-place update that reads back as null: strings and the pool through the
// update request's nullable fields or an empty string, labels as an empty
// list. tenant_ids sent empty revokes every tenant but the owner, whom
// Morpheus always keeps, so it is planned unknown and reads back as the
// owner alone. cidr is Optional without Computed — Terraform plans the null
// on its own, and before this the apply failed with an inconsistent result
// because Update never sent it.
func TestAccMorpheusNetworkRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkDHCP)

	name := acctest.RandomWithPrefix(t.Name())
	renamed := name + "-renamed"
	const rn = "hpe_morpheus_network.host"

	// Everything but cidr is cleared in one step; cidr, the Optional-only
	// case, in its own.
	others := []string{
		"description", "gateway", "dns_primary", "dns_secondary", "search_domains", "pool_id", "labels",
		"network_domain_id", "zone_pool_id", "vlan_id",
	}
	all := append([]string{"cidr", "tenant_ids"}, others...)

	nullAfterApply := func(attrs ...string) []plancheck.PlanCheck {
		checks := []plancheck.PlanCheck{plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate)}
		for _, a := range attrs {
			checks = append(checks, plancheck.ExpectKnownValue(rn, tfjsonpath.New(a), knownvalue.Null()))
		}

		return checks
	}

	noAttrs := func(attrs ...string) resource.TestCheckFunc {
		var checks []resource.TestCheckFunc
		for _, a := range attrs {
			key := a
			if a == "labels" {
				key = a + ".#"
			}

			checks = append(checks, resource.TestCheckNoResourceAttr(rn, key))
		}

		return resource.ComposeAggregateTestCheckFunc(checks...)
	}

	// After the clear only the owner has access.
	ownerOnly := resource.ComposeAggregateTestCheckFunc(
		attr(rn, "tenant_ids.#", "1"),
		resource.TestCheckTypeSetElemAttrPair(rn, "tenant_ids.*", "data.hpe_morpheus_whoami.current", "tenant_id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: clearableHostConfig(name, all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					attr(rn, "description", "to be cleared"),
					attr(rn, "gateway", "10.212.0.1"),
					attr(rn, "dns_primary", "10.212.0.53"),
					attr(rn, "dns_secondary", "10.212.0.54"),
					attr(rn, "search_domains", "example.test"),
					attr(rn, "pool_id", "207"),
					attr(rn, "labels.#", "4"),
					attr(rn, "tenant_ids.#", "2"),
					attr(rn, "cidr", "10.212.0.0/24"),
					attr(rn, "network_domain_id", "57"),
					attr(rn, "zone_pool_id", "153047"),
					attr(rn, "vlan_id", "2212"),
					// The server leaves the display name empty on create.
					resource.TestCheckNoResourceAttr(rn, "display_name"),
				),
			},
			{
				// Everything but cidr removed; the rename makes the server
				// fill in display_name, which the user never configured.
				Config: clearableHostConfig(renamed, "cidr"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: append(nullAfterApply(others...),
						plancheck.ExpectUnknownValue(rn, tfjsonpath.New("tenant_ids"))),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					noAttrs(others...),
					ownerOnly,
					attr(rn, "cidr", "10.212.0.0/24"),
					attr(rn, "display_name", renamed),
				),
			},
			{
				// A plan with nothing changed stays empty: the cleared
				// attributes are no longer recorded as configured, and the
				// server-filled display_name is not the provider's to clear.
				Config:   clearableHostConfig(renamed, "cidr"),
				PlanOnly: true,
			},
			{
				// cidr removed: Optional without Computed, the loud case.
				Config: clearableHostConfig(renamed),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: nullAfterApply("cidr"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					noAttrs("cidr"),
					attr(rn, "display_name", renamed),
				),
			},
		},
	})
}

// A clear is only planned for an attribute the provider recorded as configured
// at its last apply. A network imported into state has no such record — the
// same situation as one last applied by a provider that predates the record —
// so removing an attribute from its configuration changes nothing until an
// apply has recorded the configuration; after that the removal clears.
func TestAccMorpheusNetworkClearNeedsARecordedApply(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkDHCP)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network.host"

	forget := testhelpers.ProviderBlock() + `
removed {
  from = hpe_morpheus_network.host

  lifecycle {
    destroy = false
  }
}
`

	var networkID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: clearableHostConfig(name, "cidr", "description", "labels"),
				Check: resource.ComposeAggregateTestCheckFunc(
					attr(rn, "description", "to be cleared"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[rn]
						if !ok {
							return fmt.Errorf("%s not in state", rn)
						}

						networkID = rs.Primary.ID

						return nil
					},
				),
			},
			{
				// Forget the network without destroying it, then import it:
				// state without private state, so without the record.
				Config: forget,
			},
			{
				Config:             clearableHostConfig(name, "cidr", "labels"),
				ResourceName:       rn,
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return networkID, nil
				},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("imported %d instances, want 1", len(states))
					}

					if got := states[0].Attributes["description"]; got != "to be cleared" {
						return fmt.Errorf("imported description = %q, want the server's value", got)
					}

					return nil
				},
			},
			{
				// description is in state and not in the configuration, but
				// nothing recorded it as configured: not a removal, no plan.
				Config:   clearableHostConfig(name, "cidr", "labels"),
				PlanOnly: true,
			},
			{
				// An apply that records the configuration; description is
				// configured again and search_domains makes it a real update.
				Config: clearableHostConfig(name, "cidr", "labels", "description", "search_domains"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate)},
				},
				Check: attr(rn, "search_domains", "example.test"),
			},
			{
				// Now recorded, the removal clears.
				Config: clearableHostConfig(name, "cidr", "labels", "search_domains"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
					},
				},
				Check: resource.TestCheckNoResourceAttr(rn, "description"),
			},
		},
	})
}

// display_name is the one attribute a clear and the server race over in a
// single update: Morpheus fills an empty display name in from the name when a
// rename is sent. That refill is guarded by the network's *current* display
// name being empty, so a display name the user set — and is now removing —
// is cleared even when the same apply renames the network. (A display name
// the user never set is not recorded and is never cleared; see the test
// above, where the refilled value survives.)
func TestAccMorpheusNetworkRemovedDisplayNameClearsDespiteRename(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkDHCP)

	name := acctest.RandomWithPrefix(t.Name())
	renamed := name + "-renamed"
	const rn = "hpe_morpheus_network.host"

	config := func(name string, displayName bool) string {
		extra := ""
		if displayName {
			extra = `  display_name = "shown"`
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network" "host" {
  name       = %q
  cloud_id   = 17
  group_id   = 1
  type_code  = "host"
  cidr       = "10.212.1.0/24"
  visibility = "private"
%s
}
`, name, extra)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(name, true),
				Check:  attr(rn, "display_name", "shown"),
			},
			{
				// Removed and renamed in the same apply: cleared, not refilled.
				Config: config(renamed, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("display_name"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					attr(rn, "name", renamed),
					resource.TestCheckNoResourceAttr(rn, "display_name"),
				),
			},
			{
				Config:   config(renamed, false),
				PlanOnly: true,
			},
		},
	})
}
