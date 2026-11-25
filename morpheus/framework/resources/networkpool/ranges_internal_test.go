// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool

import (
	"context"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func TestUnitOrderLikePrior(t *testing.T) {
	t.Parallel()

	api := []ipRange{
		{id: 3, start: "10.0.0.50", end: "10.0.0.60"}, // edited: was .59 in prior
		{id: 1, start: "10.0.0.10", end: "10.0.0.20"},
		{id: 9, start: "10.0.0.90", end: "10.0.0.99"}, // new on the appliance
		{id: 2, start: "10.0.0.30", end: "10.0.0.40"},
	}
	prior := []ipRange{
		{id: 1, start: "10.0.0.10", end: "10.0.0.20"},
		{id: 0, start: "10.0.0.30", end: "10.0.0.40"}, // created last apply, id unknown in prior -> by addresses
		{id: 3, start: "10.0.0.50", end: "10.0.0.59"}, // by id although addresses changed
	}

	got := orderLikePrior(api, prior)
	want := []int64{1, 2, 3, 9}
	if len(got) != len(want) {
		t.Fatalf("got %d ranges, want %d: %+v", len(got), len(want), got)
	}
	for i, id := range want {
		if got[i].id != id {
			t.Errorf("position %d: got id %d, want %d (%+v)", i, got[i].id, id, got)
		}
	}
	if got[2].end != "10.0.0.60" {
		t.Errorf("edited range must carry the API's addresses, got %+v", got[2])
	}
}

func TestUnitWithIDsFromPrior(t *testing.T) {
	t.Parallel()

	prior := []ipRange{{id: 1, start: "10.0.0.1", end: "10.0.0.10"}, {id: 2, start: "10.0.0.20", end: "10.0.0.30"}}

	// Same shape: an edit keeps its (positional) id, an unchanged range matches
	// by addresses even without an id in the plan.
	planned := []ipRange{
		{id: 1, start: "10.0.0.1", end: "10.0.0.12"}, // edited in place (overlaps range 1)
		{start: "10.0.0.20", end: "10.0.0.30"},       // no id in plan, matches prior by addresses
	}
	got := withIDsFromPrior(planned, prior)
	if got[0].id != 1 || got[1].id != 2 {
		t.Errorf("ids: %+v", got)
	}
	req := rangesToUpdateRequest(got)
	if req[0].Id == nil || *req[0].Id != 1 || *req[0].EndAddress != "10.0.0.12" || req[1].Id == nil || *req[1].Id != 2 {
		t.Errorf("update request: %+v %+v", req[0], req[1])
	}

	// Shape changed (a range added) together with an edit: positional hints
	// are not trusted, so the edited range is recreated and the new one has
	// no id; the unchanged range still matches by addresses.
	planned = []ipRange{
		{id: 1, start: "10.0.0.1", end: "10.0.0.12"},
		{start: "10.0.0.20", end: "10.0.0.30"},
		{start: "10.0.0.40", end: "10.0.0.50"},
	}
	got = withIDsFromPrior(planned, prior)
	if got[0].id != 0 || got[1].id != 2 || got[2].id != 0 {
		t.Errorf("ids with shape change: %+v", got)
	}
}

func TestUnitRangesFromAPISkipsCIDR(t *testing.T) {
	t.Parallel()

	id1, id2 := int64(1), int64(2)
	pool := &sdk.GetNetworkPool200ResponseNetworkPool{IpRanges: []sdk.GetNetworkPool200ResponseNetworkPoolIpRangesInner{
		{Id: &id1, StartAddress: *sdk.NewNullableString(ptr("10.0.0.1")), EndAddress: *sdk.NewNullableString(ptr("10.0.0.9"))},
		{Id: &id2, CidrIPv6: *sdk.NewNullableString(ptr("fd00::/64"))},
	}}
	ranges, unmanaged := rangesFromAPI(pool)
	if len(ranges) != 1 || ranges[0].id != 1 {
		t.Errorf("ranges: %+v", ranges)
	}
	if len(unmanaged) != 1 || unmanaged[0].id != 2 || unmanaged[0].String() != "range 2 (fd00::/64)" {
		t.Errorf("unmanaged: %+v", unmanaged)
	}

	// An unmanaged range is sent back on update with its id and cidr, so
	// Morpheus keeps it: a range missing from the request would be deleted.
	req := unmanagedToUpdateRequest(unmanaged)
	if len(req) != 1 || req[0].Id == nil || *req[0].Id != 2 || req[0].CidrIPv6 == nil || *req[0].CidrIPv6 != "fd00::/64" {
		t.Errorf("unmanaged update item: %+v", req)
	}
	if !rangesFromAPIIsNilSafe() {
		t.Error("rangesFromAPI(nil) must return nil, nil")
	}
}

func rangesFromAPIIsNilSafe() bool {
	r, u := rangesFromAPI(nil)

	return r == nil && u == nil
}

func TestUnitRangesListRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	in := []ipRange{{id: 5, start: "a", end: "b"}, {start: "c", end: "d"}}
	list, diags := rangesToList(ctx, in)
	if diags.HasError() {
		t.Fatal(diags)
	}
	out, diags := rangesFromList(ctx, list)
	if diags.HasError() {
		t.Fatal(diags)
	}
	if len(out) != 2 || out[0] != in[0] || out[1] != in[1] {
		t.Errorf("round trip: %+v", out)
	}
	if !rangesFromListMustNull(ctx) {
		t.Error("nil ranges must round-trip to a null list")
	}
}

