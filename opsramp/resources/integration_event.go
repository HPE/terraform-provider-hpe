// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &IntegrationEventResource{}
	_ resource.ResourceWithImportState = &IntegrationEventResource{}
)

// IntegrationEventResource manages outbound events on an installed integration.
type IntegrationEventResource struct {
	BaseResource
}

// IntegrationEventModel maps Terraform schema attributes to the provider model.
type IntegrationEventModel struct {
	Id                     types.String                   `tfsdk:"id"`
	IntegrationId          types.String                   `tfsdk:"integration_id"`
	Client                 types.String                   `tfsdk:"client"`
	Name                   types.String                   `tfsdk:"name"`
	Entity                 types.String                   `tfsdk:"entity"`
	EventType              types.String                   `tfsdk:"event_type"`
	UseBaseNotifier        types.Bool                     `tfsdk:"use_base_notifier"`
	Notifier               *IntegrationEventNotifierModel `tfsdk:"notifier"`
	ThirdPartyEventType    types.String                   `tfsdk:"third_party_event_type"`
	Headers                types.Map                      `tfsdk:"headers"`
	EventPayload           types.String                   `tfsdk:"event_payload"`
	EndPointURI            types.String                   `tfsdk:"endpoint_uri"`
	ResponseHeaders        types.Map                      `tfsdk:"response_headers"`
	ResourceGroupAllowed   types.Bool                     `tfsdk:"resource_group_allowed"`
	CustomAttributeAllowed types.Bool                     `tfsdk:"custom_attribute_allowed"`
	Active                 types.Bool                     `tfsdk:"active"`
	EventLevel             types.String                   `tfsdk:"event_level"`
}

// IntegrationEventNotifierModel is the per-event notifier override.
type IntegrationEventNotifierModel struct {
	Type           types.String `tfsdk:"type"`
	BaseURI        types.String `tfsdk:"base_uri"`
	AuthType       types.String `tfsdk:"auth_type"`
	UserName       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	APIKey         types.String `tfsdk:"api_key"`
	APISecret      types.String `tfsdk:"api_secret"`
	AccessTokenURL types.String `tfsdk:"access_token_url"`
	GrantType      types.String `tfsdk:"grant_type"`
	Scope          types.String `tfsdk:"scope"`
}

// NewIntegrationEvent creates a new instance of the resource.
func NewIntegrationEvent() resource.Resource {
	return &IntegrationEventResource{}
}

// Metadata returns the resource type name.
func (r *IntegrationEventResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_event"
}

