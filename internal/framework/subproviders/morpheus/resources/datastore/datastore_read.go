// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
	"github.com/HPE/terraform-provider-hpe/internal/framework/utils"
)

const (
	nfsDatastoreCode = "libvirt-netfs-nfs"
	alletraMPHVMCode = "hpedatastore-alletra-mp"
	// gfs2DatastoreCode = "libvirt-dir-gfs2"
)

// TODO: when implementing the datasource, note that Config can be empty.  If it is then there won't be a DatastoreType.
// populate datastore resource model with current API values
func getDatastoreAsState(
	ctx context.Context,
	id int64,
	plan DatastoreModel,
	client *sdk.APIClient,
) (DatastoreModel, diag.Diagnostics) {
	var state DatastoreModel
	var diags diag.Diagnostics

	d, hresp, err := client.DatastoresAPI.GetDatastores(ctx, id).Execute()
	if err != nil || d == nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate datastore resource",
			fmt.Sprintf("datastore %d GET failed: ", id)+errors.ErrMsg(err, hresp),
		)

		return state, diags
	}

	datastore, ok := d.GetDatastoreOk()
	if !ok || datastore == nil {
		diags.AddError(
			"populate datastore resource",
			fmt.Sprintf("datastore %d not found in response", id),
		)

		return state, diags
	}

	state.Name = convert.StrToType(&datastore.Name)
	state.AssociatedResourceId = convert.Int64ToType(datastore.RefId)

	// Set DatastoreType, we get the code and id back from the API
	datastoreType, ok := datastore.GetDatastoreTypeOk()
	if !ok || datastoreType == nil {
		diags.AddError(
			"populate datastore resource",
			fmt.Sprintf("datastore %d missing datastore type in response", id),
		)

		return state, diags
	}
	datastoreTypeValue := DatastoreTypeValue{}
	datastoreTypeValue.Id = convert.Int64ToType(&datastoreType.Id)
	datastoreTypeValue.Code = convert.StrToType(&datastoreType.Code)
	datastoreTypeValue.state = attr.ValueStateKnown
	state.DatastoreType = datastoreTypeValue

	// Set AssociatedResourceType based on RefType
	refType, ok := datastore.GetRefTypeOk()
	if !ok || refType == nil {
		diags.AddError(
			"populate datastore resource",
			fmt.Sprintf("datastore %d missing ref type in response", id),
		)

		return state, diags
	}
	switch *refType {
	case cloudRefType:
		state.AssociatedResourceType = types.StringValue(associatedResourceTypeCloud)
		// For imports, we need to get the tenants and resource permissions from the cloud datastore API
		if plan.Name.IsNull() || plan.Name.IsUnknown() {
			tenants, resourcePermissions, pdiags := populateCloudDatastoreInformation(
				ctx, id, *datastore.RefId, client)
			diags = append(diags, pdiags...)
			state.Tenants = tenants
			state.ResourcePermissions = resourcePermissions

			break
		}

		// If not importing, set from plan
		state.ResourcePermissions = plan.ResourcePermissions
		state.Tenants = plan.Tenants

	case clusterRefType:
		state.AssociatedResourceType = types.StringValue(associatedResourceTypeCluster)
		// For imports, we need to get the tenants and resource permissions from the cluster datastore API
		if plan.Name.IsNull() || plan.Name.IsUnknown() {
			tenants, resourcePermissions, pdiags := populateClusterDatastoreInformation(
				ctx, id, *datastore.RefId, client)
			diags = append(diags, pdiags...)
			state.Tenants = tenants
			state.ResourcePermissions = resourcePermissions

			break
		}

		// If not importing, set from plan
		state.ResourcePermissions = plan.ResourcePermissions
		state.Tenants = plan.Tenants

	default:
		diags.AddError(
			"populate datastore resource",
			fmt.Sprintf("datastore %d has invalid ref type '%s' in response", id, *refType),
		)

		return state, diags
	}

	// Populate Config
	switch datastoreType.GetCode() {
	case nfsDatastoreCode:
		// Check returned config against plan
		keysMap := map[string]string{
			"source_dir_path": "sourceDirPath",
			"source_hostname": "sourceHostname",
		}
		_, pdiags := utils.CheckPlanAttributeAgainstAPIAttribute(ctx, plan.ConfigNfs, datastore.Config, keysMap)
		diags = append(diags, pdiags...)

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
		if !plan.ConfigNfs.SourceVersion.IsNull() && !plan.ConfigNfs.SourceVersion.IsUnknown() {
			// If the user has set a value in the plan, but the API hasn't returned one,
			// use the plan value
			configNfsValue.SourceVersion = plan.ConfigNfs.SourceVersion
		}
		configNfsValue.state = attr.ValueStateKnown

		state.ConfigNfs = configNfsValue

	case alletraMPHVMCode:
		// Check returned config against plan
		keysMap := map[string]string{
			"enable_ransomware": "enableransomware",
			"protocol_type":     "protocolType",
		}
		keysFromMap, pdiags := utils.CheckPlanAttributeAgainstAPIAttribute(
			ctx, plan.ConfigAlletrampHvm, datastore.Config, keysMap)
		diags = append(diags, pdiags...)

		var configAlletraMPHVM ConfigAlletrampHvmValue
		for k, v := range datastore.Config {
			switch k {
			case "enableransomware":
				s := v.(string)
				b, _ := convertEnableRansomwareString(s)
				configAlletraMPHVM.EnableRansomware = convert.BoolToType(&b)
			case "protocolType":
				str := v.(string)
				configAlletraMPHVM.ProtocolType = convert.StrToType(&str)
			}
		}
		configAlletraMPHVM.state = attr.ValueStateKnown
		state.ConfigAlletrampHvm = configAlletraMPHVM

		configFromAPI, pdiags := createConfigFromApiDynamic(ctx, id, datastore.Config, keysFromMap)
		diags = append(diags, pdiags...)
		if diags.HasError() {
			return state, diags
		}

		state.ConfigFromApi = configFromAPI

	// removing for now
	/*
		case gfs2DatastoreCode:
			// Check returned config against plan
			keysMap := map[string]string{
				"block_device": "blockDevice"}
			_, pdiags := utils.CheckPlanAttributeAgainstAPIAttribute(ctx, plan.ConfigGfs2, datastore.Config, keysMap)
			diags = append(diags, pdiags...)

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

	*/

	default:
		// Dynamic config for unknown types
		// Note that some plugins return information not set by the user, so we need to filter those out
		// and store them in config_from_api
		keysFromMap, pdiags := utils.CheckPlanAttributeAgainstAPIAttribute(ctx, plan.Config, datastore.Config, nil)
		diags = append(diags, pdiags...)
		apiConfigForConfig := datastore.Config
		for k := range keysFromMap {
			if _, ok := datastore.Config[k]; ok {
				delete(apiConfigForConfig, k)
			}
		}

		configDynamic, err := convert.MapToDynamic(ctx, apiConfigForConfig)
		if err != nil {
			diags.AddError(
				"populate datastore resource",
				fmt.Sprintf("datastore %d: failed to convert config to dynamic config: %s", id, err.Error()),
			)

			return state, diags
		}
		state.Config = configDynamic

		configFromApi, pdiags := createConfigFromApiDynamic(ctx, id, datastore.Config, keysFromMap)
		diags = append(diags, pdiags...)
		if diags.HasError() {
			return state, diags
		}

		state.ConfigFromApi = configFromApi
	}

	// Set StorageServer to that returned by the API
	if server, ok := datastore.GetStorageServerOk(); ok && server != nil {
		storageServer := StorageServerValue{}
		storageServer.Id = convert.Int64ToType(server.Id)
		storageServer.state = attr.ValueStateKnown
		state.StorageServer = storageServer
	}

	state.Visibility = convert.StrToType(datastore.Visibility)
	state.Active = convert.BoolToType(datastore.Active)
	state.DefaultStore = convert.BoolToType(datastore.DefaultStore)

	state.Type = convert.StrToType(&datastore.Type)
	state.Id = convert.Int64ToType(&datastore.Id)
	state.StorageSize = convert.Int64ToType(datastore.StorageSize.Get())
	state.DrsEnabled = convert.BoolToType(datastore.DrsEnabled)
	state.AllowWrite = convert.BoolToType(datastore.AllowWrite)
	state.Online = convert.BoolToType(datastore.Online)
	state.AllowRead = convert.BoolToType(datastore.AllowRead)
	state.AllowProvision = convert.BoolToType(datastore.AllowProvision)
	state.HeartBeatTarget = convert.BoolToType(datastore.HeartBeatTarget)
	state.ExternalId = convert.StrToType(datastore.ExternalId)
	state.ExternalPath = convert.StrToType(datastore.ExternalPath)
	state.ExternalType = convert.StrToType(datastore.ExternalType)
	state.FreeSpace = convert.Int64ToType(datastore.FreeSpace.Get())
	state.StorageSize = convert.Int64ToType(datastore.StorageSize.Get())
	state.Code = convert.StrToType(datastore.Code.Get())
	state.Status = convert.StrToType(&datastore.Status)
	state.StatusMessage = convert.StrToType(datastore.StatusMessage)

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
				in sdk.ListDatastores200ResponseAllOfDatastoresInnerLocationsInner,
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
		diags = append(diags, d...)
		state.Locations = locationsSet
	}

	state.Datastores = basetypes.NewSetNull(DatastoresValue{}.Type(ctx))
	if datastores, ok := datastore.GetDatastoresOk(); ok && datastores != nil {
		datastoresSet, d := convert.ToSetType(
			ctx,
			datastores,
			func(
				in sdk.ListDatastores200ResponseAllOfDatastoresInnerDatastoresInner,
			) DatastoresValue {
				id64 := int64(*in.Id)

				return DatastoresValue{
					Id:    convert.Int64ToType(&id64),
					Name:  convert.StrToType(in.Name),
					state: attr.ValueStateKnown,
				}
			},
		)
		diags = append(diags, d...)
		state.Datastores = datastoresSet
	}

	return state, diags
}

