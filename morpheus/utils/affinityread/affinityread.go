// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package affinityread works around a defect in the affinity group single-item
// endpoint.
//
// TODO(MORPH-15806): remove this package once the appliance defect is fixed.
//
// Supplying tenant permissions to an affinity group returns HTTP 500 and,
// from that point on, the group's single-item GET returns 500 permanently.
// The stored record is undamaged — the list endpoint continues to render the
// same group correctly — so the fault is only in serving the single response.
// Reproduced on appliance versions 9.0.0 and 9.0.2.95.
//
// Without a workaround, one such group makes every plan fail for the resources
// that reference it, and the group cannot be repaired: only deleted and
// recreated.
//
// The list endpoint is not a complete substitute. It returns 11 of the 13
// fields the single-item endpoint does, omitting resourcePermissions and
// tenants. For callers that need only membership, which is what this package
// serves, the fallback is lossless: servers is present and identical in both.
// Callers needing the full object must handle the two missing fields
// themselves, and should say so rather than silently reporting them as unset.
package affinityread

import (
	"net/http"
)

// IsSingleItemRenderFailure reports whether a response looks like the defect
// rather than an ordinary error.
//
// The check is deliberately narrow. A 500 from this endpoint is the observed
// signature, and widening it risks papering over unrelated server errors that
// should surface to the practitioner.
func IsSingleItemRenderFailure(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusInternalServerError
}

// ServersFromList extracts membership from a list response.
//
// ids maps each element to its affinity group id, and servers maps it to the
// member server ids; both are supplied by the caller because the generated
// list types differ between cloud and cluster affinity groups while their
// shapes do not.
func ServersFromList[T any](
	items []T,
	wantGroupID int64,
	id func(T) (int64, bool),
	servers func(T) []int64,
) ([]int64, bool) {
	for _, item := range items {
		got, ok := id(item)
		if !ok || got != wantGroupID {
			continue
		}

		found := servers(item)
		if found == nil {
			found = []int64{}
		}

		return found, true
	}

	return nil, false
}

// IDsFromList collects every affinity group id in a list response.
//
// Used to record what exists before a create that is expected to fail to
// render, so the new group can be identified afterwards by difference.
//
// The result is non-nil whenever the listing was read, so a nil map means the
// listing could not be read at all. That is a different thing from a cloud or
// cluster which legitimately holds no groups, and the two must not be
// conflated: one is recoverable, the other is not.
func IDsFromList[T any](items []T, id func(T) (int64, bool)) map[int64]struct{} {
	ids := make(map[int64]struct{}, len(items))

	for _, item := range items {
		if got, ok := id(item); ok {
			ids[got] = struct{}{}
		}
	}

	return ids
}

// NewIDFromList returns the one id present in items but absent from before.
//
// It reports false unless exactly one such id exists. If another client created
// a group at the same moment, more than one id is new and picking between them
// would be a guess; adopting the wrong group into state is worse than reporting
// that the group cannot be identified, because the practitioner can recover
// from an error but not from silently managing someone else's resource.
func NewIDFromList[T any](
	items []T,
	before map[int64]struct{},
	id func(T) (int64, bool),
) (int64, bool) {
	var (
		found int64
		count int
	)

	for _, item := range items {
		got, ok := id(item)
		if !ok {
			continue
		}

		if _, existed := before[got]; existed {
			continue
		}

		found = got
		count++
	}

	return found, count == 1
}
