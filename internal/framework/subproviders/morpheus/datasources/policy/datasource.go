// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/datasources/policy/consts"
	internalErrors "github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

const summary = "read policy data source"

// apiTypeToResourceType converts API types back to user-facing resource types
func apiTypeToResourceType(apiType string) string {
	switch apiType {
	case "ComputeZone":
		return "Cloud"
	case "ComputeSite":
		return "Group"
	default:
		// For other types (User, Role, Network, Plan), pass through as-is
		return apiType
	}
}

// mapPolicyConfigToState maps the API config structure to the datasource schema structure
func mapPolicyConfigToState(
	ctx context.Context,
	data *PolicyModel,
	apiConfig *sdk.AddPolicies200ResponseAllOfPolicyConfig,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	// Map each API config field to the corresponding schema field - only populate non-null configurations
	// 1. ApprovePolicyTypeConfiguration2 -> approval
	if apiConfig.ApprovePolicyTypeConfiguration2 != nil {
		approvalValue, approvalDiags := NewConfigApprovalValue(
			ConfigApprovalValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"account_integration_id": convert.StrToType(&apiConfig.ApprovePolicyTypeConfiguration2.AccountIntegrationId),
				"flow_id":                convert.StrToType(apiConfig.ApprovePolicyTypeConfiguration2.FlowId),
				"workflow_id":            convert.StrToType(apiConfig.ApprovePolicyTypeConfiguration2.WorkflowId),
				"workflow_type":          convert.StrToType(apiConfig.ApprovePolicyTypeConfiguration2.WorkflowType),
			},
		)
		if approvalDiags.HasError() {
			diags.Append(approvalDiags...)
		} else {
			data.ConfigApproval = approvalValue
		}
	}

	// 2. BackupTargetsPolicyTypeConfiguration -> backup_storage
	if apiConfig.BackupTargetsPolicyTypeConfiguration != nil {
		// Handle BackupStorageIds as a set of strings
		var backupStorageIDsSet types.Set
		if len(apiConfig.BackupTargetsPolicyTypeConfiguration.BackupStorageIds) == 0 {
			backupStorageIDsSet = types.SetValueMust(types.StringType, []attr.Value{})
		} else {
			stringValues := make([]attr.Value, len(apiConfig.BackupTargetsPolicyTypeConfiguration.BackupStorageIds))
			for i, id := range apiConfig.BackupTargetsPolicyTypeConfiguration.BackupStorageIds {
				stringValues[i] = types.StringValue(id)
			}
			var setDiags diag.Diagnostics
			backupStorageIDsSet, setDiags = types.SetValueFrom(ctx, types.StringType, stringValues)
			if setDiags.HasError() {
				diags.Append(setDiags...)
			}
		}

		backupStorageAttrs := map[string]attr.Value{
			"backup_storage_ids": backupStorageIDsSet,
		}

		backupStorageValue, backupStorageDiags := NewConfigBackupStorageValue(
			ConfigBackupStorageValue{}.AttributeTypes(ctx),
			backupStorageAttrs,
		)
		if backupStorageDiags.HasError() {
			diags.Append(backupStorageDiags...)
		} else {
			data.ConfigBackupStorage = backupStorageValue
		}
	}

	// 3. BackupCreationPolicyTypeConfiguration2 -> create_backup
	if apiConfig.BackupCreationPolicyTypeConfiguration2 != nil {
		createBackupValue, createBackupDiags := NewConfigCreateBackupValue(
			ConfigCreateBackupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"create_backup":      convert.BoolToType(apiConfig.BackupCreationPolicyTypeConfiguration2.CreateBackup),
				"create_backup_type": convert.StrToType(&apiConfig.BackupCreationPolicyTypeConfiguration2.CreateBackupType),
			},
		)
		if createBackupDiags.HasError() {
			diags.Append(createBackupDiags...)
		} else {
			data.ConfigCreateBackup = createBackupValue
		}
	}

	// 4. UserCreationPolicyTypeConfiguration -> create_user
	if apiConfig.UserCreationPolicyTypeConfiguration != nil {
		createUserValue, createUserDiags := NewConfigCreateUserValue(
			ConfigCreateUserValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"create_user":      convert.BoolToType(apiConfig.UserCreationPolicyTypeConfiguration.CreateUser),
				"create_user_type": convert.StrToType(&apiConfig.UserCreationPolicyTypeConfiguration.CreateUserType),
			},
		)
		if createUserDiags.HasError() {
			diags.Append(createUserDiags...)
		} else {
			data.ConfigCreateUser = createUserValue
		}
	}

	// 5. UserGroupCreationPolicyTypeConfiguration -> create_user_group
	if apiConfig.UserGroupCreationPolicyTypeConfiguration != nil {
		createUserGroupValue, createUserGroupDiags := NewConfigCreateUserGroupValue(
			ConfigCreateUserGroupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"user_group": convert.StrToType(&apiConfig.UserGroupCreationPolicyTypeConfiguration.UserGroup),
			},
		)
		if createUserGroupDiags.HasError() {
			diags.Append(createUserGroupDiags...)
		} else {
			data.ConfigCreateUserGroup = createUserGroupValue
		}
	}

	// 6. CypherAccessPolicyTypeConfiguration -> cypher
	if apiConfig.CypherAccessPolicyTypeConfiguration != nil {
		cypherValue, cypherDiags := NewConfigCypherValue(
			ConfigCypherValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"delete":      convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration.Delete),
				"key_pattern": convert.StrToType(&apiConfig.CypherAccessPolicyTypeConfiguration.KeyPattern),
				"list":        convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration.List),
				"read":        convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration.Read),
				"update":      convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration.Update),
				"write":       convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration.Write),
			},
		)
		if cypherDiags.HasError() {
			diags.Append(cypherDiags...)
		} else {
			data.ConfigCypher = cypherValue
		}
	}

	// 7. BudgetPolicyTypeConfiguration2 -> max_price
	if apiConfig.BudgetPolicyTypeConfiguration2 != nil {
		maxPriceAttrs := map[string]attr.Value{
			"max_price":          convert.NumToType(&apiConfig.BudgetPolicyTypeConfiguration2.MaxPrice),
			"max_price_currency": convert.StrToType(apiConfig.BudgetPolicyTypeConfiguration2.MaxPriceCurrency),
			"max_price_unit":     convert.StrToType(apiConfig.BudgetPolicyTypeConfiguration2.MaxPriceUnit),
		}

		maxPriceValue, maxPriceDiags := NewConfigMaxPriceValue(ConfigMaxPriceValue{}.AttributeTypes(ctx), maxPriceAttrs)
		if maxPriceDiags.HasError() {
			diags.Append(maxPriceDiags...)
		} else {
			data.ConfigMaxPrice = maxPriceValue
		}
	}

	// 8. MaxMemoryPolicyTypeConfiguration -> max_memory
	if apiConfig.MaxMemoryPolicyTypeConfiguration != nil {
		maxMemoryAttrs := map[string]attr.Value{
			"max_memory":         convert.StrToType(&apiConfig.MaxMemoryPolicyTypeConfiguration.MaxMemory),
			"exclude_containers": convert.StringToBool(ctx, apiConfig.MaxMemoryPolicyTypeConfiguration.GetExcludeContainers()),
		}

		maxMemoryValue, maxMemoryDiags := NewConfigMaxMemoryValue(ConfigMaxMemoryValue{}.AttributeTypes(ctx), maxMemoryAttrs)
		if maxMemoryDiags.HasError() {
			diags.Append(maxMemoryDiags...)
		} else {
			data.ConfigMaxMemory = maxMemoryValue
		}
	}

	// 9. MaxCoresPolicyTypeConfiguration -> max_cores
	if apiConfig.MaxCoresPolicyTypeConfiguration != nil {
		maxCoresAttrs := map[string]attr.Value{
			"max_cores":          convert.StrToType(&apiConfig.MaxCoresPolicyTypeConfiguration.MaxCores),
			"exclude_containers": convert.StringToBool(ctx, apiConfig.MaxCoresPolicyTypeConfiguration.GetExcludeContainers()),
		}

		maxCoresValue, maxCoresDiags := NewConfigMaxCoresValue(ConfigMaxCoresValue{}.AttributeTypes(ctx), maxCoresAttrs)
		if maxCoresDiags.HasError() {
			diags.Append(maxCoresDiags...)
		} else {
			data.ConfigMaxCores = maxCoresValue
		}
	}

	// 10. DelayedDeletePolicyTypeConfiguration -> delayed_removal
	if apiConfig.DelayedDeletePolicyTypeConfiguration != nil {
		delayedRemovalAttrs := map[string]attr.Value{
			"removal_age": convert.StrToType(&apiConfig.DelayedDeletePolicyTypeConfiguration.RemovalAge),
		}

		delayedRemovalValue, delayedRemovalDiags := NewConfigDelayedRemovalValue(
			ConfigDelayedRemovalValue{}.AttributeTypes(ctx),
			delayedRemovalAttrs,
		)
		if delayedRemovalDiags.HasError() {
			diags.Append(delayedRemovalDiags...)
		} else {
			data.ConfigDelayedRemoval = delayedRemovalValue
		}
	}

	// 11. ExpirationPolicyTypeConfiguration2 -> lifecycle
	if apiConfig.ExpirationPolicyTypeConfiguration2 != nil {
		lifecycleAttrs := map[string]attr.Value{
			"account_integration_id": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration2.AccountIntegrationId,
			),
			"flow_id":       convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.FlowId),
			"lifecycle_age": convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleAge),
			"lifecycle_allow_extend": convert.StringToBool(
				ctx,
				apiConfig.ExpirationPolicyTypeConfiguration2.GetLifecycleAllowExtend(),
			),
			"lifecycle_auto_renew": convert.StringToBool(
				ctx,
				apiConfig.ExpirationPolicyTypeConfiguration2.GetLifecycleAutoRenew(),
			),
			"lifecycle_extensions_before_approval": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleExtensionsBeforeApproval,
			),
			"lifecycle_hide_fixed": convert.BoolToType(
				apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleHideFixed,
			),
			"lifecycle_message": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleMessage,
			),
			"lifecycle_notify": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleNotify,
			),
			"lifecycle_renewal": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleRenewal,
			),
			"lifecycle_type": convert.StrToType(
				&apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleType,
			),
			"lifecycle_workflow_id": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleWorkflowId,
			),
			"workflow_type": convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.WorkflowType),
		}

		lifecycleValue, lifecycleDiags := NewConfigLifecycleValue(ConfigLifecycleValue{}.AttributeTypes(ctx), lifecycleAttrs)
		if lifecycleDiags.HasError() {
			diags.Append(lifecycleDiags...)
		} else {
			data.ConfigLifecycle = lifecycleValue
		}
	}

	// 12. HostnamePolicyTypeConfiguration -> host_naming
	if apiConfig.HostnamePolicyTypeConfiguration != nil {
		hostNamingAttrs := map[string]attr.Value{
			"host_naming_pattern": convert.StrToType(apiConfig.HostnamePolicyTypeConfiguration.HostNamingPattern),
			"host_naming_type":    convert.StrToType(&apiConfig.HostnamePolicyTypeConfiguration.HostNamingType),
		}

		hostNamingValue, hostNamingDiags := NewConfigHostNamingValue(
			ConfigHostNamingValue{}.AttributeTypes(ctx),
			hostNamingAttrs,
		)
		if hostNamingDiags.HasError() {
			diags.Append(hostNamingDiags...)
		} else {
			data.ConfigHostNaming = hostNamingValue
		}
	}

	// 13. InstanceNamePolicyTypeConfiguration -> naming
	if apiConfig.InstanceNamePolicyTypeConfiguration != nil {
		namingAttrs := map[string]attr.Value{
			"naming_conflict": convert.BoolToType(apiConfig.InstanceNamePolicyTypeConfiguration.NamingConflict),
			"naming_pattern":  convert.StrToType(apiConfig.InstanceNamePolicyTypeConfiguration.NamingPattern),
			"naming_type":     convert.StrToType(&apiConfig.InstanceNamePolicyTypeConfiguration.NamingType),
		}

		namingValue, namingDiags := NewConfigNamingValue(ConfigNamingValue{}.AttributeTypes(ctx), namingAttrs)
		if namingDiags.HasError() {
			diags.Append(namingDiags...)
		} else {
			data.ConfigNaming = namingValue
		}
	}

	// 14. MaxContainersPolicyTypeConfiguration -> max_containers
	if apiConfig.MaxContainersPolicyTypeConfiguration != nil {
		maxContainersAttrs := map[string]attr.Value{
			"max_containers": convert.StrToType(&apiConfig.MaxContainersPolicyTypeConfiguration.MaxContainers),
		}

		maxContainersValue, maxContainersDiags := NewConfigMaxContainersValue(
			ConfigMaxContainersValue{}.AttributeTypes(ctx),
			maxContainersAttrs,
		)
		if maxContainersDiags.HasError() {
			diags.Append(maxContainersDiags...)
		} else {
			data.ConfigMaxContainers = maxContainersValue
		}
	}

	// 15. MaxHostsPolicyTypeConfiguration -> max_hosts
	if apiConfig.MaxHostsPolicyTypeConfiguration != nil {
		maxHostsAttrs := map[string]attr.Value{
			"max_hosts": convert.StrToType(&apiConfig.MaxHostsPolicyTypeConfiguration.MaxHosts),
		}

		maxHostsValue, maxHostsDiags := NewConfigMaxHostsValue(ConfigMaxHostsValue{}.AttributeTypes(ctx), maxHostsAttrs)
		if maxHostsDiags.HasError() {
			diags.Append(maxHostsDiags...)
		} else {
			data.ConfigMaxHosts = maxHostsValue
		}
	}

	// 16. NetworkQuotaPolicyTypeConfiguration -> max_networks
	if apiConfig.NetworkQuotaPolicyTypeConfiguration != nil {
		maxNetworksAttrs := map[string]attr.Value{
			"max_networks": convert.StrToType(&apiConfig.NetworkQuotaPolicyTypeConfiguration.MaxNetworks),
		}

		maxNetworksValue, maxNetworksDiags := NewConfigMaxNetworksValue(
			ConfigMaxNetworksValue{}.AttributeTypes(ctx),
			maxNetworksAttrs,
		)
		if maxNetworksDiags.HasError() {
			diags.Append(maxNetworksDiags...)
		} else {
			data.ConfigMaxNetworks = maxNetworksValue
		}
	}

	// 17. MaxPoolMembersPolicyTypeConfiguration -> max_pool_members
	if apiConfig.MaxPoolMembersPolicyTypeConfiguration != nil {
		maxPoolMembersAttrs := map[string]attr.Value{
			"max_pool_members": convert.StrToType(&apiConfig.MaxPoolMembersPolicyTypeConfiguration.MaxPoolMembers),
		}

		maxPoolMembersValue, maxPoolMembersDiags := NewConfigMaxPoolMembersValue(
			ConfigMaxPoolMembersValue{}.AttributeTypes(ctx),
			maxPoolMembersAttrs,
		)
		if maxPoolMembersDiags.HasError() {
			diags.Append(maxPoolMembersDiags...)
		} else {
			data.ConfigMaxPoolMembers = maxPoolMembersValue
		}
	}

	// 18. MaxLoadBalancerPoolsPolicyTypeConfiguration -> max_pools
	if apiConfig.MaxLoadBalancerPoolsPolicyTypeConfiguration != nil {
		maxPoolsAttrs := map[string]attr.Value{
			"max_pools": convert.StrToType(&apiConfig.MaxLoadBalancerPoolsPolicyTypeConfiguration.MaxPools),
		}

		maxPoolsValue, maxPoolsDiags := NewConfigMaxPoolsValue(ConfigMaxPoolsValue{}.AttributeTypes(ctx), maxPoolsAttrs)
		if maxPoolsDiags.HasError() {
			diags.Append(maxPoolsDiags...)
		} else {
			data.ConfigMaxPools = maxPoolsValue
		}
	}

	// 19. RouterQuotaPolicyTypeConfiguration -> max_routers
	if apiConfig.RouterQuotaPolicyTypeConfiguration != nil {
		maxRoutersAttrs := map[string]attr.Value{
			"max_routers": convert.StrToType(&apiConfig.RouterQuotaPolicyTypeConfiguration.MaxRouters),
		}

		maxRoutersValue, maxRoutersDiags := NewConfigMaxRoutersValue(
			ConfigMaxRoutersValue{}.AttributeTypes(ctx),
			maxRoutersAttrs,
		)
		if maxRoutersDiags.HasError() {
			diags.Append(maxRoutersDiags...)
		} else {
			data.ConfigMaxRouters = maxRoutersValue
		}
	}

	// 20. MaxSnapshotsPolicyTypeConfiguration -> max_snapshots
	if apiConfig.MaxSnapshotsPolicyTypeConfiguration != nil {
		maxSnapshotsAttrs := map[string]attr.Value{
			"max_snapshots": convert.StrToType(&apiConfig.MaxSnapshotsPolicyTypeConfiguration.MaxSnapshots),
		}

		maxSnapshotsValue, maxSnapshotsDiags := NewConfigMaxSnapshotsValue(
			ConfigMaxSnapshotsValue{}.AttributeTypes(ctx),
			maxSnapshotsAttrs,
		)
		if maxSnapshotsDiags.HasError() {
			diags.Append(maxSnapshotsDiags...)
		} else {
			data.ConfigMaxSnapshots = maxSnapshotsValue
		}
	}

	// 21. MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration -> max_storage
	if apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration != nil {
		maxStorageAttrs := map[string]attr.Value{
			"exclude_containers": convert.StringToBool(
				ctx,
				apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration.GetExcludeContainers(),
			),
			"max_storage": convert.StrToType(
				&apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration.MaxStorage,
			),
		}

		maxStorageValue, maxStorageDiags := NewConfigMaxStorageValue(
			ConfigMaxStorageValue{}.AttributeTypes(ctx),
			maxStorageAttrs,
		)
		if maxStorageDiags.HasError() {
			diags.Append(maxStorageDiags...)
		} else {
			data.ConfigMaxStorage = maxStorageValue
		}
	}

	// 22. MaxVirtualServersPolicyTypeConfiguration -> max_virtual_servers
	if apiConfig.MaxVirtualServersPolicyTypeConfiguration != nil {
		maxVirtualServersAttrs := map[string]attr.Value{
			"max_virtual_servers": convert.StrToType(&apiConfig.MaxVirtualServersPolicyTypeConfiguration.MaxVirtualServers),
		}

		maxVirtualServersValue, maxVirtualServersDiags := NewConfigMaxVirtualServersValue(
			ConfigMaxVirtualServersValue{}.AttributeTypes(ctx),
			maxVirtualServersAttrs,
		)
		if maxVirtualServersDiags.HasError() {
			diags.Append(maxVirtualServersDiags...)
		} else {
			data.ConfigMaxVirtualServers = maxVirtualServersValue
		}
	}

	// 23. MaxVMsPolicyTypeConfiguration -> max_vms
	if apiConfig.MaxVMsPolicyTypeConfiguration != nil {
		maxVmsAttrs := map[string]attr.Value{
			"max_vms": convert.StrToType(&apiConfig.MaxVMsPolicyTypeConfiguration.MaxVms),
		}

		maxVmsValue, maxVmsDiags := NewConfigMaxVmsValue(ConfigMaxVmsValue{}.AttributeTypes(ctx), maxVmsAttrs)
		if maxVmsDiags.HasError() {
			diags.Append(maxVmsDiags...)
		} else {
			data.ConfigMaxVms = maxVmsValue
		}
	}

	// 24. MessageOfTheDayPolicyTypeConfiguration2 -> motd
	if apiConfig.MessageOfTheDayPolicyTypeConfiguration2 != nil {
		motdAttrs := map[string]attr.Value{
			"motd_fullpage": types.StringNull(),
			"motddate":      convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdDate),
			"motdmessage":   convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdMessage),
			"motdtitle":     types.StringNull(), // NullableString type
			"motdtype":      convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdType),
		}

		// Handle NullableString for MotdFullPage - check if field exists and is not nil
		if apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdFullPage != nil {
			motdAttrs["motd_fullpage"] = convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdFullPage.String)
		}
		// Handle NullableString for MotdTitle
		if apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdTitle.IsSet() {
			motdAttrs["motdtitle"] = types.StringValue(*apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdTitle.Get())
		}

		motdValue, motdDiags := NewConfigMotdValue(ConfigMotdValue{}.AttributeTypes(ctx), motdAttrs)
		if motdDiags.HasError() {
			diags.Append(motdDiags...)
		} else {
			data.ConfigMotd = motdValue
		}
	}

	// 25. PowerSchedulePolicyTypeConfiguration -> power_schedule
	if apiConfig.PowerSchedulePolicyTypeConfiguration != nil {
		powerScheduleAttrs := map[string]attr.Value{
			"power_schedule": convert.StrToType(
				apiConfig.PowerSchedulePolicyTypeConfiguration.PowerSchedule,
			),
			"power_schedule_hide_fixed": convert.BoolToType(
				apiConfig.PowerSchedulePolicyTypeConfiguration.PowerScheduleHideFixed,
			),
			"power_schedule_type": convert.StrToType(
				&apiConfig.PowerSchedulePolicyTypeConfiguration.PowerScheduleType,
			),
		}

		powerScheduleValue, powerScheduleDiags := NewConfigPowerScheduleValue(
			ConfigPowerScheduleValue{}.AttributeTypes(ctx),
			powerScheduleAttrs,
		)
		if powerScheduleDiags.HasError() {
			diags.Append(powerScheduleDiags...)
		} else {
			data.ConfigPowerSchedule = powerScheduleValue
		}
	}

	// 26. RequiredNetworkPolicyTypeConfiguration -> required_network
	if apiConfig.RequiredNetworkPolicyTypeConfiguration != nil {
		// Handle RequiredNetworks as a set of integers
		var requiredNetworksSet types.Set
		if len(apiConfig.RequiredNetworkPolicyTypeConfiguration.RequiredNetworks) == 0 {
			requiredNetworksSet = types.SetValueMust(types.Int64Type, []attr.Value{})
		} else {
			int64Values := make([]attr.Value, len(apiConfig.RequiredNetworkPolicyTypeConfiguration.RequiredNetworks))
			for i, networkId := range apiConfig.RequiredNetworkPolicyTypeConfiguration.RequiredNetworks {
				int64Values[i] = types.Int64Value(networkId)
			}
			var setDiags diag.Diagnostics
			requiredNetworksSet, setDiags = types.SetValueFrom(ctx, types.Int64Type, int64Values)
			if setDiags.HasError() {
				diags.Append(setDiags...)
			}
		}

		requiredNetworkAttrs := map[string]attr.Value{
			"required_networks": requiredNetworksSet,
		}

		requiredNetworkValue, requiredNetworkDiags := NewConfigRequiredNetworkValue(
			ConfigRequiredNetworkValue{}.AttributeTypes(ctx),
			requiredNetworkAttrs,
		)
		if requiredNetworkDiags.HasError() {
			diags.Append(requiredNetworkDiags...)
		} else {
			data.ConfigRequiredNetwork = requiredNetworkValue
		}
	}

	// 27. ClusterResourceNamePolicyTypeConfiguration -> server_naming
	if apiConfig.ClusterResourceNamePolicyTypeConfiguration != nil {
		serverNamingAttrs := map[string]attr.Value{
			"server_naming_conflict": convert.BoolToType(
				apiConfig.ClusterResourceNamePolicyTypeConfiguration.ServerNamingConflict,
			),
			"server_naming_pattern": convert.StrToType(
				apiConfig.ClusterResourceNamePolicyTypeConfiguration.ServerNamingPattern,
			),
			"server_naming_type": convert.StrToType(
				&apiConfig.ClusterResourceNamePolicyTypeConfiguration.ServerNamingType,
			),
		}

		serverNamingValue, serverNamingDiags := NewConfigServerNamingValue(
			ConfigServerNamingValue{}.AttributeTypes(ctx),
			serverNamingAttrs,
		)
		if serverNamingDiags.HasError() {
			diags.Append(serverNamingDiags...)
		} else {
			data.ConfigServerNaming = serverNamingValue
		}
	}

	// 28. ShutdownPolicyTypeConfiguration -> shutdown
	if apiConfig.ShutdownPolicyTypeConfiguration != nil {
		shutdownAttrs := map[string]attr.Value{
			"account_integration_id": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration.AccountIntegrationId,
			),
			"flow_id":      convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.FlowId),
			"shutdown_age": convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.ShutdownAge),
			"shutdown_allow_extend": convert.StringToBool(
				ctx,
				apiConfig.ShutdownPolicyTypeConfiguration.GetShutdownAllowExtend(),
			),
			"shutdown_auto_renew": convert.StringToBool(
				ctx,
				apiConfig.ShutdownPolicyTypeConfiguration.GetShutdownAutoRenew(),
			),
			"shutdown_extensions_before_approval": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration.ShutdownExtensionsBeforeApproval,
			),
			"shutdown_hide_fixed": convert.BoolToType(
				apiConfig.ShutdownPolicyTypeConfiguration.ShutdownHideFixed,
			),
			"shutdown_message": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration.ShutdownMessage,
			),
			"shutdown_notify": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration.ShutdownNotify,
			),
			"shutdown_renewal": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration.ShutdownRenewal,
			),
			"shutdown_type": convert.StrToType(
				&apiConfig.ShutdownPolicyTypeConfiguration.ShutdownType,
			),
			"shutdown_workflow_id": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration.ShutdownWorkflowId,
			),
			"workflow_type": convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.WorkflowType),
		}

		shutdownValue, shutdownDiags := NewConfigShutdownValue(ConfigShutdownValue{}.AttributeTypes(ctx), shutdownAttrs)
		if shutdownDiags.HasError() {
			diags.Append(shutdownDiags...)
		} else {
			data.ConfigShutdown = shutdownValue
		}
	}

	// 29. StorageServerStorageQuotaPolicyTypeConfiguration -> storage_server_quota
	if apiConfig.StorageServerStorageQuotaPolicyTypeConfiguration != nil {
		storageServerQuotaAttrs := map[string]attr.Value{
			"max_storage":       convert.StrToType(apiConfig.StorageServerStorageQuotaPolicyTypeConfiguration.MaxStorage),
			"storage_server_id": convert.StrToType(&apiConfig.StorageServerStorageQuotaPolicyTypeConfiguration.StorageServerId),
		}

		storageServerQuotaValue, storageServerQuotaDiags := NewConfigStorageServerQuotaValue(
			ConfigStorageServerQuotaValue{}.AttributeTypes(ctx),
			storageServerQuotaAttrs,
		)
		if storageServerQuotaDiags.HasError() {
			diags.Append(storageServerQuotaDiags...)
		} else {
			data.ConfigStorageServerQuota = storageServerQuotaValue
		}
	}

	// 30. TagsPolicyTypeConfiguration -> tags
	if apiConfig.TagsPolicyTypeConfiguration != nil {
		tagsAttrs := map[string]attr.Value{
			"key":           convert.StrToType(apiConfig.TagsPolicyTypeConfiguration.Key),
			"strict":        types.BoolValue(apiConfig.TagsPolicyTypeConfiguration.Strict),
			"value":         convert.StrToType(apiConfig.TagsPolicyTypeConfiguration.Value),
			"value_list_id": convert.StrToType(apiConfig.TagsPolicyTypeConfiguration.ValueListId),
		}

		tagsValue, tagsDiags := NewConfigTagsValue(ConfigTagsValue{}.AttributeTypes(ctx), tagsAttrs)
		if tagsDiags.HasError() {
			diags.Append(tagsDiags...)
		} else {
			data.ConfigTags = tagsValue
		}
	}

	// 31. WorkflowPolicyTypeConfiguration -> workflow
	if apiConfig.WorkflowPolicyTypeConfiguration != nil {
		workflowAttrs := map[string]attr.Value{
			"workflow_id": convert.StrToType(&apiConfig.WorkflowPolicyTypeConfiguration.WorkflowId),
		}

		workflowValue, workflowDiags := NewConfigWorkflowValue(ConfigWorkflowValue{}.AttributeTypes(ctx), workflowAttrs)
		if workflowDiags.HasError() {
			diags.Append(workflowDiags...)
		} else {
			data.ConfigWorkflow = workflowValue
		}
	}

	return diags
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
	resp.TypeName = req.ProviderTypeName + "_morpheus_policy"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = PolicyDataSourceSchema(ctx)
}

