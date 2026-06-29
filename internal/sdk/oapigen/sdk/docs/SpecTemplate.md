# SpecTemplate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**GetSpecTemplate200ResponseSpecTemplateAccount**](GetSpecTemplate200ResponseSpecTemplateAccount.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**Type** | Pointer to [**GetSpecTemplate200ResponseSpecTemplateType**](GetSpecTemplate200ResponseSpecTemplateType.md) |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalType** | Pointer to **NullableString** |  | [optional] 
**DeploymentId** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **NullableString** |  | [optional] 
**File** | Pointer to [**GetSpecTemplate200ResponseSpecTemplateFile**](GetSpecTemplate200ResponseSpecTemplateFile.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**CreatedBy** | Pointer to **string** |  | [optional] 
**UpdatedBy** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SpecTemplate{
    // Set fields directly
}
```

### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### ExternalType (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalType.IsSet()` — check if set
- `obj.ExternalType.Get()` — get the inner value (returns pointer)
- `obj.ExternalType.Set(&val)` — set the value
- `obj.ExternalType.Unset()` — clear the value
### DeploymentId (Nullable)

Use the Nullable wrapper methods:
- `obj.DeploymentId.IsSet()` — check if set
- `obj.DeploymentId.Get()` — get the inner value (returns pointer)
- `obj.DeploymentId.Set(&val)` — set the value
- `obj.DeploymentId.Unset()` — clear the value
### Status (Nullable)

Use the Nullable wrapper methods:
- `obj.Status.IsSet()` — check if set
- `obj.Status.Get()` — get the inner value (returns pointer)
- `obj.Status.Set(&val)` — set the value
- `obj.Status.Unset()` — clear the value
### UpdatedBy (Nullable)

Use the Nullable wrapper methods:
- `obj.UpdatedBy.IsSet()` — check if set
- `obj.UpdatedBy.Get()` — get the inner value (returns pointer)
- `obj.UpdatedBy.Set(&val)` — set the value
- `obj.UpdatedBy.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


