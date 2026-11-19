// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitParseJuniperSdnConfigKeepsPriorAddress pins that a config read
// without controllerAddress keeps the prior value rather than nulling an
// attribute that requires replacement, and that a returned value wins.
func TestUnitParseJuniperSdnConfigKeepsPriorAddress(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	got, err := parseJuniperSdnConfig(ctx, map[string]interface{}{"showIntrospect": "off"}, types.StringValue("10.0.0.10"))
	if err != nil {
		t.Fatal(err)
	}

	if got.ControllerAddress.ValueString() != "10.0.0.10" {
		t.Errorf("absent key: controller_address = %v, want the prior 10.0.0.10", got.ControllerAddress)
	}

	got, err = parseJuniperSdnConfig(ctx, map[string]interface{}{"controllerAddress": "10.0.0.20"}, types.StringValue("10.0.0.10"))
	if err != nil {
		t.Fatal(err)
	}

	if got.ControllerAddress.ValueString() != "10.0.0.20" {
		t.Errorf("returned key: controller_address = %v, want 10.0.0.20", got.ControllerAddress)
	}

	got, err = parseJuniperSdnConfig(ctx, map[string]interface{}{}, types.StringNull())
	if err != nil {
		t.Fatal(err)
	}

	if !got.ControllerAddress.IsNull() {
		t.Errorf("no key, no prior: controller_address = %v, want null", got.ControllerAddress)
	}
}
