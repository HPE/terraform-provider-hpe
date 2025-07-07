// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package rolepermissions

import (
	"context"
	"encoding/json"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
)

//nolint:unused
const summary = "role permissions data source"

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_role_permissions"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = RolePermissionsDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data RolePermissionsModel

	// Read config
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionsStruct := permissions{}

	if !data.FeaturePermissions.IsNull() && !data.FeaturePermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRoleFeaturePermissionsInner
		if err := json.Unmarshal([]byte(data.FeaturePermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal feature_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.FeaturePermissions = inners

	}

	if !data.CloudPermissions.IsNull() && !data.CloudPermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRoleZonesInner
		if err := json.Unmarshal([]byte(data.CloudPermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal cloud_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.Zones = inners
	}

	if !data.GroupPermissions.IsNull() && !data.GroupPermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRoleSitesInner
		if err := json.Unmarshal([]byte(data.GroupPermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal group_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.Sites = inners
	}

	if !data.BlueprintPermissions.IsNull() && !data.BlueprintPermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRoleAppTemplatePermissionsInner
		if err := json.Unmarshal([]byte(data.BlueprintPermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal blueprint_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.AppTemplatePermissions = inners
	}

	if !data.CatalogItemTypePermissions.IsNull() && !data.CatalogItemTypePermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRoleCatalogItemTypePermissionsInner
		if err := json.Unmarshal([]byte(data.CatalogItemTypePermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal catalog_item_type_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.CatalogItemTypePermissions = inners
	}

	if !data.InstanceTypePermissions.IsNull() && !data.InstanceTypePermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRoleInstanceTypePermissionsInner
		if err := json.Unmarshal([]byte(data.InstanceTypePermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal instance_type_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.InstanceTypePermissions = inners
	}

	if !data.PersonaPermissions.IsNull() && !data.PersonaPermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRolePersonaPermissionsInner
		if err := json.Unmarshal([]byte(data.PersonaPermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal persona_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.PersonaPermissions = inners
	}

	if !data.ReportTypePermissions.IsNull() && !data.ReportTypePermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRoleReportTypePermissionsInner
		if err := json.Unmarshal([]byte(data.ReportTypePermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal report_type_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.ReportTypePermissions = inners
	}

	if !data.TaskPermissions.IsNull() && !data.TaskPermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRoleTaskPermissionsInner
		if err := json.Unmarshal([]byte(data.TaskPermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal task_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.TaskPermissions = inners
	}

	if !data.WorkflowPermissions.IsNull() && !data.WorkflowPermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRoleTaskSetPermissionsInner
		if err := json.Unmarshal([]byte(data.WorkflowPermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal workflow_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.TaskSetPermissions = inners
	}

	if !data.VdiPoolPermissions.IsNull() && !data.VdiPoolPermissions.IsUnknown() {
		var inners []sdk.AddRolesRequestRoleVdiPoolPermissionsInner
		if err := json.Unmarshal([]byte(data.VdiPoolPermissions.String()), &inners); err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal vdi_pool_permissions to sdk struct",
				err.Error(),
			)

			return
		}
		permissionsStruct.VdiPoolPermissions = inners
	}

	if !data.DefaultGroupAccess.IsNull() && !data.DefaultGroupAccess.IsUnknown() {
		defaultGroupAccess := data.DefaultGroupAccess.ValueString()
		permissionsStruct.GlobalSiteAccess = &defaultGroupAccess
	}

	if !data.DefaultCloudAccess.IsNull() && !data.DefaultCloudAccess.IsUnknown() {
		defaultCloudAccess := data.DefaultCloudAccess.ValueString()
		permissionsStruct.GlobalZoneAccess = &defaultCloudAccess
	}

	if !data.DefaultBlueprintAccess.IsNull() && !data.DefaultBlueprintAccess.IsUnknown() {
		defaultBlueprintAccess := data.DefaultBlueprintAccess.ValueString()
		permissionsStruct.GlobalAppTemplateAccess = &defaultBlueprintAccess
	}

	if !data.DefaultCatalogItemTypeAccess.IsNull() && !data.DefaultCatalogItemTypeAccess.IsUnknown() {
		defaultCatalogItemTypeAccess := data.DefaultCatalogItemTypeAccess.ValueString()
		permissionsStruct.GlobalCatalogItemTypeAccess = &defaultCatalogItemTypeAccess
	}

	if !data.DefaultInstanceTypeAccess.IsNull() && !data.DefaultInstanceTypeAccess.IsUnknown() {
		defaultInstanceTypeAccess := data.DefaultInstanceTypeAccess.ValueString()
		permissionsStruct.GlobalInstanceTypeAccess = &defaultInstanceTypeAccess
	}

	if !data.DefaultPersonaAccess.IsNull() && !data.DefaultPersonaAccess.IsUnknown() {
		defaultPersonaAccess := data.DefaultPersonaAccess.ValueString()
		permissionsStruct.GlobalPersonaAccess = &defaultPersonaAccess
	}

	if !data.DefaultReportTypeAccess.IsNull() && !data.DefaultReportTypeAccess.IsUnknown() {
		defaultReportTypeAccess := data.DefaultReportTypeAccess.ValueString()
		permissionsStruct.GlobalReportTypeAccess = &defaultReportTypeAccess
	}

	if !data.DefaultTaskAccess.IsNull() && !data.DefaultTaskAccess.IsUnknown() {
		defaultTaskAccess := data.DefaultTaskAccess.ValueString()
		permissionsStruct.GlobalTaskAccess = &defaultTaskAccess
	}

	if !data.DefaultWorkflowAccess.IsNull() && !data.DefaultWorkflowAccess.IsUnknown() {
		defaultWorkflowAccess := data.DefaultWorkflowAccess.ValueString()
		permissionsStruct.GlobalTaskSetAccess = &defaultWorkflowAccess
	}

	if !data.DefaultVdiPoolAccess.IsNull() && !data.DefaultVdiPoolAccess.IsUnknown() {
		defaultVdiPoolAccess := data.DefaultVdiPoolAccess.ValueString()
		permissionsStruct.GlobalVdiPoolAccess = &defaultVdiPoolAccess
	}

	// marshal the permissions struct to JSON
	b, err := json.Marshal(&permissionsStruct)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to marshal sdk AddRole struct to json",
			err.Error(),
		)

		return
	}

	jsonBody := string(b)

	diags = resp.State.SetAttribute(ctx, path.Root("json"), jsonBody)
	resp.Diagnostics.Append(diags...)
}
