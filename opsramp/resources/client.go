// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &ClientResource{}
var _ resource.ResourceWithImportState = &ClientResource{}
var _ resource.ResourceWithModifyPlan = &ClientResource{}

// ClientResource defines the resource implementation.
type ClientResource struct {
	BaseResource
}

// ClientModel maps Terraform schema attributes to the provider model.
type ClientModel struct {
	Id          types.Int64  `tfsdk:"id"`
	UniqueId    types.String `tfsdk:"unique_id"`
	Name        types.String `tfsdk:"name"`
	Address     types.String `tfsdk:"address"`
	City        types.String `tfsdk:"city"`
	State       types.String `tfsdk:"state"`
	Country     types.String `tfsdk:"country"`
	Zip         types.String `tfsdk:"zip"`
	TimeZone    types.String `tfsdk:"time_zone"`
	PhoneNumber types.String `tfsdk:"phone_number"`
	Addons      types.Set    `tfsdk:"addons"`
	Packages    types.Set    `tfsdk:"packages"`
}

func stringSliceFromSet(ctx context.Context, set types.Set) ([]string, error) {
	if set.IsNull() || set.IsUnknown() {
		return []string{}, nil
	}

	values := make([]string, 0)
	diags := set.ElementsAs(ctx, &values, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to parse set values")
	}

	return values, nil
}

func stringSetFromSlice(values []string) (types.Set, error) {
	elements := make([]attr.Value, len(values))
	for i, value := range values {
		elements[i] = types.StringValue(value)
	}

	set, diags := types.SetValue(types.StringType, elements)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("failed to build set value")
	}

	return set, nil
}

func mapClientResponseToState(state *ClientModel, existing *client.ClientResponse) {
	state.Id = types.Int64Value(int64(existing.Id))
	state.UniqueId = types.StringValue(existing.UniqueId)
	state.Name = types.StringValue(existing.Name)
	if existing.Address != "" {
		state.Address = types.StringValue(existing.Address)
	} else {
		state.Address = types.StringNull()
	}
	if existing.City != "" {
		state.City = types.StringValue(existing.City)
	} else {
		state.City = types.StringNull()
	}
	if existing.State != "" {
		state.State = types.StringValue(existing.State)
	} else {
		state.State = types.StringNull()
	}
	if existing.Country != "" {
		state.Country = types.StringValue(existing.Country)
	} else {
		state.Country = types.StringNull()
	}
	if existing.Zip != "" {
		state.Zip = types.StringValue(existing.Zip)
	} else {
		state.Zip = types.StringNull()
	}
	if existing.TimeZone != "" {
		state.TimeZone = types.StringValue(existing.TimeZone)
	} else {
		state.TimeZone = types.StringNull()
	}
	if existing.PhoneNumber != "" {
		state.PhoneNumber = types.StringValue(existing.PhoneNumber)
	} else {
		state.PhoneNumber = types.StringNull()
	}
}

// NewClient creates a new instance of the resource.
func NewClient() resource.Resource {
	return &ClientResource{}
}

// Metadata returns the resource type name.
func (r *ClientResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

// Schema defines the schema for the resource.
func (r *ClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Client (sub-tenant) resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the client.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"unique_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A unique identifier for the client.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the client.",
			},
			"address": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The address of the client.",
			},
			"city": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The city of the client.",
			},
			"state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The state of the client.",
			},
			"country": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The country of the client.",
			},
			"zip": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The zip code of the client.",
			},
			"time_zone": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The time zone of the client.",
			},
			"phone_number": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The phone number of the client.",
			},
			"addons": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of addons enabled in the Client.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"Adapter Integrations",
							"Extended Data Retention",
							"Knowledgebase Management",
							"OS Service Start/Stop Actions",
							"Offline Alerts",
							"Process Automation",
							"Remote Access Management",
							"SMS and Voice",
							"Alert Problem Area",
						),
					),
				},
				Default: setdefault.StaticValue(
					types.SetValueMust(types.StringType, []attr.Value{}),
				),
			},
			"packages": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of packages enabled in the Client. (e.g. `Hybrid Discovery and Monitoring`, `Event and Incident Management`, `Remediation and Automation`.)",
				Validators: []validator.Set{
					SetMustContainAndAllow(
						[]string{ // mandatory
							"Hybrid Discovery and Monitoring",
						},
						[]string{ // optional
							"Hybrid Discovery and Monitoring",
							"Event and Incident Management",
							"Remediation and Automation",
						},
					),
				},

				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("Hybrid Discovery and Monitoring"),
						},
					),
				),
			},
		},
	}
}

