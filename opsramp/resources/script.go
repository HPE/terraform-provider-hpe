// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/HPE/terraform-provider-hpe/opsramp/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &ScriptResource{}
	_ resource.ResourceWithImportState = &ScriptResource{}
)

// ScriptResource defines the resource implementation.
type ScriptResource struct {
	BaseResource
}

// ScriptParameterModel represents a single script parameter in Terraform state.
type ScriptParameterModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	DefaultValue types.String `tfsdk:"default_value"`
	Type         types.String `tfsdk:"type"`
	DataType     types.String `tfsdk:"data_type"`
}

// ScriptAttachmentModel represents the script file attachment in Terraform state.
type ScriptAttachmentModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	File types.String `tfsdk:"file"`
}

// ScriptModel maps Terraform schema attributes to the provider model.
type ScriptModel struct {
	Client         types.String           `tfsdk:"client"`
	Uuid           types.String           `tfsdk:"uuid"`
	CategoryId     types.String           `tfsdk:"category_id"`
	Name           types.String           `tfsdk:"name"`
	Description    types.String           `tfsdk:"description"`
	Platforms      []types.String         `tfsdk:"platforms"`
	Parameters     []ScriptParameterModel `tfsdk:"parameters"`
	ExecutionType  types.String           `tfsdk:"execution_type"`
	InstallTimeout types.Int64            `tfsdk:"install_timeout"`
	Attachment     *ScriptAttachmentModel `tfsdk:"attachment"`
	ScriptVersion  types.String           `tfsdk:"script_version"`
}

// NewScript creates a new instance of the resource.
func NewScript() resource.Resource {
	return &ScriptResource{}
}

// Metadata returns the resource type name.
func (r *ScriptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_script"
}

// Schema defines the schema for the resource.
func (r *ScriptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp RBA Script resource. Scripts belong to RBA categories and define automation actions that can be run on managed resources.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this script should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The UUID identifier of the script (as a string).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"category_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Script Category this script belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the script.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A description of the script.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"platforms": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The platforms this script targets. Valid values: `WINDOWS`, `LINUX`.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("WINDOWS", "LINUX")),
				},
			},
			"execution_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The execution type of the script (e.g., `SHELL`, `POWERSHELL`, `PYTHON`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("SHELL", "POWERSHELL", "PYTHON"),
				},
			},
			"install_timeout": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Timeout in seconds for script execution.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"script_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The version of the script, assigned by the API.",
			},
			"parameters": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Input parameters accepted by the script.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the parameter, assigned by the API.",
						},
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The name of the parameter.",
						},
						"description": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "A description of the parameter.",
						},
						"default_value": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "The default value for the parameter.",
						},
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Whether the parameter is `REQUIRED` or `OPTIONAL`.",
							Validators: []validator.String{
								stringvalidator.OneOf("REQUIRED", "OPTIONAL"),
							},
						},
						"data_type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The data type of the parameter. Valid values: `STRING`, `INTEGER`, `PASSWORD`.",
							Validators: []validator.String{
								stringvalidator.OneOf("STRING", "INTEGER", "PASSWORD"),
							},
						},
					},
				},
			},
			"attachment": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The script file attachment. Provide `name` and `content_url` (base64-encoded script content) when creating a script with a file attachment.",
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The attachment ID assigned by the API.",
					},
					"name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The filename of the attachment.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"file": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Script content",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
		},
	}
}

// resolveTenantId determines the effective tenant ID.
func (r *ScriptResource) resolveTenantId(clientAttr types.String) string {
	if !clientAttr.IsNull() && clientAttr.ValueString() != "" {
		return clientAttr.ValueString()
	}

	return r.apiClient.TenantId
}

// buildScriptRequest converts the Terraform model into an API request body.
func buildScriptRequest(plan ScriptModel) client.Script {
	s := client.Script{
		Name:          plan.Name.ValueString(),
		Description:   plan.Description.ValueString(),
		ExecutionType: plan.ExecutionType.ValueString(),
	}

	if !plan.InstallTimeout.IsNull() && !plan.InstallTimeout.IsUnknown() {
		s.InstallTimeout = int(plan.InstallTimeout.ValueInt64())
	}

	for _, p := range plan.Platforms {
		s.Platforms = append(s.Platforms, p.ValueString())
	}

	for _, param := range plan.Parameters {
		sp := client.ScriptParameter{
			Name:         param.Name.ValueString(),
			Description:  param.Description.ValueString(),
			DefaultValue: param.DefaultValue.ValueString(),
			Type:         param.Type.ValueString(),
			DataType:     param.DataType.ValueString(),
		}
		if !param.Id.IsNull() && !param.Id.IsUnknown() && param.Id.ValueString() != "" {
			if numId, err := strconv.Atoi(param.Id.ValueString()); err == nil {
				sp.Id = numId
			}
		}
		s.Parameters = append(s.Parameters, sp)
	}

	s.Attachment = &client.ScriptAttachment{
		Name: plan.Attachment.Name.ValueString(),
		File: base64.StdEncoding.EncodeToString([]byte(plan.Attachment.File.ValueString())),
	}

	return s
}

