// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagebucket

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestAddBucketConfigAlwaysMarshals guards the create failure reported as
//
//	json: error calling MarshalJSON for type
//	sdk.AddStorageBucketsRequestStorageBucketConfig: unexpected end of JSON input
//
// storageBucket.config is a non-pointer oneOf wrapper that ToMap always
// serialises, and the generated MarshalJSON returns (nil, nil) when no variant
// is set, which encoding/json rejects. A variant must therefore be selected
// even when the practitioner supplied no credentials at all.
func TestAddBucketConfigAlwaysMarshals(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		endpoint, accessKey, secretKey types.String
		wantConfig                     map[string]any
	}{
		"no credentials configured": {
			endpoint:   types.StringNull(),
			accessKey:  types.StringNull(),
			secretKey:  types.StringNull(),
			wantConfig: map[string]any{},
		},
		"unknown values are omitted": {
			endpoint:   types.StringUnknown(),
			accessKey:  types.StringUnknown(),
			secretKey:  types.StringUnknown(),
			wantConfig: map[string]any{},
		},
		"credentials are nested under config": {
			endpoint:  types.StringValue("https://s3.example.com"),
			accessKey: types.StringValue("AKIAEXAMPLE"),
			secretKey: types.StringValue("s3cr3t"),
			wantConfig: map[string]any{
				"endpoint":  "https://s3.example.com",
				"accessKey": "AKIAEXAMPLE",
				"secretKey": "s3cr3t",
			},
		},
		"endpoint only": {
			endpoint:   types.StringValue("https://s3.example.com"),
			accessKey:  types.StringNull(),
			secretKey:  types.StringNull(),
			wantConfig: map[string]any{"endpoint": "https://s3.example.com"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bucketConfig := addBucketConfig(tc.endpoint, tc.accessKey, tc.secretKey)
			body := sdk.AddStorageBucketsRequestStorageBucket{
				Name:         "example",
				ProviderType: "s3",
				Config:       &bucketConfig,
			}

			raw, err := json.Marshal(sdk.AddStorageBucketsRequest{StorageBucket: body})
			if err != nil {
				t.Fatalf("marshalling the create request failed: %v", err)
			}

			var got struct {
				StorageBucket struct {
					Name         string         `json:"name"`
					ProviderType string         `json:"providerType"`
					Config       map[string]any `json:"config"`
				} `json:"storageBucket"`
			}

			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatalf("unmarshalling the create request failed: %v", err)
			}

			if got.StorageBucket.Config == nil {
				t.Fatal("config must always be an object, never null")
			}

			if len(got.StorageBucket.Config) != len(tc.wantConfig) {
				t.Fatalf("config = %v, want %v", got.StorageBucket.Config, tc.wantConfig)
			}

			for k, want := range tc.wantConfig {
				if got.StorageBucket.Config[k] != want {
					t.Errorf("config[%q] = %v, want %v", k, got.StorageBucket.Config[k], want)
				}
			}

			// Credentials must not leak to the top level: Morpheus binds the
			// bucket with config excluded, so anything there is discarded.
			var top map[string]any
			if err := json.Unmarshal(raw, &top); err != nil {
				t.Fatalf("unmarshalling the create request failed: %v", err)
			}

			bucket, ok := top["storageBucket"].(map[string]any)
			if !ok {
				t.Fatal("storageBucket missing from the request body")
			}

			for _, k := range []string{"accessKey", "secretKey", "endpoint"} {
				if _, found := bucket[k]; found {
					t.Errorf("%q was sent at the top level; it belongs under config", k)
				}
			}
		})
	}
}
