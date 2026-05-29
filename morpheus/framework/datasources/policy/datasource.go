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

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/policy/consts"
	internalErrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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
	apiConfig *sdk.GetPolicies200ResponseAllOfPolicyConfig,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	// Map each API config field to the corresponding schema field - only populate non-null configurations
	// 1. ApprovePolicyTypeConfiguration2 -> approval
	if apiConfig.ApprovePolicyTypeConfiguration3 != nil {
		approvalValue, approvalDiags := NewConfigApprovalValue(
			ConfigApprovalValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"account_integration_id": convert.StrToType(&apiConfig.ApprovePolicyTypeConfiguration3.AccountIntegrationId),
				"flow_id":                convert.StrToType(apiConfig.ApprovePolicyTypeConfiguration3.FlowId),
				"workflow_id":            convert.StrToType(apiConfig.ApprovePolicyTypeConfiguration3.WorkflowId),
				"workflow_type":          convert.StrToType(apiConfig.ApprovePolicyTypeConfiguration3.WorkflowType),
			},
		)
		if approvalDiags.HasError() {
			diags.Append(approvalDiags...)
		} else {
			data.ConfigApproval = approvalValue
		}
	}

	// 2. BackupTargetsPolicyTypeConfiguration -> backup_storage
	if apiConfig.BackupTargetsPolicyTypeConfiguration3 != nil {
		// Handle BackupStorageIds as a set of strings
		var backupStorageIDsSet types.Set
		if len(apiConfig.BackupTargetsPolicyTypeConfiguration3.BackupStorageIds) == 0 {
			backupStorageIDsSet = types.SetValueMust(types.StringType, []attr.Value{})
		} else {
			stringValues := make([]attr.Value, len(apiConfig.BackupTargetsPolicyTypeConfiguration3.BackupStorageIds))
			for i, id := range apiConfig.BackupTargetsPolicyTypeConfiguration3.BackupStorageIds {
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

	// 3. BackupCreationPolicyTypeConfiguration3 -> create_backup
	if apiConfig.BackupCreationPolicyTypeConfiguration3 != nil {
		createBackupValue, createBackupDiags := NewConfigCreateBackupValue(
			ConfigCreateBackupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"create_backup":      convert.BoolToType(apiConfig.BackupCreationPolicyTypeConfiguration3.CreateBackup),
				"create_backup_type": convert.StrToType(&apiConfig.BackupCreationPolicyTypeConfiguration3.CreateBackupType),
			},
		)
		if createBackupDiags.HasError() {
			diags.Append(createBackupDiags...)
		} else {
			data.ConfigCreateBackup = createBackupValue
		}
	}

	// 4. UserCreationPolicyTypeConfiguration -> create_user
	if apiConfig.UserCreationPolicyTypeConfiguration3 != nil {
		createUserValue, createUserDiags := NewConfigCreateUserValue(
			ConfigCreateUserValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"create_user":      convert.BoolToType(apiConfig.UserCreationPolicyTypeConfiguration3.CreateUser),
				"create_user_type": convert.StrToType(&apiConfig.UserCreationPolicyTypeConfiguration3.CreateUserType),
			},
		)
		if createUserDiags.HasError() {
			diags.Append(createUserDiags...)
		} else {
			data.ConfigCreateUser = createUserValue
		}
	}

	// 5. UserGroupCreationPolicyTypeConfiguration3 -> create_user_group
	if apiConfig.UserGroupCreationPolicyTypeConfiguration3 != nil {
		createUserGroupValue, createUserGroupDiags := NewConfigCreateUserGroupValue(
			ConfigCreateUserGroupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"user_group": convert.StrToType(&apiConfig.UserGroupCreationPolicyTypeConfiguration3.UserGroup),
			},
		)
		if createUserGroupDiags.HasError() {
			diags.Append(createUserGroupDiags...)
		} else {
			data.ConfigCreateUserGroup = createUserGroupValue
		}
	}

	// 6. CypherAccessPolicyTypeConfiguration3 -> cypher
	if apiConfig.CypherAccessPolicyTypeConfiguration3 != nil {
		cypherValue, cypherDiags := NewConfigCypherValue(
			ConfigCypherValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"delete":      convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration3.Delete),
				"key_pattern": convert.StrToType(&apiConfig.CypherAccessPolicyTypeConfiguration3.KeyPattern),
				"list":        convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration3.List),
				"read":        convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration3.Read),
				"update":      convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration3.Update),
				"write":       convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration3.Write),
			},
		)
		if cypherDiags.HasError() {
			diags.Append(cypherDiags...)
		} else {
			data.ConfigCypher = cypherValue
		}
	}

	// 7. BudgetPolicyTypeConfiguration3 -> max_price
	if apiConfig.BudgetPolicyTypeConfiguration3 != nil {
		maxPriceAttrs := map[string]attr.Value{
			"max_price":          convert.NumToType(&apiConfig.BudgetPolicyTypeConfiguration3.MaxPrice),
			"max_price_currency": convert.StrToType(apiConfig.BudgetPolicyTypeConfiguration3.MaxPriceCurrency),
			"max_price_unit":     convert.StrToType(apiConfig.BudgetPolicyTypeConfiguration3.MaxPriceUnit),
		}

		maxPriceValue, maxPriceDiags := NewConfigMaxPriceValue(ConfigMaxPriceValue{}.AttributeTypes(ctx), maxPriceAttrs)
		if maxPriceDiags.HasError() {
			diags.Append(maxPriceDiags...)
		} else {
			data.ConfigMaxPrice = maxPriceValue
		}
	}

	// 8. MaxMemoryPolicyTypeConfiguration -> max_memory
	if apiConfig.MaxMemoryPolicyTypeConfiguration3 != nil {
		excludeContainers := types.BoolNull()
		if apiConfig.MaxMemoryPolicyTypeConfiguration3.ExcludeContainers != nil {
			excludeContainers = convert.StringToBool(ctx, *apiConfig.MaxMemoryPolicyTypeConfiguration3.ExcludeContainers)
		}

		maxMemoryAttrs := map[string]attr.Value{
			"max_memory":         convert.StrToType(&apiConfig.MaxMemoryPolicyTypeConfiguration3.MaxMemory),
			"exclude_containers": excludeContainers,
		}

		maxMemoryValue, maxMemoryDiags := NewConfigMaxMemoryValue(ConfigMaxMemoryValue{}.AttributeTypes(ctx), maxMemoryAttrs)
		if maxMemoryDiags.HasError() {
			diags.Append(maxMemoryDiags...)
		} else {
			data.ConfigMaxMemory = maxMemoryValue
		}
	}

	// 9. MaxCoresPolicyTypeConfiguration3 -> max_cores
	if apiConfig.MaxCoresPolicyTypeConfiguration3 != nil {
		excludeContainers := types.BoolNull()
		if apiConfig.MaxCoresPolicyTypeConfiguration3.ExcludeContainers != nil {
			excludeContainers = convert.StringToBool(ctx, *apiConfig.MaxCoresPolicyTypeConfiguration3.ExcludeContainers)
		}

		maxCoresAttrs := map[string]attr.Value{
			"max_cores":          convert.StrToType(&apiConfig.MaxCoresPolicyTypeConfiguration3.MaxCores),
			"exclude_containers": excludeContainers,
		}

		maxCoresValue, maxCoresDiags := NewConfigMaxCoresValue(ConfigMaxCoresValue{}.AttributeTypes(ctx), maxCoresAttrs)
		if maxCoresDiags.HasError() {
			diags.Append(maxCoresDiags...)
		} else {
			data.ConfigMaxCores = maxCoresValue
		}
	}

	// 10. DelayedDeletePolicyTypeConfiguration3 -> delayed_removal
	if apiConfig.DelayedDeletePolicyTypeConfiguration3 != nil {
		delayedRemovalAttrs := map[string]attr.Value{
			"removal_age": convert.StrToType(&apiConfig.DelayedDeletePolicyTypeConfiguration3.RemovalAge),
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

	// 11. ExpirationPolicyTypeConfiguration3 -> lifecycle
	if apiConfig.ExpirationPolicyTypeConfiguration3 != nil {
		lifecycleAllowExtend := types.BoolNull()
		if apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleAllowExtend != nil {
			lifecycleAllowExtend = convert.StringToBool(ctx, *apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleAllowExtend)
		}
		lifecycleAutoRenew := types.BoolNull()
		if apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleAutoRenew != nil {
			lifecycleAutoRenew = convert.StringToBool(ctx, *apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleAutoRenew)
		}

		lifecycleAttrs := map[string]attr.Value{
			"account_integration_id": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration3.AccountIntegrationId,
			),
			"flow_id":                convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration3.FlowId),
			"lifecycle_age":          convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleAge),
			"lifecycle_allow_extend": lifecycleAllowExtend,
			"lifecycle_auto_renew":   lifecycleAutoRenew,
			"lifecycle_extensions_before_approval": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleExtensionsBeforeApproval,
			),
			"lifecycle_hide_fixed": convert.BoolToType(
				apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleHideFixed,
			),
			"lifecycle_message": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleMessage,
			),
			"lifecycle_notify": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleNotify,
			),
			"lifecycle_renewal": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleRenewal,
			),
			"lifecycle_type": convert.StrToType(
				&apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleType,
			),
			"lifecycle_workflow_id": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleWorkflowId,
			),
			"workflow_type": convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration3.WorkflowType),
		}

		lifecycleValue, lifecycleDiags := NewConfigLifecycleValue(ConfigLifecycleValue{}.AttributeTypes(ctx), lifecycleAttrs)
		if lifecycleDiags.HasError() {
			diags.Append(lifecycleDiags...)
		} else {
			data.ConfigLifecycle = lifecycleValue
		}
	}

	// 12. HostnamePolicyTypeConfiguration3 -> host_naming
	if apiConfig.HostnamePolicyTypeConfiguration3 != nil {
		hostNamingAttrs := map[string]attr.Value{
			"host_naming_pattern": convert.StrToType(apiConfig.HostnamePolicyTypeConfiguration3.HostNamingPattern),
			"host_naming_type":    convert.StrToType(&apiConfig.HostnamePolicyTypeConfiguration3.HostNamingType),
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

	// 13. InstanceNamePolicyTypeConfiguration3 -> naming
	if apiConfig.InstanceNamePolicyTypeConfiguration3 != nil {
		namingAttrs := map[string]attr.Value{
			"naming_conflict": convert.BoolToType(apiConfig.InstanceNamePolicyTypeConfiguration3.NamingConflict),
			"naming_pattern":  convert.StrToType(apiConfig.InstanceNamePolicyTypeConfiguration3.NamingPattern),
			"naming_type":     convert.StrToType(&apiConfig.InstanceNamePolicyTypeConfiguration3.NamingType),
		}

		namingValue, namingDiags := NewConfigNamingValue(ConfigNamingValue{}.AttributeTypes(ctx), namingAttrs)
		if namingDiags.HasError() {
			diags.Append(namingDiags...)
		} else {
			data.ConfigNaming = namingValue
		}
	}

	// 14. MaxContainersPolicyTypeConfiguration3 -> max_containers
	if apiConfig.MaxContainersPolicyTypeConfiguration3 != nil {
		maxContainersAttrs := map[string]attr.Value{
			"max_containers": convert.StrToType(&apiConfig.MaxContainersPolicyTypeConfiguration3.MaxContainers),
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

	// 15. MaxHostsPolicyTypeConfiguration3 -> max_hosts
	if apiConfig.MaxHostsPolicyTypeConfiguration3 != nil {
		maxHostsAttrs := map[string]attr.Value{
			"max_hosts": convert.StrToType(&apiConfig.MaxHostsPolicyTypeConfiguration3.MaxHosts),
		}

		maxHostsValue, maxHostsDiags := NewConfigMaxHostsValue(ConfigMaxHostsValue{}.AttributeTypes(ctx), maxHostsAttrs)
		if maxHostsDiags.HasError() {
			diags.Append(maxHostsDiags...)
		} else {
			data.ConfigMaxHosts = maxHostsValue
		}
	}

	// 16. NetworkQuotaPolicyTypeConfiguration3 -> max_networks
	if apiConfig.NetworkQuotaPolicyTypeConfiguration3 != nil {
		maxNetworksAttrs := map[string]attr.Value{
			"max_networks": convert.StrToType(&apiConfig.NetworkQuotaPolicyTypeConfiguration3.MaxNetworks),
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

	// 17. MaxPoolMembersPolicyTypeConfiguration3 -> max_pool_members
	if apiConfig.MaxPoolMembersPolicyTypeConfiguration3 != nil {
		maxPoolMembersAttrs := map[string]attr.Value{
			"max_pool_members": convert.StrToType(&apiConfig.MaxPoolMembersPolicyTypeConfiguration3.MaxPoolMembers),
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

	// 18. MaxLoadBalancerPoolsPolicyTypeConfiguration3 -> max_pools
	if apiConfig.MaxLoadBalancerPoolsPolicyTypeConfiguration3 != nil {
		maxPoolsAttrs := map[string]attr.Value{
			"max_pools": convert.StrToType(&apiConfig.MaxLoadBalancerPoolsPolicyTypeConfiguration3.MaxPools),
		}

		maxPoolsValue, maxPoolsDiags := NewConfigMaxPoolsValue(ConfigMaxPoolsValue{}.AttributeTypes(ctx), maxPoolsAttrs)
		if maxPoolsDiags.HasError() {
			diags.Append(maxPoolsDiags...)
		} else {
			data.ConfigMaxPools = maxPoolsValue
		}
	}

	// 19. RouterQuotaPolicyTypeConfiguration3 -> max_routers
	if apiConfig.RouterQuotaPolicyTypeConfiguration3 != nil {
		maxRoutersAttrs := map[string]attr.Value{
			"max_routers": convert.StrToType(&apiConfig.RouterQuotaPolicyTypeConfiguration3.MaxRouters),
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

	// 20. MaxSnapshotsPolicyTypeConfiguration3 -> max_snapshots
	if apiConfig.MaxSnapshotsPolicyTypeConfiguration3 != nil {
		maxSnapshotsAttrs := map[string]attr.Value{
			"max_snapshots": convert.StrToType(&apiConfig.MaxSnapshotsPolicyTypeConfiguration3.MaxSnapshots),
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

	// 21. MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration3 -> max_storage
	if apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration3 != nil {
		excludeContainers := types.BoolNull()
		if apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration3.ExcludeContainers != nil {
			excludeContainers = convert.StringToBool(ctx, *apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration3.ExcludeContainers)
		}

		maxStorageAttrs := map[string]attr.Value{
			"exclude_containers": excludeContainers,
			"max_storage": convert.StrToType(
				&apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration3.MaxStorage,
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

	// 22. MaxVirtualServersPolicyTypeConfiguration3 -> max_virtual_servers
	if apiConfig.MaxVirtualServersPolicyTypeConfiguration3 != nil {
		maxVirtualServersAttrs := map[string]attr.Value{
			"max_virtual_servers": convert.StrToType(&apiConfig.MaxVirtualServersPolicyTypeConfiguration3.MaxVirtualServers),
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

	// 23. MaxVMsPolicyTypeConfiguration3 -> max_vms
	if apiConfig.MaxVMsPolicyTypeConfiguration3 != nil {
		maxVmsAttrs := map[string]attr.Value{
			"max_vms": convert.StrToType(&apiConfig.MaxVMsPolicyTypeConfiguration3.MaxVms),
		}

		maxVmsValue, maxVmsDiags := NewConfigMaxVmsValue(ConfigMaxVmsValue{}.AttributeTypes(ctx), maxVmsAttrs)
		if maxVmsDiags.HasError() {
			diags.Append(maxVmsDiags...)
		} else {
			data.ConfigMaxVms = maxVmsValue
		}
	}

	// 24. MessageOfTheDayPolicyTypeConfiguration32 -> motd
	if apiConfig.MessageOfTheDayPolicyTypeConfiguration3 != nil {
		motdAttrs := map[string]attr.Value{
			"motd_fullpage": types.StringNull(),
			"motddate":      convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdDate),
			"motdmessage":   convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdMessage),
			"motdtitle":     types.StringNull(), // NullableString type
			"motdtype":      convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdType),
		}

		// Handle NullableString for MotdFullPage - check if field exists and is not nil
		if apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdFullPage != nil {
			motdAttrs["motd_fullpage"] = convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdFullPage.String)
		}
		// Handle NullableString for MotdTitle
		if apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdTitle.IsSet() {
			motdAttrs["motdtitle"] = types.StringValue(*apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdTitle.Get())
		}

		motdValue, motdDiags := NewConfigMotdValue(ConfigMotdValue{}.AttributeTypes(ctx), motdAttrs)
		if motdDiags.HasError() {
			diags.Append(motdDiags...)
		} else {
			data.ConfigMotd = motdValue
		}
	}

	// 25. PowerSchedulePolicyTypeConfiguration3 -> power_schedule
	if apiConfig.PowerSchedulePolicyTypeConfiguration3 != nil {
		powerScheduleAttrs := map[string]attr.Value{
			"power_schedule": convert.StrToType(
				apiConfig.PowerSchedulePolicyTypeConfiguration3.PowerSchedule,
			),
			"power_schedule_hide_fixed": convert.BoolToType(
				apiConfig.PowerSchedulePolicyTypeConfiguration3.PowerScheduleHideFixed,
			),
			"power_schedule_type": convert.StrToType(
				&apiConfig.PowerSchedulePolicyTypeConfiguration3.PowerScheduleType,
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

	// 26. RequiredNetworkPolicyTypeConfiguration3 -> required_network
	if apiConfig.RequiredNetworkPolicyTypeConfiguration3 != nil {
		// Handle RequiredNetworks as a set of integers
		var requiredNetworksSet types.Set
		if len(apiConfig.RequiredNetworkPolicyTypeConfiguration3.RequiredNetworks) == 0 {
			requiredNetworksSet = types.SetValueMust(types.Int64Type, []attr.Value{})
		} else {
			int64Values := make([]attr.Value, len(apiConfig.RequiredNetworkPolicyTypeConfiguration3.RequiredNetworks))
			for i, networkId := range apiConfig.RequiredNetworkPolicyTypeConfiguration3.RequiredNetworks {
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

	// 27. ClusterResourceNamePolicyTypeConfiguration3 -> server_naming
	if apiConfig.ClusterResourceNamePolicyTypeConfiguration3 != nil {
		serverNamingAttrs := map[string]attr.Value{
			"server_naming_conflict": convert.BoolToType(
				apiConfig.ClusterResourceNamePolicyTypeConfiguration3.ServerNamingConflict,
			),
			"server_naming_pattern": convert.StrToType(
				apiConfig.ClusterResourceNamePolicyTypeConfiguration3.ServerNamingPattern,
			),
			"server_naming_type": convert.StrToType(
				&apiConfig.ClusterResourceNamePolicyTypeConfiguration3.ServerNamingType,
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

	// 28. ShutdownPolicyTypeConfiguration3 -> shutdown
	if apiConfig.ShutdownPolicyTypeConfiguration3 != nil {
		shutdownAllowExtend := types.BoolNull()
		if apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownAllowExtend != nil {
			shutdownAllowExtend = convert.StringToBool(ctx, *apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownAllowExtend)
		}
		shutdownAutoRenew := types.BoolNull()
		if apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownAutoRenew != nil {
			shutdownAutoRenew = convert.StringToBool(ctx, *apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownAutoRenew)
		}

		shutdownAttrs := map[string]attr.Value{
			"account_integration_id": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration3.AccountIntegrationId,
			),
			"flow_id":               convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration3.FlowId),
			"shutdown_age":          convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownAge),
			"shutdown_allow_extend": shutdownAllowExtend,
			"shutdown_auto_renew":   shutdownAutoRenew,
			"shutdown_extensions_before_approval": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownExtensionsBeforeApproval,
			),
			"shutdown_hide_fixed": convert.BoolToType(
				apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownHideFixed,
			),
			"shutdown_message": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownMessage,
			),
			"shutdown_notify": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownNotify,
			),
			"shutdown_renewal": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownRenewal,
			),
			"shutdown_type": convert.StrToType(
				&apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownType,
			),
			"shutdown_workflow_id": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownWorkflowId,
			),
			"workflow_type": convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration3.WorkflowType),
		}

		shutdownValue, shutdownDiags := NewConfigShutdownValue(ConfigShutdownValue{}.AttributeTypes(ctx), shutdownAttrs)
		if shutdownDiags.HasError() {
			diags.Append(shutdownDiags...)
		} else {
			data.ConfigShutdown = shutdownValue
		}
	}

	// 29. StorageServerStorageQuotaPolicyTypeConfiguration3 -> storage_server_quota
	if apiConfig.StorageServerStorageQuotaPolicyTypeConfiguration3 != nil {
		storageServerQuotaAttrs := map[string]attr.Value{
			"max_storage":       convert.StrToType(apiConfig.StorageServerStorageQuotaPolicyTypeConfiguration3.MaxStorage),
			"storage_server_id": convert.StrToType(&apiConfig.StorageServerStorageQuotaPolicyTypeConfiguration3.StorageServerId),
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

	// 30. TagsPolicyTypeConfiguration3 -> tags
	if apiConfig.TagsPolicyTypeConfiguration3 != nil {
		tagsAttrs := map[string]attr.Value{
			"key":           convert.StrToType(apiConfig.TagsPolicyTypeConfiguration3.Key),
			"strict":        types.BoolValue(apiConfig.TagsPolicyTypeConfiguration3.Strict),
			"value":         convert.StrToType(apiConfig.TagsPolicyTypeConfiguration3.Value),
			"value_list_id": convert.StrToType(apiConfig.TagsPolicyTypeConfiguration3.ValueListId),
		}

		tagsValue, tagsDiags := NewConfigTagsValue(ConfigTagsValue{}.AttributeTypes(ctx), tagsAttrs)
		if tagsDiags.HasError() {
			diags.Append(tagsDiags...)
		} else {
			data.ConfigTags = tagsValue
		}
	}

	// 31. WorkflowPolicyTypeConfiguration3 -> workflow
	if apiConfig.WorkflowPolicyTypeConfiguration3 != nil {
		workflowAttrs := map[string]attr.Value{
			"workflow_id": convert.StrToType(&apiConfig.WorkflowPolicyTypeConfiguration3.WorkflowId),
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

	policies := ps.Policies

	// Additional filtering to ensure exact name match (API might return partial matches)
	var filteredPolicies []sdk.ListPolicies200ResponseAllOfPoliciesInner
	for _, p := range policies {
		if p.Name != nil && *p.Name == data.Name.ValueString() {
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
	if policy.Id == nil {
		diags.AddError(summary, consts.ErrorNoPolicyFound)

		return diags
	}

	return getPolicyByID(ctx, *policy.Id, data, apiClient)
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
	policy := p.Policy
	if policy == nil {
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
		data.AssociatedResourceId = types.Int64Value(*policy.RefId.Get())
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
	if policy.Owner.IsSet() && policy.Owner.Get() != nil && policy.Owner.Get().Id != nil {
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
