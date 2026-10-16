// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package wodigest

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// fakeResp is a stand-in for a resource's read response: a map of attribute to
// the (hash, present) an extractor would pull from it. A key absent from the
// map models a field the appliance did not return at all.
type fakeResp struct {
	hashes map[string]struct {
		hash    string
		present bool
	}
}

// extractor builds a HashExtractor for one attribute over a fakeResp.
func extractor(attr string) HashExtractor {
	return func(resp any) (string, bool) {
		r, ok := resp.(*fakeResp)
		if !ok || r == nil {
			return "", false
		}

		h, ok := r.hashes[attr]
		if !ok {
			return "", false
		}

		return h.hash, h.present
	}
}

func resp(entries map[string]struct {
	hash    string
	present bool
},
) *fakeResp {
	return &fakeResp{hashes: entries}
}

func present(hash string) struct {
	hash    string
	present bool
} {
	return struct {
		hash    string
		present bool
	}{hash, true}
}

func absent() struct {
	hash    string
	present bool
} {
	return struct {
		hash    string
		present bool
	}{"", false}
}

// driftTracker tracks password_wo and key_wo (both Resend) over the test
// schema. Both have an extractor so neither is tracked without one.
var driftTracker = Tracker{
	Resource:   "thing",
	Attributes: map[string]string{"password_wo": "password_wo_version", "key_wo": "key_wo_version"},
	APIHash: map[string]HashExtractor{
		"password_wo": extractor("password_wo"),
		"key_wo":      extractor("key_wo"),
	},
}

func seedDigests(t *testing.T, private *fakePrivate, tr Tracker, attrs Attrs) {
	t.Helper()

	if diags := Store(context.Background(), private, tr.Key(), Compute(attrs, tr.Immutable)); diags.HasError() {
		t.Fatal(diags)
	}
}

func loadRec(t *testing.T, private *fakePrivate, tr Tracker) Digests {
	t.Helper()

	d, ok, diags := Load(context.Background(), private, tr.Key())
	if diags.HasError() || !ok {
		t.Fatalf("expected a record, ok=%v diags=%v", ok, diags)
	}

	return d
}

// TestUnitV2RecordLoadsAsV3 is the compatibility guarantee: a record written
// at version 2 (before Observed/Drift existed) still loads, keeps its digests,
// and simply has empty drift maps.
func TestUnitV2RecordLoadsAsV3(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivate{}
	key := driftTracker.Key()

	v2 := Compute(Attrs{"password_wo": sha("x")}, nil)
	v2.Version = 2

	raw, _ := json.Marshal(v2)
	private.SetKey(ctx, key, raw)

	got, ok, diags := Load(ctx, private, key)
	if diags.HasError() || !ok {
		t.Fatalf("a v2 record must load, ok=%v diags=%v", ok, diags)
	}

	if got.All != v2.All || got.Attrs["password_wo"] != sha("x") {
		t.Fatalf("v2 digests must survive the load: %+v", got)
	}

	if len(got.Observed) != 0 || len(got.Drift) != 0 {
		t.Fatalf("a v2 record must have no drift baseline: %+v", got)
	}
}

func TestUnitObserve(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivate{}
	seedDigests(t, private, driftTracker, Attrs{"password_wo": sha("pw")})

	// Not present: nothing recorded for that attribute.
	diags := driftTracker.Observe(ctx, private, resp(map[string]struct {
		hash    string
		present bool
	}{"password_wo": absent()}))
	if diags.HasError() {
		t.Fatal(diags)
	}

	if rec := loadRec(t, private, driftTracker); len(rec.Observed) != 0 {
		t.Fatalf("a not-present field must not seed a baseline: %+v", rec.Observed)
	}

	// Present: baseline recorded, drift cleared.
	diags = driftTracker.Observe(ctx, private, resp(map[string]struct {
		hash    string
		present bool
	}{"password_wo": present("HASH")}))
	if diags.HasError() {
		t.Fatal(diags)
	}

	rec := loadRec(t, private, driftTracker)
	if rec.Observed["password_wo"] != "HASH" || rec.Drift["password_wo"] {
		t.Fatalf("present field must set the baseline and clear drift: %+v", rec)
	}

	// nil resp is a no-op.
	before := loadRec(t, private, driftTracker)
	if diags := driftTracker.Observe(ctx, private, nil); diags.HasError() {
		t.Fatal(diags)
	}

	after := loadRec(t, private, driftTracker)
	if after.Observed["password_wo"] != before.Observed["password_wo"] {
		t.Fatal("nil resp must not change the record")
	}
}

