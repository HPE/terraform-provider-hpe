// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package wodigest detects changes to write-only attributes.
//
// A write-only value is never stored in state, so Terraform cannot see it
// change: prior and planned state are both null. HashiCorp's documented answer,
// and the one used here, is to keep a secure digest of the value in the
// resource's private state and compare it during plan modification. Private
// state carries hex digests and nothing recoverable.
//
// The digest has two levels so that an importer can reproduce it knowing only
// the values. Each write-only attribute's value is hashed on its own -- a
// string as-is, anything else through its canonical JSON form -- and the
// attribute digest set {attribute: hex} is hashed again into the digest the
// resource stores. An importer that knows a value can therefore supply its
// sha256 in the import ID ("<id>#password_wo=<hex>") and the resource stores
// the same digests it would have computed itself; see ParseImportID.
//
// A resource keeps two digests. All covers every write-only value; Immutable
// covers the subset the API will not change after create, so that a change to
// one of those plans a replacement rather than an update the API would
// silently ignore. Resources whose write-only values are all updatable leave
// Immutable empty.
//
// A change to a write-only value cannot itself produce a plan diff, so the
// resource plans an update by marking a computed attribute unknown -- by
// convention write_only_revision, which Create and Update set to the time of
// the apply that last sent a write-only value, together with every other
// computed attribute the re-send changes (Terraform checks plan/apply
// consistency on all of them).
//
// The digests are written by Create, by Update (always, whether or not a value
// was re-sent, so an import sentinel or a missing digest never survives an
// update) and by ImportState (real digests when the import ID carries them,
// otherwise a sentinel: the configured values are taken as applied until the
// first apply records real digests, so the plan after import is clean). Read
// never touches them: it has no configuration to hash.
package wodigest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Version changes when the canonical form changes, so a digest written by an
// older provider is treated as missing rather than compared.
const Version = 2

// ImportSentinel marks digests stored by ImportState without digests: nothing
// has been hashed, and the configured values are accepted as applied until the
// first apply.
const ImportSentinel = "import"

// Digests are the sha256 hex digests kept in private state.
type Digests struct {
	// All covers every write-only value as sent to the API.
	All string `json:"all"`
	// Immutable covers only the values the API will not change after create.
	// Empty when the resource has none.
	Immutable string `json:"immutable"`
	// Sentinel is set instead of digests when nothing has been hashed yet.
	Sentinel string `json:"sentinel,omitempty"`
	// Attrs are the per-attribute digests All and Immutable were computed
	// from, kept so a resource can compare only the attributes that are
	// configured (DecideConfigured).
	Attrs   Attrs `json:"attrs,omitempty"`
	Version int   `json:"v"`
}

// Equal reports whether two Digests have the same digests, sentinel and
// version (Attrs are derived and not compared).
func (d Digests) Equal(o Digests) bool {
	return d.All == o.All && d.Immutable == o.Immutable && d.Sentinel == o.Sentinel && d.Version == o.Version
}

// Attrs are the per-attribute digests: attribute path (dotted for nested
// attributes, e.g. "config_vmware.password") to sha256 hex of the value. An
// attribute that is not set is absent.
type Attrs map[string]string

// Decision is what ModifyPlan should do about the write-only values.
type Decision int

const (
	// Unchanged: the configured values are what was last sent (or, after an
	// import, are accepted as such). Plan nothing.
	Unchanged Decision = iota
	// Update: a value changed; plan an update that re-sends the values.
	Update
	// Replace: a value the API will not change after create changed; plan a
	// replacement.
	Replace
	// Seed: nothing is recorded (a resource created by a provider without
	// digests, or a digest of a different version); plan a one-time update
	// that sends the values and records their digests. Never a replacement.
	Seed
)

func (d Decision) String() string {
	switch d {
	case Unchanged:
		return "unchanged"
	case Update:
		return "update"
	case Replace:
		return "replace"
	case Seed:
		return "seed"
	default:
		return fmt.Sprintf("decision(%d)", int(d))
	}
}

