// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	hpeErrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

var (
	tenantsClusterFunc             = func() *sdk.SaveClusterDatastoreRequestDatastoreTenantsInner { return &sdk.SaveClusterDatastoreRequestDatastoreTenantsInner{} }
	resourcePermissionsClusterFunc = func() *sdk.SaveClusterDatastoreRequestDatastoreResourcePermissions { return &sdk.SaveClusterDatastoreRequestDatastoreResourcePermissions{} }
	permissionsSitesClusterFunc    = func() *sdk.SaveClusterDatastoreRequestDatastoreResourcePermissionsSitesInner { return &sdk.SaveClusterDatastoreRequestDatastoreResourcePermissionsSitesInner{} }
	datastoreTypeClusterFunc       = func() *sdk.SaveClusterDatastoreRequestDatastoreDatastoreType { return &sdk.SaveClusterDatastoreRequestDatastoreDatastoreType{} }
	nfsConfigClusterFunc           = func() *sdk.NFSDatastoreConfiguration { return &sdk.NFSDatastoreConfiguration{} }
	alletrampHvmConfigClusterFunc  = func() *sdk.AlletraMPHVMDatastoreConfiguration { return &sdk.AlletraMPHVMDatastoreConfiguration{} }
)

func datastoreCreateCluster(ctx context.Context,
	datastoreType DatastoreTypeValue,
	name string,
	associatedResourceId int64,
	client *sdk.APIClient,
	plan DatastoreModel,
	resp *resource.CreateResponse,
) int64 {
	// datastoreCreate is used by the SDK to create the datastore
	datastoreCreate := &sdk.SaveClusterDatastoreRequestDatastore{}

	// Set the required fields
	datastoreCreate.Name = plan.Name.ValueStringPointer()

	// Set the type
	datastoreTypeForRequest := datastoreTypeClusterFunc()
	datastoreTypeForRequest.Id = plan.DatastoreType.Id.ValueInt64Pointer()
	datastoreCreate.DatastoreType = datastoreTypeForRequest

	// Set the config.  As far as I can tell you need a config object, even if empty.
	// The config can be one of several types, handled below.
	// If none of the specific types are set, then use the generic config map.
	// The specific types are mutually exclusive.
	createConfig := datastoreCreate.Config
	if createConfig == nil {
		createConfig = &sdk.SaveClusterDatastoreRequestDatastoreConfig{}
	}
	switch {
	case !plan.ConfigNfs.IsNull() && !plan.ConfigNfs.IsUnknown():
		nfsConfig := nfsConfigClusterFunc()

		if !plan.ConfigNfs.SourceHostname.IsNull() {
			nfsConfig.SourceHostname = plan.ConfigNfs.SourceHostname.ValueString()
		}

		if !plan.ConfigNfs.SourceDirPath.IsNull() {
			nfsConfig.SourceDirPath = plan.ConfigNfs.SourceDirPath.ValueString()
		}

		if !plan.ConfigNfs.SourceVersion.IsNull() {
			nfsConfig.SourceVersion = plan.ConfigNfs.SourceVersion.ValueStringPointer()
		}

		createConfig.NFSDatastoreConfiguration = nfsConfig

	case !plan.ConfigAlletrampHvm.IsNull() && !plan.ConfigAlletrampHvm.IsUnknown():
		alletrampHvmConfig := alletrampHvmConfigClusterFunc()

		if !plan.ConfigAlletrampHvm.EnableRansomware.IsUnknown() {
			enableRansomwareString := convert.BoolToStringOnOff(plan.ConfigAlletrampHvm.EnableRansomware.ValueBool())
			alletrampHvmConfig.Enableransomware = enableRansomwareString.ValueStringPointer()
		}

		if !plan.ConfigAlletrampHvm.ProtocolType.IsNull() {
			alletrampHvmConfig.ProtocolType = plan.ConfigAlletrampHvm.ProtocolType.ValueString()
		}

		createConfig.AlletraMPHVMDatastoreConfiguration = alletrampHvmConfig

		// removing for now
		/*
			case !plan.ConfigGfs2.IsNull() && !plan.ConfigGfs2.IsUnknown():
				gfs2Config := &sdk.SaveClusterDatastoreRequestDatastoreConfigAnyOf1{}

				if !plan.ConfigGfs2.BlockDevice.IsNull() {
					gfs2Config.SetBlockDevice(plan.ConfigGfs2.BlockDevice.ValueString())
				}

				if !plan.ConfigGfs2.AllowReformat.IsUnknown() {
					gfs2Config.SetAllowReformat(plan.ConfigGfs2.AllowReformat.ValueBool())
				}

				createConfig.SaveClusterDatastoreRequestDatastoreConfigAnyOf1 = gfs2Config

		*/

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"create cluster datastore resource",
				"datastore "+name+": failed to convert config: "+
					err.Error(),
			)

			return 0
		}

		configDataMap, ok := configMap.(map[string]any)
		if ok {
			createConfig.MapmapOfStringAny = &configDataMap
		} else {
			resp.Diagnostics.AddError(
				"create cluster datastore resource",
				"datastore "+name+": config must be a valid object/map",
			)

			return 0
		}

	}

	datastoreCreate.Config = createConfig

	// Optional fields
	if !plan.StorageServer.IsNull() && !plan.StorageServer.IsUnknown() {
		storageServerConfig := &sdk.SaveClusterDatastoreRequestDatastoreStorageServer{}
		storageServerConfig.Id = plan.StorageServer.Id.ValueInt64Pointer()
		datastoreCreate.StorageServer = storageServerConfig
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		datastoreCreate.Visibility = plan.Visibility.ValueStringPointer()
	}

	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		datastoreCreate.Active = plan.Active.ValueBoolPointer()
	}

	if !plan.DefaultStore.IsNull() && !plan.DefaultStore.IsUnknown() {
		datastoreCreate.DefaultStore = plan.DefaultStore.ValueBoolPointer()
	}

	if !plan.Tenants.IsNull() && !plan.Tenants.IsUnknown() {
		var tenantsValues []TenantsValue
		diags := plan.Tenants.ElementsAs(ctx, &tenantsValues, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return 0
		}

		var tenantPermissions []sdk.SaveClusterDatastoreRequestDatastoreTenantsInner
		for _, tenantsValue := range tenantsValues {
			tenantPermission := tenantsClusterFunc()
			tenantPermission.Id = tenantsValue.Id.ValueInt64Pointer()
			tenantPermissions = append(tenantPermissions, *tenantPermission)
		}
		datastoreCreate.Tenants = tenantPermissions
	}

	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		resourcePermissions := resourcePermissionsClusterFunc()
		resourcePermissions.AllGroups = plan.ResourcePermissions.AllGroups.ValueBoolPointer()
		resourcePermissions.DefaultStore = plan.ResourcePermissions.DefaultStore.ValueBoolPointer()
		resourcePermissions.CanManage = plan.ResourcePermissions.CanManage.ValueBoolPointer()
		resourcePermissions.All = plan.ResourcePermissions.All.ValueBoolPointer()
		if !plan.ResourcePermissions.Groups.IsNull() && !plan.ResourcePermissions.Groups.IsUnknown() {
			var groupsValues []GroupsValue
			diags := plan.ResourcePermissions.Groups.ElementsAs(ctx, &groupsValues, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return 0
			}

			sites := []sdk.SaveClusterDatastoreRequestDatastoreResourcePermissionsSitesInner{}
			for _, groupsValue := range groupsValues {
				site := permissionsSitesClusterFunc()
				site.Id = groupsValue.Id.ValueInt64Pointer()
				sites = append(sites, *site)
			}

			resourcePermissions.Sites = sites
		}

		// nolint:duplicate
		if !plan.ResourcePermissions.Plans.IsNull() && !plan.ResourcePermissions.Plans.IsUnknown() {
			var plansValues []PlansValue
			diags := plan.ResourcePermissions.Plans.ElementsAs(ctx, &plansValues, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return 0
			}

			var plans []sdk.SaveClusterDatastoreRequestDatastoreResourcePermissionsPlansInner
			for _, plansValue := range plansValues {
				planItem := &sdk.SaveClusterDatastoreRequestDatastoreResourcePermissionsPlansInner{}
				planItem.Id = plansValue.Id.ValueInt64Pointer()
				planItem.Code = plansValue.Code.ValueStringPointer()
				planItem.Name = plansValue.Name.ValueStringPointer()
				plans = append(plans, *planItem)
			}

			resourcePermissions.Plans = plans
		}

		datastoreCreate.ResourcePermissions = resourcePermissions
	}

	// Call API
	datastoreRequest := &sdk.SaveClusterDatastoreRequest{}
	datastoreRequest.Datastore = datastoreCreate

	response, hresp, err := client.ClustersAPI.SaveClusterDatastore(ctx, associatedResourceId).
		SaveClusterDatastoreRequest(*datastoreRequest).Execute()
	if response == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create cluster datastore resource",
			"datastore "+name+" POST failed: "+hpeErrors.ErrMsg(err, hresp),
		)

		return 0
	}

	datastore := response.Datastore
	if datastore == nil {
		resp.Diagnostics.AddError(
			"create cluster datastore resource",
			"datastore "+name+": could not get datastore from response",
		)

		return 0
	}
	id := datastore.Id
	if id == nil {
		resp.Diagnostics.AddError(
			"create cluster datastore resource",
			"datastore "+name+": could not get id",
		)

		return 0
	}

	return *id
}
