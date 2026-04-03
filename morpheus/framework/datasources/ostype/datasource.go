// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

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

func mapImages(
	ctx context.Context,
	images []sdk.GetOsType200ResponseOsTypeImagesInner,
) (types.Set, error) {
	if len(images) == 0 {
		return types.SetNull(ImagesValue{}.Type(ctx)), nil
	}

	var vals []attr.Value

	for _, img := range images {
		v := ImagesValue{
			Account:          convert.Int64ToType(img.Account.Get()),
			ComputeZoneType:  convert.Int64ToType(img.ComputeZoneType.Get()),
			Id:               convert.Int64ToType(img.Id),
			ProvisionType:    convert.Int64ToType(img.ProvisionType.Get()),
			VirtualImageId:   convert.Int64ToType(img.VirtualImageId),
			VirtualImageName: convert.StrToType(img.VirtualImageName),
			Zone:             convert.Int64ToType(img.Zone.Get()),
			state:            attr.ValueStateKnown,
		}

		vals = append(vals, v)
	}

	set, diags := types.SetValue(ImagesValue{}.Type(ctx), vals)
	if diags.HasError() {
		return types.SetNull(ImagesValue{}.Type(ctx)), fmt.Errorf("error creating images set")
	}

	return set, nil
}

func osTypeAsState(
	ctx context.Context,
	osType *sdk.GetOsType200ResponseOsType,
) (OsTypeModel, error) {
	images, err := mapImages(ctx, osType.Images)
	if err != nil {
		return OsTypeModel{}, err
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
	data *OsTypeModel,
	apiClient *sdk.APIClient,
) (*sdk.GetOsType200ResponseOsType, error) {
	if !data.Id.IsNull() {
		return getOsTypeByID(ctx, data.Id.ValueInt64(), apiClient)
	} else if !data.Name.IsNull() {
		return getOsTypeByName(ctx, data.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data OsTypeModel

	// Read config
	diags := req.Config.Get(ctx, &data)
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

	osType, err := getOsType(ctx, &data, apiClient)
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
