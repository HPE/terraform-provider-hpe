// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"
	"github.com/HPE/terraform-provider-hpe/utils/modifiers"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// templateGeneration is fixed: this resource only manages G2 templates (metricGroups).
const templateGeneration = 2

const (
	templateResourceType   = "DEVICE"
	templateSyntheticType  = "SYNTHETIC"
	templateDeploymentType = "Standard"
	templateAlertSubject   = "${severity} - ${component.name} - ${metric.value} (${threshold})"
	templateAlertContent   = "${severity} - ${component.name} Host: ${resource.name} IP: ${resource.ip} Metric Name: ${metric.name} Component: ${component.name} Severity: ${severity} Value: ${metric.value} Reason: ${metric.value} ${threshold}"
)

var (
	integerThresholdPattern      = regexp.MustCompile(`^-?\d+$`)
	integerListThresholdPattern  = regexp.MustCompile(`^-?\d+(,-?\d+)*$`)
	integerRangeThresholdPattern = regexp.MustCompile(`^-?\d+,-?\d+$`)
)

var templateThresholdOperators = map[string]string{
	"REGEX_MATCHES":          "regex",
	"NOT_REGEX_MATCHES":      "regex",
	"GREATER_THAN_EQUAL":     "integer",
	"LESS_THAN":              "integer",
	"GREATER_THAN":           "integer",
	"LESS_THAN_EQUAL":        "integer",
	"EQUAL":                  "integer",
	"EQUALS":                 "integer",
	"EQUALS_IGNORE_CASE":     "none",
	"CONTAINS":               "none",
	"NOT_EQUAL":              "none",
	"NOT_EQUALS":             "none",
	"NOT_EQUALS_IGNORE_CASE": "none",
	"NOT_CONTAINS":           "none",
	"INS":                    "none",
	"NOT_INS":                "none",
	"IN_LIST":                "integer_list",
	"NOT_IN_LIST":            "integer_list",
	"NONE":                   "none",
	"RANGE":                  "integer_range",
	"OUT_OF_RANGE":           "integer_range",
}

var (
	_ resource.Resource                = &TemplateResource{}
	_ resource.ResourceWithImportState = &TemplateResource{}
)

// TemplateResource manages an OpsRamp G2 monitoring template.
type TemplateResource struct {
	BaseResource
}

// TemplateThresholdModel maps a single threshold on a linked template metric.
type TemplateThresholdModel struct {
	CriticalThreshold   types.String `tfsdk:"critical_threshold"`
	CriticalRepeatCount types.Int64  `tfsdk:"critical_repeat_count"`
	WarningThreshold    types.String `tfsdk:"warning_threshold"`
	WarningRepeatCount  types.Int64  `tfsdk:"warning_repeat_count"`
	WarningOperator     types.String `tfsdk:"warning_operator"`
	CriticalOperator    types.String `tfsdk:"critical_operator"`
	RaiseAlert          types.Bool   `tfsdk:"raise_alert"`
	Availability        types.Bool   `tfsdk:"availability"`
	LevelType           types.String `tfsdk:"level_type"`
	CompName            types.String `tfsdk:"comp_name"`
}

// TemplateMetricValueModel maps one OID/formula input feeding a metric's expression.
type TemplateMetricValueModel struct {
	Name          types.String `tfsdk:"name"`
	Oid           types.String `tfsdk:"oid"`
	OidType       types.String `tfsdk:"oid_type"`
	DataType      types.String `tfsdk:"data_type"`
	Formula       types.String `tfsdk:"formula"`
	PollFrequency types.String `tfsdk:"poll_frequency"`
}

// TemplateEnumeratedMapEntryModel maps one raw collected value to a display label.
type TemplateEnumeratedMapEntryModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// TemplateFormatMetricValueModel controls how a metric's collected value is interpreted/displayed.
type TemplateFormatMetricValueModel struct {
	UsedBy        types.String                      `tfsdk:"used_by"`
	EnumeratedMap []TemplateEnumeratedMapEntryModel `tfsdk:"enumerated_map"`
	Alerts        types.Bool                        `tfsdk:"alerts"`
	Graph         types.Bool                        `tfsdk:"graph"`
}

// TemplateFormatMetricValueDetailsModel wraps the display/formatting configuration of a metric.
type TemplateFormatMetricValueDetailsModel struct {
	FormatMetricValue *TemplateFormatMetricValueModel `tfsdk:"format_metric_value"`
	Units             types.String                    `tfsdk:"units"`
	Graph             types.Bool                      `tfsdk:"graph"`
}

// TemplateMetricModel links an existing metric into a metric group, or (for SNMP/JMX/VMware
// metric groups) defines a new metric inline via MetricValues/Expression.
type TemplateMetricModel struct {
	Name                     types.String                           `tfsdk:"name"`
	DisplayName              types.String                           `tfsdk:"display_name"`
	Description              types.String                           `tfsdk:"description"`
	Category                 types.String                           `tfsdk:"category"`
	ThresholdType            types.String                           `tfsdk:"threshold_type"`
	Thresholds               []TemplateThresholdModel               `tfsdk:"thresholds"`
	Availability             types.Bool                             `tfsdk:"availability"`
	AlertOnMissingComponents types.Bool                             `tfsdk:"alert_on_missing_components"`
	AlertSubject             types.String                           `tfsdk:"alert_subject"`
	AlertContent             types.String                           `tfsdk:"alert_content"`
	MetricValues             []TemplateMetricValueModel             `tfsdk:"metric_values"`
	Expression               types.String                           `tfsdk:"expression"`
	FormatMetricValueDetails *TemplateFormatMetricValueDetailsModel `tfsdk:"format_metric_value_details"`
	MetricCategory           types.String                           `tfsdk:"metric_category"`
}

// TemplateConfigurationParameterModel maps a configurable input of a metric group.
// VistaraCustomParamId references a well-known, cross-template catalog parameter (e.g. the
// standard "Remote Server Operating System" prompt); leave it unset for custom parameters.
type TemplateConfigurationParameterModel struct {
	Id                   types.Int64  `tfsdk:"id"`
	VistaraCustomParamId types.Int64  `tfsdk:"vistara_custom_param_id"`
	Name                 types.String `tfsdk:"name"`
	DisplayName          types.String `tfsdk:"display_name"`
	CollectorType        types.String `tfsdk:"collector_type"`
	DefaultValue         types.String `tfsdk:"default_value"`
	Description          types.String `tfsdk:"description"`
	Configurable         types.Bool   `tfsdk:"configurable"`
	InputType            types.String `tfsdk:"input_type"`
	InputOptionValues    types.String `tfsdk:"input_option_values"`
}

// TemplateComponentPropertiesModel describes how component instances are named/identified.
type TemplateComponentPropertiesModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Mode      types.String `tfsdk:"mode"`
	Oid       types.String `tfsdk:"oid"`
	LookupOid types.String `tfsdk:"lookup_oid"`
	Format    types.String `tfsdk:"format"`
}

// TemplateOtelMetricsConfigModel configures OpenTelemetry-based metric collection.
type TemplateOtelMetricsConfigModel struct {
	ConfigType                    types.String `tfsdk:"config_type"`
	ReceiverConfigYaml            types.String `tfsdk:"receiver_config_yaml"`
	AlertDefinitionsYaml          types.String `tfsdk:"alert_definitions_yaml"`
	ExternalExporter              types.Bool   `tfsdk:"external_exporter"`
	ExternalExporterName          types.String `tfsdk:"external_exporter_name"`
	ExternalExporterCmdLineConfig types.String `tfsdk:"external_exporter_cmd_line_config"`
}

var componentPropertiesAttrTypes = map[string]attr.Type{
	"id":         types.StringType,
	"name":       types.StringType,
	"mode":       types.StringType,
	"oid":        types.StringType,
	"lookup_oid": types.StringType,
	"format":     types.StringType,
}

// TemplateDiscoveryFilterModel is a condition used to automatically match component instances,
// when a TemplateInstanceFilterModel's InstanceType is AUTOMATIC_INSTANCE.
type TemplateDiscoveryFilterModel struct {
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
}

// TemplateInstanceFilterModel scopes which discovered component instances a metric group applies to.
type TemplateInstanceFilterModel struct {
	InstanceType     types.String                   `tfsdk:"instance_type"`
	InstanceNames    []types.String                 `tfsdk:"instance_names"`
	DiscoveryFilters []TemplateDiscoveryFilterModel `tfsdk:"discovery_filters"`
}

// TemplateAdditionalPropertyModel is an extra SNMP OID collected alongside a metric group's
// metrics (e.g. for use in component naming or reference), not itself alertable.
type TemplateAdditionalPropertyModel struct {
	Id            types.String                      `tfsdk:"id"`
	Name          types.String                      `tfsdk:"name"`
	OidType       types.String                      `tfsdk:"oid_type"`
	Oid           types.String                      `tfsdk:"oid"`
	Operator      types.String                      `tfsdk:"operator"`
	Operand       types.String                      `tfsdk:"operand"`
	FormatValue   types.String                      `tfsdk:"format_value"`
	EnumeratedMap []TemplateEnumeratedMapEntryModel `tfsdk:"enumerated_map"`
	Units         types.String                      `tfsdk:"units"`
	DefaultValue  types.String                      `tfsdk:"default_value"`
	Description   types.String                      `tfsdk:"description"`
}

// TemplateMetricGroupModel maps a G2 metric group (replaces the legacy "monitors" array).
type TemplateMetricGroupModel struct {
	Id                      types.String                          `tfsdk:"id"`
	Name                    types.String                          `tfsdk:"name"`
	Frequency               types.Int64                           `tfsdk:"frequency"`
	AdapterType             types.String                          `tfsdk:"adapter_type"`
	ApplicationType         types.String                          `tfsdk:"application_type"`
	Script                  types.String                          `tfsdk:"script"`
	ComponentProperties     types.Object                          `tfsdk:"component_properties"`
	IsMultiInstance         types.Bool                            `tfsdk:"is_multi_instance"`
	InputScopeComponent     types.Bool                            `tfsdk:"input_scope_component"`
	DiscoveryMetric         types.String                          `tfsdk:"discovery_metric"`
	ConfigurationParameters []TemplateConfigurationParameterModel `tfsdk:"configuration_parameters"`
	OtelMetricsConfig       *TemplateOtelMetricsConfigModel       `tfsdk:"otel_metrics_config"`
	InstanceFilter          *TemplateInstanceFilterModel          `tfsdk:"instance_filter"`
	AdditionalProperties    []TemplateAdditionalPropertyModel     `tfsdk:"additional_properties"`
	Metrics                 []TemplateMetricModel                 `tfsdk:"metrics"`
}

