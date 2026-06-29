# SecurityPackage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Type** | Pointer to [**AddSecurityPackages200ResponseAllOfSecurityPackageType**](AddSecurityPackages200ResponseAllOfSecurityPackageType.md) |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Url** | Pointer to **string** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SecurityPackage{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