// populateCloudDatastoreInformation gets the cloud datastore information for a datastore
// and populates the Tenants and ResourcePermissions fields of the datastore resource model
// only to be called on Import
func populateCloudDatastoreInformation(
	ctx context.Context,
	id, cloudId int64,
	client *sdk.APIClient,
) (types.Set, ResourcePermissionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get cloud datastore information
	cdResp, hresp, err := client.CloudsAPI.GetCloudDatastores(ctx, cloudId, id).Execute()
	if err != nil {
		diags.AddWarning(
			"populate datastore resource",
			fmt.Sprintf("datastore %d cloud datastore GET failed: %s", id, errors.ErrMsg(err, hresp)),
		)

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), diags
	}

	cloudDatastore, ok := cdResp.GetDatastoreOk()
	if !ok || cloudDatastore == nil {
		diags.AddWarning(
			"populate datastore resource",
			fmt.Sprintf("datastore %d missing cloud datastore in response", id),
		)

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), diags
	}

	// Populate Tenants
	tenants, ok := cloudDatastore.GetTenantsOk()
	if !ok || tenants == nil {
		diags.AddWarning(
			"populate datastore resource",
			fmt.Sprintf("datastore %d missing tenants in response", id),
		)

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), diags
	}

	tenantsSet, tdiag := convert.ToSetType(
		ctx,
		tenants,
		func(
			in sdk.ListCloudDatastores200ResponseAllOfDatastoresInnerTenantsInner,
		) TenantsValue {
			return TenantsValue{
				Id:    convert.Int64ToType(in.Id),
				state: attr.ValueStateKnown,
			}
		},
	)

	diags = append(diags, tdiag...)

	// Populate ResourcePermissions, we'll only do Groups for now
	rp, ok := cloudDatastore.GetResourcePermissionOk()
	if !ok || rp == nil {
		diags.AddWarning(
			"populate datastore resource",
			fmt.Sprintf("datastore %d missing resource permissions in response", id),
		)

		return tenantsSet, NewResourcePermissionsValueNull(), diags
	}

	sites, ok := rp.GetSitesOk()
	if !ok || sites == nil {
		diags.AddWarning(
			"populate datastore resource",
			fmt.Sprintf("datastore %d missing resource permission sites in response", id),
		)

		return tenantsSet, NewResourcePermissionsValueNull(), diags
	}
	groupsSet, gdiag := convert.ToSetType(
		ctx,
		sites,
		func(
			in sdk.ListCloudDatastores200ResponseAllOfDatastoresInnerResourcePermissionSitesInner,
		) GroupsValue {
			return GroupsValue{
				Id:    convert.Int64ToType(in.Id),
				state: attr.ValueStateKnown,
			}
		},
	)

	diags = append(diags, gdiag...)

	resourcePermissions, rdiags := populateResourcePermissionsFromApi(ctx, groupsSet)
	diags = append(diags, rdiags...)

	return tenantsSet, resourcePermissions, diags
}