// TemplateModel maps the top-level hpe_opsramp_template resource.
type TemplateModel struct {
	Id                  types.String               `tfsdk:"id"`
	Client              types.String               `tfsdk:"client"`
	Name                types.String               `tfsdk:"name"`
	Description         types.String               `tfsdk:"description"`
	CollectorType       types.String               `tfsdk:"collector_type"`
	Status              types.String               `tfsdk:"status"`
	EolDate             types.String               `tfsdk:"eol_date"`
	TemplateFamily      types.String               `tfsdk:"template_family"`
	Prerequisites       types.String               `tfsdk:"prerequisites"`
	ApplicationType     types.String               `tfsdk:"application_type"`
	TemplateType        types.String               `tfsdk:"template_type"`
	CreatedDate         types.String               `tfsdk:"created_date"`
	Scope               types.String               `tfsdk:"scope"`
	AssignedToResources types.Bool                 `tfsdk:"assigned_to_resources"`
	SdkTag              types.String               `tfsdk:"sdk_tag"`
	MonTemplateVer      types.Int64                `tfsdk:"mon_template_ver"`
	MetricGroups        []TemplateMetricGroupModel `tfsdk:"metric_groups"`
}

func componentPropertiesModel(value types.Object) *TemplateComponentPropertiesModel {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	attributes := value.Attributes()

	return &TemplateComponentPropertiesModel{
		Id:        attributes["id"].(types.String),
		Name:      attributes["name"].(types.String),
		Mode:      attributes["mode"].(types.String),
		Oid:       attributes["oid"].(types.String),
		LookupOid: attributes["lookup_oid"].(types.String),
		Format:    attributes["format"].(types.String),
	}
}

func componentPropertiesObject(model TemplateComponentPropertiesModel) types.Object {
	return types.ObjectValueMust(componentPropertiesAttrTypes, map[string]attr.Value{
		"id":         model.Id,
		"name":       model.Name,
		"mode":       model.Mode,
		"oid":        model.Oid,
		"lookup_oid": model.LookupOid,
		"format":     model.Format,
	})
}

func NewTemplate() resource.Resource {
	return &TemplateResource{}
}

