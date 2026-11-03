// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// digestKey is the private-state key under which the write-only key digest is
// stored.
const digestKey = "license.wo_digests"

// licenseKeyID is the short identifier Morpheus uses in the API for an
// installed licence: the first eight characters of the trimmed key. Ported
// from the SDKv2 resource.
func licenseKeyID(key string) string {
	k := strings.TrimSpace(key)
	if len(k) > 8 {
		return k[:8]
	}

	return k
}

// keyDigest is the hex sha256 of the trimmed key -- identical to Morpheus's
// own license "hash" (ApplianceLicenseService.getLicenseKeyHash). That
// identity is what makes the installed-licence list an oracle for the digest:
// Read can record the digest from the API's hash without ever seeing the key.
func keyDigest(key string) string {
	// Value never returns an error for a string argument.
	h, _ := wodigest.Value(strings.TrimSpace(key))

	return h
}

// digestsFor builds the Digests record stored in private state for a licence
// whose key hashes to hash. The hash is lowercased to match the API, which
// reports it in lowercase hex.
func digestsFor(hash string) wodigest.Digests {
	return wodigest.Compute(wodigest.Attrs{"key": strings.ToLower(hash)}, nil)
}

// storedDigest reads the key digest kept in private state. ok is true only
// when a non-empty digest is present and no import sentinel is set. Load
// diagnostics are appended to diags (a pointer argument, not a result, so this
// helper is not mistaken for a Read not-found helper by the read-rules test).
func storedDigest(
	ctx context.Context,
	private wodigest.PrivateState,
	diags *diag.Diagnostics,
) (hash string, ok bool) {
	d, loaded, loadDiags := wodigest.Load(ctx, private, digestKey)
	diags.Append(loadDiags...)
	if !loaded || d.Sentinel != "" {
		return "", false
	}

	h := d.Attrs["key"]
	if h == "" {
		return "", false
	}

	return h, true
}

// storeDigest records the key digest for hash in private state.
func storeDigest(
	ctx context.Context,
	private wodigest.PrivateState,
	hash string,
) diag.Diagnostics {
	return wodigest.Store(ctx, private, digestKey, digestsFor(hash))
}

// planDecision decides what a plan should do about the write-only key given
// the digest stored in private state (stored/ok) and the digest of the
// currently configured key. It is a thin, unit-testable wrapper over
// wodigest.DecideConfigured: no stored digest -> Seed; equal (case-insensitive)
// -> Unchanged; otherwise Update. The key is never immutable, so Replace is
// never returned.
func planDecision(stored string, ok bool, configured string) wodigest.Decision {
	return wodigest.DecideConfigured(
		wodigest.Digests{Attrs: wodigest.Attrs{"key": strings.ToLower(stored)}},
		ok,
		wodigest.Attrs{"key": strings.ToLower(configured)},
		nil,
	)
}
