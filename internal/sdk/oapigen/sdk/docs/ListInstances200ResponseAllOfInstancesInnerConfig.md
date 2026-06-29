# ListInstances200ResponseAllOfInstancesInnerConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreateUser** | Pointer to **bool** |  | [optional] 
**IsEC2** | Pointer to **string** |  | [optional] 
**IsVpcSelectable** | Pointer to **bool** |  | [optional] 
**NoAgent** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerConfigNoAgent**](ListInstances200ResponseAllOfInstancesInnerConfigNoAgent.md) |  | [optional] 
**SecurityGroups** | Pointer to [**[]ListInstances200ResponseAllOfInstancesInnerConfigSecurityGroupsInner**](ListInstances200ResponseAllOfInstancesInnerConfigSecurityGroupsInner.md) |  | [optional] 
**KvmHostId** | Pointer to **NullableInt64** |  | [optional] 
**CloneInstanceId** | Pointer to **int64** | The ID of the source instance this instance was cloned from. Set by Morpheus when the instance is created via the clone endpoint; absent for instances that were not cloned. | [optional] 
**SmbiosAssetTag** | Pointer to **NullableString** |  | [optional] 
**NestedVirtualization** | Pointer to **NullableString** |  | [optional] 
**VmwareFolderId** | Pointer to **string** |  | [optional] 
**CustomOptions** | Pointer to **map[string]interface{}** |  | [optional] 
**ResourcePoolId** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerConfigResourcePoolId**](ListInstances200ResponseAllOfInstancesInnerConfigResourcePoolId.md) |  | [optional] 
**PoolProviderType** | Pointer to **NullableString** |  | [optional] 
**UserGroup** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerConfigUserGroup**](ListInstances200ResponseAllOfInstancesInnerConfigUserGroup.md) |  | [optional] 
**ExpireDays** | Pointer to **string** |  | [optional] 
**ShutdownDays** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**HostName** | Pointer to **string** |  | [optional] 
**InstanceType** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerConfigInstanceType**](ListInstances200ResponseAllOfInstancesInnerConfigInstanceType.md) |  | [optional] 
**Site** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerConfigSite**](ListInstances200ResponseAllOfInstancesInnerConfigSite.md) |  | [optional] 
**EnvironmentPrefix** | Pointer to **NullableString** |  | [optional] 
**Layout** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerConfigLayout**](ListInstances200ResponseAllOfInstancesInnerConfigLayout.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**KmsKeyId** | Pointer to **NullableString** |  | [optional] 
**InstanceProfile** | Pointer to **NullableString** |  | [optional] 
**AvailabilityId** | Pointer to **NullableString** |  | [optional] 
**PublicIpType** | Pointer to **NullableString** |  | [optional] 
**InstanceContext** | Pointer to **string** |  | [optional] 
**MemoryDisplay** | Pointer to **string** |  | [optional] 
**Expose** | Pointer to **[]int64** |  | [optional] 
**CreateBackup** | Pointer to **bool** |  | [optional] 
**Backup** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerConfigBackup**](ListInstances200ResponseAllOfInstancesInnerConfigBackup.md) |  | [optional] 
**ReplicationGroup** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerConfigReplicationGroup**](ListInstances200ResponseAllOfInstancesInnerConfigReplicationGroup.md) |  | [optional] 
**LayoutSize** | Pointer to **int64** |  | [optional] 
**LbInstances** | Pointer to **[]map[string]interface{}** |  | [optional] 
**NetworkDomain** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerConfigNetworkDomain**](ListInstances200ResponseAllOfInstancesInnerConfigNetworkDomain.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListInstances200ResponseAllOfInstancesInnerConfig{
    // Set fields directly
}
```

### KvmHostId (Nullable)

Use the Nullable wrapper methods:
- `obj.KvmHostId.IsSet()` — check if set
- `obj.KvmHostId.Get()` — get the inner value (returns pointer)
- `obj.KvmHostId.Set(&val)` — set the value
- `obj.KvmHostId.Unset()` — clear the value
### SmbiosAssetTag (Nullable)

Use the Nullable wrapper methods:
- `obj.SmbiosAssetTag.IsSet()` — check if set
- `obj.SmbiosAssetTag.Get()` — get the inner value (returns pointer)
- `obj.SmbiosAssetTag.Set(&val)` — set the value
- `obj.SmbiosAssetTag.Unset()` — clear the value
### NestedVirtualization (Nullable)

Use the Nullable wrapper methods:
- `obj.NestedVirtualization.IsSet()` — check if set
- `obj.NestedVirtualization.Get()` — get the inner value (returns pointer)
- `obj.NestedVirtualization.Set(&val)` — set the value
- `obj.NestedVirtualization.Unset()` — clear the value
### PoolProviderType (Nullable)

Use the Nullable wrapper methods:
- `obj.PoolProviderType.IsSet()` — check if set
- `obj.PoolProviderType.Get()` — get the inner value (returns pointer)
- `obj.PoolProviderType.Set(&val)` — set the value
- `obj.PoolProviderType.Unset()` — clear the value
### EnvironmentPrefix (Nullable)

Use the Nullable wrapper methods:
- `obj.EnvironmentPrefix.IsSet()` — check if set
- `obj.EnvironmentPrefix.Get()` — get the inner value (returns pointer)
- `obj.EnvironmentPrefix.Set(&val)` — set the value
- `obj.EnvironmentPrefix.Unset()` — clear the value
### KmsKeyId (Nullable)

Use the Nullable wrapper methods:
- `obj.KmsKeyId.IsSet()` — check if set
- `obj.KmsKeyId.Get()` — get the inner value (returns pointer)
- `obj.KmsKeyId.Set(&val)` — set the value
- `obj.KmsKeyId.Unset()` — clear the value
### InstanceProfile (Nullable)

Use the Nullable wrapper methods:
- `obj.InstanceProfile.IsSet()` — check if set
- `obj.InstanceProfile.Get()` — get the inner value (returns pointer)
- `obj.InstanceProfile.Set(&val)` — set the value
- `obj.InstanceProfile.Unset()` — clear the value
### AvailabilityId (Nullable)

Use the Nullable wrapper methods:
- `obj.AvailabilityId.IsSet()` — check if set
- `obj.AvailabilityId.Get()` — get the inner value (returns pointer)
- `obj.AvailabilityId.Set(&val)` — set the value
- `obj.AvailabilityId.Unset()` — clear the value
### PublicIpType (Nullable)

Use the Nullable wrapper methods:
- `obj.PublicIpType.IsSet()` — check if set
- `obj.PublicIpType.Get()` — get the inner value (returns pointer)
- `obj.PublicIpType.Set(&val)` — set the value
- `obj.PublicIpType.Unset()` — clear the value
### LbInstances (Nullable)

Use the Nullable wrapper methods:
- `obj.LbInstances.IsSet()` — check if set
- `obj.LbInstances.Get()` — get the inner value (returns pointer)
- `obj.LbInstances.Set(&val)` — set the value
- `obj.LbInstances.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


