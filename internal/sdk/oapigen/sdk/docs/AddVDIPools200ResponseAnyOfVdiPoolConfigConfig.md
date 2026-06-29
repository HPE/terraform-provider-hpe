# AddVDIPools200ResponseAnyOfVdiPoolConfigConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreateUser** | Pointer to **bool** |  | [optional] 
**IsEC2** | Pointer to **bool** |  | [optional] 
**IsVpcSelectable** | Pointer to **bool** |  | [optional] 
**NoAgent** | Pointer to **bool** |  | [optional] 
**SmbiosAssetTag** | Pointer to **NullableString** |  | [optional] 
**NestedVirtualization** | Pointer to **NullableString** |  | [optional] 
**VmwareFolderId** | Pointer to **string** |  | [optional] 
**ResourcePoolId** | Pointer to **int64** |  | [optional] 
**PoolProviderType** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddVDIPools200ResponseAnyOfVdiPoolConfigConfig{
    // Set fields directly
}
```

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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


