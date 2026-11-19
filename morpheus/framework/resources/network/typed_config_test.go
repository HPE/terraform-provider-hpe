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
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// The typed configuration blocks are write-only, so these tests exercise what
// state cannot show: the block reaches Morpheus (config_from_api reflects it),
// an editable change is applied in place, a change Morpheus does not apply
// after create is planned as a replacement, and an imported network carries
// no configuration in state but plans cleanly once the block is applied.
//
// Each test targets one network type on the appliance's fixtures.

// typedLifecycle runs the shared step sequence for one typed block.
//
//   - create: the initial configuration; checks run after apply.
//   - update: the same network with an editable value changed, applied in
//     place; nil when the type has no editable configuration.
//   - replace: the same network with an immutable value changed, applied —
//     the testing framework only checks the planned action on a step it
//     applies — so the value must be a real alternative on the appliance; nil
//     when the fixtures offer none.
//   - importIgnore: attributes Morpheus does not return, kept from the
//     configuration, so an imported state cannot carry them.
type typedLifecycle struct {
	resourceName string
	create       string
	createChecks []resource.TestCheckFunc
	update       string
	updateChecks []resource.TestCheckFunc
	replace      string
	importIgnore []string
}

func (l typedLifecycle) steps() []resource.TestStep {
	steps := []resource.TestStep{
		{
			Config: l.create,
			Check:  resource.ComposeTestCheckFunc(l.createChecks...),
		},
		{
			// Nothing changed: the write-only block must not plan anything.
			Config:             l.create,
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
	}

	if l.update != "" {
		steps = append(steps,
			resource.TestStep{
				Config: l.update,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(l.resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(l.updateChecks...),
			},
			resource.TestStep{
				Config:             l.update,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		)
	}

	current := l.create
	if l.update != "" {
		current = l.update
	}

	if l.replace != "" {
		current = l.replace

		steps = append(steps,
			resource.TestStep{
				Config: l.replace,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(l.resourceName, plancheck.ResourceActionReplace),
					},
				},
			},
			resource.TestStep{
				Config:             l.replace,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		)
	}

	steps = append(steps,
		resource.TestStep{
			Config:                  current,
			ResourceName:            l.resourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: l.importIgnore,
		},
		resource.TestStep{
			// After import the private-state digests are empty, so the first
			// plan re-applies the block once; the plan after that is clean.
			Config:             current,
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
	)

	return steps
}

func attr(name, key, value string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(name, key, value)
}

// --- HVM networks (OVS port group, HVM data network) -----------------------

// hvmNetworkFixture is the HVM fixture (testhelpers/hvm_fixture.go,
// testhelpers/network_server_fixture.go) the typed HVM network tests resolve by
// name: the cloud, the cluster whose resource pool the network belongs to,
// the cluster's network server, two of its virtual switches and a host
// interface, plus the caller's tenant and a group of the test's own. Every id
// in the network configuration comes from these.
type hvmNetworkFixture struct {
	cloudName, clusterName, networkServerName string
	switchName, switchName2                   string
	// hostInterface is a host interface for the HOST_INTERFACE layout, which
	// only a cluster without virtual switches (HVM layout before 2.0) accepts;
	// it is read only by the test that needs it.
	hostInterface string
}

// loadHVMNetworkFixture reads the fixture from the environment, skipping the
// test for the first variable that is unset.
func loadHVMNetworkFixture(t *testing.T) hvmNetworkFixture {
	t.Helper()

	f := hvmNetworkFixture{
		cloudName:         testhelpers.HVMCloudName(t),
		clusterName:       testhelpers.HVMClusterName(t),
		networkServerName: testhelpers.HVMNetworkServerName(t),
	}
	f.switchName, f.switchName2 = testhelpers.HVMSwitchNames(t)

	return f
}

const (
	hvmCloudIDRef         = "data.hpe_morpheus_cloud.hvm_cloud.id"
	hvmPoolIDRef          = "data.hpe_morpheus_cluster.hvm.permissions.resource_pool.id"
	hvmNetworkServerIDRef = "data.hpe_morpheus_network_server.hvm.id"
	hvmGroupIDRef         = "hpe_morpheus_group.hvm_network.id"
)

// config renders the fixture's data sources and the test's group, to go
// before a network resource. groupName is the test's name.
func (f hvmNetworkFixture) config(groupName string) string {
	return testhelpers.ProviderBlock() + testhelpers.WhoamiBlock() + fmt.Sprintf(`
data "hpe_morpheus_cloud" "hvm_cloud" {
  name = %s
}

data "hpe_morpheus_cluster" "hvm" {
  name = %s
}

resource "hpe_morpheus_group" "hvm_network" {
  name = %q
  code = %q
}

`, testhelpers.HCLQuote(f.cloudName), testhelpers.HCLQuote(f.clusterName), groupName, strings.ToLower(groupName)) +
		testhelpers.NetworkServerBlock("hvm", f.networkServerName) +
		testhelpers.HVMSwitchBlock("switch", f.switchName, hvmCloudIDRef, hvmPoolIDRef) +
		testhelpers.HVMSwitchBlock("switch_2", f.switchName2, hvmCloudIDRef, hvmPoolIDRef)
}

// ovsPortGroupConfig is an OVS port group on one of the cluster's switches,
// switchRef being the option-source data source for it.
func (f hvmNetworkFixture) ovsPortGroupConfig(name, trunks, switchRef string) string {
	return f.config(name) + fmt.Sprintf(`
resource "hpe_morpheus_network" "ovs" {
  name              = %q
  cloud_id          = %s
  group_id          = %s
  type_code         = "ovsPortGroup"
  switch_id         = %s.value
  zone_pool_id      = %s
  network_server_id = %s
  vlan_id           = 2103
  cidr              = "10.209.103.0/24"
  visibility        = "private"
  tenant_ids        = [%s]
  labels            = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]

  config_ovs_port_group = {
    vlan_trunks = %q
    trunk_only  = false
  }

  lifecycle {
    ignore_changes = [display_name]
  }
}
`, name, hvmCloudIDRef, hvmGroupIDRef, switchRef, hvmPoolIDRef, hvmNetworkServerIDRef,
		testhelpers.WhoamiTenantIDRef, trunks)
}

func TestAccMorpheusNetworkTypedOVSPortGroupLifecycle(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.NetworkServer)

	f := loadHVMNetworkFixture(t)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network.ovs"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: typedLifecycle{
			resourceName: rn,
			create:       f.ovsPortGroupConfig(name, "100,200", "data.hpe_morpheus_option_source_value.switch"),
			createChecks: []resource.TestCheckFunc{
				// The type's id is per appliance; that it is populated is what matters.
				resource.TestCheckResourceAttrSet(rn, "type_id"),
				attr(rn, "type_code", "ovsPortGroup"),
				attr(rn, "config_from_api.hvmVlanTrunks", "100,200"),
				attr(rn, "config_from_api.hvmTrunkOnly", "false"),
				resource.TestCheckNoResourceAttr(rn, "config_ovs_port_group"),
			},
			update: f.ovsPortGroupConfig(name, "100,200,300", "data.hpe_morpheus_option_source_value.switch"),
			updateChecks: []resource.TestCheckFunc{
				attr(rn, "config_from_api.hvmVlanTrunks", "100,200,300"),
			},
			// switchId is an option type Morpheus does not edit after create;
			// the cluster's other switch is a real alternative.
			replace:      f.ovsPortGroupConfig(name, "100,200,300", "data.hpe_morpheus_option_source_value.switch_2"),
			importIgnore: []string{"network_server_id"},
		}.steps(),
	})
}

