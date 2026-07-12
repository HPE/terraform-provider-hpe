// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &SiteResource{}
var _ resource.ResourceWithImportState = &SiteResource{}
var _ resource.ResourceWithModifyPlan = &SiteResource{}

type SiteResource struct {
	BaseResource
}

type SiteModel struct {
	Client           types.String `tfsdk:"client"`
	Id               types.Int64  `tfsdk:"id"`
	Uuid             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	ParentId         types.Int64  `tfsdk:"parent_id"`
	Description      types.String `tfsdk:"description"`
	Address          types.String `tfsdk:"address"`
	State            types.String `tfsdk:"state"`
	City             types.String `tfsdk:"city"`
	Country          types.String `tfsdk:"country"`
	Zip              types.String `tfsdk:"zip"`
	PrimaryContactId types.String `tfsdk:"primary_contact_id"`
	PhoneNumber      types.String `tfsdk:"phone_number"`
	PhoneExtension   types.String `tfsdk:"phone_extension"`
	SearchQuery      types.String `tfsdk:"search_query"`
	Resources        types.Set    `tfsdk:"resources"`
}

func NewSite() resource.Resource {
	return &SiteResource{}
}

func (r *SiteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *SiteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Site resource for grouping and organizing monitored resources by physical or logical location.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this site should be created. If not specified, uses the provider's tenant.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The numeric identifier of the site.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the site.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{Required: true, MarkdownDescription: "The name of the site."},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The description of the site.",
			},
			"address": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The physical address of the site.",
			},
			"city": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The city where the site is located.",
			},
			"state": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The state or region where the site is located.",
			},
			"country": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The country where the site is located.",
			},
			"zip": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The postal code for the site.",
			},
			"phone_number": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The primary phone number for the site.",
			},
			"phone_extension": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The phone extension for the site.",
			},
			"parent_id": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The parent site ID, if this site is nested under another.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"primary_contact_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The primary contact person for the site.",
			},
			"search_query": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The search query to filter resources for this site.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"resources": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Set of resource IDs to assign to this site.",
			},
		},
	}
}

func translatePlanToSite(plan SiteModel) client.Site {
	site := client.Site{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		Address:        plan.Address.ValueString(),
		City:           plan.City.ValueString(),
		State:          plan.State.ValueString(),
		Country:        plan.Country.ValueString(),
		Zip:            plan.Zip.ValueString(),
		PhoneNumber:    plan.PhoneNumber.ValueString(),
		PhoneExtension: plan.PhoneExtension.ValueString(),
	}

	if !plan.ParentId.IsNull() && !plan.ParentId.IsUnknown() {
		site.Parent = &client.SiteParentRef{Id: plan.ParentId.ValueInt64()}
	}

	if !plan.PrimaryContactId.IsNull() && !plan.PrimaryContactId.IsUnknown() {
		site.PrimaryContact = &client.SiteContact{Id: plan.PrimaryContactId.ValueString()}
	}

	if !plan.SearchQuery.IsNull() && !plan.SearchQuery.IsUnknown() && plan.SearchQuery.ValueString() != "" {
		site.FilterCriteria = &client.SiteFilter{SearchQuery: plan.SearchQuery.ValueString()}
	}

	return site
}

func mapSiteResponseToModel(api *client.Site, model *SiteModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.Id = types.Int64Value(api.Id)
	model.Uuid = types.StringValue(api.Uuid)
	model.Name = types.StringValue(api.Name)
	model.Description = types.StringValue(api.Description)
	model.Address = types.StringValue(api.Address)
	model.City = types.StringValue(api.City)
	model.State = types.StringValue(api.State)
	model.Country = types.StringValue(api.Country)
	model.Zip = types.StringValue(api.Zip)
	model.PhoneNumber = types.StringValue(api.PhoneNumber)
	model.PhoneExtension = types.StringValue(api.PhoneExtension)

	if api.Parent != nil {
		model.ParentId = types.Int64Value(api.Parent.Id)
	} else {
		model.ParentId = types.Int64Null()
	}

	if api.PrimaryContact != nil && api.PrimaryContact.Id != "" {
		model.PrimaryContactId = types.StringValue(api.PrimaryContact.Id)
	} else {
		model.PrimaryContactId = types.StringNull()
	}

	model.SearchQuery = types.StringValue("")
	if api.FilterCriteria != nil && api.FilterCriteria.SearchQuery != "" {
		model.SearchQuery = types.StringValue(api.FilterCriteria.SearchQuery)
	}

	model.Resources = types.SetNull(types.StringType)

	return diags
}

