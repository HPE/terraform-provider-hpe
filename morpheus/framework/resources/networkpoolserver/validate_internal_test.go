// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func validateWith(t *testing.T, attrs map[string]string) *resource.ValidateConfigResponse {
	t.Helper()

	ctx := context.Background()
	r := &networkPoolServerResource{}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	objType := schema.Type().TerraformType(ctx).(tftypes.Object)
	values := map[string]tftypes.Value{}
	for name, at := range objType.AttributeTypes {
		if v, ok := attrs[name]; ok {
			values[name] = tftypes.NewValue(tftypes.String, v)
		} else {
			values[name] = tftypes.NewValue(at, nil)
		}
	}

	resp := &resource.ValidateConfigResponse{}
	r.ValidateConfig(ctx, resource.ValidateConfigRequest{
		Config: tfsdk.Config{Schema: schema, Raw: tftypes.NewValue(objType, values)},
	}, resp)

	return resp
}

func TestUnitValidateConfigAppID(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		attrs        map[string]string
		wantErrors   int
		wantWarnings int
	}{
		"phpipam with app_id":      {map[string]string{"type_code": "phpipam", "app_id": "m"}, 0, 0},
		"phpipam without app_id":   {map[string]string{"type_code": "phpipam"}, 1, 0},
		"infoblox with app_id":     {map[string]string{"type_code": "infoblox", "app_id": "m"}, 1, 0},
		"type unknown is skipped":  {map[string]string{"app_id": "m"}, 0, 0},
		"bluecat with zone_filter": {map[string]string{"type_code": "bluecat", "zone_filter": "z"}, 0, 1},
		"solarwinds network_filter": {
			map[string]string{"type_code": "solarwindsipam", "network_filter": "n", "service_mode": "static"}, 0, 2,
		},
		"phpipam network_filter": {
			map[string]string{"type_code": "phpipam", "app_id": "m", "network_filter": "n"}, 0, 0,
		},
		"fallback type not warned": {map[string]string{"type_code": "solidserver", "zone_filter": "z"}, 0, 0},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := validateWith(t, tc.attrs)
			if got := resp.Diagnostics.ErrorsCount(); got != tc.wantErrors {
				t.Errorf("errors = %d, want %d: %v", got, tc.wantErrors, resp.Diagnostics)
			}
			if got := resp.Diagnostics.WarningsCount(); got != tc.wantWarnings {
				t.Errorf("warnings = %d, want %d: %v", got, tc.wantWarnings, resp.Diagnostics)
			}
		})
	}
}
