// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// mapPrivate is a minimal wodigest.PrivateState backed by a map, so the
// tracker's record read/write can be exercised without the framework.
type mapPrivate struct{ m map[string][]byte }

func newMapPrivate() *mapPrivate { return &mapPrivate{m: map[string][]byte{}} }

func (p *mapPrivate) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return p.m[key], nil
}

func (p *mapPrivate) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	p.m[key] = value

	return nil
}

func serverHash(s string) *sdk.GetHost200ResponseServer {
	return &sdk.GetHost200ResponseServer{SshPasswordHash: *sdk.NewNullableString(&s)}
}

// TestUnitInstanceNodeDriftWarnPath covers the Warn-path detection for the
// instance_node tracker: a baseline seeded by wo.Observe, then a later read
// whose sshPasswordHash differs, must be recorded as drift by wo.CheckDrift and
// returned as the drifted attribute. Because ssh_password is Immutable,
// planDrift (covered in wodigest's own tests) forces that to a warning; here we
// prove the detection half for this resource's real tracker.
func TestUnitInstanceNodeDriftWarnPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := newMapPrivate()

	// Seed configured-value digests directly, then a baseline off a read.
	seed := wodigest.Compute(wodigest.Attrs{sshPasswordAttr: "digest"}, wo.Immutable)
	if d := wodigest.Store(ctx, private, wo.Key(), seed); d.HasError() {
		t.Fatalf("seed: %v", d)
	}

	if d := wo.Observe(ctx, private, serverHash("h-a")); d.HasError() {
		t.Fatalf("Observe: %v", d)
	}

	if !hasDriftBaseline(ctx, private) {
		t.Fatal("a baseline should be recorded after Observe")
	}

	drifted, d := wo.CheckDrift(ctx, private, serverHash("DIFFERENT"))
	if d.HasError() {
		t.Fatalf("CheckDrift: %v", d)
	}

	if len(drifted) != 1 || drifted[0] != sshPasswordAttr {
		t.Fatalf("expected %q drifted, got %v", sshPasswordAttr, drifted)
	}
}

// TestUnitInstanceNodeNoBaselineSeedsOnRead proves the create-time rule: with
// no baseline and a hash now present, Read seeds the baseline (the resource
// calls wo.Observe) rather than comparing. We assert the discriminator Read
// uses (hasDriftBaseline) flips from false to true once Observe runs.
func TestUnitInstanceNodeNoBaselineSeedsOnRead(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := newMapPrivate()

	seed := wodigest.Compute(wodigest.Attrs{sshPasswordAttr: "digest"}, wo.Immutable)
	if d := wodigest.Store(ctx, private, wo.Key(), seed); d.HasError() {
		t.Fatalf("seed: %v", d)
	}

	if hasDriftBaseline(ctx, private) {
		t.Fatal("no baseline should exist before Observe")
	}

	if d := wo.Observe(ctx, private, serverHash("h-a")); d.HasError() {
		t.Fatalf("Observe: %v", d)
	}

	if !hasDriftBaseline(ctx, private) {
		t.Fatal("a baseline should exist after Observe")
	}
}