func (r *TemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (r *TemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	thresholdAttrs := map[string]schema.Attribute{
		"critical_threshold": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The critical threshold value. The API expects threshold values as strings.",
		},
		"critical_repeat_count": schema.Int64Attribute{
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(1),
			MarkdownDescription: "Number of consecutive breaches before raising a critical alert.",
		},
		"warning_threshold": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The warning threshold value. The API expects threshold values as strings.",
		},
		"warning_repeat_count": schema.Int64Attribute{
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(1),
			MarkdownDescription: "Number of consecutive breaches before raising a warning alert.",
		},
		"warning_operator": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("NONE"),
			MarkdownDescription: "Comparison operator applied to the warning threshold.",
			Validators: []validator.String{
				thresholdOperatorValidator(),
			},
		},
		"critical_operator": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("NONE"),
			MarkdownDescription: "Comparison operator applied to the critical threshold (e.g. `GREATER_THAN_EQUAL`, `EQUALS_IGNORE_CASE`).",
			Validators: []validator.String{
				thresholdOperatorValidator(),
			},
		},
		"raise_alert": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
			MarkdownDescription: "Whether breaching this threshold raises an alert.",
		},
		"availability": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether this metric contributes to availability.",
		},
		"level_type": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("TEMPLATE"),
			MarkdownDescription: "The threshold level type.",
			Validators: []validator.String{
				stringvalidator.OneOf("TEMPLATE", "TEMPLATE_COMPONENT", "RESOURCE", "COMPONENT"),
			},
		},
		"comp_name": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "Component name used when `level_type` is `TEMPLATE_COMPONENT` or `COMPONENT`.",
		},
	}

	metricValueAttrs := map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The variable name, referenced by `expression` (e.g. `hrStorageUsed`).",
		},
		"oid": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The SNMP OID collected for this variable.",
		},
		"oid_type": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("SCALAROID"),
			MarkdownDescription: "The OID type (e.g. `SCALAROID`, `TABLEOID`).",
			Validators: []validator.String{
				stringvalidator.OneOf("SCALAROID", "TABLEOID"),
			},
		},
		"data_type": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("NONE"),
			MarkdownDescription: "The collected value's data type (e.g. `COUNTER64`).",
			Validators: []validator.String{
				stringvalidator.OneOf(
					"NONE", "COUNTER64", "COUNTER32", "GAUGE32", "INTEGER", "HEX_TO_STRING", "DATE_AND_TIME",
					"DISPLAY_STRING", "IP_ADDRESS", "MAC_ADDRESS", "OBJECT_IDENTIFIER", "OCTET_STRING",
					"TIME_TICKS", "UNSIGNED32", "HEX_TO_DECIMAL", "HEX_TO_IPADDRESS",
				),
			},
		},
		"formula": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("None"),
			MarkdownDescription: "An optional per-value transformation formula.",
			Validators: []validator.String{
				stringvalidator.OneOf("None", "CounterDelta", "CounterRate", "DELTA", "RATE"),
			},
		},
		"poll_frequency": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("NORMAL"),
			MarkdownDescription: "The polling frequency for this metric value (e.g. `NORMAL`, `LESSFREQUENCY`).",
			Validators: []validator.String{
				stringvalidator.OneOf("NORMAL", "LESSFREQUENCY"),
			},
		},
	}

	enumeratedMapEntryAttrs := map[string]schema.Attribute{
		"key": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The raw collected value.",
		},
		"value": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The display label for `key`.",
		},
	}

	formatMetricValueAttrs := map[string]schema.Attribute{
		"used_by": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "How the enumerated map, if any, is used to interpret the metric value.",
		},
		"enumerated_map": schema.ListNestedAttribute{
			Optional:            true,
			MarkdownDescription: "Maps raw collected values to display labels (e.g. when `used_by` is `ENUMARATED_MAP`).",
			NestedObject: schema.NestedAttributeObject{
				Attributes: enumeratedMapEntryAttrs,
			},
		},
		"alerts": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether alerts are enabled for this metric's formatted value.",
		},
		"graph": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether this formatted value is graphed.",
		},
	}

	formatMetricValueDetailsAttrs := map[string]schema.Attribute{
		"format_metric_value": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          formatMetricValueAttrs,
			MarkdownDescription: "Formatting/interpretation rules for the metric value.",
		},
		"units": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The unit label displayed alongside the metric value (e.g. `MBytes`).",
		},
		"graph": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether this metric is graphed.",
		},
	}

	metricAttrs := map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The monitor name shown for this linked metric in the template.",
		},
		"display_name": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The label shown in the UI for this metric. For metrics referencing a library metric, defaults to the value returned by the API.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
		},
		"description": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "A description of this metric. For metrics referencing a library metric, defaults to the value returned by the API.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
		},
		"category": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The metric category (e.g. `METRIC`, `EVENT`). For metrics referencing a library metric, defaults to the value returned by the API.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
		},
		"threshold_type": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Alerting condition type.",
			Validators: []validator.String{
				stringvalidator.OneOf("CHANGE_DETECTION", "FORECAST", "STATIC"),
			},
		},
		"availability": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether this metric contributes to availability.",
		},
		"alert_on_missing_components": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
			MarkdownDescription: "Whether to alert when components are missing.",
		},
		"thresholds": schema.ListNestedAttribute{
			Optional:            true,
			MarkdownDescription: "Threshold definitions for this metric.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: thresholdAttrs,
			},
		},
		"alert_subject": schema.StringAttribute{
			Optional: true,
			Computed: true,
			MarkdownDescription: "The subject line for alerts related to this metric. " +
				"Omit (do not set to \"\") to use the built-in default.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"alert_content": schema.StringAttribute{
			Optional: true,
			Computed: true,
			MarkdownDescription: "The content/body for alerts related to this metric. " +
				"Omit (do not set to \"\") to use the built-in default.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"metric_values": schema.ListNestedAttribute{
			Optional:            true,
			MarkdownDescription: "Inline OID/formula definitions feeding `expression`, for SNMP/JMX/VMware metric groups.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: metricValueAttrs,
			},
		},
		"expression": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "A formula referencing `metric_values` names, for SNMP/JMX/VMware metric groups.",
		},
		"format_metric_value_details": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          formatMetricValueDetailsAttrs,
			MarkdownDescription: "Display/formatting configuration, for SNMP/JMX/VMware metric groups.",
		},
		"metric_category": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Reserved for future use; not currently populated by the API.",
		},
	}

	configParamAttrs := map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "The parameter ID assigned by the API once saved on the template.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"vistara_custom_param_id": schema.Int64Attribute{
			Optional: true,
			MarkdownDescription: "References a well-known, cross-template catalog parameter (e.g. `738` for the standard " +
				"\"Application Type\" prompt). Leave unset for custom parameters.",
		},
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The parameter name, referenced in the script as `${name}`.",
		},
		"display_name": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The label shown in the UI when applying the template.",
		},
		"collector_type": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The collector type this parameter applies to (e.g. `Application`).",
		},
		"default_value": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The default value for the parameter.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "A description of the parameter.",
		},
		"configurable": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
			MarkdownDescription: "Whether the parameter value can be overridden when applying the template.",
		},
		"input_type": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The input widget type (e.g. `text`, `radio`).",
		},
		"input_option_values": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "JSON-encoded list of selectable values, required when `input_type` is `radio`.",
		},
	}

	componentPropertiesAttrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The component property's identifier, assigned by the API.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
		},
		"name": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The property used to name discovered components (e.g. `Name`).",
		},
		"mode": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "How the component name is resolved (e.g. `VALUE`, `INDEX`, `INDEX_LOOKUP`, `INDEX_VALUE`).",
			Validators: []validator.String{
				stringvalidator.OneOf("VALUE", "INDEX", "INDEX_LOOKUP", "INDEX_VALUE"),
			},
		},
		"oid": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The SNMP OID used to resolve the component name, for SNMP-based metric groups.",
		},
		"lookup_oid": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The SNMP OID used to look up the component index, for SNMP-based metric groups.",
		},
		"format": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("NONE"),
			MarkdownDescription: "The format applied to the resolved component name (e.g. `NONE`, `COUNTER64`).",
			Validators: []validator.String{
				stringvalidator.OneOf(
					"NONE", "COUNTER64", "COUNTER32", "GAUGE32", "INTEGER", "HEX_TO_STRING", "DATE_AND_TIME",
					"DISPLAY_STRING", "IP_ADDRESS", "MAC_ADDRESS", "OBJECT_IDENTIFIER", "OCTET_STRING",
					"TIME_TICKS", "UNSIGNED32", "HEX_TO_DECIMAL", "HEX_TO_IPADDRESS",
				),
			},
		},
	}

	otelMetricsConfigAttrs := map[string]schema.Attribute{
		"config_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("", "nativeOtel", "prometheus"),
			},
			MarkdownDescription: "The OpenTelemetry configuration type (e.g., `nativeOtel` or `prometheus`).",
		},
		"receiver_config_yaml": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The OpenTelemetry receiver configuration, as YAML.",
		},
		"alert_definitions_yaml": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Alert definitions for OpenTelemetry-collected metrics, as YAML.",
		},
		"external_exporter": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "Whether metrics are exported by an external OpenTelemetry exporter rather than collected directly.",
		},
		"external_exporter_name": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString(""),
			Validators: []validator.String{
				stringvalidator.OneOf("", "windows_exporter.exe", "node_exporter"),
			},
			MarkdownDescription: "The name of the external OpenTelemetry exporter, if `external_exporter` is `true`.",
		},
		"external_exporter_cmd_line_config": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The command line configuration for the external OpenTelemetry exporter, if `external_exporter` is `true`.",
		},
	}

	instanceFilterAttrs := map[string]schema.Attribute{
		"instance_type": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "Which discovered instances this metric group applies to.",
			Validators: []validator.String{
				stringvalidator.OneOf("", "ALL", "MANUAL_INSTANCE", "AUTOMATIC_INSTANCE"),
			},
		},
		"instance_names": schema.ListAttribute{
			Optional:            true,
			ElementType:         types.StringType,
			MarkdownDescription: "The specific instance names this metric group applies to. Required when `instance_type` is `MANUAL_INSTANCE`.",
		},
		"discovery_filters": schema.ListNestedAttribute{
			Optional:            true,
			MarkdownDescription: "Conditions used to automatically match component instances. Required when `instance_type` is `AUTOMATIC_INSTANCE`.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"operator": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The comparison operator (e.g. `CONTAINS`, `NOT_REGEX_MATCHES`).",
						Validators: []validator.String{
							thresholdOperatorValidator(),
						},
					},
					"value": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The value compared against the discovered instance name.",
					},
				},
			},
		},
	}

	additionalPropertyAttrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The additional property's identifier, assigned by the API.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
		},
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The property name (e.g. `ifAdminStatus`).",
		},
		"oid_type": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("SCALAROID"),
			MarkdownDescription: "The OID type (e.g. `SCALAROID`, `TABULAROID`).",
			Validators: []validator.String{
				stringvalidator.OneOf("SCALAROID", "TABULAROID"),
			},
		},
		"oid": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			MarkdownDescription: "The SNMP OID collected for this property.",
		},
		"operator": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "An optional comparison operator applied to the collected value (e.g. `NONE`, `ADD`, `MULTIPLY`, `DEVIDE`).", //nolint:misspell // DEVIDE is the literal value expected by the API
			Default:             stringdefault.StaticString("NONE"),
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("NONE", "ADD", "MULTIPLY", "DEVIDE"), //nolint:misspell // DEVIDE is the literal value expected by the API
			},
		},
		"operand": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "An optional operand paired with `operator`.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
		},
		"format_value": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "How the collected value is formatted (e.g. `NONE`, `ENUMARATED_MAP`).",
			Default:             stringdefault.StaticString("NONE"),
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("NONE", "ENUMARATED_MAP"),
			},
		},
		"enumerated_map": schema.ListNestedAttribute{
			Optional:            true,
			MarkdownDescription: "Maps raw collected values to display labels.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: enumeratedMapEntryAttrs,
			},
		},
		"units": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The unit label displayed alongside the collected value.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
		},
		"default_value": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The default value for this property.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
		},
		"description": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "A description of this property.",
			PlanModifiers: []planmodifier.String{
				modifiers.StringUnknownIfEmpty(),
			},
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp G2 monitoring template. G2 templates describe their monitors " +
			"using `metric_groups`. This resource does not support generation-1 templates.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier (UUID) of the template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique ID of the client (sub-tenant). If not provided, the template is managed at the provider tenant level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The template name.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A description of the template.",
			},
			"collector_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The collector type. Valid values: `OpsRamp Agent`, `OpsRamp Gateway`, `Synthetics`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The template activity status. Valid values: `ACTIVE`, `INACTIVE`, `EOL`. Defaults to `ACTIVE`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "INACTIVE", "EOL"),
				},
				Default: stringdefault.StaticString("ACTIVE"),
			},
			"eol_date": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The end-of-life date for the template, mandatory when status = `EOL`.",
			},
			"template_family": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The template family grouping.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"prerequisites": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Pre-requisite details needed before applying the template.",
			},
			"template_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Template category.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"OTEL Based Templates-Agent",
						"APPLICATION",
						"Remote Script Executor-Gateway",
						"Remote Script Executor-Agent",
						"SNMP",
						"SNMP Extended",
						"JMX",
						"VMware",
					),
				},
			},
			"application_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				MarkdownDescription: "Application type(s) associated with the template, required when `template_type` = " +
					"`Application`. Sent to the API as a single-element `applicationType` array. Valid values are the `name` field returned by `GET /api/v2/monitoring/adapterTypes` " +
					"(e.g. `APACHE_TOMCAT` for Apache Tomcat templates).",
			},
			"created_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when the template was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The template scope (e.g. `GLOBAL`, `PARTNER`, `CLIENT`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"assigned_to_resources": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the template is currently assigned to any resources.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sdk_tag": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SDK tag associated with exported template payloads.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mon_template_ver": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The monitoring template version.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"metric_groups": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "The G2 metric groups that make up this template.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The metric group ID assigned by the API.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The metric group name.",
						},
						"frequency": schema.Int64Attribute{
							Required: true,
							Validators: []validator.Int64{
								int64validator.OneOf(1, 2, 5, 10, 15, 30, 60, 120, 1440),
							},
							MarkdownDescription: "Polling frequency, in minutes.",
						},
						"adapter_type": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
							MarkdownDescription: "The adapter type (e.g. `APPLICATION`).",
						},
						"application_type": schema.StringAttribute{
							Optional: true,
							Computed: true,
							MarkdownDescription: "The application type for script-based (RSE) metric groups. Valid values are the " +
								"`name` field returned by `GET /api/v2/monitoring/adapterTypes` (e.g. `RSE`).",
							PlanModifiers: []planmodifier.String{
								modifiers.StringUnknownIfEmpty(),
							},
						},
						"script": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
							MarkdownDescription: "The script content executed to collect this metric group's data, for RSE metric groups.",
						},
						"component_properties": schema.SingleNestedAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Controls how component instances discovered by this metric group are named.",
							Attributes:          componentPropertiesAttrs,
							PlanModifiers: []planmodifier.Object{
								objectplanmodifier.UseStateForUnknown(),
							},
						},
						"is_multi_instance": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							MarkdownDescription: "Whether this metric group can report multiple component instances. When `true`, `component_properties`, `additional_properties`, and `instance_filter` are required.",
						},
						"input_scope_component": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether users can scope this metric group to specific components when applying the template.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"discovery_metric": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
							MarkdownDescription: "The metric used to auto-discover component instances for this metric group.",
						},
						"configuration_parameters": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Configurable input parameters for this metric group.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: configParamAttrs,
							},
						},
						"otel_metrics_config": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "OpenTelemetry-based metric collection configuration.",
							Attributes:          otelMetricsConfigAttrs,
						},
						"instance_filter": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Scopes which discovered component instances this metric group applies to, for SNMP/JMX/VMware metric groups.",
							Attributes:          instanceFilterAttrs,
						},
						"additional_properties": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Extra SNMP OIDs collected alongside this metric group's metrics.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: additionalPropertyAttrs,
							},
						},
						"metrics": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "The metrics collected by this metric group.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: metricAttrs,
							},
						},
					},
				},
			},
		},
	}
}

func (r *TemplateResource) getTenantId(clientId types.String) string {
	if !clientId.IsNull() && clientId.ValueString() != "" {
		return clientId.ValueString()
	}

	return r.apiClient.TenantId
}

func (r *TemplateResource) validateApplicationTypes(plan *TemplateModel) error {
	if plan.TemplateType.ValueString() != "APPLICATION" {
		// if the template type is not "APPLICATION", the application type should be ignored
		return nil
	}

	if plan.TemplateType.ValueString() == "APPLICATION" && plan.ApplicationType.ValueString() == "" {
		return fmt.Errorf("application_type is required when template_type is \"APPLICATION\"")
	}

	adapterTypes, err := r.apiClient.GetAdapterTypes()
	if err != nil {
		return fmt.Errorf("failed to retrieve adapter types: %w", err)
	}

	validNames := make(map[string]struct{}, len(adapterTypes))
	validList := make([]string, 0, len(adapterTypes))
	for _, adapterType := range adapterTypes {
		if adapterType.Name == "" {
			continue
		}
		validNames[adapterType.Name] = struct{}{}
		validList = append(validList, adapterType.Name)
	}

	invalid := []string{}
	addInvalidApplicationType("application_type", plan.ApplicationType, validNames, &invalid)

	for i, metricGroup := range plan.MetricGroups {
		path := fmt.Sprintf("metric_groups[%d].application_type", i)
		addInvalidApplicationType(path, metricGroup.ApplicationType, validNames, &invalid)
	}

	if len(invalid) > 0 {
		return fmt.Errorf(
			"invalid application_type value(s): %s. Valid values from /api/v2/monitoring/adapterTypes are: %s",
			strings.Join(invalid, ", "),
			strings.Join(validList, ", "),
		)
	}

	return nil
}

