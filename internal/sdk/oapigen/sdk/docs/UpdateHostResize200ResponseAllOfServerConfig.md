# UpdateHostResize200ResponseAllOfServerConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PoolProviderType** | Pointer to **NullableString** |  | [optional] 
**IsVpcSelectable** | Pointer to **bool** |  | [optional] 
**SmbiosAssetTag** | Pointer to **NullableString** |  | [optional] 
**IsEC2** | Pointer to **bool** |  | [optional] 
**ResourcePoolId** | Pointer to [**UpdateHostResize200ResponseAllOfServerConfigResourcePoolId**](UpdateHostResize200ResponseAllOfServerConfigResourcePoolId.md) |  | [optional] 
**HostId** | Pointer to **NullableInt64** |  | [optional] 
**CreateUser** | Pointer to [**UpdateHostResize200ResponseAllOfServerConfigCreateUser**](UpdateHostResize200ResponseAllOfServerConfigCreateUser.md) |  | [optional] 
**NestedVirtualization** | Pointer to **NullableString** |  | [optional] 
**VmwareFolderId** | Pointer to **string** |  | [optional] 
**NoAgent** | Pointer to **bool** |  | [optional] 
**PowerScheduleType** | Pointer to **NullableInt64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateHostResize200ResponseAllOfServerConfig{
    // Set fields directly
}
```

### PoolProviderType (Nullable)

Use the Nullable wrapper methods:
- `obj.PoolProviderType.IsSet()` — check if set
- `obj.PoolProviderType.Get()` — get the inner value (returns pointer)
- `obj.PoolProviderType.Set(&val)` — set the value
- `obj.PoolProviderType.Unset()` — clear the value
### SmbiosAssetTag (Nullable)

Use the Nullable wrapper methods:
- `obj.SmbiosAssetTag.IsSet()` — check if set
- `obj.SmbiosAssetTag.Get()` — get the inner value (returns pointer)
- `obj.SmbiosAssetTag.Set(&val)` — set the value
- `obj.SmbiosAssetTag.Unset()` — clear the value
### HostId (Nullable)

Use the Nullable wrapper methods:
- `obj.HostId.IsSet()` — check if set
- `obj.HostId.Get()` — get the inner value (returns pointer)
- `obj.HostId.Set(&val)` — set the value
- `obj.HostId.Unset()` — clear the value
### NestedVirtualization (Nullable)

Use the Nullable wrapper methods:
- `obj.NestedVirtualization.IsSet()` — check if set
- `obj.NestedVirtualization.Get()` — get the inner value (returns pointer)
- `obj.NestedVirtualization.Set(&val)` — set the value
- `obj.NestedVirtualization.Unset()` — clear the value
### PowerScheduleType (Nullable)

Use the Nullable wrapper methods:
- `obj.PowerScheduleType.IsSet()` — check if set
- `obj.PowerScheduleType.Get()` — get the inner value (returns pointer)
- `obj.PowerScheduleType.Set(&val)` — set the value
- `obj.PowerScheduleType.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


