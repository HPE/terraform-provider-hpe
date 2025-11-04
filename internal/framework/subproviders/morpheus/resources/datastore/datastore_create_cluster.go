// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	hpeErrors "github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

var (
	tenantsFuncCluster             = sdk.NewListCloudDatastores200ResponseAllOfDatastoresInnerTenantsInnerWithDefaults
	resourcePermissionsFuncCluster = sdk.NewSaveClusterDatastoreRequestDatastoreResourcePermissionsWithDefaults
	permissionsSitesFuncCluster    = sdk.NewGetAlerts200ResponseAllOfChecksInnerAccountWithDefaults
	datastoreTypeFuncCluster       = sdk.NewGetAlerts200ResponseAllOfChecksInnerAccountWithDefaults
	nfsConfigFuncCluster           = sdk.NewNFSDatastoreConfigurationWithDefaults
	alletrampHvmConfigFuncCluster  = sdk.NewAlletraMPHVMDatastoreConfigurationWithDefaults
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
	datastoreCreate := sdk.NewSaveClusterDatastoreRequestDatastoreWithDefaults()

	// Set the required fields
	datastoreCreate.SetName(name)

	// Set the type
	datastoreTypeForRequest := datastoreTypeFuncCluster()
	datastoreTypeForRequest.SetId(datastoreType.Id.ValueInt64())
	datastoreCreate.SetDatastoreType(*datastoreTypeForRequest)

	// Set the config.  As far as I can tell you need a config object, even if empty.
	// The config can be one of several types, handled below.
	// If none of the specific types are set, then use the generic config map.
	// The specific types are mutually exclusive.
	createConfig := datastoreCreate.GetConfig()
	switch {
	case !plan.ConfigNfs.IsNull() && !plan.ConfigNfs.IsUnknown():
		nfsConfig := nfsConfigFuncCluster()

		if !plan.ConfigNfs.SourceHostname.IsNull() {
			nfsConfig.SetSourceHostname(plan.ConfigNfs.SourceHostname.ValueString())
		}

		if !plan.ConfigNfs.SourceDirPath.IsNull() {
			nfsConfig.SetSourceDirPath(plan.ConfigNfs.SourceDirPath.ValueString())
		}

		if !plan.ConfigNfs.SourceVersion.IsNull() {
			nfsConfig.SetSourceVersion(plan.ConfigNfs.SourceVersion.ValueString())
		}

		createConfig.NFSDatastoreConfiguration = nfsConfig

	case !plan.ConfigAlletrampHvm.IsNull() && !plan.ConfigAlletrampHvm.IsUnknown():
		alletrampHvmConfig := alletrampHvmConfigFuncCluster()

		if !plan.ConfigAlletrampHvm.EnableRansomware.IsUnknown() {
			enableRansomwareString, _ := convertEnableRansomwareBool(plan.ConfigAlletrampHvm.EnableRansomware.ValueBool())
			alletrampHvmConfig.SetEnableransomware(enableRansomwareString)
		}

		if !plan.ConfigAlletrampHvm.ProtocolType.IsNull() {
			alletrampHvmConfig.SetProtocolType(plan.ConfigAlletrampHvm.ProtocolType.ValueString())
		}

		createConfig.AlletraMPHVMDatastoreConfiguration = alletrampHvmConfig

		// removing for now
		/*
			case !plan.ConfigGfs2.IsNull() && !plan.ConfigGfs2.IsUnknown():
				gfs2Config := sdk.NewSaveClusterDatastoreRequestDatastoreConfigAnyOf1WithDefaults()

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

	datastoreCreate.SetConfig(createConfig)

	// Optional fields
	if !plan.StorageServer.IsNull() && !plan.StorageServer.IsUnknown() {
		storageServerConfig := sdk.NewGetAlerts200ResponseAllOfChecksInnerAccountWithDefaults()
		storageServerConfig.SetId(plan.StorageServer.Id.ValueInt64())
		datastoreCreate.SetStorageServer(*storageServerConfig)
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		datastoreCreate.SetVisibility(plan.Visibility.ValueString())
	}

	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		datastoreCreate.SetActive(plan.Active.ValueBool())
	}

	if !plan.DefaultStore.IsNull() && !plan.DefaultStore.IsUnknown() {
		datastoreCreate.SetDefaultStore(plan.DefaultStore.ValueBool())
	}

	if !plan.Tenants.IsNull() && !plan.Tenants.IsUnknown() {
		var tenantsValues []TenantsValue
		diags := plan.Tenants.ElementsAs(ctx, &tenantsValues, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return 0
		}

		var tenantPermissions []sdk.ListCloudDatastores200ResponseAllOfDatastoresInnerTenantsInner
		for _, tenantsValue := range tenantsValues {
			tenantPermission := tenantsFuncCluster()
			tenantPermission.SetId(tenantsValue.Id.ValueInt64())
			tenantPermissions = append(tenantPermissions, *tenantPermission)
		}
		datastoreCreate.SetTenants(tenantPermissions)
	}

	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		resourcePermissions := resourcePermissionsFuncCluster()
		resourcePermissions.SetAllGroups(plan.ResourcePermissions.AllGroups.ValueBool())
		resourcePermissions.SetDefaultStore(plan.ResourcePermissions.DefaultStore.ValueBool())
		resourcePermissions.SetCanManage(plan.ResourcePermissions.CanManage.ValueBool())
		resourcePermissions.SetAll(plan.ResourcePermissions.All.ValueBool())
		if !plan.ResourcePermissions.Groups.IsNull() && !plan.ResourcePermissions.Groups.IsUnknown() {
			var groupsValues []GroupsValue
			diags := plan.ResourcePermissions.Groups.ElementsAs(ctx, &groupsValues, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return 0
			}

			sites := []sdk.GetAlerts200ResponseAllOfChecksInnerAccount{}
			for _, groupsValue := range groupsValues {
				site := permissionsSitesFuncCluster()
				site.SetId(groupsValue.Id.ValueInt64())
				sites = append(sites, *site)
			}

			resourcePermissions.SetSites(sites)
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
				planItem.SetId(plansValue.Id.ValueInt64())
				planItem.SetCode(plansValue.Code.ValueString())
				planItem.SetName(plansValue.Name.ValueString())
				plans = append(plans, *planItem)
			}

			resourcePermissions.SetPlans(plans)
		}

		datastoreCreate.SetResourcePermissions(*resourcePermissions)
	}

	// Call API
	datastoreRequest := sdk.NewSaveClusterDatastoreRequestWithDefaults()
	datastoreRequest.SetDatastore(*datastoreCreate)

	response, hresp, err := client.ClustersAPI.SaveClusterDatastore(ctx, associatedResourceId).
		SaveClusterDatastoreRequest(*datastoreRequest).Execute()
	if response == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create cluster datastore resource",
			"datastore "+name+" POST failed: "+hpeErrors.ErrMsg(err, hresp),
		)

		return 0
	}

	datastore, ok := response.GetDatastoreOk()
	if !ok {
		resp.Diagnostics.AddError(
			"create cluster datastore resource",
			"datastore "+name+": could not get datastore from response",
		)

		return 0
	}
	id, ok := datastore.GetIdOk()
	if !ok || id == nil {
		resp.Diagnostics.AddError(
			"create cluster datastore resource",
			"datastore "+name+": could not get id",
		)

		return 0
	}

	return *id
}
