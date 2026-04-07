// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                 = "read os type data source"
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
	ErrorNoOsTypeFound      = `no os type found`
	ErrorMultipleOsTypes    = `multiple os types were returned`
)

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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_os_type"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = OsTypeDataSourceSchema(ctx)
}

func osTypeAsState(
	ctx context.Context,
	osType *sdk.GetOsType200ResponseOsType,
) (OsTypeModel, error) {
	images, diags := convert.ToSetType(ctx, osType.Images, mapImage)
	if diags.HasError() {
		return OsTypeModel{}, fmt.Errorf("error creating images set")
	}

	return OsTypeModel{
		BitCount:         convert.Int64ToType(osType.BitCount),
		Category:         convert.StrToType(osType.Category.Get()),
		CloudInitVersion: convert.StrToType(osType.CloudInitVersion),
		Code:             convert.StrToType(osType.Code),
		Description:      convert.StrToType(osType.Description.Get()),
		Id:               convert.Int64ToType(osType.Id),
		Images:           images,
		InstallAgent:     convert.BoolToType(osType.InstallAgent.Get()),
		Name:             convert.StrToType(osType.Name),
		OsCodename:       convert.StrToType(osType.OsCodename.Get()),
		OsFamily:         convert.StrToType(osType.OsFamily.Get()),
		OsName:           convert.StrToType(osType.OsName.Get()),
		OsVersion:        convert.StrToType(osType.OsVersion.Get()),
		Owner:            convert.StrToType(osType.Owner.Get()),
		Platform:         convert.StrToType(osType.Platform),
		Vendor:           convert.StrToType(osType.Vendor.Get()),
	}, nil
}

func getOsTypeByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetOsType200ResponseOsType, error) {
	r, hresp, err := apiClient.LibraryAPI.GetOsType(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for os type %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	osType := r.GetOsType()

	return &osType, nil
}

func getOsTypeByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetOsType200ResponseOsType, error) {
	rs, hresp, err := apiClient.LibraryAPI.ListOsTypes(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for os type %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	var matched []sdk.ListOsTypes200ResponseAllOfOsTypesInner

	for _, o := range rs.OsTypes {
		if o.GetName() == name {
			matched = append(matched, o)
		}
	}

	if len(matched) == 0 {
		return nil, errors.New(ErrorNoOsTypeFound)
	} else if len(matched) > 1 {
		return nil, errors.New(ErrorMultipleOsTypes)
	}

	return getOsTypeByID(ctx, matched[0].GetId(), apiClient)
}

func getOsType(
	ctx context.Context,
	config *OsTypeModel,
	apiClient *sdk.APIClient,
) (*sdk.GetOsType200ResponseOsType, error) {
	if !config.Id.IsNull() {
		return getOsTypeByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getOsTypeByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config OsTypeModel

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	osType, err := getOsType(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state, err := osTypeAsState(ctx, osType)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