func TestUnitCheckDrift(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	e := func(k, v string, p bool) map[string]struct {
		hash    string
		present bool
	} {
		return map[string]struct {
			hash    string
			present bool
		}{k: {v, p}}
	}

	cases := []struct {
		name       string
		baseline   string // Observed value seeded; "" means no baseline
		readHash   string
		readPresnt bool
		wantDrift  bool
		wantReturn bool // attribute returned as drifted
	}{
		{"equal", "H", "H", true, false, false},
		{"different", "H", "H2", true, true, true},
		{"cleared after non-empty baseline", "H", "", true, true, true},
		{"not present leaves unchanged", "H", "ignored", false, false, false},
		{"no baseline is untouched", "", "H", true, false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			private := &fakePrivate{}
			seedDigests(t, private, driftTracker, Attrs{"password_wo": sha("pw")})

			if tc.baseline != "" {
				rec := loadRec(t, private, driftTracker)
				rec.Observed = map[string]string{"password_wo": tc.baseline}
				if diags := Store(ctx, private, driftTracker.Key(), rec); diags.HasError() {
					t.Fatal(diags)
				}
			}

			drifted, diags := driftTracker.CheckDrift(ctx, private, resp(e("password_wo", tc.readHash, tc.readPresnt)))
			if diags.HasError() {
				t.Fatal(diags)
			}

			gotReturn := len(drifted) == 1 && drifted[0] == "password_wo"
			if gotReturn != tc.wantReturn {
				t.Fatalf("drifted = %v, want returned=%v", drifted, tc.wantReturn)
			}

			if tc.baseline == "" {
				// No baseline: the record's Drift must stay empty.
				if rec := loadRec(t, private, driftTracker); rec.Drift["password_wo"] {
					t.Fatal("an attribute with no baseline must not be flagged")
				}

				return
			}

			if rec := loadRec(t, private, driftTracker); rec.Drift["password_wo"] != tc.wantDrift {
				t.Fatalf("Drift[password_wo] = %v, want %v", rec.Drift["password_wo"], tc.wantDrift)
			}
		})
	}
}

// driftedRec builds a stored record whose digests match the configured value
// "pw" (so the digest decision is Unchanged) and whose drift baseline for attr
// is recorded as drifted or not.
func driftedRec(attr string, drifted bool) Digests {
	d := Compute(Attrs{attr: sha("pw")}, nil)
	d.Observed = map[string]string{attr: "H"}
	d.Drift = map[string]bool{attr: drifted}

	return d
}

// configuredPw is the configured attrs for a set password_wo.
var configuredPw = Attrs{"password_wo": sha("pw")}

func TestUnitPlanDrift(t *testing.T) {
	t.Parallel()

	resendTr := driftTracker

	warnTr := Tracker{
		Resource:   "thing",
		Attributes: map[string]string{"password_wo": "password_wo_version"},
		APIHash:    map[string]HashExtractor{"password_wo": extractor("password_wo")},
		OnDrift:    WarnPolicy,
	}

	immutableTr := Tracker{
		Resource:   "thing",
		Attributes: map[string]string{"password_wo": "password_wo_version"},
		Immutable:  []string{"password_wo"},
		APIHash:    map[string]HashExtractor{"password_wo": extractor("password_wo")},
		OnDrift:    ResendPolicy, // forced to Warn because immutable
	}

	plainTr := Tracker{
		Resource:   "thing",
		Attributes: map[string]string{"password_wo": "password_wo_version"},
	}

	cases := []struct {
		name        string
		tr          Tracker
		stored      Digests
		current     Attrs
		wantResend  bool
		wantWarning int
	}{
		{
			"drift + configured + Resend: resend + informational diagnostic",
			resendTr, driftedRec("password_wo", true), configuredPw, true, 1,
		},
		{
			"drift + Warn: warning, no resend",
			warnTr, driftedRec("password_wo", true), configuredPw, false, 1,
		},
		{
			"drift on an immutable attribute: warning even with OnDrift Resend",
			immutableTr, driftedRec("password_wo", true), configuredPw, false, 1,
		},
		{
			"drift + attribute not configured: nothing",
			resendTr, driftedRec("password_wo", true),
			Attrs{},
			false, 0,
		},
		{
			"no drift: nothing, as today",
			resendTr, driftedRec("password_wo", false), configuredPw, false, 0,
		},
		{
			"no APIHash: drift map ignored entirely",
			plainTr, driftedRec("password_wo", true), configuredPw, false, 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resend, diags := tc.tr.planDrift(tc.stored, tc.current)
			if resend != tc.wantResend {
				t.Fatalf("resend = %v, want %v", resend, tc.wantResend)
			}

			if diags.WarningsCount() != tc.wantWarning {
				t.Fatalf("warnings = %d, want %d: %v", diags.WarningsCount(), tc.wantWarning, diags)
			}

			if diags.HasError() {
				t.Fatalf("planDrift must not error: %v", diags)
			}
		})
	}
}