// Schema defines the schema for the resource.
func (r *IntegrationEventResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an outbound event on an OpsRamp integration. Events define how OpsRamp notifies external systems when specific entity/event-type combinations occur.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the event (e.g. INTG-EVENT-...).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"integration_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the parent integration (e.g. INTG-...).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique ID of the client (sub-tenant). If not provided, uses the provider tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the event.",
			},
			"entity": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The OpsRamp entity type this event applies to.",
			},
			"event_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The event type that triggers this notification (e.g. `UPDATE`, `CREATE`, `DELETE`).",
				Validators: []validator.String{
					stringvalidator.OneOf("UPDATE", "CREATE", "DELETE"),
				},
			},
			"use_base_notifier": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "When `true`, the event uses the integration's base outbound notifier. Set to `false` and provide a `notifier` block to override.",
			},
			"notifier": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Per-event notifier override. Only used when `use_base_notifier` is `false`.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("REST_API"),
						MarkdownDescription: "The notifier type. Allowed values: `REST_API`, `SOAP_API`.",
						Validators: []validator.String{
							stringvalidator.OneOf("REST_API", "SOAP_API"),
						},
					},
					"base_uri": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The base URI of the external system.",
					},
					"auth_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Authentication type. Allowed values: `NONE`, `BASIC`, `OAUTH2`.",
						Validators: []validator.String{
							stringvalidator.OneOf("NONE", "BASIC", "OAUTH2"),
						},
					},
					"username": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Username for BASIC or OAUTH2 (PASSWORD/REFRESH_TOKEN) authentication.",
					},
					"password": schema.StringAttribute{
						Optional:            true,
						Sensitive:           true,
						MarkdownDescription: "Password for BASIC authentication.",
					},
					"api_key": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Client ID / API key for OAUTH2 authentication.",
					},
					"api_secret": schema.StringAttribute{
						Optional:            true,
						Sensitive:           true,
						MarkdownDescription: "Client secret for OAUTH2 authentication.",
					},
					"access_token_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The token endpoint URL for OAUTH2 authentication.",
					},
					"grant_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "OAuth2 grant type. Allowed values: `CLIENT_CREDENTIALS`, `PASSWORD`, `REFRESH_TOKEN`.",
						Validators: []validator.String{
							stringvalidator.OneOf("CLIENT_CREDENTIALS", "PASSWORD", "REFRESH_TOKEN"),
						},
					},
					"scope": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The OAuth2 scope.",
					},
				},
			},
			"third_party_event_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The HTTP method used to call the external endpoint (e.g. `POST`, `PUT`).",
				Validators: []validator.String{
					stringvalidator.OneOf("POST", "PUT", "PATCH", "GET", "DELETE"),
				},
			},
			"headers": schema.MapAttribute{
				Optional:            true,
				MarkdownDescription: "HTTP headers to include in the outbound request (key-value pairs).",
				ElementType:         types.StringType,
			},
			"event_payload": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The payload template to send in the outbound request body.",
			},
			"endpoint_uri": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The specific endpoint URI for this event (overrides the base URI).",
			},
			"response_headers": schema.MapAttribute{
				Optional:            true,
				MarkdownDescription: "Headers to extract from the response (key-value pairs). Allowed keys: `Status Message`, `extTicketURL`, `extTicketId`, `Error Message`.",
				ElementType:         types.StringType,
				Validators: []validator.Map{
					mapvalidator.KeysAre(stringvalidator.OneOf("Status Message", "extTicketURL", "extTicketId", "Error Message")),
				},
			},
			"resource_group_allowed": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether resource group context is allowed for this event.",
			},
			"custom_attribute_allowed": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether custom attributes are allowed in the payload.",
			},
			"active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether this event is active.",
			},
			"event_level": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The scope level of the event (e.g. `client`, `partner`).",
			},
		},
	}
}

// getTenantId determines which tenant ID to use.
func (r *IntegrationEventResource) getTenantId(clientId types.String) string {
	if !clientId.IsNull() && clientId.ValueString() != "" {
		return clientId.ValueString()
	}

	return r.apiClient.TenantId
}

