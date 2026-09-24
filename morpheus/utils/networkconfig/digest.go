// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkconfig holds the pieces of hpe_morpheus_network's typed,
// write-only configuration that do not depend on the Terraform schema: the
// catalogue of what each network type's configuration fields are, the routing
// of a configured block into the config map Morpheus reads, and the digests
// through which a write-only block's changes are detected.
//
// Morpheus keeps a network type's configuration in two places. Fields whose
// option type has fieldContext "network" are first-class columns of the
// network and are sent, and returned, at the top level of the record; fields
// with fieldContext "config" go into the network's config map and come back
// under "config". The provider mirrors that split: the top-level attributes are
// ordinary, readable attributes; the config-map fields live in write-only
// blocks, one per network type, plus the free-form write-only "config" for
// types without a block.
//
// What a type's fields are, and which of them Morpheus still edits after
// create, comes from the appliance the provider is talking to: at plan and
// apply time the resource reads the network type (GET /api/network-types) and
// FromAPI normalises it into a TypeSpec, so the provider follows the release
// and the plugins that appliance actually has. The generated table in
// known_gen.go — the same normalisation, run by the tooling's harvester
// against a reference appliance — is the fallback when the appliance cannot be
// read, and the fixture the hand-written typed blocks are checked against in
// the unit tests. A type in neither place gets the generic behaviour.
//
// Because a write-only value is never in state, Terraform cannot see it change.
// HashiCorp's documented answer, and the one used here, is to keep a secure
// digest of the value in the resource's private state and compare it during
// plan modification: a change to an editable field marks the read-only
// config_from_api unknown, which plans an update; a change to a field the
// server will not edit after creation requires replacement, since Morpheus
// applies only editable option types on update and would otherwise ignore the
// change without a word.
package networkconfig

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// PrivateKey is the key under which the digests are kept in private state.
const PrivateKey = "networkconfig.digests"

// digestVersion changes when the canonical form changes, so a digest written by
// an older provider is treated as missing rather than compared.
const digestVersion = 1

// Digests are the two sha256 hex digests kept in private state.
type Digests struct {
	// All covers every write-only configuration value as sent to Morpheus.
	All string `json:"all"`
	// Immutable covers only the values Morpheus will not change after create.
	Immutable string `json:"immutable"`
	Version   int    `json:"v"`
}

// privateState is the subset of the framework's private-state type the digests
// need. The concrete type lives in an internal framework package, so it cannot
// be named here; the request and response Private fields satisfy this
// interface, which the resource's use of Load and Store establishes at compile
// time.
type privateState interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// Digest returns the sha256 hex digest of the canonical JSON form of m. Keys are
// sorted, so two maps with the same content digest equally; the values are what
// the provider sends, so a digest changes exactly when the request would.
func Digest(m map[string]any) (string, error) {
	canonical, err := canonicalJSON(m)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(canonical)

	return hex.EncodeToString(sum[:]), nil
}

// Compute returns the digests for a full config map and its immutable subset.
func Compute(all, immutable map[string]any) (Digests, error) {
	allDigest, err := Digest(all)
	if err != nil {
		return Digests{}, fmt.Errorf("digest of configuration: %w", err)
	}

	immutableDigest, err := Digest(immutable)
	if err != nil {
		return Digests{}, fmt.Errorf("digest of immutable configuration: %w", err)
	}

	return Digests{All: allDigest, Immutable: immutableDigest, Version: digestVersion}, nil
}

// Load reads the digests from private state. ok is false when none are stored,
// or when they were written by a different digest version — either way the
// caller has nothing to compare against.
func Load(ctx context.Context, private privateState) (d Digests, ok bool, diags diag.Diagnostics) {
	if private == nil {
		return Digests{}, false, nil
	}

	raw, diags := private.GetKey(ctx, PrivateKey)
	if diags.HasError() || len(raw) == 0 {
		return Digests{}, false, diags
	}

	if err := json.Unmarshal(raw, &d); err != nil || d.Version != digestVersion || d.All == "" {
		return Digests{}, false, diags
	}

	return d, true, diags
}

// Store writes the digests to private state.
func Store(ctx context.Context, private privateState, d Digests) diag.Diagnostics {
	if private == nil {
		var diags diag.Diagnostics
		diags.AddError("network configuration digests", "private state is unavailable; the change to the "+
			"write-only configuration could not be recorded")

		return diags
	}

	raw, err := json.Marshal(d)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("network configuration digests", "encoding: "+err.Error())

		return diags
	}

	return private.SetKey(ctx, PrivateKey, raw)
}

// canonicalJSON marshals m with sorted keys at every level. encoding/json sorts
// map keys already; this normalises nested maps and slices so a value that
// arrives as map[string]any or []any hashes the same however it was built.
func canonicalJSON(m map[string]any) ([]byte, error) {
	return json.Marshal(normalise(m))
}

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
