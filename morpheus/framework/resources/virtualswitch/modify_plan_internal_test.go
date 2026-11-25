// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// planModel is a switch model with the attributes ModifyPlan reads; the rest
// stay null.
func planModel(t *testing.T, clusterID int64, uplinkServers ...int64) VirtualSwitchModel {
	t.Helper()

	ups := make([]UplinksValue, 0, len(uplinkServers))
	for _, id := range uplinkServers {
		ups = append(ups, newUplink(id, "eth2", nil))
	}

	return VirtualSwitchModel{
		ClusterId: types.Int64Value(clusterID),
		Name:      tfstr("vsw"),
		BaseType:  tfstr("general"),
		Mtu:       tfint(1500),
		Uplinks:   uplinkSet(t, ups),
		Segments:  buildSegmentsMapForTest(t, map[string]SegmentsValue{"mgmt": newSeg(t, 1, "management", nil)}),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType,
		})},
	}
}

func modifyPlan(t *testing.T, state, plan VirtualSwitchModel) resource.ModifyPlanResponse {
	t.Helper()

	ctx := context.Background()
	sch := VirtualSwitchResourceSchema(ctx)

	st := tfsdk.State{Schema: sch}
	if d := st.Set(ctx, &state); d.HasError() {
		t.Fatalf("state: %v", d)
	}

	pl := tfsdk.Plan{Schema: sch}
	if d := pl.Set(ctx, &plan); d.HasError() {
		t.Fatalf("plan: %v", d)
	}

	resp := resource.ModifyPlanResponse{Plan: pl}
	(&Resource{}).ModifyPlan(ctx, resource.ModifyPlanRequest{State: st, Plan: pl}, &resp)

	return resp
}

// Dropping an uplink is refused on an in-place update, but not on a plan that
// replaces the switch anyway (a changed cluster_id here): the replacement is
// judged first, so the removal rule does not fail a legitimate recreate.
func TestUnitModifyPlanUplinkRemovalVersusReplacement(t *testing.T) {
	t.Parallel()

	state := planModel(t, 4242, 7, 8)

	inPlace := modifyPlan(t, state, planModel(t, 4242, 7))
	if !inPlace.Diagnostics.HasError() || !strings.Contains(inPlace.Diagnostics.Errors()[0].Summary(), "Cannot remove a host uplink") {
		t.Fatalf("in-place removal: want the uplink error, got %v", inPlace.Diagnostics)
	}

	recreate := modifyPlan(t, state, planModel(t, 9999, 7))
	if recreate.Diagnostics.HasError() {
		t.Fatalf("replacement must not be refused for its dropped uplinks: %v", recreate.Diagnostics)
	}

	// The replacement itself is the schema's RequiresReplace on cluster_id;
	// ModifyPlan adds none of its own and leaves the plan as it was.
	if len(recreate.RequiresReplace) != 0 {
		t.Errorf("ModifyPlan added RequiresReplace %v; the schema's plan modifiers own the replacement", recreate.RequiresReplace)
	}

	if len(inPlace.RequiresReplace) != 0 {
		t.Errorf("a refused in-place plan must not also plan a replacement: %v", inPlace.RequiresReplace)
	}

	var planned VirtualSwitchModel
	if d := recreate.Plan.Get(context.Background(), &planned); d.HasError() {
		t.Fatalf("plan: %v", d)
	}

	if planned.ClusterId.ValueInt64() != 9999 || planned.Name.ValueString() != "vsw" {
		t.Errorf("ModifyPlan altered the plan: cluster_id=%v name=%v", planned.ClusterId, planned.Name)
	}
}

func TestUnitReplacesSwitch(t *testing.T) {
	t.Parallel()

	base := planModel(t, 1)

	same := base
	if replacesSwitch(base, same) {
		t.Error("equal models must not replace")
	}

	cluster := base
	cluster.ClusterId = types.Int64Value(2)
	if !replacesSwitch(base, cluster) {
		t.Error("a changed cluster_id replaces")
	}

	baseType := base
	baseType.BaseType = tfstr("iscsi")
	if !replacesSwitch(base, baseType) {
		t.Error("a changed base_type replaces")
	}

	code := base
	code.Code = tfstr("c1")
	if !replacesSwitch(base, code) {
		t.Error("a code set where it was null replaces")
	}

	unknown := base
	unknown.Code = types.StringUnknown()
	if !replacesSwitch(base, unknown) {
		t.Error("an unknown planned code may differ and counts as a replacement")
	}

	unknownCluster := base
	unknownCluster.ClusterId = types.Int64Unknown()
	if !replacesSwitch(base, unknownCluster) {
		t.Error("an unknown planned cluster_id may differ and counts as a replacement")
	}

	unknownBase := base
	unknownBase.BaseType = types.StringUnknown()
	if !replacesSwitch(base, unknownBase) {
		t.Error("an unknown planned base_type may differ and counts as a replacement")
	}
}