// hvmDataConfig is an HVM data network on one of the hosts' interfaces.
func (f hvmNetworkFixture) hvmDataConfig(name string, vlan int) string {
	return f.config(name) + fmt.Sprintf(`
resource "hpe_morpheus_network" "hvm_data" {
  name              = %q
  cloud_id          = %s
  group_id          = %s
  type_code         = "hvmPortGroupDataNetwork"
  zone_pool_id      = %s
  network_server_id = %s
  vlan_id           = %d
  visibility        = "private"
  tenant_ids        = [%s]
  labels            = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]

  config_hvm_data = {
    cluster_network_layout = "HOST_INTERFACE"
    host_network_interface = %s
    mtu                    = 1500
  }

  lifecycle {
    ignore_changes = [display_name]
  }
}
`, name, hvmCloudIDRef, hvmGroupIDRef, hvmPoolIDRef, hvmNetworkServerIDRef, vlan,
		testhelpers.WhoamiTenantIDRef, testhelpers.HCLQuote(f.hostInterface))
}

func TestAccMorpheusNetworkTypedHVMDataNetwork(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.NetworkServer)

	f := loadHVMNetworkFixture(t)
	// The HOST_INTERFACE layout is for a cluster without virtual switches; on
	// an HVM 2.0 cluster Morpheus requires a segment instead, so leave the
	// interface unset there.
	f.hostInterface = testhelpers.HVMHostInterface(t)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network.hvm_data"

	// The MTU is the type's one editable field; it is left alone rather than
	// put jumbo frames on the shared cluster. vlanIdEditable is off for the
	// type, so a VLAN change is the catalogue-driven replacement.
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: typedLifecycle{
			resourceName: rn,
			create:       f.hvmDataConfig(name, 2104),
			createChecks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttrSet(rn, "type_id"),
				attr(rn, "type_code", "hvmPortGroupDataNetwork"),
				attr(rn, "config_from_api.hvmHostNetworkInterface", f.hostInterface),
				attr(rn, "config_from_api.hvmMtu", "1500"),
			},
			replace:      f.hvmDataConfig(name, 2105),
			importIgnore: []string{"network_server_id"},
		}.steps(),
	})
}

