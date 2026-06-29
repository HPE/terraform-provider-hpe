# CreateOsTypeImage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**OsType** | **int64** | id of osType | 
**VirtualImage** | **int64** | id of virtualImage | 
**ProvisionType** | Pointer to **NullableInt64** | id of provisionType | [optional] 
**Zone** | Pointer to **NullableInt64** | id of cloud/zone | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateOsTypeImage{
    // Set fields directly
}
```

### ProvisionType (Nullable)

Use the Nullable wrapper methods:
- `obj.ProvisionType.IsSet()` — check if set
- `obj.ProvisionType.Get()` — get the inner value (returns pointer)
- `obj.ProvisionType.Set(&val)` — set the value
- `obj.ProvisionType.Unset()` — clear the value
### Zone (Nullable)

Use the Nullable wrapper methods:
- `obj.Zone.IsSet()` — check if set
- `obj.Zone.Get()` — get the inner value (returns pointer)
- `obj.Zone.Set(&val)` — set the value
- `obj.Zone.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