func addInvalidApplicationType(
	attributePath string,
	applicationType types.String,
	validNames map[string]struct{},
	invalid *[]string,
) {
	if applicationType.IsNull() || applicationType.IsUnknown() || applicationType.ValueString() == "" {
		return
	}

	value := applicationType.ValueString()
	if _, ok := validNames[value]; ok {
		return
	}

	*invalid = append(*invalid, fmt.Sprintf("%s=%q", attributePath, value))
}

func validateTemplateStatus(plan *TemplateModel) error {
	status := plan.Status.ValueString()
	eolDate := plan.EolDate.ValueString()

	if status == "EOL" && eolDate == "" {
		return fmt.Errorf("eol_date is required when status is EOL")
	}

	if eolDate != "" {
		if _, err := time.Parse(time.RFC3339Nano, eolDate); err != nil {
			return fmt.Errorf("eol_date must be an RFC3339 timestamp like 2026-09-16T22:00:00.000Z: %w", err)
		}
	}

	return nil
}

func validateTemplateOtelMetricsConfig(plan *TemplateModel) error {
	for i, metricGroup := range plan.MetricGroups {
		config := metricGroup.OtelMetricsConfig
		if config == nil || config.ConfigType.ValueString() != "prometheus" {
			continue
		}

		path := fmt.Sprintf("metric_groups[%d].otel_metrics_config", i)
		if config.ExternalExporter.IsNull() || config.ExternalExporter.IsUnknown() {
			return fmt.Errorf("%s.external_exporter is required when config_type is prometheus", path)
		}

		if !config.ExternalExporter.ValueBool() {
			continue
		}

		if config.ExternalExporterName.ValueString() == "" {
			return fmt.Errorf("%s.external_exporter_name is required when external_exporter is true", path)
		}

		if config.ExternalExporterCmdLineConfig.ValueString() == "" {
			return fmt.Errorf(
				"%s.external_exporter_cmd_line_config is required when external_exporter is true",
				path,
			)
		}
	}

	return nil
}

// validateTemplateMetricGroups enforces multi-instance metric groups carry the metadata
// the API needs to identify each instance, and that instance_filter's own fields are
// consistent with its instance_type.
func validateTemplateMetricGroups(plan *TemplateModel) error {
	for groupIndex, metricGroup := range plan.MetricGroups {
		path := fmt.Sprintf("metric_groups[%d]", groupIndex)

		if !metricGroup.IsMultiInstance.ValueBool() {
			continue
		}

		componentProperties := componentPropertiesModel(metricGroup.ComponentProperties)
		if componentProperties == nil {
			return fmt.Errorf("%s.component_properties is required when is_multi_instance is true", path)
		}

		if componentProperties.Name.ValueString() == "" || componentProperties.Mode.ValueString() == "" {
			return fmt.Errorf(
				"%s.component_properties.name and mode are required when is_multi_instance is true",
				path,
			)
		}

		if metricGroup.AdditionalProperties == nil {
			return fmt.Errorf("%s.additional_properties is required (can be an empty list) when is_multi_instance is true", path)
		}

		filter := metricGroup.InstanceFilter
		if filter == nil {
			return fmt.Errorf("%s.instance_filter is required when is_multi_instance is true", path)
		}

		filterPath := path + ".instance_filter"
		switch filter.InstanceType.ValueString() {
		case "MANUAL_INSTANCE":
			if len(filter.InstanceNames) == 0 {
				return fmt.Errorf("%s.instance_names is required when instance_type is MANUAL_INSTANCE", filterPath)
			}
		case "AUTOMATIC_INSTANCE":
			if len(filter.DiscoveryFilters) == 0 {
				return fmt.Errorf("%s.discovery_filters is required when instance_type is AUTOMATIC_INSTANCE", filterPath)
			}
		case "ALL":
			// No further requirements.
		default:
			return fmt.Errorf("%s.instance_type must be one of ALL, MANUAL_INSTANCE, AUTOMATIC_INSTANCE", filterPath)
		}
	}

	return nil
}

