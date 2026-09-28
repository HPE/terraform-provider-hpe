// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package wodigest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func sha(s string) string {
	sum := sha256.Sum256([]byte(s))

	return hex.EncodeToString(sum[:])
}

func TestUnitValueMatchesPlainSha256ForStrings(t *testing.T) {
	t.Parallel()

	// An importer reproduces a string value's digest with sha256 alone.
	got, err := Value("Secret123!")
	if err != nil || got != sha("Secret123!") {
		t.Fatalf("Value(string) = %s, %v; want sha256 of the string", got, err)
	}

	s := "x"
	if got, _ := Value(&s); got != sha("x") {
		t.Fatal("a *string hashes like its string")
	}

	var nilStr *string
	if got, _ := Value(nilStr); got != "" {
		t.Fatal("a nil *string is not set")
	}

	if got, _ := Value(nil); got != "" {
		t.Fatal("nil is not set")
	}

	if got, _ := Value(""); got != sha("") {
		t.Fatal("an empty string is set, and distinct from nil")
	}
}

func TestUnitValueIsCanonicalForObjects(t *testing.T) {
	t.Parallel()

	a := map[string]any{"b": 1, "a": map[string]any{"y": []any{1, "two"}, "x": true}}
	b := map[string]any{"a": map[string]any{"x": true, "y": []any{1, "two"}}, "b": 1}

	da, _ := Value(a)
	db, _ := Value(b)

	if da != db {
		t.Fatalf("same content, different digests: %s vs %s", da, db)
	}

	// The same as sha256 of the sorted-key JSON an importer would produce.
	if da != sha(`{"a":{"x":true,"y":[1,"two"]},"b":1}`) {
		t.Fatal("object digest is not sha256 of its canonical JSON")
	}

	c := map[string]any{"a": map[string]any{"x": true, "y": []any{"two", 1}}, "b": 1}
	if dc, _ := Value(c); dc == da {
		t.Fatal("slice order is significant and must change the digest")
	}
}

func TestUnitComputeAndImmutable(t *testing.T) {
	t.Parallel()

	attrs, err := Values(map[string]any{"password_wo": "s3cret", "key": "k", "unset": nil})
	if err != nil {
		t.Fatal(err)
	}

	if len(attrs) != 2 {
		t.Fatalf("unset values must be dropped: %v", attrs)
	}

	d := Compute(attrs, nil)
	if d.All == "" || d.Immutable != "" || d.Version != Version {
		t.Fatalf("unexpected digests without an immutable subset: %+v", d)
	}

	// The digest is over the attribute digest set, so an importer with only
	// the hashes reproduces it.
	if d.All != Digest(Attrs{"password_wo": sha("s3cret"), "key": sha("k")}) {
		t.Fatal("All must be the digest of the {attribute: sha256} set")
	}

	d2 := Compute(attrs, []string{"key"})
	if d2.All != d.All || d2.Immutable == "" || d2.Immutable == d2.All {
		t.Fatalf("unexpected digests with an immutable subset: %+v", d2)
	}

	if d2.Immutable != Digest(Attrs{"key": sha("k")}) {
		t.Fatal("Immutable must cover only the named attributes")
	}

	// An immutable attribute that is unset contributes nothing.
	d3 := Compute(attrs, []string{"unset"})
	if d3.Immutable != Digest(Attrs{}) {
		t.Fatal("an unset immutable attribute must digest as empty")
	}

	// ComputeValues is the same thing from raw values.
	d4, err := ComputeValues(map[string]any{"password_wo": "s3cret", "key": "k", "unset": nil}, []string{"key"})
	if err != nil || d4.All != d2.All || d4.Immutable != d2.Immutable {
		t.Fatalf("ComputeValues = %+v, %v; want %+v", d4, err, d2)
	}
}

