// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import "testing"

func TestAssignSegmentKeys(t *testing.T) {
	t.Parallel()

	t.Run("match on id", func(t *testing.T) {
		t.Parallel()

		rows := []segRow{
			{id: 10, hasID: true, trafficType: "vm"},
			{id: 11, hasID: true, trafficType: "data-nfs", vlanID: 100, hasVlan: true},
		}
		priors := []priorSeg{
			{key: "myvm", id: 10, hasID: true, trafficType: "vm"},
			{key: "nfs_a", id: 11, hasID: true, trafficType: "data-nfs", vlanID: 100, hasVlan: true},
		}

		got := assignSegmentKeys(rows, priors)
		if got[0] != "myvm" || got[1] != "nfs_a" {
			t.Fatalf("got %v", got)
		}
	})

	t.Run("match on traffic_type+vlan when no id (first read after create)", func(t *testing.T) {
		t.Parallel()

		rows := []segRow{
			{trafficType: "vm"},
			{trafficType: "data-nfs", vlanID: 100, hasVlan: true},
		}
		priors := []priorSeg{
			{key: "myvm", trafficType: "vm"},
			{key: "nfs_a", trafficType: "data-nfs", vlanID: 100, hasVlan: true},
		}

		got := assignSegmentKeys(rows, priors)
		if got[0] != "myvm" || got[1] != "nfs_a" {
			t.Fatalf("got %v", got)
		}
	})

	t.Run("synthesise keys on import", func(t *testing.T) {
		t.Parallel()

		rows := []segRow{
			{id: 1, hasID: true, trafficType: "management"},
			{id: 2, hasID: true, trafficType: "data-nfs", vlanID: 100, hasVlan: true},
		}

		got := assignSegmentKeys(rows, nil)
		if got[0] != "management" || got[1] != "data-nfs_100" {
			t.Fatalf("got %v", got)
		}
	})

	t.Run("dedupe synthesised keys", func(t *testing.T) {
		t.Parallel()

		rows := []segRow{
			{id: 1, hasID: true, trafficType: "vm"},
			{id: 2, hasID: true, trafficType: "vm"},
		}

		got := assignSegmentKeys(rows, nil)
		if got[0] != "vm" || got[1] != "vm_2" {
			t.Fatalf("got %v", got)
		}
	})

	// After a re-key following import, the id is in state under the user's new
	// key. The id-match pass must keep that user key rather than re-synthesising
	// the data-nfs_100 convention, so subsequent reads are stable.
	t.Run("keep user key by id after rekey", func(t *testing.T) {
		t.Parallel()

		rows := []segRow{
			{id: 2, hasID: true, trafficType: "data-nfs", vlanID: 100, hasVlan: true},
		}
		priors := []priorSeg{
			{key: "nfs", id: 2, hasID: true, trafficType: "data-nfs", vlanID: 100, hasVlan: true},
		}

		got := assignSegmentKeys(rows, priors)
		if got[0] != "nfs" {
			t.Fatalf("got %v, want [nfs]", got)
		}
	})
}