// --- NSX-T segment ----------------------------------------------------------

func nsxtConfig(name, gateway string) string {
	return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network" "nsxt" {
  name              = %q
  cloud_id          = 45
  group_id          = 1
  type_code         = "nsxtLogicalSwitch"
  zone_pool_id      = 281
  network_server_id = 16
  cidr              = "10.32.201.1/24"
  dhcp_server       = false
  visibility        = "private"
  tenant_ids        = [1]
  labels            = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]

  config_nsxt = {
    connected_gateway = %q
  }

  lifecycle {
    ignore_changes = [display_name]
  }
}
`, name, gateway)
}

func TestAccMorpheusNetworkTypedNSXTSegment(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT, capabilities.NetworkServer)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network.nsxt"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: typedLifecycle{
			resourceName: rn,
			create:       nsxtConfig(name, "/infra/tier-1s/qa-t1-01"),
			createChecks: []resource.TestCheckFunc{
				attr(rn, "type_code", "nsxtLogicalSwitch"),
				attr(rn, "config_from_api.connectedGateway", "/infra/tier-1s/qa-t1-01"),
				// The deprecated top-level attribute mirrors the block's value.
				attr(rn, "connected_gateway", "/infra/tier-1s/qa-t1-01"),
			},
			update: nsxtConfig(name, "/infra/tier-1s/dd-test"),
			updateChecks: []resource.TestCheckFunc{
				attr(rn, "config_from_api.connectedGateway", "/infra/tier-1s/dd-test"),
			},
			importIgnore: []string{"network_server_id"},
		}.steps(),
	})
}

// --- Custom child network (VMware parent) ---------------------------------

func childNetworkConfig(name, cidr, gateway string) string {
	return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network" "child" {
  name              = %q
  cloud_id          = 2
  group_id          = 1
  type_code         = "childNetwork"
  parent_network_id = 86657
  cidr              = %q
  gateway           = %q
  visibility        = "private"
  tenant_ids        = [1]
  labels            = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
}
`, name, cidr, gateway)
}

