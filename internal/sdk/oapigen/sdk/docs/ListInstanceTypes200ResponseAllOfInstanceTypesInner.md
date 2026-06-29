# ListInstanceTypes200ResponseAllOfInstanceTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**ListInstanceTypes200ResponseAllOfInstanceTypesInnerAccount**](ListInstanceTypes200ResponseAllOfInstanceTypesInnerAccount.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ProvisionTypeCode** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**EnvironmentPrefix** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Featured** | Pointer to **bool** |  | [optional] 
**Versions** | Pointer to **[]string** |  | [optional] 
**InstanceTypeLayouts** | Pointer to [**[]ListInstanceTypes200ResponseAllOfInstanceTypesInnerInstanceTypeLayoutsInner**](ListInstanceTypes200ResponseAllOfInstanceTypesInnerInstanceTypeLayoutsInner.md) |  | [optional] 
**ImagePath** | Pointer to **NullableString** | Logo image URL | [optional] 
**DarkImagePath** | Pointer to **NullableString** | Dark logo image URL | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListInstanceTypes200ResponseAllOfInstanceTypesInner{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ProvisionTypeCode (Nullable)

Use the Nullable wrapper methods:
- `obj.ProvisionTypeCode.IsSet()` — check if set
- `obj.ProvisionTypeCode.Get()` — get the inner value (returns pointer)
- `obj.ProvisionTypeCode.Set(&val)` — set the value
- `obj.ProvisionTypeCode.Unset()` — clear the value
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


