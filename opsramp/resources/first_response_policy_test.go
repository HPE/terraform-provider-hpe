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

func makeFirstResponsePolicyPlan(t *testing.T, model FirstResponsePolicyModel) tfsdk.Plan {
	t.Helper()

	ctx := context.Background()
	r := &FirstResponsePolicyResource{}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := plan.Set(ctx, &model); diags.HasError() {
		t.Fatalf("build plan: %v", diags)
	}

	return plan
}

func newFirstResponsePolicyModel(attributeActions, patternActions types.Object) FirstResponsePolicyModel {
	return FirstResponsePolicyModel{
		Client:           types.StringNull(),
		Id:               types.StringNull(),
		Name:             types.StringValue("policy"),
		EnabledMode:      types.StringValue("OBSERVED"),
		FilterQuery:      types.StringValue(""),
		AttributeActions: attributeActions,
		PatternActions:   patternActions,
	}
}

func TestUnitBuildFirstResponsePolicyRequestOmitsUnsetComputedBlocks(t *testing.T) {
	t.Parallel()

	policy := buildFirstResponsePolicyRequest(FirstResponsePolicyModel{
		Name:             types.StringValue("Permanent Suppression"),
		EnabledMode:      types.StringValue("OBSERVED"),
		FilterQuery:      types.StringValue("objectName = \"LAB-WINDOWS-SERVER01\""),
		AttributeActions: types.ObjectNull(attributeActionsAttrTypes),
		PatternActions:   types.ObjectNull(patternActionsAttrTypes),
	})

	if policy.AttributeActions != nil {
		t.Fatalf("expected attribute_actions to be omitted, got %#v", policy.AttributeActions)
	}

	if policy.PatternActions != nil {
		t.Fatalf("expected pattern_actions to be omitted, got %#v", policy.PatternActions)
	}
}

func TestUnitMapFirstResponsePolicyToStateLeavesOmittedNestedBlocksNull(t *testing.T) {
	t.Parallel()

	state := &FirstResponsePolicyModel{}
	mapFirstResponsePolicyToState(&client.FirstResponsePolicy{
		Id:          "POLICY-AC-123",
		Name:        "Permanent Suppression",
		EnabledMode: "OBSERVED",
		FilterQuery: "objectName = \"LAB-WINDOWS-SERVER01\"",
		AttributeActions: &client.FirstResponseAttrActions{
			ContinuousLearning: false,
			Suppress: &client.FirstResponseAttrSuppress{
				LearnedConfiguration: false,
				SuppressDuration:     -1,
			},
			Insights:   nil,
			RunProcess: nil,
		},
		PatternActions: &client.FirstResponsePatternActions{
			SeasonalityTimeFrame: "",
			Suppress:             nil,
		},
	}, state)

	if state.AttributeActions.IsNull() {
		t.Fatal("expected attribute_actions object to be present")
	}

	attributeActions := state.AttributeActions.Attributes()
	if !attributeActions["insights"].IsNull() {
		t.Fatalf("expected insights to remain null when omitted by API, got %#v", attributeActions["insights"])
	}
	if !attributeActions["run_process"].IsNull() {
		t.Fatalf("expected run_process to remain null when omitted by API, got %#v", attributeActions["run_process"])
	}

	if state.PatternActions.IsNull() {
		t.Fatal("expected pattern_actions object to be present")
	}

	patternActions := state.PatternActions.Attributes()
	if !patternActions["seasonality_time_frame"].IsNull() {
		t.Fatalf("expected seasonality_time_frame to remain null when omitted by API, got %#v", patternActions["seasonality_time_frame"])
	}
	if !patternActions["suppress"].IsNull() {
		t.Fatalf("expected pattern suppress block to remain null when omitted by API, got %#v", patternActions["suppress"])
	}
}

