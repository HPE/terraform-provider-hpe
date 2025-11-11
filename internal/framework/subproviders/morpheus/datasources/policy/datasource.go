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
func mapPolicyConfigToState(ctx context.Context, apiConfig *sdk.AddPolicies200ResponseAllOfPolicyConfig) (ConfigValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	// Initialize all config nested objects as null - required by Terraform Plugin Framework
	configAttrs := map[string]attr.Value{
		"approval":             types.ObjectNull(ApprovalValue{}.AttributeTypes(ctx)),
		"backup_storage":       types.ObjectNull(BackupStorageValue{}.AttributeTypes(ctx)),
		"create_backup":        types.ObjectNull(CreateBackupValue{}.AttributeTypes(ctx)),
		"create_user":          types.ObjectNull(CreateUserValue{}.AttributeTypes(ctx)),
		"create_user_group":    types.ObjectNull(CreateUserGroupValue{}.AttributeTypes(ctx)),
		"cypher":               types.ObjectNull(CypherValue{}.AttributeTypes(ctx)),
		"delayed_removal":      types.ObjectNull(DelayedRemovalValue{}.AttributeTypes(ctx)),
		"host_naming":          types.ObjectNull(HostNamingValue{}.AttributeTypes(ctx)),
		"lifecycle":            types.ObjectNull(LifecycleValue{}.AttributeTypes(ctx)),
		"max_containers":       types.ObjectNull(MaxContainersValue{}.AttributeTypes(ctx)),
		"max_cores":            types.ObjectNull(MaxCoresValue{}.AttributeTypes(ctx)),
		"max_hosts":            types.ObjectNull(MaxHostsValue{}.AttributeTypes(ctx)),
		"max_memory":           types.ObjectNull(MaxMemoryValue{}.AttributeTypes(ctx)),
		"max_networks":         types.ObjectNull(MaxNetworksValue{}.AttributeTypes(ctx)),
		"max_pool_members":     types.ObjectNull(MaxPoolMembersValue{}.AttributeTypes(ctx)),
		"max_pools":            types.ObjectNull(MaxPoolsValue{}.AttributeTypes(ctx)),
		"max_price":            types.ObjectNull(MaxPriceValue{}.AttributeTypes(ctx)),
		"max_routers":          types.ObjectNull(MaxRoutersValue{}.AttributeTypes(ctx)),
		"max_storage":          types.ObjectNull(MaxStorageValue{}.AttributeTypes(ctx)),
		"max_virtual_servers":  types.ObjectNull(MaxVirtualServersValue{}.AttributeTypes(ctx)),
		"max_vms":              types.ObjectNull(MaxVmsValue{}.AttributeTypes(ctx)),
		"motd":                 types.ObjectNull(MotdValue{}.AttributeTypes(ctx)),
		"naming":               types.ObjectNull(NamingValue{}.AttributeTypes(ctx)),
		"power_schedule":       types.ObjectNull(PowerScheduleValue{}.AttributeTypes(ctx)),
		"required_network":     types.ObjectNull(RequiredNetworkValue{}.AttributeTypes(ctx)),
		"server_naming":        types.ObjectNull(ServerNamingValue{}.AttributeTypes(ctx)),
		"shutdown":             types.ObjectNull(ShutdownValue{}.AttributeTypes(ctx)),
		"storage_server_quota": types.ObjectNull(StorageServerQuotaValue{}.AttributeTypes(ctx)),
		"tags":                 types.ObjectNull(TagsValue{}.AttributeTypes(ctx)),
		"workflow":             types.ObjectNull(WorkflowValue{}.AttributeTypes(ctx)),
	}

	// Map each API config field to the corresponding schema field - only populate non-null configurations
	if apiConfig.ApprovePolicyTypeConfiguration != nil {
		approvalValue, approvalDiags := NewApprovalValue(
			ApprovalValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"account_integration_id": convert.StrToType(apiConfig.ApprovePolicyTypeConfiguration.AccountIntegrationId),
			},
		)
		if approvalDiags.HasError() {
			diags.Append(approvalDiags...)
		} else {
			objectValue, objectDiags := approvalValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["approval"] = objectValue
			}
		}
	}

	if apiConfig.BackupTargetsPolicyTypeConfiguration != nil {
		// Handle BackupStorageIds as a set of strings
		var backupStorageIdsSet types.Set
		if len(apiConfig.BackupTargetsPolicyTypeConfiguration.BackupStorageIds) == 0 {
			backupStorageIdsSet = types.SetValueMust(types.StringType, []attr.Value{})
		} else {
			stringValues := make([]attr.Value, len(apiConfig.BackupTargetsPolicyTypeConfiguration.BackupStorageIds))
			for i, id := range apiConfig.BackupTargetsPolicyTypeConfiguration.BackupStorageIds {
				stringValues[i] = types.StringValue(id)
			}
			var setDiags diag.Diagnostics
			backupStorageIdsSet, setDiags = types.SetValueFrom(ctx, types.StringType, stringValues)
			if setDiags.HasError() {
				diags.Append(setDiags...)
			}
		}

		backupStorageAttrs := map[string]attr.Value{
			"backup_storage_ids": backupStorageIdsSet,
		}

		backupStorageValue, backupStorageDiags := NewBackupStorageValue(BackupStorageValue{}.AttributeTypes(ctx), backupStorageAttrs)
		if backupStorageDiags.HasError() {
			diags.Append(backupStorageDiags...)
		} else {
			objectValue, objectDiags := backupStorageValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["backup_storage"] = objectValue
			}
		}
	}

	if apiConfig.BackupCreationPolicyTypeConfiguration != nil {
		createBackupValue, createBackupDiags := NewCreateBackupValue(
			CreateBackupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"account_integration_id": convert.StrToType(apiConfig.BackupCreationPolicyTypeConfiguration.AccountIntegrationId),
				"create_backup":          convert.BoolToType(apiConfig.BackupCreationPolicyTypeConfiguration.CreateBackup),
				"create_backup_type":     types.StringValue(apiConfig.BackupCreationPolicyTypeConfiguration.CreateBackupType),
			},
		)
		if createBackupDiags.HasError() {
			diags.Append(createBackupDiags...)
		} else {
			objectValue, objectDiags := createBackupValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["create_backup"] = objectValue
			}
		}
	}

	if apiConfig.UserCreationPolicyTypeConfiguration != nil {
		createUserValue, createUserDiags := NewCreateUserValue(
			CreateUserValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"create_user":      convert.BoolToType(apiConfig.UserCreationPolicyTypeConfiguration.CreateUser),
				"create_user_type": types.StringValue(apiConfig.UserCreationPolicyTypeConfiguration.CreateUserType),
			},
		)
		if createUserDiags.HasError() {
			diags.Append(createUserDiags...)
		} else {
			objectValue, objectDiags := createUserValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["create_user"] = objectValue
			}
		}
	}

	if apiConfig.UserGroupCreationPolicyTypeConfiguration != nil {
		createUserGroupValue, createUserGroupDiags := NewCreateUserGroupValue(
			CreateUserGroupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"user_group": types.StringValue(apiConfig.UserGroupCreationPolicyTypeConfiguration.UserGroup),
			},
		)
		if createUserGroupDiags.HasError() {
			diags.Append(createUserGroupDiags...)
		} else {
			objectValue, objectDiags := createUserGroupValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["create_user_group"] = objectValue
			}
		}
	}

	if apiConfig.CypherAccessPolicyTypeConfiguration != nil {
		cypherValue, cypherDiags := NewCypherValue(
			CypherValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"account_integration_id": convert.StrToType(apiConfig.CypherAccessPolicyTypeConfiguration.AccountIntegrationId),
				"delete":                 convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration.Delete),
				"key_pattern":            types.StringValue(apiConfig.CypherAccessPolicyTypeConfiguration.KeyPattern),
				"list":                   convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration.List),
				"read":                   convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration.Read),
				"update":                 convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration.Update),
				"write":                  convert.BoolToType(apiConfig.CypherAccessPolicyTypeConfiguration.Write),
			},
		)
		if cypherDiags.HasError() {
			diags.Append(cypherDiags...)
		} else {
			objectValue, objectDiags := cypherValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["cypher"] = objectValue
			}
		}
	}

	// Map remaining common config types
	if apiConfig.BudgetPolicyTypeConfiguration != nil {
		maxPriceAttrs := map[string]attr.Value{
			"max_price":          types.StringNull(),
			"max_price_currency": types.StringNull(),
			"max_price_unit":     types.StringNull(),
		}

		if apiConfig.BudgetPolicyTypeConfiguration.MaxPrice != "" {
			maxPriceAttrs["max_price"] = types.StringValue(apiConfig.BudgetPolicyTypeConfiguration.MaxPrice)
		}
		maxPriceAttrs["max_price_currency"] = convert.StrToType(apiConfig.BudgetPolicyTypeConfiguration.MaxPriceCurrency)
		maxPriceAttrs["max_price_unit"] = convert.StrToType(apiConfig.BudgetPolicyTypeConfiguration.MaxPriceUnit)

		maxPriceValue, maxPriceDiags := NewMaxPriceValue(MaxPriceValue{}.AttributeTypes(ctx), maxPriceAttrs)
		if maxPriceDiags.HasError() {
			diags.Append(maxPriceDiags...)
		} else {
			objectValue, objectDiags := maxPriceValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_price"] = objectValue
			}
		}
	}

	if apiConfig.MaxMemoryPolicyTypeConfiguration != nil {
		maxMemoryAttrs := map[string]attr.Value{
			"max_memory":         types.StringNull(),
			"exclude_containers": types.StringNull(),
		}

		if apiConfig.MaxMemoryPolicyTypeConfiguration.MaxMemory != "" {
			maxMemoryAttrs["max_memory"] = types.StringValue(apiConfig.MaxMemoryPolicyTypeConfiguration.MaxMemory)
		}
		maxMemoryAttrs["exclude_containers"] = convert.StrToType(apiConfig.MaxMemoryPolicyTypeConfiguration.ExcludeContainers)

		maxMemoryValue, maxMemoryDiags := NewMaxMemoryValue(MaxMemoryValue{}.AttributeTypes(ctx), maxMemoryAttrs)
		if maxMemoryDiags.HasError() {
			diags.Append(maxMemoryDiags...)
		} else {
			objectValue, objectDiags := maxMemoryValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_memory"] = objectValue
			}
		}
	}

	if apiConfig.MaxCoresPolicyTypeConfiguration != nil {
		maxCoresAttrs := map[string]attr.Value{
			"max_cores":          types.StringNull(),
			"exclude_containers": types.StringNull(),
		}

		if apiConfig.MaxCoresPolicyTypeConfiguration.MaxCores != "" {
			maxCoresAttrs["max_cores"] = types.StringValue(apiConfig.MaxCoresPolicyTypeConfiguration.MaxCores)
		}
		maxCoresAttrs["exclude_containers"] = convert.StrToType(apiConfig.MaxCoresPolicyTypeConfiguration.ExcludeContainers)

		maxCoresValue, maxCoresDiags := NewMaxCoresValue(MaxCoresValue{}.AttributeTypes(ctx), maxCoresAttrs)
		if maxCoresDiags.HasError() {
			diags.Append(maxCoresDiags...)
		} else {
			objectValue, objectDiags := maxCoresValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_cores"] = objectValue
			}
		}
	}

	// 7. DelayedDeletePolicyTypeConfiguration -> delayed_removal
	if apiConfig.DelayedDeletePolicyTypeConfiguration != nil {
		delayedRemovalAttrs := map[string]attr.Value{
			"account_integration_id": convert.StrToType(apiConfig.DelayedDeletePolicyTypeConfiguration.AccountIntegrationId),
			"removal_age":            types.StringValue(apiConfig.DelayedDeletePolicyTypeConfiguration.RemovalAge),
		}

		delayedRemovalValue, delayedRemovalDiags := NewDelayedRemovalValue(DelayedRemovalValue{}.AttributeTypes(ctx), delayedRemovalAttrs)
		if delayedRemovalDiags.HasError() {
			diags.Append(delayedRemovalDiags...)
		} else {
			objectValue, objectDiags := delayedRemovalValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["delayed_removal"] = objectValue
			}
		}
	}

	// 8. ExpirationPolicyTypeConfiguration2 -> lifecycle
	if apiConfig.ExpirationPolicyTypeConfiguration2 != nil {
		lifecycleAttrs := map[string]attr.Value{
			"account_integration_id":               convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.AccountIntegrationId),
			"lifecycle_age":                        convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleAge),
			"lifecycle_allow_extend":               convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleAllowExtend),
			"lifecycle_auto_renew":                 convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleAutoRenew),
			"lifecycle_extensions_before_approval": convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleExtensionsBeforeApproval),
			"lifecycle_hide_fixed":                 convert.BoolToType(apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleHideFixed),
			"lifecycle_message":                    convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleMessage),
			"lifecycle_notify":                     convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleNotify),
			"lifecycle_renewal":                    convert.StrToType(apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleRenewal),
			"lifecycle_type":                       types.StringValue(apiConfig.ExpirationPolicyTypeConfiguration2.LifecycleType),
		}

		lifecycleValue, lifecycleDiags := NewLifecycleValue(LifecycleValue{}.AttributeTypes(ctx), lifecycleAttrs)
		if lifecycleDiags.HasError() {
			diags.Append(lifecycleDiags...)
		} else {
			objectValue, objectDiags := lifecycleValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["lifecycle"] = objectValue
			}
		}
	}

	// 9. HostnamePolicyTypeConfiguration -> host_naming
	if apiConfig.HostnamePolicyTypeConfiguration != nil {
		hostNamingAttrs := map[string]attr.Value{
			"host_naming_pattern": convert.StrToType(apiConfig.HostnamePolicyTypeConfiguration.HostNamingPattern),
			"host_naming_type":    types.StringValue(apiConfig.HostnamePolicyTypeConfiguration.HostNamingType),
		}

		hostNamingValue, hostNamingDiags := NewHostNamingValue(HostNamingValue{}.AttributeTypes(ctx), hostNamingAttrs)
		if hostNamingDiags.HasError() {
			diags.Append(hostNamingDiags...)
		} else {
			objectValue, objectDiags := hostNamingValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["host_naming"] = objectValue
			}
		}
	}

	// 10. InstanceNamePolicyTypeConfiguration -> naming
	if apiConfig.InstanceNamePolicyTypeConfiguration != nil {
		namingAttrs := map[string]attr.Value{
			"naming_conflict": convert.BoolToType(apiConfig.InstanceNamePolicyTypeConfiguration.NamingConflict),
			"naming_pattern":  convert.StrToType(apiConfig.InstanceNamePolicyTypeConfiguration.NamingPattern),
			"naming_type":     types.StringValue(apiConfig.InstanceNamePolicyTypeConfiguration.NamingType),
		}

		namingValue, namingDiags := NewNamingValue(NamingValue{}.AttributeTypes(ctx), namingAttrs)
		if namingDiags.HasError() {
			diags.Append(namingDiags...)
		} else {
			objectValue, objectDiags := namingValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["naming"] = objectValue
			}
		}
	}

	// 11. MaxContainersPolicyTypeConfiguration -> max_containers
	if apiConfig.MaxContainersPolicyTypeConfiguration != nil {
		maxContainersAttrs := map[string]attr.Value{
			"max_containers": types.StringValue(apiConfig.MaxContainersPolicyTypeConfiguration.MaxContainers),
		}

		maxContainersValue, maxContainersDiags := NewMaxContainersValue(MaxContainersValue{}.AttributeTypes(ctx), maxContainersAttrs)
		if maxContainersDiags.HasError() {
			diags.Append(maxContainersDiags...)
		} else {
			objectValue, objectDiags := maxContainersValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_containers"] = objectValue
			}
		}
	}

	// 12. MaxHostsPolicyTypeConfiguration -> max_hosts
	if apiConfig.MaxHostsPolicyTypeConfiguration != nil {
		maxHostsAttrs := map[string]attr.Value{
			"max_hosts": types.StringValue(apiConfig.MaxHostsPolicyTypeConfiguration.MaxHosts),
		}

		maxHostsValue, maxHostsDiags := NewMaxHostsValue(MaxHostsValue{}.AttributeTypes(ctx), maxHostsAttrs)
		if maxHostsDiags.HasError() {
			diags.Append(maxHostsDiags...)
		} else {
			objectValue, objectDiags := maxHostsValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_hosts"] = objectValue
			}
		}
	}

	// 13. NetworkQuotaPolicyTypeConfiguration -> max_networks
	if apiConfig.NetworkQuotaPolicyTypeConfiguration != nil {
		maxNetworksAttrs := map[string]attr.Value{
			"max_networks": types.StringValue(apiConfig.NetworkQuotaPolicyTypeConfiguration.MaxNetworks),
		}

		maxNetworksValue, maxNetworksDiags := NewMaxNetworksValue(MaxNetworksValue{}.AttributeTypes(ctx), maxNetworksAttrs)
		if maxNetworksDiags.HasError() {
			diags.Append(maxNetworksDiags...)
		} else {
			objectValue, objectDiags := maxNetworksValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_networks"] = objectValue
			}
		}
	}

	// 14. MaxPoolMembersPolicyTypeConfiguration -> max_pool_members
	if apiConfig.MaxPoolMembersPolicyTypeConfiguration != nil {
		maxPoolMembersAttrs := map[string]attr.Value{
			"max_pool_members": types.StringValue(apiConfig.MaxPoolMembersPolicyTypeConfiguration.MaxPoolMembers),
		}

		maxPoolMembersValue, maxPoolMembersDiags := NewMaxPoolMembersValue(MaxPoolMembersValue{}.AttributeTypes(ctx), maxPoolMembersAttrs)
		if maxPoolMembersDiags.HasError() {
			diags.Append(maxPoolMembersDiags...)
		} else {
			objectValue, objectDiags := maxPoolMembersValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_pool_members"] = objectValue
			}
		}
	}

	// 15. MaxLoadBalancerPoolsPolicyTypeConfiguration -> max_pools
	if apiConfig.MaxLoadBalancerPoolsPolicyTypeConfiguration != nil {
		maxPoolsAttrs := map[string]attr.Value{
			"max_pools": types.StringValue(apiConfig.MaxLoadBalancerPoolsPolicyTypeConfiguration.MaxPools),
		}

		maxPoolsValue, maxPoolsDiags := NewMaxPoolsValue(MaxPoolsValue{}.AttributeTypes(ctx), maxPoolsAttrs)
		if maxPoolsDiags.HasError() {
			diags.Append(maxPoolsDiags...)
		} else {
			objectValue, objectDiags := maxPoolsValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_pools"] = objectValue
			}
		}
	}

	// 16. RouterQuotaPolicyTypeConfiguration -> max_routers
	if apiConfig.RouterQuotaPolicyTypeConfiguration != nil {
		maxRoutersAttrs := map[string]attr.Value{
			"max_routers": types.StringValue(apiConfig.RouterQuotaPolicyTypeConfiguration.MaxRouters),
		}

		maxRoutersValue, maxRoutersDiags := NewMaxRoutersValue(MaxRoutersValue{}.AttributeTypes(ctx), maxRoutersAttrs)
		if maxRoutersDiags.HasError() {
			diags.Append(maxRoutersDiags...)
		} else {
			objectValue, objectDiags := maxRoutersValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_routers"] = objectValue
			}
		}
	}

	// 17. MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration -> max_storage
	if apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration != nil {
		maxStorageAttrs := map[string]attr.Value{
			"exclude_containers": convert.StrToType(apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration.ExcludeContainers),
			"max_storage":        types.StringValue(apiConfig.MaxStorageAndObjectStorageQuotaPolicyTypeConfiguration.MaxStorage),
		}

		maxStorageValue, maxStorageDiags := NewMaxStorageValue(MaxStorageValue{}.AttributeTypes(ctx), maxStorageAttrs)
		if maxStorageDiags.HasError() {
			diags.Append(maxStorageDiags...)
		} else {
			objectValue, objectDiags := maxStorageValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_storage"] = objectValue
			}
		}
	}

	// 18. MaxVirtualServersPolicyTypeConfiguration -> max_virtual_servers
	if apiConfig.MaxVirtualServersPolicyTypeConfiguration != nil {
		maxVirtualServersAttrs := map[string]attr.Value{
			"max_virtual_servers": types.StringValue(apiConfig.MaxVirtualServersPolicyTypeConfiguration.MaxVirtualServers),
		}

		maxVirtualServersValue, maxVirtualServersDiags := NewMaxVirtualServersValue(MaxVirtualServersValue{}.AttributeTypes(ctx), maxVirtualServersAttrs)
		if maxVirtualServersDiags.HasError() {
			diags.Append(maxVirtualServersDiags...)
		} else {
			objectValue, objectDiags := maxVirtualServersValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_virtual_servers"] = objectValue
			}
		}
	}

	// 19. MaxVMsPolicyTypeConfiguration -> max_vms
	if apiConfig.MaxVMsPolicyTypeConfiguration != nil {
		maxVmsAttrs := map[string]attr.Value{
			"max_vms": types.StringValue(apiConfig.MaxVMsPolicyTypeConfiguration.MaxVms),
		}

		maxVmsValue, maxVmsDiags := NewMaxVmsValue(MaxVmsValue{}.AttributeTypes(ctx), maxVmsAttrs)
		if maxVmsDiags.HasError() {
			diags.Append(maxVmsDiags...)
		} else {
			objectValue, objectDiags := maxVmsValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["max_vms"] = objectValue
			}
		}
	}

	// 20. MessageOfTheDayPolicyTypeConfiguration2 -> motd
	if apiConfig.MessageOfTheDayPolicyTypeConfiguration2 != nil {
		// Handle MotdFullPage
		var motdFullPageValue MotdFullPageValue
		if apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdFullPage != nil && apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdFullPage.String != nil {
			motdFullPageAttrs := map[string]attr.Value{
				"oneof0": types.StringValue(*apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdFullPage.String),
			}
			motdFullPageValue, _ = NewMotdFullPageValue(MotdFullPageValue{}.AttributeTypes(ctx), motdFullPageAttrs)
		} else {
			motdFullPageValue = NewMotdFullPageValueNull()
		}

		motdFullPageObjectValue, motdFullPageDiags := motdFullPageValue.ToObjectValue(ctx)
		if motdFullPageDiags.HasError() {
			diags.Append(motdFullPageDiags...)
		}

		motdAttrs := map[string]attr.Value{
			"motd_full_page": motdFullPageObjectValue,
			"motddate":       convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdDate),
			"motdmessage":    convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdMessage),
			"motdtitle":      types.StringNull(), // NullableString type
			"motdtype":       convert.StrToType(apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdType),
		}
		// Handle NullableString for MotdTitle
		if apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdTitle.IsSet() {
			motdAttrs["motdtitle"] = types.StringValue(*apiConfig.MessageOfTheDayPolicyTypeConfiguration2.MotdTitle.Get())
		}

		motdValue, motdDiags := NewMotdValue(MotdValue{}.AttributeTypes(ctx), motdAttrs)
		if motdDiags.HasError() {
			diags.Append(motdDiags...)
		} else {
			objectValue, objectDiags := motdValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["motd"] = objectValue
			}
		}
	}

	// 21. PowerSchedulePolicyTypeConfiguration -> power_schedule
	if apiConfig.PowerSchedulePolicyTypeConfiguration != nil {
		powerScheduleAttrs := map[string]attr.Value{
			"power_schedule":            convert.StrToType(apiConfig.PowerSchedulePolicyTypeConfiguration.PowerSchedule),
			"power_schedule_hide_fixed": convert.BoolToType(apiConfig.PowerSchedulePolicyTypeConfiguration.PowerScheduleHideFixed),
			"power_schedule_type":       types.StringValue(apiConfig.PowerSchedulePolicyTypeConfiguration.PowerScheduleType),
		}

		powerScheduleValue, powerScheduleDiags := NewPowerScheduleValue(PowerScheduleValue{}.AttributeTypes(ctx), powerScheduleAttrs)
		if powerScheduleDiags.HasError() {
			diags.Append(powerScheduleDiags...)
		} else {
			objectValue, objectDiags := powerScheduleValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["power_schedule"] = objectValue
			}
		}
	}

	// 22. RequiredNetworkPolicyTypeConfiguration -> required_network
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

		requiredNetworkValue, requiredNetworkDiags := NewRequiredNetworkValue(RequiredNetworkValue{}.AttributeTypes(ctx), requiredNetworkAttrs)
		if requiredNetworkDiags.HasError() {
			diags.Append(requiredNetworkDiags...)
		} else {
			objectValue, objectDiags := requiredNetworkValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["required_network"] = objectValue
			}
		}
	}

	// 23. ClusterResourceNamePolicyTypeConfiguration -> server_naming
	if apiConfig.ClusterResourceNamePolicyTypeConfiguration != nil {
		serverNamingAttrs := map[string]attr.Value{
			"account_integration_id": convert.StrToType(apiConfig.ClusterResourceNamePolicyTypeConfiguration.AccountIntegrationId),
			"server_naming_conflict": convert.BoolToType(apiConfig.ClusterResourceNamePolicyTypeConfiguration.ServerNamingConflict),
			"server_naming_pattern":  convert.StrToType(apiConfig.ClusterResourceNamePolicyTypeConfiguration.ServerNamingPattern),
			"server_naming_type":     types.StringValue(apiConfig.ClusterResourceNamePolicyTypeConfiguration.ServerNamingType),
		}

		serverNamingValue, serverNamingDiags := NewServerNamingValue(ServerNamingValue{}.AttributeTypes(ctx), serverNamingAttrs)
		if serverNamingDiags.HasError() {
			diags.Append(serverNamingDiags...)
		} else {
			objectValue, objectDiags := serverNamingValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["server_naming"] = objectValue
			}
		}
	}

	// 24. ShutdownPolicyTypeConfiguration -> shutdown
	if apiConfig.ShutdownPolicyTypeConfiguration != nil {
		shutdownAttrs := map[string]attr.Value{
			"account_integration_id":              convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.AccountIntegrationId),
			"shutdown_age":                        convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.ShutdownAge),
			"shutdown_allow_extend":               convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.ShutdownAllowExtend),
			"shutdown_auto_renew":                 convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.ShutdownAutoRenew),
			"shutdown_extensions_before_approval": convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.ShutdownExtensionsBeforeApproval),
			"shutdown_hide_fixed":                 convert.BoolToType(apiConfig.ShutdownPolicyTypeConfiguration.ShutdownHideFixed),
			"shutdown_message":                    convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.ShutdownMessage),
			"shutdown_notify":                     convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.ShutdownNotify),
			"shutdown_renewal":                    convert.StrToType(apiConfig.ShutdownPolicyTypeConfiguration.ShutdownRenewal),
			"shutdown_type":                       types.StringValue(apiConfig.ShutdownPolicyTypeConfiguration.ShutdownType),
		}

		shutdownValue, shutdownDiags := NewShutdownValue(ShutdownValue{}.AttributeTypes(ctx), shutdownAttrs)
		if shutdownDiags.HasError() {
			diags.Append(shutdownDiags...)
		} else {
			objectValue, objectDiags := shutdownValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["shutdown"] = objectValue
			}
		}
	}

	// 25. StorageServerStorageQuotaPolicyTypeConfiguration -> storage_server_quota
	if apiConfig.StorageServerStorageQuotaPolicyTypeConfiguration != nil {
		storageServerQuotaAttrs := map[string]attr.Value{
			"max_storage":       convert.StrToType(apiConfig.StorageServerStorageQuotaPolicyTypeConfiguration.MaxStorage),
			"storage_server_id": types.StringValue(apiConfig.StorageServerStorageQuotaPolicyTypeConfiguration.StorageServerId),
		}

		storageServerQuotaValue, storageServerQuotaDiags := NewStorageServerQuotaValue(StorageServerQuotaValue{}.AttributeTypes(ctx), storageServerQuotaAttrs)
		if storageServerQuotaDiags.HasError() {
			diags.Append(storageServerQuotaDiags...)
		} else {
			objectValue, objectDiags := storageServerQuotaValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["storage_server_quota"] = objectValue
			}
		}
	}

	// 26. TagsPolicyTypeConfiguration -> tags
	if apiConfig.TagsPolicyTypeConfiguration != nil {
		tagsAttrs := map[string]attr.Value{
			"key":           convert.StrToType(apiConfig.TagsPolicyTypeConfiguration.Key),
			"strict":        types.BoolValue(apiConfig.TagsPolicyTypeConfiguration.Strict),
			"value":         convert.StrToType(apiConfig.TagsPolicyTypeConfiguration.Value),
			"value_list_id": convert.StrToType(apiConfig.TagsPolicyTypeConfiguration.ValueListId),
		}

		tagsValue, tagsDiags := NewTagsValue(TagsValue{}.AttributeTypes(ctx), tagsAttrs)
		if tagsDiags.HasError() {
			diags.Append(tagsDiags...)
		} else {
			objectValue, objectDiags := tagsValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["tags"] = objectValue
			}
		}
	}

	// 27. WorkflowPolicyTypeConfiguration -> workflow
	if apiConfig.WorkflowPolicyTypeConfiguration != nil {
		workflowAttrs := map[string]attr.Value{
			"workflow_id": types.StringValue(apiConfig.WorkflowPolicyTypeConfiguration.WorkflowId),
		}

		workflowValue, workflowDiags := NewWorkflowValue(WorkflowValue{}.AttributeTypes(ctx), workflowAttrs)
		if workflowDiags.HasError() {
			diags.Append(workflowDiags...)
		} else {
			objectValue, objectDiags := workflowValue.ToObjectValue(ctx)
			if objectDiags.HasError() {
				diags.Append(objectDiags...)
			} else {
				configAttrs["workflow"] = objectValue
			}
		}
	}

	// Create the config value
	configValue, configValueDiags := NewConfigValue(ConfigValue{}.AttributeTypes(ctx), configAttrs)
	if configValueDiags.HasError() {
		diags.Append(configValueDiags...)
		return NewConfigValueNull(), diags
	}

	return configValue, diags
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

	// Handle Config - map API config to schema structure
	if policy.Config != nil {
		configValue, configDiags := mapPolicyConfigToState(ctx, policy.Config)
		if configDiags.HasError() {
			diags.Append(configDiags...)
			return diags
		}
		data.Config = configValue
	} else {
		data.Config = NewConfigValueNull()
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
		for _, account := range policy.Accounts {
			tenantValue, tenantDiags := types.ObjectValueFrom(ctx,
				map[string]attr.Type{
					"id":   types.Int64Type,
					"name": types.StringType,
				},
				map[string]attr.Value{
					"id":   convert.Int64ToType(account.Id),
					"name": convert.StrToType(account.Name),
				},
			)
			if tenantDiags.HasError() {
				diags.Append(tenantDiags...)
				return diags
			}
			tenantValues = append(tenantValues, tenantValue)
		}

		tenantsSet, setDiags := types.SetValueFrom(ctx,
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":   types.Int64Type,
					"name": types.StringType,
				},
			},
			tenantValues,
		)
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
