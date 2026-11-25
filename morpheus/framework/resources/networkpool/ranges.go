// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool

import (
	"context"
	"fmt"
	"net/netip"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// ipRange is one address range of a pool as the provider reasons about it.
// id is 0 when the range does not exist in Morpheus yet.
type ipRange struct {
	id    int64
	start string
	end   string
}

func (r ipRange) key() string { return r.start + "-" + r.end }

// Morpheus reconciles a pool's ranges BY ID on update: a range sent with its
// id is edited in place, a range without an id is created, and a range of the
// pool that is absent from the request is deleted (refused when it still has
// reservations). The provider therefore carries each range's id in state and
// sends it back, so that changing a range's addresses never becomes a
// delete-and-create — which would fail for any range with allocated IPs.

// rangesFromList converts the ip_address_ranges list value into ranges.
// Unknown or null ids become 0.
func rangesFromList(ctx context.Context, list types.List) ([]ipRange, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var values []IpAddressRangesValue
	diags.Append(list.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]ipRange, 0, len(values))
	for _, v := range values {
		r := ipRange{start: v.StartingAddress.ValueString(), end: v.EndingAddress.ValueString()}
		if !v.Id.IsNull() && !v.Id.IsUnknown() {
			r.id = v.Id.ValueInt64()
		}
		out = append(out, r)
	}

	return out, diags
}

// rangesToList converts ranges into the ip_address_ranges list value.
func rangesToList(ctx context.Context, ranges []ipRange) (types.List, diag.Diagnostics) {
	elemType := IpAddressRangesValue{}.Type(ctx)
	if ranges == nil {
		return types.ListNull(elemType), nil
	}

	elems := make([]attr.Value, 0, len(ranges))
	for _, r := range ranges {
		id := types.Int64Null()
		if r.id != 0 {
			id = types.Int64Value(r.id)
		}
		elems = append(elems, IpAddressRangesValue{
			Id:              id,
			StartingAddress: types.StringValue(r.start),
			EndingAddress:   types.StringValue(r.end),
			state:           attr.ValueStateKnown,
		})
	}

	return types.ListValue(elemType, elems)
}

// unmanagedRange is a range of the pool that ip_address_ranges cannot express
// (no start/end address: an IPv6 or CIDR range). It is never shown to the
// practitioner, but it must be sent back on every update with its id and cidr
// so Morpheus keeps it — a range missing from the request is deleted.
type unmanagedRange struct {
	id       int64
	cidr     string
	cidrIPv6 string
}

func (u unmanagedRange) String() string {
	switch {
	case u.cidrIPv6 != "":
		return fmt.Sprintf("range %d (%s)", u.id, u.cidrIPv6)
	case u.cidr != "":
		return fmt.Sprintf("range %d (%s)", u.id, u.cidr)
	}

	return fmt.Sprintf("range %d", u.id)
}

// rangesFromAPI splits a pool's ranges into the start/end ranges that
// ip_address_ranges manages and the unmanaged ones (see unmanagedRange).
func rangesFromAPI(pool *sdk.GetNetworkPool200ResponseNetworkPool) (ranges []ipRange, unmanaged []unmanagedRange) {
	if pool == nil {
		return nil, nil
	}
	for _, r := range pool.IpRanges {
		if r.Id == nil {
			continue // cannot be addressed at all; nothing we can do with it
		}
		if r.StartAddress.Get() == nil || r.EndAddress.Get() == nil {
			u := unmanagedRange{id: *r.Id}
			if c := r.Cidr.Get(); c != nil {
				u.cidr = *c
			}
			if c := r.CidrIPv6.Get(); c != nil {
				u.cidrIPv6 = *c
			}
			unmanaged = append(unmanaged, u)

			continue
		}
		ranges = append(ranges, ipRange{id: *r.Id, start: *r.StartAddress.Get(), end: *r.EndAddress.Get()})
	}

	return ranges, unmanaged
}

