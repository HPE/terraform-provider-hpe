// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AlertCorrelationPolicyModelV0 is the state shape produced by schema version 0, every released
// version of the resource before alerts_time_window became a number. Its only difference from
// AlertCorrelationPolicyModel is AlgorithmCorrelation's type.
type AlertCorrelationPolicyModelV0 struct {
	Client                   types.String                 `tfsdk:"client"`
	Id                       types.String                 `tfsdk:"id"`
	Name                     types.String                 `tfsdk:"name"`
	EnabledMode              types.String                 `tfsdk:"enabled_mode"`
	OrganizationMatchingType types.String                 `tfsdk:"organization_matching_type"`
	IncludedClients          types.List                   `tfsdk:"included_clients"`
	Precedence               types.Int64                  `tfsdk:"precedence"`
	Type                     types.String                 `tfsdk:"type"`
	FilterQuery              types.String                 `tfsdk:"filter_query"`
	InferenceQuery           types.String                 `tfsdk:"inference_query"`
	Review                   types.Bool                   `tfsdk:"review"`
	InferenceSubject         types.String                 `tfsdk:"inference_subject"`
	AlgorithmCorrelation     *AlgorithmCorrelationModelV0 `tfsdk:"algorithm_correlation"`
	MachineLearning          *MachineLearningModel        `tfsdk:"machine_learning"`
}

// AlgorithmCorrelationModelV0 mirrors AlgorithmCorrelationModel from schema version 0, where
// alerts_time_window was a string holding a number of minutes.
type AlgorithmCorrelationModelV0 struct {
	MatchingConditions []MatchingConditionModel `tfsdk:"matching_conditions"`
	AlertsTimeWindow   types.String             `tfsdk:"alerts_time_window"`
	AlertTrigger       *AlertTriggerModel       `tfsdk:"alert_trigger"`
}

// alertCorrelationPolicySchemaV0 describes schema version 0, the shape of every released state
// before alerts_time_window became a number, by calling alertCorrelationPolicySchema with version
// 0 and the string attribute it used to be.
func alertCorrelationPolicySchemaV0() schema.Schema {
	return alertCorrelationPolicySchema(0, schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Time window in minutes for correlating alerts.",
	})
}

// UpgradeState upgrades state written when alerts_time_window was a string (schema version 0) to
// the current schema, where it is a number.
func (r *AlertCorrelationPolicyResource) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	priorSchema := alertCorrelationPolicySchemaV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &priorSchema,
			StateUpgrader: upgradeAlertCorrelationPolicyStateV0toV1,
		},
	}
}

// upgradeAlertCorrelationPolicyStateV0toV1 parses the legacy string alerts_time_window into a
// number and carries every other attribute across unchanged.
func upgradeAlertCorrelationPolicyStateV0toV1(
	ctx context.Context,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	var old AlertCorrelationPolicyModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upgraded := AlertCorrelationPolicyModel{
		Client:                   old.Client,
		Id:                       old.Id,
		Name:                     old.Name,
		EnabledMode:              old.EnabledMode,
		OrganizationMatchingType: old.OrganizationMatchingType,
		IncludedClients:          old.IncludedClients,
		Precedence:               old.Precedence,
		Type:                     old.Type,
		FilterQuery:              old.FilterQuery,
		InferenceQuery:           old.InferenceQuery,
		Review:                   old.Review,
		InferenceSubject:         old.InferenceSubject,
		MachineLearning:          old.MachineLearning,
	}

	if old.AlgorithmCorrelation != nil {
		ac := &AlgorithmCorrelationModel{
			MatchingConditions: old.AlgorithmCorrelation.MatchingConditions,
			AlertTrigger:       old.AlgorithmCorrelation.AlertTrigger,
			AlertsTimeWindow:   types.Int64Null(),
		}

		if !old.AlgorithmCorrelation.AlertsTimeWindow.IsNull() && !old.AlgorithmCorrelation.AlertsTimeWindow.IsUnknown() {
			raw := strings.TrimSpace(old.AlgorithmCorrelation.AlertsTimeWindow.ValueString())
			if raw != "" {
				parsed, err := strconv.ParseInt(raw, 10, 64)
				if err != nil {
					resp.Diagnostics.AddError(
						"upgrade alert_correlation_policy resource state",
						"could not convert legacy alerts_time_window value "+raw+" to a number: "+err.Error(),
					)

					return
				}

				ac.AlertsTimeWindow = types.Int64Value(parsed)
			}
		}

		upgraded.AlgorithmCorrelation = ac
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, upgraded)...)
}
