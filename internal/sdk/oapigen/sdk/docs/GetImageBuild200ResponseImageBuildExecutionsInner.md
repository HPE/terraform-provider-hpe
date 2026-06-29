# GetImageBuild200ResponseImageBuildExecutionsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**ImageBuild** | Pointer to [**NullableGetImageBuild200ResponseImageBuildExecutionsInnerImageBuild**](GetImageBuild200ResponseImageBuildExecutionsInnerImageBuild.md) |  | [optional] 
**BuildNumber** | Pointer to **int64** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **NullableTime** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**StatusPercent** | Pointer to **int64** |  | [optional] 
**StatusEta** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**ErrorMessage** | Pointer to **NullableString** |  | [optional] 
**CreatedBy** | Pointer to [**GetImageBuild200ResponseImageBuildExecutionsInnerCreatedBy**](GetImageBuild200ResponseImageBuildExecutionsInnerCreatedBy.md) |  | [optional] 
**TempInstance** | Pointer to **NullableString** |  | [optional] 
**VirtualImages** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetImageBuild200ResponseImageBuildExecutionsInner{
    // Set fields directly
}
```

### ImageBuild (Nullable)

Use the Nullable wrapper methods:
- `obj.ImageBuild.IsSet()` — check if set
- `obj.ImageBuild.Get()` — get the inner value (returns pointer)
- `obj.ImageBuild.Set(&val)` — set the value
- `obj.ImageBuild.Unset()` — clear the value
### EndDate (Nullable)

Use the Nullable wrapper methods:
- `obj.EndDate.IsSet()` — check if set
- `obj.EndDate.Get()` — get the inner value (returns pointer)
- `obj.EndDate.Set(&val)` — set the value
- `obj.EndDate.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### StatusEta (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusEta.IsSet()` — check if set
- `obj.StatusEta.Get()` — get the inner value (returns pointer)
- `obj.StatusEta.Set(&val)` — set the value
- `obj.StatusEta.Unset()` — clear the value
### ErrorMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorMessage.IsSet()` — check if set
- `obj.ErrorMessage.Get()` — get the inner value (returns pointer)
- `obj.ErrorMessage.Set(&val)` — set the value
- `obj.ErrorMessage.Unset()` — clear the value
### TempInstance (Nullable)

Use the Nullable wrapper methods:
- `obj.TempInstance.IsSet()` — check if set
- `obj.TempInstance.Get()` — get the inner value (returns pointer)
- `obj.TempInstance.Set(&val)` — set the value
- `obj.TempInstance.Unset()` — clear the value
### VirtualImages (Nullable)

Use the Nullable wrapper methods:
- `obj.VirtualImages.IsSet()` — check if set
- `obj.VirtualImages.Get()` — get the inner value (returns pointer)
- `obj.VirtualImages.Set(&val)` — set the value
- `obj.VirtualImages.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