// unmanagedToUpdateRequest re-sends unmanaged ranges unchanged so an update
// never deletes them.
func unmanagedToUpdateRequest(unmanaged []unmanagedRange) []sdk.UpdateNetworkPoolRequestNetworkPoolIpRangesInner {
	out := make([]sdk.UpdateNetworkPoolRequestNetworkPoolIpRangesInner, 0, len(unmanaged))
	for _, u := range unmanaged {
		u := u
		item := sdk.UpdateNetworkPoolRequestNetworkPoolIpRangesInner{Id: &u.id}
		if u.cidrIPv6 != "" {
			item.CidrIPv6 = &u.cidrIPv6
		}
		if u.cidr != "" {
			item.AdditionalProperties = map[string]interface{}{"cidr": u.cidr}
		}
		out = append(out, item)
	}

	return out
}

// orderLikePrior returns the API's ranges in the order of prior (matched by id,
// then by addresses), with ranges unknown to prior appended in API order. This
// keeps the configured order stable across reads so a list does not churn.
func orderLikePrior(api, prior []ipRange) []ipRange {
	if len(prior) == 0 {
		// No configured order to follow (import, or a pool managed through
		// the deprecated single-range attribute): a deterministic order, by
		// id, so the result does not depend on how the appliance happens to
		// list them. A configuration written in a different order plans one
		// reorder after import and is then stable.
		sorted := append([]ipRange(nil), api...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].id < sorted[j].id })

		return sorted
	}

	byID := map[int64]int{}
	byKey := map[string]int{}
	for i, r := range api {
		if r.id != 0 {
			byID[r.id] = i
		}
		byKey[r.key()] = i
	}

	used := make([]bool, len(api))
	out := make([]ipRange, 0, len(api))
	for _, p := range prior {
		if i, ok := byID[p.id]; ok && p.id != 0 && !used[i] {
			out = append(out, api[i])
			used[i] = true

			continue
		}
		if i, ok := byKey[p.key()]; ok && !used[i] {
			out = append(out, api[i])
			used[i] = true
		}
	}
	for i, r := range api {
		if !used[i] {
			out = append(out, r)
		}
	}

	return out
}

// withIDsFromPrior decides which existing range each planned range IS:
//
//  1. a planned range whose addresses exactly match a prior range is that
//     range (its id), wherever either sits in its list;
//  2. when the list has the same length as before, a planned range matching
//     nothing by addresses but carrying an id (the nested id plan modifier
//     copies ids by POSITION, so this is only a hint) is an in-place edit of
//     that range — provided the id was not claimed in step 1 AND the planned
//     addresses overlap the prior range's addresses. An edit that shares no
//     address with the range it would edit cannot keep that range's
//     reservations anyway, and a remove+add that keeps the list length
//     (prior [A,B], plan [B′,C]) puts A's id at B′'s position: without the
//     overlap test B′ would silently become A and C would become B;
//  3. anything else is new (id 0) and Morpheus creates it. In particular, when
//     ranges were added or removed, positional hints are ignored: after a
//     removal the hint at a given position names a neighbour, and adopting it
//     would silently rewrite one range into another. An edit combined with a
//     removal therefore deletes and recreates the edited range — and if that
//     range holds reservations Morpheus keeps it and the provider says so, so
//     the practitioner can do the two changes in two applies.
//
// Matching by addresses first is what keeps "remove the first range" from
// turning into "edit the first range to look like the second and delete the
// second", which would break reservations in the second range.
func withIDsFromPrior(planned, prior []ipRange) []ipRange {
	byKey := map[string]int64{}
	byID := map[int64]ipRange{}
	for _, p := range prior {
		if p.id != 0 {
			byKey[p.key()] = p.id
			byID[p.id] = p
		}
	}

	out := make([]ipRange, len(planned))
	claimed := map[int64]bool{}
	pending := []int{}
	for i, r := range planned {
		if id, ok := byKey[r.key()]; ok && !claimed[id] {
			r.id = id
			claimed[id] = true
		} else {
			pending = append(pending, i)
			r.id = planned[i].id // positional hint, resolved below
		}
		out[i] = r
	}
	sameShape := len(planned) == len(prior)
	for _, i := range pending {
		hint := out[i].id
		if hinted, known := byID[hint]; sameShape && hint != 0 && known && !claimed[hint] && overlaps(out[i], hinted) {
			claimed[hint] = true

			continue
		}
		out[i].id = 0
	}

	return out
}

