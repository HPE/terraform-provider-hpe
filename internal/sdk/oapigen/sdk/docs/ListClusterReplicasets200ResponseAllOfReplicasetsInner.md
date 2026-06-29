# ListClusterReplicasets200ResponseAllOfReplicasetsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**ResourceLevel** | Pointer to **NullableString** |  | [optional] 
**ResourceType** | Pointer to **string** |  | [optional] 
**Managed** | Pointer to **NullableBool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**LastUpdated** | Pointer to **NullableTime** |  | [optional] 
**Owner** | Pointer to [**ListClusterReplicasets200ResponseAllOfReplicasetsInnerOwner**](ListClusterReplicasets200ResponseAllOfReplicasetsInnerOwner.md) |  | [optional] 
**TotalCpuUsage** | Pointer to **NullableInt64** |  | [optional] 
**Stats** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListClusterReplicasets200ResponseAllOfReplicasetsInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ResourceLevel (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourceLevel.IsSet()` — check if set
- `obj.ResourceLevel.Get()` — get the inner value (returns pointer)
- `obj.ResourceLevel.Set(&val)` — set the value
- `obj.ResourceLevel.Unset()` — clear the value
### Managed (Nullable)

Use the Nullable wrapper methods:
- `obj.Managed.IsSet()` — check if set
- `obj.Managed.Get()` — get the inner value (returns pointer)
- `obj.Managed.Set(&val)` — set the value
- `obj.Managed.Unset()` — clear the value
### LastUpdated (Nullable)

Use the Nullable wrapper methods:
- `obj.LastUpdated.IsSet()` — check if set
- `obj.LastUpdated.Get()` — get the inner value (returns pointer)
- `obj.LastUpdated.Set(&val)` — set the value
- `obj.LastUpdated.Unset()` — clear the value
### TotalCpuUsage (Nullable)

Use the Nullable wrapper methods:
- `obj.TotalCpuUsage.IsSet()` — check if set
- `obj.TotalCpuUsage.Get()` — get the inner value (returns pointer)
- `obj.TotalCpuUsage.Set(&val)` — set the value
- `obj.TotalCpuUsage.Unset()` — clear the value
### Stats (Nullable)

Use the Nullable wrapper methods:
- `obj.Stats.IsSet()` — check if set
- `obj.Stats.Get()` — get the inner value (returns pointer)
- `obj.Stats.Set(&val)` — set the value
- `obj.Stats.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


