// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storageserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                     = "read storage server data source"
	cloudListMax                = 1000
	refTypeComputeZone          = "ComputeZone"
	ErrorNoValidSearchTerms     = `no valid search terms - an id or name is required`
	ErrorNoStorageServerFound   = `no storage server found`
	ErrorMultipleStorageServers = `multiple storage servers were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "storage_server"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = StorageServerDataSourceSchema(ctx)
	resp.Schema.Description = "Retrieves information about a single Morpheus storage server by id or name."
	resp.Schema.MarkdownDescription = "Retrieves information about a single Morpheus storage server by id or name."
}

// int32PtrToType widens a nullable int32 (Grails Integer) to a Terraform Int64,
// which is the only integer width the framework supports.
func int32PtrToType(i *int32) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*i))
}

// cloudIDName derives the cloud (zone) id and name for a storage server. A
// storage server is cloud-scoped only when its refType is ComputeZone, in which
// case refId is the cloud id; the name is looked up from the supplied map.
func cloudIDName(
	refType *string,
	refID *int64,
	names map[int64]string,
) (types.Int64, types.String) {
	if refType == nil || *refType != refTypeComputeZone || refID == nil {
		return types.Int64Null(), types.StringNull()
	}

	name := types.StringNull()
	if n, ok := names[*refID]; ok {
		name = types.StringValue(n)
	}

	return types.Int64Value(*refID), name
}

// listCloudNames returns a map of cloud id to cloud name, used to resolve the
// cloud_id / cloud_name attributes from a storage server's refType/refId.
func listCloudNames(
	ctx context.Context,
	apiClient *sdk.APIClient,
) (map[int64]string, error) {
	rs, hresp, err := apiClient.CloudsAPI.ListClouds(ctx).Max(cloudListMax).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LIST failed for clouds: %s", providererrors.ErrMsg(err, hresp))
	}

	names := make(map[int64]string, len(rs.Zones))
	for i := range rs.Zones {
		z := rs.Zones[i]
		if z.Id != nil && z.Name != nil {
			names[*z.Id] = *z.Name
		}
	}

	return names, nil
}

// storageServerAsState maps an API storage server into the datasource model.
func storageServerAsState(
	ss *sdk.GetStorageServers200ResponseStorageServer,
	cloudNames map[int64]string,
) StorageServerModel {
	state := StorageServerModel{
		Id:              convert.Int64ToType(ss.Id),
		Name:            convert.StrToType(ss.Name),
		Description:     convert.StrToType(ss.Description.Get()),
		Visibility:      convert.StrToType(ss.Visibility),
		Enabled:         convert.BoolToType(ss.Enabled),
		Status:          convert.StrToType(ss.Status.Get()),
		StatusMessage:   convert.StrToType(ss.StatusMessage.Get()),
		StatusDate:      convert.TimeToType(ss.StatusDate),
		ErrorMessage:    convert.StrToType(ss.ErrorMessage.Get()),
		InternalId:      convert.StrToType(ss.InternalId.Get()),
		ExternalId:      convert.StrToType(ss.ExternalId.Get()),
		InternalIp:      convert.StrToType(ss.InternalIp.Get()),
		ExternalIp:      convert.StrToType(ss.ExternalIp.Get()),
		ServiceUrl:      convert.StrToType(ss.ServiceUrl.Get()),
		ServiceHost:     convert.StrToType(ss.ServiceHost.Get()),
		ServicePath:     convert.StrToType(ss.ServicePath.Get()),
		ServiceVersion:  convert.StrToType(ss.ServiceVersion.Get()),
		ServiceUsername: convert.StrToType(ss.ServiceUsername.Get()),
		ApiPort:         int32PtrToType(ss.ApiPort.Get()),
		AdminPort:       int32PtrToType(ss.AdminPort.Get()),
		Category:        convert.StrToType(ss.Category.Get()),
		ServerVendor:    convert.StrToType(ss.ServerVendor.Get()),
		ServerModel:     convert.StrToType(ss.ServerModel.Get()),
		SerialNumber:    convert.StrToType(ss.SerialNumber.Get()),
		MaxStorage:      convert.Int64ToType(ss.MaxStorage.Get()),
		UsedStorage:     convert.Int64ToType(ss.UsedStorage.Get()),
		DiskCount:       int32PtrToType(ss.DiskCount.Get()),
		DateCreated:     convert.TimeToType(ss.DateCreated),
		LastUpdated:     convert.TimeToType(ss.LastUpdated),
		RefType:         convert.StrToType(ss.RefType),
		RefId:           convert.Int64ToType(ss.RefId),
	}

	state.TypeId = types.Int64Null()
	state.TypeCode = types.StringNull()
	state.TypeName = types.StringNull()
	if ss.Type != nil {
		state.TypeId = convert.Int64ToType(ss.Type.Id)
		state.TypeCode = convert.StrToType(ss.Type.Code)
		state.TypeName = convert.StrToType(ss.Type.Name)
	}

	state.CloudId, state.CloudName = cloudIDName(ss.RefType, ss.RefId, cloudNames)

	return state
}

func getStorageServerByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetStorageServers200ResponseStorageServer, error) {
	r, hresp, err := apiClient.StorageAPI.GetStorageServers(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for storage server %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.StorageServer == nil {
		return nil, fmt.Errorf("GET failed for storage server %d: response missing storageServer", id)
	}

	return r.StorageServer, nil
}

func getStorageServerByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetStorageServers200ResponseStorageServer, error) {
	rs, hresp, err := apiClient.StorageAPI.ListStorageServers(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for storage server %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	var matchedID int64

	var matchCount int

	for _, ss := range rs.StorageServers {
		if ss.Name != nil && *ss.Name == name {
			if ss.Id != nil {
				matchedID = *ss.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoStorageServerFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleStorageServers)
	}

	return getStorageServerByID(ctx, matchedID, apiClient)
}

func getStorageServer(
	ctx context.Context,
	config *StorageServerModel,
	apiClient *sdk.APIClient,
) (*sdk.GetStorageServers200ResponseStorageServer, error) {
	if !config.Id.IsNull() {
		return getStorageServerByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getStorageServerByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config StorageServerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	ss, err := getStorageServer(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	cloudNames, err := listCloudNames(ctx, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := storageServerAsState(ss, cloudNames)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
