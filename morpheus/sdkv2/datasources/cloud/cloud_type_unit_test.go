// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sdklegacy "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/client"
)

// TestUnitCloudTypeDataSourceUsesZoneTypesEndpoint_MORPH16399 verifies the
// cloud_type data source reads the tenant-reachable /api/zone-types endpoint
// (not the appliance-scoped /api/appliance-settings/zone-types, which 403s for
// sub-tenants) and resolves the matching type id.
func TestUnitCloudTypeDataSourceUsesZoneTypesEndpoint_MORPH16399(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			fmt.Fprint(w, `{"zoneTypes":[{"id":1,"name":"Private Cloud"}]}`)
		}))
	defer srv.Close()

	cl := client.NewLegacyClient(
		context.Background(),
		srv.URL,
		"",
		"",
		"",
		"abc",
		sdklegacy.SkipLogin(),
	)

	res := DataSourceCloudType()
	d := res.TestResourceData()
	d.Set("name", "Private Cloud")

	if diags := dataSourceCloudTypeRead(context.Background(), d, cl); diags.HasError() {
		t.Fatalf("read returned diagnostics: %v", diags)
	}

	if gotPath != "/api/zone-types" {
		t.Errorf("request path = %q, want /api/zone-types", gotPath)
	}
	if d.Id() != "1" {
		t.Errorf("id = %q, want \"1\"", d.Id())
	}
}
