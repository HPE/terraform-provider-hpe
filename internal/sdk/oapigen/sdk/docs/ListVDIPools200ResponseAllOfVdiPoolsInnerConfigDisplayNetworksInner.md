# ListVDIPools200ResponseAllOfVdiPoolsInnerConfigDisplayNetworksInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**IpMode** | Pointer to **string** |  | [optional] 
**Pool** | Pointer to [**NullableListVDIPools200ResponseAllOfVdiPoolsInnerConfigDisplayNetworksInnerPool**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigDisplayNetworksInnerPool.md) |  | [optional] 
**UsePool** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListVDIPools200ResponseAllOfVdiPoolsInnerConfigDisplayNetworksInner{
    // Set fields directly
}
```

### Pool (Nullable)

Use the Nullable wrapper methods:
- `obj.Pool.IsSet()` — check if set
- `obj.Pool.Get()` — get the inner value (returns pointer)
- `obj.Pool.Set(&val)` — set the value
- `obj.Pool.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