// Create handles creating the event.
func (r *IntegrationEventResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IntegrationEventModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.UseBaseNotifier.ValueBool() && plan.Notifier == nil {
		resp.Diagnostics.AddError("Validation Error", "A `notifier` block is required when `use_base_notifier` is false.")

		return
	}

	tenantId := r.getTenantId(plan.Client)

	eventReq := r.modelToRequest(plan)

	created, err := r.apiClient.CreateIntegrationEvent(tenantId, plan.IntegrationId.ValueString(), eventReq)
	if err != nil {
		resp.Diagnostics.AddError("Event Create Error", fmt.Sprintf("Could not create event: %s", err.Error()))

		return
	}

	// Capture desired active state before responseToModel overwrites it with the
	// create API response (which does not reflect the activate/deactivate action).
	desiredActive := plan.Active.ValueBool()

	r.responseToModel(created, &plan)

	// Set active/inactive state via the dedicated activate/deactivate endpoint.
	if err := r.apiClient.SetIntegrationEventActive(tenantId, plan.IntegrationId.ValueString(), created.ID, desiredActive); err != nil {
		resp.Diagnostics.AddError("Event Activate Error", fmt.Sprintf("Could not set active state for event %s: %s", created.ID, err.Error()))

		return
	}
	// Restore the value we actually applied – the create response predates the action.
	plan.Active = types.BoolValue(desiredActive)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the event state.
func (r *IntegrationEventResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IntegrationEventModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	existing, err := r.apiClient.GetIntegrationEvent(tenantId, state.IntegrationId.ValueString(), state.Id.ValueString())
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "404") || strings.Contains(errStr, "not found") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Event Read Error", err.Error())

		return
	}

	r.responseToModel(existing, &state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the event.
func (r *IntegrationEventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IntegrationEventModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IntegrationEventModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.UseBaseNotifier.ValueBool() && plan.Notifier == nil {
		resp.Diagnostics.AddError("Validation Error", "A `notifier` block is required when `use_base_notifier` is false.")

		return
	}

	tenantId := r.getTenantId(plan.Client)

	// Capture desired active state before responseToModel overwrites it with the
	// update API response (which does not reflect the activate/deactivate action).
	desiredActive := plan.Active.ValueBool()

	eventReq := r.modelToRequest(plan)

	updated, err := r.apiClient.UpdateIntegrationEvent(tenantId, state.IntegrationId.ValueString(), state.Id.ValueString(), eventReq)
	if err != nil {
		resp.Diagnostics.AddError("Event Update Error", err.Error())

		return
	}

	// If active state changed, call the dedicated activate/deactivate endpoint.
	if desiredActive != state.Active.ValueBool() {
		if err := r.apiClient.SetIntegrationEventActive(tenantId, state.IntegrationId.ValueString(), state.Id.ValueString(), desiredActive); err != nil {
			resp.Diagnostics.AddError("Event Activate Error", fmt.Sprintf("Could not set active state for event %s: %s", state.Id.ValueString(), err.Error()))

			return
		}
	}

	plan.Id = state.Id
	r.responseToModel(updated, &plan)
	// Restore the value we actually applied – the update response predates the action.
	plan.Active = types.BoolValue(desiredActive)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the event.
func (r *IntegrationEventResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IntegrationEventModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	err := r.apiClient.DeleteIntegrationEvent(tenantId, state.IntegrationId.ValueString(), state.Id.ValueString())
	if err != nil {
		if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError("Event Delete Error", err.Error())

			return
		}
	}
}

// ImportState handles import.
// Import ID format: integration_id:event_id or client_id:integration_id:event_id (MSP only)
func (r *IntegrationEventResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parsed, err := r.ParseImportID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())

		return
	}

	var state IntegrationEventModel
	state.Client = parsed.Client
	state.IntegrationId = types.StringValue(parsed.Parts[0])
	state.Id = types.StringValue(parsed.Parts[1])

	state.Name = types.StringUnknown()
	state.Entity = types.StringUnknown()
	state.EventType = types.StringUnknown()
	state.EventLevel = types.StringUnknown()
	state.Headers = types.MapNull(types.StringType)
	state.ResponseHeaders = types.MapNull(types.StringType)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// modelToRequest converts the Terraform model to the API request struct.
func (r *IntegrationEventResource) modelToRequest(plan IntegrationEventModel) client.IntegrationEventRequest {
	eventReq := client.IntegrationEventRequest{
		Name:                   plan.Name.ValueString(),
		Entity:                 plan.Entity.ValueString(),
		EventType:              plan.EventType.ValueString(),
		UseBaseNotifier:        plan.UseBaseNotifier.ValueBool(),
		ThirdPartyEventType:    plan.ThirdPartyEventType.ValueString(),
		EventPayload:           plan.EventPayload.ValueString(),
		EndPointURI:            plan.EndPointURI.ValueString(),
		ResourceGroupAllowed:   plan.ResourceGroupAllowed.ValueBool(),
		CustomAttributeAllowed: plan.CustomAttributeAllowed.ValueBool(),
	}

	if plan.Notifier != nil && !plan.UseBaseNotifier.ValueBool() {
		n := plan.Notifier
		notifier := &client.NotifierRequest{
			Type:     n.Type.ValueString(),
			BaseURI:  n.BaseURI.ValueString(),
			AuthType: n.AuthType.ValueString(),
		}
		if !n.UserName.IsNull() {
			notifier.UserName = n.UserName.ValueString()
		}
		if !n.Password.IsNull() {
			notifier.Password = n.Password.ValueString()
		}
		if !n.APIKey.IsNull() {
			notifier.APIKey = n.APIKey.ValueString()
		}
		if !n.APISecret.IsNull() {
			notifier.APISecret = n.APISecret.ValueString()
		}
		if !n.AccessTokenURL.IsNull() {
			notifier.AccessTokenURL = n.AccessTokenURL.ValueString()
		}
		if !n.GrantType.IsNull() {
			notifier.GrantType = n.GrantType.ValueString()
		}
		if !n.Scope.IsNull() {
			notifier.Scope = n.Scope.ValueString()
		}
		eventReq.Notifier = notifier
	}

	if !plan.Headers.IsNull() {
		for k, v := range plan.Headers.Elements() {
			eventReq.Headers = append(eventReq.Headers, client.KeyValuePair{
				Key:   k,
				Value: v.(types.String).ValueString(),
			})
		}
	}

	if !plan.ResponseHeaders.IsNull() {
		for k, v := range plan.ResponseHeaders.Elements() {
			eventReq.ResponseHeaders = append(eventReq.ResponseHeaders, client.KeyValuePair{
				Key:   k,
				Value: v.(types.String).ValueString(),
			})
		}
	}

	return eventReq
}

