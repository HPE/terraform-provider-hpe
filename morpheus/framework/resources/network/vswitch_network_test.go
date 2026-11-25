// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/virtualswitch"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// vswitchNetworksFixture is an HVM cluster on which a virtual switch may be
// created: the network fixture plus the NICs the switch's uplinks take, as
// the virtual switch tests read them.
type vswitchNetworksFixture struct {
	hvmNetworkFixture
	nic1, nic2 string
}

func loadVswitchNetworksFixture(t *testing.T) vswitchNetworksFixture {
	t.Helper()
	f := vswitchNetworksFixture{
		hvmNetworkFixture: loadHVMNetworkFixture(t),
		nic1:              os.Getenv("TF_ACC_VSW_NIC1"),
		nic2:              os.Getenv("TF_ACC_VSW_NIC2"),
	}
	if f.nic1 == "" {
		t.Skip("TF_ACC_VSW_NIC1 must be set")
	}

	return f
}

const (
	// nfsSegmentVlan is the VLAN of the switch's data-nfs segment; a network
	// on the switch may not use it.
	nfsSegmentVlan = 100
	// vswitchMTU is the switch's MTU; changing it while networks are bound is
	// what Morpheus refuses.
	vswitchMTU = 9000
)

// vswitchBlock is a virtual switch on the fixture's cluster with a vm segment
// and a tagged data-nfs segment, the shape the virtual switch resource's own
// test creates.
func (f vswitchNetworksFixture) vswitchBlock(name string, mtu int) string {
	bond := ""
	nic2 := ""
	if f.nic2 != "" {
		bond = `bond_mode  = "active-backup"`
		nic2 = fmt.Sprintf("nic2      = %q", f.nic2)
	}

	return fmt.Sprintf(`
data "hpe_morpheus_cluster_available_nics" "nics" {
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

# Each host's nfs address is its position among the hosts by name: small, and
# stable across the order available_nics happens to list them in.
locals {
  host_names = sort([for host in data.hpe_morpheus_cluster_available_nics.nics.available_nics : host.server_name])
}

resource "hpe_morpheus_virtual_switch" "vs" {
  cluster_id = data.hpe_morpheus_cluster.hvm.id
  name       = %q
  %s
  mtu        = %d

  uplinks = [
    for host in data.hpe_morpheus_cluster_available_nics.nics.available_nics : {
      server_id = host.server_id
      nic1      = %q
      %s
    }
  ]

  segments = {
    vm = {
      traffic_type = "vm"
    }
    nfs = {
      traffic_type = "data-nfs"
      vlan_id      = %d
      netmask      = "255.255.255.0"
      host_ips = [
        for host in data.hpe_morpheus_cluster_available_nics.nics.available_nics : {
          server_id = host.server_id
          ip        = cidrhost("10.20.0.0/24", 10 + index(local.host_names, host.server_name))
        }
      ]
    }
  }
}
`, name, bond, mtu, f.nic1, nic2, nfsSegmentVlan)
}

// switchRefs are the HCL expressions the networks use for the switch's name
// and the nfs segment's id: normally the resource's, and in the step that
// drops the switch from the configuration a data source's, so the networks
// stay as they are while Terraform tries to destroy the switch.
type switchRefs struct {
	name, nfsSegmentID string
}

var (
	resourceSwitchRefs = switchRefs{
		name:         "hpe_morpheus_virtual_switch.vs.name",
		nfsSegmentID: `hpe_morpheus_virtual_switch.vs.segments["nfs"].id`,
	}
	dataSwitchRefs = switchRefs{
		name:         "data.hpe_morpheus_virtual_switch.vs.name",
		nfsSegmentID: `one([for s in data.hpe_morpheus_virtual_switch.vs.segments : s.id if s.traffic_type == "data-nfs"])`,
	}
)

// vswitchDataBlock reads the test's switch by cluster and name.
func vswitchDataBlock(name string) string {
	return fmt.Sprintf(`
data "hpe_morpheus_virtual_switch" "vs" {
  cluster_id = data.hpe_morpheus_cluster.hvm.id
  name       = %q
}
`, name)
}

// standardOnVswitch is an ovsPortGroup network bound to the test's switch,
// with a VLAN and trunks that may or may not collide with the nfs segment.
// switch_id is the switch's libvirt network name, which only the
// hvm-vswitch-or-router option source exposes, looked up by the switch's name.
func (f vswitchNetworksFixture) standardOnVswitch(label, name string, vlan int, trunks string, refs switchRefs) string {
	return fmt.Sprintf(`
data "hpe_morpheus_option_source_value" %[1]q {
  source       = "hvm-vswitch-or-router"
  service_type = "com.hpe.morpheus.hvm"
  cloud_id     = %[3]s
  name         = %[11]s
  parameters = {
    "network.zonePool.id" = tostring(%[5]s)
  }
}

resource "hpe_morpheus_network" %[1]q {
  name              = %[2]q
  cloud_id          = %[3]s
  group_id          = %[4]s
  type_code         = "ovsPortGroup"
  switch_id         = data.hpe_morpheus_option_source_value.%[1]s.value
  zone_pool_id      = %[5]s
  network_server_id = %[6]s
  vlan_id           = %[7]d
  cidr              = "10.209.%[8]d.0/24"
  visibility        = "private"
  tenant_ids        = [%[9]s]
  labels            = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]

  config_ovs_port_group = {
    vlan_trunks = %[10]q
    trunk_only  = false
  }

  lifecycle {
    ignore_changes = [display_name]
  }
}
`, label, name, hvmCloudIDRef, hvmGroupIDRef, hvmPoolIDRef, hvmNetworkServerIDRef, vlan, vlan%256,
		testhelpers.WhoamiTenantIDRef, trunks, refs.name)
}

