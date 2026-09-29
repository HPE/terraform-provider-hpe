// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestGetInstanceToleratesStringBoolViaClient is the end-to-end guard for
// MORPH-16966.
//
// A live appliance cannot reproduce the defect — no real instance emits a
// non-boolean string in a boolean config field — so this stands in for an
// acceptance test by serving a crafted response through a mock server and
// driving it with the real SDK client. It exercises the full transport and
// decode path the provider read uses, not just the response type's
// UnmarshalJSON.
//
// Before the fix, `createBackup: "banana"` failed the whole decode and the
// client returned no instance at all.
func TestGetInstanceToleratesStringBoolViaClient(t *testing.T) {
	t.Parallel()

	const body = `{
		"instance": {
			"id": 42,
			"name": "web-01",
			"config": {
				"createUser": "yes",
				"createBackup": "banana",
				"isVpcSelectable": "off"
			}
		}
	}`

	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(body))
		}))
	defer srv.Close()

	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: srv.URL}}
	client := sdk.NewAPIClient(cfg)

	resp, hresp, err := client.InstancesAPI.GetInstance(context.Background(), 42).Execute()
	if err != nil {
		t.Fatalf("GetInstance returned an error: %v", err)
	}

	if hresp == nil || hresp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected HTTP status: %v", hresp)
	}

	if resp == nil || resp.Instance == nil {
		t.Fatal("no instance returned — the stray bool value sank the whole response")
	}

	if resp.Instance.Id == nil || *resp.Instance.Id != 42 {
		t.Errorf("id = %v, want 42", resp.Instance.Id)
	}

	if resp.Instance.Name == nil || *resp.Instance.Name != "web-01" {
		t.Errorf("name = %v, want web-01", resp.Instance.Name)
	}

	cfgOut := resp.Instance.Config
	if cfgOut == nil {
		t.Fatal("config not decoded")
	}

	// The recognised value coerced to true.
	if cfgOut.CreateUser == nil || !*cfgOut.CreateUser {
		t.Error(`createUser: want true from "yes"`)
	}

	// The unrecognised value degraded to false instead of failing the response.
	if cfgOut.CreateBackup == nil || *cfgOut.CreateBackup {
		t.Error(`createBackup: want false from the unrecognised "banana"`)
	}

	// The "on"/"off" coercion still holds — assert the field the test seeds.
	if cfgOut.IsVpcSelectable == nil || *cfgOut.IsVpcSelectable {
		t.Error(`isVpcSelectable: want false from "off"`)
	}
}
