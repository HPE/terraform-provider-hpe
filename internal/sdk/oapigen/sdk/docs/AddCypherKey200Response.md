# AddCypherKey200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**LeaseDuration** | Pointer to **NullableInt64** |  | [optional] 
**Cypher** | Pointer to [**AddCypherKey200ResponseAllOfCypher**](AddCypherKey200ResponseAllOfCypher.md) |  | [optional] 
**Success** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCypherKey200Response{
    // Set fields directly
}
```

### LeaseDuration (Nullable)

Use the Nullable wrapper methods:
- `obj.LeaseDuration.IsSet()` — check if set
- `obj.LeaseDuration.Get()` — get the inner value (returns pointer)
- `obj.LeaseDuration.Set(&val)` — set the value
- `obj.LeaseDuration.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