// PrivateState is the subset of the framework's private-state type the digests
// need. The concrete type lives in an internal framework package, so it cannot
// be named here; the request and response Private fields satisfy this
// interface.
type PrivateState interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// Value returns the sha256 hex digest of one write-only value: a string is
// hashed as-is (so an importer can use sha256 of the value directly), anything
// else through its canonical JSON form, with keys sorted at every level. A nil
// value returns "", meaning not set.
func Value(v any) (string, error) {
	if v == nil {
		return "", nil
	}

	var raw []byte

	switch t := v.(type) {
	case string:
		raw = []byte(t)
	case *string:
		if t == nil {
			return "", nil
		}

		raw = []byte(*t)
	default:
		var err error

		raw, err = json.Marshal(normalise(v))
		if err != nil {
			return "", err
		}
	}

	sum := sha256.Sum256(raw)

	return hex.EncodeToString(sum[:]), nil
}

// ValueOf is Value for a framework attribute value: a string is hashed
// as-is; a dynamic value through its underlying value; anything else through
// the JSON form of its Terraform value (objects with sorted keys, numbers as
// JSON numbers), which is what Terraform's jsonencode() produces for the same
// value -- so an importer can reproduce the digest with sha256(jsonencode(v)).
func ValueOf(ctx context.Context, v attr.Value) (string, error) {
	if v == nil || v.IsNull() || v.IsUnknown() {
		return "", nil
	}

	if d, ok := v.(basetypes.DynamicValue); ok {
		return ValueOf(ctx, d.UnderlyingValue())
	}

	if s, ok := v.(basetypes.StringValue); ok {
		return Value(s.ValueString())
	}

	tv, err := v.ToTerraformValue(ctx)
	if err != nil {
		return "", err
	}

	var raw any
	if err := tftypesToGo(tv, &raw); err != nil {
		return "", err
	}

	return Value(raw)
}

// tftypesToGo converts a tftypes.Value into plain Go values (map[string]any,
// []any, string, float64/int64 via json.Number, bool) so it can be hashed
// through the canonical JSON form. Numbers are rendered as JSON numbers
// without a trailing ".0" when integral, matching jsonencode().
func tftypesToGo(v tftypes.Value, out *any) error {
	if v.IsNull() {
		*out = nil

		return nil
	}

	switch {
	case v.Type().Is(tftypes.String):
		var s string
		if err := v.As(&s); err != nil {
			return err
		}

		*out = s
	case v.Type().Is(tftypes.Bool):
		var b bool
		if err := v.As(&b); err != nil {
			return err
		}

		*out = b
	case v.Type().Is(tftypes.Number):
		// Terraform numbers are arbitrary precision; jsonencode() renders them
		// exactly. Every write-only number in the provider is an integer id,
		// rendered here as int64 -- identical to jsonencode(). A fractional or
		// beyond-int64 value is rendered through its shortest float64 text,
		// which may differ from jsonencode(): the digests then disagree and
		// the cost is one harmless re-send, never a missed change.
		var f *big.Float
		if err := v.As(&f); err != nil {
			return err
		}

		if i, acc := f.Int64(); acc == big.Exact {
			*out = i
		} else {
			fl, _ := f.Float64()
			*out = fl
		}
	case v.Type().Is(tftypes.Object{}) || v.Type().Is(tftypes.Map{}):
		var m map[string]tftypes.Value
		if err := v.As(&m); err != nil {
			return err
		}

		res := make(map[string]any, len(m))

		for k, e := range m {
			var ev any
			if err := tftypesToGo(e, &ev); err != nil {
				return err
			}

			res[k] = ev
		}

		*out = res
	case v.Type().Is(tftypes.List{}) || v.Type().Is(tftypes.Set{}) || v.Type().Is(tftypes.Tuple{}):
		var l []tftypes.Value
		if err := v.As(&l); err != nil {
			return err
		}

		res := make([]any, len(l))

		for i, e := range l {
			if err := tftypesToGo(e, &res[i]); err != nil {
				return err
			}
		}

		*out = res
	default:
		return fmt.Errorf("unsupported value type %s", v.Type())
	}

	return nil
}

