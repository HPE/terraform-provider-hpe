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
)

//nolint:unused
const summary = "role permissions data source"

type permissions sdk.AddRolesRequestRole

// custom JSON marshaler override to ignore required authority field
// while still using generated SDK POST structs for permissions
func (p *permissions) MarshalJSON() ([]byte, error) {
	res := make(map[string]any)

	if len(p.FeaturePermissions) > 0 {
		res["featurePermissions"] = p.FeaturePermissions
	}

	if p.GlobalSiteAccess != nil && *p.GlobalSiteAccess != "" {
		res["globalSiteAccess"] = *p.GlobalSiteAccess
	}

	if len(p.Sites) > 0 {
		res["sites"] = p.Sites
	}

	if p.GlobalZoneAccess != nil && *p.GlobalZoneAccess != "" {
		res["globalZoneAccess"] = *p.GlobalZoneAccess
	}

	if len(p.Zones) > 0 {
		res["zones"] = p.Zones
	}

	if p.GlobalInstanceTypeAccess != nil && *p.GlobalInstanceTypeAccess != "" {
		res["globalInstanceTypeAccess"] = *p.GlobalInstanceTypeAccess
	}

	if len(p.InstanceTypePermissions) > 0 {
		res["instanceTypePermissions"] = p.InstanceTypePermissions
	}

	if p.GlobalAppTemplateAccess != nil && *p.GlobalAppTemplateAccess != "" {
		res["globalAppTemplateAccess"] = *p.GlobalAppTemplateAccess
	}

	if len(p.AppTemplatePermissions) > 0 {
		res["appTemplatePermissions"] = p.AppTemplatePermissions
	}

	if p.GlobalCatalogItemTypeAccess != nil && *p.GlobalCatalogItemTypeAccess != "" {
		res["globalCatalogItemTypeAccess"] = *p.GlobalCatalogItemTypeAccess
	}

	if len(p.CatalogItemTypePermissions) > 0 {
		res["catalogItemTypePermissions"] = p.CatalogItemTypePermissions
	}

	if p.GlobalPersonaAccess != nil && *p.GlobalPersonaAccess != "" {
		res["globalPersonaAccess"] = *p.GlobalPersonaAccess
	}

	if len(p.PersonaPermissions) > 0 {
		res["personaPermissions"] = p.PersonaPermissions
	}

	if p.GlobalVdiPoolAccess != nil && *p.GlobalVdiPoolAccess != "" {
		res["globalVdiPoolAccess"] = *p.GlobalVdiPoolAccess
	}

	if len(p.VdiPoolPermissions) > 0 {
		res["vdiPoolPermissions"] = p.VdiPoolPermissions
	}

	if p.GlobalReportTypeAccess != nil && *p.GlobalReportTypeAccess != "" {
		res["globalReportTypeAccess"] = *p.GlobalReportTypeAccess
	}

	if len(p.ReportTypePermissions) > 0 {
		res["reportTypePermissions"] = p.ReportTypePermissions
	}

	if p.GlobalTaskAccess != nil && *p.GlobalTaskAccess != "" {
		res["globalTaskAccess"] = *p.GlobalTaskAccess
	}

	if len(p.TaskPermissions) > 0 {
		res["taskPermissions"] = p.TaskPermissions
	}

	if p.GlobalTaskSetAccess != nil && *p.GlobalTaskSetAccess != "" {
		res["globalTaskSetAccess"] = *p.GlobalTaskSetAccess
	}

	if len(p.TaskSetPermissions) > 0 {
		res["taskSetPermissions"] = p.TaskSetPermissions
	}

	// TODO: Add Cluster Permissions support (not yet documented in OpenAPI spec)

	return json.Marshal(res)
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &DataSource{}
)

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