// TestUnitResendLearnsDriftThroughRevision documents the channel: ModifyPlan
// marks write_only_revision unknown for a drift-driven re-send exactly as it
// does for a digest change, and Resend reads that unknown. planDrift returning
// resend=true is what drives ModifyPlan down that same path, so no new signal
// is needed. Here we assert the revision-unknown plan that follows is read back
// as a re-send by Resend.
func TestUnitResendReadsUnknownRevision(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// A plan whose write_only_revision is unknown -- what ModifyPlan sets for a
	// drift-driven re-send.
	raw := trackerValue(t, map[string]tftypes.Value{
		"name":            str("x"),
		RevisionAttribute: tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
	})
	plan := tfsdk.Plan{Schema: trackerSchema, Raw: raw}

	if !Resend(ctx, plan) {
		t.Fatal("Resend must be true when the revision is unknown")
	}

	// And false for a known revision.
	known := tfsdk.Plan{Schema: trackerSchema, Raw: trackerValue(t, map[string]tftypes.Value{
		"name": str("x"), RevisionAttribute: num(7),
	})}
	if Resend(ctx, known) {
		t.Fatal("Resend must be false when the revision is known")
	}
}

// cfgFor builds a tfsdk.Config over trackerSchema with password_wo set to v.
func cfgFor(t *testing.T, v string) tfsdk.Config {
	t.Helper()

	return tfsdk.Config{Schema: trackerSchema, Raw: trackerValue(t, map[string]tftypes.Value{
		"name": str("x"), "password_wo": str(v),
	})}
}

// TestUnitStoreCarriesDriftForwardThroughStore exercises the carry-forward
// through Tracker.Store itself (not the Store helper seedDigests uses): a Store
// that records the configured digests must preserve the Observed/Drift baseline
// a previous CheckDrift wrote, so an unrelated apply does not wipe it.
func TestUnitStoreCarriesDriftForwardThroughStore(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivate{}

	// Seed a real record carrying a drift baseline, as CheckDrift would leave.
	rec := Compute(Attrs{"password_wo": sha("pw")}, nil)
	rec.Observed = map[string]string{"password_wo": "HASH"}
	rec.Drift = map[string]bool{"password_wo": true}

	if diags := Store(ctx, private, driftTracker.Key(), rec); diags.HasError() {
		t.Fatal(diags)
	}

	// An apply that stores the (unchanged) configured digests.
	if diags := driftTracker.Store(ctx, private, cfgFor(t, "pw")); diags.HasError() {
		t.Fatal(diags)
	}

	got := loadRec(t, private, driftTracker)
	if got.Observed["password_wo"] != "HASH" {
		t.Fatalf("Store must carry the Observed baseline forward: %+v", got.Observed)
	}

	if !got.Drift["password_wo"] {
		t.Fatalf("Store must carry the Drift flag forward: %+v", got.Drift)
	}
}

// TestUnitStoreCarriesDriftForwardOverSentinel is the regression for defect 1:
// on an import-sentinel record (a non-merge path), Store must still carry the
// Observed/Drift a CheckDrift ran after import recorded, not wipe them to nil.
func TestUnitStoreCarriesDriftForwardOverSentinel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivate{}

	// An import sentinel that has since gained a drift baseline via CheckDrift.
	rec := Digests{Sentinel: ImportSentinel, Version: Version}
	rec.Observed = map[string]string{"password_wo": "HASH"}
	rec.Drift = map[string]bool{"password_wo": true}

	if diags := Store(ctx, private, driftTracker.Key(), rec); diags.HasError() {
		t.Fatal(diags)
	}

	if diags := driftTracker.Store(ctx, private, cfgFor(t, "pw")); diags.HasError() {
		t.Fatal(diags)
	}

	got := loadRec(t, private, driftTracker)
	if got.Observed["password_wo"] != "HASH" || !got.Drift["password_wo"] {
		t.Fatalf("a sentinel record's Observed/Drift must be preserved: %+v", got)
	}
}

