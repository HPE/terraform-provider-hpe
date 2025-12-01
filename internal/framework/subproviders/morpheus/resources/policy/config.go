// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
)

// mapStateToAddPolicyConfig converts state config fields to API config structure for Add (Create)
func mapStateToAddPolicyConfig(
	ctx context.Context,
	plan *PolicyModel,
) (*sdk.AddPoliciesRequestPolicyConfig, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var configMap map[string]interface{}

	// First, check if the legacy dynamic config field is set
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configValue := plan.Config.UnderlyingValue()
		anyValue, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			diags.AddError(
				"map state to API config",
				"failed to convert dynamic config: "+err.Error(),
			)
			return nil, diags
		}
		var ok bool
		configMap, ok = anyValue.(map[string]interface{})
		if !ok {
			diags.AddError(
				"map state to API config",
				"config must be a map",
			)
			return nil, diags
		}

		// Check if config is empty
		if len(configMap) == 0 {
			diags.AddError(
				"map state to API config",
				"config cannot be empty. Please provide the required configuration fields for this policy type.",
			)
			return nil, diags
		}
	} else {
		// Use static config fields
		configMap = make(map[string]interface{})

		// Map each static config field to the API structure
		// 1. config_approval -> ApprovePolicyTypeConfiguration
		if !plan.ConfigApproval.IsNull() && !plan.ConfigApproval.IsUnknown() {
			approval := make(map[string]interface{})
			if !plan.ConfigApproval.AccountIntegrationId.IsNull() {
				approval["accountIntegrationId"] = plan.ConfigApproval.AccountIntegrationId.ValueString()
			}
			if !plan.ConfigApproval.FlowId.IsNull() {
				approval["flowId"] = plan.ConfigApproval.FlowId.ValueString()
			}
			if !plan.ConfigApproval.WorkflowId.IsNull() {
				approval["workflowId"] = plan.ConfigApproval.WorkflowId.ValueString()
			}
			if !plan.ConfigApproval.WorkflowType.IsNull() {
				approval["workflowType"] = plan.ConfigApproval.WorkflowType.ValueString()
			}
			if len(approval) > 0 {
				configMap = approval
			}
		}

		// 2. config_backup_storage -> BackupTargetsPolicyTypeConfiguration
		if !plan.ConfigBackupStorage.IsNull() {
			backupStorage := make(map[string]interface{})
			if !plan.ConfigBackupStorage.BackupStorageIds.IsNull() {
				var backupStorageIDs []int64
				diagsSet := plan.ConfigBackupStorage.BackupStorageIds.ElementsAs(ctx, &backupStorageIDs, false)
				if diagsSet.HasError() {
					diags.Append(diagsSet...)
					return nil, diags
				}
				backupStorage["backupStorageIds"] = backupStorageIDs
			}
			if len(backupStorage) > 0 {
				configMap = backupStorage
			}
		}

		// 3. config_create_backup -> BackupCreationPolicyTypeConfiguration
		if !plan.ConfigCreateBackup.IsNull() {
			createBackup := make(map[string]interface{})
			if !plan.ConfigCreateBackup.CreateBackup.IsNull() {
				createBackup["createBackup"] = plan.ConfigCreateBackup.CreateBackup.ValueBool()
			}
			if !plan.ConfigCreateBackup.CreateBackupType.IsNull() {
				createBackup["createBackupType"] = plan.ConfigCreateBackup.CreateBackupType.ValueString()
			}
			if len(createBackup) > 0 {
				configMap = createBackup
			}
		}

		// 4. config_create_user -> UserCreationPolicyTypeConfiguration
		if !plan.ConfigCreateUser.IsNull() {
			createUser := make(map[string]interface{})
			if !plan.ConfigCreateUser.CreateUser.IsNull() {
				createUser["createUser"] = plan.ConfigCreateUser.CreateUser.ValueBool()
			}
			if !plan.ConfigCreateUser.CreateUserType.IsNull() {
				createUser["createUserType"] = plan.ConfigCreateUser.CreateUserType.ValueString()
			}
			if len(createUser) > 0 {
				configMap = createUser
			}
		}

		// 5. config_create_user_group -> UserGroupCreationPolicyTypeConfiguration
		if !plan.ConfigCreateUserGroup.IsNull() {
			createUserGroup := make(map[string]interface{})
			if !plan.ConfigCreateUserGroup.UserGroup.IsNull() {
				createUserGroup["userGroup"] = plan.ConfigCreateUserGroup.UserGroup.ValueString()
			}
			if len(createUserGroup) > 0 {
				configMap = createUserGroup
			}
		}

		// 6. config_cypher -> CypherAccessPolicyTypeConfiguration
		if !plan.ConfigCypher.IsNull() {
			cypher := make(map[string]interface{})
			if !plan.ConfigCypher.Delete.IsNull() {
				cypher["delete"] = plan.ConfigCypher.Delete.ValueBool()
			}
			if !plan.ConfigCypher.KeyPattern.IsNull() {
				cypher["keyPattern"] = plan.ConfigCypher.KeyPattern.ValueString()
			}
			if !plan.ConfigCypher.List.IsNull() {
				cypher["list"] = plan.ConfigCypher.List.ValueBool()
			}
			if !plan.ConfigCypher.Read.IsNull() {
				cypher["read"] = plan.ConfigCypher.Read.ValueBool()
			}
			if !plan.ConfigCypher.Update.IsNull() {
				cypher["update"] = plan.ConfigCypher.Update.ValueBool()
			}
			if !plan.ConfigCypher.Write.IsNull() {
				cypher["write"] = plan.ConfigCypher.Write.ValueBool()
			}
			if len(cypher) > 0 {
				configMap = cypher
			}
		}

		// 7. config_delayed_removal -> DelayedDeletePolicyTypeConfiguration
		if !plan.ConfigDelayedRemoval.IsNull() {
			delayedRemoval := make(map[string]interface{})
			if !plan.ConfigDelayedRemoval.RemovalAge.IsNull() {
				delayedRemoval["removalAge"] = plan.ConfigDelayedRemoval.RemovalAge.ValueString()
			}
			if len(delayedRemoval) > 0 {
				configMap = delayedRemoval
			}
		}

		// 8. config_host_naming -> HostnamePolicyTypeConfiguration
		if !plan.ConfigHostNaming.IsNull() {
			hostNaming := make(map[string]interface{})
			if !plan.ConfigHostNaming.HostNamingPattern.IsNull() {
				hostNaming["hostNamingPattern"] = plan.ConfigHostNaming.HostNamingPattern.ValueString()
			}
			if !plan.ConfigHostNaming.HostNamingType.IsNull() {
				hostNaming["hostNamingType"] = plan.ConfigHostNaming.HostNamingType.ValueString()
			}
			if len(hostNaming) > 0 {
				configMap = hostNaming
			}
		}

		// 9. config_lifecycle -> ExpirationPolicyTypeConfiguration2
		if !plan.ConfigLifecycle.IsNull() {
			lifecycle := make(map[string]interface{})
			if !plan.ConfigLifecycle.AccountIntegrationId.IsNull() {
				lifecycle["accountIntegrationId"] = plan.ConfigLifecycle.AccountIntegrationId.ValueString()
			}
			if !plan.ConfigLifecycle.FlowId.IsNull() {
				lifecycle["flowId"] = plan.ConfigLifecycle.FlowId.ValueString()
			}
			if !plan.ConfigLifecycle.LifecycleAge.IsNull() {
				lifecycle["lifecycleAge"] = plan.ConfigLifecycle.LifecycleAge.ValueString()
			}
			if !plan.ConfigLifecycle.LifecycleAllowExtend.IsNull() {
				lifecycle["lifecycleAllowExtend"] = plan.ConfigLifecycle.LifecycleAllowExtend.ValueBool()
			}
			if !plan.ConfigLifecycle.LifecycleAutoRenew.IsNull() {
				lifecycle["lifecycleAutoRenew"] = plan.ConfigLifecycle.LifecycleAutoRenew.ValueBool()
			}
			if !plan.ConfigLifecycle.LifecycleExtensionsBeforeApproval.IsNull() {
				lifecycle["lifecycleExtensionsBeforeApproval"] = plan.ConfigLifecycle.LifecycleExtensionsBeforeApproval.ValueString()
			}
			if !plan.ConfigLifecycle.LifecycleHideFixed.IsNull() {
				lifecycle["lifecycleHideFixed"] = plan.ConfigLifecycle.LifecycleHideFixed.ValueBool()
			}
			if !plan.ConfigLifecycle.LifecycleMessage.IsNull() {
				lifecycle["lifecycleMessage"] = plan.ConfigLifecycle.LifecycleMessage.ValueString()
			}
			if !plan.ConfigLifecycle.LifecycleNotify.IsNull() {
				lifecycle["lifecycleNotify"] = plan.ConfigLifecycle.LifecycleNotify.ValueString()
			}
			if !plan.ConfigLifecycle.LifecycleRenewal.IsNull() {
				lifecycle["lifecycleRenewal"] = plan.ConfigLifecycle.LifecycleRenewal.ValueString()
			}
			if !plan.ConfigLifecycle.LifecycleType.IsNull() {
				lifecycle["lifecycleType"] = plan.ConfigLifecycle.LifecycleType.ValueString()
			}
			if !plan.ConfigLifecycle.LifecycleWorkflowId.IsNull() {
				lifecycle["lifecycleWorkflowId"] = plan.ConfigLifecycle.LifecycleWorkflowId.ValueString()
			}
			if !plan.ConfigLifecycle.WorkflowType.IsNull() {
				lifecycle["workflowType"] = plan.ConfigLifecycle.WorkflowType.ValueString()
			}
			if len(lifecycle) > 0 {
				configMap = lifecycle
			}
		}

		// 10. config_max_containers -> MaxContainersPolicyTypeConfiguration
		if !plan.ConfigMaxContainers.IsNull() {
			maxContainers := make(map[string]interface{})
			if !plan.ConfigMaxContainers.MaxContainers.IsNull() {
				maxContainers["maxContainers"] = plan.ConfigMaxContainers.MaxContainers.ValueString()
			}
			if len(maxContainers) > 0 {
				configMap = maxContainers
			}
		}

		// 11. config_max_cores -> MaxCoresPolicyTypeConfiguration
		if !plan.ConfigMaxCores.IsNull() {
			maxCores := make(map[string]interface{})
			if !plan.ConfigMaxCores.MaxCores.IsNull() {
				maxCores["maxCores"] = plan.ConfigMaxCores.MaxCores.ValueString()
			}
			if !plan.ConfigMaxCores.ExcludeContainers.IsNull() {
				maxCores["excludeContainers"] = fmt.Sprintf("%t", plan.ConfigMaxCores.ExcludeContainers.ValueBool())
			}
			if len(maxCores) > 0 {
				configMap = maxCores
			}
		}

		// 12. config_max_hosts -> MaxHostsPolicyTypeConfiguration
		if !plan.ConfigMaxHosts.IsNull() {
			maxHosts := make(map[string]interface{})
			if !plan.ConfigMaxHosts.MaxHosts.IsNull() {
				maxHosts["maxHosts"] = plan.ConfigMaxHosts.MaxHosts.ValueString()
			}
			if len(maxHosts) > 0 {
				configMap = maxHosts
			}
		}

		// 13. config_max_memory -> MaxMemoryPolicyTypeConfiguration
		if !plan.ConfigMaxMemory.IsNull() {
			maxMemory := make(map[string]interface{})
			if !plan.ConfigMaxMemory.MaxMemory.IsNull() {
				maxMemory["maxMemory"] = plan.ConfigMaxMemory.MaxMemory.ValueString()
			}
			if !plan.ConfigMaxMemory.ExcludeContainers.IsNull() {
				maxMemory["excludeContainers"] = fmt.Sprintf("%t", plan.ConfigMaxMemory.ExcludeContainers.ValueBool())
			}
			if len(maxMemory) > 0 {
				configMap = maxMemory
			}
		}

		// 14. config_max_networks -> NetworkQuotaPolicyTypeConfiguration
		if !plan.ConfigMaxNetworks.IsNull() {
			maxNetworks := make(map[string]interface{})
			if !plan.ConfigMaxNetworks.MaxNetworks.IsNull() {
				maxNetworks["maxNetworks"] = plan.ConfigMaxNetworks.MaxNetworks.ValueString()
			}
			if len(maxNetworks) > 0 {
				configMap = maxNetworks
			}
		}

		// 15. config_max_pool_members -> MaxPoolMembersPolicyTypeConfiguration
		if !plan.ConfigMaxPoolMembers.IsNull() {
			maxPoolMembers := make(map[string]interface{})
			if !plan.ConfigMaxPoolMembers.MaxPoolMembers.IsNull() {
				maxPoolMembers["maxPoolMembers"] = plan.ConfigMaxPoolMembers.MaxPoolMembers.ValueString()
			}
			if len(maxPoolMembers) > 0 {
				configMap = maxPoolMembers
			}
		}

		// 16. config_max_pools -> MaxLoadBalancerPoolsPolicyTypeConfiguration
		if !plan.ConfigMaxPools.IsNull() {
			maxPools := make(map[string]interface{})
			if !plan.ConfigMaxPools.MaxPools.IsNull() {
				maxPools["maxPools"] = plan.ConfigMaxPools.MaxPools.ValueString()
			}
			if len(maxPools) > 0 {
				configMap = maxPools
			}
		}

		// 17. config_max_price -> BudgetPolicyTypeConfiguration
		if !plan.ConfigMaxPrice.IsNull() {
			maxPrice := make(map[string]interface{})
			if !plan.ConfigMaxPrice.MaxPrice.IsNull() {
				// MaxPrice is a Number type, get the float value
				f, _ := plan.ConfigMaxPrice.MaxPrice.ValueBigFloat().Float64()
				maxPrice["maxPrice"] = f
			}
			if !plan.ConfigMaxPrice.MaxPriceCurrency.IsNull() {
				maxPrice["maxPriceCurrency"] = plan.ConfigMaxPrice.MaxPriceCurrency.ValueString()
			}
			if !plan.ConfigMaxPrice.MaxPriceUnit.IsNull() {
				maxPrice["maxPriceUnit"] = plan.ConfigMaxPrice.MaxPriceUnit.ValueString()
			}
			if len(maxPrice) > 0 {
				configMap = maxPrice
			}
		}

		// 18. config_max_routers -> RouterQuotaPolicyTypeConfiguration
		if !plan.ConfigMaxRouters.IsNull() {
			maxRouters := make(map[string]interface{})
			if !plan.ConfigMaxRouters.MaxRouters.IsNull() {
				maxRouters["maxRouters"] = plan.ConfigMaxRouters.MaxRouters.ValueString()
			}
			if len(maxRouters) > 0 {
				configMap = maxRouters
			}
		}

		// 19. config_max_snapshots -> MaxSnapshotsPolicyTypeConfiguration
		if !plan.ConfigMaxSnapshots.IsNull() {
			maxSnapshots := make(map[string]interface{})
			if !plan.ConfigMaxSnapshots.MaxSnapshots.IsNull() {
				maxSnapshots["maxSnapshots"] = plan.ConfigMaxSnapshots.MaxSnapshots.ValueString()
			}
			if len(maxSnapshots) > 0 {
				configMap = maxSnapshots
			}
		}

		// 20. config_max_storage -> MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration
		if !plan.ConfigMaxStorage.IsNull() {
			maxStorage := make(map[string]interface{})
			if !plan.ConfigMaxStorage.MaxStorage.IsNull() {
				maxStorage["maxStorage"] = plan.ConfigMaxStorage.MaxStorage.ValueString()
			}
			if !plan.ConfigMaxStorage.ExcludeContainers.IsNull() {
				maxStorage["excludeContainers"] = fmt.Sprintf("%t", plan.ConfigMaxStorage.ExcludeContainers.ValueBool())
			}
			if len(maxStorage) > 0 {
				configMap = maxStorage
			}
		}

		// 21. config_max_virtual_servers -> MaxVirtualServersPolicyTypeConfiguration
		if !plan.ConfigMaxVirtualServers.IsNull() {
			maxVirtualServers := make(map[string]interface{})
			if !plan.ConfigMaxVirtualServers.MaxVirtualServers.IsNull() {
				maxVirtualServers["maxVirtualServers"] = plan.ConfigMaxVirtualServers.MaxVirtualServers.ValueString()
			}
			if len(maxVirtualServers) > 0 {
				configMap = maxVirtualServers
			}
		}

		// 22. config_max_vms -> MaxVMsPolicyTypeConfiguration
		if !plan.ConfigMaxVms.IsNull() {
			maxVms := make(map[string]interface{})
			if !plan.ConfigMaxVms.MaxVms.IsNull() {
				maxVms["maxVms"] = plan.ConfigMaxVms.MaxVms.ValueString()
			}
			if len(maxVms) > 0 {
				configMap = maxVms
			}
		}

		// 23. config_motd -> MessageOfTheDayPolicyTypeConfiguration2
		if !plan.ConfigMotd.IsNull() {
			motd := make(map[string]interface{})
			if !plan.ConfigMotd.Motddate.IsNull() {
				motd["motdDate"] = plan.ConfigMotd.Motddate.ValueString()
			}
			if !plan.ConfigMotd.Motdmessage.IsNull() {
				motd["motdMessage"] = plan.ConfigMotd.Motdmessage.ValueString()
			}
			if !plan.ConfigMotd.Motdtitle.IsNull() {
				motd["motdTitle"] = plan.ConfigMotd.Motdtitle.ValueString()
			}
			if !plan.ConfigMotd.Motdtype.IsNull() {
				motd["motdType"] = plan.ConfigMotd.Motdtype.ValueString()
			}
			if len(motd) > 0 {
				configMap = motd
			}
		}

		// 24. config_naming -> InstanceNamePolicyTypeConfiguration
		if !plan.ConfigNaming.IsNull() {
			naming := make(map[string]interface{})
			if !plan.ConfigNaming.NamingConflict.IsNull() {
				naming["namingConflict"] = plan.ConfigNaming.NamingConflict.ValueBool()
			}
			if !plan.ConfigNaming.NamingPattern.IsNull() {
				naming["namingPattern"] = plan.ConfigNaming.NamingPattern.ValueString()
			}
			if !plan.ConfigNaming.NamingType.IsNull() {
				naming["namingType"] = plan.ConfigNaming.NamingType.ValueString()
			}
			if len(naming) > 0 {
				configMap = naming
			}
		}

		// 25. config_power_schedule -> PowerSchedulePolicyTypeConfiguration
		if !plan.ConfigPowerSchedule.IsNull() {
			powerSchedule := make(map[string]interface{})
			if !plan.ConfigPowerSchedule.PowerSchedule.IsNull() {
				powerSchedule["powerSchedule"] = plan.ConfigPowerSchedule.PowerSchedule.ValueString()
			}
			if !plan.ConfigPowerSchedule.PowerScheduleHideFixed.IsNull() {
				powerSchedule["powerScheduleHideFixed"] = plan.ConfigPowerSchedule.PowerScheduleHideFixed.ValueBool()
			}
			if !plan.ConfigPowerSchedule.PowerScheduleType.IsNull() {
				powerSchedule["powerScheduleType"] = plan.ConfigPowerSchedule.PowerScheduleType.ValueString()
			}
			if len(powerSchedule) > 0 {
				configMap = powerSchedule
			}
		}

		// 26. config_required_network -> RequiredNetworkPolicyTypeConfiguration
		if !plan.ConfigRequiredNetwork.IsNull() {
			requiredNetwork := make(map[string]interface{})
			if !plan.ConfigRequiredNetwork.RequiredNetworks.IsNull() {
				var requiredNetworks []int64
				diagsSet := plan.ConfigRequiredNetwork.RequiredNetworks.ElementsAs(ctx, &requiredNetworks, false)
				if diagsSet.HasError() {
					diags.Append(diagsSet...)
					return nil, diags
				}
				requiredNetwork["requiredNetworks"] = requiredNetworks
			}
			if len(requiredNetwork) > 0 {
				configMap = requiredNetwork
			}
		}

		// 27. config_server_naming -> ClusterResourceNamePolicyTypeConfiguration
		if !plan.ConfigServerNaming.IsNull() {
			serverNaming := make(map[string]interface{})
			if !plan.ConfigServerNaming.ServerNamingConflict.IsNull() {
				serverNaming["serverNamingConflict"] = plan.ConfigServerNaming.ServerNamingConflict.ValueBool()
			}
			if !plan.ConfigServerNaming.ServerNamingPattern.IsNull() {
				serverNaming["serverNamingPattern"] = plan.ConfigServerNaming.ServerNamingPattern.ValueString()
			}
			if !plan.ConfigServerNaming.ServerNamingType.IsNull() {
				serverNaming["serverNamingType"] = plan.ConfigServerNaming.ServerNamingType.ValueString()
			}
			if len(serverNaming) > 0 {
				configMap = serverNaming
			}
		}

		// 28. config_shutdown -> ShutdownPolicyTypeConfiguration
		if !plan.ConfigShutdown.IsNull() {
			shutdown := make(map[string]interface{})
			if !plan.ConfigShutdown.AccountIntegrationId.IsNull() {
				shutdown["accountIntegrationId"] = plan.ConfigShutdown.AccountIntegrationId.ValueString()
			}
			if !plan.ConfigShutdown.FlowId.IsNull() {
				shutdown["flowId"] = plan.ConfigShutdown.FlowId.ValueString()
			}
			if !plan.ConfigShutdown.ShutdownAge.IsNull() {
				shutdown["shutdownAge"] = plan.ConfigShutdown.ShutdownAge.ValueString()
			}
			if !plan.ConfigShutdown.ShutdownAllowExtend.IsNull() {
				shutdown["shutdownAllowExtend"] = plan.ConfigShutdown.ShutdownAllowExtend.ValueBool()
			}
			if !plan.ConfigShutdown.ShutdownAutoRenew.IsNull() {
				shutdown["shutdownAutoRenew"] = plan.ConfigShutdown.ShutdownAutoRenew.ValueBool()
			}
			if !plan.ConfigShutdown.ShutdownExtensionsBeforeApproval.IsNull() {
				shutdown["shutdownExtensionsBeforeApproval"] = plan.ConfigShutdown.ShutdownExtensionsBeforeApproval.ValueString()
			}
			if !plan.ConfigShutdown.ShutdownHideFixed.IsNull() {
				shutdown["shutdownHideFixed"] = plan.ConfigShutdown.ShutdownHideFixed.ValueBool()
			}
			if !plan.ConfigShutdown.ShutdownMessage.IsNull() {
				shutdown["shutdownMessage"] = plan.ConfigShutdown.ShutdownMessage.ValueString()
			}
			if !plan.ConfigShutdown.ShutdownNotify.IsNull() {
				shutdown["shutdownNotify"] = plan.ConfigShutdown.ShutdownNotify.ValueString()
			}
			if !plan.ConfigShutdown.ShutdownRenewal.IsNull() {
				shutdown["shutdownRenewal"] = plan.ConfigShutdown.ShutdownRenewal.ValueString()
			}
			if !plan.ConfigShutdown.ShutdownType.IsNull() {
				shutdown["shutdownType"] = plan.ConfigShutdown.ShutdownType.ValueString()
			}
			if !plan.ConfigShutdown.ShutdownWorkflowId.IsNull() {
				shutdown["shutdownWorkflowId"] = plan.ConfigShutdown.ShutdownWorkflowId.ValueString()
			}
			if !plan.ConfigShutdown.WorkflowType.IsNull() {
				shutdown["workflowType"] = plan.ConfigShutdown.WorkflowType.ValueString()
			}
			if len(shutdown) > 0 {
				configMap = shutdown
			}
		}

		// 29. config_storage_server_quota -> StorageServerStorageQuotaPolicyTypeConfiguration
		if !plan.ConfigStorageServerQuota.IsNull() {
			storageServerQuota := make(map[string]interface{})
			if !plan.ConfigStorageServerQuota.MaxStorage.IsNull() {
				storageServerQuota["maxStorage"] = plan.ConfigStorageServerQuota.MaxStorage.ValueString()
			}
			if !plan.ConfigStorageServerQuota.StorageServerId.IsNull() {
				storageServerQuota["storageServerId"] = plan.ConfigStorageServerQuota.StorageServerId.ValueString()
			}
			if len(storageServerQuota) > 0 {
				configMap = storageServerQuota
			}
		}

		// 30. config_tags -> TagsPolicyTypeConfiguration
		if !plan.ConfigTags.IsNull() {
			tags := make(map[string]interface{})
			if !plan.ConfigTags.Key.IsNull() {
				tags["key"] = plan.ConfigTags.Key.ValueString()
			}
			if !plan.ConfigTags.Strict.IsNull() {
				tags["strict"] = plan.ConfigTags.Strict.ValueBool()
			}
			if !plan.ConfigTags.Value.IsNull() {
				tags["value"] = plan.ConfigTags.Value.ValueString()
			}
			if !plan.ConfigTags.ValueListId.IsNull() {
				tags["valueListId"] = plan.ConfigTags.ValueListId.ValueString()
			}
			if len(tags) > 0 {
				configMap = tags
			}
		}

		// 31. config_workflow -> WorkflowPolicyTypeConfiguration
		if !plan.ConfigWorkflow.IsNull() {
			workflow := make(map[string]interface{})
			if !plan.ConfigWorkflow.WorkflowId.IsNull() {
				workflow["workflowId"] = plan.ConfigWorkflow.WorkflowId.ValueString()
			}
			if len(workflow) > 0 {
				configMap = workflow
			}
		}

		// Check if any config was provided
		if len(configMap) == 0 {
			diags.AddError(
				"map state to API config",
				"No configuration provided. Please provide at least one config_* field for this policy type.",
			)
			return nil, diags
		}
	}

	// Marshal to JSON then unmarshal to SDK config structure
	configJSON, err := json.Marshal(configMap)
	if err != nil {
		diags.AddError(
			"map state to API config",
			"failed to marshal config to JSON: "+err.Error(),
		)
		return nil, diags
	}

	var sdkConfig sdk.AddPoliciesRequestPolicyConfig
	if err := json.Unmarshal(configJSON, &sdkConfig); err != nil {
		diags.AddError(
			"map state to API config",
			"invalid config: "+err.Error(),
		)
		return nil, diags
	}

	return &sdkConfig, diags
}

// mapStateToUpdatePolicyConfig converts state config fields to API config structure for Update
func mapStateToUpdatePolicyConfig(
	ctx context.Context,
	plan *PolicyModel,
) (*sdk.UpdatePoliciesRequestPolicyConfig, diag.Diagnostics) {
	addConfig, diags := mapStateToAddPolicyConfig(ctx, plan)
	if diags.HasError() {
		return nil, diags
	}

	// Convert AddPoliciesRequestPolicyConfig to UpdatePoliciesRequestPolicyConfig via JSON
	configJSON, err := json.Marshal(addConfig)
	if err != nil {
		diags.AddError(
			"map state to update API config",
			"failed to marshal config to JSON: "+err.Error(),
		)
		return nil, diags
	}

	var updateConfig sdk.UpdatePoliciesRequestPolicyConfig
	if err := json.Unmarshal(configJSON, &updateConfig); err != nil {
		diags.AddError(
			"map state to update API config",
			"invalid config: "+err.Error(),
		)
		return nil, diags
	}

	return &updateConfig, diags
}