func TestUnitFirstResponsePolicyModifyPlanRequiresAction(t *testing.T) {
	t.Parallel()

	makeAttributeActions := func(suppressDuration *int64, suppressLearning bool, runProcessIDs []string, runProcessLearning bool, createPrcInsights bool) types.Object {
		suppressValue := types.ObjectNull(attrSuppressAttrTypes)
		if suppressDuration != nil || suppressLearning {
			learnedConfiguration := types.BoolNull()
			if suppressLearning || suppressDuration != nil {
				learnedConfiguration = types.BoolValue(suppressLearning)
			}

			durationValue := types.Int64Null()
			if suppressDuration != nil {
				durationValue = types.Int64Value(*suppressDuration)
			}

			suppressValue = types.ObjectValueMust(attrSuppressAttrTypes, map[string]attr.Value{
				"learned_configuration": learnedConfiguration,
				"suppress_duration":     durationValue,
			})
		}

		runProcessValue := types.ObjectNull(attrRunProcessAttrTypes)
		if runProcessIDs != nil || runProcessLearning {
			processIDValues := types.ListNull(types.StringType)
			if runProcessIDs != nil {
				elements := make([]attr.Value, len(runProcessIDs))
				for i, processID := range runProcessIDs {
					elements[i] = types.StringValue(processID)
				}
				processIDValues = types.ListValueMust(types.StringType, elements)
			}

			runProcessValue = types.ObjectValueMust(attrRunProcessAttrTypes, map[string]attr.Value{
				"learned_configuration": types.BoolValue(runProcessLearning),
				"run_immediately":       types.BoolValue(false),
				"process_ids":           processIDValues,
			})
		}

		insightsValue := types.ObjectNull(insightsAttrTypes)
		if createPrcInsights {
			insightsValue = types.ObjectValueMust(insightsAttrTypes, map[string]attr.Value{
				"create_prc_insights": types.BoolValue(true),
			})
		}

		return types.ObjectValueMust(attributeActionsAttrTypes, map[string]attr.Value{
			"continuous_learning": types.BoolValue(false),
			"suppress":            suppressValue,
			"insights":            insightsValue,
			"run_process":         runProcessValue,
		})
	}

	makePatternActions := func(seasonalAlerts bool) types.Object {
		return types.ObjectValueMust(patternActionsAttrTypes, map[string]attr.Value{
			"seasonality_time_frame": types.StringValue("7D"),
			"suppress": types.ObjectValueMust(patternSuppressAttrTypes, map[string]attr.Value{
				"seasonal_alerts": types.BoolValue(seasonalAlerts),
			}),
		})
	}

	t.Run("rejects policies without an enabled action", func(t *testing.T) {
		plan := makeFirstResponsePolicyPlan(t, newFirstResponsePolicyModel(
			makeAttributeActions(nil, false, nil, false, false),
			makePatternActions(false),
		))

		r := &FirstResponsePolicyResource{}
		resp := &resource.ModifyPlanResponse{Plan: plan}
		r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan}, resp)

		if !resp.Diagnostics.HasError() {
			t.Fatal("expected plan diagnostics when no supported action is enabled")
		}
	})

	t.Run("accepts attribute suppress action", func(t *testing.T) {
		suppressDuration := int64(-1)
		plan := makeFirstResponsePolicyPlan(t, newFirstResponsePolicyModel(
			makeAttributeActions(&suppressDuration, false, nil, false, false),
			types.ObjectNull(patternActionsAttrTypes),
		))

		r := &FirstResponsePolicyResource{}
		resp := &resource.ModifyPlanResponse{Plan: plan}
		r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for attribute suppress action, got: %v", resp.Diagnostics)
		}
	})

	t.Run("accepts suppress learned configuration action", func(t *testing.T) {
		plan := makeFirstResponsePolicyPlan(t, newFirstResponsePolicyModel(
			makeAttributeActions(nil, true, nil, false, false),
			types.ObjectNull(patternActionsAttrTypes),
		))

		r := &FirstResponsePolicyResource{}
		resp := &resource.ModifyPlanResponse{Plan: plan}
		r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for suppress learned_configuration action, got: %v", resp.Diagnostics)
		}
	})

	t.Run("accepts run process action", func(t *testing.T) {
		plan := makeFirstResponsePolicyPlan(t, newFirstResponsePolicyModel(
			makeAttributeActions(nil, false, []string{"process-1"}, false, false),
			types.ObjectNull(patternActionsAttrTypes),
		))

		r := &FirstResponsePolicyResource{}
		resp := &resource.ModifyPlanResponse{Plan: plan}
		r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for run_process action, got: %v", resp.Diagnostics)
		}
	})

	t.Run("accepts run process learned configuration action", func(t *testing.T) {
		plan := makeFirstResponsePolicyPlan(t, newFirstResponsePolicyModel(
			makeAttributeActions(nil, false, nil, true, false),
			types.ObjectNull(patternActionsAttrTypes),
		))

		r := &FirstResponsePolicyResource{}
		resp := &resource.ModifyPlanResponse{Plan: plan}
		r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for run_process learned_configuration action, got: %v", resp.Diagnostics)
		}
	})

	t.Run("accepts prc insights action", func(t *testing.T) {
		plan := makeFirstResponsePolicyPlan(t, newFirstResponsePolicyModel(
			makeAttributeActions(nil, false, nil, false, true),
			types.ObjectNull(patternActionsAttrTypes),
		))

		r := &FirstResponsePolicyResource{}
		resp := &resource.ModifyPlanResponse{Plan: plan}
		r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for PRC insights action, got: %v", resp.Diagnostics)
		}
	})

	t.Run("accepts seasonal alerts action", func(t *testing.T) {
		plan := makeFirstResponsePolicyPlan(t, newFirstResponsePolicyModel(
			types.ObjectNull(attributeActionsAttrTypes),
			makePatternActions(true),
		))

		r := &FirstResponsePolicyResource{}
		resp := &resource.ModifyPlanResponse{Plan: plan}
		r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no diagnostics for seasonal alerts action, got: %v", resp.Diagnostics)
		}
	})
}
