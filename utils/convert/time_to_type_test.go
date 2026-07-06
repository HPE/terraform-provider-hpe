// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package convert

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestTimeToType(t *testing.T) {
	t.Parallel()

	utc := time.Date(2026, time.July, 6, 13, 45, 30, 0, time.UTC)
	offset := time.FixedZone("UTC+2", 2*60*60)
	withOffset := time.Date(2026, time.July, 6, 13, 45, 30, 0, offset)

	tests := []struct {
		name     string
		input    *time.Time
		expected types.String
	}{
		{
			name:     "nil returns null",
			input:    nil,
			expected: types.StringNull(),
		},
		{
			name:     "utc formats as RFC3339",
			input:    &utc,
			expected: types.StringValue("2026-07-06T13:45:30Z"),
		},
		{
			name:     "non-utc preserves offset",
			input:    &withOffset,
			expected: types.StringValue("2026-07-06T13:45:30+02:00"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := TimeToType(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