func (r *SiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SiteModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	site := translatePlanToSite(plan)

	// Save planned resources before mapSiteResponseToModel overwrites them with null.
	resources := setToStringSlice(plan.Resources)

	created, err := r.apiClient.CreateSite(tenantId, site)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	diags = mapSiteResponseToModel(created, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(resources) > 0 {
		err = r.apiClient.AddSiteChilds(tenantId, created.Uuid, resources)
		if err != nil {
			resp.Diagnostics.AddError("Error adding resources to site", err.Error())
			return
		}
	}

	planResources := make([]attr.Value, len(resources))
	for i, resId := range resources {
		planResources[i] = types.StringValue(resId)
	}

	plan.Resources = types.SetValueMust(types.StringType, planResources)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *SiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SiteModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetSite(tenantId, state.Uuid.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "No site found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	// Preserve resources before mapSiteResponseToModel overwrites them with null.
	// The API response includes all resources (explicit + query-matched) so we
	// cannot use it to reconstruct the explicitly-managed set.
	savedResources := state.Resources

	diags = mapSiteResponseToModel(existing, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Resources = savedResources

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func siteResourcesSet(ctx context.Context, ids []string) (types.Set, diag.Diagnostics) {
	if ids == nil {
		ids = []string{}
	}

	return types.SetValueFrom(ctx, types.StringType, ids)
}

func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SiteModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	// Save resource sets before mapSiteResponseToModel overwrites them with null.
	oldIds := setToStringSlice(state.Resources)
	newIds := setToStringSlice(plan.Resources)

	updated, err := r.apiClient.UpdateSite(tenantId, state.Uuid.ValueString(), translatePlanToSite(plan))
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	diags = mapSiteResponseToModel(updated, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Client = plan.Client

	// Unset resources is equivalent to []: always reconcile.
	toAdd := stringSetDiff(newIds, oldIds)
	toRemove := stringSetDiff(oldIds, newIds)

	if len(toAdd) > 0 {
		if err := r.apiClient.AddSiteChilds(tenantId, updated.Uuid, toAdd); err != nil {
			resp.Diagnostics.AddError("Error adding resources to site", err.Error())
			return
		}
	}
	if len(toRemove) > 0 {
		if err := r.apiClient.RemoveSiteChilds(tenantId, updated.Uuid, toRemove); err != nil {
			resp.Diagnostics.AddError("Error removing resources from site", err.Error())
			return
		}
	}

	if !plan.Resources.IsNull() && !plan.Resources.IsUnknown() {
		state.Resources = plan.Resources
	} else {
		state.Resources, diags = siteResourcesSet(ctx, newIds)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *SiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SiteModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	if err := r.apiClient.DeleteSite(tenantId, state.Uuid.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing site.
// Import ID format: <site_id> or <client_id>:<site_id> (MSP only)
func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parsed, err := r.ParseImportID(req.ID, 1)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())
		return
	}

	tenantId := r.TenantForImport(parsed)
	siteId := parsed.Parts[0]
	var diags diag.Diagnostics

	existing, err := r.apiClient.GetSite(tenantId, siteId)
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Site", fmt.Sprintf("Could not import site with ID '%s': %s", siteId, err))
		return
	}

	state := SiteModel{Client: parsed.Client}
	diags = mapSiteResponseToModel(existing, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *SiteResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	r.BaseResource.ModifyPlan(ctx, req, resp)

	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state SiteModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	req.State.Get(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	// Removing the resources attribute from config is equivalent to resources = [].
	if plan.Resources.IsUnknown() {
		plan.Resources, _ = siteResourcesSet(ctx, nil)
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	hasClientOverride := !plan.Client.IsNull() && (plan.Client.IsUnknown() || strings.TrimSpace(plan.Client.ValueString()) != "")
	if strings.ToUpper(r.apiClient.Scope) != "CLIENT" && !hasClientOverride {
		resp.Diagnostics.AddError("Sites can only be created at Client level", "Use a client-scoped provider configuration or specify the client using unique ID.")
	}
}