// Values hashes each value in m with Value, dropping unset ones.
func Values(m map[string]any) (Attrs, error) {
	attrs := make(Attrs, len(m))

	for k, v := range m {
		h, err := Value(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", k, err)
		}

		if h != "" {
			attrs[k] = h
		}
	}

	return attrs, nil
}

// Digest returns the digest of an attribute digest set: sha256 hex of its
// canonical JSON form.
func Digest(attrs Attrs) string {
	if attrs == nil {
		attrs = Attrs{}
	}

	raw, _ := json.Marshal(attrs) // a map[string]string cannot fail to marshal
	sum := sha256.Sum256(raw)

	return hex.EncodeToString(sum[:])
}

// Compute returns the digests for an attribute digest set. immutable names the
// attributes the API will not change after create; nil or empty when there are
// none.
func Compute(attrs Attrs, immutable []string) Digests {
	d := Digests{All: Digest(attrs), Attrs: attrs, Version: Version}

	if len(immutable) > 0 {
		sub := Attrs{}

		for _, k := range immutable {
			if h, ok := attrs[k]; ok {
				sub[k] = h
			}
		}

		d.Immutable = Digest(sub)
	}

	return d
}

// ComputeValues is Compute over raw values: Values then Compute.
func ComputeValues(values map[string]any, immutable []string) (Digests, error) {
	attrs, err := Values(values)
	if err != nil {
		return Digests{}, err
	}

	return Compute(attrs, immutable), nil
}

// Decide compares the stored digests with the current ones.
func Decide(stored Digests, ok bool, current Digests) Decision {
	switch {
	case !ok:
		return Seed
	case stored.Sentinel == ImportSentinel:
		return Unchanged
	case stored.All == current.All:
		return Unchanged
	case stored.Immutable != current.Immutable:
		return Replace
	default:
		return Update
	}
}

// DecideConfigured is Decide for a resource whose write-only attributes are
// independent optional values (passwords, keys) rather than one configuration
// map: only the attributes present in current are compared, so an attribute
// left out of the configuration is not managed and never a change. immutable
// names the attributes whose change requires replacement.
func DecideConfigured(stored Digests, ok bool, current Attrs, immutable []string) Decision {
	switch {
	case !ok:
		return Seed
	case stored.Sentinel == ImportSentinel:
		return Unchanged
	case stored.Attrs == nil:
		// A record without per-attribute digests: compare whole.
		cur := Compute(current, immutable)
		if stored.Immutable == "" {
			cur.Immutable = ""
		}

		return Decide(stored, ok, cur)
	}

	decision := Unchanged

	for k, h := range current {
		if stored.Attrs[k] == h {
			continue
		}

		if contains(immutable, k) {
			return Replace
		}

		decision = Update
	}

	return decision
}

// Load reads the digests stored under key. ok is false when none are stored,
// or when they were written by a different digest version -- either way the
// caller has nothing to compare against.
func Load(ctx context.Context, private PrivateState, key string) (d Digests, ok bool, diags diag.Diagnostics) {
	if private == nil {
		return Digests{}, false, nil
	}

	raw, diags := private.GetKey(ctx, key)
	if diags.HasError() || len(raw) == 0 {
		return Digests{}, false, diags
	}

	if err := json.Unmarshal(raw, &d); err != nil || d.Version != Version {
		return Digests{}, false, diags
	}

	if d.All == "" && d.Sentinel == "" {
		return Digests{}, false, diags
	}

	return d, true, diags
}

