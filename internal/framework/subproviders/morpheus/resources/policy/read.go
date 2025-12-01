// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

// mapPolicyConfigToState maps the API config structure to the resource schema structure
func mapPolicyConfigToState(
	ctx context.Context,
	state *PolicyModel,
	apiConfig *sdk.AddPolicies200ResponseAllOfPolicyConfig,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	// Map each API config field to the corresponding schema field - only populate non-null configurations
	// 1. ApprovePolicyTypeConfiguration -> config_approval
	if apiConfig.ApprovePolicyTypeConfiguration != nil {
		approvalValue, approvalDiags := NewConfigApprovalValue(
			ConfigApprovalValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"account_integration_id": convert.StrToType(&apiConfig.ApprovePolicyTypeConfiguration.AccountIntegrationId),
				"flow_id":                convert.StrToType(apiConfig.ApprovePolicyTypeConfiguration.FlowId),
				"workflow_id":            convert.StrToType(apiConfig.ApprovePolicyTypeConfiguration.WorkflowId),
				"workflow_type":          convert.StrToType(apiConfig.ApprovePolicyTypeConfiguration.WorkflowType),
			},
		)
		if approvalDiags.HasError() {
			diags.Append(approvalDiags...)
		} else {
			state.ConfigApproval = approvalValue
		}
	}

	// 2. BackupTargetsPolicyTypeConfiguration -> config_backup_storage
	if apiConfig.BackupTargetsPolicyTypeConfiguration != nil {
		// Handle BackupStorageIds as a set of int64
		var backupStorageIDsSet types.Set
		var setDiags diag.Diagnostics
		if len(apiConfig.BackupTargetsPolicyTypeConfiguration.BackupStorageIds) == 0 {
			backupStorageIDsSet = types.SetValueMust(types.Int64Type, []attr.Value{})
		} else {
			// BackupStorageIds come from API as []string, convert to []int64
			backupStorageIDsSet, setDiags = types.SetValueFrom(
				ctx,
				types.Int64Type,
				apiConfig.BackupTargetsPolicyTypeConfiguration.BackupStorageIds,
			)
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
			state.ConfigBackupStorage = backupStorageValue
		}
	}

	// 3. BackupCreationPolicyTypeConfiguration -> config_create_backup
	if apiConfig.BackupCreationPolicyTypeConfiguration != nil {
		createBackupValue, createBackupDiags := NewConfigCreateBackupValue(
			ConfigCreateBackupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"create_backup":      convert.BoolToType(apiConfig.BackupCreationPolicyTypeConfiguration.CreateBackup),
				"create_backup_type": convert.StrToType(&apiConfig.BackupCreationPolicyTypeConfiguration.CreateBackupType),
			},
		)
		if createBackupDiags.HasError() {
			diags.Append(createBackupDiags...)
		} else {
			state.ConfigCreateBackup = createBackupValue
		}
	}

	// 4. UserCreationPolicyTypeConfiguration -> config_create_user
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
			state.ConfigCreateUser = createUserValue
		}
	}

	// 5. UserGroupCreationPolicyTypeConfiguration -> config_create_user_group
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
			state.ConfigCreateUserGroup = createUserGroupValue
		}
	}

	// 6. CypherAccessPolicyTypeConfiguration -> config_cypher
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
			state.ConfigCypher = cypherValue
		}
	}

	// 7. BudgetPolicyTypeConfiguration -> config_max_price
	if apiConfig.BudgetPolicyTypeConfiguration != nil {
		maxPriceAttrs := map[string]attr.Value{
			"max_price":          convert.StrToNumber(&apiConfig.BudgetPolicyTypeConfiguration.MaxPrice),
			"max_price_currency": convert.StrToType(apiConfig.BudgetPolicyTypeConfiguration.MaxPriceCurrency),
			"max_price_unit":     convert.StrToType(apiConfig.BudgetPolicyTypeConfiguration.MaxPriceUnit),
		}

		maxPriceValue, maxPriceDiags := NewConfigMaxPriceValue(ConfigMaxPriceValue{}.AttributeTypes(ctx), maxPriceAttrs)
		if maxPriceDiags.HasError() {
			diags.Append(maxPriceDiags...)
		} else {
			state.ConfigMaxPrice = maxPriceValue
		}
	}

	// 8. MaxMemoryPolicyTypeConfiguration -> config_max_memory
	if apiConfig.MaxMemoryPolicyTypeConfiguration != nil {
		maxMemoryAttrs := map[string]attr.Value{
			"max_memory":         convert.StrToType(&apiConfig.MaxMemoryPolicyTypeConfiguration.MaxMemory),
			"exclude_containers": convert.StringToBool(ctx, apiConfig.MaxMemoryPolicyTypeConfiguration.GetExcludeContainers()),
		}

		maxMemoryValue, maxMemoryDiags := NewConfigMaxMemoryValue(ConfigMaxMemoryValue{}.AttributeTypes(ctx), maxMemoryAttrs)
		if maxMemoryDiags.HasError() {
			diags.Append(maxMemoryDiags...)
		} else {
			state.ConfigMaxMemory = maxMemoryValue
		}
	}

	// 9. MaxCoresPolicyTypeConfiguration -> config_max_cores
	if apiConfig.MaxCoresPolicyTypeConfiguration != nil {
		maxCoresAttrs := map[string]attr.Value{
			"max_cores":          convert.StrToType(&apiConfig.MaxCoresPolicyTypeConfiguration.MaxCores),
			"exclude_containers": convert.StringToBool(ctx, apiConfig.MaxCoresPolicyTypeConfiguration.GetExcludeContainers()),
		}

		maxCoresValue, maxCoresDiags := NewConfigMaxCoresValue(ConfigMaxCoresValue{}.AttributeTypes(ctx), maxCoresAttrs)
		if maxCoresDiags.HasError() {
			diags.Append(maxCoresDiags...)
		} else {
			state.ConfigMaxCores = maxCoresValue
		}
	}

	// 10. DelayedDeletePolicyTypeConfiguration -> config_delayed_removal
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
			state.ConfigDelayedRemoval = delayedRemovalValue
		}
	}

	// 11. ExpirationPolicyTypeConfiguration2 -> config_lifecycle
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
			state.ConfigLifecycle = lifecycleValue
		}
	}

	// 12. HostnamePolicyTypeConfiguration -> config_host_naming
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
			state.ConfigHostNaming = hostNamingValue
		}
	}

	// 13. InstanceNamePolicyTypeConfiguration -> config_naming
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
			state.ConfigNaming = namingValue
		}
	}

	// 14. MaxContainersPolicyTypeConfiguration -> config_max_containers
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
			state.ConfigMaxContainers = maxContainersValue
		}
	}

	// 15. MaxHostsPolicyTypeConfiguration -> config_max_hosts
	if apiConfig.MaxHostsPolicyTypeConfiguration != nil {
		maxHostsAttrs := map[string]attr.Value{
			"max_hosts": convert.StrToType(&apiConfig.MaxHostsPolicyTypeConfiguration.MaxHosts),
		}

		maxHostsValue, maxHostsDiags := NewConfigMaxHostsValue(ConfigMaxHostsValue{}.AttributeTypes(ctx), maxHostsAttrs)
		if maxHostsDiags.HasError() {
			diags.Append(maxHostsDiags...)
		} else {
			state.ConfigMaxHosts = maxHostsValue
		}
	}

	// 16. NetworkQuotaPolicyTypeConfiguration -> config_max_networks
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
			state.ConfigMaxNetworks = maxNetworksValue
		}
	}

	// 17. MaxPoolMembersPolicyTypeConfiguration -> config_max_pool_members
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
			state.ConfigMaxPoolMembers = maxPoolMembersValue
		}
	}

	// 18. MaxLoadBalancerPoolsPolicyTypeConfiguration -> config_max_pools
	if apiConfig.MaxLoadBalancerPoolsPolicyTypeConfiguration != nil {
		maxPoolsAttrs := map[string]attr.Value{
			"max_pools": convert.StrToType(&apiConfig.MaxLoadBalancerPoolsPolicyTypeConfiguration.MaxPools),
		}

		maxPoolsValue, maxPoolsDiags := NewConfigMaxPoolsValue(ConfigMaxPoolsValue{}.AttributeTypes(ctx), maxPoolsAttrs)
		if maxPoolsDiags.HasError() {
			diags.Append(maxPoolsDiags...)
		} else {
			state.ConfigMaxPools = maxPoolsValue
		}
	}

	// 19. RouterQuotaPolicyTypeConfiguration -> config_max_routers
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
			state.ConfigMaxRouters = maxRoutersValue
		}
	}

	// 20. MaxSnapshotsPolicyTypeConfiguration -> config_max_snapshots
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
			state.ConfigMaxSnapshots = maxSnapshotsValue
		}
	}

	// 21. MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration -> config_max_storage
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
			state.ConfigMaxStorage = maxStorageValue
		}
	}

	// 22. MaxVirtualServersPolicyTypeConfiguration -> config_max_virtual_servers
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
			state.ConfigMaxVirtualServers = maxVirtualServersValue
		}
	}

	// 23. MaxVMsPolicyTypeConfiguration -> config_max_vms
	if apiConfig.MaxVMsPolicyTypeConfiguration != nil {
		maxVmsAttrs := map[string]attr.Value{
			"max_vms": convert.StrToType(&apiConfig.MaxVMsPolicyTypeConfiguration.MaxVms),
		}

		maxVmsValue, maxVmsDiags := NewConfigMaxVmsValue(ConfigMaxVmsValue{}.AttributeTypes(ctx), maxVmsAttrs)
		if maxVmsDiags.HasError() {
			diags.Append(maxVmsDiags...)
		} else {
			state.ConfigMaxVms = maxVmsValue
		}
	}

	// 24. MessageOfTheDayPolicyTypeConfiguration2 -> config_motd
	if apiConfig.MessageOfTheDayPolicyTypeConfiguration2 != nil {
		motdAttrs := map[string]attr.Value{
			"motddate":    convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdDate),
			"motdmessage": convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdMessage),
			"motdtitle":   types.StringNull(),
			"motdtype":    convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdType),
		}

		// Handle NullableString for MotdTitle
		if apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdTitle.IsSet() {
			motdAttrs["motdtitle"] = types.StringValue(*apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdTitle.Get())
		}

		motdValue, motdDiags := NewConfigMotdValue(ConfigMotdValue{}.AttributeTypes(ctx), motdAttrs)
		if motdDiags.HasError() {
			diags.Append(motdDiags...)
		} else {
			state.ConfigMotd = motdValue
		}
	}

	// 25. PowerSchedulePolicyTypeConfiguration -> config_power_schedule
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
			state.ConfigPowerSchedule = powerScheduleValue
		}
	}

	// 26. RequiredNetworkPolicyTypeConfiguration -> config_required_network
	if apiConfig.RequiredNetworkPolicyTypeConfiguration != nil {
		// Handle RequiredNetworks as a set of integers
		var requiredNetworksSet types.Set
		if len(apiConfig.RequiredNetworkPolicyTypeConfiguration.RequiredNetworks) == 0 {
			requiredNetworksSet = types.SetValueMust(types.Int64Type, []attr.Value{})
		} else {
			int64Values := make([]attr.Value, len(apiConfig.RequiredNetworkPolicyTypeConfiguration.RequiredNetworks))
			for i, networkID := range apiConfig.RequiredNetworkPolicyTypeConfiguration.RequiredNetworks {
				int64Values[i] = types.Int64Value(networkID)
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
			state.ConfigRequiredNetwork = requiredNetworkValue
		}
	}

	// 27. ClusterResourceNamePolicyTypeConfiguration -> config_server_naming
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
			state.ConfigServerNaming = serverNamingValue
		}
	}

	// 28. ShutdownPolicyTypeConfiguration -> config_shutdown
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
			state.ConfigShutdown = shutdownValue
		}
	}

	// 29. StorageServerStorageQuotaPolicyTypeConfiguration -> config_storage_server_quota
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
			state.ConfigStorageServerQuota = storageServerQuotaValue
		}
	}

	// 30. TagsPolicyTypeConfiguration -> config_tags
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
			state.ConfigTags = tagsValue
		}
	}

	// 31. WorkflowPolicyTypeConfiguration -> config_workflow
	if apiConfig.WorkflowPolicyTypeConfiguration != nil {
		workflowAttrs := map[string]attr.Value{
			"workflow_id": convert.StrToType(&apiConfig.WorkflowPolicyTypeConfiguration.WorkflowId),
		}

		workflowValue, workflowDiags := NewConfigWorkflowValue(ConfigWorkflowValue{}.AttributeTypes(ctx), workflowAttrs)
		if workflowDiags.HasError() {
			diags.Append(workflowDiags...)
		} else {
			state.ConfigWorkflow = workflowValue
		}
	}

	return diags
}

