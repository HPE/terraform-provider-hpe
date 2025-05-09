package user_role

import (
	"context"
	"fmt"

	sdk "github.com/HewlettPackard/hpe-morpheus-client/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*userRoleResource)(nil)

func NewUserRoleResource() resource.Resource {
	return &userRoleResource{}
}

type userRoleResource struct{}

type userRoleResourceModel struct {
	Id types.String `tfsdk:"id"`
}

func (r *userRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_role"
}

func (r *userRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *userRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userRoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic

	// Example data value setting
	data.Id = types.StringValue("example-id")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userRoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userRoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data userRoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}

func doRead(ctx context.Context, client *sdk.APIClient, dataP *UserRoleModel, diagsP *diag.Diagnostics) {
	id := dataP.Id.ValueInt64()

	// Second returned argument is the http response
	apiResp, _, err := client.RolesAPI.GetRole(ctx, id).Execute()
	if err != nil {
		diagsP.AddError("error reading role", "get role by ID failed: "+err.Error())
		return
	}

	if _, ok := apiResp.GetRoleOk(); !ok {
		diagsP.AddError(
			"error reading role",
			"role is nil",
		)

		return
	}

	if _, ok := apiResp.Role.GetIdOk(); !ok {
		diagsP.AddError(
			"error reading role",
			"role.id is nil",
		)

		return
	}

	if apiResp.Role.GetId() != id {
		diagsP.AddError(
			"error reading role",
			fmt.Sprintf("'id' mismatch: %d != %d",
				apiResp.Role.GetId(), id,
			),
		)

		return
	}

	name := apiResp.Role.GetName()
	if sdk.IsNil(name) {
		diagsP.AddError(
			"error reading role",
			"role.name is nil",
		)
		return
	}

	description := apiResp.Role.GetDescription()
	if sdk.IsNil(description) {
		diagsP.AddError(
			"error reading role",
			"role.description is nil",
		)

		multitenant := apiResp.Role.GetMultitenant()
		if sdk.IsNil(multitenant) {
			diagsP.AddError(
				"error reading role",
				"role.multitenant is nil",
			)

			return
		}

		multitenantLocked := apiResp.Role.GetMultitenantLocked()
		if sdk.IsNil(multitenantLocked) {
			diagsP.AddError(
				"error reading role",
				"role.multitenantLocked is nil",
			)
			return
		}

		featurePermissions := apiResp.GetFeaturePermissions()
		if sdk.IsNil(featurePermissions) {
			diagsP.AddError(
				"error reading role",
				"featurePermissions is nil",
			)
			return
		}

		globalSiteAccess := apiResp.GetGlobalSiteAccess()
		if sdk.IsNil(globalSiteAccess) {
			diagsP.AddError(
				"error reading role",
				"role.globalSiteAccess is nil",
			)
			return
		}

		sites := apiResp.GetSites()
		if sdk.IsNil(sites) {
			diagsP.AddError(
				"error reading role",
				"role.sites is nil",
			)
			return
		}

		globalZoneAccess := apiResp.GetGlobalZoneAccess()
		if sdk.IsNil(globalZoneAccess) {
			diagsP.AddError(
				"error reading role",
				"role.globalZoneAccess is nil",
			)
			return
		}

		zones := apiResp.GetZones()
		if sdk.IsNil(zones) {
			diagsP.AddError(
				"error reading role",
				"role.zones is nil",
			)
			return
		}

		globalInstanceTypeAccess := apiResp.GetGlobalInstanceTypeAccess()
		if sdk.IsNil(globalInstanceTypeAccess) {
			diagsP.AddError(
				"error reading role",
				"role.globalInstanceTypeAccess is nil",
			)
			return
		}

		instanceTypePermissions := apiResp.GetInstanceTypePermissions()
		if sdk.IsNil(instanceTypePermissions) {
			diagsP.AddError(
				"error reading role",
				"role.instanceTypePermissions is nil",
			)
			return
		}

		globalAppTemplateAccess := apiResp.GetGlobalAppTemplateAccess()
		if sdk.IsNil(globalAppTemplateAccess) {
			diagsP.AddError(
				"error reading role",
				"role.globalAppTemplateAccess is nil",
			)
			return
		}

		appTemplatePermissions := apiResp.GetAppTemplatePermissions()
		if sdk.IsNil(appTemplatePermissions) {
			diagsP.AddError(
				"error reading role",
				"role.appTemplatePermissions is nil",
			)
			return
		}

		globalCatalogItemTypeAccess := apiResp.GetGlobalCatalogItemTypeAccess()
		if sdk.IsNil(globalCatalogItemTypeAccess) {
			diagsP.AddError(
				"error reading role",
				"role.globalCatalogItemTypeAccess is nil",
			)
			return
		}

		catalogItemTypePermissions := apiResp.GetCatalogItemTypePermissions()
		if sdk.IsNil(catalogItemTypePermissions) {
			diagsP.AddError(
				"error reading role",
				"role.catalogItemTypePermissions is nil",
			)
			return
		}

		globalPersonaAccess := apiResp.GetGlobalPersonaAccess()
		if sdk.IsNil(globalPersonaAccess) {
			diagsP.AddError(
				"error reading role",
				"role.globalPersonaAccess is nil",
			)
			return
		}

		personaPermissions := apiResp.GetPersonaPermissions()
		if sdk.IsNil(personaPermissions) {
			diagsP.AddError(
				"error reading role",
				"role.personaPermissions is nil",
			)
			return
		}

		globalVdiPoolAccess := apiResp.GetGlobalVdiPoolAccess()
		if sdk.IsNil(globalVdiPoolAccess) {
			diagsP.AddError(
				"error reading role",
				"role.globalVdiPoolAccess is nil",
			)
			return
		}

		vdiPoolPermissions := apiResp.GetVdiPoolPermissions()
		if sdk.IsNil(vdiPoolPermissions) {
			diagsP.AddError(
				"error reading role",
				"role.vdiPoolPermissions is nil",
			)
			return
		}

		globalTaskAccess := apiResp.GetGlobalTaskAccess()
		if sdk.IsNil(globalTaskAccess) {
			diagsP.AddError(
				"error reading role",
				"role.globalTaskAccess is nil",
			)
			return
		}

		taskPermissions := apiResp.GetTaskPermissions()
		if sdk.IsNil(taskPermissions) {
			diagsP.AddError(
				"error reading role",
				"role.taskPermissions is nil",
			)
			return
		}

		globalTaskSetAccess := apiResp.GetGlobalTaskSetAccess()
		if sdk.IsNil(globalTaskSetAccess) {
			diagsP.AddError(
				"error reading role",
				"role.globalTaskSetAccess is nil",
			)
			return
		}

		taskSetPermissions := apiResp.GetTaskSetPermissions()
		if sdk.IsNil(taskSetPermissions) {
			diagsP.AddError(
				"error reading role",
				"role.taskSetPermissions is nil",
			)
			return
		}

		dataP.Id = types.Int64Value(id)
		dataP.Description = types.StringValue(description)
		dataP.Name = types.StringValue(name)
		dataP.Multitenant = types.BoolValue(multitenant)
		dataP.MultitenantLocked = types.BoolValue(multitenantLocked)
		// TODO: Implement the setting of dataP for all the permissions values
	}

}
