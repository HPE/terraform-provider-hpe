// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

// TestUnitTaintResourceState_WritesIDToRealSchema exercises the taint path
// against the real resource schema. This is the gap that let the missing-id
// bug through: a test that only asserts an error is returned would not catch
// a SetAttribute failure on a schema that lacks the "id" attribute.
func TestUnitTaintResourceState_WritesIDToRealSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	r := &Resource{}

	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	// Initialize the state with a null object of the correct type, matching
	// what the framework provides in a real CreateResponse.
	schemaType := schemaResp.Schema.Type().TerraformType(ctx)
	nullVal := tftypes.NewValue(schemaType, nil)

	state := tfsdk.State{
		Schema: schemaResp.Schema,
		Raw:    nullVal,
	}

	var diags diag.Diagnostics

	cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
		ResourceType: "instance_node",
		ResourceID:   42,
		StateWriter:  &state,
		Diagnostics:  &diags,
	})

	// The taint call must not produce errors — only a warning is expected.
	if diags.HasError() {
		for _, d := range diags.Errors() {
			t.Errorf("TaintResourceState produced an error against the real schema: %s: %s",
				d.Summary(), d.Detail())
		}

		t.FailNow()
	}

	// Verify the id landed in state.
	var id types.Int64

	getDiags := state.GetAttribute(ctx, path.Root("id"), &id)
	if getDiags.HasError() {
		t.Fatalf("failed to read id from state after taint: %v", getDiags)
	}

	if id.IsNull() || id.IsUnknown() {
		t.Fatal("expected id to be set after taint, got null/unknown")
	}

	if id.ValueInt64() != 42 {
		t.Errorf("expected id=42 in state, got %d", id.ValueInt64())
	}
}
