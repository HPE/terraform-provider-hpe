# ListClusterLayouts200ResponseAllOfLayoutsInnerSpecTemplatesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Type** | Pointer to [**ListClusterLayouts200ResponseAllOfLayoutsInnerSpecTemplatesInnerType**](ListClusterLayouts200ResponseAllOfLayoutsInnerSpecTemplatesInnerType.md) |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalType** | Pointer to **NullableString** |  | [optional] 
**DeploymentId** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **NullableString** |  | [optional] 
**File** | Pointer to [**ListClusterLayouts200ResponseAllOfLayoutsInnerSpecTemplatesInnerFile**](ListClusterLayouts200ResponseAllOfLayoutsInnerSpecTemplatesInnerFile.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**CreatedBy** | Pointer to **NullableString** |  | [optional] 
**UpdatedBy** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListClusterLayouts200ResponseAllOfLayoutsInnerSpecTemplatesInner{
    // Set fields directly
}
```

### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
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
### CreatedBy (Nullable)

Use the Nullable wrapper methods:
- `obj.CreatedBy.IsSet()` — check if set
- `obj.CreatedBy.Get()` — get the inner value (returns pointer)
- `obj.CreatedBy.Set(&val)` — set the value
- `obj.CreatedBy.Unset()` — clear the value
### UpdatedBy (Nullable)

Use the Nullable wrapper methods:
- `obj.UpdatedBy.IsSet()` — check if set
- `obj.UpdatedBy.Get()` — get the inner value (returns pointer)
- `obj.UpdatedBy.Set(&val)` — set the value
- `obj.UpdatedBy.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


