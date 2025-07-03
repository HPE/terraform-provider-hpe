// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package rolepermissions

import (
	"context"
	"encoding/json"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
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
		var fpInners []sdk.AddRolesRequestRoleFeaturePermissionsInner
		featurePermissions := data.FeaturePermissions.String()
		err := json.Unmarshal([]byte(featurePermissions), &fpInners)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to unmarshal feature_permissions to sdk struct",
				err.Error(),
			)

			return
		}

		permissionsStruct.FeaturePermissions = fpInners

	}

	if !data.DefaultGroupAccess.IsNull() && !data.DefaultGroupAccess.IsUnknown() {
		v, err := convert.ValueToAny(ctx, data.DefaultGroupAccess)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to convert default_group_access to any",
				err.Error(),
			)

			return
		}

		defaultGroupAccess, ok := v.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"failed to convert default_group_access to string",
				"type assertion failed",
			)

			return

		}

		permissionsStruct.GlobalSiteAccess = &defaultGroupAccess
	}

	if !data.DefaultCloudAccess.IsNull() && !data.DefaultCloudAccess.IsUnknown() {
		v, err := convert.ValueToAny(ctx, data.DefaultCloudAccess)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to convert default_cloud_access to any",
				err.Error(),
			)

			return
		}

		defaultCloudAccess, ok := v.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"failed to convert default_cloud_access to string",
				"type assertion failed",
			)

			return
		}

		permissionsStruct.GlobalZoneAccess = &defaultCloudAccess
	}

	if !data.DefaultBlueprintAccess.IsNull() && !data.DefaultBlueprintAccess.IsUnknown() {
		v, err := convert.ValueToAny(ctx, data.DefaultBlueprintAccess)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to convert default_blueprint_access to any",
				err.Error(),
			)

			return
		}

		defaultBlueprintAccess, ok := v.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"failed to convert default_blueprint_access to string",
				"type assertion failed",
			)

			return
		}

		permissionsStruct.GlobalAppTemplateAccess = &defaultBlueprintAccess
	}

	if !data.DefaultCatalogItemTypeAccess.IsNull() && !data.DefaultCatalogItemTypeAccess.IsUnknown() {
		v, err := convert.ValueToAny(ctx, data.DefaultCatalogItemTypeAccess)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to convert default_catalog_item_type_access to any",
				err.Error(),
			)

			return
		}

		defaultCatalogItemTypeAccess, ok := v.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"failed to convert default_catalog_item_type_access to string",
				"type assertion failed",
			)

			return
		}

		permissionsStruct.GlobalCatalogItemTypeAccess = &defaultCatalogItemTypeAccess
	}

	if !data.DefaultInstanceTypeAccess.IsNull() && !data.DefaultInstanceTypeAccess.IsUnknown() {
		v, err := convert.ValueToAny(ctx, data.DefaultInstanceTypeAccess)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to convert default_instance_type_access to any",
				err.Error(),
			)

			return
		}

		defaultInstanceTypeAccess, ok := v.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"failed to convert default_instance_type_access to string",
				"type assertion failed",
			)

			return
		}

		permissionsStruct.GlobalInstanceTypeAccess = &defaultInstanceTypeAccess
	}

	if !data.DefaultPersonaAccess.IsNull() && !data.DefaultPersonaAccess.IsUnknown() {
		v, err := convert.ValueToAny(ctx, data.DefaultPersonaAccess)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to convert default_persona_access to any",
				err.Error(),
			)

			return
		}

		defaultPersonaAccess, ok := v.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"failed to convert default_persona_access to string",
				"type assertion failed",
			)

			return
		}

		permissionsStruct.GlobalPersonaAccess = &defaultPersonaAccess
	}

	if !data.DefaultReportTypeAccess.IsNull() && !data.DefaultReportTypeAccess.IsUnknown() {
		v, err := convert.ValueToAny(ctx, data.DefaultReportTypeAccess)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to convert default_report_type_access to any",
				err.Error(),
			)

			return
		}

		defaultReportTypeAccess, ok := v.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"failed to convert default_report_type_access to string",
				"type assertion failed",
			)

			return
		}

		permissionsStruct.GlobalReportTypeAccess = &defaultReportTypeAccess
	}

	if !data.DefaultTaskAccess.IsNull() && !data.DefaultTaskAccess.IsUnknown() {
		v, err := convert.ValueToAny(ctx, data.DefaultTaskAccess)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to convert default_task_access to any",
				err.Error(),
			)

			return
		}

		defaultTaskAccess, ok := v.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"failed to convert default_task_access to string",
				"type assertion failed",
			)

			return
		}

		permissionsStruct.GlobalTaskAccess = &defaultTaskAccess
	}

	if !data.DefaultWorkflowAccess.IsNull() && !data.DefaultWorkflowAccess.IsUnknown() {
		v, err := convert.ValueToAny(ctx, data.DefaultWorkflowAccess)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to convert default_workflow_access to any",
				err.Error(),
			)
			return
		}
		defaultWorkflowAccess, ok := v.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"failed to convert default_workflow_access to string",
				"type assertion failed",
			)
			return
		}
		permissionsStruct.GlobalTaskSetAccess = &defaultWorkflowAccess
	}

	if !data.DefaultVdiPoolAccess.IsNull() && !data.DefaultVdiPoolAccess.IsUnknown() {
		v, err := convert.ValueToAny(ctx, data.DefaultVdiPoolAccess)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to convert default_vdi_pool_access to any",
				err.Error(),
			)

			return
		}

		defaultVdiPoolAccess, ok := v.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"failed to convert default_vdi_pool_access to string",
				"type assertion failed",
			)

			return
		}

		permissionsStruct.GlobalVdiPoolAccess = &defaultVdiPoolAccess
	}

	// TODO: Do the same as above for the other permissions fields

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
