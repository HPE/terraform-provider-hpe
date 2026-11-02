// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagebucket

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

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

			body := sdk.AddStorageBucketsRequestStorageBucket{
				Name:         "example",
				ProviderType: "s3",
				Config:       addBucketConfig(types.StringValue("s3"), types.StringNull(), tc.endpoint, tc.accessKey, tc.secretKey),
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

// TestUnitLocalBucketConfig: a local bucket's config is the basePath variant
// on create and on update -- never the S3-shaped one, whose empty object would
// leave Morpheus a bucket with no directory -- and every other provider type
// keeps the S3 shape.
func TestUnitLocalBucketConfig(t *testing.T) {
	t.Parallel()

	local := types.StringValue(providerTypeLocal)
	basePath := types.StringValue("/var/opt/morpheus/morpheus-ui/images")

	t.Run("create sends basePath for local", func(t *testing.T) {
		t.Parallel()

		raw, err := json.Marshal(sdk.AddStorageBucketsRequest{StorageBucket: sdk.AddStorageBucketsRequestStorageBucket{
			Name:         "images",
			ProviderType: providerTypeLocal,
			Config:       addBucketConfig(local, basePath, types.StringNull(), types.StringNull(), types.StringNull()),
		}})
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}

		var got struct {
			StorageBucket struct {
				Config map[string]any `json:"config"`
			} `json:"storageBucket"`
		}
		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		if want := (map[string]any{"basePath": "/var/opt/morpheus/morpheus-ui/images"}); !equalMaps(got.StorageBucket.Config, want) {
			t.Errorf("config = %v, want %v", got.StorageBucket.Config, want)
		}
	})

	t.Run("update sends basePath for local and ignores credentials", func(t *testing.T) {
		t.Parallel()

		ak, ep := "AKIAEXAMPLE", "https://s3.example.com"
		raw, err := json.Marshal(sdk.UpdateStorageBucketsRequest{StorageBucket: sdk.UpdateStorageBucketsRequestStorageBucket{
			Config: updateBucketConfig(local, basePath, &ak, nil, &ep),
		}})
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}

		var got struct {
			StorageBucket struct {
				Config map[string]any `json:"config"`
			} `json:"storageBucket"`
		}
		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		if want := (map[string]any{"basePath": "/var/opt/morpheus/morpheus-ui/images"}); !equalMaps(got.StorageBucket.Config, want) {
			t.Errorf("config = %v, want %v", got.StorageBucket.Config, want)
		}
	})

	t.Run("update keeps the S3 shape for other types", func(t *testing.T) {
		t.Parallel()

		ep := "https://s3.example.com"
		raw, err := json.Marshal(sdk.UpdateStorageBucketsRequest{StorageBucket: sdk.UpdateStorageBucketsRequestStorageBucket{
			Config: updateBucketConfig(types.StringValue("s3"), types.StringNull(), nil, nil, &ep),
		}})
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}

		var got struct {
			StorageBucket struct {
				Config map[string]any `json:"config"`
			} `json:"storageBucket"`
		}
		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		if want := (map[string]any{"endpoint": "https://s3.example.com"}); !equalMaps(got.StorageBucket.Config, want) {
			t.Errorf("config = %v, want %v", got.StorageBucket.Config, want)
		}
	})
}

// TestUnitBasePathFromResponse: the response's config model is typed for the
// S3 shape, so a local bucket's basePath arrives as an additional property;
// it is read back for local buckets only, and an absent or empty one is null.
func TestUnitBasePathFromResponse(t *testing.T) {
	t.Parallel()

	local, s3 := providerTypeLocal, "s3"
	withPath := &sdk.GetStorageBuckets200ResponseStorageBucketConfig{
		AdditionalProperties: map[string]any{"basePath": "/srv/images"},
	}
	emptyPath := &sdk.GetStorageBuckets200ResponseStorageBucketConfig{
		AdditionalProperties: map[string]any{"basePath": ""},
	}

	cases := map[string]struct {
		sb   sdk.GetStorageBuckets200ResponseStorageBucket
		want types.String
	}{
		"local with basePath": {
			sdk.GetStorageBuckets200ResponseStorageBucket{ProviderType: &local, Config: withPath},
			types.StringValue("/srv/images"),
		},
		"local without config": {
			sdk.GetStorageBuckets200ResponseStorageBucket{ProviderType: &local},
			types.StringNull(),
		},
		"local with empty basePath": {
			sdk.GetStorageBuckets200ResponseStorageBucket{ProviderType: &local, Config: emptyPath},
			types.StringNull(),
		},
		"s3 carrying a stray basePath": {
			sdk.GetStorageBuckets200ResponseStorageBucket{ProviderType: &s3, Config: withPath},
			types.StringNull(),
		},
		"no provider type": {
			sdk.GetStorageBuckets200ResponseStorageBucket{Config: withPath},
			types.StringNull(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := basePathFromResponse(&tc.sb); !got.Equal(tc.want) {
				t.Errorf("basePathFromResponse = %v, want %v", got, tc.want)
			}
		})
	}
}

func equalMaps(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}

	return true
}

// TestUnitDefaultVirtualImageTargetKeepsState pins the plan behaviour of an
// unset default_virtual_image_target: the state value carries over, so a
// bucket that is the appliance's default image store stays so when the
// attribute is not configured. The spec's request body declares a default of
// false for the field; a schema default built from it would plan true ->
// false for such a bucket, and the update would take the appliance's default
// image store away (the regression PR #142's review caught).
func TestUnitDefaultVirtualImageTargetKeepsState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := StorageBucketResourceSchema(ctx)

	attr, ok := s.Attributes["default_virtual_image_target"].(schema.BoolAttribute)
	if !ok {
		t.Fatal("default_virtual_image_target is not a BoolAttribute")
	}

	if attr.Default != nil {
		t.Fatal("default_virtual_image_target has a schema default; it must not")
	}

	// An existing resource: the modifiers treat a null state as a create and
	// do nothing, so the request carries a (non-null) state object.
	objType, ok := s.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}

	attrs := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, at := range objType.AttributeTypes {
		attrs[name] = tftypes.NewValue(at, nil)
	}

	// An unset Optional+Computed attribute is unknown in the plan; the plan
	// modifiers must resolve it to the state value.
	req := planmodifier.BoolRequest{
		State:       tfsdk.State{Schema: s, Raw: tftypes.NewValue(objType, attrs)},
		ConfigValue: types.BoolNull(),
		StateValue:  types.BoolValue(true),
		PlanValue:   types.BoolUnknown(),
	}
	resp := planmodifier.BoolResponse{PlanValue: req.PlanValue}

	for _, m := range attr.PlanModifiers {
		m.PlanModifyBool(ctx, req, &resp)
	}

	if resp.Diagnostics.HasError() {
		t.Fatalf("plan modifiers: %v", resp.Diagnostics)
	}

	if !resp.PlanValue.Equal(types.BoolValue(true)) {
		t.Fatalf("planned default_virtual_image_target = %v, want the state value true", resp.PlanValue)
	}
}