// responseToModel populates the model from an API response, preserving sensitive values from state.
func (r *IntegrationEventResource) responseToModel(apiResp *client.IntegrationEventResponse, model *IntegrationEventModel) {
	model.Id = types.StringValue(apiResp.ID)
	model.Name = types.StringValue(apiResp.Name)
	model.Entity = types.StringValue(apiResp.Entity)
	model.EventType = types.StringValue(apiResp.EventType)
	model.UseBaseNotifier = types.BoolValue(apiResp.UseBaseNotifier)
	model.ThirdPartyEventType = types.StringValue(apiResp.ThirdPartyEventType)
	model.EventPayload = types.StringValue(apiResp.EventPayload)
	model.EndPointURI = types.StringValue(apiResp.EndPointURI)
	model.ResourceGroupAllowed = types.BoolValue(apiResp.ResourceGroupAllowed)
	model.CustomAttributeAllowed = types.BoolValue(apiResp.CustomAttributeAllowed)
	model.Active = types.BoolValue(apiResp.Active)
	model.EventLevel = types.StringValue(apiResp.EventLevel)

	// Rebuild headers map
	if len(apiResp.Headers) > 0 {
		elems := make(map[string]attr.Value, len(apiResp.Headers))
		for _, kv := range apiResp.Headers {
			elems[kv.Key] = types.StringValue(kv.Value)
		}
		model.Headers = types.MapValueMust(types.StringType, elems)
	} else {
		model.Headers = types.MapNull(types.StringType)
	}

	// Rebuild response headers map
	if len(apiResp.ResponseHeaders) > 0 {
		elems := make(map[string]attr.Value, len(apiResp.ResponseHeaders))
		for _, kv := range apiResp.ResponseHeaders {
			elems[kv.Key] = types.StringValue(kv.Value)
		}
		model.ResponseHeaders = types.MapValueMust(types.StringType, elems)
	} else {
		model.ResponseHeaders = types.MapNull(types.StringType)
	}

	// Notifier — the API response doesn't return sensitive credential fields.
	// Preserve the existing notifier block from the model (already populated from state/plan).
	// Only update non-sensitive computed fields.
	if apiResp.Notifier != nil && model.Notifier != nil {
		model.Notifier.Type = types.StringValue(apiResp.Notifier.Type)
		model.Notifier.BaseURI = types.StringValue(apiResp.Notifier.BaseURI)
		model.Notifier.AuthType = types.StringValue(apiResp.Notifier.AuthType)
		if apiResp.Notifier.GrantType != "" {
			model.Notifier.GrantType = types.StringValue(apiResp.Notifier.GrantType)
		}
		if apiResp.Notifier.AccessTokenURL != "" {
			model.Notifier.AccessTokenURL = types.StringValue(apiResp.Notifier.AccessTokenURL)
		}
		if apiResp.Notifier.Scope != "" {
			model.Notifier.Scope = types.StringValue(apiResp.Notifier.Scope)
		}
	}
}

func (r *IntegrationEventResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Call base implementation
	r.BaseResource.ModifyPlan(ctx, req, resp)

	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state IntegrationEventModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	req.State.Get(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.UseBaseNotifier.ValueBool() && plan.Notifier != nil {
		resp.Diagnostics.AddError("Validation Error", "Cannot specify a `notifier` block when `use_base_notifier` is true. Either set `use_base_notifier` to false or remove the `notifier` block.")

		return
	}

	if !plan.UseBaseNotifier.ValueBool() && plan.Notifier == nil {
		resp.Diagnostics.AddError("Validation Error", "A `notifier` block is required when `use_base_notifier` is false.")

		return
	}
}