// populate policy resource model with current API values
func getPolicyAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan *PolicyModel,
) (PolicyModel, diag.Diagnostics) {
	var state PolicyModel
	var diags diag.Diagnostics

	policy, hresp, err := client.PoliciesAPI.GetPolicies(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK || policy == nil {
		diags.AddError(
			"populate policy resource",
			fmt.Sprintf("policy %d GET failed", id)+errors.ErrMsg(err, hresp),
		)

		return state, diags
	}

	if policy.Policy == nil {
		diags.AddError(
			"populate policy resource",
			fmt.Sprintf("policy %d is nil", id),
		)

		return state, diags
	}

	p := policy.Policy

	// Set basic fields
	state.Id = convert.Int64ToType(p.Id)
	state.Name = convert.StrToType(p.Name)

	// Handle nullable fields properly
	if p.Description.IsSet() {
		state.Description = convert.StrToType(p.Description.Get())
	}

	state.Enabled = convert.BoolToType(p.Enabled)

	if p.EachUser.IsSet() {
		state.EachUser = convert.BoolToType(p.EachUser.Get())
	}

	// Handle RefId - convert string to int64
	if p.RefId.IsSet() && p.RefId.Get() != nil {
		state.AssociatedResourceId = types.Int64Value(p.GetRefId())
	}

	// Handle RefType
	// If RefType is null or not set, it means it's a Global policy
	if p.RefType.IsSet() && p.RefType.Get() != nil {
		apiType := *p.RefType.Get()
		// Convert API type to user-facing resource type
		state.AssociatedResourceType = types.StringValue(apiTypeToResourceType(apiType))
	} else {
		state.AssociatedResourceType = types.StringValue(AssociatedResourceTypeGlobal)
	}

	// Set Tenant IDs
	if len(p.Accounts) > 0 {
		tenantIDs := make([]int64, 0, len(p.Accounts))
		for _, acc := range p.Accounts {
			if acc.Id != nil {
				tenantIDs = append(tenantIDs, *acc.Id)
			}
		}
		tenantsSet, setDiags := types.SetValueFrom(ctx, types.Int64Type, tenantIDs)
		if setDiags.HasError() {
			diags.Append(setDiags...)

			return state, diags
		}
		state.Tenants = tenantsSet
	} else {
		// Set to null if no accounts
		state.Tenants = types.SetNull(types.Int64Type)
	}

	// Set PolicyType
	if p.PolicyType != nil {
		policyTypeAttrs := map[string]attr.Value{}
		if p.PolicyType.Code != nil {
			policyTypeAttrs["code"] = types.StringValue(*p.PolicyType.Code)
		} else {
			policyTypeAttrs["code"] = types.StringNull()
		}
		if p.PolicyType.Id != nil {
			policyTypeAttrs["id"] = types.Int64Value(*p.PolicyType.Id)
		} else {
			policyTypeAttrs["id"] = types.Int64Null()
		}
		if p.PolicyType.Name != nil {
			policyTypeAttrs["name"] = types.StringValue(*p.PolicyType.Name)
		} else {
			policyTypeAttrs["name"] = types.StringNull()
		}

		policyTypeValue, policyTypeDiags := NewPolicyTypeValue(
			PolicyTypeValue{}.AttributeTypes(ctx), policyTypeAttrs)
		if policyTypeDiags.HasError() {
			diags.Append(policyTypeDiags...)

			return state, diags
		}
		state.PolicyType = policyTypeValue
	}

	// Initialize all config_* fields to null by default
	state.ConfigApproval = NewConfigApprovalValueNull()
	state.ConfigBackupStorage = NewConfigBackupStorageValueNull()
	state.ConfigCreateBackup = NewConfigCreateBackupValueNull()
	state.ConfigCreateUser = NewConfigCreateUserValueNull()
	state.ConfigCreateUserGroup = NewConfigCreateUserGroupValueNull()
	state.ConfigCypher = NewConfigCypherValueNull()
	state.ConfigDelayedRemoval = NewConfigDelayedRemovalValueNull()
	state.ConfigHostNaming = NewConfigHostNamingValueNull()
	state.ConfigLifecycle = NewConfigLifecycleValueNull()
	state.ConfigMaxContainers = NewConfigMaxContainersValueNull()
	state.ConfigMaxCores = NewConfigMaxCoresValueNull()
	state.ConfigMaxHosts = NewConfigMaxHostsValueNull()
	state.ConfigMaxMemory = NewConfigMaxMemoryValueNull()
	state.ConfigMaxNetworks = NewConfigMaxNetworksValueNull()
	state.ConfigMaxPoolMembers = NewConfigMaxPoolMembersValueNull()
	state.ConfigMaxPools = NewConfigMaxPoolsValueNull()
	state.ConfigMaxPrice = NewConfigMaxPriceValueNull()
	state.ConfigMaxRouters = NewConfigMaxRoutersValueNull()
	state.ConfigMaxSnapshots = NewConfigMaxSnapshotsValueNull()
	state.ConfigMaxStorage = NewConfigMaxStorageValueNull()
	state.ConfigMaxVirtualServers = NewConfigMaxVirtualServersValueNull()
	state.ConfigMaxVms = NewConfigMaxVmsValueNull()
	state.ConfigMotd = NewConfigMotdValueNull()
	state.ConfigNaming = NewConfigNamingValueNull()
	state.ConfigPowerSchedule = NewConfigPowerScheduleValueNull()
	state.ConfigRequiredNetwork = NewConfigRequiredNetworkValueNull()
	state.ConfigServerNaming = NewConfigServerNamingValueNull()
	state.ConfigShutdown = NewConfigShutdownValueNull()
	state.ConfigStorageServerQuota = NewConfigStorageServerQuotaValueNull()
	state.ConfigTags = NewConfigTagsValueNull()
	state.ConfigWorkflow = NewConfigWorkflowValueNull()

	// Handle Config - use static schema fields when available, fallback to dynamic
	if p.Config != nil {
		// Map API config to static schema fields
		configDiags := mapPolicyConfigToState(ctx, &state, p.Config)
		if configDiags.HasError() {
			diags.Append(configDiags...)

			return state, diags
		}

		// Also preserve the dynamic config field if it was set in plan
		if plan != nil && !plan.Config.IsNull() && !plan.Config.IsUnknown() {
			state.Config = plan.Config
		} else {
			// Convert API config to dynamic type as fallback
			var err error
			state.Config, err = convert.StructToDynamic(ctx, p.Config)
			if err != nil {
				diags.AddError(
					"populate policy resource",
					fmt.Sprintf("policy %d: failed to convert config: %s", id, err.Error()),
				)

				return state, diags
			}
		}
	}

	// Computed types
	// Set Cloud if present
	if p.Zone != nil {
		cloudAttrs := map[string]attr.Value{}
		if p.Zone.Id != nil {
			cloudAttrs["id"] = types.Int64Value(*p.Zone.Id)
		} else {
			cloudAttrs["id"] = types.Int64Null()
		}
		if p.Zone.Name != nil {
			cloudAttrs["name"] = types.StringValue(*p.Zone.Name)
		} else {
			cloudAttrs["name"] = types.StringNull()
		}

		cloudValue, cloudDiags := NewCloudValue(CloudValue{}.AttributeTypes(ctx), cloudAttrs)
		if cloudDiags.HasError() {
			diags.Append(cloudDiags...)

			return state, diags
		}
		state.Cloud = cloudValue
	} else {
		state.Cloud = NewCloudValueNull()
	}

	// Set Group if present
	if p.Site != nil {
		groupAttrs := map[string]attr.Value{}
		if p.Site.Id != nil {
			groupAttrs["id"] = types.Int64Value(*p.Site.Id)
		} else {
			groupAttrs["id"] = types.Int64Null()
		}
		if p.Site.Name != nil {
			groupAttrs["name"] = types.StringValue(*p.Site.Name)
		} else {
			groupAttrs["name"] = types.StringNull()
		}

		groupValue, groupDiags := NewGroupValue(GroupValue{}.AttributeTypes(ctx), groupAttrs)
		if groupDiags.HasError() {
			diags.Append(groupDiags...)

			return state, diags
		}
		state.Group = groupValue
	} else {
		state.Group = NewGroupValueNull()
	}

	// Set Owner if present
	state.Owner = NewOwnerValueNull()
	if p.Owner.IsSet() && p.Owner.Get() != nil {
		owner := p.Owner.Get()
		ownerAttrs := map[string]attr.Value{}
		if owner.Id != nil {
			ownerAttrs["id"] = types.Int64Value(*owner.Id)
		} else {
			ownerAttrs["id"] = types.Int64Null()
		}
		if owner.Name != nil {
			ownerAttrs["name"] = types.StringValue(*owner.Name)
		} else {
			ownerAttrs["name"] = types.StringNull()
		}

		ownerValue, ownerDiags := NewOwnerValue(OwnerValue{}.AttributeTypes(ctx), ownerAttrs)
		if ownerDiags.HasError() {
			diags.Append(ownerDiags...)

			return state, diags
		}
		state.Owner = ownerValue
	}

	// Set Role if present
	if p.Role != nil {
		roleAttrs := map[string]attr.Value{}
		if p.Role.Id != nil {
			roleAttrs["id"] = types.Int64Value(*p.Role.Id)
		} else {
			roleAttrs["id"] = types.Int64Null()
		}
		if p.Role.Authority != nil {
			roleAttrs["authority"] = types.StringValue(*p.Role.Authority)
		} else {
			roleAttrs["authority"] = types.StringNull()
		}

		roleValue, roleDiags := NewRoleValue(RoleValue{}.AttributeTypes(ctx), roleAttrs)
		if roleDiags.HasError() {
			diags.Append(roleDiags...)

			return state, diags
		}
		state.Role = roleValue
	} else {
		state.Role = NewRoleValueNull()
	}

	// Set User if present
	if p.User != nil {
		userAttrs := map[string]attr.Value{}
		if p.User.Id != nil {
			userAttrs["id"] = types.Int64Value(*p.User.Id)
		} else {
			userAttrs["id"] = types.Int64Null()
		}
		if p.User.Username != nil {
			userAttrs["username"] = types.StringValue(*p.User.Username)
		} else {
			userAttrs["username"] = types.StringNull()
		}

		userValue, userDiags := NewUserValue(UserValue{}.AttributeTypes(ctx), userAttrs)
		if userDiags.HasError() {
			diags.Append(userDiags...)

			return state, diags
		}
		state.User = userValue
	} else {
		state.User = NewUserValueNull()
	}

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan PolicyModel

	diags := req.State.Get(ctx, &plan)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read policy resource",
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	state, diags := getPolicyAsState(ctx, id, client, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
