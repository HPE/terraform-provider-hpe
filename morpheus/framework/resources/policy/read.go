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

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// mapPolicyConfigToState maps the API config structure to the resource schema structure
func mapPolicyConfigToState(
	ctx context.Context,
	state *PolicyModel,
	apiConfig *sdk.GetPolicies200ResponseAllOfPolicyConfig,
	policyTypeCode string,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	// Map API config field to the corresponding schema field based on policy type code
	switch policyTypeCode {
	case "deleteApproval", "provisionApproval", "reconfigureApproval", "workflowApproval":
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
			state.ConfigApproval = approvalValue
		}

	case "backupStorage":
		// Handle BackupStorageIds as a set of int64
		var backupStorageIDsSet types.Set
		var setDiags diag.Diagnostics
		if len(apiConfig.BackupTargetsPolicyTypeConfiguration3.BackupStorageIds) == 0 {
			backupStorageIDsSet = types.SetValueMust(types.Int64Type, []attr.Value{})
		} else {
			// BackupStorageIds come from API as []string, convert to []int64
			backupStorageIDsSet, setDiags = types.SetValueFrom(
				ctx,
				types.Int64Type,
				apiConfig.BackupTargetsPolicyTypeConfiguration3.BackupStorageIds,
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

	case "createBackup":
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
			state.ConfigCreateBackup = createBackupValue
		}

	case "createUser":
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
			state.ConfigCreateUser = createUserValue
		}

	case "createUserGroup":
		createUserGroupValue, createUserGroupDiags := NewConfigCreateUserGroupValue(
			ConfigCreateUserGroupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"user_group": convert.StrToType(&apiConfig.UserGroupCreationPolicyTypeConfiguration3.UserGroup),
			},
		)
		if createUserGroupDiags.HasError() {
			diags.Append(createUserGroupDiags...)
		} else {
			state.ConfigCreateUserGroup = createUserGroupValue
		}

	case "cypher":
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
			state.ConfigCypher = cypherValue
		}

	case "maxPrice":
		maxPriceAttrs := map[string]attr.Value{
			"max_price":          convert.NumToType(&apiConfig.BudgetPolicyTypeConfiguration3.MaxPrice),
			"max_price_currency": convert.StrToType(apiConfig.BudgetPolicyTypeConfiguration3.MaxPriceCurrency),
			"max_price_unit":     convert.StrToType(apiConfig.BudgetPolicyTypeConfiguration3.MaxPriceUnit),
		}

		maxPriceValue, maxPriceDiags := NewConfigMaxPriceValue(ConfigMaxPriceValue{}.AttributeTypes(ctx), maxPriceAttrs)
		if maxPriceDiags.HasError() {
			diags.Append(maxPriceDiags...)
		} else {
			state.ConfigMaxPrice = maxPriceValue
		}

	case "maxMemory":
		maxMemoryAttrs := map[string]attr.Value{
			"max_memory":         convert.StrToType(&apiConfig.MaxMemoryPolicyTypeConfiguration3.MaxMemory),
			"exclude_containers": convert.StringToBool(ctx, apiConfig.MaxMemoryPolicyTypeConfiguration3.GetExcludeContainers()),
		}

		maxMemoryValue, maxMemoryDiags := NewConfigMaxMemoryValue(ConfigMaxMemoryValue{}.AttributeTypes(ctx), maxMemoryAttrs)
		if maxMemoryDiags.HasError() {
			diags.Append(maxMemoryDiags...)
		} else {
			state.ConfigMaxMemory = maxMemoryValue
		}

	case "maxCores":
		maxCoresAttrs := map[string]attr.Value{
			"max_cores":          convert.StrToType(&apiConfig.MaxCoresPolicyTypeConfiguration3.MaxCores),
			"exclude_containers": convert.StringToBool(ctx, apiConfig.MaxCoresPolicyTypeConfiguration3.GetExcludeContainers()),
		}

		maxCoresValue, maxCoresDiags := NewConfigMaxCoresValue(ConfigMaxCoresValue{}.AttributeTypes(ctx), maxCoresAttrs)
		if maxCoresDiags.HasError() {
			diags.Append(maxCoresDiags...)
		} else {
			state.ConfigMaxCores = maxCoresValue
		}

	case "delayedRemoval":
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
			state.ConfigDelayedRemoval = delayedRemovalValue
		}

	case "lifecycle":
		lifecycleAttrs := map[string]attr.Value{
			"account_integration_id": convert.StrToType(
				apiConfig.ExpirationPolicyTypeConfiguration3.AccountIntegrationId,
			),
			"flow_id":       convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration3.FlowId),
			"lifecycle_age": convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration3.LifecycleAge),
			"lifecycle_allow_extend": convert.StringToBool(
				ctx,
				apiConfig.ExpirationPolicyTypeConfiguration3.GetLifecycleAllowExtend(),
			),
			"lifecycle_auto_renew": convert.StringToBool(
				ctx,
				apiConfig.ExpirationPolicyTypeConfiguration3.GetLifecycleAutoRenew(),
			),
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
			state.ConfigLifecycle = lifecycleValue
		}

	case "hostNaming":
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
			state.ConfigHostNaming = hostNamingValue
		}

	case "naming":
		namingAttrs := map[string]attr.Value{
			"naming_conflict": convert.BoolToType(apiConfig.InstanceNamePolicyTypeConfiguration3.NamingConflict),
			"naming_pattern":  convert.StrToType(apiConfig.InstanceNamePolicyTypeConfiguration3.NamingPattern),
			"naming_type":     convert.StrToType(&apiConfig.InstanceNamePolicyTypeConfiguration3.NamingType),
		}

		namingValue, namingDiags := NewConfigNamingValue(ConfigNamingValue{}.AttributeTypes(ctx), namingAttrs)
		if namingDiags.HasError() {
			diags.Append(namingDiags...)
		} else {
			state.ConfigNaming = namingValue
		}

	case "maxContainers":
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
			state.ConfigMaxContainers = maxContainersValue
		}

	case "maxHosts":
		maxHostsAttrs := map[string]attr.Value{
			"max_hosts": convert.StrToType(&apiConfig.MaxHostsPolicyTypeConfiguration3.MaxHosts),
		}

		maxHostsValue, maxHostsDiags := NewConfigMaxHostsValue(ConfigMaxHostsValue{}.AttributeTypes(ctx), maxHostsAttrs)
		if maxHostsDiags.HasError() {
			diags.Append(maxHostsDiags...)
		} else {
			state.ConfigMaxHosts = maxHostsValue
		}

	case "maxNetworks":
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
			state.ConfigMaxNetworks = maxNetworksValue
		}

	case "maxPoolMembers":
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
			state.ConfigMaxPoolMembers = maxPoolMembersValue
		}

	case "maxPools":
		maxPoolsAttrs := map[string]attr.Value{
			"max_pools": convert.StrToType(&apiConfig.MaxLoadBalancerPoolsPolicyTypeConfiguration3.MaxPools),
		}

		maxPoolsValue, maxPoolsDiags := NewConfigMaxPoolsValue(ConfigMaxPoolsValue{}.AttributeTypes(ctx), maxPoolsAttrs)
		if maxPoolsDiags.HasError() {
			diags.Append(maxPoolsDiags...)
		} else {
			state.ConfigMaxPools = maxPoolsValue
		}

	case "maxRouters":
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
			state.ConfigMaxRouters = maxRoutersValue
		}

	case "maxSnapshots":
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
			state.ConfigMaxSnapshots = maxSnapshotsValue
		}

	case "maxStorage", "storageBucketQuota", "storageShareQuota":
		maxStorageAttrs := map[string]attr.Value{
			"exclude_containers": convert.StringToBool(
				ctx,
				apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration3.GetExcludeContainers(),
			),
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
			state.ConfigMaxStorage = maxStorageValue
		}

	case "maxVirtualServers":
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
			state.ConfigMaxVirtualServers = maxVirtualServersValue
		}

	case "maxVms":
		maxVmsAttrs := map[string]attr.Value{
			"max_vms": convert.StrToType(&apiConfig.MaxVMsPolicyTypeConfiguration3.MaxVms),
		}

		maxVmsValue, maxVmsDiags := NewConfigMaxVmsValue(ConfigMaxVmsValue{}.AttributeTypes(ctx), maxVmsAttrs)
		if maxVmsDiags.HasError() {
			diags.Append(maxVmsDiags...)
		} else {
			state.ConfigMaxVms = maxVmsValue
		}

	case "motd":
		motdAttrs := map[string]attr.Value{
			"motddate":    convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdDate),
			"motdmessage": convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdMessage),
			"motdtitle":   types.StringNull(),
			"motdtype":    convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdType),
		}

		// Handle NullableString for MotdTitle
		if apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdTitle.IsSet() {
			motdAttrs["motdtitle"] = types.StringValue(*apiConfig.MessageOfTheDayPolicyTypeConfiguration3.MotdTitle.Get())
		}

		motdValue, motdDiags := NewConfigMotdValue(ConfigMotdValue{}.AttributeTypes(ctx), motdAttrs)
		if motdDiags.HasError() {
			diags.Append(motdDiags...)
		} else {
			state.ConfigMotd = motdValue
		}

	case "powerSchedule":
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
			state.ConfigPowerSchedule = powerScheduleValue
		}

	case "requiredNetwork":
		// Handle RequiredNetworks as a set of integers
		var requiredNetworksSet types.Set
		if len(apiConfig.RequiredNetworkPolicyTypeConfiguration3.RequiredNetworks) == 0 {
			requiredNetworksSet = types.SetValueMust(types.Int64Type, []attr.Value{})
		} else {
			int64Values := make([]attr.Value, len(apiConfig.RequiredNetworkPolicyTypeConfiguration3.RequiredNetworks))
			for i, networkID := range apiConfig.RequiredNetworkPolicyTypeConfiguration3.RequiredNetworks {
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

	case "serverNaming":
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
			state.ConfigServerNaming = serverNamingValue
		}

	case "shutdown":
		shutdownAttrs := map[string]attr.Value{
			"account_integration_id": convert.StrToType(
				apiConfig.ShutdownPolicyTypeConfiguration3.AccountIntegrationId,
			),
			"flow_id":      convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration3.FlowId),
			"shutdown_age": convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration3.ShutdownAge),
			"shutdown_allow_extend": convert.StringToBool(
				ctx,
				apiConfig.ShutdownPolicyTypeConfiguration3.GetShutdownAllowExtend(),
			),
			"shutdown_auto_renew": convert.StringToBool(
				ctx,
				apiConfig.ShutdownPolicyTypeConfiguration3.GetShutdownAutoRenew(),
			),
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
			state.ConfigShutdown = shutdownValue
		}

	case "storageServerQuota":
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
			state.ConfigStorageServerQuota = storageServerQuotaValue
		}

	case "tags":
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
			state.ConfigTags = tagsValue
		}

	case "workflow":
		workflowAttrs := map[string]attr.Value{
			"workflow_id": convert.StrToType(&apiConfig.WorkflowPolicyTypeConfiguration3.WorkflowId),
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
			fmt.Sprintf("policy %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
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
			PolicyTypeValue{}.AttributeTypes(ctx), policyTypeAttrs,
		)
		if policyTypeDiags.HasError() {
			diags.Append(policyTypeDiags...)

			return state, diags
		}
		state.PolicyType = policyTypeValue
	}

	// Handle Config - use static schema fields when available, fallback to dynamic
	if p.Config != nil {
		// Check if user is using dynamic config field
		usingDynamicConfig := plan != nil && !plan.Config.IsNull() && !plan.Config.IsUnknown()

		if usingDynamicConfig {
			// User is using dynamic config - only populate the dynamic field
			state.Config = plan.Config
		} else {
			// User is using static config fields - populate them from API response
			// Get policy type code for mapping (required field, but API could return nil)
			policyTypeCode := ""
			if p.PolicyType != nil && p.PolicyType.Code != nil {
				policyTypeCode = *p.PolicyType.Code
			}

			// Map API config to static schema fields
			configDiags := mapPolicyConfigToState(ctx, &state, p.Config, policyTypeCode)
			if configDiags.HasError() {
				diags.Append(configDiags...)

				return state, diags
			}

			// Don't populate dynamic config when using static fields to avoid drift
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
	owner := p.GetOwner()

	if owner.Id != nil {
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
