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
	permissionsFunc        = func() *sdk.SaveDatastoreRequestDatastoreResourcePermissions { return &sdk.SaveDatastoreRequestDatastoreResourcePermissions{} }
	permissionsPlansFunc   = func() *sdk.SaveDatastoreRequestDatastoreResourcePermissionsPlansInner { return &sdk.SaveDatastoreRequestDatastoreResourcePermissionsPlansInner{} }
	permissionsSitesFunc   = func() *sdk.SaveDatastoreRequestDatastoreResourcePermissionsSitesInner { return &sdk.SaveDatastoreRequestDatastoreResourcePermissionsSitesInner{} }
	tenantsFunc            = func() *sdk.SaveDatastoreRequestDatastoreTenantPermissions { return &sdk.SaveDatastoreRequestDatastoreTenantPermissions{} }
	nfsConfigFunc          = func() *sdk.NFSDatastoreConfiguration1 { return &sdk.NFSDatastoreConfiguration1{} }
	alletrampHvmConfigFunc = func() *sdk.AlletraMPHVMDatastoreConfiguration1 { return &sdk.AlletraMPHVMDatastoreConfiguration1{} }
	storageServerFunc      = func() *sdk.SaveDatastoreRequestDatastoreStorageServer { return &sdk.SaveDatastoreRequestDatastoreStorageServer{} }
)

type (
	permissionsPlans = sdk.SaveDatastoreRequestDatastoreResourcePermissionsPlansInner
	permissionsSites = sdk.SaveDatastoreRequestDatastoreResourcePermissionsSitesInner
)

func datastoreCreateDatastore(ctx context.Context,
	datastoreType DatastoreTypeValue,
	associatedResourceType, name string,
	associatedResourceId int64,
	client *sdk.APIClient,
	plan DatastoreModel,
	resp *resource.CreateResponse,
) int64 {
	// datastoreCreate is used by the SDK to create the datastore
	datastoreCreate := &sdk.SaveDatastoreRequestDatastore{}

	// Set the required fields
	datastoreCreate.Name = name
	datastoreCreate.DatastoreType = datastoreType.Code.ValueString()
	// Set the associated resource - this is refType for the API
	switch associatedResourceType {
	case associatedResourceTypeCloud:
		datastoreCreate.RefType = cloudRefType
		// TODO allow the following when API has been fixed
	// case associatedResourceTypeCluster:
	//	datastoreCreate.RefType = clusterRefType
	default:
		resp.Diagnostics.AddError(
			"create datastore resource",
			"datastore "+name+": invalid associated_resource_type "+associatedResourceType+", must be 'Cloud' or 'Cluster'",
		)
	}
	datastoreCreate.RefId = associatedResourceId

	// Set the config.  As far as I can tell you need a config object, even if empty.
	// The config can be one of several types, handled below.
	// If none of the specific types are set, then use the generic config map.
	// The specific types are mutually exclusive.
	createConfig := datastoreCreate.Config
	switch {
	case !plan.ConfigNfs.IsNull() && !plan.ConfigNfs.IsUnknown():
		nfsConfig := nfsConfigFunc()

		if !plan.ConfigNfs.SourceHostname.IsNull() {
			nfsConfig.SourceHostname = plan.ConfigNfs.SourceHostname.ValueString()
		}

		if !plan.ConfigNfs.SourceDirPath.IsNull() {
			nfsConfig.SourceDirPath = plan.ConfigNfs.SourceDirPath.ValueString()
		}

		if !plan.ConfigNfs.SourceVersion.IsNull() {
			nfsConfig.SourceVersion = plan.ConfigNfs.SourceVersion.ValueStringPointer()
		}

		createConfig.NFSDatastoreConfiguration1 = nfsConfig
	case !plan.ConfigAlletrampHvm.IsNull() && !plan.ConfigAlletrampHvm.IsUnknown():
		alletrampHvmConfig := alletrampHvmConfigFunc()

		if !plan.ConfigAlletrampHvm.EnableRansomware.IsUnknown() {
			enableRansomwareString := convert.BoolToStringOnOff(plan.ConfigAlletrampHvm.EnableRansomware.ValueBool())
			alletrampHvmConfig.Enableransomware = enableRansomwareString.ValueStringPointer()
		}

		if !plan.ConfigAlletrampHvm.ProtocolType.IsNull() {
			alletrampHvmConfig.ProtocolType = plan.ConfigAlletrampHvm.ProtocolType.ValueString()
		}
		createConfig.AlletraMPHVMDatastoreConfiguration1 = alletrampHvmConfig

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
				"create datastore resource",
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
				"create datastore resource",
				"datastore "+name+": config must be a valid object/map",
			)

			return 0
		}

	}

	datastoreCreate.Config = createConfig

	// Optional fields
	if !plan.StorageServer.IsNull() && !plan.StorageServer.IsUnknown() {
		storageServerConfig := storageServerFunc()
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

		tenantPermissions := tenantsFunc()
		var accounts []sdk.SaveDatastoreRequestDatastoreTenantPermissionsAccountsInner
		for _, tenantsValue := range tenantsValues {
			account := &sdk.SaveDatastoreRequestDatastoreTenantPermissionsAccountsInner{}
			account.Id = tenantsValue.Id.ValueInt64Pointer()
			accounts = append(accounts, *account)
		}
		tenantPermissions.Accounts = accounts
		datastoreCreate.TenantPermissions = tenantPermissions
	}

	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		resourcePermissions := permissionsFunc()
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

			var sites []permissionsSites
			for _, groupsValue := range groupsValues {
				site := permissionsSitesFunc()
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

			var plans []permissionsPlans
			for _, plansValue := range plansValues {
				planItem := permissionsPlansFunc()
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
	datastoreRequest := &sdk.SaveDatastoreRequest{}
	datastoreRequest.Datastore = datastoreCreate

	response, hresp, err := client.DatastoresAPI.SaveDatastore(ctx).SaveDatastoreRequest(*datastoreRequest).Execute()
	if response == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create datastore resource",
			"datastore "+name+" POST failed: "+hpeErrors.ErrMsg(err, hresp),
		)

		return 0
	}

	datastore := response.Datastore
	if datastore == nil {
		resp.Diagnostics.AddError(
			"create datastore resource",
			"datastore "+name+": could not get datastore from response",
		)

		return 0
	}
	id := datastore.Id
	if id == 0 {
		resp.Diagnostics.AddError(
			"create datastore resource",
			"datastore "+name+": could not get id",
		)

		return 0
	}

	return id
}
