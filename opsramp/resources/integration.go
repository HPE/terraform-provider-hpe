// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

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
	_ resource.Resource                = &IntegrationResource{}
	_ resource.ResourceWithImportState = &IntegrationResource{}
	_ resource.ResourceWithModifyPlan  = &IntegrationResource{}
)

// IntegrationResource defines the resource implementation.
type IntegrationResource struct {
	BaseResource
}

// IntegrationModel maps Terraform schema attributes to the provider model.
type IntegrationModel struct {
	Id                           types.String `tfsdk:"id"`
	DisplayName                  types.String `tfsdk:"display_name"`
	Description                  types.String `tfsdk:"description"`
	Application                  types.String `tfsdk:"application"`
	Category                     types.String `tfsdk:"category"`
	Client                       types.String `tfsdk:"client"`
	Status                       types.String `tfsdk:"status"`
	BypassResourceReconciliation types.Bool   `tfsdk:"bypass_resource_reconciliation"`

	// Alert source for event integrations (CUSTOM-EVENT)
	AlertSourceID types.Int64 `tfsdk:"alert_source_id"`

	// Inbound configuration
	Inbound *IntegrationInboundModel `tfsdk:"inbound"`

	// Outbound configuration
	Outbound *IntegrationOutboundModel `tfsdk:"outbound"`

	ProfileId types.String `tfsdk:"profile_id"`
}

// IntegrationInboundModel represents the inbound configuration block
type IntegrationInboundModel struct {
	AuthType   types.String `tfsdk:"auth_type"`
	RoleId     types.String `tfsdk:"role_id"`
	Token      types.String `tfsdk:"token"`
	WebhookURL types.String `tfsdk:"webhook_url"`

	// Mapping attributes
	MapAttributes []IntegrationMapAttributes `tfsdk:"map_attributes"`

	// Drop alerts configuration
	EnableDropAlerts types.Bool `tfsdk:"enable_drop_alerts"`

	// Process definition assignments
	ProcessDefinitionIds types.Set `tfsdk:"process_definition_ids"`

	// Webhook handshake properties (JSON-encoded map)
	WebhookHandshake types.String `tfsdk:"webhook_handshake"`

	// Additional properties (CUSTOM application type only)
	AdditionalProperties types.Map `tfsdk:"additional_properties"`
}

type IntegrationMapAttributes struct {
	ThirdPartyAttribute types.String                  `tfsdk:"third_party_attribute"`
	OpsRampAttribute    types.String                  `tfsdk:"opsramp_attribute"`
	EntityType          types.String                  `tfsdk:"entity_type"`
	AttributeValues     types.Map                     `tfsdk:"attribute_values"`
	DefaultParsingValue types.String                  `tfsdk:"default_parsing_value"`
	ParsingOperators    []IntegrationParsingOperators `tfsdk:"parsing_operators"`
}

type IntegrationParsingOperators struct {
	Operator  string `tfsdk:"operator"`
	StartWord string `tfsdk:"start_word"`
	EndWord   string `tfsdk:"end_word"`
	RegexStr  string `tfsdk:"regex_str"`
}

// IntegrationOutboundModel represents the outbound configuration block
type IntegrationOutboundModel struct {
	Type     types.String `tfsdk:"type"`
	BaseURI  types.String `tfsdk:"base_uri"`
	AuthType types.String `tfsdk:"auth_type"`

	// OAuth2 / Basic auth fields
	UserName       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	APIKey         types.String `tfsdk:"api_key"`
	APISecret      types.String `tfsdk:"api_secret"`
	AccessTokenURL types.String `tfsdk:"access_token_url"`
	GrantType      types.String `tfsdk:"grant_type"`
	Scope          types.String `tfsdk:"scope"`

	// Outbound attribute mappings (read from API, not set directly)
	MapAttributes []IntegrationMapAttributes `tfsdk:"map_attributes"`
}

// NewIntegration creates a new instance of the resource.
func NewIntegration() resource.Resource {
	return &IntegrationResource{}
}

