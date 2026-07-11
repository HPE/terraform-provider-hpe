// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"regexp"
	"strings"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &AlertEscalationPolicyResource{}
var _ resource.ResourceWithModifyPlan = &AlertEscalationPolicyResource{}

// AlertEscalationPolicyResource defines the resource implementation.
type AlertEscalationPolicyResource struct {
	BaseResource
}

// AlertEscalationPolicyModel maps Terraform schema attributes to the provider model.
type AlertEscalationPolicyModel struct {
	Client          types.String           `tfsdk:"client"`
	Id              types.String           `tfsdk:"id"`
	Name            types.String           `tfsdk:"name"`
	Description     types.String           `tfsdk:"description"`
	Precedence      types.Int64            `tfsdk:"precedence"`
	EscalationType  types.String           `tfsdk:"escalation_type"`
	PolicyType      types.String           `tfsdk:"policy_type"`
	EnabledMode     types.String           `tfsdk:"enabled_mode"`
	IncludedClients types.Set              `tfsdk:"included_clients"`
	Escalations     []EscalationLevelModel `tfsdk:"escalations"`

	SearchQuery         types.String `tfsdk:"search_query"`
	ResourceSearchQuery types.String `tfsdk:"resource_search_query"`
}

type EscalationResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type EscalationLevelModel struct {
	WaitMins               types.Int64                    `tfsdk:"wait_mins"`
	Action                 types.String                   `tfsdk:"action"`
	Priority               types.String                   `tfsdk:"priority"`
	RepeatFrequency        types.Int64                    `tfsdk:"repeat_frequency"`
	NotifyLimitCount       types.Int64                    `tfsdk:"notify_limit_count"`
	NotificationType       types.String                   `tfsdk:"notification_type"`
	NotificationTemplateId types.String                   `tfsdk:"notification_template_id"`
	Recipients             []EscalationRecipientModel     `tfsdk:"recipients"`
	Incident               *EscalationIncidentModel       `tfsdk:"incident"`
	UpdateIncident         *EscalationUpdateIncidentModel `tfsdk:"update_incident"`
}

type EscalationRecipientModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type EscalationIncidentModel struct {
	Priority            types.String   `tfsdk:"priority"`
	Subject             types.String   `tfsdk:"subject"`
	Description         types.String   `tfsdk:"description"`
	AssigneeGroupId     types.String   `tfsdk:"assignee_group_id"`
	AssignedUserId      types.String   `tfsdk:"assigned_user_id"`
	CategoryId          types.String   `tfsdk:"category_id"`
	SubCategoryId       types.String   `tfsdk:"sub_category_id"`
	BusinessImpactId    types.String   `tfsdk:"business_impact_id"`
	UrgencyId           types.String   `tfsdk:"urgency_id"`
	RosterId            types.String   `tfsdk:"roster_id"`
	KnowledgeArticleIds []types.String `tfsdk:"knowledge_article_ids"`
	Cc                  types.String   `tfsdk:"cc"`
}

type EscalationUpdateIncidentModel struct {
	UpdateIncidentMode              types.String                  `tfsdk:"update_incident_mode"`
	UpdateIncidentSubjectMode       types.String                  `tfsdk:"update_incident_subject_mode"`
	AutoResolveIncidentMode         types.String                  `tfsdk:"auto_resolve_incident_mode"`
	AutoHealWaitTime                types.Int64                   `tfsdk:"auto_heal_wait_time"`
	UpdatePriorityByMLConfiguration types.Bool                    `tfsdk:"update_priority_by_ml_configuration"`
	PriorityRules                   []EscalationPriorityRuleModel `tfsdk:"priority_rules"`
}

type EscalationPriorityRuleModel struct {
	AlertState       types.String `tfsdk:"alert_state"`
	BusinessImpactId types.String `tfsdk:"business_impact_id"`
	UrgencyId        types.String `tfsdk:"urgency_id"`
	Priority         types.String `tfsdk:"priority"`
}

// NewAlertEscalationPolicy creates a new instance of the resource.
func NewAlertEscalationPolicy() resource.Resource {
	return &AlertEscalationPolicyResource{}
}

