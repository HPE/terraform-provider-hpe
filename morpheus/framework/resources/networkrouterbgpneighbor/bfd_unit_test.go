// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestBfdTimers guards the create failure where NSX-T rejected the body with
// the opaque "General error has occurred.".
//
// bfd_interval and bfd_multiple are Optional+Computed with no schema default,
// so omitting them left them unknown and they were not sent. Morpheus builds
// the outbound NSX-T bfd object unconditionally and does not strip nulls on
// create, producing "bfd":{"enabled":false,"interval":null,"multiple":null}.
func TestBfdTimers(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		interval, multiple         types.Int64
		wantInterval, wantMultiple int64
	}{
		"omitted from configuration (unknown)": {
			interval:     types.Int64Unknown(),
			multiple:     types.Int64Unknown(),
			wantInterval: defaultBfdIntervalMs,
			wantMultiple: defaultBfdMultiple,
		},
		"null": {
			interval:     types.Int64Null(),
			multiple:     types.Int64Null(),
			wantInterval: defaultBfdIntervalMs,
			wantMultiple: defaultBfdMultiple,
		},
		"configured values win": {
			interval:     types.Int64Value(500),
			multiple:     types.Int64Value(5),
			wantInterval: 500,
			wantMultiple: 5,
		},
		"partially configured": {
			interval:     types.Int64Value(300),
			multiple:     types.Int64Unknown(),
			wantInterval: 300,
			wantMultiple: defaultBfdMultiple,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			interval, multiple := bfdTimers(tc.interval, tc.multiple)

			if interval == nil || multiple == nil {
				t.Fatal("BFD timers must always be sent, never nil")

				return
			}

			if *interval != tc.wantInterval {
				t.Errorf("bfd_interval = %d, want %d", *interval, tc.wantInterval)
			}

			if *multiple != tc.wantMultiple {
				t.Errorf("bfd_multiple = %d, want %d", *multiple, tc.wantMultiple)
			}
		})
	}
}