func validateTemplateThresholds(plan *TemplateModel) error {
	for groupIndex, metricGroup := range plan.MetricGroups {
		for metricIndex, metric := range metricGroup.Metrics {
			for thresholdIndex, threshold := range metric.Thresholds {
				path := fmt.Sprintf(
					"metric_groups[%d].metrics[%d].thresholds[%d]",
					groupIndex,
					metricIndex,
					thresholdIndex,
				)

				if threshold.RaiseAlert.ValueBool() {
					if threshold.CriticalOperator.ValueString() == "" {
						return fmt.Errorf("%s.critical_operator is required when raise_alert is true", path)
					}
					if threshold.WarningOperator.ValueString() == "" {
						return fmt.Errorf("%s.warning_operator is required when raise_alert is true", path)
					}
				}

				if threshold.LevelType.ValueString() == "TEMPLATE_COMPONENT" && threshold.CompName.ValueString() == "" {
					return fmt.Errorf("%s.comp_name is required when level_type is TEMPLATE_COMPONENT", path)
				}

				if err := validateThresholdSide(
					path,
					"critical",
					threshold.CriticalOperator,
					threshold.CriticalThreshold,
					threshold.CriticalRepeatCount,
				); err != nil {
					return err
				}

				if err := validateThresholdSide(
					path,
					"warning",
					threshold.WarningOperator,
					threshold.WarningThreshold,
					threshold.WarningRepeatCount,
				); err != nil {
					return err
				}

				if err := validateThresholdOperatorValue(
					path+".critical_threshold",
					threshold.CriticalOperator.ValueString(),
					threshold.CriticalThreshold.ValueString(),
				); err != nil {
					return err
				}

				if err := validateThresholdOperatorValue(
					path+".warning_threshold",
					threshold.WarningOperator.ValueString(),
					threshold.WarningThreshold.ValueString(),
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func validateThresholdSide(
	path string,
	severity string,
	operator types.String,
	threshold types.String,
	repeatCount types.Int64,
) error {
	operatorValue := operator.ValueString()
	if operatorValue == "" || operatorValue == "NONE" {
		return nil
	}

	if threshold.IsNull() || threshold.IsUnknown() || threshold.ValueString() == "" {
		return fmt.Errorf("%s.%s_threshold is required when %s_operator is %s", path, severity, severity, operatorValue)
	}

	if repeatCount.IsNull() || repeatCount.IsUnknown() {
		return fmt.Errorf("%s.%s_repeat_count is required when %s_operator is %s", path, severity, severity, operatorValue)
	}

	return nil
}

func validateThresholdOperatorValue(attributePath string, operator string, value string) error {
	if operator == "" {
		return nil
	}

	valueKind, ok := templateThresholdOperators[operator]
	if !ok {
		return fmt.Errorf("%s has unsupported operator %q", attributePath, operator)
	}

	switch valueKind {
	case "regex":
		if value == "" {
			return fmt.Errorf("%s must be a regex pattern when operator is %s", attributePath, operator)
		}
		if _, err := regexp.Compile(value); err != nil {
			return fmt.Errorf("%s must be a valid regex when operator is %s: %w", attributePath, operator, err)
		}
	case "integer":
		if !integerThresholdPattern.MatchString(value) {
			return fmt.Errorf("%s must be one integer when operator is %s", attributePath, operator)
		}
	case "integer_list":
		if !integerListThresholdPattern.MatchString(value) {
			return fmt.Errorf("%s must be a comma-separated integer list when operator is %s", attributePath, operator)
		}
	case "integer_range":
		if !integerRangeThresholdPattern.MatchString(value) {
			return fmt.Errorf("%s must be two comma-separated integers when operator is %s", attributePath, operator)
		}
	case "none":
		return nil
	default:
		return fmt.Errorf("%s has unsupported validation kind %q for operator %s", attributePath, valueKind, operator)
	}

	return nil
}

func thresholdOperatorValidator() validator.String {
	operators := []string{""}
	for operator := range templateThresholdOperators {
		operators = append(operators, operator)
	}

	return stringvalidator.OneOf(operators...)
}

func templateFamily(plan *TemplateModel) string {
	if !plan.TemplateFamily.IsNull() && !plan.TemplateFamily.IsUnknown() && plan.TemplateFamily.ValueString() != "" {
		return plan.TemplateFamily.ValueString()
	}

	if plan.TemplateType.IsNull() || plan.TemplateType.IsUnknown() || plan.TemplateType.ValueString() == "" {
		return ""
	}

	return plan.TemplateType.ValueString()
}

func (r *TemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TemplateModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateTemplateStatus(&plan); err != nil {
		resp.Diagnostics.AddError("Invalid Template Status", err.Error())

		return
	}

	if err := validateTemplateOtelMetricsConfig(&plan); err != nil {
		resp.Diagnostics.AddError("Invalid OpenTelemetry Metrics Config", err.Error())

		return
	}

	if err := validateTemplateThresholds(&plan); err != nil {
		resp.Diagnostics.AddError("Invalid Template Threshold", err.Error())

		return
	}

	if err := validateTemplateMetricGroups(&plan); err != nil {
		resp.Diagnostics.AddError("Invalid Metric Group", err.Error())

		return
	}

	if err := r.validateApplicationTypes(&plan); err != nil {
		resp.Diagnostics.AddError("Invalid Application Type", err.Error())

		return
	}

	tenantId := r.getTenantId(plan.Client)
	apiReq := buildTemplateRequest(&plan)

	created, err := r.apiClient.CreateTemplate(tenantId, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Template", err.Error())

		return
	}

	plan.Id = types.StringValue(created.TemplateId)

	// The create response is minimal (id + template-name); fetch the full
	// representation so the resource can populate computed attributes.
	fullTemplate, err := r.apiClient.GetTemplate(tenantId, created.TemplateId)
	if err != nil {
		resp.Diagnostics.AddError("Error Fetching Created Template", err.Error())

		return
	}

	// The OpsRamp API can return metricGroups in a different order than submitted.
	fullTemplate.MetricGroups = orderMetricGroupsByReference(fullTemplate.MetricGroups, metricGroupRequestNames(apiReq.MetricGroups))

	populateTemplateModelFromAPI(&plan, fullTemplate)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *TemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	existing, err := r.apiClient.GetTemplate(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Error Reading Template", fmt.Sprintf("Could not read template %s: %s", state.Id.ValueString(), err))

		return
	}

	// The OpsRamp API can return metricGroups in a different order than they were previously stored.
	existing.MetricGroups = orderMetricGroupsByReference(existing.MetricGroups, metricGroupModelNames(state.MetricGroups))

	refreshed := TemplateModel{Client: state.Client}
	populateTemplateModelFromAPI(&refreshed, existing)

	diags = resp.State.Set(ctx, &refreshed)
	resp.Diagnostics.Append(diags...)
}

func (r *TemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TemplateModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	var state TemplateModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateTemplateStatus(&plan); err != nil {
		resp.Diagnostics.AddError("Invalid Template Status", err.Error())

		return
	}

	if err := validateTemplateOtelMetricsConfig(&plan); err != nil {
		resp.Diagnostics.AddError("Invalid OpenTelemetry Metrics Config", err.Error())

		return
	}

	if err := validateTemplateThresholds(&plan); err != nil {
		resp.Diagnostics.AddError("Invalid Template Threshold", err.Error())

		return
	}

	if err := validateTemplateMetricGroups(&plan); err != nil {
		resp.Diagnostics.AddError("Invalid Metric Group", err.Error())

		return
	}

	if err := r.validateApplicationTypes(&plan); err != nil {
		resp.Diagnostics.AddError("Invalid Application Type", err.Error())

		return
	}

	tenantId := r.getTenantId(plan.Client)
	requestPlan := plan
	existing, err := r.apiClient.GetTemplate(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Template",
			fmt.Sprintf("Could not read template %s before update: %s", state.Id.ValueString(), err),
		)

		return
	}
	enrichTemplateRequestIDs(&requestPlan, existing)
	apiReq := buildTemplateRequest(&requestPlan)

	updated, err := r.apiClient.UpdateTemplate(tenantId, state.Id.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Template", fmt.Sprintf("Could not update template %s: %s", state.Id.ValueString(), err))

		return
	}

	updated = orderTemplateResponse(updated, apiReq)

	refreshed := TemplateModel{Client: plan.Client}
	refreshed.Id = state.Id
	populateTemplateModelFromAPI(&refreshed, updated)

	diags = resp.State.Set(ctx, &refreshed)
	resp.Diagnostics.Append(diags...)
}

// orderMetricGroupsByReference reorders API-returned metric groups to match the order of names in
// reference, since the OpsRamp API can return metricGroups in an order different from how they were submitted.
// Groups whose name isn't found in reference are appended in their original order.
func orderMetricGroupsByReference(groups []client.TemplateMetricGroup, reference []string) []client.TemplateMetricGroup {
	ordered := make([]client.TemplateMetricGroup, 0, len(groups))
	used := make([]bool, len(groups))

	for _, name := range reference {
		for i := range groups {
			if used[i] || groups[i].Name != name {
				continue
			}

			ordered = append(ordered, groups[i])
			used[i] = true

			break
		}
	}

	for i := range groups {
		if !used[i] {
			ordered = append(ordered, groups[i])
		}
	}

	return ordered
}

func metricGroupRequestNames(groups []client.TemplateMetricGroup) []string {
	names := make([]string, len(groups))
	for i, group := range groups {
		names[i] = group.Name
	}

	return names
}

func metricGroupModelNames(groups []TemplateMetricGroupModel) []string {
	names := make([]string, len(groups))
	for i, group := range groups {
		names[i] = group.Name.ValueString()
	}

	return names
}

func orderTemplateResponse(response *client.Template, request client.TemplateRequest) *client.Template {
	ordered := *response
	ordered.MetricGroups = make([]client.TemplateMetricGroup, 0, len(request.MetricGroups))

	for _, requestedGroup := range request.MetricGroups {
		actualGroup := findRequestedMetricGroup(response.MetricGroups, requestedGroup)
		if actualGroup == nil {
			continue
		}

		group := *actualGroup
		group.Metrics = orderMetricResponse(actualGroup.Metrics, requestedGroup.Metrics)
		ordered.MetricGroups = append(ordered.MetricGroups, group)
	}

	return &ordered
}

func findRequestedMetricGroup(
	groups []client.TemplateMetricGroup, requested client.TemplateMetricGroup,
) *client.TemplateMetricGroup {
	for index := range groups {
		if requested.Id != "" && groups[index].Id == requested.Id {
			return &groups[index]
		}
		if requested.Id == "" && groups[index].Name == requested.Name {
			return &groups[index]
		}
	}

	return nil
}

func orderMetricResponse(actual []client.TemplateMetric, requested []client.TemplateMetric) []client.TemplateMetric {
	ordered := make([]client.TemplateMetric, 0, len(requested))
	for _, requestedMetric := range requested {
		for _, actualMetric := range actual {
			if actualMetric.Name == requestedMetric.Name {
				ordered = append(ordered, actualMetric)

				break
			}
		}
	}

	return ordered
}

func enrichTemplateRequestIDs(plan *TemplateModel, existing *client.Template) {
	for groupIndex := range plan.MetricGroups {
		group := &plan.MetricGroups[groupIndex]
		apiGroup := findAPIMetricGroup(existing.MetricGroups, group)
		if apiGroup == nil {
			continue
		}

		if group.Id.IsNull() || group.Id.IsUnknown() || group.Id.ValueString() == "" {
			group.Id = types.StringValue(apiGroup.Id)
		}

		for parameterIndex := range group.ConfigurationParameters {
			parameter := &group.ConfigurationParameters[parameterIndex]
			if !parameter.Id.IsNull() && !parameter.Id.IsUnknown() && parameter.Id.ValueInt64() != 0 {
				continue
			}

			for _, apiParameter := range apiGroup.ConfigurationParameters {
				if apiParameter.Name == parameter.Name.ValueString() {
					parameter.Id = types.Int64Value(int64(apiParameter.Id))

					break
				}
			}
		}
	}
}

func (r *TemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TemplateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	err := r.apiClient.DeleteTemplate(tenantId, state.Id.ValueString())
	if err != nil && !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
		resp.Diagnostics.AddError("Error Deleting Template", fmt.Sprintf("Could not delete template %s: %s", state.Id.ValueString(), err))

		return
	}
}

// ImportState handles importing an existing template.
// Import ID format: <template_id> or <client_id>:<template_id> (MSP only)
func (r *TemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parsed, err := r.ParseImportID(req.ID, 1)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())

		return
	}

	tenantId := r.TenantForImport(parsed)
	templateId := parsed.Parts[0]

	existing, err := r.apiClient.GetTemplate(tenantId, templateId)
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Template", fmt.Sprintf("Could not import template: %s", err))

		return
	}

	var state TemplateModel
	state.Client = parsed.Client
	populateTemplateModelFromAPI(&state, existing)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// buildTemplateRequest converts the Terraform plan into an API request. This resource is
// G2-only, so generation/monitorType/isG2MonitorType are always hardcoded rather than
// exposed as configurable attributes. applicationType/metrics/eventLogMetrics/otherMetrics
// are sent as empty arrays (never null) to match what the API expects.
func buildTemplateRequest(plan *TemplateModel) client.TemplateRequest {
	resourceType := templateResourceType

	apiReq := client.TemplateRequest{
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		ResourceType: resourceType,
		Status: client.TemplateStatus{
			State:   plan.Status.ValueString(),
			EolDate: plan.EolDate.ValueString(),
		},
		Tags:            "",
		Generation:      templateGeneration,
		TemplateFamily:  templateFamily(plan),
		Notes:           "",
		Prerequisites:   plan.Prerequisites.ValueString(),
		DeploymentType:  templateDeploymentType,
		ApplicationType: []string{},
		MonitorType:     "METRICS",
		Metrics:         []interface{}{},
		EventLogMetrics: []interface{}{},
		OtherMetrics:    []interface{}{},
		ClonedTemplate:  false,
	}

	if plan.ApplicationType.ValueString() != "" {
		apiReq.ApplicationType = append(apiReq.ApplicationType, plan.ApplicationType.ValueString())
	}

	if plan.TemplateType.ValueString() != "" {
		apiReq.TemplateType = append(apiReq.TemplateType, plan.TemplateType.ValueString())
	}

	for _, mg := range plan.MetricGroups {
		apiReq.MetricGroups = append(apiReq.MetricGroups, buildTemplateMetricGroup(plan, mg))
	}

	return apiReq
}

func buildTemplateMetricGroup(plan *TemplateModel, mg TemplateMetricGroupModel) client.TemplateMetricGroup {
	group := client.TemplateMetricGroup{
		Id:                   mg.Id.ValueString(),
		Name:                 mg.Name.ValueString(),
		Frequency:            int(mg.Frequency.ValueInt64()),
		AdapterType:          metricGroupAdapterType(plan, mg),
		ApplicationType:      metricGroupApplicationType(plan, mg),
		Script:               mg.Script.ValueString(),
		DiscoveryMetric:      mg.DiscoveryMetric.ValueString(),
		IsMultiInstance:      mg.IsMultiInstance.ValueBool(),
		AdditionalProperties: []client.TemplateAdditionalProperty{},
		IsG2MonitorType:      true,
	}

	componentProperties := componentPropertiesModel(mg.ComponentProperties)
	if componentProperties != nil {
		group.ComponentProperties = &client.TemplateComponentProperties{
			Id:        componentProperties.Id.ValueString(),
			Name:      componentProperties.Name.ValueString(),
			Mode:      componentProperties.Mode.ValueString(),
			Oid:       componentProperties.Oid.ValueString(),
			LookupOid: componentProperties.LookupOid.ValueString(),
			Format:    componentProperties.Format.ValueString(),
		}
	}

	otelConfig := mg.OtelMetricsConfig
	if otelConfig != nil {
		group.OtelMetricsConfig = &client.TemplateOtelMetricsConfig{
			ConfigType:                    otelConfig.ConfigType.ValueString(),
			ReceiverConfigYaml:            otelConfig.ReceiverConfigYaml.ValueString(),
			AlertDefinitionsYaml:          otelConfig.AlertDefinitionsYaml.ValueString(),
			ExternalExporter:              otelConfig.ExternalExporter.ValueBool(),
			ExternalExporterName:          otelConfig.ExternalExporterName.ValueString(),
			ExternalExporterCmdLineConfig: otelConfig.ExternalExporterCmdLineConfig.ValueString(),
		}
	}

	if filter := mg.InstanceFilter; filter != nil {
		group.InstanceFilter = &client.TemplateInstanceFilter{
			InstanceType:     filter.InstanceType.ValueString(),
			DiscoveryFilters: buildDiscoveryFilters(filter.DiscoveryFilters),
			InstanceNames:    stringValuesToStrings(filter.InstanceNames),
		}
	}

	for _, ap := range mg.AdditionalProperties {
		group.AdditionalProperties = append(group.AdditionalProperties, client.TemplateAdditionalProperty{
			Id:            ap.Id.ValueString(),
			Name:          ap.Name.ValueString(),
			OidType:       ap.OidType.ValueString(),
			Oid:           ap.Oid.ValueString(),
			Operator:      ap.Operator.ValueString(),
			Operand:       ap.Operand.ValueString(),
			FormatValue:   ap.FormatValue.ValueString(),
			EnumeratedMap: buildEnumeratedMap(ap.EnumeratedMap),
			Units:         ap.Units.ValueString(),
			DefaultValue:  ap.DefaultValue.ValueString(),
			Description:   ap.Description.ValueString(),
		})
	}

	for _, cp := range mg.ConfigurationParameters {
		group.ConfigurationParameters = append(group.ConfigurationParameters, client.TemplateConfigurationParameter{
			Id:                   int(cp.Id.ValueInt64()),
			VistaraCustomParamId: int(cp.VistaraCustomParamId.ValueInt64()),
			Name:                 cp.Name.ValueString(),
			DisplayName:          cp.DisplayName.ValueString(),
			CollectorType:        cp.CollectorType.ValueString(),
			DefaultValue:         cp.DefaultValue.ValueString(),
			Description:          cp.Description.ValueString(),
			Configurable:         cp.Configurable.ValueBool(),
			Mandatory:            true,
			InputType:            cp.InputType.ValueString(),
			InputOptionValues:    cp.InputOptionValues.ValueString(),
		})
	}

	for _, m := range mg.Metrics {
		group.Metrics = append(group.Metrics, buildTemplateMetric(m))
	}

	return group
}

func stringValuesToStrings(values []types.String) []string {
	result := make([]string, 0, len(values))
	for _, v := range values {
		result = append(result, v.ValueString())
	}

	return result
}

func buildDiscoveryFilters(filters []TemplateDiscoveryFilterModel) []client.TemplateDiscoveryFilter {
	result := make([]client.TemplateDiscoveryFilter, 0, len(filters))
	for _, f := range filters {
		result = append(result, client.TemplateDiscoveryFilter{
			Operator: f.Operator.ValueString(),
			Value:    f.Value.ValueString(),
		})
	}

	return result
}

func buildEnumeratedMap(entries []TemplateEnumeratedMapEntryModel) []client.TemplateEnumeratedMapEntry {
	result := make([]client.TemplateEnumeratedMapEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, client.TemplateEnumeratedMapEntry{
			Key:   entry.Key.ValueString(),
			Value: entry.Value.ValueString(),
		})
	}

	return result
}

