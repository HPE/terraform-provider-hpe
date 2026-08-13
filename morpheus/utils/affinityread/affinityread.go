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
