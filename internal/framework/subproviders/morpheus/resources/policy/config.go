// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"encoding/json"

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
		// Use static config fields based on policy type code
		configMap = make(map[string]interface{})

		policyTypeCode := plan.PolicyType.Code.ValueString()

		// Map config based on policy type code
		switch policyTypeCode {
		case "deleteApproval", "provisionApproval", "reconfigureApproval", "workflowApproval":
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

		case "backupStorage":
			backupStorage := make(map[string]interface{})
			if !plan.ConfigBackupStorage.BackupStorageIds.IsNull() {
				var backupStorageIDs []int64
				diagSet := plan.ConfigBackupStorage.BackupStorageIds.ElementsAs(ctx, &backupStorageIDs, false)
				if diagSet.HasError() {
					diags.Append(diagSet...)

					return nil, diags
				}
				backupStorage["backupStorageIds"] = backupStorageIDs
			}
			if len(backupStorage) > 0 {
				configMap = backupStorage
			}

		case "createBackup":
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

		case "createUser":
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

		case "createUserGroup":
			createUserGroup := make(map[string]interface{})
			if !plan.ConfigCreateUserGroup.UserGroup.IsNull() {
				createUserGroup["userGroup"] = plan.ConfigCreateUserGroup.UserGroup.ValueString()
			}
			if len(createUserGroup) > 0 {
				configMap = createUserGroup
			}

		case "cypher":
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

		case "delayedRemoval":
			delayedRemoval := make(map[string]interface{})
			if !plan.ConfigDelayedRemoval.RemovalAge.IsNull() {
				delayedRemoval["removalAge"] = plan.ConfigDelayedRemoval.RemovalAge.ValueString()
			}
			if len(delayedRemoval) > 0 {
				configMap = delayedRemoval
			}

		case "hostNaming":
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

		case "lifecycle":
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
				lifecycle["lifecycleAllowExtend"] = convert.BoolToStringOnOff(
					plan.ConfigLifecycle.LifecycleAllowExtend.ValueBool(),
				).ValueString()
			}
			if !plan.ConfigLifecycle.LifecycleAutoRenew.IsNull() {
				lifecycle["lifecycleAutoRenew"] = convert.BoolToStringOnOff(
					plan.ConfigLifecycle.LifecycleAutoRenew.ValueBool(),
				).ValueString()
			}
			if !plan.ConfigLifecycle.LifecycleExtensionsBeforeApproval.IsNull() {
				lifecycle["lifecycleExtensionsBeforeApproval"] = plan.ConfigLifecycle.
					LifecycleExtensionsBeforeApproval.ValueString()
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

		case "maxContainers":
			maxContainers := make(map[string]interface{})
			if !plan.ConfigMaxContainers.MaxContainers.IsNull() {
				maxContainers["maxContainers"] = plan.ConfigMaxContainers.MaxContainers.ValueString()
			}
			if len(maxContainers) > 0 {
				configMap = maxContainers
			}

		case "maxCores":
			maxCores := make(map[string]interface{})
			if !plan.ConfigMaxCores.MaxCores.IsNull() {
				maxCores["maxCores"] = plan.ConfigMaxCores.MaxCores.ValueString()
			}
			if !plan.ConfigMaxCores.ExcludeContainers.IsNull() {
				maxCores["excludeContainers"] = convert.BoolToStringOnOff(
					plan.ConfigMaxCores.ExcludeContainers.ValueBool(),
				).ValueString()
			}
			if len(maxCores) > 0 {
				configMap = maxCores
			}

		case "maxHosts":
			maxHosts := make(map[string]interface{})
			if !plan.ConfigMaxHosts.MaxHosts.IsNull() {
				maxHosts["maxHosts"] = plan.ConfigMaxHosts.MaxHosts.ValueString()
			}
			if len(maxHosts) > 0 {
				configMap = maxHosts
			}

		case "maxMemory":
			maxMemory := make(map[string]interface{})
			if !plan.ConfigMaxMemory.MaxMemory.IsNull() {
				maxMemory["maxMemory"] = plan.ConfigMaxMemory.MaxMemory.ValueString()
			}
			if !plan.ConfigMaxMemory.ExcludeContainers.IsNull() {
				maxMemory["excludeContainers"] = convert.BoolToStringOnOff(
					plan.ConfigMaxMemory.ExcludeContainers.ValueBool(),
				).ValueString()
			}
			if len(maxMemory) > 0 {
				configMap = maxMemory
			}

		case "maxNetworks":
			maxNetworks := make(map[string]interface{})
			if !plan.ConfigMaxNetworks.MaxNetworks.IsNull() {
				maxNetworks["maxNetworks"] = plan.ConfigMaxNetworks.MaxNetworks.ValueString()
			}
			if len(maxNetworks) > 0 {
				configMap = maxNetworks
			}

		case "maxPoolMembers":
			maxPoolMembers := make(map[string]interface{})
			if !plan.ConfigMaxPoolMembers.MaxPoolMembers.IsNull() {
				maxPoolMembers["maxPoolMembers"] = plan.ConfigMaxPoolMembers.MaxPoolMembers.ValueString()
			}
			if len(maxPoolMembers) > 0 {
				configMap = maxPoolMembers
			}

		case "maxPools":
			maxPools := make(map[string]interface{})
			if !plan.ConfigMaxPools.MaxPools.IsNull() {
				maxPools["maxPools"] = plan.ConfigMaxPools.MaxPools.ValueString()
			}
			if len(maxPools) > 0 {
				configMap = maxPools
			}

		case "maxPrice":
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

		case "maxRouters":
			maxRouters := make(map[string]interface{})
			if !plan.ConfigMaxRouters.MaxRouters.IsNull() {
				maxRouters["maxRouters"] = plan.ConfigMaxRouters.MaxRouters.ValueString()
			}
			if len(maxRouters) > 0 {
				configMap = maxRouters
			}

		case "maxSnapshots":
			maxSnapshots := make(map[string]interface{})
			if !plan.ConfigMaxSnapshots.MaxSnapshots.IsNull() {
				maxSnapshots["maxSnapshots"] = plan.ConfigMaxSnapshots.MaxSnapshots.ValueString()
			}
			if len(maxSnapshots) > 0 {
				configMap = maxSnapshots
			}

		case "maxStorage":
			maxStorage := make(map[string]interface{})
			if !plan.ConfigMaxStorage.MaxStorage.IsNull() {
				maxStorage["maxStorage"] = plan.ConfigMaxStorage.MaxStorage.ValueString()
			}
			if !plan.ConfigMaxStorage.ExcludeContainers.IsNull() {
				maxStorage["excludeContainers"] = convert.BoolToStringOnOff(
					plan.ConfigMaxStorage.ExcludeContainers.ValueBool(),
				).ValueString()
			}
			if len(maxStorage) > 0 {
				configMap = maxStorage
			}

		case "maxVirtualServers":
			maxVirtualServers := make(map[string]interface{})
			if !plan.ConfigMaxVirtualServers.MaxVirtualServers.IsNull() {
				maxVirtualServers["maxVirtualServers"] = plan.ConfigMaxVirtualServers.MaxVirtualServers.ValueString()
			}
			if len(maxVirtualServers) > 0 {
				configMap = maxVirtualServers
			}

		case "maxVms":
			maxVms := make(map[string]interface{})
			if !plan.ConfigMaxVms.MaxVms.IsNull() {
				maxVms["maxVms"] = plan.ConfigMaxVms.MaxVms.ValueString()
			}
			if len(maxVms) > 0 {
				configMap = maxVms
			}

		case "motd":
			motd := make(map[string]interface{})
			if !plan.ConfigMotd.Motdtitle.IsNull() {
				motd["motd.title"] = plan.ConfigMotd.Motdtitle.ValueString()
			}
			if !plan.ConfigMotd.Motdmessage.IsNull() {
				motd["motd.message"] = plan.ConfigMotd.Motdmessage.ValueString()
			}
			if !plan.ConfigMotd.Motdtype.IsNull() {
				motd["motd.type"] = plan.ConfigMotd.Motdtype.ValueString()
			}
			if !plan.ConfigMotd.Motddate.IsNull() {
				motd["motd.date"] = plan.ConfigMotd.Motddate.ValueString()
			}
			if len(motd) > 0 {
				configMap = motd
			}

		case "naming":
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

		case "powerSchedule":
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

		case "requiredNetwork":
			requiredNetwork := make(map[string]interface{})
			if !plan.ConfigRequiredNetwork.RequiredNetworks.IsNull() {
				var requiredNetworks []int64
				diagSet := plan.ConfigRequiredNetwork.RequiredNetworks.ElementsAs(ctx, &requiredNetworks, false)
				if diagSet.HasError() {
					diags.Append(diagSet...)

					return nil, diags
				}
				requiredNetwork["requiredNetworks"] = requiredNetworks
			}
			if len(requiredNetwork) > 0 {
				configMap = requiredNetwork
			}

		case "serverNaming":
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

		case "shutdown":
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
				shutdown["shutdownAllowExtend"] = convert.BoolToStringOnOff(
					plan.ConfigShutdown.ShutdownAllowExtend.ValueBool(),
				).ValueString()
			}
			if !plan.ConfigShutdown.ShutdownAutoRenew.IsNull() {
				shutdown["shutdownAutoRenew"] = convert.BoolToStringOnOff(
					plan.ConfigShutdown.ShutdownAutoRenew.ValueBool(),
				).ValueString()
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

		case "storageServerQuota":
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

		case "tags":
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

		case "workflow":
			workflow := make(map[string]interface{})
			if !plan.ConfigWorkflow.WorkflowId.IsNull() {
				workflow["workflowId"] = plan.ConfigWorkflow.WorkflowId.ValueString()
			}
			if len(workflow) > 0 {
				configMap = workflow
			}

		default:

			// No config was provided
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