func TestAccMorpheusNetworkTypedChildNetwork(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network.child"

	// childNetwork has nameEditable and cidrEditable on, so both a rename and
	// a re-addressing are in-place updates — the rename is the one case the
	// static RequiresReplace on name used to hide, and the update request has
	// to carry the name for it to happen.
	steps := typedLifecycle{
		resourceName: rn,
		create:       childNetworkConfig(name, "10.77.0.0/24", "10.77.0.1"),
		createChecks: []resource.TestCheckFunc{
			attr(rn, "type_code", "childNetwork"),
			attr(rn, "parent_network_id", "86657"),
			attr(rn, "cidr", "10.77.0.0/24"),
		},
		update: childNetworkConfig(name+"-renamed", "10.78.0.0/24", "10.78.0.1"),
		updateChecks: []resource.TestCheckFunc{
			attr(rn, "name", name+"-renamed"),
			attr(rn, "cidr", "10.78.0.0/24"),
		},
		importIgnore: []string{"parent_network_id"},
	}.steps()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// --- OpenStack private network ---------------------------------------------

func openstackConfig(name, dhcpRange, hostRoutes string) string {
	return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network" "openstack" {
  name              = %q
  cloud_id          = 5472
  group_id          = 1
  type_code         = "openstackPrivate"
  zone_pool_id      = 46930
  network_server_id = 2294
  cidr              = "192.168.177.0/24"
  gateway           = "192.168.177.1"
  dhcp_server       = true
  visibility        = "private"
  tenant_ids        = [1]
  labels            = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]

  config_openstack = {
    dhcp_range  = %q
    host_routes = %q
  }

  lifecycle {
    ignore_changes = [display_name]
  }
}
`, name, dhcpRange, hostRoutes)
}

func TestAccMorpheusNetworkTypedOpenStackPrivate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.OpenStack, capabilities.NetworkServer)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network.openstack"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: typedLifecycle{
			resourceName: rn,
			create:       openstackConfig(name, "192.168.177.10,192.168.177.100", ""),
			createChecks: []resource.TestCheckFunc{
				attr(rn, "type_code", "openstackPrivate"),
				attr(rn, "config_from_api.dhcpRange", "192.168.177.10,192.168.177.100"),
			},
			update: openstackConfig(name, "192.168.177.20,192.168.177.100", ""),
			updateChecks: []resource.TestCheckFunc{
				attr(rn, "config_from_api.dhcpRange", "192.168.177.20,192.168.177.100"),
			},
			// hostRoutes is not editable after create.
			replace:      openstackConfig(name, "192.168.177.20,192.168.177.100", "10.0.0.0/8,192.168.177.254"),
			importIgnore: []string{"network_server_id"},
		}.steps(),
	})
}

// --- Azure virtual network ---------------------------------------------------

func azureConfig(name, subnetCIDR string) string {
	return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network" "azure" {
  name       = %q
  cloud_id   = 4617
  group_id   = 1
  type_code  = "azureNetwork"
  cidr       = "10.61.0.0/16"
  visibility = "private"
  tenant_ids = [1]
  labels     = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]

  config_azure = {
    resource_group_id = "morph-qa"
    subnet_name       = "%s-subnet"
    subnet_cidr       = %q
  }
}
`, name, name, subnetCIDR)
}

func TestAccMorpheusNetworkTypedAzure(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Azure)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network.azure"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: typedLifecycle{
			resourceName: rn,
			create:       azureConfig(name, "10.61.1.0/24"),
			createChecks: []resource.TestCheckFunc{
				attr(rn, "type_id", "35"),
				attr(rn, "type_code", "azureNetwork"),
				attr(rn, "config_from_api.resourceGroupId", "morph-qa"),
				attr(rn, "config_from_api.subnetCidr", "10.61.1.0/24"),
			},
			// All three Azure fields are immutable.
			replace: azureConfig(name, "10.61.2.0/24"),
		}.steps(),
	})
}

