// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package helpers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	sdklegacy "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/client"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

// newTestClient builds a legacy client (token auth, login skipped) pointed at a
// whoami test server. Each call constructs a fresh client pointer so the
// per-pointer master-tenant cache does not leak across cases.
func newTestClient(t *testing.T, url string) *sdklegacy.Client {
	t.Helper()

	return client.NewLegacyClient(
		context.Background(),
		url,
		"",
		"",
		"",
		"abc",
		sdklegacy.SkipLogin(),
	)
}

// TestUnitVisibilityCustomizeDiff_MORPH16419 verifies the SDKv2 visibility
// CustomizeDiff: a master caller (isMasterAccount present-true) allows public,
// a sub-tenant caller (isMasterAccount absent) rejects it, and a non-public /
// unknown value skips the whoami call entirely.
func TestUnitVisibilityCustomizeDiff_MORPH16419(t *testing.T) {
	cases := []struct {
		name       string
		whoamiBody string
		visibility any
		wantErr    bool
		wantCall   bool
	}{
		{
			name:       "master public allowed",
			whoamiBody: `{"isMasterAccount":true}`,
			visibility: "public",
			wantErr:    false,
			wantCall:   true,
		},
		{
			name:       "sub-tenant public rejected",
			whoamiBody: `{}`,
			visibility: "public",
			wantErr:    true,
			wantCall:   true,
		},
		{
			name:       "private skips whoami",
			whoamiBody: `{}`,
			visibility: "private",
			wantErr:    false,
			wantCall:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var called bool
			srv := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					called = true
					fmt.Fprint(w, tc.whoamiBody)
				}))
			defer srv.Close()

			cl := newTestClient(t, srv.URL)

			res := &schema.Resource{
				Schema: map[string]*schema.Schema{
					"visibility": {Type: schema.TypeString, Optional: true},
				},
				CustomizeDiff: helpers.VisibilityCustomizeDiff,
			}

			diff, err := res.Diff(
				context.Background(),
				nil,
				terraform.NewResourceConfigRaw(
					map[string]any{"visibility": tc.visibility},
				),
				cl,
			)
			_ = diff

			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if called != tc.wantCall {
				t.Errorf("whoami called = %v, want %v", called, tc.wantCall)
			}
		})
	}
}