func rangesFromListMustNull(ctx context.Context) bool {
	list, _ := rangesToList(ctx, nil)

	return list.IsNull()
}

func ptr(s string) *string { return &s }

func TestUnitWithIDsFromPriorRemovalDoesNotShiftIdentity(t *testing.T) {
	t.Parallel()

	// Prior: three ranges. The first is removed from the configuration. The
	// positional id modifier has copied prior ids onto the plan BY POSITION,
	// so plan[0] (really range 2) carries range 1's id. Address matching must
	// win: ranges 2 and 3 keep their own ids; range 1's id is sent by nobody,
	// so Morpheus deletes range 1.
	prior := []ipRange{
		{id: 10, start: "10.0.0.1", end: "10.0.0.10"},
		{id: 18, start: "10.0.0.20", end: "10.0.0.30"},
		{id: 25, start: "10.0.0.40", end: "10.0.0.50"},
	}
	planned := []ipRange{{id: 10, start: "10.0.0.20", end: "10.0.0.30"}, {id: 18, start: "10.0.0.40", end: "10.0.0.50"}}
	got := withIDsFromPrior(planned, prior)
	if got[0].id != 18 || got[1].id != 25 {
		t.Errorf("identities shifted: %+v", got)
	}

	// An edit of the middle range (addresses changed) keeps its positional id,
	// because that id was not claimed by an address match.
	planned = []ipRange{
		{id: 10, start: "10.0.0.1", end: "10.0.0.10"},
		{id: 18, start: "10.0.0.20", end: "10.0.0.35"},
		{id: 25, start: "10.0.0.40", end: "10.0.0.50"},
	}
	got = withIDsFromPrior(planned, prior)
	if got[1].id != 18 || got[0].id != 10 || got[2].id != 25 {
		t.Errorf("edit must keep its id: %+v", got)
	}

	// Remove the first AND edit the (now first) second: the hint points at the
	// removed range's id, which belongs to nobody in the plan — but using it
	// would silently turn range 1 into range 2. Treated as new instead.
	planned = []ipRange{{id: 10, start: "10.0.0.20", end: "10.0.0.35"}, {id: 18, start: "10.0.0.40", end: "10.0.0.50"}}
	got = withIDsFromPrior(planned, prior)
	if got[1].id != 25 || got[0].id != 0 {
		t.Errorf("ambiguous edit after removal must not adopt a neighbour's id: %+v", got)
	}
}