func TestUnitDecide(t *testing.T) {
	t.Parallel()

	stored := Digests{All: "A", Immutable: "I", Version: Version}

	cases := []struct {
		name    string
		stored  Digests
		ok      bool
		current Digests
		want    Decision
	}{
		{"nothing stored", Digests{}, false, Digests{All: "A"}, Seed},
		{"nothing stored, immutable differs too", Digests{}, false, Digests{All: "B", Immutable: "J"}, Seed},
		{"import sentinel", Digests{Sentinel: ImportSentinel, Version: Version}, true, Digests{All: "B", Immutable: "J"}, Unchanged},
		{"unchanged", stored, true, Digests{All: "A", Immutable: "I"}, Unchanged},
		{"unchanged whatever immutable says", stored, true, Digests{All: "A", Immutable: "X"}, Unchanged},
		{"updatable value changed", stored, true, Digests{All: "B", Immutable: "I"}, Update},
		{"immutable value changed", stored, true, Digests{All: "B", Immutable: "J"}, Replace},
		{"no immutable subset", Digests{All: "A", Version: Version}, true, Digests{All: "B"}, Update},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := Decide(tc.stored, tc.ok, tc.current); got != tc.want {
				t.Fatalf("Decide = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestUnitDecideConfigured(t *testing.T) {
	t.Parallel()

	stored := Compute(Attrs{"password_wo": "P", "linux_password_wo": "L"}, nil)
	imm := []string{"linux_password_wo"}
	noAttrs := Digests{All: Digest(Attrs{"password_wo": "P"}), Version: Version}

	cases := []struct {
		name    string
		stored  Digests
		ok      bool
		current Attrs
		want    Decision
	}{
		{"nothing stored", Digests{}, false, Attrs{"password_wo": "P"}, Seed},
		{"sentinel", Digests{Sentinel: ImportSentinel, Version: Version}, true, Attrs{"password_wo": "X"}, Unchanged},
		{"same", stored, true, Attrs{"password_wo": "P", "linux_password_wo": "L"}, Unchanged},
		{"one left unset is not a change", stored, true, Attrs{"password_wo": "P"}, Unchanged},
		{"all left unset is not a change", stored, true, Attrs{}, Unchanged},
		{"changed", stored, true, Attrs{"password_wo": "Q"}, Update},
		{"newly set", stored, true, Attrs{"password_wo": "P", "windows_password_wo": "W"}, Update},
		{"immutable changed", stored, true, Attrs{"linux_password_wo": "M"}, Replace},
		{"immutable and mutable changed", stored, true, Attrs{"password_wo": "Q", "linux_password_wo": "M"}, Replace},
		{"record without attrs compares whole", noAttrs, true, Attrs{"password_wo": "P"}, Unchanged},
		{"record without attrs, changed", noAttrs, true, Attrs{"password_wo": "Q"}, Update},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := DecideConfigured(tc.stored, tc.ok, tc.current, imm); got != tc.want {
				t.Fatalf("DecideConfigured = %s, want %s", got, tc.want)
			}
		})
	}
}

// fakePrivate is the framework's private state, as far as the digests use it.
type fakePrivate struct {
	data map[string][]byte
}

func (f *fakePrivate) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f.data[key], nil
}

func (f *fakePrivate) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if f.data == nil {
		f.data = map[string][]byte{}
	}

	f.data[key] = value

	return nil
}

func TestUnitDigestsRoundTripThroughPrivateState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivate{}
	const key = "test.wo_digests"

	if _, ok, _ := Load(ctx, private, key); ok {
		t.Fatal("expected no digests before any are stored")
	}

	want := Digests{All: "a", Immutable: "b"} // Version filled by Store
	if diags := Store(ctx, private, key, want); diags.HasError() {
		t.Fatal(diags)
	}

	got, ok, diags := Load(ctx, private, key)
	if diags.HasError() || !ok {
		t.Fatalf("expected stored digests, got ok=%v diags=%v", ok, diags)
	}

	if got.All != "a" || got.Immutable != "b" || got.Version != Version || got.Sentinel != "" || got.Attrs != nil {
		t.Fatalf("round trip changed the digests: %+v", got)
	}

	// Two keys, two independent records.
	if diags := Store(ctx, private, "other", Digests{All: "z"}); diags.HasError() {
		t.Fatal(diags)
	}

	if got, _, _ := Load(ctx, private, key); got.All != "a" {
		t.Fatal("a second key must not disturb the first")
	}

	// A digest of a different version -- including the pre-release v1 form --
	// is treated as missing.
	private.data[key] = []byte(`{"all":"a","immutable":"b","v":1}`)
	if _, ok, _ := Load(ctx, private, key); ok {
		t.Fatal("expected a different digest version to read as missing")
	}

	// As is anything unparseable, or a record with nothing in it.
	private.data[key] = []byte(`not json`)
	if _, ok, _ := Load(ctx, private, key); ok {
		t.Fatal("expected garbage to read as missing")
	}

	private.data[key] = []byte(`{"v":2}`)
	if _, ok, _ := Load(ctx, private, key); ok {
		t.Fatal("expected an empty record to read as missing")
	}

	// Nil private state: Load says missing; Store reports an error rather than
	// silently dropping the digests.
	if _, ok, _ := Load(ctx, nil, key); ok {
		t.Fatal("expected nil private state to read as missing")
	}

	if diags := Store(ctx, nil, key, want); !diags.HasError() {
		t.Fatal("expected an error storing to nil private state")
	}
}