// dataOnVswitchSegment is an hvmPortGroupDataNetwork on the switch's nfs
// segment, the VSWITCH layout.
func (f vswitchNetworksFixture) dataOnVswitchSegment(name string, refs switchRefs) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_network" "nfs" {
  name              = %q
  cloud_id          = %s
  group_id          = %s
  type_code         = "hvmPortGroupDataNetwork"
  zone_pool_id      = %s
  network_server_id = %s
  visibility        = "private"
  tenant_ids        = [%s]
  labels            = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]

  config_hvm_data = {
    cluster_network_layout    = "VSWITCH"
    virtual_switch_segment_id = %s
  }

  lifecycle {
    ignore_changes = [display_name]
  }
}
`, name, hvmCloudIDRef, hvmGroupIDRef, hvmPoolIDRef, hvmNetworkServerIDRef, testhelpers.WhoamiTenantIDRef, refs.nfsSegmentID)
}

// TestAccMorpheusNetworkOnVirtualSwitch creates a virtual switch and binds a
// standard network (switch_id) and a data network (virtual_switch_segment_id)
// to it, then exercises what Morpheus enforces across the two resources: a
// network VLAN colliding with a segment VLAN on the switch is refused, as are
// changing the switch's MTU and deleting the switch while networks are bound
// to it. The references give network-then-switch destroy order, which is
// what lets the final destroy succeed.
func TestAccMorpheusNetworkOnVirtualSwitch(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.NetworkServer, capabilities.VirtualSwitch)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	f := loadVswitchNetworksFixture(t)

	name := acctest.RandomWithPrefix(t.Name())
	// Switch names are unique per cluster and limited to 16 characters; the
	// prefix is what the virtual switch sweeper matches.
	switchName := virtualswitch.TestNamePrefix + acctest.RandString(5)
	const (
		vs  = "hpe_morpheus_virtual_switch.vs"
		vm  = "hpe_morpheus_network.vm"
		nfs = "hpe_morpheus_network.nfs"
	)

	base := func(mtu int) string {
		return f.config(name) + f.vswitchBlock(switchName, mtu) +
			f.standardOnVswitch("vm", name+"-vm", 2113, "200,201", resourceSwitchRefs) +
			f.dataOnVswitchSegment(name+"-nfs", resourceSwitchRefs)
	}
	// The same networks with the switch resource gone from the configuration,
	// referencing the switch through a data source instead: Terraform plans
	// only the switch's destroy, which Morpheus refuses.
	withoutSwitch := f.config(name) + vswitchDataBlock(switchName) +
		f.standardOnVswitch("vm", name+"-vm", 2113, "200,201", dataSwitchRefs) +
		f.dataOnVswitchSegment(name+"-nfs", dataSwitchRefs)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: base(vswitchMTU),
				Check: resource.ComposeTestCheckFunc(
					attr(vm, "type_code", "ovsPortGroup"),
					// switch_id is the switch's libvirt network name.
					resource.TestMatchResourceAttr(vm, "switch_id", regexp.MustCompile(`^hvmnet\d+$`)),
					attr(vm, "config_from_api.hvmVlanTrunks", "200,201"),
					attr(nfs, "type_code", "hvmPortGroupDataNetwork"),
					resource.TestCheckResourceAttrPair(nfs, "config_from_api.virtualSwitchSegmentId", vs, "segments.nfs.id"),
					// The plugin records the switch on both networks, which is
					// what makes Morpheus refuse the switch's delete and MTU change.
					resource.TestCheckResourceAttrPair(nfs, "config_from_api.virtualSwitchId", vs, "id"),
					resource.TestCheckResourceAttrPair(vm, "config_from_api.virtualSwitchId", vs, "id"),
				),
			},
			{
				Config:             base(vswitchMTU),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// A second standard network whose VLAN is the nfs segment's.
				Config:      base(vswitchMTU) + f.standardOnVswitch("clash", name+"-clash", nfsSegmentVlan, "", resourceSwitchRefs),
				ExpectError: regexp.MustCompile(`(?is)in\s+use\s+by\s+a\s+segment`),
			},
			{
				// A trunk VLAN colliding with the segment's is refused too.
				Config: base(vswitchMTU) +
					f.standardOnVswitch("clash", name+"-clash", 2114, fmt.Sprintf("%d,300", nfsSegmentVlan), resourceSwitchRefs),
				ExpectError: regexp.MustCompile(`(?is)in\s+use\s+by\s+a\s+segment`),
			},
			{
				// The switch's MTU cannot change while networks are bound to it.
				Config:      base(1500),
				ExpectError: regexp.MustCompile(`(?is)cannot\s+change\s+the\s+MTU`),
			},
			{
				// Nor can the switch be deleted while they are.
				Config:      withoutSwitch,
				ExpectError: regexp.MustCompile(`(?is)cannot\s+remove\s+virtual\s+switches\s+that\s+have\s+active\s+networks`),
			},
			{
				// Back to the applied configuration, as an apply rather than a
				// plan: the previous step rewrote the networks' recorded
				// dependencies to the data source, and only an apply writes
				// the resource references back to state. The test's destroy
				// then removes the networks before the switch, as Morpheus
				// requires. Nothing changes on the appliance here.
				Config: base(vswitchMTU),
			},
		},
	})
}
