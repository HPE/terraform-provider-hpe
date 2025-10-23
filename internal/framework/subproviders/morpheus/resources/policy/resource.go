// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// policy provides the package for hpe_morpheus_policy
package policy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_policy"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = PolicyResourceSchema(ctx)
}

// populate policy resource model with current API values
func getPolicyAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
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

	// Handle RefId - it might be a string or int64 depending on the API
	// For now, we'll skip complex parsing and handle basic cases

	// Set account IDs
	if p.Accounts != nil {
		accountIDs := make([]int64, len(p.Accounts))
		for i, acc := range p.Accounts {
			if acc.Id != nil {
				accountIDs[i] = *acc.Id
			}
		}
		accountSet, setDiags := types.SetValueFrom(ctx, types.Int64Type, accountIDs)
		if setDiags.HasError() {
			diags.Append(setDiags...)
			return state, diags
		}
		state.Accounts = accountSet
	}

	// Set PolicyType
	if p.PolicyType != nil {
		policyTypeAttrs := map[string]attr.Value{}
		if p.PolicyType.Code != nil {
			policyTypeAttrs["code"] = types.StringValue(*p.PolicyType.Code)
		} else {
			policyTypeAttrs["code"] = types.StringNull()
		}

		policyTypeValue, policyTypeDiags := NewPolicyTypeValue(
			PolicyTypeValue{}.AttributeTypes(ctx), policyTypeAttrs)
		if policyTypeDiags.HasError() {
			diags.Append(policyTypeDiags...)
			return state, diags
		}
		state.PolicyType = policyTypeValue
	}

	// TODO: Set complex config fields based on API response
	// This requires mapping from API response config to the various schema config structures
	// Each config type (approval, backup creation, budget, etc.) needs to be handled
	// For now, leaving as null values since the exact API response structure may vary

	return state, diags
}