// populateClusterDatastoreInformation gets the cluster datastore information for a datastore
// and populates the Tenants and ResourcePermissions fields of the datastore resource model
// only to be called on Import
func populateClusterDatastoreInformation(
	ctx context.Context,
	id, clusterId int64,
	client *sdk.APIClient,
) (types.Set, ResourcePermissionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get cluster datastore information
	cdResp, hresp, err := client.ClustersAPI.GetClusterDatastore(ctx, clusterId, id).Execute()
	if err != nil {
		diags.AddWarning(
			"populate datastore resource",
			fmt.Sprintf("datastore %d cluster datastore GET failed: %s", id, errors.ErrMsg(err, hresp)),
		)

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), diags
	}
	clusterDatastore, ok := cdResp.GetDatastoreOk()
	if !ok || clusterDatastore == nil {
		diags.AddWarning(
			"populate datastore resource",
			fmt.Sprintf("datastore %d missing cluster datastore in response", id),
		)

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), diags
	}

	// Populate Tenants
	tenants, ok := clusterDatastore.GetTenantsOk()
	if !ok || tenants == nil {
		diags.AddWarning(
			"populate datastore resource",
			fmt.Sprintf("datastore %d missing tenants in response", id),
		)

		return types.SetNull(TenantsType{}), NewResourcePermissionsValueNull(), diags
	}

	tenantsSet, tdiag := convert.ToSetType(
		ctx,
		tenants,
		func(
			in sdk.ListCloudDatastores200ResponseAllOfDatastoresInnerTenantsInner,
		) TenantsValue {
			return TenantsValue{
				Id:    convert.Int64ToType(in.Id),
				state: attr.ValueStateKnown,
			}
		},
	)

	diags = append(diags, tdiag...)

	// Populate ResourcePermissions, we'll only do Groups for now
	rp, ok := clusterDatastore.GetResourcePermissionsOk()
	if !ok || rp == nil {
		diags.AddWarning(
			"populate datastore resource",
			fmt.Sprintf("datastore %d missing resource permissions in response", id),
		)

		return tenantsSet, NewResourcePermissionsValueNull(), diags
	}
	sites, ok := rp.GetSitesOk()
	if !ok || sites == nil {
		diags.AddWarning(
			"populate datastore resource",
			fmt.Sprintf("datastore %d missing resource permission sites in response", id),
		)

		return tenantsSet, NewResourcePermissionsValueNull(), diags
	}

	groupsSet, gdiag := convert.ToSetType(
		ctx,
		sites,
		func(
			in sdk.ListCloudDatastores200ResponseAllOfDatastoresInnerResourcePermissionSitesInner,
		) GroupsValue {
			return GroupsValue{
				Id:    convert.Int64ToType(in.Id),
				state: attr.ValueStateKnown,
			}
		},
	)

	diags = append(diags, gdiag...)

	resourcePermissions, rdiags := populateResourcePermissionsFromApi(ctx, groupsSet)
	diags = append(diags, rdiags...)

	return tenantsSet, resourcePermissions, diags
}

