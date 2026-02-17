// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/compare"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	nfsDatastoreCode  = "libvirt-netfs-nfs"
	alletraMPHVMCode  = "hpedatastore-alletra-mp"
	gfs2DatastoreCode = "libvirt-dir-gfs2"

	cloudRefType   = "ComputeZone"
	clusterRefType = "ComputeServerGroup"

	associatedResourceTypeCluster = "Cluster"
	associatedResourceTypeCloud   = "Cloud"
)

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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_datastore"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = DatastoreDataSourceSchema(ctx)
}

func getDatastoreById(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (*DatastoreModel, error) {
	response, hresp, err := client.DatastoresAPI.GetDatastores(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("datastore %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	state := &DatastoreModel{}

	datastore, ok := response.GetDatastoreOk()
	if !ok {
		return nil, fmt.Errorf("datastore %d is nil", id)
	}

	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(&datastore.Name)
	state.Type = convert.StrToType(&datastore.Type)
	state.Code = convert.StrToType(datastore.Code.Get())
	state.Visibility = convert.StrToType(datastore.Visibility)
	state.Status = convert.StrToType(&datastore.Status)
	state.Active = convert.BoolToType(datastore.Active)
	state.ExternalId = convert.StrToType(datastore.ExternalId)
	state.ExternalType = convert.StrToType(datastore.ExternalType)
	state.StorageSize = convert.Int64ToType(datastore.StorageSize.Get())
	state.ExternalPath = convert.StrToType(datastore.ExternalPath)
	state.FreeSpace = convert.Int64ToType(datastore.FreeSpace.Get())
	state.HeartBeatTarget = convert.BoolToType(datastore.HeartBeatTarget)
	state.Online = convert.BoolToType(datastore.Online)
	state.StatusMessage = convert.StrToType(datastore.StatusMessage)
	state.AllowProvision = convert.BoolToType(datastore.AllowProvision)
	state.AllowWrite = convert.BoolToType(datastore.AllowWrite)
	state.AllowRead = convert.BoolToType(datastore.AllowRead)
	state.DrsEnabled = convert.BoolToType(datastore.DrsEnabled)
	state.DefaultStore = convert.BoolToType(datastore.DefaultStore)

	// Set DatastoreType, we get the code and id back from the API
	datastoreType, ok := datastore.GetDatastoreTypeOk()
	if !ok || datastoreType == nil {
		return nil, fmt.Errorf("datastore %d missing datastore type in response", id)
	}
	datastoreTypeValue := DatastoreTypeValue{}
	datastoreTypeValue.Id = convert.Int64ToType(&datastoreType.Id)
	datastoreTypeValue.Code = convert.StrToType(&datastoreType.Code)
	datastoreTypeValue.state = attr.ValueStateKnown
	state.DatastoreType = datastoreTypeValue

	// Set AssociatedResourceType based on RefType
	refType, ok := datastore.GetRefTypeOk()
	if !ok || refType == nil {
		return nil, fmt.Errorf("datastore %d missing ref type in response", id)
	}

	state.AssociatedResourceId = convert.Int64ToType(datastore.RefId)
	switch *refType {
	case cloudRefType:
		state.AssociatedResourceType = types.StringValue(associatedResourceTypeCloud)
		tenants, resourcePermissions, err := populateCloudDatastoreInformation(
			ctx, id, *datastore.RefId, client)
		if err != nil {
			return nil, err
		}
		state.Tenants = tenants
		state.ResourcePermissions = resourcePermissions

	case clusterRefType:
		state.AssociatedResourceType = types.StringValue(associatedResourceTypeCluster)
		tenants, resourcePermissions, err := populateClusterDatastoreInformation(
			ctx, id, *datastore.RefId, client)
		if err != nil {
			return nil, err
		}
		state.Tenants = tenants
		state.ResourcePermissions = resourcePermissions

	default:
		return nil, fmt.Errorf("datastore %d has invalid ref type '%s' in response", id, *refType)
	}

	// Populate Config
	switch datastoreType.GetCode() {
	case nfsDatastoreCode:
		var configNfsValue ConfigNfsValue
		for k, v := range datastore.Config {
			switch k {
			case "sourceDirPath":
				str := v.(string)
				configNfsValue.SourceDirPath = convert.StrToType(&str)
			case "sourceHostname":
				str := v.(string)
				configNfsValue.SourceHostname = convert.StrToType(&str)
			}
		}
		configNfsValue.state = attr.ValueStateKnown
		state.ConfigNfs = configNfsValue

	case alletraMPHVMCode:
		var configAlletraMPHVM ConfigAlletrampHvmValue
		for k, v := range datastore.Config {
			switch k {
			case "enableransomware":
				switch t := v.(type) {
				case bool:
					configAlletraMPHVM.EnableRansomware = convert.BoolToType(&t)
				case string:
					configAlletraMPHVM.EnableRansomware = convert.StringToBool(ctx, t)
				}
			case "protocolType":
				str := v.(string)
				configAlletraMPHVM.ProtocolType = convert.StrToType(&str)
			}
		}
		configAlletraMPHVM.state = attr.ValueStateKnown
		state.ConfigAlletrampHvm = configAlletraMPHVM

		// The Alletra MP HVM plugin returns extra values in config other than those set in the API POST
		// We put these extra values into "config_from_api"
		keysMap := map[string]string{
			"enable_ransomware": "enableransomware",
			"protocol_type":     "protocolType",
		}
		keysFromMap, _ := compare.CheckPlanAttributeAgainstAPIAttribute(
			ctx, configAlletraMPHVM, datastore.Config, keysMap)

		configFromAPI, _ := createConfigFromApiDynamic(ctx, datastore.Config, keysFromMap)
		state.ConfigFromApi = configFromAPI

	case gfs2DatastoreCode:
		var configGfs2 ConfigGfs2Value
		for k, v := range datastore.Config {
			switch k {
			case "blockDevice":
				str := v.(string)
				configGfs2.BlockDevice = convert.StrToType(&str)
			}
		}
		configGfs2.state = attr.ValueStateKnown
		state.ConfigGfs2 = configGfs2

	default:
		// Generic config
		config, _ := convert.MapToDynamic(ctx, datastore.Config)
		state.Config = config
	}

	// Set StorageServer to that returned by the API
	if server, ok := datastore.GetStorageServerOk(); ok && server != nil {
		storageServer := StorageServerValue{}
		storageServer.Id = convert.Int64ToType(server.Id)
		storageServer.state = attr.ValueStateKnown
		state.StorageServer = storageServer
	}

	state.Cloud = NewCloudValueNull()
	if cloudId, ok := datastore.Zone.GetIdOk(); ok && cloudId != nil {
		cloud := CloudValue{}
		cloud.Id = convert.Int64ToType(cloudId)
		cloud.state = attr.ValueStateKnown
		state.Cloud = cloud
	}

	state.ResourcePool = NewResourcePoolValueNull()
	if poolId, ok := datastore.ZonePool.GetIdOk(); ok && poolId != nil {
		resourcePool := ResourcePoolValue{}
		resourcePool.Id = convert.Int64ToType(poolId)
		resourcePool.state = attr.ValueStateKnown
		state.ResourcePool = resourcePool
	}

	state.Owner = NewOwnerValueNull()
	if ownerId, ok := datastore.Owner.GetIdOk(); ok && ownerId != nil {
		owner := OwnerValue{}
		owner.Id = convert.Int64ToType(ownerId)
		owner.state = attr.ValueStateKnown
		state.Owner = owner
	}

	state.Locations = basetypes.NewSetNull(LocationsValue{}.Type(ctx))
	if locations, ok := datastore.GetLocationsOk(); ok && locations != nil {
		locationsSet, d := convert.ToSetType(
			ctx,
			locations,
			func(
				in sdk.GetDatastores200ResponseAllOfDatastoreLocationsInner,
			) LocationsValue {
				return LocationsValue{
					RefId:         convert.Int64ToType(in.RefId),
					RefType:       convert.StrToType(in.RefType),
					Status:        convert.StrToType(in.Status),
					StatusMessage: convert.StrToType(in.StatusMessage),
					state:         attr.ValueStateKnown,
				}
			},
		)
		if d.HasError() {
			return nil, fmt.Errorf("datastore %d error in creating locations set", id)
		}

		state.Locations = locationsSet
	}

	state.Datastores = basetypes.NewSetNull(DatastoresValue{}.Type(ctx))
	if datastores, ok := datastore.GetDatastoresOk(); ok && datastores != nil {
		datastoresSet, d := convert.ToSetType(
			ctx,
			datastores,
			func(
				in sdk.GetDatastores200ResponseAllOfDatastoreDatastoresInner,
			) DatastoresValue {
				id64 := int64(*in.Id)

				return DatastoresValue{
					Id:    convert.Int64ToType(&id64),
					Name:  convert.StrToType(in.Name),
					state: attr.ValueStateKnown,
				}
			},
		)
		if d.HasError() {
			return nil, fmt.Errorf("datastore %d error in creating datastores set", id)
		}
		state.Datastores = datastoresSet
	}

	return state, nil
}

// populateCloudDatastoreInformation gets the cloud datastore information for a datastore
// and populates the Tenants and ResourcePermissions fields of the datastore resource model
func populateCloudDatastoreInformation(
	ctx context.Context,
	id, cloudId int64,
	client *sdk.APIClient,
) (types.Set, ResourcePermissionsValue, error) {
	// Get cloud datastore information
	cdResp, hresp, err := client.CloudsAPI.GetCloudDatastores(ctx, cloudId, id).Execute()
	if err != nil {
		retErr := fmt.Errorf("datastore %d cloud datastore GET failed: %s", id, errfmt.ErrMsg(err, hresp))

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), retErr
	}

	cloudDatastore, ok := cdResp.GetDatastoreOk()
	if !ok || cloudDatastore == nil {
		retErr := fmt.Errorf("datastore %d missing cloud datastore in response", id)

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), retErr
	}

	// Populate Tenants
	var tenantsSet types.Set
	tenants, ok := cloudDatastore.GetTenantsOk()
	if ok && tenants != nil {

		var tdiag diag.Diagnostics
		tenantsSet, tdiag = convert.ToSetType(
			ctx,
			tenants,
			func(
				in sdk.GetCloudDatastores200ResponseAllOfDatastoreTenantsInner,
			) TenantsValue {
				return TenantsValue{
					Id:    convert.Int64ToType(in.Id),
					state: attr.ValueStateKnown,
				}
			},
		)

		if tdiag.HasError() {
			retErr := fmt.Errorf("datastore %d error in creating tenants set", id)

			return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), retErr
		}
	}

	// Populate ResourcePermissions, we'll only do Groups for now
	var resourcePermissions ResourcePermissionsValue
	rp, ok := cloudDatastore.GetResourcePermissionOk()
	if ok && rp != nil {
		sites, ok := rp.GetSitesOk()
		if !ok || sites == nil {
			tflog.Debug(ctx, fmt.Sprintf("datastore %d missing resource permission sites in response", id))

			return tenantsSet, NewResourcePermissionsValueNull(), nil
		}
		groupsSet, gdiag := convert.ToSetType(
			ctx,
			sites,
			func(
				in sdk.GetCloudDatastores200ResponseAllOfDatastoreResourcePermissionSitesInner,
			) GroupsValue {
				return GroupsValue{
					Id:    convert.Int64ToType(in.Id),
					state: attr.ValueStateKnown,
				}
			},
		)

		if gdiag.HasError() {
			retErr := fmt.Errorf("datastore %d error in creating groups set", id)

			return tenantsSet, NewResourcePermissionsValueNull(), retErr
		}

		resourcePermissions, err = populateResourcePermissionsFromApi(ctx, id, groupsSet)
	}

	return tenantsSet, resourcePermissions, err
}

