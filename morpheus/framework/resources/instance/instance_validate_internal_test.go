// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestHvmHostAffinityConflictValidator_BothSet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, resp := hvmValidateReq(ctx, t, ptr(int64(42)), ptr(int64(7)))

	v := hvmHostAffinityConflictValidator{}
	v.ValidateResource(ctx, req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when both kvm_host_id and affinity_group_id are set")
	}

	found := false
	for _, d := range resp.Diagnostics.Errors() {
		if d.Summary() == "Conflicting host placement attributes" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 'Conflicting host placement attributes' error, got: %v", resp.Diagnostics.Errors())
	}
}

func TestHvmHostAffinityConflictValidator_OnlyHost(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, resp := hvmValidateReq(ctx, t, ptr(int64(42)), nil)

	v := hvmHostAffinityConflictValidator{}
	v.ValidateResource(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics.Errors())
	}
}

func TestHvmHostAffinityConflictValidator_OnlyAffinity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, resp := hvmValidateReq(ctx, t, nil, ptr(int64(7)))

	v := hvmHostAffinityConflictValidator{}
	v.ValidateResource(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics.Errors())
	}
}

func TestHvmHostAffinityConflictValidator_ConfigHvmNull(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := getTestSchema(ctx)
	rootType := schemaRootTfType(ctx, s)

	rawAttrs := map[string]tftypes.Value{}
	for name, attr := range s.Attributes {
		rawAttrs[name] = tftypes.NewValue(attr.GetType().TerraformType(ctx), nil)
	}

	req := fwresource.ValidateConfigRequest{
		Config: tfsdk.Config{
			Schema: s,
			Raw:    tftypes.NewValue(rootType, rawAttrs),
		},
	}
	resp := &fwresource.ValidateConfigResponse{}

	v := hvmHostAffinityConflictValidator{}
	v.ValidateResource(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error when config_hvm is null: %v", resp.Diagnostics.Errors())
	}
}

// --- helpers ---

func ptr[T any](v T) *T { return &v }

func getTestSchema(ctx context.Context) schema.Schema {
	r := &Resource{}
	schemaResp := fwresource.SchemaResponse{}
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)

	return schemaResp.Schema
}

func hvmValidateReq(
	ctx context.Context,
	t *testing.T,
	kvmHostId *int64,
	affinityGroupId *int64,
) (fwresource.ValidateConfigRequest, *fwresource.ValidateConfigResponse) {
	t.Helper()

	s := getTestSchema(ctx)

	hvmType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"affinity_group_id":     tftypes.Number,
		"create_user":           tftypes.Bool,
		"kvm_host_id":           tftypes.Number,
		"nested_virtualization": tftypes.String,
		"no_agent":              tftypes.Bool,
		"resource_pool_id":      tftypes.String,
	}}

	hvmAttrs := map[string]tftypes.Value{
		"create_user":           tftypes.NewValue(tftypes.Bool, nil),
		"nested_virtualization": tftypes.NewValue(tftypes.String, nil),
		"no_agent":              tftypes.NewValue(tftypes.Bool, nil),
		"resource_pool_id":      tftypes.NewValue(tftypes.String, nil),
	}

	if kvmHostId != nil {
		hvmAttrs["kvm_host_id"] = tftypes.NewValue(tftypes.Number, *kvmHostId)
	} else {
		hvmAttrs["kvm_host_id"] = tftypes.NewValue(tftypes.Number, nil)
	}

	if affinityGroupId != nil {
		hvmAttrs["affinity_group_id"] = tftypes.NewValue(tftypes.Number, *affinityGroupId)
	} else {
		hvmAttrs["affinity_group_id"] = tftypes.NewValue(tftypes.Number, nil)
	}

	hvmVal := tftypes.NewValue(hvmType, hvmAttrs)

	rootType := schemaRootTfType(ctx, s)
	rawAttrs := map[string]tftypes.Value{}
	for name, attr := range s.Attributes {
		if name == "config_hvm" {
			rawAttrs[name] = hvmVal
		} else {
			rawAttrs[name] = tftypes.NewValue(attr.GetType().TerraformType(ctx), nil)
		}
	}

	req := fwresource.ValidateConfigRequest{
		Config: tfsdk.Config{
			Schema: s,
			Raw:    tftypes.NewValue(rootType, rawAttrs),
		},
	}

	return req, &fwresource.ValidateConfigResponse{}
}

func schemaRootTfType(ctx context.Context, s schema.Schema) tftypes.Object {
	rootType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}
	for name, attr := range s.Attributes {
		rootType.AttributeTypes[name] = attr.GetType().TerraformType(ctx)
	}

	return rootType
}