// Metadata returns the resource type name.
func (r *AlertEscalationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_escalation_policy"
}

// Schema defines the schema for the resource.
func (r *AlertEscalationPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Alert Escalation Policy. Defines automated notifications and incident creation for escalated alerts.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this policy should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the alert escalation policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the alert escalation policy.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The description of the alert escalation policy.",
			},
			"precedence": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The execution order of the policy.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"escalation_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Automated notification type. Valid values: `AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED_TICKETED`, `AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED`, `MANUAL`, `AUTOMATIC`.",
				Validators: []validator.String{
					stringvalidator.OneOf("AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED_TICKETED", "AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED", "MANUAL", "AUTOMATIC"),
				},
			},
			"policy_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The policy type. Value: `ESCALATION_POLICY`.",
				Default:             stringdefault.StaticString("ESCALATION_POLICY"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("ESCALATION_POLICY"),
				},
			},
			"enabled_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The enabled mode. Valid values: `ON`, `OFF`, `RECOMMEND`, `OBSERVED`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ON", "OFF", "RECOMMEND", "OBSERVED"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"included_clients": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of included client UUIDs. Only applicable when the provider is in MSP scope. If empty (or omitted), the policy applies to all clients.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"escalations": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "Escalation steps. Each step can be a NOTIFICATION or INCIDENT action.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"wait_mins": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "Minutes to wait before executing this escalation step.",
						},
						"action": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The action type. Valid values: `NOTIFICATION`, `INCIDENT`.",
							Validators: []validator.String{
								stringvalidator.OneOf("NOTIFICATION", "INCIDENT"),
							},
						},

						"priority": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "The priority for notifications (e.g. `Normal`, `High`, `Low`, `Urgent`).",
							Validators: []validator.String{
								stringvalidator.OneOf("Normal", "Low", "High", "Urgent"),
							},
							Default: stringdefault.StaticString("Low"),
						},
						"repeat_frequency": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Number of minutes between repeat notifications.",
							Default:             int64default.StaticInt64(5),
						},
						"notify_limit_count": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Maximum number of notifications to send.",
							Default:             int64default.StaticInt64(2),
						},
						"notification_type": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Notification template type (e.g. `basic`, `advanced`).",
							Validators: []validator.String{
								stringvalidator.OneOf("basic", "advanced"),
							},
							Default: stringdefault.StaticString("basic"),
						},
						"notification_template_id": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "The ID of the notification template.",
						},
						"recipients": schema.SetNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Recipients for notification actions.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The recipient ID.",
									},
									"type": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The recipient type (e.g. `USERGROUP`, `USER`, `USERGROUP_DL`).",
										Validators: []validator.String{
											stringvalidator.OneOf("USER", "USERGROUP", "USERGROUP_DL"),
										},
									},
								},
							},
						},

						"incident": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Incident creation settings. Used when action is `INCIDENT`.",
							Attributes: map[string]schema.Attribute{
								"priority": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Incident priority (e.g. `Normal`, `Low`, `High`, `Urgent`).",
									Validators: []validator.String{
										stringvalidator.OneOf("Normal", "Low", "High", "Urgent"),
									},
								},
								"subject": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Incident subject. Supports placeholders like `$alert.subject`.",
								},
								"description": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString(""),
									MarkdownDescription: "The description of the Incident. Supports placeholders like `$alert.description`.",
								},
								"assignee_group_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The uniqueId of the assignee user group.",
								},
								"assigned_user_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The ID of the assigned user.",
								},
								"category_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The uniqueId of the service desk category.",
								},
								"sub_category_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The uniqueId of the service desk sub-category.",
								},
								"business_impact_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The uniqueId of the business impact.",
								},
								"urgency_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The uniqueId of the urgency.",
								},
								"roster_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The ID of the roster.",
								},
								"knowledge_article_ids": schema.SetAttribute{
									Optional:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of knowledge base article IDs to attach.",
								},
								"cc": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "CC email addresses for the incident.",
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^([a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,})(\s*,\s*[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,})*$`),
											"must be a valid email address or a comma-separated list of valid email addresses",
										),
									},
								},
							},
						},
						"update_incident": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Incident update settings. Used when action is `INCIDENT`.",
							Attributes: map[string]schema.Attribute{
								"update_incident_mode": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "The mode for updating incidents. Valid values: `UpdateWhenAlertStateChange`, `UpdateWithRuleWhenAlertStateChange`, `UpdateForEveryRepeatAlert`, `UpdateWithRuleForEveryRepeatAlert`.",
									Validators: []validator.String{
										stringvalidator.OneOf("UpdateWhenAlertStateChange", "UpdateWithRuleWhenAlertStateChange", "UpdateForEveryRepeatAlert", "UpdateWithRuleForEveryRepeatAlert"),
									},
								},

								"update_incident_subject_mode": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "The mode for updating incident subject. Valid values: `UpdateIncidentSubject`, `UpdateIncidentSubjectWithRule`.",
									Validators: []validator.String{
										stringvalidator.OneOf("UpdateIncidentSubject", "UpdateIncidentSubjectWithRule"),
									},
								},

								"auto_resolve_incident_mode": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "The mode for resolving incidents. Valid values: `AutoResolveIncident`, `AutoResolveUnassignedIncident`.",
									Validators: []validator.String{
										stringvalidator.OneOf("AutoResolveIncident", "AutoResolveUnassignedIncident"),
									},
								},
								"auto_heal_wait_time": schema.Int64Attribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Wait time in minutes before auto-healing. Valid values: 0, 15, 30, 45, 60.",
									Validators: []validator.Int64{
										int64validator.OneOf(0, 15, 30, 45, 60),
									},
									Default: int64default.StaticInt64(0),
								},

								"update_priority_by_ml_configuration": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Update priority by ML configuration.",
								},
								"priority_rules": schema.ListNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Priority rules for incident updates.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"alert_state": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												MarkdownDescription: "The value to match (e.g. `WARNING`, `CRITICAL`).",
												Validators: []validator.String{
													stringvalidator.OneOf("WARNING", "CRITICAL"),
												},
												Default: stringdefault.StaticString("WARNING"),
											},
											"business_impact_id": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The uniqueId of the resulting business impact for this rule.",
											},
											"urgency_id": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The uniqueId of the resulting urgency for this rule.",
											},
											"priority": schema.StringAttribute{
												Optional:            true,
												Computed:            true,
												MarkdownDescription: "The resulting priority (e.g. `Very Low`, `Low`, `Normal`, `High`, `Urgent`).",
												Validators: []validator.String{
													stringvalidator.OneOf("Very Low", "Low", "Normal", "High", "Urgent"),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"search_query": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Filter query to scope which alerts are evaluated by this policy.",
			},
			"resource_search_query": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Filter query to scope which resources are evaluated by this policy.",
			},
		},
	}
}

func buildAlertEscalationPolicyRequest(plan AlertEscalationPolicyModel, apiClient *client.OpsRampClient) client.AlertEscalationPolicy {
	// default values
	policy := client.AlertEscalationPolicy{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		EscalationType: plan.EscalationType.ValueString(),
		PolicyType:     plan.PolicyType.ValueString(),
		EnabledMode:    plan.EnabledMode.ValueString(),
		Scope:          &client.EscalationScope{Uuid: apiClient.TenantId},
		Precedence:     nil,                           // unless determined by logic below
		AllClients:     false,                         // unless determined by logic below
		Resources:      []client.EscalationResource{}, // unless determined by logic below
	}

	if !plan.Precedence.IsNull() && !plan.Precedence.IsUnknown() {
		value := int(plan.Precedence.ValueInt64())
		policy.Precedence = &value
	}

	// retrieve tenantScope based on logic (tenantScope means whether the policy is created at partner-level or client-level)
	isMSP := apiClient.Scope == "MSP" && (plan.Client.IsNull() || plan.Client.ValueString() == "")
	if isMSP {
		// tenantScope = MSP

		var elems []types.String
		_ = plan.IncludedClients.ElementsAs(context.Background(), &elems, false)
		for _, elem := range elems {
			policy.IncludedClients = append(policy.IncludedClients, client.ClientRef{UniqueId: elem.ValueString()})
		}
		policy.AllClients = len(elems) == 0
	} else {
		// tenantScope = Client

		// Use client parameter if set, otherwise provider's tenant ID
		tenantId := apiClient.TenantId
		if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
			tenantId = plan.Client.ValueString()
		}

		policy.Resources = []client.EscalationResource{{Id: tenantId, Type: "CLIENT"}}
	}

	// Create FilterCriteria if either search query are provided
	if plan.SearchQuery.ValueString() != "" || plan.ResourceSearchQuery.ValueString() != "" {
		policy.FilterCriteria = &client.EscalationFilterCriteria{
			MatchingType:        "ALL",
			SearchQuery:         plan.SearchQuery.ValueString(),
			ResourceSearchQuery: plan.ResourceSearchQuery.ValueString(),
		}
	}

	// Manage escalation list
	for _, e := range plan.Escalations {
		level := client.EscalationLevel{
			WaitMins: int(e.WaitMins.ValueInt64()),
			Action:   e.Action.ValueString(),
		}

		if level.Action == "NOTIFICATION" {
			level.NotificationTemplateId = e.NotificationTemplateId.ValueString()
			level.NotificationType = e.NotificationType.ValueString()
			level.Priority = e.Priority.ValueString()
			level.RepeatFrequency = int(e.RepeatFrequency.ValueInt64())
			level.NotifyLimitCount = int(e.NotifyLimitCount.ValueInt64())

			for _, rec := range e.Recipients {
				level.Recipients = append(level.Recipients, client.EscalationRecipient{
					Id:   rec.Id.ValueString(),
					Type: rec.Type.ValueString(),
				})
			}
		}

		if level.Action == "INCIDENT" {
			inc := &client.EscalationIncident{
				Priority:    e.Incident.Priority.ValueString(),
				Subject:     e.Incident.Subject.ValueString(),
				Description: e.Incident.Description.ValueString(),
			}

			if !e.Incident.AssigneeGroupId.IsNull() && e.Incident.AssigneeGroupId.ValueString() != "" {
				inc.AssigneeGroup = &client.EscalationUniqueRef{UniqueId: e.Incident.AssigneeGroupId.ValueString()}
			}
			if !e.Incident.AssignedUserId.IsNull() && e.Incident.AssignedUserId.ValueString() != "" {
				inc.AssignedUser = &client.EscalationUserRef{
					Id: e.Incident.AssignedUserId.ValueString(),
				}
			}

			if !e.Incident.CategoryId.IsNull() && e.Incident.CategoryId.ValueString() != "" {
				inc.Category = &client.EscalationUniqueRef{UniqueId: e.Incident.CategoryId.ValueString()}
			}
			if !e.Incident.SubCategoryId.IsNull() && e.Incident.SubCategoryId.ValueString() != "" {
				inc.SubCategory = &client.EscalationUniqueRef{UniqueId: e.Incident.SubCategoryId.ValueString()}
			}
			if !e.Incident.BusinessImpactId.IsNull() && e.Incident.BusinessImpactId.ValueString() != "" {
				inc.BusinessImpact = &client.EscalationUniqueRef{UniqueId: e.Incident.BusinessImpactId.ValueString()}
			}
			if !e.Incident.UrgencyId.IsNull() && e.Incident.UrgencyId.ValueString() != "" {
				inc.Urgency = &client.EscalationUniqueRef{UniqueId: e.Incident.UrgencyId.ValueString()}
			}
			if !e.Incident.RosterId.IsNull() && e.Incident.RosterId.ValueString() != "" {
				inc.NotifyRoster = &client.EscalationRosterRef{Id: e.Incident.RosterId.ValueString()}
			}

			for _, kid := range e.Incident.KnowledgeArticleIds {
				inc.KnowledgeArticleIds = append(inc.KnowledgeArticleIds, kid.ValueString())
				inc.AttachedArticles = append(inc.AttachedArticles, client.EscalationArticleRef{Id: kid.ValueString()})
			}
			if !e.Incident.Cc.IsNull() && e.Incident.Cc.ValueString() != "" {
				inc.Cc = e.Incident.Cc.ValueString()
			}

			level.Incident = inc
		}

		if e.UpdateIncident != nil {

			ui := &client.EscalationUpdateIncident{
				UpdateWhenAlertStateChange:         false,
				UpdateForEveryRepeatAlert:          false,
				UpdateWithRuleWhenAlertStateChange: false,
				UpdateWithRuleForEveryRepeatAlert:  false,
				UpdateIncidentSubject:              false,
				UpdateIncidentSubjectWithRule:      false,
				AutoResolveIncident:                false,
				AutoResolveUnassignedIncident:      false,
				AutoHealWaitTime:                   int(e.UpdateIncident.AutoHealWaitTime.ValueInt64()),
				UpdatePriorityByMLConfiguration:    e.UpdateIncident.UpdatePriorityByMLConfiguration.ValueBool(),
			}

			// select UpdateIncidentMode
			switch e.UpdateIncident.UpdateIncidentMode.ValueString() {
			case "UpdateWhenAlertStateChange":
				ui.UpdateWhenAlertStateChange = true
			case "UpdateWithRuleWhenAlertStateChange":
				ui.UpdateWithRuleWhenAlertStateChange = true
			case "UpdateForEveryRepeatAlert":
				ui.UpdateForEveryRepeatAlert = true
			case "UpdateWithRuleForEveryRepeatAlert":
				ui.UpdateWithRuleForEveryRepeatAlert = true
			}

			// select UpdateIncidentSubjectMode
			switch e.UpdateIncident.UpdateIncidentSubjectMode.ValueString() {
			case "UpdateIncidentSubject":
				ui.UpdateIncidentSubject = true
			case "UpdateIncidentSubjectWithRule":
				ui.UpdateIncidentSubjectWithRule = true
			}

			// select AutoResolveIncidentMode
			switch e.UpdateIncident.AutoResolveIncidentMode.ValueString() {
			case "AutoResolveIncident":
				ui.AutoResolveIncident = true
			case "AutoResolveUnassignedIncident":
				ui.AutoResolveUnassignedIncident = true
			}

			// manage priority rules
			for _, pr := range e.UpdateIncident.PriorityRules {
				rule := client.EscalationPriorityRule{
					Key:      "currentState.code",
					Operator: "Is",
					Value:    pr.AlertState.ValueString(), // alert_state is transformed into value for the rule with key "currentState.code"
					Priority: pr.Priority.ValueString(),
				}

				if !pr.BusinessImpactId.IsNull() && pr.BusinessImpactId.ValueString() != "" {
					rule.BusinessImpact = &client.EscalationUniqueRef{UniqueId: pr.BusinessImpactId.ValueString()}
				}

				if !pr.UrgencyId.IsNull() && pr.UrgencyId.ValueString() != "" {
					rule.Urgency = &client.EscalationUniqueRef{UniqueId: pr.UrgencyId.ValueString()}
				}

				ui.PriorityRules = append(ui.PriorityRules, rule)
			}
			level.UpdateIncident = ui
		}
		policy.Escalations = append(policy.Escalations, level)
	}

	return policy
}

func mapAlertEscalationPolicyToState(resp *client.AlertEscalationPolicy, state *AlertEscalationPolicyModel, isMSP bool) {
	state.Id = types.StringValue(resp.Id)
	state.Name = types.StringValue(resp.Name)
	state.Description = types.StringValue(resp.Description)
	state.EscalationType = types.StringValue(resp.EscalationType)
	state.PolicyType = types.StringValue(resp.PolicyType)
	state.EnabledMode = types.StringValue(resp.EnabledMode)

	if isMSP {
		// MSP scope: included_clients is always present ([] means "all clients").
		elems := make([]attr.Value, 0, len(resp.IncludedClients))
		for _, c := range resp.IncludedClients {
			elems = append(elems, types.StringValue(c.UniqueId))
		}
		set, _ := types.SetValue(types.StringType, elems)
		state.IncludedClients = set
	} else {
		// CLIENT scope: not applicable.
		state.IncludedClients = types.SetNull(types.StringType)
	}

	if resp.Precedence != nil {
		state.Precedence = types.Int64Value(int64(*resp.Precedence))
	} else {
		state.Precedence = types.Int64Null()
	}

	if len(resp.Escalations) > 0 {
		state.Escalations = make([]EscalationLevelModel, len(resp.Escalations))
		for i, e := range resp.Escalations {
			level := EscalationLevelModel{
				WaitMins:         types.Int64Value(int64(e.WaitMins)),
				Action:           types.StringValue(e.Action),
				Priority:         types.StringValue("Low"),
				RepeatFrequency:  types.Int64Value(5),
				NotifyLimitCount: types.Int64Value(2),
				NotificationType: types.StringValue("basic"),
			}

			if e.Priority != "" {
				level.Priority = types.StringValue(e.Priority)
			}
			if e.RepeatFrequency != 0 {
				level.RepeatFrequency = types.Int64Value(int64(e.RepeatFrequency))
			}
			if e.NotifyLimitCount != 0 {
				level.NotifyLimitCount = types.Int64Value(int64(e.NotifyLimitCount))
			}
			if e.NotificationType != "" {
				level.NotificationType = types.StringValue(e.NotificationType)
			}

			if e.Action == "NOTIFICATION" && e.NotificationTemplateId != "" {
				level.NotificationTemplateId = types.StringValue(e.NotificationTemplateId)
			} else {
				// notification_template_id is only meaningful for NOTIFICATION actions.
				level.NotificationTemplateId = types.StringNull()
			}

			for _, rec := range e.Recipients {
				level.Recipients = append(level.Recipients, EscalationRecipientModel{
					Id:   types.StringValue(rec.Id),
					Type: types.StringValue(rec.Type),
				})
			}

			if e.Incident != nil {
				inc := &EscalationIncidentModel{
					Priority:            types.StringValue(e.Incident.Priority),
					Subject:             types.StringValue(e.Incident.Subject),
					Description:         types.StringValue(e.Incident.Description),
					AssigneeGroupId:     types.StringValue(""),
					AssignedUserId:      types.StringNull(),
					CategoryId:          types.StringValue(""),
					SubCategoryId:       types.StringValue(""),
					BusinessImpactId:    types.StringValue(""),
					UrgencyId:           types.StringValue(""),
					RosterId:            types.StringNull(),
					KnowledgeArticleIds: []types.String{},
					Cc:                  types.StringValue(e.Incident.Cc),
				}
				if e.Incident.AssigneeGroup != nil {
					inc.AssigneeGroupId = types.StringValue(e.Incident.AssigneeGroup.UniqueId)
				}
				if e.Incident.AssignedUser != nil {
					inc.AssignedUserId = types.StringValue(e.Incident.AssignedUser.Id)
				}
				if e.Incident.Category != nil {
					inc.CategoryId = types.StringValue(e.Incident.Category.UniqueId)
				}
				if e.Incident.SubCategory != nil {
					inc.SubCategoryId = types.StringValue(e.Incident.SubCategory.UniqueId)
				}
				if e.Incident.BusinessImpact != nil {
					inc.BusinessImpactId = types.StringValue(e.Incident.BusinessImpact.UniqueId)
				}
				if e.Incident.Urgency != nil {
					inc.UrgencyId = types.StringValue(e.Incident.Urgency.UniqueId)
				}
				if e.Incident.NotifyRoster != nil {
					inc.RosterId = types.StringValue(e.Incident.NotifyRoster.Id)
				}
				for _, kid := range e.Incident.KnowledgeArticleIds {
					inc.KnowledgeArticleIds = append(inc.KnowledgeArticleIds, types.StringValue(kid))
				}

				level.Incident = inc
			}

			if e.UpdateIncident != nil {

				// Determine UpdateIncidentSubjectMode
				updateIncidentSubjectMode := "UpdateIncidentSubject"
				if e.UpdateIncident.UpdateIncidentSubjectWithRule {
					updateIncidentSubjectMode = "UpdateIncidentSubjectWithRule"
				}

				// Determine AutoResolveIncidentMode
				autoResolveIncidentMode := "AutoResolveIncident"
				if e.UpdateIncident.AutoResolveUnassignedIncident {
					autoResolveIncidentMode = "AutoResolveUnassignedIncident"
				}

				// Determine UpdateIncidentMode
				updateIncidentMode := "UpdateWhenAlertStateChange" // default
				if e.UpdateIncident.UpdateWithRuleWhenAlertStateChange {
					updateIncidentMode = "UpdateWithRuleWhenAlertStateChange"
				} else if e.UpdateIncident.UpdateForEveryRepeatAlert {
					updateIncidentMode = "UpdateForEveryRepeatAlert"
				} else if e.UpdateIncident.UpdateWithRuleForEveryRepeatAlert {
					updateIncidentMode = "UpdateWithRuleForEveryRepeatAlert"
				}

				ui := &EscalationUpdateIncidentModel{
					UpdateIncidentMode:              types.StringValue(updateIncidentMode),
					UpdateIncidentSubjectMode:       types.StringValue(updateIncidentSubjectMode),
					AutoResolveIncidentMode:         types.StringValue(autoResolveIncidentMode),
					AutoHealWaitTime:                types.Int64Value(int64(e.UpdateIncident.AutoHealWaitTime)),
					UpdatePriorityByMLConfiguration: types.BoolValue(e.UpdateIncident.UpdatePriorityByMLConfiguration),
					PriorityRules:                   []EscalationPriorityRuleModel{},
				}

				for _, pr := range e.UpdateIncident.PriorityRules {
					rule := EscalationPriorityRuleModel{
						AlertState: types.StringValue(pr.Value),
						Priority:   types.StringValue(pr.Priority),
					}
					if pr.BusinessImpact != nil {
						rule.BusinessImpactId = types.StringValue(pr.BusinessImpact.UniqueId)
					}
					if pr.Urgency != nil {
						rule.UrgencyId = types.StringValue(pr.Urgency.UniqueId)
					}
					ui.PriorityRules = append(ui.PriorityRules, rule)
				}
				level.UpdateIncident = ui
			}
			state.Escalations[i] = level
		}
	}

	if resp.FilterCriteria != nil {
		state.SearchQuery = types.StringValue(resp.FilterCriteria.SearchQuery)
		state.ResourceSearchQuery = types.StringValue(resp.FilterCriteria.ResourceSearchQuery)
	} else {
		state.SearchQuery = types.StringValue("")
		state.ResourceSearchQuery = types.StringValue("")
	}
}

// Create handles the creation of the resource.
func (r *AlertEscalationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlertEscalationPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	policy := buildAlertEscalationPolicyRequest(plan, r.apiClient)

	created, err := r.apiClient.CreateAlertEscalationPolicy(tenantId, policy)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	isMSP := r.apiClient.Scope == "MSP" && (plan.Client.IsNull() || plan.Client.ValueString() == "")
	mapAlertEscalationPolicyToState(created, &plan, isMSP)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *AlertEscalationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlertEscalationPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetAlertEscalationPolicy(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	mapAlertEscalationPolicyToState(existing, &state, r.apiClient.Scope == "MSP" && (state.Client.IsNull() || state.Client.ValueString() == ""))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *AlertEscalationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AlertEscalationPolicyModel
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

	policy := buildAlertEscalationPolicyRequest(plan, r.apiClient)

	updated, err := r.apiClient.UpdateAlertEscalationPolicy(tenantId, state.Id.ValueString(), policy)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	isMSPUpdate := r.apiClient.Scope == "MSP" && (plan.Client.IsNull() || plan.Client.ValueString() == "")
	mapAlertEscalationPolicyToState(updated, &plan, isMSPUpdate)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *AlertEscalationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlertEscalationPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteAlertEscalationPolicy(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