func metricGroupAdapterType(plan *TemplateModel, mg TemplateMetricGroupModel) string {
	if !mg.AdapterType.IsNull() && !mg.AdapterType.IsUnknown() && mg.AdapterType.ValueString() != "" {
		return mg.AdapterType.ValueString()
	}

	family := templateFamily(plan)
	if family == "APPLICATION" || strings.Contains(family, "Remote Script Executor") || strings.Contains(family, "OTEL") {
		return "APPLICATION"
	}

	return ""
}

func metricGroupApplicationType(plan *TemplateModel, mg TemplateMetricGroupModel) string {
	if !mg.ApplicationType.IsNull() && !mg.ApplicationType.IsUnknown() && mg.ApplicationType.ValueString() != "" {
		return mg.ApplicationType.ValueString()
	}

	family := templateFamily(plan)
	if strings.Contains(family, "Remote Script Executor") {
		return "RSE"
	}
	if strings.Contains(family, "OTEL") {
		return "OTEL_METRICS"
	}
	if !plan.ApplicationType.IsNull() && !plan.ApplicationType.IsUnknown() && plan.ApplicationType.ValueString() != "" {
		return plan.ApplicationType.ValueString()
	}

	return ""
}

func buildTemplateMetric(m TemplateMetricModel) client.TemplateMetric {
	metric := client.TemplateMetric{
		Name:                     m.Name.ValueString(),
		DisplayName:              m.DisplayName.ValueString(),
		Description:              m.Description.ValueString(),
		Category:                 m.Category.ValueString(),
		ThresholdType:            m.ThresholdType.ValueString(),
		Availability:             m.Availability.ValueBool(),
		AlertOnMissingComponents: m.AlertOnMissingComponents.ValueBool(),
		AlertSubject:             m.AlertSubject.ValueString(),
		AlertContent:             m.AlertContent.ValueString(),
		MetricValues:             []client.TemplateMetricValue{},
		Expression:               m.Expression.ValueString(),
		ComponentFilters:         []interface{}{},
	}

	if metric.AlertSubject == "" {
		metric.AlertSubject = templateAlertSubject
	}
	if metric.AlertContent == "" {
		metric.AlertContent = templateAlertContent
	}

	if !m.MetricCategory.IsNull() && !m.MetricCategory.IsUnknown() {
		value := m.MetricCategory.ValueString()
		metric.MetricCategory = &value
	}

	for _, mv := range m.MetricValues {
		metric.MetricValues = append(metric.MetricValues, client.TemplateMetricValue{
			Name:          mv.Name.ValueString(),
			Oid:           mv.Oid.ValueString(),
			OidType:       mv.OidType.ValueString(),
			DataType:      mv.DataType.ValueString(),
			Formula:       mv.Formula.ValueString(),
			PollFrequency: mv.PollFrequency.ValueString(),
		})
	}

	if details := m.FormatMetricValueDetails; details != nil {
		formatDetails := client.TemplateFormatMetricValueDetails{
			Units: details.Units.ValueString(),
			Graph: details.Graph.ValueBool(),
		}
		if value := details.FormatMetricValue; value != nil {
			formatDetails.FormatMetricValue = client.TemplateFormatMetricValue{
				UsedBy:        value.UsedBy.ValueString(),
				EnumeratedMap: buildEnumeratedMap(value.EnumeratedMap),
				Alerts:        value.Alerts.ValueBool(),
				Graph:         value.Graph.ValueBool(),
			}
		} else {
			formatDetails.FormatMetricValue.EnumeratedMap = []client.TemplateEnumeratedMapEntry{}
		}
		metric.FormatMetricValueDetails = &formatDetails
	}

	for _, th := range m.Thresholds {
		metric.Thresholds = append(metric.Thresholds, client.TemplateThreshold{
			CriticalThreshold:   th.CriticalThreshold.ValueString(),
			CriticalRepeatCount: int(th.CriticalRepeatCount.ValueInt64()),
			WarningThreshold:    th.WarningThreshold.ValueString(),
			WarningRepeatCount:  int(th.WarningRepeatCount.ValueInt64()),
			WarningOperator:     th.WarningOperator.ValueString(),
			CriticalOperator:    th.CriticalOperator.ValueString(),
			RaiseAlert:          th.RaiseAlert.ValueBool(),
			Availability:        th.Availability.ValueBool(),
			LevelType:           th.LevelType.ValueString(),
			CompName:            th.CompName.ValueString(),
		})
	}

	return metric
}

// populateTemplateModelFromAPI maps an API response back into the Terraform model.
func populateTemplateModelFromAPI(model *TemplateModel, t *client.Template) {
	model.Id = preserveString(model.Id, t.Id)
	model.Name = types.StringValue(t.Name)
	model.Description = preserveString(model.Description, t.Description)
	model.CollectorType = preserveString(model.CollectorType, t.CollectorType)
	model.Status = preserveString(model.Status, t.Status.State)
	model.EolDate = preserveString(model.EolDate, t.Status.EolDate)
	model.TemplateFamily = preserveString(model.TemplateFamily, t.TemplateFamily)
	model.Prerequisites = preserveString(model.Prerequisites, t.Prerequisites)
	model.CreatedDate = preserveString(model.CreatedDate, t.CreatedDate)
	model.Scope = preserveString(model.Scope, t.Scope)
	model.AssignedToResources = preserveBool(model.AssignedToResources, t.AssignedToResources)
	model.SdkTag = preserveString(model.SdkTag, t.SdkTag)
	model.MonTemplateVer = preserveInt64(model.MonTemplateVer, int64(t.MonTemplateVer))

	model.ApplicationType = preserveFirstString(model.ApplicationType, t.ApplicationType)
	model.TemplateType = preserveFirstString(model.TemplateType, t.TemplateType)

	priorGroups := model.MetricGroups
	metricGroups := make([]TemplateMetricGroupModel, 0, len(t.MetricGroups))
	for i := range t.MetricGroups {
		prior := findPriorMetricGroup(priorGroups, &t.MetricGroups[i])
		metricGroups = append(metricGroups, populateTemplateMetricGroup(t.MetricGroups[i], prior))
	}
	model.MetricGroups = metricGroups
}