func getPolicyByName(
	ctx context.Context,
	data *PolicyModel,
	apiClient *sdk.APIClient,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	name := data.Name.ValueString()
	ps, hresp, err := apiClient.PoliciesAPI.ListPolicies(ctx).Name(name).Execute()
	if ps == nil || err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(summary, fmt.Sprintf("GET failed for policy with name %s: %s",
			name, internalErrors.ErrMsg(err, hresp)))

		return diags
	}

	policies := ps.GetPolicies()

	// Additional filtering to ensure exact name match (API might return partial matches)
	var filteredPolicies []sdk.ListPolicies200ResponseAllOfPoliciesInner
	for _, p := range policies {
		if p.GetName() == data.Name.ValueString() {
			filteredPolicies = append(filteredPolicies, p)
		}
	}
	policies = filteredPolicies

	if len(policies) > 1 {
		diags.AddError(summary, consts.ErrorMultiplePolicies)

		return diags
	} else if len(policies) == 0 {
		diags.AddError(summary, consts.ErrorNoPolicyFound)

		return diags
	}

	policy := policies[0]

	return getPolicyByID(ctx, policy.GetId(), data, apiClient)
}

func getPolicyByID(
	ctx context.Context,
	id int64,
	data *PolicyModel,
	apiClient *sdk.APIClient,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	p, hresp, err := apiClient.PoliciesAPI.GetPolicies(ctx, id).Execute()
	if p == nil || err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(summary, fmt.Sprintf("GET failed for policy with id %d: %s",
			id, internalErrors.ErrMsg(err, hresp)))

		return diags
	}
	policy, ok := p.GetPolicyOk()
	if !ok {
		diags.AddError(summary, consts.ErrorNoPolicyFound)

		return diags
	}

	// Map basic policy fields
	data.Id = convert.Int64ToType(policy.Id)
	data.Name = convert.StrToType(policy.Name)
	data.Description = convert.StrToType(policy.Description.Get())
	data.Enabled = convert.BoolToType(policy.Enabled)
	data.EachUser = convert.BoolToType(policy.EachUser.Get())

	// Handle AssociatedResourceId and AssociatedResourceType
	if policy.RefId.IsSet() && policy.RefId.Get() != nil {
		data.AssociatedResourceId = types.Int64Value(policy.GetRefId())
	} else {
		data.AssociatedResourceId = types.Int64Null()
	}

	if policy.RefType.IsSet() && policy.RefType.Get() != nil {
		apiType := *policy.RefType.Get()
		data.AssociatedResourceType = types.StringValue(apiTypeToResourceType(apiType))
	} else {
		data.AssociatedResourceType = types.StringValue("Global")
	}

	// Handle PolicyType
	if policy.PolicyType != nil {
		policyTypeValue, policyTypeDiags := NewPolicyTypeValue(
			PolicyTypeValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(policy.PolicyType.Id),
				"code": convert.StrToType(policy.PolicyType.Code),
				"name": convert.StrToType(policy.PolicyType.Name),
			},
		)
		if policyTypeDiags.HasError() {
			diags.Append(policyTypeDiags...)

			return diags
		}
		data.PolicyType = policyTypeValue
	} else {
		data.PolicyType = NewPolicyTypeValueNull()
	}

	// Handle Config - map to structured schema fields
	// data.Config = types.DynamicNull()
	if policy.Config != nil {
		// Map API config to structured schema fields
		configDiags := mapPolicyConfigToState(ctx, data, policy.Config)
		if configDiags.HasError() {
			diags.Append(configDiags...)

			return diags
		}

		// Also convert API config to dynamic type
		var err error
		data.Config, err = convert.StructToDynamic(ctx, policy.Config)
		if err != nil {
			diags.AddError(
				summary,
				fmt.Sprintf("policy %d: failed to convert config: %s", id, err.Error()),
			)

			return diags
		}
	}

	// Handle Cloud (Zone)
	if policy.Zone != nil {
		cloudValue, cloudDiags := NewCloudValue(
			CloudValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(policy.Zone.Id),
				"name": convert.StrToType(policy.Zone.Name),
			},
		)
		if cloudDiags.HasError() {
			diags.Append(cloudDiags...)

			return diags
		}
		data.Cloud = cloudValue
	} else {
		data.Cloud = NewCloudValueNull()
	}

	// Handle Group (Site)
	if policy.Site != nil {
		groupValue, groupDiags := NewGroupValue(
			GroupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(policy.Site.Id),
				"name": convert.StrToType(policy.Site.Name),
			},
		)
		if groupDiags.HasError() {
			diags.Append(groupDiags...)

			return diags
		}
		data.Group = groupValue
	} else {
		data.Group = NewGroupValueNull()
	}

	// Handle Owner
	if policy.Owner.IsSet() && policy.Owner.Get() != nil {
		owner := policy.Owner.Get()
		ownerValue, ownerDiags := NewOwnerValue(
			OwnerValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(owner.Id),
				"name": convert.StrToType(owner.Name),
			},
		)
		if ownerDiags.HasError() {
			diags.Append(ownerDiags...)

			return diags
		}
		data.Owner = ownerValue
	} else {
		data.Owner = NewOwnerValueNull()
	}

	// Handle Role
	if policy.Role != nil {
		roleValue, roleDiags := NewRoleValue(
			RoleValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":        convert.Int64ToType(policy.Role.Id),
				"authority": convert.StrToType(policy.Role.Authority),
			},
		)
		if roleDiags.HasError() {
			diags.Append(roleDiags...)

			return diags
		}
		data.Role = roleValue
	} else {
		data.Role = NewRoleValueNull()
	}

	// Handle User
	if policy.User != nil {
		userValue, userDiags := NewUserValue(
			UserValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":       convert.Int64ToType(policy.User.Id),
				"username": convert.StrToType(policy.User.Username),
			},
		)
		if userDiags.HasError() {
			diags.Append(userDiags...)

			return diags
		}
		data.User = userValue
	} else {
		data.User = NewUserValueNull()
	}

	// Handle Tenants (Accounts)
	if len(policy.Accounts) > 0 {
		tenantValues := []attr.Value{}
		tenantObjectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.Int64Type,
				"name": types.StringType,
			},
		}
		for _, account := range policy.Accounts {
			tenantAttrs := map[string]attr.Value{
				"id":   convert.Int64ToType(account.Id),
				"name": convert.StrToType(account.Name),
			}
			tenantValue, tenantDiags := types.ObjectValue(tenantObjectType.AttrTypes, tenantAttrs)
			if tenantDiags.HasError() {
				diags.Append(tenantDiags...)

				return diags
			}
			tenantValues = append(tenantValues, tenantValue)
		}

		tenantsSet, setDiags := types.SetValue(tenantObjectType, tenantValues)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return diags
		}
		data.Tenants = tenantsSet
	} else {
		data.Tenants = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.Int64Type,
				"name": types.StringType,
			},
		})
	}

	return diags
}

func getPolicy(
	ctx context.Context,
	data *PolicyModel,
	apiClient *sdk.APIClient,
) diag.Diagnostics {
	if !data.Id.IsNull() {
		return getPolicyByID(ctx, data.Id.ValueInt64(), data, apiClient)
	} else if !data.Name.IsNull() {
		return getPolicyByName(ctx, data, apiClient)
	}

	diags := diag.Diagnostics{}
	diags.AddError(summary, consts.ErrorNoValidPolicyTerms)

	return diags
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data PolicyModel

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

	diags = getPolicy(ctx, &data, apiClient)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
