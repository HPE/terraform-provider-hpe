// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure implementation satisfies the expected interfaces
var (
	_ resource.Resource                   = &ServicemapResource{}
	_ resource.ResourceWithImportState    = &ServicemapResource{}
	_ resource.ResourceWithModifyPlan     = &ServicemapResource{}
	_ resource.ResourceWithValidateConfig = &ServicemapResource{}
)

// ServicemapResource defines the resource implementation.
type ServicemapResource struct {
	BaseResource
}

// NewServicemap creates a new instance of the resource.
func NewServicemap() resource.Resource {
	return &ServicemapResource{}
}

// Metadata returns the resource type name.
func (r *ServicemapResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicemap"
}

// Recursive schema definition example
func (r *ServicemapResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			// parent attribute
			"parent": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// link attribute
			"link": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// resources at root-level
			"resources": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			// search query attribute (e.g., "name CONTAINS pro")
			"search_query": schema.StringAttribute{
				Optional: true,
			},
			// alert type configuration attribute
			"alert_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("any-critical-alert", "availability", "critical-alert"),
				},
			},
			// metrics for critical-alert type
			"metrics": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Required and non-empty for critical-alert; must match the metric_name values in metric_criteria. Must be empty or omitted for any-critical-alert.",
				Validators:          []validator.Set{setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1))},
			},
			"metric_criteria": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Required and non-empty for critical-alert; must be omitted for other alert types. Each metric must have exactly one criterion.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"metric_name": schema.StringAttribute{
							Required:   true,
							Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"operator": schema.StringAttribute{
							Required:   true,
							Validators: []validator.String{stringvalidator.OneOf("-", "=")},
						},
						"components": schema.SetAttribute{
							Optional:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "Required and non-empty when operator is \"=\". Must be omitted when operator is \"-\".",
							Validators:          []validator.Set{setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1))},
						},
					},
				},
			},
			"threshold_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("count"),
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.OneOf("count", "percentage"),
				},
			},
			"threshold_limit": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				// Usa el validador de int64 del paquete de validators
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				Default: int64default.StaticInt64(1),
			},
		},
	}
}

func translatePlanToModel(ctx context.Context, plan ServicemapModel) (client.CreateServicemap, diag.Diagnostics) {
	var resources []string
	for _, i := range plan.Resources {
		resources = append(resources, i.ValueString())
	}

	var parent *client.Parent

	if plan.Parent.ValueString() != "" {
		parent = &client.Parent{
			Id: plan.Parent.ValueString(),
		}
	}

	availabilityThreshold := client.AvailabilityThreshold{
		ThresholdType:  plan.ThresholdType.ValueString(),
		ThresholdLimit: int(plan.ThresholdLimit.ValueInt64()),
	}

	monitor, diags := buildServicemapAvailabilityMonitor(ctx, plan)

	return client.CreateServicemap{
		Name:                       plan.Name.ValueString(),
		Type:                       plan.Type.ValueString(),
		Id:                         plan.Id.ValueString(),
		Parent:                     parent,
		Resources:                  resources,
		SearchQuery:                plan.SearchQuery.ValueString(),
		ServiceAvailabilityMonitor: monitor,
		AvailabilityThreshold:      availabilityThreshold,
	}, diags
}

func buildServicemapAvailabilityMonitor(
	ctx context.Context,
	plan ServicemapModel,
) (*client.ServiceAvailabilityMonitor, diag.Diagnostics) {
	var diags diag.Diagnostics
	if plan.AlertType.IsNull() || plan.AlertType.ValueString() == "" {
		return nil, diags
	}

	monitor := &client.ServiceAvailabilityMonitor{
		AlertType: plan.AlertType.ValueString(),
		Metrics:   make([]string, 0, len(plan.Metrics)),
		MatchType: "ANY",
	}
	if monitor.AlertType != "any-critical-alert" {
		for _, metric := range plan.Metrics {
			monitor.Metrics = append(monitor.Metrics, metric.ValueString())
		}
	}
	if monitor.AlertType == "critical-alert" {
		var criteria []ServicemapMetricCriterionModel
		diags.Append(plan.MetricCriteria.ElementsAs(ctx, &criteria, false)...)
		for _, criterion := range criteria {
			entry := client.ServicemapMetricCriterion{
				MetricName: criterion.MetricName.ValueString(),
				Operator:   criterion.Operator.ValueString(),
			}
			if entry.Operator == "=" {
				diags.Append(criterion.Components.ElementsAs(ctx, &entry.Components, false)...)
			}
			monitor.MetricCriteria = append(monitor.MetricCriteria, entry)
		}
	}

	return monitor, diags
}