// buildPolicyConfigForCreate maps the schema config fields to the SDK config structure for create operations
func buildPolicyConfigForCreate(ctx context.Context, plan *PolicyModel) (*sdk.AddPoliciesRequestPolicyConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	config := &sdk.AddPoliciesRequestPolicyConfig{}
	hasConfig := false

	// ConfigApproval mapping
	if !plan.ConfigApproval.IsNull() && !plan.ConfigApproval.IsUnknown() {
		approveConfig := &sdk.ApprovePolicyTypeConfiguration{}
		if !plan.ConfigApproval.AccountIntegrationId.IsNull() && !plan.ConfigApproval.AccountIntegrationId.IsUnknown() {
			approveConfig.SetAccountIntegrationId(plan.ConfigApproval.AccountIntegrationId.ValueString())
		}
		config.ApprovePolicyTypeConfiguration = approveConfig
		hasConfig = true
	}

	// ConfigBackupCreation mapping
	if !plan.ConfigBackupCreation.IsNull() && !plan.ConfigBackupCreation.IsUnknown() {
		backupConfig := &sdk.BackupCreationPolicyTypeConfiguration{}
		if !plan.ConfigBackupCreation.CreateBackup.IsNull() && !plan.ConfigBackupCreation.CreateBackup.IsUnknown() {
			backupConfig.SetCreateBackup(plan.ConfigBackupCreation.CreateBackup.ValueBool())
		}
		if !plan.ConfigBackupCreation.CreateBackupType.IsNull() && !plan.ConfigBackupCreation.CreateBackupType.IsUnknown() {
			backupConfig.SetCreateBackupType(plan.ConfigBackupCreation.CreateBackupType.ValueString())
		}
		config.BackupCreationPolicyTypeConfiguration = backupConfig
		hasConfig = true
	}

	// ConfigBackupTargets mapping
	if !plan.ConfigBackupTargets.IsNull() && !plan.ConfigBackupTargets.IsUnknown() {
		backupTargetsConfig := &sdk.BackupTargetsPolicyTypeConfiguration{}
		if !plan.ConfigBackupTargets.BackupStorageIds.IsNull() && !plan.ConfigBackupTargets.BackupStorageIds.IsUnknown() {
			var storageIds []int64
			elemDiags := plan.ConfigBackupTargets.BackupStorageIds.ElementsAs(ctx, &storageIds, false)
			if elemDiags.HasError() {
				diags.Append(elemDiags...)
			} else {
				backupTargetsConfig.SetBackupStorageIds(storageIds)
			}
		}
		config.BackupTargetsPolicyTypeConfiguration = backupTargetsConfig
		hasConfig = true
	}

	// ConfigBudget mapping
	if !plan.ConfigBudget.IsNull() && !plan.ConfigBudget.IsUnknown() {
		budgetConfig := &sdk.BudgetPolicyTypeConfiguration{}
		if !plan.ConfigBudget.MaxPrice.IsNull() && !plan.ConfigBudget.MaxPrice.IsUnknown() {
			maxPrice, _ := plan.ConfigBudget.MaxPrice.ValueBigFloat().Float64()
			budgetConfig.SetMaxPrice(float32(maxPrice))
		}
		if !plan.ConfigBudget.MaxPriceCurrency.IsNull() && !plan.ConfigBudget.MaxPriceCurrency.IsUnknown() {
			budgetConfig.SetMaxPriceCurrency(plan.ConfigBudget.MaxPriceCurrency.ValueString())
		}
		if !plan.ConfigBudget.MaxPriceUnit.IsNull() && !plan.ConfigBudget.MaxPriceUnit.IsUnknown() {
			budgetConfig.SetMaxPriceUnit(plan.ConfigBudget.MaxPriceUnit.ValueString())
		}
		config.BudgetPolicyTypeConfiguration = budgetConfig
		hasConfig = true
	}

	// ConfigClusterResourceName mapping
	if !plan.ConfigClusterResourceName.IsNull() && !plan.ConfigClusterResourceName.IsUnknown() {
		clusterConfig := &sdk.ClusterResourceNamePolicyTypeConfiguration{}
		if !plan.ConfigClusterResourceName.ServerNamingType.IsNull() && !plan.ConfigClusterResourceName.ServerNamingType.IsUnknown() {
			clusterConfig.SetServerNamingType(plan.ConfigClusterResourceName.ServerNamingType.ValueString())
		}
		if !plan.ConfigClusterResourceName.ServerNamingPattern.IsNull() && !plan.ConfigClusterResourceName.ServerNamingPattern.IsUnknown() {
			clusterConfig.SetServerNamingPattern(plan.ConfigClusterResourceName.ServerNamingPattern.ValueString())
		}
		if !plan.ConfigClusterResourceName.ServerNamingConflict.IsNull() && !plan.ConfigClusterResourceName.ServerNamingConflict.IsUnknown() {
			clusterConfig.SetServerNamingConflict(plan.ConfigClusterResourceName.ServerNamingConflict.ValueBool())
		}
		config.ClusterResourceNamePolicyTypeConfiguration = clusterConfig
		hasConfig = true
	}

	// ConfigCypherAccess mapping
	if !plan.ConfigCypherAccess.IsNull() && !plan.ConfigCypherAccess.IsUnknown() {
		cypherConfig := &sdk.CypherAccessPolicyTypeConfiguration{}
		if !plan.ConfigCypherAccess.KeyPattern.IsNull() && !plan.ConfigCypherAccess.KeyPattern.IsUnknown() {
			cypherConfig.SetKeyPattern(plan.ConfigCypherAccess.KeyPattern.ValueString())
		}
		if !plan.ConfigCypherAccess.Read.IsNull() && !plan.ConfigCypherAccess.Read.IsUnknown() {
			cypherConfig.SetRead(plan.ConfigCypherAccess.Read.ValueBool())
		}
		if !plan.ConfigCypherAccess.Write.IsNull() && !plan.ConfigCypherAccess.Write.IsUnknown() {
			cypherConfig.SetWrite(plan.ConfigCypherAccess.Write.ValueBool())
		}
		if !plan.ConfigCypherAccess.Update.IsNull() && !plan.ConfigCypherAccess.Update.IsUnknown() {
			cypherConfig.SetUpdate(plan.ConfigCypherAccess.Update.ValueBool())
		}
		if !plan.ConfigCypherAccess.Delete.IsNull() && !plan.ConfigCypherAccess.Delete.IsUnknown() {
			cypherConfig.SetDelete(plan.ConfigCypherAccess.Delete.ValueBool())
		}
		if !plan.ConfigCypherAccess.List.IsNull() && !plan.ConfigCypherAccess.List.IsUnknown() {
			cypherConfig.SetList(plan.ConfigCypherAccess.List.ValueBool())
		}
		config.CypherAccessPolicyTypeConfiguration = cypherConfig
		hasConfig = true
	}

	// ConfigDelayedDelete mapping
	if !plan.ConfigDelayedDelete.IsNull() && !plan.ConfigDelayedDelete.IsUnknown() {
		delayedDeleteConfig := &sdk.DelayedDeletePolicyTypeConfiguration{}
		if !plan.ConfigDelayedDelete.RemovalAge.IsNull() && !plan.ConfigDelayedDelete.RemovalAge.IsUnknown() {
			delayedDeleteConfig.SetRemovalAge(plan.ConfigDelayedDelete.RemovalAge.ValueString())
		}
		config.DelayedDeletePolicyTypeConfiguration = delayedDeleteConfig
		hasConfig = true
	}

	// Continue with more config mappings...

	// ConfigExpiration mapping
	if !plan.ConfigExpiration.IsNull() && !plan.ConfigExpiration.IsUnknown() {
		expirationConfig := &sdk.ExpirationPolicyTypeConfiguration{}
		if !plan.ConfigExpiration.LifecycleType.IsNull() && !plan.ConfigExpiration.LifecycleType.IsUnknown() {
			expirationConfig.SetLifecycleType(plan.ConfigExpiration.LifecycleType.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleAge.IsNull() && !plan.ConfigExpiration.LifecycleAge.IsUnknown() {
			expirationConfig.SetLifecycleAge(plan.ConfigExpiration.LifecycleAge.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleAllowExtend.IsNull() && !plan.ConfigExpiration.LifecycleAllowExtend.IsUnknown() {
			expirationConfig.SetLifecycleAllowExtend(plan.ConfigExpiration.LifecycleAllowExtend.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleAutoRenew.IsNull() && !plan.ConfigExpiration.LifecycleAutoRenew.IsUnknown() {
			expirationConfig.SetLifecycleAutoRenew(plan.ConfigExpiration.LifecycleAutoRenew.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleExtensionsBeforeApproval.IsNull() && !plan.ConfigExpiration.LifecycleExtensionsBeforeApproval.IsUnknown() {
			expirationConfig.SetLifecycleExtensionsBeforeApproval(plan.ConfigExpiration.LifecycleExtensionsBeforeApproval.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleHideFixed.IsNull() && !plan.ConfigExpiration.LifecycleHideFixed.IsUnknown() {
			expirationConfig.SetLifecycleHideFixed(plan.ConfigExpiration.LifecycleHideFixed.ValueBool())
		}
		if !plan.ConfigExpiration.LifecycleMessage.IsNull() && !plan.ConfigExpiration.LifecycleMessage.IsUnknown() {
			expirationConfig.SetLifecycleMessage(plan.ConfigExpiration.LifecycleMessage.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleNotify.IsNull() && !plan.ConfigExpiration.LifecycleNotify.IsUnknown() {
			expirationConfig.SetLifecycleNotify(plan.ConfigExpiration.LifecycleNotify.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleRenewal.IsNull() && !plan.ConfigExpiration.LifecycleRenewal.IsUnknown() {
			expirationConfig.SetLifecycleRenewal(plan.ConfigExpiration.LifecycleRenewal.ValueString())
		}
		if !plan.ConfigExpiration.AccountIntegrationId.IsNull() && !plan.ConfigExpiration.AccountIntegrationId.IsUnknown() {
			expirationConfig.SetAccountIntegrationId(plan.ConfigExpiration.AccountIntegrationId.ValueString())
		}
		config.ExpirationPolicyTypeConfiguration = expirationConfig
		hasConfig = true
	}

	// ConfigHostname mapping
	if !plan.ConfigHostname.IsNull() && !plan.ConfigHostname.IsUnknown() {
		hostnameConfig := &sdk.HostnamePolicyTypeConfiguration{}
		if !plan.ConfigHostname.HostNamingType.IsNull() && !plan.ConfigHostname.HostNamingType.IsUnknown() {
			hostnameConfig.SetHostNamingType(plan.ConfigHostname.HostNamingType.ValueString())
		}
		if !plan.ConfigHostname.HostNamingPattern.IsNull() && !plan.ConfigHostname.HostNamingPattern.IsUnknown() {
			hostnameConfig.SetHostNamingPattern(plan.ConfigHostname.HostNamingPattern.ValueString())
		}
		config.HostnamePolicyTypeConfiguration = hostnameConfig
		hasConfig = true
	}

	// ConfigInstanceName mapping
	if !plan.ConfigInstanceName.IsNull() && !plan.ConfigInstanceName.IsUnknown() {
		instanceNameConfig := &sdk.InstanceNamePolicyTypeConfiguration{}
		if !plan.ConfigInstanceName.NamingType.IsNull() && !plan.ConfigInstanceName.NamingType.IsUnknown() {
			instanceNameConfig.SetNamingType(plan.ConfigInstanceName.NamingType.ValueString())
		}
		if !plan.ConfigInstanceName.NamingPattern.IsNull() && !plan.ConfigInstanceName.NamingPattern.IsUnknown() {
			instanceNameConfig.SetNamingPattern(plan.ConfigInstanceName.NamingPattern.ValueString())
		}
		if !plan.ConfigInstanceName.NamingConflict.IsNull() && !plan.ConfigInstanceName.NamingConflict.IsUnknown() {
			instanceNameConfig.SetNamingConflict(plan.ConfigInstanceName.NamingConflict.ValueBool())
		}
		config.InstanceNamePolicyTypeConfiguration = instanceNameConfig
		hasConfig = true
	}

	// ConfigMaxContainers mapping
	if !plan.ConfigMaxContainers.IsNull() && !plan.ConfigMaxContainers.IsUnknown() {
		maxContainersConfig := &sdk.MaxContainersPolicyTypeConfiguration{}
		if !plan.ConfigMaxContainers.MaxContainers.IsNull() && !plan.ConfigMaxContainers.MaxContainers.IsUnknown() {
			maxContainersConfig.SetMaxContainers(plan.ConfigMaxContainers.MaxContainers.ValueString())
		}
		config.MaxContainersPolicyTypeConfiguration = maxContainersConfig
		hasConfig = true
	}

	// ConfigMaxCores mapping
	if !plan.ConfigMaxCores.IsNull() && !plan.ConfigMaxCores.IsUnknown() {
		maxCoresConfig := &sdk.MaxCoresPolicyTypeConfiguration{}
		if !plan.ConfigMaxCores.MaxCores.IsNull() && !plan.ConfigMaxCores.MaxCores.IsUnknown() {
			maxCoresConfig.SetMaxCores(plan.ConfigMaxCores.MaxCores.ValueString())
		}
		if !plan.ConfigMaxCores.ExcludeContainers.IsNull() && !plan.ConfigMaxCores.ExcludeContainers.IsUnknown() {
			maxCoresConfig.SetExcludeContainers(plan.ConfigMaxCores.ExcludeContainers.ValueString())
		}
		config.MaxCoresPolicyTypeConfiguration = maxCoresConfig
		hasConfig = true
	}

	// ConfigMaxHosts mapping
	if !plan.ConfigMaxHosts.IsNull() && !plan.ConfigMaxHosts.IsUnknown() {
		maxHostsConfig := &sdk.MaxHostsPolicyTypeConfiguration{}
		if !plan.ConfigMaxHosts.MaxHosts.IsNull() && !plan.ConfigMaxHosts.MaxHosts.IsUnknown() {
			maxHostsConfig.SetMaxHosts(plan.ConfigMaxHosts.MaxHosts.ValueString())
		}
		config.MaxHostsPolicyTypeConfiguration = maxHostsConfig
		hasConfig = true
	}

	// ConfigMaxLoadBalancerPools mapping
	if !plan.ConfigMaxLoadBalancerPools.IsNull() && !plan.ConfigMaxLoadBalancerPools.IsUnknown() {
		maxLBPoolsConfig := &sdk.MaxLoadBalancerPoolsPolicyTypeConfiguration{}
		if !plan.ConfigMaxLoadBalancerPools.MaxPools.IsNull() && !plan.ConfigMaxLoadBalancerPools.MaxPools.IsUnknown() {
			maxLBPoolsConfig.SetMaxPools(plan.ConfigMaxLoadBalancerPools.MaxPools.ValueString())
		}
		config.MaxLoadBalancerPoolsPolicyTypeConfiguration = maxLBPoolsConfig
		hasConfig = true
	}

	// ConfigMaxMemory mapping (complex nested structure)
	if !plan.ConfigMaxMemory.IsNull() && !plan.ConfigMaxMemory.IsUnknown() {
		maxMemoryConfig := &sdk.MaxMemoryPolicyTypeConfiguration{}
		maxMemorySet := false

		if !plan.ConfigMaxMemory.ExcludeContainers.IsNull() && !plan.ConfigMaxMemory.ExcludeContainers.IsUnknown() {
			maxMemoryConfig.SetExcludeContainers(plan.ConfigMaxMemory.ExcludeContainers.ValueString())
		}

		// Handle MaxMemory field - it's an ObjectValue with anyof0 (string) and anyof1 (int)
		if !plan.ConfigMaxMemory.MaxMemory.IsNull() && !plan.ConfigMaxMemory.MaxMemory.IsUnknown() {
			maxMemoryValue := plan.ConfigMaxMemory.MaxMemory
			// Extract the MaxMemoryValue from the ObjectValue
			if maxMemoryAttrs := maxMemoryValue.Attributes(); len(maxMemoryAttrs) > 0 {
				var maxMemory sdk.MaxMemoryPolicyTypeConfigurationMaxMemory
				// Check anyof0 (string) first
				if anyof0Val, ok := maxMemoryAttrs["anyof0"]; ok {
					if stringVal, ok := anyof0Val.(basetypes.StringValue); ok && !stringVal.IsNull() && !stringVal.IsUnknown() {
						// Use string representation (anyof0)
						str := stringVal.ValueString()
						maxMemory.String = &str
						maxMemoryConfig.SetMaxMemory(maxMemory)
						maxMemorySet = true
					}
				}
				// Check anyof1 (int64) if anyof0 wasn't set
				if !maxMemorySet {
					if anyof1Val, ok := maxMemoryAttrs["anyof1"]; ok {
						if intVal, ok := anyof1Val.(basetypes.Int64Value); ok && !intVal.IsNull() && !intVal.IsUnknown() {
							// Use int representation (anyof1)
							intValue := intVal.ValueInt64()
							maxMemory.Int64 = &intValue
							maxMemoryConfig.SetMaxMemory(maxMemory)
							maxMemorySet = true
						}
					}
				}
			}
		}

		// Only add the config if MaxMemory was actually set (it's required)
		if maxMemorySet {
			config.MaxMemoryPolicyTypeConfiguration = maxMemoryConfig
			hasConfig = true
		}
	}

	// ConfigMaxPoolMembers mapping
	if !plan.ConfigMaxPoolMembers.IsNull() && !plan.ConfigMaxPoolMembers.IsUnknown() {
		maxPoolMembersConfig := &sdk.MaxPoolMembersPolicyTypeConfiguration{}
		if !plan.ConfigMaxPoolMembers.MaxPoolMembers.IsNull() && !plan.ConfigMaxPoolMembers.MaxPoolMembers.IsUnknown() {
			maxPoolMembersConfig.SetMaxPoolMembers(plan.ConfigMaxPoolMembers.MaxPoolMembers.ValueString())
		}
		config.MaxPoolMembersPolicyTypeConfiguration = maxPoolMembersConfig
		hasConfig = true
	}

	// ConfigMaxStorage mapping
	if !plan.ConfigMaxStorage.IsNull() && !plan.ConfigMaxStorage.IsUnknown() {
		maxStorageConfig := &sdk.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration{}
		if !plan.ConfigMaxStorage.MaxStorage.IsNull() && !plan.ConfigMaxStorage.MaxStorage.IsUnknown() {
			maxStorageConfig.SetMaxStorage(plan.ConfigMaxStorage.MaxStorage.ValueString())
		}
		if !plan.ConfigMaxStorage.ExcludeContainers.IsNull() && !plan.ConfigMaxStorage.ExcludeContainers.IsUnknown() {
			maxStorageConfig.SetExcludeContainers(plan.ConfigMaxStorage.ExcludeContainers.ValueString())
		}
		config.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration = maxStorageConfig
		hasConfig = true
	}

	// ConfigMaxVirtualServers mapping
	if !plan.ConfigMaxVirtualServers.IsNull() && !plan.ConfigMaxVirtualServers.IsUnknown() {
		maxVirtualServersConfig := &sdk.MaxVirtualServersPolicyTypeConfiguration{}
		if !plan.ConfigMaxVirtualServers.MaxVirtualServers.IsNull() && !plan.ConfigMaxVirtualServers.MaxVirtualServers.IsUnknown() {
			maxVirtualServersConfig.SetMaxVirtualServers(plan.ConfigMaxVirtualServers.MaxVirtualServers.ValueString())
		}
		config.MaxVirtualServersPolicyTypeConfiguration = maxVirtualServersConfig
		hasConfig = true
	}

	// ConfigMaxVms mapping
	if !plan.ConfigMaxVms.IsNull() && !plan.ConfigMaxVms.IsUnknown() {
		maxVmsConfig := &sdk.MaxVMsPolicyTypeConfiguration{}
		if !plan.ConfigMaxVms.MaxVms.IsNull() && !plan.ConfigMaxVms.MaxVms.IsUnknown() {
			maxVmsConfig.SetMaxVms(plan.ConfigMaxVms.MaxVms.ValueString())
		}
		config.MaxVMsPolicyTypeConfiguration = maxVmsConfig
		hasConfig = true
	}

	// ConfigMessageOfTheDay mapping
	if !plan.ConfigMessageOfTheDay.IsNull() && !plan.ConfigMessageOfTheDay.IsUnknown() {
		motdConfig := &sdk.MessageOfTheDayPolicyTypeConfiguration{}
		if !plan.ConfigMessageOfTheDay.Motddate.IsNull() && !plan.ConfigMessageOfTheDay.Motddate.IsUnknown() {
			// Parse the date string - assuming RFC3339 format
			if dateStr := plan.ConfigMessageOfTheDay.Motddate.ValueString(); dateStr != "" {
				if parsedDate, err := time.Parse(time.RFC3339, dateStr); err == nil {
					motdConfig.SetMotdDate(parsedDate)
				}
				// If parsing fails, we could add to diagnostics, but for now continue without setting
			}
		}
		if !plan.ConfigMessageOfTheDay.Motdmessage.IsNull() && !plan.ConfigMessageOfTheDay.Motdmessage.IsUnknown() {
			motdConfig.SetMotdMessage(plan.ConfigMessageOfTheDay.Motdmessage.ValueString())
		}
		if !plan.ConfigMessageOfTheDay.Motdtitle.IsNull() && !plan.ConfigMessageOfTheDay.Motdtitle.IsUnknown() {
			motdConfig.SetMotdTitle(plan.ConfigMessageOfTheDay.Motdtitle.ValueString())
		}
		if !plan.ConfigMessageOfTheDay.Motdtype.IsNull() && !plan.ConfigMessageOfTheDay.Motdtype.IsUnknown() {
			motdConfig.SetMotdType(plan.ConfigMessageOfTheDay.Motdtype.ValueString())
		}
		if !plan.ConfigMessageOfTheDay.MotdFullPage.IsNull() && !plan.ConfigMessageOfTheDay.MotdFullPage.IsUnknown() {
			motdConfig.SetMotdFullPage(plan.ConfigMessageOfTheDay.MotdFullPage.ValueBool())
		}
		config.MessageOfTheDayPolicyTypeConfiguration = motdConfig
		hasConfig = true
	}

	// ConfigNetworkQuota mapping
	if !plan.ConfigNetworkQuota.IsNull() && !plan.ConfigNetworkQuota.IsUnknown() {
		networkQuotaConfig := &sdk.NetworkQuotaPolicyTypeConfiguration{}
		if !plan.ConfigNetworkQuota.MaxNetworks.IsNull() && !plan.ConfigNetworkQuota.MaxNetworks.IsUnknown() {
			networkQuotaConfig.SetMaxNetworks(plan.ConfigNetworkQuota.MaxNetworks.ValueString())
		}
		config.NetworkQuotaPolicyTypeConfiguration = networkQuotaConfig
		hasConfig = true
	}

	// ConfigPowerSchedule mapping
	if !plan.ConfigPowerSchedule.IsNull() && !plan.ConfigPowerSchedule.IsUnknown() {
		powerScheduleConfig := &sdk.PowerSchedulePolicyTypeConfiguration{}
		if !plan.ConfigPowerSchedule.PowerSchedule.IsNull() && !plan.ConfigPowerSchedule.PowerSchedule.IsUnknown() {
			powerScheduleConfig.SetPowerSchedule(plan.ConfigPowerSchedule.PowerSchedule.ValueString())
		}
		if !plan.ConfigPowerSchedule.PowerScheduleType.IsNull() && !plan.ConfigPowerSchedule.PowerScheduleType.IsUnknown() {
			powerScheduleConfig.SetPowerScheduleType(plan.ConfigPowerSchedule.PowerScheduleType.ValueString())
		}
		if !plan.ConfigPowerSchedule.PowerScheduleHideFixed.IsNull() && !plan.ConfigPowerSchedule.PowerScheduleHideFixed.IsUnknown() {
			powerScheduleConfig.SetPowerScheduleHideFixed(plan.ConfigPowerSchedule.PowerScheduleHideFixed.ValueBool())
		}
		config.PowerSchedulePolicyTypeConfiguration = powerScheduleConfig
		hasConfig = true
	}

	// ConfigRouterQuota mapping
	if !plan.ConfigRouterQuota.IsNull() && !plan.ConfigRouterQuota.IsUnknown() {
		routerQuotaConfig := &sdk.RouterQuotaPolicyTypeConfiguration{}
		if !plan.ConfigRouterQuota.MaxRouters.IsNull() && !plan.ConfigRouterQuota.MaxRouters.IsUnknown() {
			routerQuotaConfig.SetMaxRouters(plan.ConfigRouterQuota.MaxRouters.ValueString())
		}
		config.RouterQuotaPolicyTypeConfiguration = routerQuotaConfig
		hasConfig = true
	}

	// ConfigShutdown mapping
	if !plan.ConfigShutdown.IsNull() && !plan.ConfigShutdown.IsUnknown() {
		shutdownConfig := &sdk.ShutdownPolicyTypeConfiguration{}
		// ConfigShutdown appears to be a simple config without specific fields to map
		config.ShutdownPolicyTypeConfiguration = shutdownConfig
		hasConfig = true
	}

	// ConfigStorageServerStorageQuota mapping
	if !plan.ConfigStorageServerStorageQuota.IsNull() && !plan.ConfigStorageServerStorageQuota.IsUnknown() {
		storageServerQuotaConfig := &sdk.StorageServerStorageQuotaPolicyTypeConfiguration{}
		if !plan.ConfigStorageServerStorageQuota.MaxStorage.IsNull() && !plan.ConfigStorageServerStorageQuota.MaxStorage.IsUnknown() {
			storageServerQuotaConfig.SetMaxStorage(plan.ConfigStorageServerStorageQuota.MaxStorage.ValueString())
		}
		config.StorageServerStorageQuotaPolicyTypeConfiguration = storageServerQuotaConfig
		hasConfig = true
	}

	// ConfigTags mapping
	if !plan.ConfigTags.IsNull() && !plan.ConfigTags.IsUnknown() {
		tagsConfig := &sdk.TagsPolicyTypeConfiguration{}
		if !plan.ConfigTags.Key.IsNull() && !plan.ConfigTags.Key.IsUnknown() {
			tagsConfig.SetKey(plan.ConfigTags.Key.ValueString())
		}
		if !plan.ConfigTags.Value.IsNull() && !plan.ConfigTags.Value.IsUnknown() {
			tagsConfig.SetValue(plan.ConfigTags.Value.ValueString())
		}
		if !plan.ConfigTags.Strict.IsNull() && !plan.ConfigTags.Strict.IsUnknown() {
			tagsConfig.SetStrict(plan.ConfigTags.Strict.ValueBool())
		}
		if !plan.ConfigTags.ValueListId.IsNull() && !plan.ConfigTags.ValueListId.IsUnknown() {
			tagsConfig.SetValueListId(plan.ConfigTags.ValueListId.ValueString())
		}
		config.TagsPolicyTypeConfiguration = tagsConfig
		hasConfig = true
	}

	// ConfigUserCreation mapping
	if !plan.ConfigUserCreation.IsNull() && !plan.ConfigUserCreation.IsUnknown() {
		userCreationConfig := &sdk.UserCreationPolicyTypeConfiguration{}
		if !plan.ConfigUserCreation.CreateUser.IsNull() && !plan.ConfigUserCreation.CreateUser.IsUnknown() {
			userCreationConfig.SetCreateUser(plan.ConfigUserCreation.CreateUser.ValueBool())
		}
		if !plan.ConfigUserCreation.CreateUserType.IsNull() && !plan.ConfigUserCreation.CreateUserType.IsUnknown() {
			userCreationConfig.SetCreateUserType(plan.ConfigUserCreation.CreateUserType.ValueString())
		}
		config.UserCreationPolicyTypeConfiguration = userCreationConfig
		hasConfig = true
	}

	// ConfigUserGroupCreation mapping
	if !plan.ConfigUserGroupCreation.IsNull() && !plan.ConfigUserGroupCreation.IsUnknown() {
		userGroupCreationConfig := &sdk.UserGroupCreationPolicyTypeConfiguration{}
		if !plan.ConfigUserGroupCreation.UserGroup.IsNull() && !plan.ConfigUserGroupCreation.UserGroup.IsUnknown() {
			userGroupCreationConfig.SetUserGroup(plan.ConfigUserGroupCreation.UserGroup.ValueString())
		}
		config.UserGroupCreationPolicyTypeConfiguration = userGroupCreationConfig
		hasConfig = true
	}

	// ConfigWorkflow mapping
	if !plan.ConfigWorkflow.IsNull() && !plan.ConfigWorkflow.IsUnknown() {
		workflowConfig := &sdk.WorkflowPolicyTypeConfiguration{}
		if !plan.ConfigWorkflow.WorkflowId.IsNull() && !plan.ConfigWorkflow.WorkflowId.IsUnknown() {
			workflowConfig.SetWorkflowId(plan.ConfigWorkflow.WorkflowId.ValueString())
		}
		config.WorkflowPolicyTypeConfiguration = workflowConfig
		hasConfig = true
	}

	// ConfigFileShareStorageQuota mapping
	if !plan.ConfigFileShareStorageQuota.IsNull() && !plan.ConfigFileShareStorageQuota.IsUnknown() {
		fileShareQuotaConfig := &sdk.FileShareStorageQuotaPolicyTypeConfiguration{}
		if !plan.ConfigFileShareStorageQuota.MaxStorage.IsNull() && !plan.ConfigFileShareStorageQuota.MaxStorage.IsUnknown() {
			fileShareQuotaConfig.SetMaxStorage(plan.ConfigFileShareStorageQuota.MaxStorage.ValueString())
		}
		config.FileShareStorageQuotaPolicyTypeConfiguration = fileShareQuotaConfig
		hasConfig = true
	}

	if !hasConfig {
		return nil, diags
	}

	return config, diags
}

// determinePolicyTypeFromConfig automatically determines the policy type based on which config field is provided
func determinePolicyTypeFromConfig(plan *PolicyModel) string {
	// Map config fields to their corresponding policy type codes
	if !plan.ConfigApproval.IsNull() && !plan.ConfigApproval.IsUnknown() {
		return "approval"
	}
	if !plan.ConfigBackupCreation.IsNull() && !plan.ConfigBackupCreation.IsUnknown() {
		return "backupCreation"
	}
	if !plan.ConfigBackupTargets.IsNull() && !plan.ConfigBackupTargets.IsUnknown() {
		return "backupTargets"
	}
	if !plan.ConfigBudget.IsNull() && !plan.ConfigBudget.IsUnknown() {
		return "budget"
	}
	if !plan.ConfigClusterResourceName.IsNull() && !plan.ConfigClusterResourceName.IsUnknown() {
		return "clusterResourceName"
	}
	if !plan.ConfigCypherAccess.IsNull() && !plan.ConfigCypherAccess.IsUnknown() {
		return "cypherAccess"
	}
	if !plan.ConfigDelayedDelete.IsNull() && !plan.ConfigDelayedDelete.IsUnknown() {
		return "delayedDelete"
	}
	if !plan.ConfigExpiration.IsNull() && !plan.ConfigExpiration.IsUnknown() {
		return "expiration"
	}
	if !plan.ConfigFileShareStorageQuota.IsNull() && !plan.ConfigFileShareStorageQuota.IsUnknown() {
		return "fileShareStorageQuota"
	}
	if !plan.ConfigHostname.IsNull() && !plan.ConfigHostname.IsUnknown() {
		return "hostname"
	}
	if !plan.ConfigInstanceName.IsNull() && !plan.ConfigInstanceName.IsUnknown() {
		return "instanceName"
	}
	if !plan.ConfigMaxContainers.IsNull() && !plan.ConfigMaxContainers.IsUnknown() {
		return "maxContainers"
	}
	if !plan.ConfigMaxCores.IsNull() && !plan.ConfigMaxCores.IsUnknown() {
		return "maxCores"
	}
	if !plan.ConfigMaxHosts.IsNull() && !plan.ConfigMaxHosts.IsUnknown() {
		return "maxHosts"
	}
	if !plan.ConfigMaxLoadBalancerPools.IsNull() && !plan.ConfigMaxLoadBalancerPools.IsUnknown() {
		return "maxLoadBalancerPools"
	}
	if !plan.ConfigMaxMemory.IsNull() && !plan.ConfigMaxMemory.IsUnknown() {
		return "maxMemory"
	}
	if !plan.ConfigMaxPoolMembers.IsNull() && !plan.ConfigMaxPoolMembers.IsUnknown() {
		return "maxPoolMembers"
	}
	if !plan.ConfigMaxStorage.IsNull() && !plan.ConfigMaxStorage.IsUnknown() {
		return "maxStorage"
	}
	if !plan.ConfigMaxVirtualServers.IsNull() && !plan.ConfigMaxVirtualServers.IsUnknown() {
		return "maxVirtualServers"
	}
	if !plan.ConfigMaxVms.IsNull() && !plan.ConfigMaxVms.IsUnknown() {
		return "maxVms"
	}
	if !plan.ConfigMessageOfTheDay.IsNull() && !plan.ConfigMessageOfTheDay.IsUnknown() {
		return "messageOfTheDay"
	}
	if !plan.ConfigNetworkQuota.IsNull() && !plan.ConfigNetworkQuota.IsUnknown() {
		return "networkQuota"
	}
	if !plan.ConfigPowerSchedule.IsNull() && !plan.ConfigPowerSchedule.IsUnknown() {
		return "powerSchedule"
	}
	if !plan.ConfigRouterQuota.IsNull() && !plan.ConfigRouterQuota.IsUnknown() {
		return "routerQuota"
	}
	if !plan.ConfigShutdown.IsNull() && !plan.ConfigShutdown.IsUnknown() {
		return "shutdown"
	}
	if !plan.ConfigStorageServerStorageQuota.IsNull() && !plan.ConfigStorageServerStorageQuota.IsUnknown() {
		return "storageServerStorageQuota"
	}
	if !plan.ConfigTags.IsNull() && !plan.ConfigTags.IsUnknown() {
		return "tags"
	}
	if !plan.ConfigUserCreation.IsNull() && !plan.ConfigUserCreation.IsUnknown() {
		return "userCreation"
	}
	if !plan.ConfigUserGroupCreation.IsNull() && !plan.ConfigUserGroupCreation.IsUnknown() {
		return "userGroupCreation"
	}
	if !plan.ConfigWorkflow.IsNull() && !plan.ConfigWorkflow.IsUnknown() {
		return "workflow"
	}

	// Return empty string if no config is provided - this will cause an error which is appropriate
	return ""
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan PolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	addPolicy := sdk.NewAddPoliciesRequestPolicyWithDefaults()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create policy resource",
			"policy "+name+": failed to create client: "+err.Error(),
		)
		return
	}

	// Set required fields
	addPolicy.SetName(name)

	// Set optional fields
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		addPolicy.SetDescription(plan.Description.ValueString())
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		addPolicy.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.EachUser.IsNull() && !plan.EachUser.IsUnknown() {
		addPolicy.SetEachUser(plan.EachUser.ValueBool())
	}

	if !plan.RefId.IsNull() && !plan.RefId.IsUnknown() {
		addPolicy.SetRefId(plan.RefId.ValueInt64())
	}

	// Set account IDs if provided
	if !plan.Accounts.IsNull() && !plan.Accounts.IsUnknown() {
		var accountIDs []int64
		diags := plan.Accounts.ElementsAs(ctx, &accountIDs, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		addPolicy.SetAccounts(accountIDs)
	}

	// Set PolicyType - either from plan or automatically determined from config
	policyTypeCode := ""
	if !plan.PolicyType.IsNull() && !plan.PolicyType.IsUnknown() && !plan.PolicyType.Code.IsNull() && !plan.PolicyType.Code.IsUnknown() {
		policyTypeCode = plan.PolicyType.Code.ValueString()
	} else {
		// Automatically determine policy type based on which config field is provided
		policyTypeCode = determinePolicyTypeFromConfig(&plan)
	}

	if policyTypeCode != "" {
		policyType := sdk.NewAddPoliciesRequestPolicyPolicyTypeWithDefaults()
		policyType.SetCode(policyTypeCode)
		addPolicy.SetPolicyType(*policyType)
	}

	// Set Config based on the schema fields
	config, configDiags := buildPolicyConfigForCreate(ctx, &plan)
	if configDiags.HasError() {
		resp.Diagnostics.Append(configDiags...)
		return
	}
	if config != nil {
		addPolicy.SetConfig(*config)
	}

	// Set RefType if provided
	if !plan.RefType.IsNull() && !plan.RefType.IsUnknown() {
		if !plan.RefType.Oneof0.IsNull() && !plan.RefType.Oneof0.IsUnknown() {
			addPolicy.SetRefType(plan.RefType.Oneof0.ValueString())
		}
	}

	addPolicyRequest := sdk.NewAddPoliciesRequest(*addPolicy)

	policy, hresp, err := client.PoliciesAPI.AddPolicies(ctx).AddPoliciesRequest(*addPolicyRequest).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create policy resource",
			"policy "+name+" POST failed: "+errors.ErrMsg(err, hresp),
		)
		return
	}

	if policy.Policy == nil || policy.Policy.Id == nil {
		resp.Diagnostics.AddError(
			"create policy resource",
			"policy "+name+" id is nil",
		)
		return
	}

	id := *policy.Policy.Id
	plan.Id = types.Int64Value(id)

	// Write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the created policy to get full state
	state, diags := getPolicyAsState(ctx, id, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// buildPolicyConfigForUpdate maps the schema config fields to the SDK config structure for update operations
func buildPolicyConfigForUpdate(ctx context.Context, plan *PolicyModel) (*sdk.UpdatePoliciesRequestPolicyConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	config := &sdk.UpdatePoliciesRequestPolicyConfig{}
	hasConfig := false

	// ConfigApproval mapping
	if !plan.ConfigApproval.IsNull() && !plan.ConfigApproval.IsUnknown() {
		approveConfig := &sdk.ApprovePolicyTypeConfiguration{}
		if !plan.ConfigApproval.AccountIntegrationId.IsNull() && !plan.ConfigApproval.AccountIntegrationId.IsUnknown() {
			approveConfig.SetAccountIntegrationId(plan.ConfigApproval.AccountIntegrationId.ValueString())
		}
		config.ApprovePolicyTypeConfiguration = approveConfig
		hasConfig = true
	}

	// ConfigBackupCreation mapping
	if !plan.ConfigBackupCreation.IsNull() && !plan.ConfigBackupCreation.IsUnknown() {
		backupConfig := &sdk.BackupCreationPolicyTypeConfiguration{}
		if !plan.ConfigBackupCreation.CreateBackup.IsNull() && !plan.ConfigBackupCreation.CreateBackup.IsUnknown() {
			backupConfig.SetCreateBackup(plan.ConfigBackupCreation.CreateBackup.ValueBool())
		}
		if !plan.ConfigBackupCreation.CreateBackupType.IsNull() && !plan.ConfigBackupCreation.CreateBackupType.IsUnknown() {
			backupConfig.SetCreateBackupType(plan.ConfigBackupCreation.CreateBackupType.ValueString())
		}
		config.BackupCreationPolicyTypeConfiguration = backupConfig
		hasConfig = true
	}

	// ConfigBackupTargets mapping
	if !plan.ConfigBackupTargets.IsNull() && !plan.ConfigBackupTargets.IsUnknown() {
		backupTargetsConfig := &sdk.BackupTargetsPolicyTypeConfiguration{}
		if !plan.ConfigBackupTargets.BackupStorageIds.IsNull() && !plan.ConfigBackupTargets.BackupStorageIds.IsUnknown() {
			var storageIds []int64
			elemDiags := plan.ConfigBackupTargets.BackupStorageIds.ElementsAs(ctx, &storageIds, false)
			if elemDiags.HasError() {
				diags.Append(elemDiags...)
			} else {
				backupTargetsConfig.SetBackupStorageIds(storageIds)
			}
		}
		config.BackupTargetsPolicyTypeConfiguration = backupTargetsConfig
		hasConfig = true
	}

	// ConfigBudget mapping
	if !plan.ConfigBudget.IsNull() && !plan.ConfigBudget.IsUnknown() {
		budgetConfig := &sdk.BudgetPolicyTypeConfiguration{}
		if !plan.ConfigBudget.MaxPrice.IsNull() && !plan.ConfigBudget.MaxPrice.IsUnknown() {
			maxPrice, _ := plan.ConfigBudget.MaxPrice.ValueBigFloat().Float64()
			budgetConfig.SetMaxPrice(float32(maxPrice))
		}
		if !plan.ConfigBudget.MaxPriceCurrency.IsNull() && !plan.ConfigBudget.MaxPriceCurrency.IsUnknown() {
			budgetConfig.SetMaxPriceCurrency(plan.ConfigBudget.MaxPriceCurrency.ValueString())
		}
		if !plan.ConfigBudget.MaxPriceUnit.IsNull() && !plan.ConfigBudget.MaxPriceUnit.IsUnknown() {
			budgetConfig.SetMaxPriceUnit(plan.ConfigBudget.MaxPriceUnit.ValueString())
		}
		config.BudgetPolicyTypeConfiguration = budgetConfig
		hasConfig = true
	}

	// ConfigClusterResourceName mapping
	if !plan.ConfigClusterResourceName.IsNull() && !plan.ConfigClusterResourceName.IsUnknown() {
		clusterConfig := &sdk.ClusterResourceNamePolicyTypeConfiguration{}
		if !plan.ConfigClusterResourceName.ServerNamingType.IsNull() && !plan.ConfigClusterResourceName.ServerNamingType.IsUnknown() {
			clusterConfig.SetServerNamingType(plan.ConfigClusterResourceName.ServerNamingType.ValueString())
		}
		if !plan.ConfigClusterResourceName.ServerNamingPattern.IsNull() && !plan.ConfigClusterResourceName.ServerNamingPattern.IsUnknown() {
			clusterConfig.SetServerNamingPattern(plan.ConfigClusterResourceName.ServerNamingPattern.ValueString())
		}
		if !plan.ConfigClusterResourceName.ServerNamingConflict.IsNull() && !plan.ConfigClusterResourceName.ServerNamingConflict.IsUnknown() {
			clusterConfig.SetServerNamingConflict(plan.ConfigClusterResourceName.ServerNamingConflict.ValueBool())
		}
		config.ClusterResourceNamePolicyTypeConfiguration = clusterConfig
		hasConfig = true
	}

	// ConfigCypherAccess mapping
	if !plan.ConfigCypherAccess.IsNull() && !plan.ConfigCypherAccess.IsUnknown() {
		cypherConfig := &sdk.CypherAccessPolicyTypeConfiguration{}
		if !plan.ConfigCypherAccess.KeyPattern.IsNull() && !plan.ConfigCypherAccess.KeyPattern.IsUnknown() {
			cypherConfig.SetKeyPattern(plan.ConfigCypherAccess.KeyPattern.ValueString())
		}
		if !plan.ConfigCypherAccess.Read.IsNull() && !plan.ConfigCypherAccess.Read.IsUnknown() {
			cypherConfig.SetRead(plan.ConfigCypherAccess.Read.ValueBool())
		}
		if !plan.ConfigCypherAccess.Write.IsNull() && !plan.ConfigCypherAccess.Write.IsUnknown() {
			cypherConfig.SetWrite(plan.ConfigCypherAccess.Write.ValueBool())
		}
		if !plan.ConfigCypherAccess.Update.IsNull() && !plan.ConfigCypherAccess.Update.IsUnknown() {
			cypherConfig.SetUpdate(plan.ConfigCypherAccess.Update.ValueBool())
		}
		if !plan.ConfigCypherAccess.Delete.IsNull() && !plan.ConfigCypherAccess.Delete.IsUnknown() {
			cypherConfig.SetDelete(plan.ConfigCypherAccess.Delete.ValueBool())
		}
		if !plan.ConfigCypherAccess.List.IsNull() && !plan.ConfigCypherAccess.List.IsUnknown() {
			cypherConfig.SetList(plan.ConfigCypherAccess.List.ValueBool())
		}
		config.CypherAccessPolicyTypeConfiguration = cypherConfig
		hasConfig = true
	}

	// ConfigDelayedDelete mapping
	if !plan.ConfigDelayedDelete.IsNull() && !plan.ConfigDelayedDelete.IsUnknown() {
		delayedDeleteConfig := &sdk.DelayedDeletePolicyTypeConfiguration{}
		if !plan.ConfigDelayedDelete.RemovalAge.IsNull() && !plan.ConfigDelayedDelete.RemovalAge.IsUnknown() {
			delayedDeleteConfig.SetRemovalAge(plan.ConfigDelayedDelete.RemovalAge.ValueString())
		}
		config.DelayedDeletePolicyTypeConfiguration = delayedDeleteConfig
		hasConfig = true
	}

	// Continue with more config mappings...

	// ConfigExpiration mapping
	if !plan.ConfigExpiration.IsNull() && !plan.ConfigExpiration.IsUnknown() {
		expirationConfig := &sdk.ExpirationPolicyTypeConfiguration{}
		if !plan.ConfigExpiration.LifecycleType.IsNull() && !plan.ConfigExpiration.LifecycleType.IsUnknown() {
			expirationConfig.SetLifecycleType(plan.ConfigExpiration.LifecycleType.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleAge.IsNull() && !plan.ConfigExpiration.LifecycleAge.IsUnknown() {
			expirationConfig.SetLifecycleAge(plan.ConfigExpiration.LifecycleAge.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleAllowExtend.IsNull() && !plan.ConfigExpiration.LifecycleAllowExtend.IsUnknown() {
			expirationConfig.SetLifecycleAllowExtend(plan.ConfigExpiration.LifecycleAllowExtend.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleAutoRenew.IsNull() && !plan.ConfigExpiration.LifecycleAutoRenew.IsUnknown() {
			expirationConfig.SetLifecycleAutoRenew(plan.ConfigExpiration.LifecycleAutoRenew.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleExtensionsBeforeApproval.IsNull() && !plan.ConfigExpiration.LifecycleExtensionsBeforeApproval.IsUnknown() {
			expirationConfig.SetLifecycleExtensionsBeforeApproval(plan.ConfigExpiration.LifecycleExtensionsBeforeApproval.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleHideFixed.IsNull() && !plan.ConfigExpiration.LifecycleHideFixed.IsUnknown() {
			expirationConfig.SetLifecycleHideFixed(plan.ConfigExpiration.LifecycleHideFixed.ValueBool())
		}
		if !plan.ConfigExpiration.LifecycleMessage.IsNull() && !plan.ConfigExpiration.LifecycleMessage.IsUnknown() {
			expirationConfig.SetLifecycleMessage(plan.ConfigExpiration.LifecycleMessage.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleNotify.IsNull() && !plan.ConfigExpiration.LifecycleNotify.IsUnknown() {
			expirationConfig.SetLifecycleNotify(plan.ConfigExpiration.LifecycleNotify.ValueString())
		}
		if !plan.ConfigExpiration.LifecycleRenewal.IsNull() && !plan.ConfigExpiration.LifecycleRenewal.IsUnknown() {
			expirationConfig.SetLifecycleRenewal(plan.ConfigExpiration.LifecycleRenewal.ValueString())
		}
		if !plan.ConfigExpiration.AccountIntegrationId.IsNull() && !plan.ConfigExpiration.AccountIntegrationId.IsUnknown() {
			expirationConfig.SetAccountIntegrationId(plan.ConfigExpiration.AccountIntegrationId.ValueString())
		}
		config.ExpirationPolicyTypeConfiguration = expirationConfig
		hasConfig = true
	}

	// ConfigHostname mapping
	if !plan.ConfigHostname.IsNull() && !plan.ConfigHostname.IsUnknown() {
		hostnameConfig := &sdk.HostnamePolicyTypeConfiguration{}
		if !plan.ConfigHostname.HostNamingType.IsNull() && !plan.ConfigHostname.HostNamingType.IsUnknown() {
			hostnameConfig.SetHostNamingType(plan.ConfigHostname.HostNamingType.ValueString())
		}
		if !plan.ConfigHostname.HostNamingPattern.IsNull() && !plan.ConfigHostname.HostNamingPattern.IsUnknown() {
			hostnameConfig.SetHostNamingPattern(plan.ConfigHostname.HostNamingPattern.ValueString())
		}
		config.HostnamePolicyTypeConfiguration = hostnameConfig
		hasConfig = true
	}

	// ConfigInstanceName mapping
	if !plan.ConfigInstanceName.IsNull() && !plan.ConfigInstanceName.IsUnknown() {
		instanceNameConfig := &sdk.InstanceNamePolicyTypeConfiguration{}
		if !plan.ConfigInstanceName.NamingType.IsNull() && !plan.ConfigInstanceName.NamingType.IsUnknown() {
			instanceNameConfig.SetNamingType(plan.ConfigInstanceName.NamingType.ValueString())
		}
		if !plan.ConfigInstanceName.NamingPattern.IsNull() && !plan.ConfigInstanceName.NamingPattern.IsUnknown() {
			instanceNameConfig.SetNamingPattern(plan.ConfigInstanceName.NamingPattern.ValueString())
		}
		if !plan.ConfigInstanceName.NamingConflict.IsNull() && !plan.ConfigInstanceName.NamingConflict.IsUnknown() {
			instanceNameConfig.SetNamingConflict(plan.ConfigInstanceName.NamingConflict.ValueBool())
		}
		config.InstanceNamePolicyTypeConfiguration = instanceNameConfig
		hasConfig = true
	}

	// ConfigMaxContainers mapping
	if !plan.ConfigMaxContainers.IsNull() && !plan.ConfigMaxContainers.IsUnknown() {
		maxContainersConfig := &sdk.MaxContainersPolicyTypeConfiguration{}
		if !plan.ConfigMaxContainers.MaxContainers.IsNull() && !plan.ConfigMaxContainers.MaxContainers.IsUnknown() {
			maxContainersConfig.SetMaxContainers(plan.ConfigMaxContainers.MaxContainers.ValueString())
		}
		config.MaxContainersPolicyTypeConfiguration = maxContainersConfig
		hasConfig = true
	}

	// ConfigMaxCores mapping
	if !plan.ConfigMaxCores.IsNull() && !plan.ConfigMaxCores.IsUnknown() {
		maxCoresConfig := &sdk.MaxCoresPolicyTypeConfiguration1{}
		if !plan.ConfigMaxCores.MaxCores.IsNull() && !plan.ConfigMaxCores.MaxCores.IsUnknown() {
			maxCoresConfig.SetMaxCores(plan.ConfigMaxCores.MaxCores.ValueString())
		}
		if !plan.ConfigMaxCores.ExcludeContainers.IsNull() && !plan.ConfigMaxCores.ExcludeContainers.IsUnknown() {
			// ConfigMaxCores.ExcludeContainers in update expects bool, but schema has string - need to convert
			excludeStr := plan.ConfigMaxCores.ExcludeContainers.ValueString()
			excludeBool := excludeStr == "true" || excludeStr == "on" || excludeStr == "1"
			maxCoresConfig.SetExcludeContainers(excludeBool)
		}
		config.MaxCoresPolicyTypeConfiguration1 = maxCoresConfig
		hasConfig = true
	}

	// ConfigMaxHosts mapping
	if !plan.ConfigMaxHosts.IsNull() && !plan.ConfigMaxHosts.IsUnknown() {
		maxHostsConfig := &sdk.MaxHostsPolicyTypeConfiguration{}
		if !plan.ConfigMaxHosts.MaxHosts.IsNull() && !plan.ConfigMaxHosts.MaxHosts.IsUnknown() {
			maxHostsConfig.SetMaxHosts(plan.ConfigMaxHosts.MaxHosts.ValueString())
		}
		config.MaxHostsPolicyTypeConfiguration = maxHostsConfig
		hasConfig = true
	}

	// ConfigMaxLoadBalancerPools mapping
	if !plan.ConfigMaxLoadBalancerPools.IsNull() && !plan.ConfigMaxLoadBalancerPools.IsUnknown() {
		maxLBPoolsConfig := &sdk.MaxLoadBalancerPoolsPolicyTypeConfiguration{}
		if !plan.ConfigMaxLoadBalancerPools.MaxPools.IsNull() && !plan.ConfigMaxLoadBalancerPools.MaxPools.IsUnknown() {
			maxLBPoolsConfig.SetMaxPools(plan.ConfigMaxLoadBalancerPools.MaxPools.ValueString())
		}
		config.MaxLoadBalancerPoolsPolicyTypeConfiguration = maxLBPoolsConfig
		hasConfig = true
	}

	// ConfigMaxMemory mapping (complex nested structure)
	if !plan.ConfigMaxMemory.IsNull() && !plan.ConfigMaxMemory.IsUnknown() {
		maxMemoryConfig := &sdk.MaxMemoryPolicyTypeConfiguration1{}
		maxMemorySet := false

		if !plan.ConfigMaxMemory.ExcludeContainers.IsNull() && !plan.ConfigMaxMemory.ExcludeContainers.IsUnknown() {
			maxMemoryConfig.SetExcludeContainers(plan.ConfigMaxMemory.ExcludeContainers.ValueString())
		}

		// Handle MaxMemory field - it's an ObjectValue with anyof0 (string) and anyof1 (int)
		if !plan.ConfigMaxMemory.MaxMemory.IsNull() && !plan.ConfigMaxMemory.MaxMemory.IsUnknown() {
			maxMemoryValue := plan.ConfigMaxMemory.MaxMemory
			// Extract the MaxMemoryValue from the ObjectValue
			if maxMemoryAttrs := maxMemoryValue.Attributes(); len(maxMemoryAttrs) > 0 {
				var maxMemory sdk.MaxMemoryPolicyTypeConfigurationMaxMemory
				// Check anyof0 (string) first
				if anyof0Val, ok := maxMemoryAttrs["anyof0"]; ok {
					if stringVal, ok := anyof0Val.(basetypes.StringValue); ok && !stringVal.IsNull() && !stringVal.IsUnknown() {
						// Use string representation (anyof0)
						str := stringVal.ValueString()
						maxMemory.String = &str
						maxMemoryConfig.SetMaxMemory(maxMemory)
						maxMemorySet = true
					}
				}
				// Check anyof1 (int64) if anyof0 wasn't set
				if !maxMemorySet {
					if anyof1Val, ok := maxMemoryAttrs["anyof1"]; ok {
						if intVal, ok := anyof1Val.(basetypes.Int64Value); ok && !intVal.IsNull() && !intVal.IsUnknown() {
							// Use int representation (anyof1)
							intValue := intVal.ValueInt64()
							maxMemory.Int64 = &intValue
							maxMemoryConfig.SetMaxMemory(maxMemory)
							maxMemorySet = true
						}
					}
				}
			}
		}

		// Only add the config if MaxMemory was actually set (it's required)
		if maxMemorySet {
			config.MaxMemoryPolicyTypeConfiguration1 = maxMemoryConfig
			hasConfig = true
		}
	}

	// ConfigMaxPoolMembers mapping
	if !plan.ConfigMaxPoolMembers.IsNull() && !plan.ConfigMaxPoolMembers.IsUnknown() {
		maxPoolMembersConfig := &sdk.MaxPoolMembersPolicyTypeConfiguration{}
		if !plan.ConfigMaxPoolMembers.MaxPoolMembers.IsNull() && !plan.ConfigMaxPoolMembers.MaxPoolMembers.IsUnknown() {
			maxPoolMembersConfig.SetMaxPoolMembers(plan.ConfigMaxPoolMembers.MaxPoolMembers.ValueString())
		}
		config.MaxPoolMembersPolicyTypeConfiguration = maxPoolMembersConfig
		hasConfig = true
	}

	// ConfigMaxStorage mapping
	if !plan.ConfigMaxStorage.IsNull() && !plan.ConfigMaxStorage.IsUnknown() {
		maxStorageConfig := &sdk.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration1{}
		if !plan.ConfigMaxStorage.MaxStorage.IsNull() && !plan.ConfigMaxStorage.MaxStorage.IsUnknown() {
			maxStorageConfig.SetMaxStorage(plan.ConfigMaxStorage.MaxStorage.ValueString())
		}
		if !plan.ConfigMaxStorage.ExcludeContainers.IsNull() && !plan.ConfigMaxStorage.ExcludeContainers.IsUnknown() {
			// ConfigMaxStorage.ExcludeContainers in update expects bool, but schema has string - need to convert
			excludeStr := plan.ConfigMaxStorage.ExcludeContainers.ValueString()
			excludeBool := excludeStr == "true" || excludeStr == "on" || excludeStr == "1"
			maxStorageConfig.SetExcludeContainers(excludeBool)
		}
		config.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration1 = maxStorageConfig
		hasConfig = true
	}

	// ConfigMaxVirtualServers mapping
	if !plan.ConfigMaxVirtualServers.IsNull() && !plan.ConfigMaxVirtualServers.IsUnknown() {
		maxVirtualServersConfig := &sdk.MaxVirtualServersPolicyTypeConfiguration{}
		if !plan.ConfigMaxVirtualServers.MaxVirtualServers.IsNull() && !plan.ConfigMaxVirtualServers.MaxVirtualServers.IsUnknown() {
			maxVirtualServersConfig.SetMaxVirtualServers(plan.ConfigMaxVirtualServers.MaxVirtualServers.ValueString())
		}
		config.MaxVirtualServersPolicyTypeConfiguration = maxVirtualServersConfig
		hasConfig = true
	}

	// ConfigMaxVms mapping
	if !plan.ConfigMaxVms.IsNull() && !plan.ConfigMaxVms.IsUnknown() {
		maxVmsConfig := &sdk.MaxVMsPolicyTypeConfiguration{}
		if !plan.ConfigMaxVms.MaxVms.IsNull() && !plan.ConfigMaxVms.MaxVms.IsUnknown() {
			maxVmsConfig.SetMaxVms(plan.ConfigMaxVms.MaxVms.ValueString())
		}
		config.MaxVMsPolicyTypeConfiguration = maxVmsConfig
		hasConfig = true
	}

	// ConfigMessageOfTheDay mapping
	if !plan.ConfigMessageOfTheDay.IsNull() && !plan.ConfigMessageOfTheDay.IsUnknown() {
		motdConfig := &sdk.MessageOfTheDayPolicyTypeConfiguration{}
		if !plan.ConfigMessageOfTheDay.Motddate.IsNull() && !plan.ConfigMessageOfTheDay.Motddate.IsUnknown() {
			// Parse the date string - assuming RFC3339 format
			if dateStr := plan.ConfigMessageOfTheDay.Motddate.ValueString(); dateStr != "" {
				if parsedDate, err := time.Parse(time.RFC3339, dateStr); err == nil {
					motdConfig.SetMotdDate(parsedDate)
				}
				// If parsing fails, we could add to diagnostics, but for now continue without setting
			}
		}
		if !plan.ConfigMessageOfTheDay.Motdmessage.IsNull() && !plan.ConfigMessageOfTheDay.Motdmessage.IsUnknown() {
			motdConfig.SetMotdMessage(plan.ConfigMessageOfTheDay.Motdmessage.ValueString())
		}
		if !plan.ConfigMessageOfTheDay.Motdtitle.IsNull() && !plan.ConfigMessageOfTheDay.Motdtitle.IsUnknown() {
			motdConfig.SetMotdTitle(plan.ConfigMessageOfTheDay.Motdtitle.ValueString())
		}
		if !plan.ConfigMessageOfTheDay.Motdtype.IsNull() && !plan.ConfigMessageOfTheDay.Motdtype.IsUnknown() {
			motdConfig.SetMotdType(plan.ConfigMessageOfTheDay.Motdtype.ValueString())
		}
		if !plan.ConfigMessageOfTheDay.MotdFullPage.IsNull() && !plan.ConfigMessageOfTheDay.MotdFullPage.IsUnknown() {
			motdConfig.SetMotdFullPage(plan.ConfigMessageOfTheDay.MotdFullPage.ValueBool())
		}
		config.MessageOfTheDayPolicyTypeConfiguration = motdConfig
		hasConfig = true
	}

	// ConfigNetworkQuota mapping
	if !plan.ConfigNetworkQuota.IsNull() && !plan.ConfigNetworkQuota.IsUnknown() {
		networkQuotaConfig := &sdk.NetworkQuotaPolicyTypeConfiguration{}
		if !plan.ConfigNetworkQuota.MaxNetworks.IsNull() && !plan.ConfigNetworkQuota.MaxNetworks.IsUnknown() {
			networkQuotaConfig.SetMaxNetworks(plan.ConfigNetworkQuota.MaxNetworks.ValueString())
		}
		config.NetworkQuotaPolicyTypeConfiguration = networkQuotaConfig
		hasConfig = true
	}

	// ConfigPowerSchedule mapping
	if !plan.ConfigPowerSchedule.IsNull() && !plan.ConfigPowerSchedule.IsUnknown() {
		powerScheduleConfig := &sdk.PowerSchedulePolicyTypeConfiguration{}
		if !plan.ConfigPowerSchedule.PowerSchedule.IsNull() && !plan.ConfigPowerSchedule.PowerSchedule.IsUnknown() {
			powerScheduleConfig.SetPowerSchedule(plan.ConfigPowerSchedule.PowerSchedule.ValueString())
		}
		if !plan.ConfigPowerSchedule.PowerScheduleType.IsNull() && !plan.ConfigPowerSchedule.PowerScheduleType.IsUnknown() {
			powerScheduleConfig.SetPowerScheduleType(plan.ConfigPowerSchedule.PowerScheduleType.ValueString())
		}
		if !plan.ConfigPowerSchedule.PowerScheduleHideFixed.IsNull() && !plan.ConfigPowerSchedule.PowerScheduleHideFixed.IsUnknown() {
			powerScheduleConfig.SetPowerScheduleHideFixed(plan.ConfigPowerSchedule.PowerScheduleHideFixed.ValueBool())
		}
		config.PowerSchedulePolicyTypeConfiguration = powerScheduleConfig
		hasConfig = true
	}

	// ConfigRouterQuota mapping
	if !plan.ConfigRouterQuota.IsNull() && !plan.ConfigRouterQuota.IsUnknown() {
		routerQuotaConfig := &sdk.RouterQuotaPolicyTypeConfiguration{}
		if !plan.ConfigRouterQuota.MaxRouters.IsNull() && !plan.ConfigRouterQuota.MaxRouters.IsUnknown() {
			routerQuotaConfig.SetMaxRouters(plan.ConfigRouterQuota.MaxRouters.ValueString())
		}
		config.RouterQuotaPolicyTypeConfiguration = routerQuotaConfig
		hasConfig = true
	}

	// ConfigShutdown mapping
	if !plan.ConfigShutdown.IsNull() && !plan.ConfigShutdown.IsUnknown() {
		shutdownConfig := &sdk.ShutdownPolicyTypeConfiguration{}
		config.ShutdownPolicyTypeConfiguration = shutdownConfig
		hasConfig = true
	}

	// ConfigStorageServerStorageQuota mapping
	if !plan.ConfigStorageServerStorageQuota.IsNull() && !plan.ConfigStorageServerStorageQuota.IsUnknown() {
		storageServerQuotaConfig := &sdk.StorageServerStorageQuotaPolicyTypeConfiguration{}
		config.StorageServerStorageQuotaPolicyTypeConfiguration = storageServerQuotaConfig
		hasConfig = true
	}

	// ConfigTags mapping
	if !plan.ConfigTags.IsNull() && !plan.ConfigTags.IsUnknown() {
		tagsConfig := &sdk.TagsPolicyTypeConfiguration{}
		if !plan.ConfigTags.Key.IsNull() && !plan.ConfigTags.Key.IsUnknown() {
			tagsConfig.SetKey(plan.ConfigTags.Key.ValueString())
		}
		if !plan.ConfigTags.Value.IsNull() && !plan.ConfigTags.Value.IsUnknown() {
			tagsConfig.SetValue(plan.ConfigTags.Value.ValueString())
		}
		if !plan.ConfigTags.Strict.IsNull() && !plan.ConfigTags.Strict.IsUnknown() {
			tagsConfig.SetStrict(plan.ConfigTags.Strict.ValueBool())
		}
		if !plan.ConfigTags.ValueListId.IsNull() && !plan.ConfigTags.ValueListId.IsUnknown() {
			tagsConfig.SetValueListId(plan.ConfigTags.ValueListId.ValueString())
		}
		config.TagsPolicyTypeConfiguration = tagsConfig
		hasConfig = true
	}

	// ConfigUserCreation mapping
	if !plan.ConfigUserCreation.IsNull() && !plan.ConfigUserCreation.IsUnknown() {
		userCreationConfig := &sdk.UserCreationPolicyTypeConfiguration{}
		if !plan.ConfigUserCreation.CreateUser.IsNull() && !plan.ConfigUserCreation.CreateUser.IsUnknown() {
			userCreationConfig.SetCreateUser(plan.ConfigUserCreation.CreateUser.ValueBool())
		}
		if !plan.ConfigUserCreation.CreateUserType.IsNull() && !plan.ConfigUserCreation.CreateUserType.IsUnknown() {
			userCreationConfig.SetCreateUserType(plan.ConfigUserCreation.CreateUserType.ValueString())
		}
		config.UserCreationPolicyTypeConfiguration = userCreationConfig
		hasConfig = true
	}

	// ConfigUserGroupCreation mapping
	if !plan.ConfigUserGroupCreation.IsNull() && !plan.ConfigUserGroupCreation.IsUnknown() {
		userGroupCreationConfig := &sdk.UserGroupCreationPolicyTypeConfiguration{}
		if !plan.ConfigUserGroupCreation.UserGroup.IsNull() && !plan.ConfigUserGroupCreation.UserGroup.IsUnknown() {
			userGroupCreationConfig.SetUserGroup(plan.ConfigUserGroupCreation.UserGroup.ValueString())
		}
		config.UserGroupCreationPolicyTypeConfiguration = userGroupCreationConfig
		hasConfig = true
	}

	// ConfigWorkflow mapping
	if !plan.ConfigWorkflow.IsNull() && !plan.ConfigWorkflow.IsUnknown() {
		workflowConfig := &sdk.WorkflowPolicyTypeConfiguration{}
		if !plan.ConfigWorkflow.WorkflowId.IsNull() && !plan.ConfigWorkflow.WorkflowId.IsUnknown() {
			workflowConfig.SetWorkflowId(plan.ConfigWorkflow.WorkflowId.ValueString())
		}
		config.WorkflowPolicyTypeConfiguration = workflowConfig
		hasConfig = true
	}

	// ConfigFileShareStorageQuota mapping
	if !plan.ConfigFileShareStorageQuota.IsNull() && !plan.ConfigFileShareStorageQuota.IsUnknown() {
		fileShareQuotaConfig := &sdk.FileShareStorageQuotaPolicyTypeConfiguration{}
		if !plan.ConfigFileShareStorageQuota.MaxStorage.IsNull() && !plan.ConfigFileShareStorageQuota.MaxStorage.IsUnknown() {
			fileShareQuotaConfig.SetMaxStorage(plan.ConfigFileShareStorageQuota.MaxStorage.ValueString())
		}
		config.FileShareStorageQuotaPolicyTypeConfiguration = fileShareQuotaConfig
		hasConfig = true
	}

	if !hasConfig {
		return nil, diags
	}

	return config, diags
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
	state, diags := getPolicyAsState(ctx, id, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state PolicyModel

	// Get prior state (has the ID) before touching any Value* methods.
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()
	updatePolicy := sdk.NewUpdatePoliciesRequestPolicyWithDefaults()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update policy resource",
			"failed to create client: "+err.Error(),
		)
		return
	}

	// Set updateable fields from plan
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updatePolicy.SetName(plan.Name.ValueString())
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updatePolicy.SetDescription(plan.Description.ValueString())
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		updatePolicy.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.EachUser.IsNull() && !plan.EachUser.IsUnknown() {
		updatePolicy.SetEachUser(plan.EachUser.ValueBool())
	}

	if !plan.RefId.IsNull() && !plan.RefId.IsUnknown() {
		updatePolicy.SetRefId(plan.RefId.ValueInt64())
	}

	// Set account IDs if provided
	if !plan.Accounts.IsNull() && !plan.Accounts.IsUnknown() {
		var accountIDs []int64
		diags := plan.Accounts.ElementsAs(ctx, &accountIDs, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		updatePolicy.SetAccounts(accountIDs)
	}

	// Set PolicyType if provided (not supported in update - typically policies don't change type)
	// Commenting out for now since the SDK doesn't support PolicyType updates
	// if !plan.PolicyType.IsNull() && !plan.PolicyType.IsUnknown() {
	//	 updatePolicy.SetPolicyType(...)
	// }

	// Set Config based on the schema fields for update
	config, configDiags := buildPolicyConfigForUpdate(ctx, &plan)
	if configDiags.HasError() {
		resp.Diagnostics.Append(configDiags...)
		return
	}
	if config != nil {
		updatePolicy.SetConfig(*config)
	}

	// Set RefType if provided
	if !plan.RefType.IsNull() && !plan.RefType.IsUnknown() {
		if !plan.RefType.Oneof0.IsNull() && !plan.RefType.Oneof0.IsUnknown() {
			updatePolicy.SetRefType(plan.RefType.Oneof0.ValueString())
		}
	}

	updatePolicyRequest := sdk.NewUpdatePoliciesRequest(*updatePolicy)

	_, hresp, err := client.PoliciesAPI.UpdatePolicies(ctx, id).UpdatePoliciesRequest(*updatePolicyRequest).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update policy resource",
			fmt.Sprintf("policy %d UPDATE failed: %s",
				id, errors.ErrMsg(err, hresp)),
		)
		return
	}

	updatedState, diags := getPolicyAsState(ctx, id, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"update policy resource",
			fmt.Sprintf("policy %d: failed to read from api", id),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data PolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete policy resource",
			fmt.Sprintf("policy %d: failed to create client: %s", id, err.Error()),
		)
		return
	}

	_, hresp, err := client.PoliciesAPI.RemovePolicies(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete policy resource",
			fmt.Sprintf("policy %d: DELETE failed ", id)+errors.ErrMsg(err, hresp),
		)
		return
	}
}
