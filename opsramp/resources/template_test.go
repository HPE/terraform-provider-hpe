// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnitValidateTemplateThresholds(t *testing.T) {
	t.Parallel()

	validSide := TemplateThresholdModel{
		RaiseAlert:          types.BoolValue(true),
		CriticalOperator:    types.StringValue("GREATER_THAN"),
		CriticalThreshold:   types.StringValue("13"),
		CriticalRepeatCount: types.Int64Value(1),
		WarningOperator:     types.StringValue("GREATER_THAN"),
		WarningThreshold:    types.StringValue("12"),
		WarningRepeatCount:  types.Int64Value(1),
	}

	planWith := func(threshold TemplateThresholdModel) *TemplateModel {
		return &TemplateModel{MetricGroups: []TemplateMetricGroupModel{{
			Metrics: []TemplateMetricModel{{Thresholds: []TemplateThresholdModel{threshold}}},
		}}}
	}

	tests := map[string]struct {
		threshold TemplateThresholdModel
		wantError bool
	}{
		"valid": {
			threshold: validSide,
		},
		"raise alert requires critical operator": {
			threshold: func() TemplateThresholdModel {
				v := validSide
				v.CriticalOperator = types.StringValue("")

				return v
			}(),
			wantError: true,
		},
		"raise alert requires warning operator": {
			threshold: func() TemplateThresholdModel {
				v := validSide
				v.WarningOperator = types.StringValue("")

				return v
			}(),
			wantError: true,
		},
		"critical operator requires threshold": {
			threshold: func() TemplateThresholdModel {
				v := validSide
				v.CriticalThreshold = types.StringValue("")

				return v
			}(),
			wantError: true,
		},
		"warning operator requires repeat count": {
			threshold: func() TemplateThresholdModel {
				v := validSide
				v.WarningRepeatCount = types.Int64Null()

				return v
			}(),
			wantError: true,
		},
		"none needs no values": {
			threshold: TemplateThresholdModel{
				RaiseAlert:       types.BoolValue(true),
				CriticalOperator: types.StringValue("NONE"),
				WarningOperator:  types.StringValue("NONE"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := validateTemplateThresholds(planWith(test.threshold))
			if (err != nil) != test.wantError {
				t.Fatalf("expected error=%t, got %v", test.wantError, err)
			}
		})
	}
}

func TestUnitFindAPIMetricGroup(t *testing.T) {
	t.Parallel()

	groups := []client.TemplateMetricGroup{
		{Id: "first-id", Name: "first"},
		{Id: "second-id", Name: "second"},
	}

	byID := findAPIMetricGroup(groups, &TemplateMetricGroupModel{
		Id:   types.StringValue("second-id"),
		Name: types.StringValue("wrong-name"),
	})
	if byID == nil || byID.Name != "second" {
		t.Fatalf("expected second group by ID, got %#v", byID)
	}

	byName := findAPIMetricGroup(groups, &TemplateMetricGroupModel{
		Id:   types.StringNull(),
		Name: types.StringValue("first"),
	})
	if byName == nil || byName.Id != "first-id" {
		t.Fatalf("expected first group by name, got %#v", byName)
	}
}

func TestUnitPopulateTemplateMetricGroupResolvesOptionalObjects(t *testing.T) {
	t.Parallel()

	prior := &TemplateMetricGroupModel{
		ComponentProperties: types.ObjectUnknown(componentPropertiesAttrTypes),
		OtelMetricsConfig: &TemplateOtelMetricsConfigModel{
			ConfigType:         types.StringValue(""),
			ReceiverConfigYaml: types.StringValue(""),
		},
	}

	result := populateTemplateMetricGroup(client.TemplateMetricGroup{}, prior)

	if result.ComponentProperties.IsUnknown() || !result.ComponentProperties.IsNull() {
		t.Fatalf("expected component_properties to be known null, got %s", result.ComponentProperties)
	}
	if result.OtelMetricsConfig == nil {
		t.Fatal("expected structured otel_metrics_config to be preserved")
	}
	if result.OtelMetricsConfig.ConfigType.ValueString() != "" {
		t.Fatalf("expected empty config_type, got %q", result.OtelMetricsConfig.ConfigType.ValueString())
	}
}

func TestUnitValidateThresholdOperatorValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		operator  string
		value     string
		wantError bool
	}{
		{name: "none allows empty", operator: "NONE", value: ""},
		{name: "equals ignore case allows empty", operator: "EQUALS_IGNORE_CASE", value: ""},
		{name: "contains allows empty", operator: "CONTAINS", value: ""},
		{name: "integer operator accepts integer", operator: "GREATER_THAN_EQUAL", value: "10"},
		{name: "integer operator accepts negative integer", operator: "LESS_THAN", value: "-10"},
		{name: "integer operator rejects decimal", operator: "GREATER_THAN", value: "10.5", wantError: true},
		{name: "integer operator rejects list", operator: "LESS_THAN_EQUAL", value: "2,10", wantError: true},
		{name: "in list accepts comma-separated integers", operator: "IN_LIST", value: "2,5,10"},
		{name: "not in list accepts single integer", operator: "NOT_IN_LIST", value: "2"},
		{name: "in list rejects empty segment", operator: "IN_LIST", value: "2,,10", wantError: true},
		{name: "range accepts two integers", operator: "RANGE", value: "2,10"},
		{name: "out of range accepts negative bounds", operator: "OUT_OF_RANGE", value: "-2,10"},
		{name: "range rejects one integer", operator: "RANGE", value: "2", wantError: true},
		{name: "range rejects three integers", operator: "RANGE", value: "2,5,10", wantError: true},
		{name: "regex accepts valid regex", operator: "REGEX_MATCHES", value: "^web-[0-9]+$"},
		{name: "regex rejects invalid regex", operator: "NOT_REGEX_MATCHES", value: "[", wantError: true},
		{name: "unsupported operator rejects", operator: "BAD", value: "10", wantError: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateThresholdOperatorValue("threshold", tt.operator, tt.value)
			if tt.wantError && err == nil {
				t.Fatalf("expected error")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