// populateResourcePermissionsFromApi populates the ResourcePermissionsValue
// from the API returned groups set
// we only populate Groups for now
func populateResourcePermissionsFromApi(
	ctx context.Context,
	groupsSet basetypes.SetValue,
) (ResourcePermissionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !groupsSet.IsNull() {
		attrTypes := make(map[string]attr.Type)
		attrValues := make(map[string]attr.Value)

		attrTypes["all"] = types.BoolType
		attrValues["all"] = types.BoolNull()

		attrTypes["all_groups"] = types.BoolType
		attrValues["all_groups"] = types.BoolNull()

		attrTypes["all_plans"] = types.BoolType
		attrValues["all_plans"] = types.BoolNull()

		attrTypes["can_manage"] = types.BoolType
		attrValues["can_manage"] = types.BoolNull()

		attrTypes["default_store"] = types.BoolType
		attrValues["default_store"] = types.BoolNull()

		attrTypes["default_target"] = types.BoolType
		attrValues["default_target"] = types.BoolNull()

		plansNull := NewPlansValueNull()
		plansSetNull := basetypes.NewSetNull(plansNull.Type(ctx))
		attrTypes["plans"] = plansSetNull.Type(ctx)
		attrValues["plans"] = plansSetNull

		attrTypes["groups"] = groupsSet.Type(ctx)
		attrValues["groups"] = groupsSet

		resourcePermissions, rdiags := NewResourcePermissionsValue(attrTypes, attrValues)
		diags = append(diags, rdiags...)

		return resourcePermissions, diags
	}

	return NewResourcePermissionsValueNull(), diags
}

func createConfigFromApiDynamic(
	ctx context.Context,
	id int64,
	config, keysFromMap map[string]any,
) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics
	// We get more information back from the API, we put this in config_from_api
	apiConfigForConfigApi := make(map[string]any)
	for k := range keysFromMap {
		if _, ok := config[k]; ok {
			apiConfigForConfigApi[k] = config[k]
		}
	}

	configFromApi, err := convert.MapToDynamic(ctx, apiConfigForConfigApi)
	if err != nil {
		diags.AddError(
			"populate datastore resource",
			fmt.Sprintf("datastore %d: failed to convert config_from_api to dynamic config: %s", id, err.Error()),
		)

		return types.DynamicNull(), diags
	}

	return configFromApi, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan DatastoreModel

	diags := req.State.Get(ctx, &plan)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read datastore resource",
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()

	state, pdiags := getDatastoreAsState(ctx, id, plan, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"read datastore resource",
			fmt.Sprintf("datastore %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func convertEnableRansomwareString(input string) (bool, bool) {
	enableRansomwareMap := map[string]bool{
		"on":  true,
		"off": false,
	}

	ret, ok := enableRansomwareMap[input]

	return ret, ok
}
