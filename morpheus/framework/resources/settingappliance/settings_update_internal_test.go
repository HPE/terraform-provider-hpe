// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// updateSettings must have the update confirmed by the response body: only
// success true confirms it. Morpheus renders a refusal as HTTP 200 with
// success false; an empty body (the generated client leaves the result nil)
// and a body without the success flag confirm nothing. All are errors, so
// Create and Update record no credential digests for an update the appliance
// did not confirm.
func TestUnitUpdateSettingsConfirmsTheUpdate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		status  int
		body    string
		wantErr string // empty: no error
	}{
		{"applied", http.StatusOK, `{"success":true}`, ""},
		{
			"refused", http.StatusOK,
			`{"success":false,"errors":{"defaultRoleId":"Default Tenant Role not found for '99'"}}`,
			"refused the settings: defaultRoleId: Default Tenant Role not found for '99'",
		},
		{"empty body", http.StatusOK, "", "empty body"},
		{"no success flag", http.StatusOK, `{}`, "no success flag"},
		{"server error", http.StatusInternalServerError, `{"msg":"boom"}`, "500"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut || r.URL.Path != "/api/appliance-settings" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
				}

				if tc.body != "" {
					w.Header().Set("Content-Type", "application/json")
				}

				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}))
			t.Cleanup(srv.Close)

			cfg := sdk.NewConfiguration()
			cfg.Servers = sdk.ServerConfigurations{{URL: srv.URL}}

			host := "10.0.0.100"

			var body settingsUpdate
			body.ProxyHost.Set(&host)

			err := updateSettings(context.Background(), sdk.NewAPIClient(cfg), body)

			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tc.wantErr != "" && err == nil:
				t.Fatalf("want an error containing %q, got none", tc.wantErr)
			case tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr):
				t.Fatalf("error %q does not contain %q", err, tc.wantErr)
			}
		})
	}
}