// --- tenant_ids unset (MORPH-17908) ------------------------------------------

func hostNetworkNoTenantsConfig(name, description string) string {
	return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network" "host" {
  name        = %q
  description = %q
  cloud_id    = 17
  group_id    = 1
  type_code   = "host"
  cidr        = "10.0.0.0/8"
  visibility  = "private"
  labels      = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
}
`, name, description)
}

// With tenant_ids not configured, a change to anything else used to plan a
// replacement: the attribute is Optional+Computed, so it planned as unknown,
// and its unconditional RequiresReplace compared unknown with state. It now
// plans an in-place update, and — with UseStateForUnknown — leaves tenant_ids
// out of the plan entirely (MORPH-17908).
func TestAccMorpheusNetworkTenantIdsUnsetUpdatesInPlace(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkDHCP)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network.host"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: hostNetworkNoTenantsConfig(name, "first"),
				Check: resource.ComposeTestCheckFunc(
					// No tenants given, none granted: the attribute is null in state.
					resource.TestCheckNoResourceAttr(rn, "tenant_ids.#"),
					attr(rn, "description", "first"),
				),
			},
			{
				Config: hostNetworkNoTenantsConfig(name, "second"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						// Not "(known after apply)": the prior (null) value is kept.
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("tenant_ids"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					attr(rn, "description", "second"),
					resource.TestCheckNoResourceAttr(rn, "tenant_ids.#"),
				),
			},
		},
	})
}

// --- HPE Morpheus SDN virtual network (juniperSdn) --------------------------

func juniperSdnNetworkConfig(name, cloudID, groupID, networkServerID, zonePoolID, cidr, gateway string) string {
	return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network" "sdn" {
  name              = %q
  cloud_id          = %s
  group_id          = %s
  type_code         = "juniperSdn"
  network_server_id = %s
  zone_pool_id      = %s
  cidr              = %q
  gateway           = %q
  visibility        = "private"
  tenant_ids        = [1]
  labels            = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]

  lifecycle {
    ignore_changes = [display_name]
  }
}
`, name, cloudID, groupID, networkServerID, zonePoolID, cidr, gateway)
}

// TestAccMorpheusNetworkTypedJuniperSdn exercises the SDN plugin's juniperSdn
// virtual network type: a network scoped to an HVM cluster resource pool
// (zonePool.id) on a CN2 controller. Only the gateway is editable after
// create, so the update changes it; CIDR and the resource pool are immutable.
func TestAccMorpheusNetworkTypedJuniperSdn(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)
	f := testhelpers.RequireSdnFixture(t)

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network.sdn"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: typedLifecycle{
			resourceName: rn,
			create: juniperSdnNetworkConfig(
				name, f.CloudID, f.GroupID, f.NetworkServerID, f.ClusterPoolID,
				"10.90.0.0/24", "10.90.0.1"),
			createChecks: []resource.TestCheckFunc{
				attr(rn, "type_code", "juniperSdn"),
				attr(rn, "zone_pool_id", f.ClusterPoolID),
				attr(rn, "cidr", "10.90.0.0/24"),
				attr(rn, "gateway", "10.90.0.1"),
			},
			// The plugin refuses a gateway change after creation (verified live:
			// "Gateway cannot be modified after network creation"), and edits
			// nothing else on this type, so the only lifecycle change is a
			// replacement.
			replace: juniperSdnNetworkConfig(
				name, f.CloudID, f.GroupID, f.NetworkServerID, f.ClusterPoolID,
				"10.90.0.0/24", "10.90.0.2"),
			importIgnore: []string{"network_server_id"},
		}.steps(),
	})
}
