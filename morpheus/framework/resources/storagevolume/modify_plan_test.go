// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

// number builds a tftypes numeric value for an int64.
func number(n int64) tftypes.Value {
	return tftypes.NewValue(tftypes.Number, new(big.Float).SetInt64(n))
}

// TestUnitStorageVolumeModifyPlanTypeIDSize verifies that ModifyPlan resolves a
// type_id to its code and enforces the Alletra 1-65536 GiB max_storage bound at
// plan time, while skipping the API call when type_code is supplied or
// max_storage is absent.
func TestUnitStorageVolumeModifyPlanTypeIDSize(t *testing.T) {
	ctx := context.Background()
	sch := StorageVolumeResourceSchema(ctx)
	objType := sch.Type().TerraformType(ctx).(tftypes.Object)

	plan := func(overrides map[string]tftypes.Value) tfsdk.Plan {
		vals := map[string]tftypes.Value{}
		for name, at := range objType.AttributeTypes {
			vals[name] = tftypes.NewValue(at, nil)
		}
		for name, v := range overrides {
			vals[name] = v
		}

		return tfsdk.Plan{Schema: sch, Raw: tftypes.NewValue(objType, vals)}
	}

	cases := []struct {
		name      string
		overrides map[string]tftypes.Value
		wantError bool
		wantCall  bool
	}{
		{
			name: "type_id alletra oversize is rejected at plan",
			overrides: map[string]tftypes.Value{
				"type_id":     number(1),
				"max_storage": number(99999),
			},
			wantError: true,
			wantCall:  true,
		},
		{
			name: "type_id alletra in range is accepted",
			overrides: map[string]tftypes.Value{
				"type_id":     number(1),
				"max_storage": number(100),
			},
			wantError: false,
			wantCall:  true,
		},
		{
			name: "type_code path skips the api call",
			overrides: map[string]tftypes.Value{
				"type_code":   tftypes.NewValue(tftypes.String, "hpealletraMPLUN"),
				"max_storage": number(99999),
			},
			wantError: false,
			wantCall:  false,
		},
		{
			name: "no max_storage skips the api call",
			overrides: map[string]tftypes.Value{
				"type_id": number(1),
			},
			wantError: false,
			wantCall:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var called bool
			srv := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					called = true
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(
						`{"storageVolumeType":{"id":1,"code":"hpealletraMPLUN"}}`))
				}))
			defer srv.Close()

			r := &storageVolumeResource{}
			cfgResp := &resource.ConfigureResponse{}
			r.Configure(ctx, resource.ConfigureRequest{
				ProviderData: clientfactory.New(model.MorpheusProviderModel{
					URL:         types.StringValue(srv.URL),
					AccessToken: types.StringValue("test-token"),
				}),
			}, cfgResp)
			if cfgResp.Diagnostics.HasError() {
				t.Fatalf("configure: %v", cfgResp.Diagnostics)
			}

			p := plan(tc.overrides)
			resp := &resource.ModifyPlanResponse{Plan: p}
			r.ModifyPlan(ctx, resource.ModifyPlanRequest{Plan: p}, resp)

			if got := resp.Diagnostics.HasError(); got != tc.wantError {
				t.Fatalf("ModifyPlan error=%v want %v (diags: %v)",
					got, tc.wantError, resp.Diagnostics)
			}
			if called != tc.wantCall {
				t.Fatalf("api called=%v want %v", called, tc.wantCall)
			}
		})
	}
}