func TestUnitRangesNotRemoved(t *testing.T) {
	t.Parallel()
	planned := []ipRange{{start: "a", end: "b"}}
	api := []ipRange{{id: 1, start: "a", end: "b"}, {id: 2, start: "c", end: "d"}}
	kept := rangesNotRemoved(planned, api)
	if len(kept) != 1 || kept[0].id != 2 {
		t.Errorf("kept: %+v", kept)
	}
}

// TestUnitWithIDsFromPriorRemoveAddKeepsLength: prior [A(1),B(2)], plan
// [B′,C] — a remove and an add that leave the list the same length. The
// positional hints put A's id on B′ and B's id on C; neither overlaps the
// range it would adopt, so both are new and A and B are deleted (B, if it
// holds reservations, is kept by Morpheus and reported).
func TestUnitWithIDsFromPriorRemoveAddKeepsLength(t *testing.T) {
	t.Parallel()

	prior := []ipRange{{id: 1, start: "10.0.0.1", end: "10.0.0.10"}, {id: 2, start: "10.0.0.20", end: "10.0.0.30"}}
	planned := []ipRange{{id: 1, start: "10.0.0.20", end: "10.0.0.35"}, {id: 2, start: "10.0.0.40", end: "10.0.0.50"}}
	got := withIDsFromPrior(planned, prior)
	if got[0].id != 0 || got[1].id != 0 {
		t.Errorf("neighbour ids adopted: %+v", got)
	}

	// Same shape, an edit that moves a range to a disjoint span: it cannot
	// keep its reservations, so it is a new range rather than an in-place edit.
	planned = []ipRange{{id: 1, start: "10.0.0.100", end: "10.0.0.110"}, {start: "10.0.0.20", end: "10.0.0.30"}}
	got = withIDsFromPrior(planned, prior)
	if got[0].id != 0 || got[1].id != 2 {
		t.Errorf("disjoint edit must be new: %+v", got)
	}

	if !overlaps(ipRange{start: "10.0.0.5", end: "10.0.0.1"}, ipRange{start: "10.0.0.5", end: "10.0.0.9"}) {
		t.Error("reversed bounds must still overlap")
	}
	if overlaps(ipRange{start: "nope", end: "10.0.0.1"}, ipRange{start: "10.0.0.1", end: "10.0.0.1"}) {
		t.Error("unparsable address must not overlap")
	}
}

// TestUnitSingleRangeUpdatable: the deprecated ip_ranges attribute may update
// a pool only when it holds at most one range and that range is a start/end
// range. A lone CIDR range is refused too: ip_ranges would add a second range
// beside it and every later apply would be refused.
func TestUnitSingleRangeUpdatable(t *testing.T) {
	t.Parallel()

	one := []ipRange{{id: 1, start: "10.0.0.1", end: "10.0.0.10"}}
	cidr := []unmanagedRange{{id: 2, cidrIPv6: "fd00::/64"}}
	cases := []struct {
		name      string
		managed   []ipRange
		unmanaged []unmanagedRange
		total     int
		ok        bool
	}{
		{"empty pool", nil, nil, 0, true},
		{"one managed", one, nil, 1, true},
		{"two managed", append(one, ipRange{id: 3, start: "10.0.0.20", end: "10.0.0.30"}), nil, 2, false},
		{"one managed + cidr", one, cidr, 2, false},
		{"only cidr", nil, cidr, 1, false},
	}
	for _, c := range cases {
		total, ok := singleRangeUpdatable(c.managed, c.unmanaged)
		if total != c.total || ok != c.ok {
			t.Errorf("%s: got total=%d ok=%v, want %d %v", c.name, total, ok, c.total, c.ok)
		}
	}
}