// overlaps reports whether two ranges share at least one address. Addresses
// that do not parse (or mix families) never overlap.
func overlaps(a, b ipRange) bool {
	aStart, err1 := netip.ParseAddr(a.start)
	aEnd, err2 := netip.ParseAddr(a.end)
	bStart, err3 := netip.ParseAddr(b.start)
	bEnd, err4 := netip.ParseAddr(b.end)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return false
	}
	if aStart.Compare(aEnd) > 0 {
		aStart, aEnd = aEnd, aStart
	}
	if bStart.Compare(bEnd) > 0 {
		bStart, bEnd = bEnd, bStart
	}

	return aStart.Compare(bEnd) <= 0 && bStart.Compare(aEnd) <= 0
}

// singleRangeUpdatable decides whether the deprecated ip_ranges attribute may
// update a pool holding the given ranges: only when the pool has at most one
// range and that range is one ip_ranges can express. A pool whose only range
// is CIDR/IPv6 is refused too — ip_ranges would ADD a range beside it and the
// pool would be refused forever after. total is the number of ranges reported
// in the refusal.
func singleRangeUpdatable(managed []ipRange, unmanaged []unmanagedRange) (total int, ok bool) {
	total = len(managed) + len(unmanaged)

	return total, len(unmanaged) == 0 && len(managed) <= 1
}

func rangesToCreateRequest(ranges []ipRange) []sdk.CreateNetworkPoolRequestNetworkPoolIpRangesInner {
	out := make([]sdk.CreateNetworkPoolRequestNetworkPoolIpRangesInner, 0, len(ranges))
	for _, r := range ranges {
		r := r
		out = append(out, sdk.CreateNetworkPoolRequestNetworkPoolIpRangesInner{StartAddress: &r.start, EndAddress: &r.end})
	}

	return out
}

func rangesToUpdateRequest(ranges []ipRange) []sdk.UpdateNetworkPoolRequestNetworkPoolIpRangesInner {
	out := make([]sdk.UpdateNetworkPoolRequestNetworkPoolIpRangesInner, 0, len(ranges))
	for _, r := range ranges {
		r := r
		item := sdk.UpdateNetworkPoolRequestNetworkPoolIpRangesInner{StartAddress: &r.start, EndAddress: &r.end}
		if r.id != 0 {
			item.Id = &r.id
		}
		out = append(out, item)
	}

	return out
}

// singleRangeFromObject reads the deprecated ip_ranges object.
func singleRangeFromObject(obj types.Object) (ipRange, bool) {
	if obj.IsNull() || obj.IsUnknown() {
		return ipRange{}, false
	}
	attrs := obj.Attributes()
	start, _ := attrs["starting_address"].(basetypes.StringValue)
	end, _ := attrs["ending_address"].(basetypes.StringValue)

	return ipRange{start: start.ValueString(), end: end.ValueString()}, true
}

func singleRangeToObject(r *ipRange) types.Object {
	attrTypes := map[string]attr.Type{"starting_address": types.StringType, "ending_address": types.StringType}
	if r == nil {
		return types.ObjectNull(attrTypes)
	}
	obj, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
		"starting_address": types.StringValue(r.start),
		"ending_address":   types.StringValue(r.end),
	})

	return obj
}

// isReservationConflict recognises Morpheus's refusal to delete a range that
// still has allocated addresses.
func isReservationConflict(body string) bool {
	b := strings.ToLower(body)

	// "An ip range has a reservation attached that would become invalid"
	// (gomorpheus.networkPools.ipRangeHasReservation), under the ipRanges key.
	return strings.Contains(b, "reservation") && (strings.Contains(b, "ip range") || strings.Contains(b, "iprange"))
}

// rangesNotRemoved lists planned-away ranges the appliance still reports.
// Morpheus does not fail an update that tries to delete a range holding
// reservations; it keeps the range and answers 200, so the only way to see the
// refusal is to compare what came back with what was asked for.
func rangesNotRemoved(planned, api []ipRange) []ipRange {
	want := map[string]bool{}
	for _, r := range planned {
		want[r.key()] = true
	}
	var kept []ipRange
	for _, r := range api {
		if !want[r.key()] {
			kept = append(kept, r)
		}
	}

	return kept
}
