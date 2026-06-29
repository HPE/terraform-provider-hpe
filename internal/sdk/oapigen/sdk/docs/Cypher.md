# Cypher

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int32** |  | [optional] 
**ItemKey** | Pointer to **string** |  | [optional] 
**LeaseTimeout** | Pointer to **int64** |  | [optional] 
**ExpireDate** | Pointer to **NullableTime** |  | [optional] 
**DateCreated** | Pointer to **NullableTime** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**LastAccessed** | Pointer to **time.Time** |  | [optional] 
**CreatedBy** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &Cypher{
    // Set fields directly
}
```

### ExpireDate (Nullable)

Use the Nullable wrapper methods:
- `obj.ExpireDate.IsSet()` — check if set
- `obj.ExpireDate.Get()` — get the inner value (returns pointer)
- `obj.ExpireDate.Set(&val)` — set the value
- `obj.ExpireDate.Unset()` — clear the value
### DateCreated (Nullable)

Use the Nullable wrapper methods:
- `obj.DateCreated.IsSet()` — check if set
- `obj.DateCreated.Get()` — get the inner value (returns pointer)
- `obj.DateCreated.Set(&val)` — set the value
- `obj.DateCreated.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


