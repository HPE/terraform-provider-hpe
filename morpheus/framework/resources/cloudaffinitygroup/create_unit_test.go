// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroup

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

// testClient points the SDK at a stub appliance.
func testClient(t *testing.T, handler http.HandlerFunc) *sdk.APIClient {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	return clientfactory.NewAPIClient(
		context.Background(), srv.URL, "", "", "", "test-token",
	)
}

// listJSON renders a listing response.
//
// Built from the generated types and marshalled, rather than assembled as a
// string, so the fixture cannot drift from the shape the SDK actually parses.
func listJSON(t *testing.T, ids ...int64) []byte {
	t.Helper()

	groups := make(
		[]sdk.ListCloudAffinityGroups200ResponseAllOfAffinityGroupsInner, 0, len(ids),
	)

	for _, id := range ids {
		groups = append(
			groups,
			sdk.ListCloudAffinityGroups200ResponseAllOfAffinityGroupsInner{
				Id:   ptr(id),
				Name: ptrStr("g" + strconv.FormatInt(id, 10)),
			},
		)
	}

	body, err := json.Marshal(
		sdk.ListCloudAffinityGroups200Response{AffinityGroups: groups},
	)
	if err != nil {
		t.Fatalf("marshalling the listing fixture: %v", err)
	}

	return body
}

func ptrStr(s string) *string { return &s }

// TestCreatedGroupIDPrefersTheResponse: the ordinary path must not list
// anything. A create that rendered correctly already carries the id.
func TestCreatedGroupIDPrefersTheResponse(t *testing.T) {
	t.Parallel()

	client := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		t.Error("the listing must not be called when the response carries an id")
		w.WriteHeader(http.StatusInternalServerError)
	})

	result := &sdk.SaveCloudAffinityGroup200Response{
		AffinityGroup: &sdk.SaveCloudAffinityGroup200ResponseAllOfAffinityGroup{Id: ptr(int64(42))},
	}

	var diags diag.Diagnostics

	got, ok := createdGroupID(context.Background(), client, 1, result, nil, false, &diags)
	if !ok {
		t.Fatalf("not ok: %v", diags.Errors())
	}

	if got != 42 {
		t.Errorf("id = %d, want 42", got)
	}
}

// TestCreatedGroupIDRecoversAfterRenderFailure is the regression guard for
// MORPH-15806.
//
// A create carrying tenant permissions answers 500 with no body, so the id has
// to be recovered from the listing. Before this, the create failed outright and
// the group it had just made was left on the appliance, in nobody's state.
func TestCreatedGroupIDRecoversAfterRenderFailure(t *testing.T) {
	t.Parallel()

	client := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(listJSON(t, 1, 2, 7))
	})

	prior := map[int64]struct{}{1: {}, 2: {}}

	var diags diag.Diagnostics

	got, ok := createdGroupID(context.Background(), client, 1, nil, prior, true, &diags)
	if !ok {
		t.Fatalf("not ok: %v", diags.Errors())
	}

	if got != 7 {
		t.Errorf("id = %d, want 7", got)
	}
}

// TestCreatedGroupIDRefusesToGuess: concurrent creation makes the new group
// ambiguous. Adopting the wrong one into state is worse than an error.
func TestCreatedGroupIDRefusesToGuess(t *testing.T) {
	t.Parallel()

	client := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(listJSON(t, 1, 8, 9))
	})

	prior := map[int64]struct{}{1: {}}

	var diags diag.Diagnostics

	if _, ok := createdGroupID(context.Background(), client, 1, nil, prior, true, &diags); ok {
		t.Fatal("ok with two candidate ids; must refuse to guess")
	}

	if !diags.HasError() {
		t.Error("no diagnostic raised")
	}
}

// TestCreatedGroupIDReportsUnlistablePrecondition: when the pre-create listing
// could not be read there is nothing to diff against, so the group cannot be
// identified. The practitioner has to be told it exists.
func TestCreatedGroupIDReportsUnlistablePrecondition(t *testing.T) {
	t.Parallel()

	client := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		t.Error("the listing must not be called when there is nothing to diff against")
		w.WriteHeader(http.StatusInternalServerError)
	})

	var diags diag.Diagnostics

	if _, ok := createdGroupID(context.Background(), client, 1, nil, nil, true, &diags); ok {
		t.Fatal("ok without a prior listing")
	}

	if !diags.HasError() {
		t.Fatal("no diagnostic raised")
	}

	if got := diags.Errors()[0].Summary(); got != "affinity group created but could not be identified" {
		t.Errorf("summary = %q", got)
	}
}

// TestCreatedGroupIDStillRejectsAGenuinelyEmptyResponse: the nil-id guard must
// survive. A 200 with no id is a real fault and must not be swept into the
// workaround.
func TestCreatedGroupIDStillRejectsAGenuinelyEmptyResponse(t *testing.T) {
	t.Parallel()

	client := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		t.Error("the listing must not be called when the create did not fail to render")
		w.WriteHeader(http.StatusInternalServerError)
	})

	var diags diag.Diagnostics

	if _, ok := createdGroupID(context.Background(), client, 1, nil, nil, false, &diags); ok {
		t.Fatal("ok with a nil id and no render failure")
	}

	if got := diags.Errors()[0].Summary(); got != "API returned nil ID" {
		t.Errorf("summary = %q", got)
	}
}
