// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func makeMetricAlertDefinitionPlan(t *testing.T, model MetricAlertDefinitionModel) tfsdk.Plan {
	t.Helper()

	ctx := context.Background()
	r := &MetricAlertDefinitionResource{}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := plan.Set(ctx, &model); diags.HasError() {
		t.Fatalf("build plan: %v", diags)
	}

	return plan
}

func newMetricAlertDefinitionModel(entityType string, attributes []NameValuePairModel) MetricAlertDefinitionModel {
	return MetricAlertDefinitionModel{
		Client:             types.StringNull(),
		Name:               types.StringValue("example-alert"),
		Query:              types.StringValue("metrics_samples_count"),
		AlertType:          types.StringValue("METRICS"),
		AlertThresholdType: types.StringValue("DYNAMIC_THRESHOLD"),
		AlertThresholdData: AlertThresholdDataModel{
			WarningCondition:  types.StringNull(),
			CriticalCondition: types.StringNull(),
			Limit:             types.Int64Value(2),
			Direction:         types.StringNull(),
			LearningPeriod:    types.StringNull(),
			StandardDeviation: types.Int64Null(),
			LowerLimit:        types.StringNull(),
			UpperLimit:        types.StringNull(),
		},
		NoDataCondition:      types.StringValue("WARNING_ALERT"),
		AlertTriggerDuration: types.StringValue("1m"),
		Subject:              types.StringValue("subject"),
		Description:          types.StringValue("description"),
		EntityType:           types.StringValue(entityType),
		Component:            types.StringValue("$$__name__"),
		Status:               types.BoolValue(true),
		IsObsolete:           types.BoolNull(),
		Labels:               nil,
		Attributes:           attributes,
	}
}

func TestUnitMetricAlertDefinitionBuildRequestWrapsStringInputs(t *testing.T) {
	t.Parallel()

	r := &MetricAlertDefinitionResource{}
	plan := newMetricAlertDefinitionModel("RESOURCE", []NameValuePairModel{{
		Name:  types.StringValue("host"),
		Value: types.StringValue("$$__name__"),
	}})

	request, err := r.buildRequest(context.Background(), &plan)
	if err != nil {
		t.Fatalf("buildRequest returned error: %v", err)
	}

	if len(request.EntityType) != 1 || request.EntityType[0] != "RESOURCE" {
		t.Fatalf("unexpected entity_type payload: %#v", request.EntityType)
	}

	if len(request.Component) != 1 || request.Component[0] != "$$__name__" {
		t.Fatalf("unexpected component payload: %#v", request.Component)
	}
}

func TestUnitMetricAlertDefinitionModifyPlanAttributesRequirement(t *testing.T) {
	t.Parallel()

	t.Run("resource requires attributes", func(t *testing.T) {
		plan := makeMetricAlertDefinitionPlan(t, newMetricAlertDefinitionModel("RESOURCE", nil))
		r := &MetricAlertDefinitionResource{}
		resp := &resource.ModifyPlanResponse{Plan: plan}

		r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan}, resp)

		if !resp.Diagnostics.HasError() {
			t.Fatal("expected plan diagnostics when RESOURCE has no attributes")
		}
	})

	t.Run("client allows empty attributes", func(t *testing.T) {
		plan := makeMetricAlertDefinitionPlan(t, newMetricAlertDefinitionModel("CLIENT", nil))
		r := &MetricAlertDefinitionResource{}
		resp := &resource.ModifyPlanResponse{Plan: plan}

		r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan}, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("expected no plan diagnostics for CLIENT without attributes, got: %v", resp.Diagnostics)
		}
	})
}

func TestUnitMetricAlertDefinitionModifyPlanAllowsEmptyNoDataConditionForStaticThreshold(t *testing.T) {
	t.Parallel()

	model := newMetricAlertDefinitionModel("RESOURCE", []NameValuePairModel{{
		Name:  types.StringValue("host"),
		Value: types.StringValue("$$__name__"),
	}})
	model.AlertThresholdType = types.StringValue("STATIC_THRESHOLD")
	model.NoDataCondition = types.StringValue("")
	model.AlertThresholdData.WarningCondition = types.StringValue(">0")

	plan := makeMetricAlertDefinitionPlan(t, model)
	r := &MetricAlertDefinitionResource{}
	resp := &resource.ModifyPlanResponse{Plan: plan}

	r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no plan diagnostics for STATIC_THRESHOLD with empty no_data_condition, got: %v", resp.Diagnostics)
	}
}
