// Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// The read response models config loosely, so the provisioning image arrives in
// AdditionalProperties rather than as a typed field. Morpheus resolves the image
// from imageId first and template second, and template has a map form that the
// UI typeahead submits, so all of those shapes have to survive a read.
//
// Getting this wrong is not a cosmetic bug: image_id forces replacement, so a
// value that fails to read back leaves the next plan proposing to destroy and
// recreate the instance after an import.
func TestImageIDFromConfig(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		props map[string]interface{}
		want  int64
		null  bool
	}{
		{
			name:  "absent",
			props: map[string]interface{}{},
			null:  true,
		},
		{
			// The shape a live HVM appliance returns, verified against Zodiac:
			// imageId echoed back verbatim as a JSON number.
			name:  "imageId as float64",
			props: map[string]interface{}{"imageId": float64(1405)},
			want:  1405,
		},
		{
			name:  "imageId as string",
			props: map[string]interface{}{"imageId": "1405"},
			want:  1405,
		},
		{
			name:  "template scalar fallback",
			props: map[string]interface{}{"template": float64(567)},
			want:  567,
		},
		{
			// The UI typeahead form of template.
			name:  "template map fallback",
			props: map[string]interface{}{"template": map[string]interface{}{"value": float64(567)}},
			want:  567,
		},
		{
			// imageId wins, matching getContainerVirtualImageId's precedence.
			name: "imageId takes precedence over template",
			props: map[string]interface{}{
				"imageId":  float64(1405),
				"template": float64(567),
			},
			want: 1405,
		},
		{
			name:  "unusable value is null, not zero",
			props: map[string]interface{}{"imageId": "not-a-number"},
			null:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := &sdk.GetInstance200ResponseInstanceConfig{
				AdditionalProperties: tc.props,
			}

			got := imageIDFromConfig(cfg)

			if tc.null {
				if !got.IsNull() {
					t.Fatalf("expected null, got %v", got)
				}

				return
			}

			if got.IsNull() {
				t.Fatalf("expected %d, got null", tc.want)
			}

			if got.ValueInt64() != tc.want {
				t.Fatalf("expected %d, got %d", tc.want, got.ValueInt64())
			}
		})
	}
}
