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
		refIdStr := *p.RefId.Get()
		// Try to parse as int64
		var refIdInt int64
		fmt.Sscanf(refIdStr, "%d", &refIdInt)
		state.RefId = types.Int64Value(refIdInt)
	}

	// Handle RefType
	if p.RefType.IsSet() && p.RefType.Get() != nil {
		refTypeStr := *p.RefType.Get()
		refTypeAttrs := map[string]attr.Value{
			"oneof0": types.StringValue(refTypeStr),
		}
		refTypeValue, refTypeDiags := NewRefTypeValue(
			RefTypeValue{}.AttributeTypes(ctx), refTypeAttrs)
		if refTypeDiags.HasError() {
			diags.Append(refTypeDiags...)
			return state, diags
		}
		state.RefType = refTypeValue
	}

	// Set account IDs
	if p.Accounts != nil && len(p.Accounts) > 0 {
		accountIDs := make([]int64, 0, len(p.Accounts))
		for _, acc := range p.Accounts {
			if acc.Id != nil {
				accountIDs = append(accountIDs, *acc.Id)
			}
		}
		accountSet, setDiags := types.SetValueFrom(ctx, types.Int64Type, accountIDs)
		if setDiags.HasError() {
			diags.Append(setDiags...)
			return state, diags
		}
		state.Accounts = accountSet
	} else {
		// Set to null if no accounts
		state.Accounts = types.SetNull(types.Int64Type)
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

	// Preserve config from plan - the API doesn't always return the full config structure
	// so we trust the plan's config since that's what was applied
	if plan != nil && !plan.Config.IsNull() {
		state.Config = plan.Config
	} else {
		// Initialize Config with null values for all configuration types if no plan available
		configAttrs := map[string]attr.Value{
			"approve_policy_type_configuration":                             types.ObjectNull(ApprovePolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"backup_creation_policy_type_configuration":                     types.ObjectNull(BackupCreationPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"backup_targets_policy_type_configuration":                      types.ObjectNull(BackupTargetsPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"budget_policy_type_configuration":                              types.ObjectNull(BudgetPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"cluster_resource_name_policy_type_configuration":               types.ObjectNull(ClusterResourceNamePolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"cypher_access_policy_type_configuration":                       types.ObjectNull(CypherAccessPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"delayed_delete_policy_type_configuration":                      types.ObjectNull(DelayedDeletePolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"expiration_policy_type_configuration":                          types.ObjectNull(ExpirationPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"file_share_storage_quota_policy_type_configuration":            types.ObjectNull(FileShareStorageQuotaPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"hostname_policy_type_configuration":                            types.ObjectNull(HostnamePolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"instance_name_policy_type_configuration":                       types.ObjectNull(InstanceNamePolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"max_containers_policy_type_configuration":                      types.ObjectNull(MaxContainersPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"max_cores_policy_type_configuration":                           types.ObjectNull(MaxCoresPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"max_hosts_policy_type_configuration":                           types.ObjectNull(MaxHostsPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"max_load_balancer_pools_policy_type_configuration":             types.ObjectNull(MaxLoadBalancerPoolsPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"max_memory_policy_type_configuration":                          types.ObjectNull(MaxMemoryPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"max_pool_members_policy_type_configuration":                    types.ObjectNull(MaxPoolMembersPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"max_storageand_object_storage_quota_policy_type_configuration": types.ObjectNull(MaxStorageandObjectStorageQuotaPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"max_virtual_servers_policy_type_configuration":                 types.ObjectNull(MaxVirtualServersPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"max_vms_policy_type_configuration":                             types.ObjectNull(MaxVmsPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"messageofthe_day_policy_type_configuration":                    types.ObjectNull(MessageoftheDayPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"network_quota_policy_type_configuration":                       types.ObjectNull(NetworkQuotaPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"power_schedule_policy_type_configuration":                      types.ObjectNull(PowerSchedulePolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"router_quota_policy_type_configuration":                        types.ObjectNull(RouterQuotaPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"shutdown_policy_type_configuration":                            types.ObjectNull(ShutdownPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"storage_server_storage_quota_policy_type_configuration":        types.ObjectNull(StorageServerStorageQuotaPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"tags_policy_type_configuration":                                types.ObjectNull(TagsPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"user_creation_policy_type_configuration":                       types.ObjectNull(UserCreationPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"user_group_creation_policy_type_configuration":                 types.ObjectNull(UserGroupCreationPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
			"workflow_policy_type_configuration":                            types.ObjectNull(WorkflowPolicyTypeConfigurationValue{}.AttributeTypes(ctx)),
		}

		configValue, configDiags := NewConfigValue(ConfigValue{}.AttributeTypes(ctx), configAttrs)
		if configDiags.HasError() {
			diags.Append(configDiags...)
			return state, diags
		}
		state.Config = configValue
	}

	return state, diags
}

// buildPolicyConfigForCreate maps the schema config fields to the SDK config structure for create operations
func buildPolicyConfigForCreate(ctx context.Context, plan *PolicyModel) (*sdk.AddPoliciesRequestPolicyConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	config := &sdk.AddPoliciesRequestPolicyConfig{}
	hasConfig := false

	// ConfigApproval mapping
	if !plan.Config.ApprovePolicyTypeConfiguration.IsNull() && !plan.Config.ApprovePolicyTypeConfiguration.IsUnknown() {
		approveConfig := &sdk.ApprovePolicyTypeConfiguration{}
		attrs := plan.Config.ApprovePolicyTypeConfiguration.Attributes()
		if accountIntegrationId, ok := attrs["account_integration_id"].(basetypes.StringValue); ok && !accountIntegrationId.IsNull() && !accountIntegrationId.IsUnknown() {
			approveConfig.SetAccountIntegrationId(accountIntegrationId.ValueString())
		}
		config.ApprovePolicyTypeConfiguration = approveConfig
		hasConfig = true
	}

	// ConfigBackupCreation mapping
	if !plan.Config.BackupCreationPolicyTypeConfiguration.IsNull() && !plan.Config.BackupCreationPolicyTypeConfiguration.IsUnknown() {
		backupConfig := &sdk.BackupCreationPolicyTypeConfiguration{}
		attrs := plan.Config.BackupCreationPolicyTypeConfiguration.Attributes()
		if createBackup, ok := attrs["create_backup"].(basetypes.BoolValue); ok && !createBackup.IsNull() && !createBackup.IsUnknown() {
			backupConfig.SetCreateBackup(createBackup.ValueBool())
		}
		if createBackupType, ok := attrs["create_backup_type"].(basetypes.StringValue); ok && !createBackupType.IsNull() && !createBackupType.IsUnknown() {
			backupConfig.SetCreateBackupType(createBackupType.ValueString())
		}
		config.BackupCreationPolicyTypeConfiguration = backupConfig
		hasConfig = true
	}

	// ConfigBackupTargets mapping
	if !plan.Config.BackupTargetsPolicyTypeConfiguration.IsNull() && !plan.Config.BackupTargetsPolicyTypeConfiguration.IsUnknown() {
		backupTargetsConfig := &sdk.BackupTargetsPolicyTypeConfiguration{}
		attrs := plan.Config.BackupTargetsPolicyTypeConfiguration.Attributes()
		if backupStorageIds, ok := attrs["backup_storage_ids"].(basetypes.SetValue); ok && !backupStorageIds.IsNull() && !backupStorageIds.IsUnknown() {
			var storageIds []int64
			elemDiags := backupStorageIds.ElementsAs(ctx, &storageIds, false)
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
	if !plan.Config.BudgetPolicyTypeConfiguration.IsNull() && !plan.Config.BudgetPolicyTypeConfiguration.IsUnknown() {
		budgetConfig := &sdk.BudgetPolicyTypeConfiguration{}
		attrs := plan.Config.BudgetPolicyTypeConfiguration.Attributes()
		if maxPrice, ok := attrs["max_price"].(basetypes.NumberValue); ok && !maxPrice.IsNull() && !maxPrice.IsUnknown() {
			maxPriceFloat, _ := maxPrice.ValueBigFloat().Float64()
			budgetConfig.SetMaxPrice(float32(maxPriceFloat))
		}
		if maxPriceCurrency, ok := attrs["max_price_currency"].(basetypes.StringValue); ok && !maxPriceCurrency.IsNull() && !maxPriceCurrency.IsUnknown() {
			budgetConfig.SetMaxPriceCurrency(maxPriceCurrency.ValueString())
		}
		if maxPriceUnit, ok := attrs["max_price_unit"].(basetypes.StringValue); ok && !maxPriceUnit.IsNull() && !maxPriceUnit.IsUnknown() {
			budgetConfig.SetMaxPriceUnit(maxPriceUnit.ValueString())
		}
		config.BudgetPolicyTypeConfiguration = budgetConfig
		hasConfig = true
	}

	// ConfigClusterResourceName mapping
	if !plan.Config.ClusterResourceNamePolicyTypeConfiguration.IsNull() && !plan.Config.ClusterResourceNamePolicyTypeConfiguration.IsUnknown() {
		clusterConfig := &sdk.ClusterResourceNamePolicyTypeConfiguration{}
		attrs := plan.Config.ClusterResourceNamePolicyTypeConfiguration.Attributes()
		if serverNamingType, ok := attrs["server_naming_type"].(basetypes.StringValue); ok && !serverNamingType.IsNull() && !serverNamingType.IsUnknown() {
			clusterConfig.SetServerNamingType(serverNamingType.ValueString())
		}
		if serverNamingPattern, ok := attrs["server_naming_pattern"].(basetypes.StringValue); ok && !serverNamingPattern.IsNull() && !serverNamingPattern.IsUnknown() {
			clusterConfig.SetServerNamingPattern(serverNamingPattern.ValueString())
		}
		if serverNamingConflict, ok := attrs["server_naming_conflict"].(basetypes.BoolValue); ok && !serverNamingConflict.IsNull() && !serverNamingConflict.IsUnknown() {
			clusterConfig.SetServerNamingConflict(serverNamingConflict.ValueBool())
		}
		config.ClusterResourceNamePolicyTypeConfiguration = clusterConfig
		hasConfig = true
	}

	// ConfigCypherAccess mapping
	if !plan.Config.CypherAccessPolicyTypeConfiguration.IsNull() && !plan.Config.CypherAccessPolicyTypeConfiguration.IsUnknown() {
		cypherConfig := &sdk.CypherAccessPolicyTypeConfiguration{}
		attrs := plan.Config.CypherAccessPolicyTypeConfiguration.Attributes()
		if keyPattern, ok := attrs["key_pattern"].(basetypes.StringValue); ok && !keyPattern.IsNull() && !keyPattern.IsUnknown() {
			cypherConfig.SetKeyPattern(keyPattern.ValueString())
		}
		if read, ok := attrs["read"].(basetypes.BoolValue); ok && !read.IsNull() && !read.IsUnknown() {
			cypherConfig.SetRead(read.ValueBool())
		}
		if write, ok := attrs["write"].(basetypes.BoolValue); ok && !write.IsNull() && !write.IsUnknown() {
			cypherConfig.SetWrite(write.ValueBool())
		}
		if update, ok := attrs["update"].(basetypes.BoolValue); ok && !update.IsNull() && !update.IsUnknown() {
			cypherConfig.SetUpdate(update.ValueBool())
		}
		if del, ok := attrs["delete"].(basetypes.BoolValue); ok && !del.IsNull() && !del.IsUnknown() {
			cypherConfig.SetDelete(del.ValueBool())
		}
		if list, ok := attrs["list"].(basetypes.BoolValue); ok && !list.IsNull() && !list.IsUnknown() {
			cypherConfig.SetList(list.ValueBool())
		}
		config.CypherAccessPolicyTypeConfiguration = cypherConfig
		hasConfig = true
	}

	// ConfigDelayedDelete mapping
	if !plan.Config.DelayedDeletePolicyTypeConfiguration.IsNull() && !plan.Config.DelayedDeletePolicyTypeConfiguration.IsUnknown() {
		delayedDeleteConfig := &sdk.DelayedDeletePolicyTypeConfiguration{}
		attrs := plan.Config.DelayedDeletePolicyTypeConfiguration.Attributes()
		if removalAge, ok := attrs["removal_age"].(basetypes.StringValue); ok && !removalAge.IsNull() && !removalAge.IsUnknown() {
			delayedDeleteConfig.SetRemovalAge(removalAge.ValueString())
		}
		config.DelayedDeletePolicyTypeConfiguration = delayedDeleteConfig
		hasConfig = true
	}

	// Continue with more config mappings...

	// ConfigExpiration mapping
	if !plan.Config.ExpirationPolicyTypeConfiguration.IsNull() && !plan.Config.ExpirationPolicyTypeConfiguration.IsUnknown() {
		expirationConfig := &sdk.ExpirationPolicyTypeConfiguration{}
		attrs := plan.Config.ExpirationPolicyTypeConfiguration.Attributes()
		if lifecycle_type, ok := attrs["lifecycle_type"].(basetypes.StringValue); ok && !lifecycle_type.IsNull() && !lifecycle_type.IsUnknown() {
			expirationConfig.SetLifecycleType(lifecycle_type.ValueString())
		}
		if lifecycle_age, ok := attrs["lifecycle_age"].(basetypes.StringValue); ok && !lifecycle_age.IsNull() && !lifecycle_age.IsUnknown() {
			expirationConfig.SetLifecycleAge(lifecycle_age.ValueString())
		}
		if lifecycle_allow_extend, ok := attrs["lifecycle_allow_extend"].(basetypes.StringValue); ok && !lifecycle_allow_extend.IsNull() && !lifecycle_allow_extend.IsUnknown() {
			expirationConfig.SetLifecycleAllowExtend(lifecycle_allow_extend.ValueString())
		}
		if lifecycle_auto_renew, ok := attrs["lifecycle_auto_renew"].(basetypes.StringValue); ok && !lifecycle_auto_renew.IsNull() && !lifecycle_auto_renew.IsUnknown() {
			expirationConfig.SetLifecycleAutoRenew(lifecycle_auto_renew.ValueString())
		}
		if lifecycle_extensions_before_approval, ok := attrs["lifecycle_extensions_before_approval"].(basetypes.StringValue); ok && !lifecycle_extensions_before_approval.IsNull() && !lifecycle_extensions_before_approval.IsUnknown() {
			expirationConfig.SetLifecycleExtensionsBeforeApproval(lifecycle_extensions_before_approval.ValueString())
		}
		if lifecycle_hide_fixed, ok := attrs["lifecycle_hide_fixed"].(basetypes.BoolValue); ok && !lifecycle_hide_fixed.IsNull() && !lifecycle_hide_fixed.IsUnknown() {
			expirationConfig.SetLifecycleHideFixed(lifecycle_hide_fixed.ValueBool())
		}
		if lifecycle_message, ok := attrs["lifecycle_message"].(basetypes.StringValue); ok && !lifecycle_message.IsNull() && !lifecycle_message.IsUnknown() {
			expirationConfig.SetLifecycleMessage(lifecycle_message.ValueString())
		}
		if lifecycle_notify, ok := attrs["lifecycle_notify"].(basetypes.StringValue); ok && !lifecycle_notify.IsNull() && !lifecycle_notify.IsUnknown() {
			expirationConfig.SetLifecycleNotify(lifecycle_notify.ValueString())
		}
		if lifecycle_renewal, ok := attrs["lifecycle_renewal"].(basetypes.StringValue); ok && !lifecycle_renewal.IsNull() && !lifecycle_renewal.IsUnknown() {
			expirationConfig.SetLifecycleRenewal(lifecycle_renewal.ValueString())
		}
		if account_integration_id, ok := attrs["account_integration_id"].(basetypes.StringValue); ok && !account_integration_id.IsNull() && !account_integration_id.IsUnknown() {
			expirationConfig.SetAccountIntegrationId(account_integration_id.ValueString())
		}
		config.ExpirationPolicyTypeConfiguration = expirationConfig
		hasConfig = true
	}

	// ConfigHostname mapping
	if !plan.Config.HostnamePolicyTypeConfiguration.IsNull() && !plan.Config.HostnamePolicyTypeConfiguration.IsUnknown() {
		hostnameConfig := &sdk.HostnamePolicyTypeConfiguration{}
		attrs := plan.Config.HostnamePolicyTypeConfiguration.Attributes()
		if host_naming_type, ok := attrs["host_naming_type"].(basetypes.StringValue); ok && !host_naming_type.IsNull() && !host_naming_type.IsUnknown() {
			hostnameConfig.SetHostNamingType(host_naming_type.ValueString())
		}
		if host_naming_pattern, ok := attrs["host_naming_pattern"].(basetypes.StringValue); ok && !host_naming_pattern.IsNull() && !host_naming_pattern.IsUnknown() {
			hostnameConfig.SetHostNamingPattern(host_naming_pattern.ValueString())
		}
		config.HostnamePolicyTypeConfiguration = hostnameConfig
		hasConfig = true
	}

	// ConfigInstanceName mapping
	if !plan.Config.InstanceNamePolicyTypeConfiguration.IsNull() && !plan.Config.InstanceNamePolicyTypeConfiguration.IsUnknown() {
		instanceNameConfig := &sdk.InstanceNamePolicyTypeConfiguration{}
		attrs := plan.Config.InstanceNamePolicyTypeConfiguration.Attributes()
		if naming_type, ok := attrs["naming_type"].(basetypes.StringValue); ok && !naming_type.IsNull() && !naming_type.IsUnknown() {
			instanceNameConfig.SetNamingType(naming_type.ValueString())
		}
		if naming_pattern, ok := attrs["naming_pattern"].(basetypes.StringValue); ok && !naming_pattern.IsNull() && !naming_pattern.IsUnknown() {
			instanceNameConfig.SetNamingPattern(naming_pattern.ValueString())
		}
		if naming_conflict, ok := attrs["naming_conflict"].(basetypes.BoolValue); ok && !naming_conflict.IsNull() && !naming_conflict.IsUnknown() {
			instanceNameConfig.SetNamingConflict(naming_conflict.ValueBool())
		}
		config.InstanceNamePolicyTypeConfiguration = instanceNameConfig
		hasConfig = true
	}

	// ConfigMaxContainers mapping
	if !plan.Config.MaxContainersPolicyTypeConfiguration.IsNull() && !plan.Config.MaxContainersPolicyTypeConfiguration.IsUnknown() {
		maxContainersConfig := &sdk.MaxContainersPolicyTypeConfiguration{}
		attrs := plan.Config.MaxContainersPolicyTypeConfiguration.Attributes()
		if max_containers, ok := attrs["max_containers"].(basetypes.StringValue); ok && !max_containers.IsNull() && !max_containers.IsUnknown() {
			maxContainersConfig.SetMaxContainers(max_containers.ValueString())
		}
		config.MaxContainersPolicyTypeConfiguration = maxContainersConfig
		hasConfig = true
	}

	// ConfigMaxCores mapping
	if !plan.Config.MaxCoresPolicyTypeConfiguration.IsNull() && !plan.Config.MaxCoresPolicyTypeConfiguration.IsUnknown() {
		maxCoresConfig := &sdk.MaxCoresPolicyTypeConfiguration{}
		attrs := plan.Config.MaxCoresPolicyTypeConfiguration.Attributes()
		if max_cores, ok := attrs["max_cores"].(basetypes.StringValue); ok && !max_cores.IsNull() && !max_cores.IsUnknown() {
			maxCoresConfig.SetMaxCores(max_cores.ValueString())
		}
		if exclude_containers, ok := attrs["exclude_containers"].(basetypes.BoolValue); ok && !exclude_containers.IsNull() && !exclude_containers.IsUnknown() {
			if exclude_containers.ValueBool() {
				maxCoresConfig.SetExcludeContainers("on")
			} else {
				maxCoresConfig.SetExcludeContainers("off")
			}
		}
		config.MaxCoresPolicyTypeConfiguration = maxCoresConfig
		hasConfig = true
	}

	// ConfigMaxHosts mapping
	if !plan.Config.MaxHostsPolicyTypeConfiguration.IsNull() && !plan.Config.MaxHostsPolicyTypeConfiguration.IsUnknown() {
		maxHostsConfig := &sdk.MaxHostsPolicyTypeConfiguration{}
		attrs := plan.Config.MaxHostsPolicyTypeConfiguration.Attributes()
		if max_hosts, ok := attrs["max_hosts"].(basetypes.StringValue); ok && !max_hosts.IsNull() && !max_hosts.IsUnknown() {
			maxHostsConfig.SetMaxHosts(max_hosts.ValueString())
		}
		config.MaxHostsPolicyTypeConfiguration = maxHostsConfig
		hasConfig = true
	}

	// ConfigMaxLoadBalancerPools mapping
	if !plan.Config.MaxLoadBalancerPoolsPolicyTypeConfiguration.IsNull() && !plan.Config.MaxLoadBalancerPoolsPolicyTypeConfiguration.IsUnknown() {
		maxLBPoolsConfig := &sdk.MaxLoadBalancerPoolsPolicyTypeConfiguration{}
		attrs := plan.Config.MaxLoadBalancerPoolsPolicyTypeConfiguration.Attributes()
		if max_pools, ok := attrs["max_pools"].(basetypes.StringValue); ok && !max_pools.IsNull() && !max_pools.IsUnknown() {
			maxLBPoolsConfig.SetMaxPools(max_pools.ValueString())
		}
		config.MaxLoadBalancerPoolsPolicyTypeConfiguration = maxLBPoolsConfig
		hasConfig = true
	}

	// ConfigMaxMemory mapping (complex nested structure)
	if !plan.Config.MaxMemoryPolicyTypeConfiguration.IsNull() && !plan.Config.MaxMemoryPolicyTypeConfiguration.IsUnknown() {
		maxMemoryConfig := &sdk.MaxMemoryPolicyTypeConfiguration{}
		attrs := plan.Config.MaxMemoryPolicyTypeConfiguration.Attributes()
		maxMemorySet := false

		if exclude_containers, ok := attrs["exclude_containers"].(basetypes.BoolValue); ok && !exclude_containers.IsNull() && !exclude_containers.IsUnknown() {
			if exclude_containers.ValueBool() {
				maxMemoryConfig.SetExcludeContainers("on")
			} else {
				maxMemoryConfig.SetExcludeContainers("off")
			}
		}

		// Handle MaxMemory field - it's an ObjectValue with anyof0 (string) and anyof1 (int)
		if maxMemory, ok := attrs["max_memory"].(basetypes.ObjectValue); ok && !maxMemory.IsNull() && !maxMemory.IsUnknown() {
			maxMemoryValue := maxMemory
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
	if !plan.Config.MaxPoolMembersPolicyTypeConfiguration.IsNull() && !plan.Config.MaxPoolMembersPolicyTypeConfiguration.IsUnknown() {
		maxPoolMembersConfig := &sdk.MaxPoolMembersPolicyTypeConfiguration{}
		attrs := plan.Config.MaxPoolMembersPolicyTypeConfiguration.Attributes()
		if max_pool_members, ok := attrs["max_pool_members"].(basetypes.StringValue); ok && !max_pool_members.IsNull() && !max_pool_members.IsUnknown() {
			maxPoolMembersConfig.SetMaxPoolMembers(max_pool_members.ValueString())
		}
		config.MaxPoolMembersPolicyTypeConfiguration = maxPoolMembersConfig
		hasConfig = true
	}

	// ConfigMaxStorage mapping
	if !plan.Config.MaxStorageandObjectStorageQuotaPolicyTypeConfiguration.IsNull() && !plan.Config.MaxStorageandObjectStorageQuotaPolicyTypeConfiguration.IsUnknown() {
		maxStorageConfig := &sdk.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration{}
		attrs := plan.Config.MaxStorageandObjectStorageQuotaPolicyTypeConfiguration.Attributes()
		if max_storage, ok := attrs["max_storage"].(basetypes.StringValue); ok && !max_storage.IsNull() && !max_storage.IsUnknown() {
			maxStorageConfig.SetMaxStorage(max_storage.ValueString())
		}
		if exclude_containers, ok := attrs["exclude_containers"].(basetypes.BoolValue); ok && !exclude_containers.IsNull() && !exclude_containers.IsUnknown() {
			if exclude_containers.ValueBool() {
				maxStorageConfig.SetExcludeContainers("on")
			} else {
				maxStorageConfig.SetExcludeContainers("off")
			}
		}
		config.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration = maxStorageConfig
		hasConfig = true
	}

	// ConfigMaxVirtualServers mapping
	if !plan.Config.MaxVirtualServersPolicyTypeConfiguration.IsNull() && !plan.Config.MaxVirtualServersPolicyTypeConfiguration.IsUnknown() {
		maxVirtualServersConfig := &sdk.MaxVirtualServersPolicyTypeConfiguration{}
		attrs := plan.Config.MaxVirtualServersPolicyTypeConfiguration.Attributes()
		if max_virtual_servers, ok := attrs["max_virtual_servers"].(basetypes.StringValue); ok && !max_virtual_servers.IsNull() && !max_virtual_servers.IsUnknown() {
			maxVirtualServersConfig.SetMaxVirtualServers(max_virtual_servers.ValueString())
		}
		config.MaxVirtualServersPolicyTypeConfiguration = maxVirtualServersConfig
		hasConfig = true
	}

	// ConfigMaxVms mapping
	if !plan.Config.MaxVmsPolicyTypeConfiguration.IsNull() && !plan.Config.MaxVmsPolicyTypeConfiguration.IsUnknown() {
		maxVmsConfig := &sdk.MaxVMsPolicyTypeConfiguration{}
		attrs := plan.Config.MaxVmsPolicyTypeConfiguration.Attributes()
		if max_vms, ok := attrs["max_vms"].(basetypes.StringValue); ok && !max_vms.IsNull() && !max_vms.IsUnknown() {
			maxVmsConfig.SetMaxVms(max_vms.ValueString())
		}
		config.MaxVMsPolicyTypeConfiguration = maxVmsConfig
		hasConfig = true
	}

	// ConfigMessageOfTheDay mapping
	if !plan.Config.MessageoftheDayPolicyTypeConfiguration.IsNull() && !plan.Config.MessageoftheDayPolicyTypeConfiguration.IsUnknown() {
		motdConfig := &sdk.MessageOfTheDayPolicyTypeConfiguration{}
		attrs := plan.Config.MessageoftheDayPolicyTypeConfiguration.Attributes()
		if motddate, ok := attrs["motddate"].(basetypes.StringValue); ok && !motddate.IsNull() && !motddate.IsUnknown() {
			// Parse the date string - assuming RFC3339 format
			if dateStr := motddate.ValueString(); dateStr != "" {
				if parsedDate, err := time.Parse(time.RFC3339, dateStr); err == nil {
					motdConfig.SetMotdDate(parsedDate)
				}
				// If parsing fails, we could add to diagnostics, but for now continue without setting
			}
		}
		if motdmessage, ok := attrs["motdmessage"].(basetypes.StringValue); ok && !motdmessage.IsNull() && !motdmessage.IsUnknown() {
			motdConfig.SetMotdMessage(motdmessage.ValueString())
		}
		if motdtitle, ok := attrs["motdtitle"].(basetypes.StringValue); ok && !motdtitle.IsNull() && !motdtitle.IsUnknown() {
			motdConfig.SetMotdTitle(motdtitle.ValueString())
		}
		if motdtype, ok := attrs["motdtype"].(basetypes.StringValue); ok && !motdtype.IsNull() && !motdtype.IsUnknown() {
			motdConfig.SetMotdType(motdtype.ValueString())
		}
		if motd_full_page, ok := attrs["motd_full_page"].(basetypes.BoolValue); ok && !motd_full_page.IsNull() && !motd_full_page.IsUnknown() {
			motdConfig.SetMotdFullPage(motd_full_page.ValueBool())
		}
		config.MessageOfTheDayPolicyTypeConfiguration = motdConfig
		hasConfig = true
	}

	// ConfigNetworkQuota mapping
	if !plan.Config.NetworkQuotaPolicyTypeConfiguration.IsNull() && !plan.Config.NetworkQuotaPolicyTypeConfiguration.IsUnknown() {
		networkQuotaConfig := &sdk.NetworkQuotaPolicyTypeConfiguration{}
		attrs := plan.Config.NetworkQuotaPolicyTypeConfiguration.Attributes()
		if max_networks, ok := attrs["max_networks"].(basetypes.StringValue); ok && !max_networks.IsNull() && !max_networks.IsUnknown() {
			networkQuotaConfig.SetMaxNetworks(max_networks.ValueString())
		}
		config.NetworkQuotaPolicyTypeConfiguration = networkQuotaConfig
		hasConfig = true
	}

	// ConfigPowerSchedule mapping
	if !plan.Config.PowerSchedulePolicyTypeConfiguration.IsNull() && !plan.Config.PowerSchedulePolicyTypeConfiguration.IsUnknown() {
		powerScheduleConfig := &sdk.PowerSchedulePolicyTypeConfiguration{}
		attrs := plan.Config.PowerSchedulePolicyTypeConfiguration.Attributes()
		if powerSchedule, ok := attrs["power_schedule"].(basetypes.StringValue); ok && !powerSchedule.IsNull() && !powerSchedule.IsUnknown() {
			powerScheduleConfig.SetPowerSchedule(powerSchedule.ValueString())
		}
		if powerScheduleType, ok := attrs["power_schedule_type"].(basetypes.StringValue); ok && !powerScheduleType.IsNull() && !powerScheduleType.IsUnknown() {
			powerScheduleConfig.SetPowerScheduleType(powerScheduleType.ValueString())
		}
		if powerScheduleHideFixed, ok := attrs["power_schedule_hide_fixed"].(basetypes.BoolValue); ok && !powerScheduleHideFixed.IsNull() && !powerScheduleHideFixed.IsUnknown() {
			powerScheduleConfig.SetPowerScheduleHideFixed(powerScheduleHideFixed.ValueBool())
		}
		config.PowerSchedulePolicyTypeConfiguration = powerScheduleConfig
		hasConfig = true
	}

	// ConfigRouterQuota mapping
	if !plan.Config.RouterQuotaPolicyTypeConfiguration.IsNull() && !plan.Config.RouterQuotaPolicyTypeConfiguration.IsUnknown() {
		routerQuotaConfig := &sdk.RouterQuotaPolicyTypeConfiguration{}
		attrs := plan.Config.RouterQuotaPolicyTypeConfiguration.Attributes()
		if max_routers, ok := attrs["max_routers"].(basetypes.StringValue); ok && !max_routers.IsNull() && !max_routers.IsUnknown() {
			routerQuotaConfig.SetMaxRouters(max_routers.ValueString())
		}
		config.RouterQuotaPolicyTypeConfiguration = routerQuotaConfig
		hasConfig = true
	}

	// ConfigShutdown mapping
	if !plan.Config.ShutdownPolicyTypeConfiguration.IsNull() && !plan.Config.ShutdownPolicyTypeConfiguration.IsUnknown() {
		shutdownConfig := &sdk.ShutdownPolicyTypeConfiguration{}
		// ConfigShutdown appears to be a simple config without specific fields to map
		config.ShutdownPolicyTypeConfiguration = shutdownConfig
		hasConfig = true
	}

	// ConfigStorageServerStorageQuota mapping
	if !plan.Config.StorageServerStorageQuotaPolicyTypeConfiguration.IsNull() && !plan.Config.StorageServerStorageQuotaPolicyTypeConfiguration.IsUnknown() {
		storageServerQuotaConfig := &sdk.StorageServerStorageQuotaPolicyTypeConfiguration{}
		attrs := plan.Config.StorageServerStorageQuotaPolicyTypeConfiguration.Attributes()
		if max_storage, ok := attrs["max_storage"].(basetypes.StringValue); ok && !max_storage.IsNull() && !max_storage.IsUnknown() {
			storageServerQuotaConfig.SetMaxStorage(max_storage.ValueString())
		}
		config.StorageServerStorageQuotaPolicyTypeConfiguration = storageServerQuotaConfig
		hasConfig = true
	}

	// ConfigTags mapping
	if !plan.Config.TagsPolicyTypeConfiguration.IsNull() && !plan.Config.TagsPolicyTypeConfiguration.IsUnknown() {
		tagsConfig := &sdk.TagsPolicyTypeConfiguration{}
		attrs := plan.Config.TagsPolicyTypeConfiguration.Attributes()
		if key, ok := attrs["key"].(basetypes.StringValue); ok && !key.IsNull() && !key.IsUnknown() {
			tagsConfig.SetKey(key.ValueString())
		}
		if value, ok := attrs["value"].(basetypes.StringValue); ok && !value.IsNull() && !value.IsUnknown() {
			tagsConfig.SetValue(value.ValueString())
		}
		if strict, ok := attrs["strict"].(basetypes.BoolValue); ok && !strict.IsNull() && !strict.IsUnknown() {
			tagsConfig.SetStrict(strict.ValueBool())
		}
		if valueListId, ok := attrs["value_list_id"].(basetypes.StringValue); ok && !valueListId.IsNull() && !valueListId.IsUnknown() {
			tagsConfig.SetValueListId(valueListId.ValueString())
		}
		config.TagsPolicyTypeConfiguration = tagsConfig
		hasConfig = true
	}

	// ConfigUserCreation mapping
	if !plan.Config.UserCreationPolicyTypeConfiguration.IsNull() && !plan.Config.UserCreationPolicyTypeConfiguration.IsUnknown() {
		userCreationConfig := &sdk.UserCreationPolicyTypeConfiguration{}
		attrs := plan.Config.UserCreationPolicyTypeConfiguration.Attributes()
		if create_user, ok := attrs["create_user"].(basetypes.BoolValue); ok && !create_user.IsNull() && !create_user.IsUnknown() {
			userCreationConfig.SetCreateUser(create_user.ValueBool())
		}
		if create_user_type, ok := attrs["create_user_type"].(basetypes.StringValue); ok && !create_user_type.IsNull() && !create_user_type.IsUnknown() {
			userCreationConfig.SetCreateUserType(create_user_type.ValueString())
		}
		config.UserCreationPolicyTypeConfiguration = userCreationConfig
		hasConfig = true
	}

	// ConfigUserGroupCreation mapping
	if !plan.Config.UserGroupCreationPolicyTypeConfiguration.IsNull() && !plan.Config.UserGroupCreationPolicyTypeConfiguration.IsUnknown() {
		userGroupCreationConfig := &sdk.UserGroupCreationPolicyTypeConfiguration{}
		attrs := plan.Config.UserGroupCreationPolicyTypeConfiguration.Attributes()
		if userGroup, ok := attrs["user_group"].(basetypes.StringValue); ok && !userGroup.IsNull() && !userGroup.IsUnknown() {
			userGroupCreationConfig.SetUserGroup(userGroup.ValueString())
		}
		config.UserGroupCreationPolicyTypeConfiguration = userGroupCreationConfig
		hasConfig = true
	}

	// ConfigWorkflow mapping
	if !plan.Config.WorkflowPolicyTypeConfiguration.IsNull() && !plan.Config.WorkflowPolicyTypeConfiguration.IsUnknown() {
		workflowConfig := &sdk.WorkflowPolicyTypeConfiguration{}
		attrs := plan.Config.WorkflowPolicyTypeConfiguration.Attributes()
		if workflowId, ok := attrs["workflow_id"].(basetypes.StringValue); ok && !workflowId.IsNull() && !workflowId.IsUnknown() {
			workflowConfig.SetWorkflowId(workflowId.ValueString())
		}
		config.WorkflowPolicyTypeConfiguration = workflowConfig
		hasConfig = true
	}

	// ConfigFileShareStorageQuota mapping
	if !plan.Config.FileShareStorageQuotaPolicyTypeConfiguration.IsNull() && !plan.Config.FileShareStorageQuotaPolicyTypeConfiguration.IsUnknown() {
		fileShareQuotaConfig := &sdk.FileShareStorageQuotaPolicyTypeConfiguration{}
		attrs := plan.Config.FileShareStorageQuotaPolicyTypeConfiguration.Attributes()
		if max_storage, ok := attrs["max_storage"].(basetypes.StringValue); ok && !max_storage.IsNull() && !max_storage.IsUnknown() {
			fileShareQuotaConfig.SetMaxStorage(max_storage.ValueString())
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
	if !plan.Config.ApprovePolicyTypeConfiguration.IsNull() && !plan.Config.ApprovePolicyTypeConfiguration.IsUnknown() {
		return "approval"
	}
	if !plan.Config.BackupCreationPolicyTypeConfiguration.IsNull() && !plan.Config.BackupCreationPolicyTypeConfiguration.IsUnknown() {
		return "backupCreation"
	}
	if !plan.Config.BackupTargetsPolicyTypeConfiguration.IsNull() && !plan.Config.BackupTargetsPolicyTypeConfiguration.IsUnknown() {
		return "backupTargets"
	}
	if !plan.Config.BudgetPolicyTypeConfiguration.IsNull() && !plan.Config.BudgetPolicyTypeConfiguration.IsUnknown() {
		return "budget"
	}
	if !plan.Config.ClusterResourceNamePolicyTypeConfiguration.IsNull() && !plan.Config.ClusterResourceNamePolicyTypeConfiguration.IsUnknown() {
		return "clusterResourceName"
	}
	if !plan.Config.CypherAccessPolicyTypeConfiguration.IsNull() && !plan.Config.CypherAccessPolicyTypeConfiguration.IsUnknown() {
		return "cypherAccess"
	}
	if !plan.Config.DelayedDeletePolicyTypeConfiguration.IsNull() && !plan.Config.DelayedDeletePolicyTypeConfiguration.IsUnknown() {
		return "delayedDelete"
	}
	if !plan.Config.ExpirationPolicyTypeConfiguration.IsNull() && !plan.Config.ExpirationPolicyTypeConfiguration.IsUnknown() {
		return "expiration"
	}
	if !plan.Config.FileShareStorageQuotaPolicyTypeConfiguration.IsNull() && !plan.Config.FileShareStorageQuotaPolicyTypeConfiguration.IsUnknown() {
		return "fileShareStorageQuota"
	}
	if !plan.Config.HostnamePolicyTypeConfiguration.IsNull() && !plan.Config.HostnamePolicyTypeConfiguration.IsUnknown() {
		return "hostname"
	}
	if !plan.Config.InstanceNamePolicyTypeConfiguration.IsNull() && !plan.Config.InstanceNamePolicyTypeConfiguration.IsUnknown() {
		return "instanceName"
	}
	if !plan.Config.MaxContainersPolicyTypeConfiguration.IsNull() && !plan.Config.MaxContainersPolicyTypeConfiguration.IsUnknown() {
		return "maxContainers"
	}
	if !plan.Config.MaxCoresPolicyTypeConfiguration.IsNull() && !plan.Config.MaxCoresPolicyTypeConfiguration.IsUnknown() {
		return "maxCores"
	}
	if !plan.Config.MaxHostsPolicyTypeConfiguration.IsNull() && !plan.Config.MaxHostsPolicyTypeConfiguration.IsUnknown() {
		return "maxHosts"
	}
	if !plan.Config.MaxLoadBalancerPoolsPolicyTypeConfiguration.IsNull() && !plan.Config.MaxLoadBalancerPoolsPolicyTypeConfiguration.IsUnknown() {
		return "maxLoadBalancerPools"
	}
	if !plan.Config.MaxMemoryPolicyTypeConfiguration.IsNull() && !plan.Config.MaxMemoryPolicyTypeConfiguration.IsUnknown() {
		return "maxMemory"
	}
	if !plan.Config.MaxPoolMembersPolicyTypeConfiguration.IsNull() && !plan.Config.MaxPoolMembersPolicyTypeConfiguration.IsUnknown() {
		return "maxPoolMembers"
	}
	if !plan.Config.MaxStorageandObjectStorageQuotaPolicyTypeConfiguration.IsNull() && !plan.Config.MaxStorageandObjectStorageQuotaPolicyTypeConfiguration.IsUnknown() {
		return "maxStorage"
	}
	if !plan.Config.MaxVirtualServersPolicyTypeConfiguration.IsNull() && !plan.Config.MaxVirtualServersPolicyTypeConfiguration.IsUnknown() {
		return "maxVirtualServers"
	}
	if !plan.Config.MaxVmsPolicyTypeConfiguration.IsNull() && !plan.Config.MaxVmsPolicyTypeConfiguration.IsUnknown() {
		return "maxVms"
	}
	if !plan.Config.MessageoftheDayPolicyTypeConfiguration.IsNull() && !plan.Config.MessageoftheDayPolicyTypeConfiguration.IsUnknown() {
		return "messageOfTheDay"
	}
	if !plan.Config.NetworkQuotaPolicyTypeConfiguration.IsNull() && !plan.Config.NetworkQuotaPolicyTypeConfiguration.IsUnknown() {
		return "networkQuota"
	}
	if !plan.Config.PowerSchedulePolicyTypeConfiguration.IsNull() && !plan.Config.PowerSchedulePolicyTypeConfiguration.IsUnknown() {
		return "powerSchedule"
	}
	if !plan.Config.RouterQuotaPolicyTypeConfiguration.IsNull() && !plan.Config.RouterQuotaPolicyTypeConfiguration.IsUnknown() {
		return "routerQuota"
	}
	if !plan.Config.ShutdownPolicyTypeConfiguration.IsNull() && !plan.Config.ShutdownPolicyTypeConfiguration.IsUnknown() {
		return "shutdown"
	}
	if !plan.Config.StorageServerStorageQuotaPolicyTypeConfiguration.IsNull() && !plan.Config.StorageServerStorageQuotaPolicyTypeConfiguration.IsUnknown() {
		return "storageServerStorageQuota"
	}
	if !plan.Config.TagsPolicyTypeConfiguration.IsNull() && !plan.Config.TagsPolicyTypeConfiguration.IsUnknown() {
		return "tags"
	}
	if !plan.Config.UserCreationPolicyTypeConfiguration.IsNull() && !plan.Config.UserCreationPolicyTypeConfiguration.IsUnknown() {
		return "userCreation"
	}
	if !plan.Config.UserGroupCreationPolicyTypeConfiguration.IsNull() && !plan.Config.UserGroupCreationPolicyTypeConfiguration.IsUnknown() {
		return "userGroupCreation"
	}
	if !plan.Config.WorkflowPolicyTypeConfiguration.IsNull() && !plan.Config.WorkflowPolicyTypeConfiguration.IsUnknown() {
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
	state, diags := getPolicyAsState(ctx, id, client, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// buildPolicyConfigForUpdate maps the schema config fields to the SDK config structure for update operations

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