// populateClusterDatastoreInformation gets the cluster datastore information for a datastore
// and populates the Tenants and ResourcePermissions fields of the datastore resource model
func populateClusterDatastoreInformation(
	ctx context.Context,
	id, clusterId int64,
	client *sdk.APIClient,
) (types.Set, ResourcePermissionsValue, error) {
	// Get cluster datastore information
	cdResp, hresp, err := client.ClustersAPI.GetClusterDatastore(ctx, clusterId, id).Execute()
	if err != nil {
		retErr := fmt.Errorf("datastore %d cluster datastore GET failed: %s", id, errfmt.ErrMsg(err, hresp))

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), retErr
	}
	clusterDatastore, ok := cdResp.GetDatastoreOk()
	if !ok || clusterDatastore == nil {
		retErr := fmt.Errorf("datastore %d missing cluster datastore in response", id)

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), retErr
	}

	// Populate Tenants
	var tenantsSet types.Set
	tenants, ok := clusterDatastore.GetTenantsOk()
	if ok && tenants != nil {

		var tdiag diag.Diagnostics
		tenantsSet, tdiag = convert.ToSetType(
			ctx,
			tenants,
			func(
				in sdk.GetClusterDatastore200ResponseDatastoreTenantsInner,
			) TenantsValue {
				return TenantsValue{
					Id:    convert.Int64ToType(in.Id),
					state: attr.ValueStateKnown,
				}
			},
		)
		if tdiag.HasError() {
			retErr := fmt.Errorf("datastore %d error in creating tenants set", id)

			return tenantsSet, NewResourcePermissionsValueNull(), retErr
		}
	}

	// Populate ResourcePermissions, we'll only do Groups for now
	var resourcePermissions ResourcePermissionsValue
	rp, ok := clusterDatastore.GetResourcePermissionsOk()
	if ok && rp != nil {

		sites, ok := rp.GetSitesOk()
		if !ok || sites == nil {
			tflog.Debug(ctx, fmt.Sprintf("datastore %d missing resource permission sites in response", id))

			return tenantsSet, NewResourcePermissionsValueNull(), nil
		}

		groupsSet, gdiag := convert.ToSetType(
			ctx,
			sites,
			func(
				in sdk.GetClusterDatastore200ResponseDatastoreResourcePermissionsSitesInner,
			) GroupsValue {
				return GroupsValue{
					Id:    convert.Int64ToType(in.Id),
					state: attr.ValueStateKnown,
				}
			},
		)
		if gdiag.HasError() {
			retErr := fmt.Errorf("datastore %d error in creating groups set", id)

			return tenantsSet, NewResourcePermissionsValueNull(), retErr
		}

		resourcePermissions, err = populateResourcePermissionsFromApi(ctx, id, groupsSet)
	}

	return tenantsSet, resourcePermissions, err
}

