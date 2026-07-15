// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

// TestUnitStorageVolumeModifyPlanTypeIDSize verifies that ModifyPlan resolves a
// type_id to its code and enforces the Alletra 1-65536 GiB max_storage bound at
// plan time, while skipping the API call when type_code is supplied or
// max_storage is absent.
func TestUnitStorageVolumeModifyPlanTypeIDSize(t *testing.T) {
	ctx := context.Background()
	sch := StorageVolumeResourceSchema(ctx)

	// makePlan builds a tfsdk.Plan from a typed model. The write-only config
	// block and dynamic config are set to explicit nulls so the plan matches the
	// schema; each case only populates the fields ModifyPlan actually reads.
	makePlan := func(t *testing.T, m StorageVolumeModel) tfsdk.Plan {
		t.Helper()
		m.ConfigAlletrampBmaas = NewConfigAlletrampBmaasValueNull()
		m.Config = types.DynamicNull()

		p := tfsdk.Plan{Schema: sch}
		if diags := p.Set(ctx, &m); diags.HasError() {
			t.Fatalf("build plan: %v", diags)
		}

		return p
	}

	cases := []struct {
		name      string
		model     StorageVolumeModel
		wantError bool
		wantCall  bool
	}{
		{
			name: "type_id alletra oversize is rejected at plan",
			model: StorageVolumeModel{
				TypeId:     types.Int64Value(1),
				MaxStorage: types.Int64Value(99999),
			},
			wantError: true,
			wantCall:  true,
		},
		{
			name: "type_id alletra in range is accepted",
			model: StorageVolumeModel{
				TypeId:     types.Int64Value(1),
				MaxStorage: types.Int64Value(100),
			},
			wantError: false,
			wantCall:  true,
		},
		{
			name: "type_code path skips the api call",
			model: StorageVolumeModel{
				TypeCode:   types.StringValue("hpealletraMPLUN"),
				MaxStorage: types.Int64Value(99999),
			},
			wantError: false,
			wantCall:  false,
		},
		{
			name: "no max_storage skips the api call",
			model: StorageVolumeModel{
				TypeId: types.Int64Value(1),
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

			p := makePlan(t, tc.model)
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
