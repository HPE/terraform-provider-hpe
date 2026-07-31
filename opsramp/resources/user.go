// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
	_ resource.ResourceWithModifyPlan  = &UserResource{}
)

// UserResource defines the resource implementation.
type UserResource struct {
	BaseResource
}

// UserNotificationModel maps Terraform schema attributes for notifications.
type UserNotificationModel struct {
	NotifyType            types.String `tfsdk:"notify_type"`
	NotifyMethod          types.String `tfsdk:"notify_method"`
	NotifyInputType       types.String `tfsdk:"notify_input_type"`
	NotifyRecurringReport types.Bool   `tfsdk:"notify_recurring_report"`
}

// UserModel maps Terraform schema attributes to the provider model.
type UserModel struct {
	Client            types.String            `tfsdk:"client"`
	Id                types.String            `tfsdk:"id"`
	LoginName         types.String            `tfsdk:"login_name"`
	Password          types.String            `tfsdk:"password"`
	FirstName         types.String            `tfsdk:"first_name"`
	LastName          types.String            `tfsdk:"last_name"`
	Designation       types.String            `tfsdk:"designation"`
	Address           types.String            `tfsdk:"address"`
	City              types.String            `tfsdk:"city"`
	State             types.String            `tfsdk:"state"`
	Zip               types.String            `tfsdk:"zip"`
	Country           types.String            `tfsdk:"country"`
	Email             types.String            `tfsdk:"email"`
	AltEmail          types.String            `tfsdk:"alt_email"`
	PhoneNumber       types.String            `tfsdk:"phone_number"`
	MobileNumber      types.String            `tfsdk:"mobile_number"`
	TimeZone          types.String            `tfsdk:"time_zone"`
	AuthType          types.String            `tfsdk:"auth_type"`
	Status            types.String            `tfsdk:"status"`
	Roles             []types.String          `tfsdk:"roles"`
	UserGroups        types.Set               `tfsdk:"user_groups"`
	UserNotifications []UserNotificationModel `tfsdk:"user_notifications"`
	ChangePassword    types.Bool              `tfsdk:"change_password"`
}

func optionalStringValue(v string) types.String {
	if strings.TrimSpace(v) == "" {
		return types.StringNull()
	}

	return types.StringValue(v)
}

// NewUser creates a new instance of the resource.
func NewUser() resource.Resource {
	return &UserResource{}
}

// Metadata returns the resource type name.
func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp User resource.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this user should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"login_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The login name (username) of the user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The password for the user. Required on creation.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(7),
					stringvalidator.RegexMatches(regexp.MustCompile(`[0-9]`), "Password must contain at least one number (0-9)."),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`[#!$*]`),
						"Password must contain at least one special character (#, !, $, *).",
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`[a-z]`),
						"Password must contain at least one lowercase letter (a-z).",
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`[A-Z]`),
						"Password must contain at least one uppercase letter (A-Z).",
					),
				},
			},
			"first_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The first name of the user.",
			},
			"last_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The last name of the user.",
			},
			"designation": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The designation/title of the user.",
			},
			"address": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The street address of the user.",
			},
			"city": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The city of the user.",
			},
			"state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The state/province of the user.",
			},
			"zip": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The zip/postal code of the user.",
			},
			"country": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The country of the user.",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The email address of the user.",
			},
			"alt_email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The alternate email address of the user.",
			},
			"phone_number": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The phone number of the user.",
			},
			"mobile_number": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The mobile number of the user.",
			},
			"time_zone": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("America/Los_Angeles"),
				MarkdownDescription: "The time zone of the user.",
			},
			"auth_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("LOCAL"),
				MarkdownDescription: "The authentication type (LOCAL, SSO, etc.).",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the user (`Active`, `DEACTIVATE`, `TERMINATE`).",
			},
			"roles": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				MarkdownDescription: "List of role names assigned to the user.",
			},
			"user_groups": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of user group IDs the user belongs to. Read-only - manage user_group assignments via the user_group resource.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"user_notifications": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "User notification preferences.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"notify_type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The notification type (`Account Information`, `Alert Notification`, `Report Notification`, `Export Notification`, `Login Activity Notification`).",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"Account Information",
									"Alert Notification",
									"Report Notification",
									"Export Notification",
									"Login Activity Notification",
								),
							},
						},
						"notify_method": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The notification method (`Email`, `No Notify`).",
							Validators: []validator.String{
								stringvalidator.OneOf("Email", "No Notify"),
							},
						},
						"notify_input_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The notification input type (`Primary Email`, `Alternate Email`, `Primary Alternate Email`).",
							Validators: []validator.String{
								stringvalidator.OneOf("Primary Email", "Alternate Email", "Primary Alternate Email"),
							},
						},
						"notify_recurring_report": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							MarkdownDescription: "Whether to notify for recurring reports.",
						},
					},
				},
			},
			"change_password": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether the user must change password on first login.",
			},
		},
	}
}