// populateModelFromAPI maps an API response back into the Terraform model.
func populateModelFromAPI(model *ScriptModel, s *client.Script) {
	model.Uuid = types.StringValue(s.Uuid)
	model.Name = types.StringValue(s.Name)
	model.Description = types.StringValue(s.Description)
	model.ExecutionType = types.StringValue(s.ExecutionType)
	model.InstallTimeout = types.Int64Value(int64(s.InstallTimeout))
	model.ScriptVersion = types.StringValue(s.ScriptVersion)

	var platforms []types.String
	for _, p := range s.Platforms {
		platforms = append(platforms, types.StringValue(p))
	}
	model.Platforms = platforms

	var parameters []ScriptParameterModel
	for _, p := range s.Parameters {
		parameters = append(parameters, ScriptParameterModel{
			Id:           types.StringValue(strconv.Itoa(p.Id)),
			Name:         types.StringValue(p.Name),
			Description:  types.StringValue(p.Description),
			DefaultValue: types.StringValue(p.DefaultValue),
			Type:         types.StringValue(p.Type),
			DataType:     types.StringValue(p.DataType),
		})
	}
	model.Parameters = parameters

	model.Attachment = &ScriptAttachmentModel{
		Id:   types.Int64Value(int64(s.Attachment.Id)),
		Name: types.StringValue(s.Attachment.Name),
		File: types.StringValue(s.Attachment.File),
	}
}

// Create handles the creation of the resource.
func (r *ScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ScriptModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)
	categoryId := plan.CategoryId.ValueString()

	scriptReq := buildScriptRequest(plan)
	created, err := r.apiClient.CreateScript(tenantId, categoryId, scriptReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Script", fmt.Sprintf("Could not create script: %s", err))

		return
	}

	populateModelFromAPI(&plan, created)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *ScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ScriptModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)
	categoryId := state.CategoryId.ValueString()
	scriptId := state.Uuid.ValueString()

	existing, err := r.apiClient.GetScript(tenantId, categoryId, scriptId)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)

			return
		}
		resp.Diagnostics.AddError("Error Reading Script", fmt.Sprintf("Could not read script %s: %s", scriptId, err))

		return
	}

	populateModelFromAPI(&state, existing)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *ScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ScriptModel
	var state ScriptModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)
	categoryId := state.CategoryId.ValueString()
	scriptId := state.Uuid.ValueString()

	// If description is unknown/null in plan, keep the current backend/state value
	// so update requests do not accidentally submit an empty description.
	if plan.Description.IsUnknown() || plan.Description.IsNull() {
		plan.Description = state.Description
	}

	scriptReq := buildScriptRequest(plan)
	updated, err := r.apiClient.UpdateScript(tenantId, categoryId, scriptId, scriptReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Script", fmt.Sprintf("Could not update script %s: %s", scriptId, err))

		return
	}

	populateModelFromAPI(&plan, updated)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *ScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ScriptModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)
	categoryId := state.CategoryId.ValueString()
	scriptId := state.Uuid.ValueString()

	err := r.apiClient.DeleteScript(tenantId, categoryId, scriptId)
	if err != nil && !strings.Contains(err.Error(), "status: 204") {
		resp.Diagnostics.AddError("Error Deleting Script", fmt.Sprintf("Could not delete script %s: %s", scriptId, err))

		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
// Import ID format: <category_id>:<script_id> or <client_id>:<category_id>:<script_id> (MSP only)
func (r *ScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parsed, err := r.ParseImportID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())

		return
	}

	tenantId := r.TenantForImport(parsed)
	categoryId := parsed.Parts[0]
	scriptId := parsed.Parts[1]

	existing, err := r.apiClient.GetScript(tenantId, categoryId, scriptId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Script",
			fmt.Sprintf("Could not import script: %s", err),
		)

		return
	}

	var state ScriptModel
	state.Client = parsed.Client
	state.CategoryId = types.StringValue(categoryId)

	populateModelFromAPI(&state, existing)

	resp.State.Set(ctx, &state)
}
