// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildCreateRequest(t *testing.T) {
	t.Parallel()

	var m VirtualSwitchModel
	m.Name = tfstr("vsw1")
	m.BaseType = tfstr("general")
	m.BondMode = tfstr("active-backup")
	m.Mtu = tfint(9000)
	m.Description = tfstr("desc")

	segs := []segmentInput{
		{key: "mgmt", trafficType: trafficManagement},
		{
			key: "nfs", trafficType: trafficDataNfs, vlanID: i32(100), netmask: sp("255.255.255.0"),
			hostIps: []hostIPInput{{serverID: 7, ip: "10.0.0.7"}},
		},
	}
	ups := []uplinkInput{{serverID: 7, nic1: "eth2", nic2: sp("eth3")}}

	body := buildCreateRequest(m, segs, ups)

	if body.VirtualSwitch.Name != "vsw1" {
		t.Fatalf("name: %s", body.VirtualSwitch.Name)
	}

	if body.VirtualSwitch.BondMode == nil || *body.VirtualSwitch.BondMode != "active-backup" {
		t.Fatal("bond_mode not set")
	}

	if len(body.VirtualSwitch.Segments) != 2 {
		t.Fatalf("segments: %d", len(body.VirtualSwitch.Segments))
	}

	nfs := body.VirtualSwitch.Segments[1]
	if nfs.VlanId == nil || *nfs.VlanId != 100 {
		t.Fatal("nfs vlan not set")
	}

	if len(nfs.HostIps) != 1 || nfs.HostIps[0].ServerId == nil || *nfs.HostIps[0].ServerId != 7 {
		t.Fatal("host ips not set")
	}

	if len(body.VirtualSwitch.Uplinks) != 1 || body.VirtualSwitch.Uplinks[0].Nic1 != "eth2" {
		t.Fatal("uplink not set")
	}

	if !body.VirtualSwitch.Uplinks[0].Nic2.IsSet() {
		t.Fatal("nic2 should be set")
	}
}

func TestBuildCompositeEditClearsBond(t *testing.T) {
	t.Parallel()

	p := updatePlan{
		bondModeChanged: true,
		bondMode:        sp(bondModeNone),
		mtuChanged:      true,
		mtu:             i32(1500),
		uplinksChanged:  true,
		uplinks:         []uplinkInput{{serverID: 1, nic1: "eth2"}},
		newSegments:     []segmentInput{{trafficType: trafficVM}},
		editedSegments: []editedSegment{
			{segID: 9, segment: segmentInput{trafficType: trafficDataNfs, vlanID: i32(200)}},
		},
	}

	body := buildCompositeEdit(p)

	if body.SwitchUpdate == nil || body.SwitchUpdate.BondMode == nil || *body.SwitchUpdate.BondMode != "none" {
		t.Fatal("bond_mode should be 'none'")
	}

	if len(body.NewSegments) != 1 || len(body.EditedSegments) != 1 {
		t.Fatal("segment lists wrong")
	}

	if body.EditedSegments[0].SegId == nil || *body.EditedSegments[0].SegId != 9 {
		t.Fatal("editedSegment segId missing")
	}

	if v := body.EditedSegments[0].VlanId.Get(); v == nil || *v != 200 {
		t.Fatal("editedSegment vlan missing")
	}

	// A cleared host_ips must reach the wire as "hostIps": [] (not omitted).
	raw, err := json.Marshal(body.EditedSegments[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"hostIps":[]`) {
		t.Fatalf("editedSegment without host_ips should send an empty hostIps list: %s", raw)
	}
}

func TestBuildRenameAndDescription(t *testing.T) {
	t.Parallel()

	if buildRename("x").Name != "x" {
		t.Fatal("rename")
	}

	body := buildDescriptionUpdate(sp("hi"))
	if body.VirtualSwitch == nil || body.VirtualSwitch.Description == nil || *body.VirtualSwitch.Description != "hi" {
		t.Fatal("description update")
	}

	// A cleared description must reach the wire as an explicit null.
	raw, err := json.Marshal(buildDescriptionUpdate(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"description":null`) {
		t.Fatalf("cleared description should send an explicit null: %s", raw)
	}
}
