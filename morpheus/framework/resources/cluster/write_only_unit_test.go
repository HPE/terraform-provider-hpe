// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// mapPrivate is a minimal wodigest.PrivateState backed by a map, so the
// tracker's record-reading and writing can be exercised without the framework.
type mapPrivate struct{ m map[string][]byte }

func newMapPrivate() *mapPrivate { return &mapPrivate{m: map[string][]byte{}} }

func (p *mapPrivate) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return p.m[key], nil
}

func (p *mapPrivate) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	p.m[key] = value

	return nil
}

func digest(present bool, hashes map[string]string) *clusterHostDigest {
	return &clusterHostDigest{hashes: hashes, present: present}
}

// TestUnitClusterHostDriftWarnPath covers the Warn-path decision for the
// cluster tracker: a baseline seeded by wo.Observe, then a later read whose
// per-host hash map differs, must be recorded as drift by wo.CheckDrift and
// returned as the drifted attribute. Because server.ssh_password_wo is
// Immutable, planDrift (exercised in wodigest's own tests) forces that to a
// warning; here we prove the detection half for this resource's real tracker.
func TestUnitClusterHostDriftWarnPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := newMapPrivate()

	// Seed configured-value digests directly (building the full nested cluster
	// schema config in a unit test is unnecessary), then a baseline off a full
	// read through the real tracker.
	seed := wodigest.Compute(
		wodigest.Attrs{clusterHostAttr: "digest"}, wo.Immutable)
	if d := wodigest.Store(ctx, private, wo.Key(), seed); d.HasError() {
		t.Fatalf("seed: %v", d)
	}

	base := digest(true, map[string]string{"10": "h-a", "11": "h-b"})
	if d := wo.Observe(ctx, private, base); d.HasError() {
		t.Fatalf("Observe: %v", d)
	}

	// Same host set, one hash changed: drift.
	changed := digest(true, map[string]string{"10": "h-a", "11": "DIFFERENT"})

	drifted, d := wo.CheckDrift(ctx, private, changed)
	if d.HasError() {
		t.Fatalf("CheckDrift: %v", d)
	}

	if len(drifted) != 1 || drifted[0] != clusterHostAttr {
		t.Fatalf("expected %q drifted, got %v", clusterHostAttr, drifted)
	}
}

// TestUnitClusterHostSetChangeNotDrift proves the host-set-change path: when
// the server-id set differs between baseline and read, Read re-baselines with
// wo.Observe rather than comparing, so no spurious drift is recorded. Here we
// assert the discriminators Read uses: observedHostIDs parses the baseline id
// set, and sameIDs reports the mismatch.
func TestUnitClusterHostSetChangeNotDrift(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := newMapPrivate()

	seed := wodigest.Compute(
		wodigest.Attrs{clusterHostAttr: "digest"}, wo.Immutable)
	if d := wodigest.Store(ctx, private, wo.Key(), seed); d.HasError() {
		t.Fatalf("seed: %v", d)
	}

	base := digest(true, map[string]string{"10": "h-a", "11": "h-b"})
	if d := wo.Observe(ctx, private, base); d.HasError() {
		t.Fatalf("Observe: %v", d)
	}

	ids, have, unknown := observedHostIDs(ctx, private)
	if !have || unknown {
		t.Fatalf("expected a known baseline, got have=%v unknown=%v", have, unknown)
	}

	// Scale out: host 12 added.
	scaled := digest(true, map[string]string{"10": "h-a", "11": "h-b", "12": "h-c"})
	if sameIDs(ids, digestHostIDs(scaled)) {
		t.Fatal("host set changed but sameIDs reported equal")
	}

	// A partial read must be reported not-present by the extractor.
	if _, present := clusterHostHashExtractor(digest(false, nil)); present {
		t.Fatal("a partial (not-present) digest must not be present")
	}
}