func findPriorMetricGroup(groups []TemplateMetricGroupModel, apiGroup *client.TemplateMetricGroup) *TemplateMetricGroupModel {
	if apiGroup.Id != "" {
		for i := range groups {
			if groups[i].Id.ValueString() == apiGroup.Id {
				return &groups[i]
			}
		}
	}

	for i := range groups {
		if groups[i].Name.ValueString() == apiGroup.Name {
			return &groups[i]
		}
	}

	return nil
}

func findAPIMetricGroup(groups []client.TemplateMetricGroup, prior *TemplateMetricGroupModel) *client.TemplateMetricGroup {
	priorID := prior.Id.ValueString()
	if priorID != "" {
		for i := range groups {
			if groups[i].Id == priorID {
				return &groups[i]
			}
		}
	}

	priorName := prior.Name.ValueString()
	for i := range groups {
		if groups[i].Name == priorName {
			return &groups[i]
		}
	}

	return nil
}

func populateTemplateMetricGroup(mg client.TemplateMetricGroup, prior *TemplateMetricGroupModel) TemplateMetricGroupModel {
	priorGroup := TemplateMetricGroupModel{}
	if prior != nil {
		priorGroup = *prior
	}

	groupModel := TemplateMetricGroupModel{
		Id:                  preserveString(priorGroup.Id, mg.Id),
		Name:                preserveString(priorGroup.Name, mg.Name),
		Frequency:           preserveInt64(priorGroup.Frequency, int64(mg.Frequency)),
		AdapterType:         preserveString(priorGroup.AdapterType, mg.AdapterType),
		ApplicationType:     preserveString(priorGroup.ApplicationType, mg.ApplicationType),
		Script:              preserveString(priorGroup.Script, mg.Script),
		IsMultiInstance:     preserveBool(priorGroup.IsMultiInstance, mg.IsMultiInstance),
		InputScopeComponent: preserveBool(priorGroup.InputScopeComponent, mg.InputScopeComponent),
		DiscoveryMetric:     preserveString(priorGroup.DiscoveryMetric, mg.DiscoveryMetric),
		ComponentProperties: types.ObjectNull(componentPropertiesAttrTypes),
		OtelMetricsConfig:   nil,
	}

	if mg.ComponentProperties != nil {
		priorProperties := TemplateComponentPropertiesModel{}
		if prior := componentPropertiesModel(priorGroup.ComponentProperties); prior != nil {
			priorProperties = *prior
		}
		groupModel.ComponentProperties = componentPropertiesObject(TemplateComponentPropertiesModel{
			Id:        apiStringOrPrior(priorProperties.Id, mg.ComponentProperties.Id),
			Name:      preserveString(priorProperties.Name, mg.ComponentProperties.Name),
			Mode:      preserveString(priorProperties.Mode, mg.ComponentProperties.Mode),
			Oid:       preserveString(priorProperties.Oid, mg.ComponentProperties.Oid),
			LookupOid: preserveString(priorProperties.LookupOid, mg.ComponentProperties.LookupOid),
			Format:    preserveString(priorProperties.Format, mg.ComponentProperties.Format),
		})
	} else if prior != nil && !prior.ComponentProperties.IsUnknown() {
		groupModel.ComponentProperties = prior.ComponentProperties
	}

	if mg.OtelMetricsConfig != nil {
		priorOtel := TemplateOtelMetricsConfigModel{}
		if priorGroup.OtelMetricsConfig != nil {
			priorOtel = *priorGroup.OtelMetricsConfig
		}
		groupModel.OtelMetricsConfig = &TemplateOtelMetricsConfigModel{
			ConfigType: preserveString(priorOtel.ConfigType, mg.OtelMetricsConfig.ConfigType),
			ReceiverConfigYaml: preserveString(
				priorOtel.ReceiverConfigYaml,
				mg.OtelMetricsConfig.ReceiverConfigYaml,
			),
			AlertDefinitionsYaml: preserveString(
				priorOtel.AlertDefinitionsYaml,
				mg.OtelMetricsConfig.AlertDefinitionsYaml,
			),
			ExternalExporter: preserveBool(
				priorOtel.ExternalExporter,
				mg.OtelMetricsConfig.ExternalExporter,
			),
			ExternalExporterName: preserveString(
				priorOtel.ExternalExporterName,
				mg.OtelMetricsConfig.ExternalExporterName,
			),
			ExternalExporterCmdLineConfig: preserveString(
				priorOtel.ExternalExporterCmdLineConfig,
				mg.OtelMetricsConfig.ExternalExporterCmdLineConfig,
			),
		}
	} else if prior != nil {
		groupModel.OtelMetricsConfig = prior.OtelMetricsConfig
	}

	if mg.InstanceFilter != nil {
		priorFilter := TemplateInstanceFilterModel{}
		if priorGroup.InstanceFilter != nil {
			priorFilter = *priorGroup.InstanceFilter
		}
		groupModel.InstanceFilter = &TemplateInstanceFilterModel{
			InstanceType:     preserveString(priorFilter.InstanceType, mg.InstanceFilter.InstanceType),
			InstanceNames:    stringsToStringValues(mg.InstanceFilter.InstanceNames),
			DiscoveryFilters: populateDiscoveryFilters(mg.InstanceFilter.DiscoveryFilters),
		}
	} else if prior != nil {
		groupModel.InstanceFilter = prior.InstanceFilter
	}

	var additionalProperties []TemplateAdditionalPropertyModel
	if prior != nil && len(prior.AdditionalProperties) > 0 {
		additionalProperties = make([]TemplateAdditionalPropertyModel, 0, len(prior.AdditionalProperties))
		for _, priorProperty := range prior.AdditionalProperties {
			apiProperty := findAPIAdditionalProperty(mg.AdditionalProperties, priorProperty)
			if apiProperty != nil {
				additionalProperties = append(additionalProperties, populateTemplateAdditionalProperty(*apiProperty, &priorProperty))
			} else {
				additionalProperties = append(additionalProperties, priorProperty)
			}
		}
	} else if len(mg.AdditionalProperties) > 0 {
		for _, ap := range mg.AdditionalProperties {
			additionalProperties = append(additionalProperties, populateTemplateAdditionalProperty(ap, nil))
		}
	}
	// Leave nil (not an empty slice) when there's nothing to report, so state matches an
	// unset (null) config rather than tripping "inconsistent result after apply".
	groupModel.AdditionalProperties = additionalProperties

	var configParams []TemplateConfigurationParameterModel
	configurationParameters := mg.ConfigurationParameters
	if prior != nil {
		if len(prior.ConfigurationParameters) > 0 {
			configParams = make([]TemplateConfigurationParameterModel, 0, len(prior.ConfigurationParameters))
			for _, priorParam := range prior.ConfigurationParameters {
				apiParam := findAPIConfigurationParameter(mg.ConfigurationParameters, priorParam)
				if apiParam != nil {
					configParams = append(configParams, populateTemplateConfigurationParameter(*apiParam, &priorParam))
				} else {
					configParams = append(configParams, priorParam)
				}
			}
		}
		groupModel.ConfigurationParameters = configParams
	} else {
		for _, cp := range configurationParameters {
			configParams = append(configParams, populateTemplateConfigurationParameter(cp, nil))
		}
		groupModel.ConfigurationParameters = configParams
	}

	var metrics []TemplateMetricModel
	if prior != nil {
		if len(prior.Metrics) > 0 {
			metrics = make([]TemplateMetricModel, 0, len(prior.Metrics))
			for _, priorMetric := range prior.Metrics {
				apiMetric := findAPIMetric(mg.Metrics, priorMetric.Name.ValueString())
				if apiMetric != nil {
					metrics = append(metrics, populateTemplateMetric(*apiMetric, &priorMetric))
				} else {
					metrics = append(metrics, priorMetric)
				}
			}
		}
	} else {
		for _, m := range mg.Metrics {
			metrics = append(metrics, populateTemplateMetric(m, nil))
		}
	}
	groupModel.Metrics = metrics

	return groupModel
}

func findAPIConfigurationParameter(
	parameters []client.TemplateConfigurationParameter,
	prior TemplateConfigurationParameterModel,
) *client.TemplateConfigurationParameter {
	for i := range parameters {
		if prior.Id.ValueInt64() != 0 && parameters[i].Id == int(prior.Id.ValueInt64()) {
			return &parameters[i]
		}
		if parameters[i].Name == prior.Name.ValueString() {
			return &parameters[i]
		}
	}

	return nil
}

func populateTemplateConfigurationParameter(
	cp client.TemplateConfigurationParameter,
	prior *TemplateConfigurationParameterModel,
) TemplateConfigurationParameterModel {
	priorValue := TemplateConfigurationParameterModel{}
	if prior != nil {
		priorValue = *prior
	}

	return TemplateConfigurationParameterModel{
		Id:                   preserveInt64(priorValue.Id, int64(cp.Id)),
		VistaraCustomParamId: preserveInt64(priorValue.VistaraCustomParamId, int64(cp.VistaraCustomParamId)),
		Name:                 preserveString(priorValue.Name, cp.Name),
		DisplayName:          preserveString(priorValue.DisplayName, cp.DisplayName),
		CollectorType:        preserveString(priorValue.CollectorType, cp.CollectorType),
		DefaultValue:         preserveString(priorValue.DefaultValue, cp.DefaultValue),
		Description:          preserveString(priorValue.Description, cp.Description),
		Configurable:         preserveBool(priorValue.Configurable, cp.Configurable),
		InputType:            preserveString(priorValue.InputType, cp.InputType),
		InputOptionValues:    preserveString(priorValue.InputOptionValues, cp.InputOptionValues),
	}
}

func findAPIMetric(metrics []client.TemplateMetric, name string) *client.TemplateMetric {
	for i := range metrics {
		if metrics[i].Name == name {
			return &metrics[i]
		}
	}

	return nil
}