// Store writes the digests under key.
func Store(ctx context.Context, private PrivateState, key string, d Digests) diag.Diagnostics {
	if private == nil {
		var diags diag.Diagnostics
		diags.AddError("write-only digests", "private state is unavailable; the change to the "+
			"write-only values could not be recorded")

		return diags
	}

	if d.Version == 0 {
		d.Version = Version
	}

	raw, err := json.Marshal(d)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("write-only digests", "encoding: "+err.Error())

		return diags
	}

	return private.SetKey(ctx, key, raw)
}

// StoreImportSentinel records, at import, that the configured write-only
// values are accepted as applied until the first apply. A write-only value
// edited between the import and that apply is therefore not detected -- the
// same blind spot a version marker has. An import ID that carries the digests
// (ParseImportID) avoids it.
func StoreImportSentinel(ctx context.Context, private PrivateState, key string) diag.Diagnostics {
	return Store(ctx, private, key, Digests{Sentinel: ImportSentinel, Version: Version})
}

// NewRevision is the write_only_revision value for an apply that sent a
// write-only value: the current time in seconds since the Unix epoch.
func NewRevision() types.Int64 {
	return types.Int64Value(time.Now().Unix())
}

// ImportIDSeparator separates the resource ID from the write-only digests in
// an import ID.
const ImportIDSeparator = "#"

var hexDigest = regexp.MustCompile(`^[0-9a-f]{64}$`)

// ParseImportID splits an import ID of the form
//
//	<id>
//	<id>#<attribute>=<sha256 hex>[&<attribute>=<sha256 hex>...]
//
// into the resource ID and the write-only attribute digests it carries (nil
// when it carries none). known lists the write-only attributes the resource
// has; any other attribute, a malformed digest, or a repeated attribute is an
// error, so a typo cannot silently import as "not set".
func ParseImportID(raw string, known []string) (id string, attrs Attrs, err error) {
	id, fragment, found := strings.Cut(raw, ImportIDSeparator)
	if !found {
		return id, nil, nil
	}

	if fragment == "" {
		return "", nil, fmt.Errorf("import ID %q: nothing after %q", raw, ImportIDSeparator)
	}

	attrs = Attrs{}

	for _, pair := range strings.Split(fragment, "&") {
		name, value, ok := strings.Cut(pair, "=")
		if !ok || name == "" {
			return "", nil, fmt.Errorf("import ID %q: %q is not <attribute>=<sha256>", raw, pair)
		}

		if !contains(known, name) {
			return "", nil, fmt.Errorf("import ID %q: %q is not a write-only attribute of this resource (expected one of %s)",
				raw, name, strings.Join(known, ", "))
		}

		if _, dup := attrs[name]; dup {
			return "", nil, fmt.Errorf("import ID %q: %q given twice", raw, name)
		}

		value = strings.ToLower(value)
		if !hexDigest.MatchString(value) {
			return "", nil, fmt.Errorf("import ID %q: the digest for %q is not 64 hexadecimal characters (sha256)", raw, name)
		}

		attrs[name] = value
	}

	return id, attrs, nil
}

// FormatImportID is the inverse of ParseImportID: "<id>#a=<hex>&b=<hex>", with
// attributes in sorted order, or just "<id>" when attrs is empty.
func FormatImportID(id string, attrs Attrs) string {
	if len(attrs) == 0 {
		return id
	}

	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+attrs[k])
	}

	return id + ImportIDSeparator + strings.Join(pairs, "&")
}

func contains(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}

	return false
}

// normalise sorts map keys at every level and recurses into slices, so a value
// that arrives as map[string]any or []any hashes the same however it was
// built. encoding/json sorts top-level map keys already; nested ones need this.
func normalise(v any) any {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		out := make(map[string]any, len(t))
		for _, k := range keys {
			out[k] = normalise(t[k])
		}

		return out
	case []any:
		out := make([]any, len(t))
		for i, e := range t {
			out[i] = normalise(e)
		}

		return out
	default:
		return v
	}
}
