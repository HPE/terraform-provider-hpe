// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Platform defaults for the BFD timers.
//
// These mirror the defaults the Morpheus UI applies when creating an NSX-T BGP
// neighbour, which it always sends with the request.
const (
	defaultBfdIntervalMs = 1000
	defaultBfdMultiple   = 3
)

// bfdTimers resolves the bfd_interval and bfd_multiple values to send.
//
// Both attributes are Optional+Computed with no schema default, so omitting
// them from the configuration leaves them unknown and they were previously not
// sent at all. The API always forwards a BFD object to NSX-T and does not drop
// unset members on create, so omitting them produces
//
//	"bfd": {"enabled": false, "interval": null, "multiple": null}
//
// which NSX-T rejects during body validation with the opaque
// "General error has occurred." -- masking any real semantic error underneath.
// Creating a neighbour through the UI never hits this, because the UI defaults
// are always sent, so the provider sends them too.
func bfdTimers(interval, multiple types.Int64) (intervalOut, multipleOut *int64) {
	i := int64(defaultBfdIntervalMs)
	if !interval.IsNull() && !interval.IsUnknown() {
		i = interval.ValueInt64()
	}

	m := int64(defaultBfdMultiple)
	if !multiple.IsNull() && !multiple.IsUnknown() {
		m = multiple.ValueInt64()
	}

	return &i, &m
}
