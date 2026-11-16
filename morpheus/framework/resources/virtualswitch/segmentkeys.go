// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"fmt"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// segRow is the minimal identity of an API segment row used for key matching.
type segRow struct {
	id          int64
	hasID       bool
	trafficType string
	vlanID      int64
	hasVlan     bool
}

// priorSeg is the minimal identity of a prior-state segment, tied to its key.
type priorSeg struct {
	key         string
	id          int64
	hasID       bool
	trafficType string
	vlanID      int64
	hasVlan     bool
}

// carrySegmentKeys decides the map key for each API segment row, in row order.
//
// The map key lives only in Terraform; Morpheus does not know it. Keys are
// carried forward by matching prior-state entries to API rows:
//   - first on the computed segment id (stable once the switch exists);
//   - then, for rows with no id yet (the first Read after create, where create
//     returns no ids), on (traffic_type, vlan_id) — unique per switch by the
//     server's own rules.
//
// Rows that match no prior key get a synthesised key (import case):
// "<traffic_type>" for an untagged segment, "<traffic_type>_<vlan_id>" for a
// tagged one, de-duplicated with a numeric suffix if it collides.
func carrySegmentKeys(
	ctx context.Context,
	rows []sdk.GetVirtualSwitch200ResponseVirtualSwitchAllOfSegmentsInner1,
	prior VirtualSwitchModel,
) []string {
	apiRows := make([]segRow, len(rows))

	for i, row := range rows {
		r := segRow{}
		if row.Id != nil {
			r.id = *row.Id
			r.hasID = true
		}

		if row.TrafficType != nil {
			r.trafficType = *row.TrafficType
		}

		if v := row.VlanId.Get(); v != nil {
			r.vlanID = int64(*v)
			r.hasVlan = true
		}

		apiRows[i] = r
	}

	return assignSegmentKeys(apiRows, priorSegsFromState(ctx, prior))
}

// priorSegsFromState extracts the (key, id, traffic_type, vlan_id) of each
// prior-state segment.
func priorSegsFromState(ctx context.Context, prior VirtualSwitchModel) []priorSeg {
	if prior.Segments.IsNull() || prior.Segments.IsUnknown() {
		return nil
	}

	segsMap := make(map[string]SegmentsValue, len(prior.Segments.Elements()))
	if d := prior.Segments.ElementsAs(ctx, &segsMap, false); d.HasError() {
		return nil
	}

	out := make([]priorSeg, 0, len(segsMap))

	for key, seg := range segsMap {
		p := priorSeg{key: key}
		if !seg.Id.IsNull() && !seg.Id.IsUnknown() {
			p.id = seg.Id.ValueInt64()
			p.hasID = p.id != 0
		}

		if !seg.TrafficType.IsNull() && !seg.TrafficType.IsUnknown() {
			p.trafficType = seg.TrafficType.ValueString()
		}

		if !seg.VlanId.IsNull() && !seg.VlanId.IsUnknown() {
			p.vlanID = seg.VlanId.ValueInt64()
			p.hasVlan = true
		}

		out = append(out, p)
	}

	return out
}

// assignSegmentKeys is the pure key-assignment algorithm, separated for unit
// testing.
func assignSegmentKeys(rows []segRow, priors []priorSeg) []string {
	keys := make([]string, len(rows))
	used := make(map[string]bool)
	priorUsed := make([]bool, len(priors))

	// Pass 1: match on id.
	for i, row := range rows {
		if !row.hasID {
			continue
		}

		for j, p := range priors {
			if priorUsed[j] || !p.hasID {
				continue
			}

			if p.id == row.id {
				keys[i] = p.key
				used[p.key] = true
				priorUsed[j] = true

				break
			}
		}
	}

	// Pass 2: match remaining rows on (traffic_type, vlan_id).
	for i, row := range rows {
		if keys[i] != "" {
			continue
		}

		for j, p := range priors {
			if priorUsed[j] {
				continue
			}

			if p.trafficType == row.trafficType && p.hasVlan == row.hasVlan && p.vlanID == row.vlanID {
				keys[i] = p.key
				used[p.key] = true
				priorUsed[j] = true

				break
			}
		}
	}

	// Pass 3: synthesise keys for unmatched rows.
	for i, row := range rows {
		if keys[i] != "" {
			continue
		}

		keys[i] = dedupeKey(synthKey(row), used)
		used[keys[i]] = true
	}

	return keys
}

func synthKey(row segRow) string {
	if row.hasVlan {
		return fmt.Sprintf("%s_%d", row.trafficType, row.vlanID)
	}

	return row.trafficType
}

func dedupeKey(base string, used map[string]bool) string {
	if !used[base] {
		return base
	}

	for n := 2; ; n++ {
		candidate := fmt.Sprintf("%s_%d", base, n)
		if !used[candidate] {
			return candidate
		}
	}
}