func TestUnitRefusalAndMutationErr(t *testing.T) {
	t.Parallel()

	yes, no := true, false

	if err := refusal(&yes, nil); err != nil {
		t.Errorf("success:true is not a refusal: %v", err)
	}
	if err := refusal(nil, nil); err != nil {
		t.Errorf("no success field is not a refusal: %v", err)
	}

	err := refusal(&no, map[string]interface{}{"msg": "Validation failed", "errors": map[string]any{"sync": "A sync is in flight"}})

	var refused refusedError
	if !errors.As(err, &refused) {
		t.Fatalf("success:false must be a refusedError, got %T %v", err, err)
	}
	if !strings.Contains(err.Error(), "Validation failed") || !strings.Contains(err.Error(), "A sync is in flight") {
		t.Errorf("refusal message = %q, want msg and the per-field error", err)
	}

	// mutationErr: a transport/status error wins; then the body.
	if mutationErr(nil, nil, &http.Response{StatusCode: 200}) != nil {
		t.Error("a 200 with no body is not an error")
	}
	var final refusedError

	err = mutationErr(&sdk.UpdateVirtualSwitch200Response{Success: &no}, nil, &http.Response{StatusCode: 200})
	if !errors.As(err, &final) || final.syncing {
		t.Errorf("a 200 with success:false is a final refusedError, got %T %v", err, err)
	}
	if mutationErr(nil, errors.New("boom"), nil) == nil {
		t.Error("a transport error is an error")
	}
}

// retryRefused: a sync-guard refusal is retried and succeeds; any other
// error, a final (non-sync) refusal, exhaustion, a deadline too close to pause
// again and a cancelled context all return at once.
func TestUnitRetryRefused(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	far := time.Now().Add(time.Hour)
	syncing := refusedError{msg: "A sync is in flight", syncing: true}

	calls := 0
	err := retryRefused(ctx, "w", far, 5, time.Millisecond, nil, func() error {
		calls++
		if calls < 3 {
			return syncing
		}

		return nil
	})
	if err != nil || calls != 3 {
		t.Errorf("success after refusals: err=%v calls=%d, want nil after 3", err, calls)
	}

	calls = 0
	err = retryRefused(ctx, "w", far, 5, time.Millisecond, nil, func() error {
		calls++

		return errors.New("boom")
	})
	if err == nil || calls != 1 {
		t.Errorf("non-refusal: err=%v calls=%d, want the error after 1 call", err, calls)
	}

	calls = 0
	final := refusedError{msg: "VLAN in use"}
	err = retryRefused(ctx, "w", far, 5, time.Millisecond, nil, func() error {
		calls++

		return final
	})
	if !errors.Is(err, final) || calls != 1 {
		t.Errorf("final refusal: err=%v calls=%d, want it unretried", err, calls)
	}

	calls = 0
	err = retryRefused(ctx, "w", far, 3, time.Millisecond, nil, func() error {
		calls++

		return syncing
	})
	if !errors.Is(err, syncing) || calls != 3 || !strings.Contains(err.Error(), "still refused after 3 attempts") {
		t.Errorf("exhaustion: err=%v calls=%d, want the refusal after 3 calls, saying so", err, calls)
	}

	// With a settle function the pause is not used: the switch is polled
	// until it stops syncing, and a settle failure ends the retry.
	calls, settled := 0, 0
	settle := func(time.Duration) error {
		settled++

		return nil
	}
	started := time.Now()
	err = retryRefused(ctx, "w", far, 5, 20*time.Millisecond, settle, func() error {
		calls++
		if calls < 3 {
			return syncing
		}

		return nil
	})
	if err != nil || calls != 3 || settled != 2 {
		t.Errorf("settle: err=%v calls=%d settled=%d, want nil after 3 calls and 2 settles", err, calls, settled)
	}

	// A settle that returns at once is topped up to the pause, so attempts
	// are never back to back.
	if elapsed := time.Since(started); elapsed < 40*time.Millisecond {
		t.Errorf("two instant settles took %s, want at least two pauses (40ms)", elapsed)
	}

	calls = 0
	failSettle := func(time.Duration) error { return errors.New("poll failed") }
	err = retryRefused(ctx, "w", far, 5, time.Millisecond, failSettle, func() error {
		calls++

		return syncing
	})
	if !errors.Is(err, syncing) || calls != 1 || !strings.Contains(err.Error(), "poll failed") {
		t.Errorf("settle failure: err=%v calls=%d, want the refusal wrapping the poll failure after one call", err, calls)
	}

	calls = 0
	err = retryRefused(ctx, "w", time.Now().Add(5*time.Millisecond), 5, time.Second, nil, func() error {
		calls++

		return syncing
	})
	if !errors.Is(err, syncing) || calls != 1 || !strings.Contains(err.Error(), "leaves no time to retry") {
		t.Errorf("deadline inside the pause: err=%v calls=%d, want one call and the timeout note", err, calls)
	}

	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	calls = 0
	err = retryRefused(cancelled, "w", far, 5, time.Minute, nil, func() error {
		calls++

		return syncing
	})
	if !errors.Is(err, syncing) || !errors.Is(err, context.Canceled) || calls != 1 {
		t.Errorf("cancelled context: err=%v calls=%d, want the refusal wrapping the cancellation after one call", err, calls)
	}
}

