// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
)

// TestUnitValidateOrganizationMatchingDefaultsNullToAllAtMSPScope covers the backward-compatible
// behavior restored for MSP-scoped alert policies: when organization_matching_type is omitted in
// config, the attribute is Computed (so Terraform core would plan it as unknown rather than
// null), and ModifyPlan must not reject it as a missing required attribute. The actual "ALL"
// default is then resolved later from the API response by the resource's map*ToState function.
func TestUnitValidateOrganizationMatchingDefaultsNullToAllAtMSPScope(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	r := &AlertCorrelationPolicyResource{}
	r.apiClient = &client.OpsRampClient{Scope: "MSP"}

	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	configModel := AlertCorrelationPolicyModel{
		Client:                   types.StringNull(),
		Id:                       types.StringNull(),
		Name:                     types.StringValue("test"),
		EnabledMode:              types.StringNull(),
		OrganizationMatchingType: types.StringNull(),
		IncludedClients:          types.ListNull(types.StringType),
		Precedence:               types.Int64Null(),
		Type:                     types.StringNull(),
		FilterQuery:              types.StringNull(),
		InferenceQuery:           types.StringNull(),
		Review:                   types.BoolNull(),
		InferenceSubject:         types.StringNull(),
	}
	configAsPlan := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := configAsPlan.Set(ctx, &configModel); diags.HasError() {
		t.Fatal(diags)
	}
	config := tfsdk.Config{Schema: schemaResp.Schema, Raw: configAsPlan.Raw}

	// Simulate how Terraform core plans a Computed attribute left null in config with no prior
	// state: both organization_matching_type and included_clients become unknown.
	planModel := configModel
	planModel.Id = types.StringUnknown()
	planModel.OrganizationMatchingType = types.StringUnknown()
	planModel.IncludedClients = types.ListUnknown(types.StringType)
	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := plan.Set(ctx, &planModel); diags.HasError() {
		t.Fatal(diags)
	}

	resp := &resource.ModifyPlanResponse{Plan: plan}
	r.ModifyPlan(ctx, resource.ModifyPlanRequest{Config: config, Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics when organization_matching_type is omitted at MSP scope, got: %v", resp.Diagnostics)
	}
}

// TestUnitValidateOrganizationMatchingStillRequiresIncludedClientsForInclude ensures the
// Computed/UseStateForUnknown restoration above didn't loosen the "INCLUDE requires a non-empty
// included_clients" rule for a known organization_matching_type.
func TestUnitValidateOrganizationMatchingStillRequiresIncludedClientsForInclude(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	r := &AlertCorrelationPolicyResource{}
	r.apiClient = &client.OpsRampClient{Scope: "MSP"}

	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	model := AlertCorrelationPolicyModel{
		Client:                   types.StringNull(),
		Id:                       types.StringNull(),
		Name:                     types.StringValue("test"),
		EnabledMode:              types.StringNull(),
		OrganizationMatchingType: types.StringValue("INCLUDE"),
		IncludedClients:          types.ListNull(types.StringType),
		Precedence:               types.Int64Null(),
		Type:                     types.StringNull(),
		FilterQuery:              types.StringNull(),
		InferenceQuery:           types.StringNull(),
		Review:                   types.BoolNull(),
		InferenceSubject:         types.StringNull(),
	}
	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := plan.Set(ctx, &model); diags.HasError() {
		t.Fatal(diags)
	}
	config := tfsdk.Config{Schema: schemaResp.Schema, Raw: plan.Raw}

	resp := &resource.ModifyPlanResponse{Plan: plan}
	r.ModifyPlan(ctx, resource.ModifyPlanRequest{Config: config, Plan: plan}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error when organization_matching_type is INCLUDE with no included_clients")
	}
}

// TestUnitApplyIncludedClientsTargetingResolvesKnownValue reproduces the apply-time bug where
// organization_matching_type (now Computed) was left unknown after Create/Update/Read for
// CLIENT-scoped policies, and for MSP-scoped ALL policies whose API response omits
// "clientsIncluded" entirely: Terraform rejects a still-unknown Computed attribute after apply
// ("Provider produced invalid result object after apply"), so this must always resolve to a
// known value.
func TestUnitApplyIncludedClientsTargetingResolvesKnownValue(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name            string
		includedClients []string
		clientsIncluded string
		isMSP           bool
		fromRead        bool
		priorMatching   types.String
		priorIncluded   types.List
		wantMatching    types.String
		wantIncluded    types.List
	}{
		{
			name:          "client scope always resolves to null regardless of prior state",
			isMSP:         false,
			fromRead:      true,
			priorMatching: types.StringUnknown(),
			priorIncluded: types.ListUnknown(types.StringType),
			wantMatching:  types.StringNull(),
			wantIncluded:  types.ListNull(types.StringType),
		},
		{
			name:            "MSP scope with no included clients and no clientsIncluded defaults to ALL",
			isMSP:           true,
			fromRead:        true,
			includedClients: nil,
			clientsIncluded: "",
			priorMatching:   types.StringUnknown(),
			priorIncluded:   types.ListUnknown(types.StringType),
			wantMatching:    types.StringValue("ALL"),
			wantIncluded:    types.ListNull(types.StringType),
		},
		{
			name:            "MSP scope with included clients resolves to INCLUDE",
			isMSP:           true,
			fromRead:        true,
			includedClients: []string{"client-a"},
			priorMatching:   types.StringUnknown(),
			priorIncluded:   types.ListUnknown(types.StringType),
			wantMatching:    types.StringValue("INCLUDE"),
			wantIncluded:    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("client-a")}),
		},
		{
			name:            "write path (not fromRead) preserves a known plan value untouched",
			isMSP:           true,
			fromRead:        false,
			includedClients: nil,
			clientsIncluded: "",
			priorMatching:   types.StringValue("INCLUDE"),
			priorIncluded:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("client-a")}),
			wantMatching:    types.StringValue("INCLUDE"),
			wantIncluded:    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("client-a")}),
		},
		{
			name:            "read path refreshes a stale known value when the remote scope changed",
			isMSP:           true,
			fromRead:        true,
			includedClients: nil,
			clientsIncluded: "",
			priorMatching:   types.StringValue("INCLUDE"),
			priorIncluded:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("client-a")}),
			wantMatching:    types.StringValue("ALL"),
			wantIncluded:    types.ListNull(types.StringType),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			matching := tc.priorMatching
			included := tc.priorIncluded
			applyIncludedClientsTargeting(tc.includedClients, tc.clientsIncluded, tc.isMSP, tc.fromRead, &matching, &included)

			if matching.IsUnknown() {
				t.Fatal("organization_matching_type must never remain unknown after resolution")
			}
			if !matching.Equal(tc.wantMatching) {
				t.Fatalf("organization_matching_type = %#v; want %#v", matching, tc.wantMatching)
			}
			if !included.Equal(tc.wantIncluded) {
				t.Fatalf("included_clients = %#v; want %#v", included, tc.wantIncluded)
			}
		})
	}
}
