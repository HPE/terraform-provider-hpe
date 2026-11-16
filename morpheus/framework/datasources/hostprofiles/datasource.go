// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package hostprofiles implements a data source listing host profiles.
package hostprofiles

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read host_profiles data source"

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

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
	resp.TypeName = req.ProviderTypeName + "_" + "host_profiles"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = HostProfilesDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config HostProfilesModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	listReq := client.HostProfilesAPI.ListHostProfiles(ctx)
	if !config.Max.IsNull() && !config.Max.IsUnknown() {
		listReq = listReq.Max(config.Max.ValueInt64())
	}

	if !config.Offset.IsNull() && !config.Offset.IsUnknown() {
		listReq = listReq.Offset(config.Offset.ValueInt64())
	}

	if !config.Phrase.IsNull() && !config.Phrase.IsUnknown() {
		listReq = listReq.Phrase(config.Phrase.ValueString())
	}

	r, httpResp, err := listReq.Execute()
	if r == nil || httpResp == nil || err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf("LIST host_profiles failed: %s", errfmt.ErrMsg(err, httpResp)),
		)

		return
	}

	elemType := HostProfilesValue{}.Type(ctx)

	elems := make([]HostProfilesValue, 0, len(r.HostProfiles))
	for i := range r.HostProfiles {
		hp := r.HostProfiles[i]
		elems = append(elems, HostProfilesValue{
			Code:             convert.StrToType(hp.Code.Get()),
			ComplianceStatus: convert.StrToType(hp.ComplianceStatus),
			Description:      convert.StrToType(hp.Description.Get()),
			Enabled:          convert.BoolToType(hp.Enabled),
			Id:               convert.Int64ToType(hp.Id),
			Name:             convert.StrToType(hp.Name),
			ServerCount:      convert.Int32ToType(hp.ServerCount),
			state:            attr.ValueStateKnown,
		})
	}

	var hostProfiles types.Set
	if len(elems) == 0 {
		hostProfiles = types.SetNull(elemType)
	} else {
		set, setDiags := types.SetValueFrom(ctx, elemType, elems)
		resp.Diagnostics.Append(setDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		hostProfiles = set
	}

	config.HostProfiles = hostProfiles

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