// Create handles the creation of the resource.
func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Validate and get timezone object
	var timeZone *client.TimeZone
	if !plan.TimeZone.IsNull() && plan.TimeZone.ValueString() != "" {
		tz, err := r.apiClient.ValidateTimezone(plan.TimeZone.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Timezone",
				fmt.Sprintf("The specified timezone '%s' is not valid: %s", plan.TimeZone.ValueString(), err),
			)

			return
		}
		timeZone = tz
	}

	// Build roles list
	var roles []client.UserRoleRef
	for _, roleName := range plan.Roles {
		roles = append(roles, client.UserRoleRef{UniqueId: roleName.ValueString()})
	}

	// Build user notifications list
	var userNotifications []client.UserNotification
	for _, notif := range plan.UserNotifications {
		userNotifications = append(userNotifications, client.UserNotification{
			NotifyType:            notif.NotifyType.ValueString(),
			NotifyMethod:          notif.NotifyMethod.ValueString(),
			NotifyInputType:       notif.NotifyInputType.ValueString(),
			NotifyRecurringReport: notif.NotifyRecurringReport.ValueBool(),
		})
	}

	createUser := client.CreateUser{
		LoginName:         plan.LoginName.ValueString(),
		Password:          plan.Password.ValueString(),
		FirstName:         plan.FirstName.ValueString(),
		LastName:          plan.LastName.ValueString(),
		Designation:       plan.Designation.ValueString(),
		Address:           plan.Address.ValueString(),
		City:              plan.City.ValueString(),
		State:             plan.State.ValueString(),
		Zip:               plan.Zip.ValueString(),
		Country:           plan.Country.ValueString(),
		Email:             plan.Email.ValueString(),
		AltEmail:          plan.AltEmail.ValueString(),
		PhoneNumber:       plan.PhoneNumber.ValueString(),
		MobileNumber:      plan.MobileNumber.ValueString(),
		TimeZone:          timeZone,
		AuthType:          plan.AuthType.ValueString(),
		Roles:             roles,
		UserNotifications: userNotifications,
		ChangePassword:    plan.ChangePassword.ValueBool(),
	}

	created, err := r.apiClient.CreateUser(tenantId, createUser)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())

		return
	}

	// Only update truly computed values (id) - preserve plan values for others
	plan.Id = types.StringValue(created.Id)
	plan.Status = types.StringValue(created.Status)

	ugElems := make([]attr.Value, len(created.UserGroups))
	for i, group := range created.UserGroups {
		ugElems[i] = types.StringValue(group.UniqueId)
	}
	plan.UserGroups, diags = types.SetValue(types.StringType, ugElems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetUser(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}

	// Update state from API response
	state.Id = types.StringValue(existing.Id)
	state.LoginName = types.StringValue(existing.LoginName)
	state.FirstName = types.StringValue(existing.FirstName)
	state.LastName = types.StringValue(existing.LastName)
	state.Designation = optionalStringValue(existing.Designation)
	state.Address = optionalStringValue(existing.Address)
	state.City = optionalStringValue(existing.City)
	state.State = optionalStringValue(existing.State)
	state.Zip = optionalStringValue(existing.Zip)
	state.Country = optionalStringValue(existing.Country)
	state.Email = types.StringValue(existing.Email)
	state.AltEmail = optionalStringValue(existing.AltEmail)
	state.PhoneNumber = optionalStringValue(existing.PhoneNumber)
	state.MobileNumber = optionalStringValue(existing.MobileNumber)
	if existing.TimeZone != nil {
		state.TimeZone = types.StringValue(existing.TimeZone.Name)
	}
	state.AuthType = types.StringValue(existing.AuthType)
	// Normalize status to Title Case to match schema default.
	// A Caser is not safe for concurrent use, so build one per call rather than
	// sharing a package-level instance.
	if existing.Status != "" {
		state.Status = types.StringValue(
			cases.Title(language.Und).String(strings.ToLower(existing.Status)),
		)
	}

	// Update roles - initialize to empty slice to avoid null vs empty diff
	// Remove user group roles from this list since those are managed via the user_group resource
	roles := make([]types.String, 0, len(existing.Roles))
	for _, role := range existing.Roles {
		if !role.UserGroupRole {
			roles = append(roles, types.StringValue(role.UniqueId))
		}
	}
	state.Roles = roles

	// Update user_groups from API response (computed field)
	ugElems := make([]attr.Value, len(existing.UserGroups))
	for i, group := range existing.UserGroups {
		ugElems[i] = types.StringValue(group.UniqueId)
	}
	state.UserGroups, diags = types.SetValue(types.StringType, ugElems)
	resp.Diagnostics.Append(diags...)

	// Preserve user_notifications from state - don't overwrite with API values
	// as the API may return additional default notifications

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Validate and get timezone object
	var timeZone *client.TimeZone
	if !plan.TimeZone.IsNull() && plan.TimeZone.ValueString() != "" {
		tz, err := r.apiClient.ValidateTimezone(plan.TimeZone.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Timezone",
				fmt.Sprintf("The specified timezone '%s' is not valid: %s", plan.TimeZone.ValueString(), err),
			)

			return
		}
		timeZone = tz
	}

	// Build roles list
	var roles []client.UserRoleRef
	for _, roleName := range plan.Roles {
		roles = append(roles, client.UserRoleRef{UniqueId: roleName.ValueString()})
	}

	// Build user notifications list
	var userNotifications []client.UserNotification
	for _, notif := range plan.UserNotifications {
		userNotifications = append(userNotifications, client.UserNotification{
			NotifyType:            notif.NotifyType.ValueString(),
			NotifyMethod:          notif.NotifyMethod.ValueString(),
			NotifyInputType:       notif.NotifyInputType.ValueString(),
			NotifyRecurringReport: notif.NotifyRecurringReport.ValueBool(),
		})
	}

	updateUser := client.UpdateUser{
		FirstName:         plan.FirstName.ValueString(),
		LastName:          plan.LastName.ValueString(),
		Designation:       plan.Designation.ValueString(),
		Address:           plan.Address.ValueString(),
		City:              plan.City.ValueString(),
		State:             plan.State.ValueString(),
		Zip:               plan.Zip.ValueString(),
		Country:           plan.Country.ValueString(),
		Email:             plan.Email.ValueString(),
		AltEmail:          plan.AltEmail.ValueString(),
		PhoneNumber:       plan.PhoneNumber.ValueString(),
		MobileNumber:      plan.MobileNumber.ValueString(),
		TimeZone:          timeZone,
		AuthType:          plan.AuthType.ValueString(),
		Roles:             roles,
		UserNotifications: userNotifications,
	}

	updated, err := r.apiClient.UpdateUser(tenantId, state.Id.ValueString(), updateUser)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating User",
			fmt.Sprintf("Could not update user: %s", err),
		)

		return
	}

	plan.Id = state.Id
	plan.Status = types.StringValue(updated.Status)

	ugElems := make([]attr.Value, len(updated.UserGroups))
	for i, group := range updated.UserGroups {
		ugElems[i] = types.StringValue(group.UniqueId)
	}
	plan.UserGroups, diags = types.SetValue(types.StringType, ugElems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteUser(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())

		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: userId (uses provider's tenant)
	userId := req.ID

	existing, err := r.apiClient.GetUser(r.apiClient.TenantId, userId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing User",
			fmt.Sprintf("Could not import user: %s", err),
		)

		return
	}

	// Build roles list
	var roles []types.String
	for _, role := range existing.Roles {
		roles = append(roles, types.StringValue(role.UniqueId))
	}

	// Build user groups list
	ugElems := make([]attr.Value, len(existing.UserGroups))
	for i, group := range existing.UserGroups {
		ugElems[i] = types.StringValue(group.UniqueId)
	}
	groups, ugDiags := types.SetValue(types.StringType, ugElems)
	resp.Diagnostics.Append(ugDiags...)

	// Get timezone name
	timeZoneName := ""
	if existing.TimeZone != nil {
		timeZoneName = existing.TimeZone.Name
	}

	state := UserModel{
		Id:                types.StringValue(existing.Id),
		LoginName:         types.StringValue(existing.LoginName),
		FirstName:         types.StringValue(existing.FirstName),
		LastName:          types.StringValue(existing.LastName),
		Email:             types.StringValue(existing.Email),
		PhoneNumber:       optionalStringValue(existing.PhoneNumber),
		MobileNumber:      optionalStringValue(existing.MobileNumber),
		TimeZone:          types.StringValue(timeZoneName),
		AuthType:          types.StringValue(existing.AuthType),
		Status:            types.StringValue(existing.Status),
		Roles:             roles,
		UserGroups:        groups,
		UserNotifications: convertUserNotificationsToModel(existing.UserNotifications),
	}

	resp.State.Set(ctx, &state)
}

// convertUserNotificationsToModel converts API user notifications to Terraform model
func convertUserNotificationsToModel(notifications []client.UserNotification) []UserNotificationModel {
	var result []UserNotificationModel
	for _, n := range notifications {
		notif := UserNotificationModel{
			NotifyType:            types.StringValue(n.NotifyType),
			NotifyMethod:          types.StringValue(n.NotifyMethod),
			NotifyRecurringReport: types.BoolValue(n.NotifyRecurringReport),
		}
		// Handle empty string as null for NotifyInputType
		if n.NotifyInputType != "" {
			notif.NotifyInputType = types.StringValue(n.NotifyInputType)
		} else {
			notif.NotifyInputType = types.StringNull()
		}
		result = append(result, notif)
	}

	return result
}
