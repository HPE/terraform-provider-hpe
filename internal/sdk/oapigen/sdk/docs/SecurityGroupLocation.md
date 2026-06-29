# SecurityGroupLocation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**IacId** | Pointer to **NullableString** |  | [optional] 
**Zone** | Pointer to [**SecurityGroupLocationZone**](SecurityGroupLocationZone.md) |  | [optional] 
**ZonePool** | Pointer to [**SecurityGroupLocationZonePool**](SecurityGroupLocationZonePool.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Priority** | Pointer to **NullableString** |  | [optional] 
**GroupLayer** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SecurityGroupLocation{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### IacId (Nullable)

Use the Nullable wrapper methods:
- `obj.IacId.IsSet()` — check if set
- `obj.IacId.Get()` — get the inner value (returns pointer)
- `obj.IacId.Set(&val)` — set the value
- `obj.IacId.Unset()` — clear the value
### Priority (Nullable)

Use the Nullable wrapper methods:
- `obj.Priority.IsSet()` — check if set
- `obj.Priority.Get()` — get the inner value (returns pointer)
- `obj.Priority.Set(&val)` — set the value
- `obj.Priority.Unset()` — clear the value
### GroupLayer (Nullable)

Use the Nullable wrapper methods:
- `obj.GroupLayer.IsSet()` — check if set
- `obj.GroupLayer.Get()` — get the inner value (returns pointer)
- `obj.GroupLayer.Set(&val)` — set the value
- `obj.GroupLayer.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


