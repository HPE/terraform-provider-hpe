# GetClusterPackage200ResponseClusterPackage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Account** | Pointer to **NullableInt64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**RepeatInstall** | Pointer to **bool** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**PackageType** | Pointer to **string** |  | [optional] 
**PackageVersion** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**IconPath** | Pointer to **NullableString** |  | [optional] 
**ImagePath** | Pointer to **NullableString** |  | [optional] 
**DarkImagePath** | Pointer to **NullableString** |  | [optional] 
**SpecTemplates** | Pointer to [**[]GetClusterPackage200ResponseClusterPackageSpecTemplatesInner**](GetClusterPackage200ResponseClusterPackageSpecTemplatesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetClusterPackage200ResponseClusterPackage{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### IconPath (Nullable)

Use the Nullable wrapper methods:
- `obj.IconPath.IsSet()` — check if set
- `obj.IconPath.Get()` — get the inner value (returns pointer)
- `obj.IconPath.Set(&val)` — set the value
- `obj.IconPath.Unset()` — clear the value
### ImagePath (Nullable)

Use the Nullable wrapper methods:
- `obj.ImagePath.IsSet()` — check if set
- `obj.ImagePath.Get()` — get the inner value (returns pointer)
- `obj.ImagePath.Set(&val)` — set the value
- `obj.ImagePath.Unset()` — clear the value
### DarkImagePath (Nullable)

Use the Nullable wrapper methods:
- `obj.DarkImagePath.IsSet()` — check if set
- `obj.DarkImagePath.Get()` — get the inner value (returns pointer)
- `obj.DarkImagePath.Set(&val)` — set the value
- `obj.DarkImagePath.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