func TestUnitImportSentinel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivate{}
	const key = "test.wo_digests"

	if diags := StoreImportSentinel(ctx, private, key); diags.HasError() {
		t.Fatal(diags)
	}

	var raw map[string]any
	if err := json.Unmarshal(private.data[key], &raw); err != nil {
		t.Fatal(err)
	}

	if raw["sentinel"] != ImportSentinel || raw["all"] != "" {
		t.Fatalf("unexpected sentinel record: %s", private.data[key])
	}

	got, ok, _ := Load(ctx, private, key)
	if !ok || got.Sentinel != ImportSentinel {
		t.Fatalf("expected the sentinel to load, got ok=%v %+v", ok, got)
	}

	if Decide(got, ok, Digests{All: "anything"}) != Unchanged {
		t.Fatal("the first plan after import must be clean")
	}

	// The first apply of any kind replaces the sentinel with real digests.
	real := Compute(Attrs{"password_wo": sha("x")}, nil)
	if diags := Store(ctx, private, key, real); diags.HasError() {
		t.Fatal(diags)
	}

	got, ok, _ = Load(ctx, private, key)
	if !ok || got.Sentinel != "" || got.All != real.All {
		t.Fatalf("expected real digests after the apply, got %+v", got)
	}
}

func TestUnitImportID(t *testing.T) {
	t.Parallel()

	known := []string{"password_wo", "linux_password_wo", "windows_password_wo"}
	h := sha("Secret123!")

	// Plain ID: no digests.
	id, attrs, err := ParseImportID("123", known)
	if err != nil || id != "123" || attrs != nil {
		t.Fatalf("plain: %q %v %v", id, attrs, err)
	}

	// One digest; upper-case hex is accepted.
	id, attrs, err = ParseImportID("123#password_wo="+strings.ToUpper(h), known)
	if err != nil || id != "123" || len(attrs) != 1 || attrs["password_wo"] != h {
		t.Fatalf("one: %q %v %v", id, attrs, err)
	}

	// The digests it yields are the ones the resource would compute from the
	// same configured value -- the whole point.
	fromImport := Compute(attrs, nil)
	fromConfig, _ := ComputeValues(map[string]any{"password_wo": "Secret123!", "linux_password_wo": nil}, nil)

	if fromImport.All != fromConfig.All {
		t.Fatal("digests from the import ID must equal digests from the configuration")
	}

	// Several, in any order; FormatImportID sorts.
	raw := "123#windows_password_wo=" + sha("w") + "&password_wo=" + h
	id, attrs, err = ParseImportID(raw, known)
	if err != nil || id != "123" || len(attrs) != 2 {
		t.Fatalf("two: %q %v %v", id, attrs, err)
	}

	if got := FormatImportID(id, attrs); got != "123#password_wo="+h+"&windows_password_wo="+sha("w") {
		t.Fatalf("FormatImportID = %s", got)
	}

	if FormatImportID("123", nil) != "123" {
		t.Fatal("FormatImportID with no digests is the plain ID")
	}

	// Errors: nothing after '#', bad pair, unknown attribute, bad hex, repeat.
	for _, bad := range []string{
		"123#",
		"123#password_wo",
		"123#=abc",
		"123#passwd=" + h,
		"123#password_wo=nothex",
		"123#password_wo=" + h[:63],
		"123#password_wo=" + h + "&password_wo=" + h,
	} {
		if _, _, err := ParseImportID(bad, known); err == nil {
			t.Errorf("ParseImportID(%q) succeeded, want an error", bad)
		}
	}
}
