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
// The change is detected through the digests kept in private state by the
// shared wodigest package: a change to an editable field marks the read-only
// config_from_api unknown, which plans an update; a change to a field the
// server will not edit after creation requires replacement, since Morpheus
// applies only editable option types on update and would otherwise ignore the
// change without a word.
package networkconfig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// PrivateKey is the key under which the digests are kept in private state.
const PrivateKey = "networkconfig.digests"

// Digests are the two sha256 hex digests kept in private state; see wodigest.
type Digests = wodigest.Digests

// Digest returns the digest of a config map: each value hashed on its own,
// then the {key: digest} set hashed again (see wodigest).
func Digest(m map[string]any) (string, error) {
	attrs, err := wodigest.Values(m)
	if err != nil {
		return "", err
	}

	return wodigest.Digest(attrs), nil
}

// Compute returns the digests for a full config map and its immutable subset
// (the immutable map's keys name the immutable attributes).
func Compute(all, immutable map[string]any) (Digests, error) {
	keys := make([]string, 0, len(immutable))
	for k := range immutable {
		keys = append(keys, k)
	}

	d, err := wodigest.ComputeValues(all, keys)
	if err != nil {
		return Digests{}, err
	}

	if d.Immutable == "" {
		// The network resource always compares the immutable digest, so an
		// empty subset digests as the empty set rather than as "none".
		d.Immutable = wodigest.Digest(nil)
	}

	return d, nil
}

// Load reads the digests from private state. ok is false when none are stored,
// or when they were written by a different digest version.
func Load(ctx context.Context, private wodigest.PrivateState) (Digests, bool, diag.Diagnostics) {
	return wodigest.Load(ctx, private, PrivateKey)
}

// Store writes the digests to private state.
func Store(ctx context.Context, private wodigest.PrivateState, d Digests) diag.Diagnostics {
	return wodigest.Store(ctx, private, PrivateKey, d)
}