// mapServiceAvailabilityMonitorToModel maps alert_type/metrics/metric_criteria from a service
// availability monitor returned by the API, the inverse of buildServicemapAvailabilityMonitor.
// A nil monitor, or one without metric criteria, maps metric_criteria to a typed null set (not
// an empty one), and a criterion without components maps its components to a typed null set, so
// Terraform doesn't report a "produced an unexpected new value" plan inconsistency.
func mapServiceAvailabilityMonitorToModel(
	monitor *client.ServiceAvailabilityMonitor,
) (alertType types.String, metrics []types.String, metricCriteria types.Set, diags diag.Diagnostics) {
	metricCriteria = types.SetNull(servicemapMetricCriterionType)

	if monitor == nil {
		return types.StringNull(), nil, metricCriteria, diags
	}

	alertType = types.StringValue(monitor.AlertType)
	for _, m := range monitor.Metrics {
		metrics = append(metrics, types.StringValue(m))
	}

	if len(monitor.MetricCriteria) == 0 {
		return alertType, metrics, metricCriteria, diags
	}

	elems := make([]attr.Value, len(monitor.MetricCriteria))
	for i, criterion := range monitor.MetricCriteria {
		components := types.SetNull(types.StringType)
		if len(criterion.Components) > 0 {
			componentElems := make([]attr.Value, len(criterion.Components))
			for j, c := range criterion.Components {
				componentElems[j] = types.StringValue(c)
			}
			componentsSet, componentsDiags := types.SetValue(types.StringType, componentElems)
			diags.Append(componentsDiags...)
			components = componentsSet
		}

		obj, objDiags := types.ObjectValue(servicemapMetricCriterionType.AttrTypes, map[string]attr.Value{
			"metric_name": types.StringValue(criterion.MetricName),
			"operator":    types.StringValue(criterion.Operator),
			"components":  components,
		})
		diags.Append(objDiags...)
		elems[i] = obj
	}

	criteriaSet, criteriaDiags := types.SetValue(servicemapMetricCriterionType, elems)
	diags.Append(criteriaDiags...)
	metricCriteria = criteriaSet

	return alertType, metrics, metricCriteria, diags
}