// Create handles the creation of the resource.
func (r *ClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClientModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addons, err := stringSliceFromSet(ctx, plan.Addons)
	if err != nil {
		resp.Diagnostics.AddError("Invalid addons value", err.Error())
		return
	}

	packages, err := stringSliceFromSet(ctx, plan.Packages)
	if err != nil {
		resp.Diagnostics.AddError("Invalid packages value", err.Error())
		return
	}

	createClient := client.CreateClient{
		Name:        plan.Name.ValueString(),
		UniqueId:    plan.UniqueId.ValueString(),
		Address:     plan.Address.ValueString(),
		City:        plan.City.ValueString(),
		State:       plan.State.ValueString(),
		Country:     plan.Country.ValueString(),
		Zip:         plan.Zip.ValueString(),
		TimeZone:    plan.TimeZone.ValueString(),
		PhoneNumber: plan.PhoneNumber.ValueString(),
		Addons:      addons,
		Packages:    packages,
	}

	created, err := r.apiClient.CreateClient(createClient)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	createdState, err := r.apiClient.GetClient(created.UniqueId)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", fmt.Sprintf("could not refresh created client %s: %s", created.UniqueId, err))
		return
	}

	mapClientResponseToState(&plan, createdState)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *ClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClientModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	existing, err := r.apiClient.GetClient(state.UniqueId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	mapClientResponseToState(&state, existing)

	// Only update optional fields if API returns non-empty values
	if existing.City != "" {
		state.City = types.StringValue(existing.City)
	}
	if existing.State != "" {
		state.State = types.StringValue(existing.State)
	}
	if existing.Zip != "" {
		state.Zip = types.StringValue(existing.Zip)
	}
	if existing.PhoneNumber != "" {
		state.PhoneNumber = types.StringValue(existing.PhoneNumber)
	}

	// Only overwrite addons/packages when the API returns values or the prior
	// state was non-null; if the user omitted the attribute and the API echoes
	// nothing back, keep null so we don't produce a null → [] drift on the
	// next plan.
	if len(existing.Addons) > 0 || !state.Addons.IsNull() {
		state.Addons, err = stringSetFromSlice(existing.Addons)
		if err != nil {
			resp.Diagnostics.AddError("Read Error", err.Error())
			return
		}
	}
	if len(existing.Packages) > 0 || !state.Packages.IsNull() {
		state.Packages, err = stringSetFromSlice(existing.Packages)
		if err != nil {
			resp.Diagnostics.AddError("Read Error", err.Error())
			return
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *ClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ClientModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addons, err := stringSliceFromSet(ctx, plan.Addons)
	if err != nil {
		resp.Diagnostics.AddError("Invalid addons value", err.Error())
		return
	}

	packages, err := stringSliceFromSet(ctx, plan.Packages)
	if err != nil {
		resp.Diagnostics.AddError("Invalid packages value", err.Error())
		return
	}

	updateClient := client.CreateClient{
		Name:        plan.Name.ValueString(),
		Address:     plan.Address.ValueString(),
		City:        plan.City.ValueString(),
		State:       plan.State.ValueString(),
		Country:     plan.Country.ValueString(),
		Zip:         plan.Zip.ValueString(),
		TimeZone:    plan.TimeZone.ValueString(),
		PhoneNumber: plan.PhoneNumber.ValueString(),
		Addons:      addons,
		Packages:    packages,
	}

	updated, err := r.apiClient.UpdateClient(state.UniqueId.ValueString(), updateClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Client",
			fmt.Sprintf("Could not update client: %s", err),
		)
		return
	}

	updatedState, err := r.apiClient.GetClient(updated.UniqueId)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Client", fmt.Sprintf("could not refresh updated client %s: %s", updated.UniqueId, err))
		return
	}

	mapClientResponseToState(&plan, updatedState)
	plan.Addons = state.Addons
	plan.Packages = state.Packages

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *ClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClientModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.apiClient.DeleteClient(state.UniqueId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
func (r *ClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importId := req.ID

	existing, err := r.apiClient.GetClient(importId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Client",
			fmt.Sprintf("Could not import client with ID '%s': %s", importId, err),
		)
		return
	}

	state := ClientModel{}
	mapClientResponseToState(&state, existing)

	// Only set optional fields if API returns non-empty values
	if existing.City != "" {
		state.City = types.StringValue(existing.City)
	}
	if existing.State != "" {
		state.State = types.StringValue(existing.State)
	}
	if existing.Zip != "" {
		state.Zip = types.StringValue(existing.Zip)
	}
	if existing.PhoneNumber != "" {
		state.PhoneNumber = types.StringValue(existing.PhoneNumber)
	}

	state.Addons, err = stringSetFromSlice(existing.Addons)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Client",
			fmt.Sprintf("Could not map addons for client with ID '%s': %s", importId, err),
		)
		return
	}

	state.Packages, err = stringSetFromSlice(existing.Packages)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Client",
			fmt.Sprintf("Could not map packages for client with ID '%s': %s", importId, err),
		)
		return
	}

	resp.State.Set(ctx, &state)
}

func (r *ClientResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Call base implementation
	r.BaseResource.ModifyPlan(ctx, req, resp)

	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state ClientModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	req.State.Get(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	if strings.ToUpper(r.apiClient.Scope) != "MSP" {
		resp.Diagnostics.AddError("Clients can only be created at MSP level", "Use an msp-scoped provider configuration.")
		return
	}
}