func populateTemplateMetric(m client.TemplateMetric, prior *TemplateMetricModel) TemplateMetricModel {
	priorMetric := TemplateMetricModel{}
	if prior != nil {
		priorMetric = *prior
	}

	metricModel := TemplateMetricModel{
		Name:                     types.StringValue(m.Name),
		DisplayName:              apiStringOrPrior(priorMetric.DisplayName, m.DisplayName),
		Description:              apiStringOrPrior(priorMetric.Description, m.Description),
		Category:                 apiStringOrPrior(priorMetric.Category, m.Category),
		ThresholdType:            types.StringValue(m.ThresholdType),
		Availability:             preserveBool(priorMetric.Availability, m.Availability),
		AlertOnMissingComponents: preserveBool(priorMetric.AlertOnMissingComponents, m.AlertOnMissingComponents),
		AlertSubject:             alertStringOrDefault(m.AlertSubject, templateAlertSubject),
		AlertContent:             alertStringOrDefault(m.AlertContent, templateAlertContent),
		Expression:               preserveString(priorMetric.Expression, m.Expression),
		MetricCategory:           priorMetric.MetricCategory,
	}

	if m.MetricCategory != nil {
		metricModel.MetricCategory = types.StringValue(*m.MetricCategory)
	}

	var metricValues []TemplateMetricValueModel
	for i, mv := range m.MetricValues {
		priorValue := TemplateMetricValueModel{}
		if prior != nil && i < len(prior.MetricValues) {
			priorValue = prior.MetricValues[i]
		}
		metricValues = append(metricValues, TemplateMetricValueModel{
			Name:          types.StringValue(mv.Name),
			Oid:           preserveString(priorValue.Oid, mv.Oid),
			OidType:       preserveString(priorValue.OidType, mv.OidType),
			DataType:      preserveString(priorValue.DataType, mv.DataType),
			Formula:       preserveString(priorValue.Formula, mv.Formula),
			PollFrequency: preserveString(priorValue.PollFrequency, mv.PollFrequency),
		})
	}
	metricModel.MetricValues = metricValues

	if m.FormatMetricValueDetails != nil {
		priorDetails := TemplateFormatMetricValueDetailsModel{}
		if priorMetric.FormatMetricValueDetails != nil {
			priorDetails = *priorMetric.FormatMetricValueDetails
		}
		priorFormatValue := TemplateFormatMetricValueModel{}
		if priorDetails.FormatMetricValue != nil {
			priorFormatValue = *priorDetails.FormatMetricValue
		}
		metricModel.FormatMetricValueDetails = &TemplateFormatMetricValueDetailsModel{
			FormatMetricValue: &TemplateFormatMetricValueModel{
				UsedBy:        preserveString(priorFormatValue.UsedBy, m.FormatMetricValueDetails.FormatMetricValue.UsedBy),
				EnumeratedMap: populateEnumeratedMap(m.FormatMetricValueDetails.FormatMetricValue.EnumeratedMap),
				Alerts:        preserveBool(priorFormatValue.Alerts, m.FormatMetricValueDetails.FormatMetricValue.Alerts),
				Graph:         preserveBool(priorFormatValue.Graph, m.FormatMetricValueDetails.FormatMetricValue.Graph),
			},
			Units: preserveString(priorDetails.Units, m.FormatMetricValueDetails.Units),
			Graph: preserveBool(priorDetails.Graph, m.FormatMetricValueDetails.Graph),
		}
	} else if prior != nil {
		metricModel.FormatMetricValueDetails = prior.FormatMetricValueDetails
	}

	var thresholds []TemplateThresholdModel
	for _, th := range m.Thresholds {
		priorThreshold := priorThreshold(prior, len(thresholds))
		priorValue := TemplateThresholdModel{}
		if priorThreshold != nil {
			priorValue = *priorThreshold
		}
		thresholds = append(thresholds, TemplateThresholdModel{
			CriticalThreshold:   apiStringOrPrior(priorValue.CriticalThreshold, th.CriticalThreshold),
			CriticalRepeatCount: apiInt64OrPrior(priorValue.CriticalRepeatCount, int64(th.CriticalRepeatCount)),
			WarningThreshold:    apiStringOrPrior(priorValue.WarningThreshold, th.WarningThreshold),
			WarningRepeatCount:  apiInt64OrPrior(priorValue.WarningRepeatCount, int64(th.WarningRepeatCount)),
			WarningOperator:     apiStringOrPrior(priorValue.WarningOperator, th.WarningOperator),
			CriticalOperator:    apiStringOrPrior(priorValue.CriticalOperator, th.CriticalOperator),
			RaiseAlert:          preserveBool(priorValue.RaiseAlert, th.RaiseAlert),
			Availability:        preserveBool(priorValue.Availability, th.Availability),
			LevelType:           apiStringOrPrior(priorValue.LevelType, th.LevelType),
			CompName:            apiStringOrPrior(priorValue.CompName, th.CompName),
		})
	}
	metricModel.Thresholds = thresholds

	return metricModel
}

// alertStringOrDefault returns the API value verbatim (the user-established string) when
// non-empty, otherwise falls back to the template's built-in default alert text.
func alertStringOrDefault(apiValue string, defaultValue string) types.String {
	if apiValue == "" {
		return types.StringValue(defaultValue)
	}

	return types.StringValue(apiValue)
}

func stringsToStringValues(values []string) []types.String {
	var result []types.String
	for _, v := range values {
		result = append(result, types.StringValue(v))
	}

	return result
}

func populateDiscoveryFilters(filters []client.TemplateDiscoveryFilter) []TemplateDiscoveryFilterModel {
	var result []TemplateDiscoveryFilterModel
	for _, f := range filters {
		result = append(result, TemplateDiscoveryFilterModel{
			Operator: types.StringValue(f.Operator),
			Value:    types.StringValue(f.Value),
		})
	}

	return result
}

// populateEnumeratedMap converts the API's key/value entries into model form, leaving the
// result nil (not an empty slice) when there's nothing to report, so state matches an unset
// (null) config rather than tripping "inconsistent result after apply".
func populateEnumeratedMap(entries []client.TemplateEnumeratedMapEntry) []TemplateEnumeratedMapEntryModel {
	var result []TemplateEnumeratedMapEntryModel
	for _, entry := range entries {
		result = append(result, TemplateEnumeratedMapEntryModel{
			Key:   types.StringValue(entry.Key),
			Value: types.StringValue(entry.Value),
		})
	}

	return result
}

// populateEnumeratedMapAlwaysList is like populateEnumeratedMap but always returns a concrete
// (possibly empty) slice, for additional_properties' enumerated_map which examples always set
// explicitly to [] rather than leaving unset.
func populateEnumeratedMapAlwaysList(entries []client.TemplateEnumeratedMapEntry) []TemplateEnumeratedMapEntryModel {
	result := make([]TemplateEnumeratedMapEntryModel, 0, len(entries))
	for _, entry := range entries {
		result = append(result, TemplateEnumeratedMapEntryModel{
			Key:   types.StringValue(entry.Key),
			Value: types.StringValue(entry.Value),
		})
	}

	return result
}

func findAPIAdditionalProperty(
	properties []client.TemplateAdditionalProperty,
	prior TemplateAdditionalPropertyModel,
) *client.TemplateAdditionalProperty {
	for i := range properties {
		if prior.Id.ValueString() != "" && properties[i].Id == prior.Id.ValueString() {
			return &properties[i]
		}
		if properties[i].Name == prior.Name.ValueString() {
			return &properties[i]
		}
	}

	return nil
}

func populateTemplateAdditionalProperty(
	ap client.TemplateAdditionalProperty,
	prior *TemplateAdditionalPropertyModel,
) TemplateAdditionalPropertyModel {
	priorValue := TemplateAdditionalPropertyModel{}
	if prior != nil {
		priorValue = *prior
	}

	return TemplateAdditionalPropertyModel{
		Id:            apiStringOrPrior(priorValue.Id, ap.Id),
		Name:          preserveString(priorValue.Name, ap.Name),
		OidType:       preserveString(priorValue.OidType, ap.OidType),
		Oid:           preserveString(priorValue.Oid, ap.Oid),
		Operator:      apiStringOrPrior(priorValue.Operator, ap.Operator),
		Operand:       apiStringOrPrior(priorValue.Operand, ap.Operand),
		FormatValue:   apiStringOrPrior(priorValue.FormatValue, ap.FormatValue),
		EnumeratedMap: populateEnumeratedMapAlwaysList(ap.EnumeratedMap),
		Units:         apiStringOrPrior(priorValue.Units, ap.Units),
		DefaultValue:  apiStringOrPrior(priorValue.DefaultValue, ap.DefaultValue),
		Description:   apiStringOrPrior(priorValue.Description, ap.Description),
	}
}

func preserveString(current types.String, apiValue string) types.String {
	if !current.IsNull() && !current.IsUnknown() {
		return current
	}

	return types.StringValue(apiValue)
}

func apiStringOrPrior(current types.String, apiValue string) types.String {
	if apiValue != "" {
		return types.StringValue(apiValue)
	}

	if !current.IsNull() && !current.IsUnknown() {
		return current
	}

	return types.StringValue("")
}

func preserveBool(current types.Bool, apiValue bool) types.Bool {
	if !current.IsNull() && !current.IsUnknown() {
		return current
	}

	return types.BoolValue(apiValue)
}

func preserveInt64(current types.Int64, apiValue int64) types.Int64 {
	if !current.IsNull() && !current.IsUnknown() {
		return current
	}

	return types.Int64Value(apiValue)
}

func apiInt64OrPrior(current types.Int64, apiValue int64) types.Int64 {
	if apiValue != 0 {
		return types.Int64Value(apiValue)
	}

	if !current.IsNull() && !current.IsUnknown() {
		return current
	}

	return types.Int64Value(0)
}

func preserveFirstString(current types.String, apiValue []string) types.String {
	if !current.IsNull() && !current.IsUnknown() {
		return current
	}

	if len(apiValue) == 0 {
		return types.StringValue("")
	}

	return types.StringValue(apiValue[0])
}

func priorThreshold(metric *TemplateMetricModel, index int) *TemplateThresholdModel {
	if metric != nil && index < len(metric.Thresholds) {
		return &metric.Thresholds[index]
	}

	return nil
}