func (r *ServicemapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServicemapModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use the JSON from the plan
	resp.Diagnostics.Append(validateServicemapMonitor(ctx, plan.AlertType, servicemapMetricsSet(plan.Metrics), plan.MetricCriteria)...)
	createServicemap, diags := translatePlanToModel(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Create the Servicemap in the backend
	newServicemap, err := r.apiClient.CreateServicemap(tenantId, createServicemap)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())

		return
	}

	// Assign the backend response directly into the state
	plan.Id = types.StringValue(newServicemap.Id)

	// Save new state back to Terraform
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// --- Read ---
func (r *ServicemapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServicemapModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	backendServicemap, err := r.apiClient.GetServicemap(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "No Service group") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}

	if backendServicemap == nil {
		resp.State.RemoveResource(ctx)

		return
	}

	// Only refresh the monitor fields when the API actually returns one: some service map node
	// types (e.g. child "Resource" nodes) omit "serviceAvailabilityMonitor" from the read
	// response even though it was accepted on create/update, so treating that omission as "the
	// monitor was removed remotely" would wipe out a still-valid configuration on every refresh.
	if backendServicemap.ServiceAvailabilityMonitor != nil {
		alertType, metrics, metricCriteria, monitorDiags := mapServiceAvailabilityMonitorToModel(
			backendServicemap.ServiceAvailabilityMonitor,
		)
		resp.Diagnostics.Append(monitorDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.AlertType = alertType
		state.Metrics = metrics
		state.MetricCriteria = metricCriteria
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update modifies an existing resource.
func (r *ServicemapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ServicemapModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Id cannot be changed; ensure Id stays the same.
	if plan.Id.ValueString() != "" && plan.Id.ValueString() != state.Id.ValueString() {
		resp.Diagnostics.AddError("Id Immutable", "Id cannot be changed after creation.")

		return
	}

	// Use the JSON from the state
	resp.Diagnostics.Append(validateServicemapMonitor(ctx, plan.AlertType, servicemapMetricsSet(plan.Metrics), plan.MetricCriteria)...)
	updateServicemap, diags := translatePlanToModel(ctx, state)
	resp.Diagnostics.Append(diags...)
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	if !plan.Name.IsNull() {
		updateServicemap.Name = plan.Name.ValueString()
	}
	if !plan.Type.IsNull() {
		updateServicemap.Type = plan.Type.ValueString()
	}
	if !reflect.DeepEqual(plan.Resources, state.Resources) {
		updateServicemap.Resources = make([]string, 0)
		for _, i := range plan.Resources {
			updateServicemap.Resources = append(updateServicemap.Resources, i.ValueString())
		}
	}
	if !plan.SearchQuery.IsNull() {
		updateServicemap.SearchQuery = plan.SearchQuery.ValueString()
	}
	if !plan.AlertType.IsNull() && plan.AlertType.ValueString() != "" {
		var monitorDiags diag.Diagnostics
		updateServicemap.ServiceAvailabilityMonitor, monitorDiags = buildServicemapAvailabilityMonitor(ctx, plan)
		resp.Diagnostics.Append(monitorDiags...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	updateServicemap.AvailabilityThreshold = client.AvailabilityThreshold{
		ThresholdType:  plan.ThresholdType.ValueString(),
		ThresholdLimit: int(plan.ThresholdLimit.ValueInt64()),
	}

	updatedServicemap, err := r.apiClient.CreateServicemap(tenantId, updateServicemap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Resource",
			fmt.Sprintf("Could not update resource: %s", err),
		)

		return
	}

	state.Id = types.StringValue(updatedServicemap.Id)
	state.Name = types.StringValue(updatedServicemap.Name)
	state.Type = types.StringValue(updatedServicemap.Type)
	state.ThresholdLimit = types.Int64Value(int64(updatedServicemap.AvailabilityThreshold.ThresholdLimit))
	state.ThresholdType = types.StringValue(updatedServicemap.AvailabilityThreshold.ThresholdType)
	state.SearchQuery = plan.SearchQuery
	state.AlertType = plan.AlertType
	state.Metrics = plan.Metrics
	state.MetricCriteria = plan.MetricCriteria

	// Save
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete removes the resource from the API.
func (r *ServicemapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServicemapModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteServicemap(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())

		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing servicemap.
// Import ID format: <servicemap_id> or <client_id>:<servicemap_id> (MSP only)
func (r *ServicemapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parsed, err := r.ParseImportID(req.ID, 1)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())

		return
	}

	tenantId := r.TenantForImport(parsed)
	id := parsed.Parts[0]

	existing, err := r.apiClient.GetServicemap(tenantId, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Servicemap",
			fmt.Sprintf("Could not import servicemap with ID '%s': %s", id, err),
		)

		return
	}

	if existing == nil {
		resp.Diagnostics.AddError(
			"Error Importing Servicemap",
			fmt.Sprintf("Servicemap with ID '%s' not found", id),
		)

		return
	}

	alertType, metrics, metricCriteria, monitorDiags := mapServiceAvailabilityMonitorToModel(existing.ServiceAvailabilityMonitor)
	resp.Diagnostics.Append(monitorDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := ServicemapModel{
		Client:         parsed.Client,
		Id:             types.StringValue(existing.Id),
		Name:           types.StringValue(existing.Name),
		Type:           types.StringValue(existing.Type),
		ThresholdType:  types.StringValue(existing.AvailabilityThreshold.ThresholdType),
		ThresholdLimit: types.Int64Value(int64(existing.AvailabilityThreshold.ThresholdLimit)),
		AlertType:      alertType,
		Metrics:        metrics,
		MetricCriteria: metricCriteria,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

var servicemapMetricCriterionType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"metric_name": types.StringType,
	"operator":    types.StringType,
	"components":  types.SetType{ElemType: types.StringType},
}}

type ServicemapModel struct {
	Client         types.String   `tfsdk:"client"`
	Id             types.String   `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	Type           types.String   `tfsdk:"type"`
	Parent         types.String   `tfsdk:"parent"`
	Link           types.String   `tfsdk:"link"`
	Resources      []types.String `tfsdk:"resources"`
	SearchQuery    types.String   `tfsdk:"search_query"`
	AlertType      types.String   `tfsdk:"alert_type"`
	Metrics        []types.String `tfsdk:"metrics"`
	MetricCriteria types.Set      `tfsdk:"metric_criteria"`
	ThresholdType  types.String   `tfsdk:"threshold_type"`
	ThresholdLimit types.Int64    `tfsdk:"threshold_limit"`
}

type ServicemapMetricCriterionModel struct {
	MetricName types.String `tfsdk:"metric_name"`
	Operator   types.String `tfsdk:"operator"`
	Components types.Set    `tfsdk:"components"`
}

func servicemapMetricsSet(metrics []types.String) types.Set {
	values := make([]attr.Value, len(metrics))
	for i, metric := range metrics {
		values[i] = metric
	}

	return types.SetValueMust(types.StringType, values)
}

func (r *ServicemapResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var alertType types.String
	var metrics, criteria types.Set
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("alert_type"), &alertType)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("metrics"), &metrics)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("metric_criteria"), &criteria)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(validateServicemapMonitor(ctx, alertType, metrics, criteria)...)
	}
}

func validateServicemapMonitor(ctx context.Context, alertType types.String, metrics, criteria types.Set) diag.Diagnostics {
	var diags diag.Diagnostics
	if alertType.IsUnknown() {
		return diags
	}
	if alertType.ValueString() != "critical-alert" {
		if !criteria.IsNull() && !criteria.IsUnknown() {
			diags.AddAttributeError(
				path.Root("metric_criteria"),
				"Unexpected Metric Criteria",
				"metric_criteria must be omitted unless alert_type is critical-alert.",
			)
		}
		if alertType.ValueString() == "any-critical-alert" && !metrics.IsUnknown() && len(metrics.Elements()) > 0 {
			diags.AddAttributeError(
				path.Root("metrics"),
				"Unexpected Metrics",
				"metrics must be empty or omitted when alert_type is any-critical-alert.",
			)
		}

		return diags
	}
	if !metrics.IsUnknown() && len(metrics.Elements()) == 0 {
		diags.AddAttributeError(path.Root("metrics"), "Missing Metrics", "critical-alert requires a non-empty metrics set.")
	}
	if criteria.IsUnknown() {
		return diags
	}
	if len(criteria.Elements()) == 0 {
		diags.AddAttributeError(
			path.Root("metric_criteria"),
			"Missing Metric Criteria",
			"critical-alert requires a non-empty metric_criteria set.",
		)

		return diags
	}
	names := make(map[string]bool)
	allNamesKnown := true
	for _, element := range criteria.Elements() {
		if element.IsUnknown() {
			allNamesKnown = false

			continue
		}
		var criterion ServicemapMetricCriterionModel
		object, ok := element.(types.Object)
		if !ok || object.IsNull() {
			diags.AddAttributeError(
				path.Root("metric_criteria"),
				"Invalid Metric Criterion",
				"Each criterion must be an object with metric_name and operator.",
			)

			continue
		}
		diags.Append(object.As(ctx, &criterion, basetypes.ObjectAsOptions{})...)
		if criterion.MetricName.IsUnknown() {
			allNamesKnown = false
		} else {
			name := criterion.MetricName.ValueString()
			if strings.TrimSpace(name) == "" || names[name] {
				diags.AddAttributeError(
					path.Root("metric_criteria"),
					"Invalid Metric Name",
					"Each criterion must have a non-empty, unique metric_name.",
				)
			}
			names[name] = true
		}
		switch criterion.Operator.ValueString() {
		case "=":
			if !criterion.Components.IsUnknown() && len(criterion.Components.Elements()) == 0 {
				diags.AddAttributeError(path.Root("metric_criteria"), "Missing Components", "operator = requires a non-empty components set.")
			}
		case "-":
			if !criterion.Components.IsNull() && !criterion.Components.IsUnknown() {
				diags.AddAttributeError(path.Root("metric_criteria"), "Unexpected Components", "components must be omitted when operator is -.")
			}
		default:
			if !criterion.Operator.IsUnknown() {
				diags.AddAttributeError(path.Root("metric_criteria"), "Invalid Operator", "operator must be - or =.")
			}
		}
	}
	if allNamesKnown && !metrics.IsUnknown() {
		metricsKnown := true
		metricNames := make(map[string]bool)
		for _, element := range metrics.Elements() {
			if element.IsUnknown() {
				metricsKnown = false

				continue
			}
			value, ok := element.(types.String)
			if ok {
				metricNames[value.ValueString()] = true
			}
		}
		if metricsKnown && !reflect.DeepEqual(names, metricNames) {
			diags.AddAttributeError(
				path.Root("metric_criteria"),
				"Mismatched Metric Names",
				"metrics and metric_criteria must contain exactly the same metric names.",
			)
		}
	}

	return diags
}

func (r *ServicemapResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Call base implementation
	r.BaseResource.ModifyPlan(ctx, req, resp)

	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state ServicemapModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	req.State.Get(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateServicemapMonitor(ctx, plan.AlertType, servicemapMetricsSet(plan.Metrics), plan.MetricCriteria)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasClientOverride := !plan.Client.IsNull() && (plan.Client.IsUnknown() || strings.TrimSpace(plan.Client.ValueString()) != "")
	if strings.ToUpper(r.apiClient.Scope) != "CLIENT" && !hasClientOverride {
		resp.Diagnostics.AddError(
			"ServiceMaps can only be created at Client level",
			"Use a client-scoped provider configuration or specify the client using unique ID.",
		)

		return
	}
}