// bodyErr stands in for the SDK's GenericOpenAPIError: an error carrying the
// response body.
type bodyErr string

func (e bodyErr) Error() string { return "400 Bad Request" }
func (e bodyErr) Body() []byte  { return []byte(e) }

// syncRefusal recognises only a 400 whose body carries errors.sync.
func TestUnitSyncRefusal(t *testing.T) {
	t.Parallel()

	guard := bodyErr(`{"success":false,"errors":{"sync":"A sync is in flight"}}`)
	other := bodyErr(`{"success":false,"errors":{"name":"Name is in use"}}`)

	err := syncRefusal(guard, &http.Response{StatusCode: 400})

	var refused refusedError
	if !errors.As(err, &refused) || !refused.syncing || !strings.Contains(err.Error(), "A sync is in flight") {
		t.Errorf("sync guard: got %T %v, want a syncing refusedError with its message", err, err)
	}

	if err := syncRefusal(other, &http.Response{StatusCode: 400}); err != nil {
		t.Errorf("a 400 without errors.sync is not the guard: %v", err)
	}

	if err := syncRefusal(guard, &http.Response{StatusCode: 500}); err != nil {
		t.Errorf("a 500 is not the guard: %v", err)
	}

	if err := syncRefusal(nil, &http.Response{StatusCode: 200}); err != nil {
		t.Errorf("success is not the guard: %v", err)
	}

	// mutationErr classifies the guard the same way.
	if err := mutationErr(nil, guard, &http.Response{StatusCode: 400}); !errors.As(err, &refused) || !refused.syncing {
		t.Errorf("mutationErr on the guard: got %T %v", err, err)
	}

	// ... and a 200 with success:false as a final refusedError.
	no := false
	err = mutationErr(&sdk.UpdateVirtualSwitch200Response{Success: &no}, nil, &http.Response{StatusCode: 200})
	if !errors.As(err, &refused) || refused.syncing {
		t.Errorf("200 success:false: got %T %v, want a final refusedError", err, err)
	}
}

// The sync guard as the generated SDK actually returns it: a
// *sdk.GenericOpenAPIError carrying the 400 body, from a real round trip.
func TestUnitSyncRefusalFromSDKError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"success":false,"errors":{"sync":"A sync is in flight"}}`))
	}))
	defer srv.Close()

	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: srv.URL}}
	client := sdk.NewAPIClient(cfg)

	resp, httpResp, err := client.VirtualSwitchesAPI.RenameVirtualSwitch(context.Background(), 1).
		RenameVirtualSwitchRequest(buildRename("x")).Execute()

	var refused refusedError
	if e := mutationErr(resp, err, httpResp); !errors.As(e, &refused) || !refused.syncing ||
		!strings.Contains(e.Error(), "A sync is in flight") {
		t.Fatalf("SDK 400 sync guard: got %T %v, want a syncing refusedError with its message", e, e)
	}
}