// Metadata returns the resource type name.
func (r *IntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

// Schema defines the schema for the resource.
func (r *IntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Integration. Supports event-based integrations (inbound), custom integrations (inbound + outbound), and configuration-based integrations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the installed integration (e.g. INTG-...).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The integration application type (e.g. CUSTOM, CUSTOM-EVENT, NEWRELIC, VMWARE, PROMETHEUSREMOTEWRITE).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The display name of the integration.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "A description of the integration.",
			},
			"category": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The integration category. Applicable for `CUSTOM` and `CUSTOM-EVENT` application types. For `CUSTOM-EVENT`, this is automatically set to `Monitoring`. Allowed values: `Custom`, `Collaboration`, `Monitoring`, `SSO`, `Automation`, `ADAPTER_INTEGRATION`.",
				Validators: []validator.String{
					stringvalidator.OneOf("Custom", "Collaboration", "Monitoring", "SSO", "Automation", "ADAPTER_INTEGRATION"),
				},
			},
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique ID of the client (sub-tenant) where the integration should be installed. If not provided, the integration is created at the provider tenant level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the integration (e.g. enabled, disabled).",
			},
			"alert_source_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The alert source ID for CUSTOM-EVENT integrations. Retrieve available IDs using the OpsRamp API: GET /api/v2/tenants/{tenantId}/cfg/alertSource/available/custIntg/CUSTOM-EVENT.",
			},
			"bypass_resource_reconciliation": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to bypass resource reconciliation for this integration.",
			},
			"inbound": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Inbound integration configuration. Used for event-based integrations that receive data via webhook.",
				Attributes: map[string]schema.Attribute{
					"auth_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The authentication type for inbound data (`WEBHOOK`, `OAUTH2`).",
						Validators: []validator.String{
							stringvalidator.OneOf("WEBHOOK", "OAUTH2"),
						},
					},
					"role_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The unique ID of the role to assign for inbound authentication. Must be available for the integration. If not specified, defaults to the first available role.",
					},
					"token": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "The webhook token generated for this integration.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"webhook_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The webhook URL for sending events to this integration.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"map_attributes": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Attribute mapping rules for inbound data transformation.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"opsramp_attribute": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The OpsRamp attribute to map to (e.g. alert.alertTime, alert.component). Use the opsramp_integration_inbound_properties data source to look up valid values.",
								},
								"third_party_attribute": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The third-party entity type name.",
								},
								"entity_type": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString("ALERT"),
									MarkdownDescription: "The OpsRamp entity type this mapping applies to (e.g. `ALERT`, `INCIDENT`, `SERVICEREQUEST`, `PROBLEM`, `CHANGE`, `TASK`). Defaults to `ALERT`.",
								},
								"attribute_values": schema.MapAttribute{
									Optional:            true,
									MarkdownDescription: "Specific attribute value mappings from third-party to OpsRamp values.",
									ElementType:         types.StringType,
								},
								"default_parsing_value": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString(""),
									MarkdownDescription: "If parsing operators are defined, this value is used when no operators match. Required when parsing_operators is set.",
								},
								"parsing_operators": schema.ListNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Attribute mapping rules for inbound data transformation.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"operator": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Operator type (e.g. `BEFORE`, `AFTER`, `BETWEEN`, `MATCHES`).",
												Validators: []validator.String{
													stringvalidator.OneOf("BEFORE", "AFTER", "BETWEEN", "MATCHES"),
												},
											},
											"start_word": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												MarkdownDescription: "Capture value after this word (for BEFORE/AFTER) or between this and end_word (for BETWEEN).",
												Default:             stringdefault.StaticString(""),
											},
											"end_word": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												MarkdownDescription: "Capture value before this word (for BEFORE/AFTER) or between this and end_word (for BETWEEN).",
												Default:             stringdefault.StaticString(""),
											},
											"regex_str": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												MarkdownDescription: "Capture value using a regular expresion.",
												Default:             stringdefault.StaticString(""),
											},
										},
									},
								},
							},
						},
					},
					"enable_drop_alerts": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Whether to enable dropping of duplicate/unwanted alerts.",
					},
					"process_definition_ids": schema.SetAttribute{
						Optional:            true,
						MarkdownDescription: "List of process definition IDs to assign to this integration.",
						ElementType:         types.StringType,
					},
					"webhook_handshake": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "JSON-encoded webhook handshake properties. Set when auth_type is WEBHOOK to configure handshake parameters (e.g. `{\"header_name\":\"value\"}`).",
					},
					"additional_properties": schema.MapAttribute{
						Optional:            true,
						MarkdownDescription: "Additional key-value properties for the integration.",
						ElementType:         types.StringType,
					},
				},
			},
			"outbound": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Outbound integration configuration. Used for integrations that push data to external systems.",
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
						MarkdownDescription: "The base URI of the external system to send notifications to.",
					},
					"auth_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The authentication type for outbound calls. Allowed values: `NONE`, `BASIC`, `OAUTH2`.",
						Validators: []validator.String{
							stringvalidator.OneOf("NONE", "BASIC", "OAUTH2"),
						},
					},
					"username": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Username for BASIC or OAUTH2 (PASSWORD/REFRESH_TOKEN grant) authentication.",
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
						MarkdownDescription: "Client secret / API secret for OAUTH2 authentication.",
					},
					"access_token_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The token endpoint URL for OAUTH2 authentication.",
					},
					"grant_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The OAuth2 grant type. Allowed values: `CLIENT_CREDENTIALS`, `PASSWORD`, `REFRESH_TOKEN`.",
						Validators: []validator.String{
							stringvalidator.OneOf("CLIENT_CREDENTIALS", "PASSWORD", "REFRESH_TOKEN"),
						},
					},
					"scope": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The OAuth2 scope.",
					},
					"map_attributes": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Attribute mapping rules for inbound data transformation.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"opsramp_attribute": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The OpsRamp attribute to map to (e.g. alert.alertTime, alert.component). Use the opsramp_integration_inbound_properties data source to look up valid values.",
								},
								"third_party_attribute": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The third-party entity type name.",
								},
								"entity_type": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString("ALERT"),
									MarkdownDescription: "The OpsRamp entity type this mapping applies to (e.g. `ALERT`, `INCIDENT`, `SERVICEREQUEST`, `PROBLEM`, `CHANGE`, `TASK`). Defaults to `ALERT`.",
								},
								"attribute_values": schema.MapAttribute{
									Optional:            true,
									MarkdownDescription: "Specific attribute value mappings from third-party to OpsRamp values.",
									ElementType:         types.StringType,
								},
								"default_parsing_value": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString(""),
									MarkdownDescription: "If parsing operators are defined, this value is used when no operators match. Required when parsing_operators is set.",
								},
								"parsing_operators": schema.ListNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Attribute mapping rules for inbound data transformation.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"operator": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Operator type (e.g. `BEFORE`, `AFTER`, `BETWEEN`, `MATCHES`).",
												Validators: []validator.String{
													stringvalidator.OneOf("BEFORE", "AFTER", "BETWEEN", "MATCHES"),
												},
											},
											"start_word": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												MarkdownDescription: "Capture value after this word (for BEFORE/AFTER) or between this and end_word (for BETWEEN).",
												Default:             stringdefault.StaticString(""),
											},
											"end_word": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												MarkdownDescription: "Capture value before this word (for BEFORE/AFTER) or between this and end_word (for BETWEEN).",
												Default:             stringdefault.StaticString(""),
											},
											"regex_str": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												MarkdownDescription: "Capture value using a regular expresion.",
												Default:             stringdefault.StaticString(""),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"profile_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The gateway/profile UUID to associate with this app installation.",
				PlanModifiers: []planmodifier.String{
					// UseStateForUnknown must come before RequiresReplace so that a null
					// value persisted from a previous apply is carried forward as-is
					// rather than being re-marked as (known after apply) by the framework,
					// which would incorrectly trigger a replacement on every plan.
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// getTenantId determines which tenant ID to use based on the optional client parameter
func (r *IntegrationResource) getTenantId(clientId types.String) string {
	if !clientId.IsNull() && clientId.ValueString() != "" {
		return clientId.ValueString()
	}

	return r.apiClient.TenantId
}

// Create handles the creation of the resource.
func (r *IntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IntegrationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)
	application := plan.Application.ValueString()

	var profile client.InstallV3Profile
	if !plan.ProfileId.IsNull() && plan.ProfileId.ValueString() != "" {
		profile = client.InstallV3Profile{UuId: plan.ProfileId.ValueString()}
	} else {
		profile = client.InstallV3Profile{UuId: ""}
	}

	// Build the install request (v2 API)
	installReq := client.InstallIntegrationRequest{
		DisplayName:               plan.DisplayName.ValueString(),
		Description:               plan.Description.ValueString(),
		Category:                  plan.Category.ValueString(),
		Profile:                   &profile,
		MultiAppsDiscoveryEnabled: plan.BypassResourceReconciliation.ValueBool(),
	}

	// For CUSTOM-EVENT integrations, include alertSource
	if !plan.AlertSourceID.IsNull() {
		installReq.AlertSource = &client.AlertSource{
			ID: int(plan.AlertSourceID.ValueInt64()),
		}
	}

	// Install the integration
	installed, err := r.apiClient.InstallIntegration(tenantId, application, installReq)
	if err != nil {
		resp.Diagnostics.AddError("Integration Install Error",
			fmt.Sprintf("Could not install integration of type '%s': %s", application, err.Error()))

		return
	}

	// Save state immediately after install so that if subsequent configuration fails,
	// Terraform knows the resource exists and will Update rather than re-Create.
	plan.Id = types.StringValue(installed.ID)
	plan.Status = types.StringValue(installed.Status)
	plan.DisplayName = types.StringValue(installed.DisplayName)
	if installed.Category != "" {
		plan.Category = types.StringValue(installed.Category)
	}

	// The v2 API does not return profile_id. Resolve it to a known value so that
	// Terraform does not report 'unknown value after apply'.
	if plan.ProfileId.IsUnknown() {
		plan.ProfileId = types.StringNull()
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Configure inbound if specified
	if plan.Inbound != nil {
		if err := r.configureInbound(tenantId, installed.ID, application, plan.Inbound, installed); err != nil {
			resp.Diagnostics.AddError("Inbound Configuration Error", err.Error())
			resp.State.Set(ctx, &plan)

			return
		}

		// If integration returned inbound config directly (e.g. NEWRELIC), use those values
		if installed.InboundConfig != nil && installed.InboundConfig.Authentication != nil {
			auth := installed.InboundConfig.Authentication
			if auth.Token != "" {
				plan.Inbound.Token = types.StringValue(auth.Token)
			}
			if auth.WebhookURL != "" {
				plan.Inbound.WebhookURL = types.StringValue(auth.WebhookURL)
			}
		}
	}

	// Configure outbound if specified
	if plan.Outbound != nil {
		if err := r.configureOutbound(tenantId, installed.ID, plan.Outbound); err != nil {
			resp.Diagnostics.AddError("Outbound Configuration Error", err.Error())
			resp.State.Set(ctx, &plan)

			return
		}
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// configureInbound sets up the inbound configuration for an integration.
// applicationName is the integration type (e.g. "CUSTOM", "CUSTOM-EVENT", "NEWRELIC").
// installed may be nil (e.g. during Update) and is only used to seed token/webhook from the install response.
func (r *IntegrationResource) configureInbound(
	tenantId, integrationId, applicationName string,
	inbound *IntegrationInboundModel,
	installed *client.IntegrationResponse,
) error {
	if installed != nil && installed.InboundConfig != nil && installed.InboundConfig.Authentication != nil {
		// Use values from install response (e.g. NEWRELIC auto-provisions auth)
		auth := installed.InboundConfig.Authentication
		if auth.Token != "" {
			inbound.Token = types.StringValue(auth.Token)
		}
		if auth.WebhookURL != "" {
			inbound.WebhookURL = types.StringValue(auth.WebhookURL)
		}
	}

	if inbound.AuthType.ValueString() != "" {
		var selectedRole *client.RoleClientRef

		if applicationName == "CUSTOM-EVENT" || applicationName == "CUSTOM" {

			// Retrieve Available Role
			integrationAvailableRoles, err := r.apiClient.GetIntegrationAvailableRoles(tenantId, integrationId)
			if err != nil {
				return fmt.Errorf("failed to retrieve available integration roles: %w", err)
			}

			if len(integrationAvailableRoles) == 0 {
				return fmt.Errorf("no available integration roles retrieved")
			}

			selectedRole = &integrationAvailableRoles[0]

			if inbound.RoleId.ValueString() != "" {
				// Validate specified role ID is in the list of available roles
				foundRole := false
				for _, role := range integrationAvailableRoles {
					if strings.EqualFold(role.UniqueId, inbound.RoleId.ValueString()) {
						foundRole = true
						selectedRole = &role

						break
					}
				}

				if !foundRole {
					return fmt.Errorf(
						"invalid role_id '%s' specified for inbound authentication. Must be one of the following: %v",
						inbound.RoleId.ValueString(),
						integrationAvailableRoles,
					)
				}
			}
		}

		authReq := client.SetInboundAuthRequest{
			AuthType: inbound.AuthType.ValueString(),
			Role:     selectedRole,
		}

		authResp, err := r.apiClient.SetInboundAuthentication(tenantId, integrationId, authReq)
		if err != nil {
			return fmt.Errorf("failed to set inbound authentication: %w", err)
		}

		inbound.Token = types.StringValue(authResp.Token)
		inbound.WebhookURL = types.StringValue(authResp.WebhookURL)
	}

	// Set mapping attributes if specified.
	// Always delete existing mappings first so that removed entries are cleaned up,
	// then POST the desired set. During create there are no existing mappings, so the
	// delete step is a no-op.
	existing, err := r.apiClient.GetInstalledMappingAttributes(tenantId, integrationId)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing mapping attributes: %w", err)
	}
	for _, m := range existing {
		if err := r.apiClient.DeleteMappingAttribute(tenantId, integrationId, m.UniqueId); err != nil {
			return fmt.Errorf("failed to delete mapping attribute '%s': %w", m.UniqueId, err)
		}
	}

	if len(inbound.MapAttributes) > 0 {
		mapAttrs, err := r.buildMappingAttributes(tenantId, integrationId, inbound.MapAttributes, false)
		if err != nil {
			return fmt.Errorf("failed to build mapping attributes: %w", err)
		}

		if err := r.apiClient.SetMappingAttributes(tenantId, integrationId, mapAttrs); err != nil {
			return fmt.Errorf("failed to set inbound mapping attributes: %w", err)
		}
	}

	// Set enable drop alerts
	if !inbound.EnableDropAlerts.IsNull() {
		if err := r.apiClient.SetEnableDropAlerts(tenantId, integrationId, inbound.EnableDropAlerts.ValueBool()); err != nil {
			return fmt.Errorf("failed to set drop alerts: %w", err)
		}
	}

	// Assign process definitions
	if !inbound.ProcessDefinitionIds.IsNull() && len(inbound.ProcessDefinitionIds.Elements()) > 0 {
		for _, elem := range inbound.ProcessDefinitionIds.Elements() {
			processId := elem.(types.String).ValueString()
			if err := r.apiClient.AssignProcessDefinition(tenantId, integrationId, processId, true); err != nil {
				return fmt.Errorf("failed to assign process definition '%s': %w", processId, err)
			}
		}
	}

	// Set webhook handshake if auth_type is WEBHOOK and handshake props are specified
	if !inbound.WebhookHandshake.IsNull() && inbound.WebhookHandshake.ValueString() != "" {
		var props map[string]any
		if err := json.Unmarshal([]byte(inbound.WebhookHandshake.ValueString()), &props); err != nil {
			return fmt.Errorf("webhook_handshake must be valid JSON: %w", err)
		}
		if err := r.apiClient.SetWebhookHandshake(tenantId, integrationId, props); err != nil {
			return fmt.Errorf("failed to set webhook handshake: %w", err)
		}
	}

	// Set additional properties (CUSTOM application type only, enforced by ModifyPlan)
	if !inbound.AdditionalProperties.IsNull() && len(inbound.AdditionalProperties.Elements()) > 0 {
		props := make(map[string]string)
		for k, v := range inbound.AdditionalProperties.Elements() {
			props[k] = v.(types.String).ValueString()
		}
		if err := r.apiClient.SetAdditionalProperties(tenantId, integrationId, props); err != nil {
			return fmt.Errorf("failed to set additional properties: %w", err)
		}
	}

	return nil
}

// configureOutbound sets up the outbound notifier configuration and any mapping attributes.
func (r *IntegrationResource) configureOutbound(tenantId, integrationId string, outbound *IntegrationOutboundModel) error {
	notifierReq := client.NotifierRequest{
		Type:         outbound.Type.ValueString(),
		BaseURI:      outbound.BaseURI.ValueString(),
		AuthType:     outbound.AuthType.ValueString(),
		TokenPayload: map[string]any{},
	}

	if !outbound.UserName.IsNull() {
		notifierReq.UserName = outbound.UserName.ValueString()
	}
	if !outbound.Password.IsNull() {
		notifierReq.Password = outbound.Password.ValueString()
	}
	if !outbound.APIKey.IsNull() {
		notifierReq.APIKey = outbound.APIKey.ValueString()
	}
	if !outbound.APISecret.IsNull() {
		notifierReq.APISecret = outbound.APISecret.ValueString()
	}
	if !outbound.AccessTokenURL.IsNull() {
		notifierReq.AccessTokenURL = outbound.AccessTokenURL.ValueString()
	}
	if !outbound.GrantType.IsNull() {
		notifierReq.GrantType = outbound.GrantType.ValueString()
	}
	if !outbound.Scope.IsNull() {
		notifierReq.Scope = outbound.Scope.ValueString()
	}

	if err := r.apiClient.SetNotifier(tenantId, integrationId, notifierReq); err != nil {
		return fmt.Errorf("failed to set notifier: %w", err)
	}

	// Reconcile outbound mapping attributes: delete all existing, then post desired set.
	existing, err := r.apiClient.GetInstalledOutboundMappingAttributes(tenantId, integrationId)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing outbound mapping attributes: %w", err)
	}
	for _, m := range existing {
		if err := r.apiClient.DeleteOutboundMappingAttribute(tenantId, integrationId, m.UniqueId); err != nil {
			return fmt.Errorf("failed to delete outbound mapping attribute '%s': %w", m.UniqueId, err)
		}
	}

	if len(outbound.MapAttributes) > 0 {
		mapAttrs, err := r.buildMappingAttributes(tenantId, integrationId, outbound.MapAttributes, true)
		if err != nil {
			return fmt.Errorf("failed to build outbound mapping attributes: %w", err)
		}
		if err := r.apiClient.SetMappingAttributes(tenantId, integrationId, mapAttrs); err != nil {
			return fmt.Errorf("failed to set outbound mapping attributes: %w", err)
		}
	}

	return nil
}

// buildMappingAttributes converts the Terraform model to the API request.
// It fetches available properties from the API to resolve the display name for each opsramp_attribute.
// Set outbound=true to wrap the result in OutboundConfig instead of InboundConfig.
func (r *IntegrationResource) buildMappingAttributes(
	tenantId, integrationId string,
	models []IntegrationMapAttributes,
	outbound bool,
) (client.MappingAttributesRequest, error) {
	// Cache of entityType -> (property identifier -> display name), fetched lazily per entity type
	propertyCache := make(map[string]map[string]string)
	getPropertyNames := func(entityType string) (map[string]string, error) {
		if names, ok := propertyCache[entityType]; ok {
			return names, nil
		}

		var properties []client.EntityProperty
		var err error
		if outbound {
			properties, err = r.apiClient.GetOutboundEntityProperties(tenantId, integrationId, entityType)
		} else {
			properties, err = r.apiClient.GetInboundEntityProperties(tenantId, integrationId, entityType)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to retrieve inbound properties for entity type '%s': %w", entityType, err)
		}

		names := make(map[string]string, len(properties))
		for _, p := range properties {
			names[p.Property] = p.Name
		}
		propertyCache[entityType] = names

		return names, nil
	}

	attrs := make([]client.MapAttribute, len(models))
	for i, m := range models {
		entityType := m.EntityType.ValueString()
		attrName := m.OpsRampAttribute.ValueString()

		propertyNames, err := getPropertyNames(entityType)
		if err != nil {
			return client.MappingAttributesRequest{}, err
		}

		displayName := propertyNames[attrName]
		if displayName == "" {
			return client.MappingAttributesRequest{}, fmt.Errorf(
				"opsramp_attribute '%s' is not a valid property for entity type '%s' on this integration",
				attrName,
				entityType,
			)
		}

		attr := client.MapAttribute{
			EntityType:           entityType,
			ThirdPartyEntityType: "EVENT",
			Name:                 displayName,
			ThirdPartyAttrName:   m.ThirdPartyAttribute.ValueString(),
			AttrName:             attrName,
			AttrValues:           []client.AttrValueMapping{},
			ParsingProperty: &client.ParsingProperty{
				DefaultValue:  "",
				OprSet:        []client.ParsingOperator{},
				ValueMappings: []any{},
			},
		}

		// Convert attribute value mappings
		if !m.AttributeValues.IsNull() && len(m.AttributeValues.Elements()) > 0 {
			for k, v := range m.AttributeValues.Elements() {
				attr.AttrValues = append(attr.AttrValues, client.AttrValueMapping{
					AttrValue:           k,
					ThirdPartyAttrValue: v.(types.String).ValueString(),
				})
			}
		}

		// Convert parsing property
		if len(m.ParsingOperators) > 0 || (!m.DefaultParsingValue.IsNull() && m.DefaultParsingValue.ValueString() != "") {
			attr.ParsingProperty.DefaultValue = m.DefaultParsingValue.ValueString()
			for _, op := range m.ParsingOperators {
				attr.ParsingProperty.OprSet = append(attr.ParsingProperty.OprSet, client.ParsingOperator{
					Operator:  op.Operator,
					StartWord: op.StartWord,
					EndWord:   op.EndWord,
					RegexStr:  op.RegexStr,
				})
			}
		}

		attrs[i] = attr
	}

	mapAttrs := &client.MappingConfig{MapAttributes: attrs}
	if outbound {
		return client.MappingAttributesRequest{OutboundConfig: mapAttrs}, nil
	}

	return client.MappingAttributesRequest{InboundConfig: mapAttrs}, nil
}

// installedMappingsToModel converts the API response for installed mappings into the Terraform model slice.
// The API returns one row per attribute-value pair; we group them back into a single model entry per
// (entityType, opsrampAttribute, thirdPartyAttribute) key and collect the attribute values.
func (r *IntegrationResource) installedMappingsToModel(mappings []client.InstalledMappingResult) []IntegrationMapAttributes {
	// Return nil (not an empty slice) so that an Optional list attribute whose
	// config value is absent (null) produces no spurious diff in the next plan.
	// Terraform distinguishes [] from null for Optional list attributes.
	if len(mappings) == 0 {
		return nil
	}

	type key struct {
		entity   string
		property string
		tpAttr   string
	}

	// Preserve insertion order using a slice of keys
	order := []key{}
	groups := map[key]*IntegrationMapAttributes{}

	for _, m := range mappings {
		k := key{entity: m.Entity, property: m.Property, tpAttr: m.TenantProperty}
		if _, exists := groups[k]; !exists {
			order = append(order, k)

			model := &IntegrationMapAttributes{
				EntityType:          types.StringValue(m.Entity),
				OpsRampAttribute:    types.StringValue(m.Property),
				ThirdPartyAttribute: types.StringValue(m.TenantProperty),
				AttributeValues:     types.MapNull(types.StringType),
			}

			// Parsing property comes from the first row (same for all rows in the group)
			if m.ParsingProperty != nil {
				model.DefaultParsingValue = types.StringValue(m.ParsingProperty.DefaultValue)
				for _, op := range m.ParsingProperty.OprSet {
					model.ParsingOperators = append(model.ParsingOperators, IntegrationParsingOperators{
						Operator:  op.Operator,
						StartWord: op.StartWord,
						EndWord:   op.EndWord,
						RegexStr:  op.RegexStr,
					})
				}
			} else {
				model.DefaultParsingValue = types.StringValue("")
			}

			groups[k] = model
		}

		// Collect attribute values (propertyValue -> tenantPropertyValue)
		if m.PropertyValue != "" {
			elems := map[string]attr.Value{}
			if !groups[k].AttributeValues.IsNull() {
				for k2, v := range groups[k].AttributeValues.Elements() {
					elems[k2] = v
				}
			}
			elems[m.PropertyValue] = types.StringValue(m.TenantPropertyValue)
			groups[k].AttributeValues = types.MapValueMust(types.StringType, elems)
		}
	}

	result := make([]IntegrationMapAttributes, 0, len(order))
	for _, k := range order {
		result = append(result, *groups[k])
	}

	return result
}

// Read handles reading the resource state.
func (r *IntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IntegrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	existing, err := r.apiClient.GetIntegration(tenantId, state.Id.ValueString())
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "404") ||
			strings.Contains(errStr, "not found") ||
			strings.Contains(errStr, "No installed integration found") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}

	// Update non-sensitive state from API
	state.DisplayName = types.StringValue(existing.DisplayName)
	state.Description = types.StringValue(existing.Description)
	if existing.AlertSource != nil {
		state.AlertSourceID = types.Int64Value(int64(existing.AlertSource.ID))
	}
	state.Status = types.StringValue(existing.Status)
	if existing.Category != "" {
		state.Category = types.StringValue(existing.Category)
	}

	// The v2 API does not return profile_id. Ensure it stays as a known null
	// (not unknown) so Terraform does not emit 'unknown value after apply'.
	if state.ProfileId.IsUnknown() {
		state.ProfileId = types.StringNull()
	}

	if existing.MultiAppsDiscoveryEnabled {
		state.BypassResourceReconciliation = types.BoolValue(existing.MultiAppsDiscoveryEnabled)
	}

	// Preserve inbound sensitive values (token, webhook_url) from state
	// as they may not always be returned in GET responses
	if state.Inbound != nil && existing.InboundConfig != nil && existing.InboundConfig.Authentication != nil {
		auth := existing.InboundConfig.Authentication
		if auth.Token != "" {
			state.Inbound.Token = types.StringValue(auth.Token)
		}
		if auth.WebhookURL != "" {
			state.Inbound.WebhookURL = types.StringValue(auth.WebhookURL)
		}
	}

	// Refresh mapping attributes from the API when inbound is configured
	if state.Inbound != nil {
		mappings, err := r.apiClient.GetInstalledMappingAttributes(tenantId, state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Read Inbound Mapping Attributes Error", err.Error())

			return
		}
		state.Inbound.MapAttributes = r.installedMappingsToModel(mappings)
	}

	// Refresh outbound mapping attributes from the API when outbound is configured
	if state.Outbound != nil {
		outMappings, err := r.apiClient.GetInstalledOutboundMappingAttributes(tenantId, state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Read Outbound Mapping Attributes Error", err.Error())

			return
		}
		state.Outbound.MapAttributes = r.installedMappingsToModel(outMappings)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *IntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IntegrationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IntegrationModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(plan.Client)
	integrationId := state.Id.ValueString()

	// Update base integration fields (displayName, description) if changed
	if plan.DisplayName.ValueString() != state.DisplayName.ValueString() ||
		plan.Description.ValueString() != state.Description.ValueString() ||
		plan.AlertSourceID.ValueInt64() != state.AlertSourceID.ValueInt64() ||
		plan.BypassResourceReconciliation.ValueBool() != state.BypassResourceReconciliation.ValueBool() {
		updateReq := client.InstallIntegrationRequest{
			DisplayName: plan.DisplayName.ValueString(),
			Description: plan.Description.ValueString(),
			AlertSource: &client.AlertSource{
				ID: int(plan.AlertSourceID.ValueInt64()),
			},
			MultiAppsDiscoveryEnabled: plan.BypassResourceReconciliation.ValueBool(),
		}
		_, err := r.apiClient.UpdateIntegration(tenantId, integrationId, updateReq)
		if err != nil {
			resp.Diagnostics.AddError("Integration Update Error", err.Error())

			return
		}
	}

	// Handle process definition unassignment (removed IDs)
	if plan.Inbound != nil && state.Inbound != nil {
		r.reconcileProcessDefinitions(
			tenantId,
			integrationId,
			state.Inbound.ProcessDefinitionIds,
			plan.Inbound.ProcessDefinitionIds,
			resp,
		)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Re-configure inbound if specified
	if plan.Inbound != nil {
		if err := r.configureInbound(tenantId, integrationId, plan.Application.ValueString(), plan.Inbound, nil); err != nil {
			resp.Diagnostics.AddError("Inbound Configuration Error", err.Error())

			return
		}

		// Handle webhook handshake
		r.reconcileWebhookHandshake(tenantId, integrationId, plan.Inbound, state.Inbound, resp)
		if resp.Diagnostics.HasError() {
			return
		}

		// Preserve token/webhook from state if not refreshed
		if plan.Inbound.Token.IsUnknown() || plan.Inbound.Token.IsNull() {
			if state.Inbound != nil {
				plan.Inbound.Token = state.Inbound.Token
				plan.Inbound.WebhookURL = state.Inbound.WebhookURL
			}
		}
	}

	// Re-configure outbound if specified
	if plan.Outbound != nil {
		if err := r.configureOutbound(tenantId, integrationId, plan.Outbound); err != nil {
			resp.Diagnostics.AddError("Outbound Configuration Error", err.Error())

			return
		}
	}

	// Preserve computed values from state when plan doesn't set them
	plan.Id = state.Id
	plan.Status = state.Status
	if plan.DisplayName.IsUnknown() {
		plan.DisplayName = state.DisplayName
	}
	if plan.Category.IsUnknown() {
		plan.Category = state.Category
	}
	if plan.Description.IsUnknown() {
		plan.Description = state.Description
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// reconcileProcessDefinitions unassigns process definitions that were removed from the plan.
func (r *IntegrationResource) reconcileProcessDefinitions(
	tenantId, integrationId string,
	oldSet types.Set,
	newSet types.Set,
	resp *resource.UpdateResponse,
) {
	// Build set of new IDs
	newIds := make(map[string]bool)
	if !newSet.IsNull() && !newSet.IsUnknown() {
		for _, elem := range newSet.Elements() {
			newIds[elem.(types.String).ValueString()] = true
		}
	}

	// Unassign any IDs that were in old state but not in new plan
	if !oldSet.IsNull() && !oldSet.IsUnknown() {
		for _, elem := range oldSet.Elements() {
			processId := elem.(types.String).ValueString()
			if !newIds[processId] {
				if err := r.apiClient.AssignProcessDefinition(tenantId, integrationId, processId, false); err != nil {
					resp.Diagnostics.AddError("Process Definition Unassign Error",
						fmt.Sprintf("Failed to unassign process definition '%s': %s", processId, err.Error()))

					return
				}
			}
		}
	}
}

// reconcileWebhookHandshake manages webhook handshake properties.
func (r *IntegrationResource) reconcileWebhookHandshake(
	tenantId, integrationId string,
	planInbound, stateInbound *IntegrationInboundModel,
	resp *resource.UpdateResponse,
) {
	planHasHandshake := planInbound != nil && !planInbound.WebhookHandshake.IsNull() &&
		planInbound.WebhookHandshake.ValueString() != ""
	stateHasHandshake := stateInbound != nil && !stateInbound.WebhookHandshake.IsNull() &&
		stateInbound.WebhookHandshake.ValueString() != ""

	if planHasHandshake {
		// Set/update webhook handshake
		var props map[string]any
		if err := json.Unmarshal([]byte(planInbound.WebhookHandshake.ValueString()), &props); err != nil {
			resp.Diagnostics.AddError("Invalid Webhook Handshake", "webhook_handshake must be valid JSON: "+err.Error())

			return
		}
		if err := r.apiClient.SetWebhookHandshake(tenantId, integrationId, props); err != nil {
			resp.Diagnostics.AddError("Webhook Handshake Error", err.Error())

			return
		}
	} else if stateHasHandshake && !planHasHandshake {
		// Remove webhook handshake - send DELETE with old props
		var props map[string]any
		if err := json.Unmarshal([]byte(stateInbound.WebhookHandshake.ValueString()), &props); err != nil {
			// If we can't parse old state, send empty map
			props = map[string]any{}
		}
		if err := r.apiClient.DeleteWebhookHandshake(tenantId, integrationId, props); err != nil {
			resp.Diagnostics.AddError("Webhook Handshake Delete Error", err.Error())

			return
		}
	}
}

// Delete handles deleting the resource.
func (r *IntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IntegrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.getTenantId(state.Client)

	err := r.apiClient.DeleteIntegration(tenantId, state.Id.ValueString(), "Terraform - Resource destroyed")
	if err != nil {
		if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError("Delete Error", err.Error())

			return
		}
	}
}

// ImportState handles resource import.
// Import ID format: <integration_id> or <client_id>:<integration_id> (MSP only)
func (r *IntegrationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parsed, err := r.ParseImportID(req.ID, 1)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())

		return
	}

	var state IntegrationModel
	state.Client = parsed.Client
	state.Id = types.StringValue(parsed.Parts[0])

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *IntegrationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Call base implementation
	r.BaseResource.ModifyPlan(ctx, req, resp)

	// Don't modify plan during destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state IntegrationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	req.State.Get(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	applicationWithDisplayName := []string{"CUSTOM", "CUSTOM-EVENT"}

	// For CUSTOM/CUSTOM-EVENT, display_name must be provided.
	if contains(applicationWithDisplayName, plan.Application.ValueString()) && plan.DisplayName.ValueString() == "" {
		resp.Diagnostics.AddError(
			"display_name is required for application types CUSTOM and CUSTOM-EVENT",
			"Please provide a display_name value.",
		)

		return
	}

	if plan.Application.ValueString() == "CUSTOM-EVENT" {
		if plan.AlertSourceID.IsNull() {
			resp.Diagnostics.AddError(
				"alert_source_id is required for application type CUSTOM-EVENT",
				"Please provide an alert_source_id value.",
			)

			return
		}

		// Auto-set category to Monitoring for CUSTOM-EVENT
		plan.Category = types.StringValue("Monitoring")
		diags = resp.Plan.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else if plan.Application.ValueString() == "CUSTOM" {
		if plan.Inbound == nil || plan.Inbound.AuthType.ValueString() != "OAUTH2" {
			resp.Diagnostics.AddError(
				"For application type CUSTOM, inbound auth_type must be OAUTH2",
				"Please set inbound auth_type to OAUTH2.",
			)

			return
		}
	} else {
		if plan.Inbound != nil && !plan.Inbound.RoleId.IsNull() && plan.Inbound.RoleId.ValueString() != "" {
			resp.Diagnostics.AddError(
				"role_id can only be set for application types CUSTOM or CUSTOM-EVENT",
				"Remove the inbound role_id value or set application to CUSTOM or CUSTOM-EVENT.",
			)

			return
		}
		if plan.Inbound != nil && !plan.Inbound.AdditionalProperties.IsNull() && len(plan.Inbound.AdditionalProperties.Elements()) > 0 {
			resp.Diagnostics.AddError(
				"inbound additional_properties can only be set for application type CUSTOM",
				"Remove the inbound additional_properties value or set application to CUSTOM.",
			)

			return
		}
	}
	if plan.Application.ValueString() == "CUSTOM-EVENT" && plan.Inbound != nil &&
		!plan.Inbound.AdditionalProperties.IsNull() && len(plan.Inbound.AdditionalProperties.Elements()) > 0 {
		resp.Diagnostics.AddError(
			"inbound additional_properties can only be set for application type CUSTOM",
			"Remove the inbound additional_properties value or set application to CUSTOM.",
		)

		return
	}

	// Validate default_parsing_value is set when parsing_operators is provided
	if plan.Inbound != nil && len(plan.Inbound.MapAttributes) > 0 {
		for _, attr := range plan.Inbound.MapAttributes {
			if len(attr.ParsingOperators) > 0 &&
				(attr.DefaultParsingValue.IsNull() || attr.DefaultParsingValue.ValueString() == "") {
				resp.Diagnostics.AddError(
					"default_parsing_value is required when parsing_operators is set",
					fmt.Sprintf(
						"Attribute mapping for '%s' has parsing_operators but no default_parsing_value.",
						attr.OpsRampAttribute.ValueString(),
					),
				)

				return
			}

			// Validate operator-specific fields
			for _, op := range attr.ParsingOperators {
				switch op.Operator {
				case "BEFORE":
					if op.EndWord == "" {
						resp.Diagnostics.AddError(
							"end_word is required for BEFORE operator",
							fmt.Sprintf(
								"Attribute mapping for '%s' uses BEFORE operator but end_word is not set.",
								attr.OpsRampAttribute.ValueString(),
							),
						)

						return
					}
				case "AFTER":
					if op.StartWord == "" {
						resp.Diagnostics.AddError(
							"start_word is required for AFTER operator",
							fmt.Sprintf(
								"Attribute mapping for '%s' uses AFTER operator but start_word is not set.",
								attr.OpsRampAttribute.ValueString(),
							),
						)

						return
					}
				case "BETWEEN":
					if op.StartWord == "" || op.EndWord == "" {
						resp.Diagnostics.AddError(
							"start_word and end_word are required for BETWEEN operator",
							fmt.Sprintf(
								"Attribute mapping for '%s' uses BETWEEN operator but start_word and/or end_word is not set.",
								attr.OpsRampAttribute.ValueString(),
							),
						)

						return
					}
				case "MATCHES":
					if op.RegexStr == "" {
						resp.Diagnostics.AddError(
							"regex_str is required for MATCHES operator",
							fmt.Sprintf(
								"Attribute mapping for '%s' uses MATCHES operator but regex_str is not set.",
								attr.OpsRampAttribute.ValueString(),
							),
						)

						return
					}
				}
			}
		}
	}
}