func createConfigFromApiDynamic(
	ctx context.Context,
	config, keysFromMap map[string]any,
) (types.Dynamic, error) {
	// We get more information back from the API, we put this in config_from_api
	apiConfigForConfigApi := make(map[string]any)
	for k := range keysFromMap {
		if _, ok := config[k]; ok {
			apiConfigForConfigApi[k] = config[k]
		}
	}

	configFromApi, err := convert.MapToDynamic(ctx, apiConfigForConfigApi)
	if err != nil {
		return types.DynamicNull(), err
	}

	return configFromApi, nil
}

// populateResourcePermissionsFromApi populates the ResourcePermissionsValue
// from the API returned groups set
// we only populate Groups for now
func populateResourcePermissionsFromApi(
	ctx context.Context,
	id int64,
	groupsSet basetypes.SetValue,
) (ResourcePermissionsValue, error) {
	var err error

	if !groupsSet.IsNull() {
		attrTypes := make(map[string]attr.Type)
		attrValues := make(map[string]attr.Value)

		attrTypes["groups"] = groupsSet.Type(ctx)
		attrValues["groups"] = groupsSet

		resourcePermissions, rdiags := NewResourcePermissionsValue(attrTypes, attrValues)
		if rdiags.HasError() {
			err = fmt.Errorf("datastore %d error in creating resource permissions", id)
		}

		return resourcePermissions, err
	}

	return NewResourcePermissionsValueNull(), err
}