// TestUnitCheckDriftPresentThenNotPresentClears is the regression for defect 2
// on the Resend/not-present side: a drift flagged true on one read is cleared
// when a later read does not present the field, so "no information" stops
// asserting a stale drift.
func TestUnitCheckDriftPresentThenNotPresentClears(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivate{}
	seedDigests(t, private, driftTracker, Attrs{"password_wo": sha("pw")})

	rec := loadRec(t, private, driftTracker)
	rec.Observed = map[string]string{"password_wo": "H"}

	if diags := Store(ctx, private, driftTracker.Key(), rec); diags.HasError() {
		t.Fatal(diags)
	}

	// Present and different: drift goes true.
	drifted, diags := driftTracker.CheckDrift(ctx, private, resp(map[string]struct {
		hash    string
		present bool
	}{"password_wo": present("H2")}))
	if diags.HasError() {
		t.Fatal(diags)
	}

	if len(drifted) != 1 || loadRec(t, private, driftTracker).Drift["password_wo"] != true {
		t.Fatalf("expected drift true after a differing present read: %v", drifted)
	}

	// Not present on the next read: the flag must be cleared, not left stale.
	drifted, diags = driftTracker.CheckDrift(ctx, private, resp(map[string]struct {
		hash    string
		present bool
	}{"password_wo": absent()}))
	if diags.HasError() {
		t.Fatal(diags)
	}

	if len(drifted) != 0 {
		t.Fatalf("a not-present read must not report drift: %v", drifted)
	}

	if _, had := loadRec(t, private, driftTracker).Drift["password_wo"]; had {
		t.Fatal("a not-present read must clear the drift flag, not leave it")
	}
}

// TestUnitCheckDriftWarnStaysTrueWhileDiffers is the regression for defect 2 on
// the Warn/immutable side: while the appliance value genuinely still differs,
// re-running CheckDrift keeps the flag true (so the warning correctly repeats
// until the user acts). CheckDrift, the single owner of the flag, re-asserts
// true on each differing read.
func TestUnitCheckDriftWarnStaysTrueWhileDiffers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivate{}
	seedDigests(t, private, driftTracker, Attrs{"password_wo": sha("pw")})

	rec := loadRec(t, private, driftTracker)
	rec.Observed = map[string]string{"password_wo": "H"}

	if diags := Store(ctx, private, driftTracker.Key(), rec); diags.HasError() {
		t.Fatal(diags)
	}

	differing := resp(map[string]struct {
		hash    string
		present bool
	}{"password_wo": present("H2")})

	// Two plans in a row, same still-differing appliance hash: still drifted.
	for i := 0; i < 2; i++ {
		drifted, diags := driftTracker.CheckDrift(ctx, private, differing)
		if diags.HasError() {
			t.Fatal(diags)
		}

		if len(drifted) != 1 || drifted[0] != "password_wo" {
			t.Fatalf("read %d: expected sustained drift, got %v", i, drifted)
		}

		if !loadRec(t, private, driftTracker).Drift["password_wo"] {
			t.Fatalf("read %d: drift flag must stay true while the hash differs", i)
		}
	}
}

// TestUnitCheckDriftKeyWo covers key_wo, which previously had no extractor and
// no test: an independent second tracked attribute drifts on its own.
func TestUnitCheckDriftKeyWo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivate{}
	seedDigests(t, private, driftTracker, Attrs{"key_wo": sha("k")})

	rec := loadRec(t, private, driftTracker)
	rec.Observed = map[string]string{"key_wo": "KH"}

	if diags := Store(ctx, private, driftTracker.Key(), rec); diags.HasError() {
		t.Fatal(diags)
	}

	drifted, diags := driftTracker.CheckDrift(ctx, private, resp(map[string]struct {
		hash    string
		present bool
	}{"key_wo": present("KH2")}))
	if diags.HasError() {
		t.Fatal(diags)
	}

	if len(drifted) != 1 || drifted[0] != "key_wo" {
		t.Fatalf("key_wo must be detected as drifted: %v", drifted)
	}
}

// nullableString is a stand-in for the SDK's generated NullableString, matching
// the IsSet/Get interface NullableStringHash takes.
type nullableString struct {
	set bool
	val *string
}

func (n nullableString) IsSet() bool  { return n.set }
func (n nullableString) Get() *string { return n.val }

func TestUnitNullableStringHash(t *testing.T) {
	t.Parallel()

	v := "HASH"

	// Omitted (not set): no information.
	if h, p := NullableStringHash(nullableString{set: false}); h != "" || p {
		t.Fatalf("omitted field must be (\"\", false), got (%q, %v)", h, p)
	}

	// Set and null: present, no secret.
	if h, p := NullableStringHash(nullableString{set: true, val: nil}); h != "" || !p {
		t.Fatalf("null field must be (\"\", true), got (%q, %v)", h, p)
	}

	// Set with a value.
	if h, p := NullableStringHash(nullableString{set: true, val: &v}); h != "HASH" || !p {
		t.Fatalf("value field must be (\"HASH\", true), got (%q, %v)", h, p)
	}
}