func getDatastoreByName(
	ctx context.Context,
	name string,
	client *sdk.APIClient,
) (*DatastoreModel, error) {
	datastores, hresp, err := client.DatastoresAPI.ListDatastores(ctx).Name(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("datastore %s list failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	matchingDatastores, ok := datastores.GetDatastoresOk()
	if !ok || len(matchingDatastores) == 0 {
		return nil, fmt.Errorf("datastore %s not found", name)
	}

	// Not sure if this can happen
	if len(matchingDatastores) > 1 {
		var datastoreIDs []string
		for _, n := range matchingDatastores {
			if id, ok := n.GetIdOk(); ok {
				datastoreIDs = append(datastoreIDs, fmt.Sprintf("%d", *id))
			}
		}

		return nil, fmt.Errorf(
			"multiple datastores found with name %s. Datastore IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(datastoreIDs, ", "),
		)
	}

	id, ok := matchingDatastores[0].GetIdOk()
	if !ok {
		return nil, fmt.Errorf("datastore %s has missing ID", name)
	}

	return getDatastoreById(ctx, *id, client)
}

func getDatastore(
	ctx context.Context,
	config DatastoreModel,
	client *sdk.APIClient,
) (*DatastoreModel, error) {
	if !config.Id.IsNull() {
		return getDatastoreById(ctx, config.Id.ValueInt64(), client)
	}
	if !config.Name.IsNull() {
		return getDatastoreByName(ctx, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config DatastoreModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read datastore data source",
			fmt.Sprintf("failed to create client: %s", err.Error()),
		)

		return
	}

	state, err := getDatastore(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(
			"read datastore data source",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
