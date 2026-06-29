# CatalogItemTypeInstanceScribePortsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Port** | Pointer to **int64** | Port number. | [optional] 
**Name** | Pointer to **string** | A name for the port. | [optional] 
**Lb** | Pointer to **NullableString** | The load balancer protocol. HTTP, HTTPS, or TCP. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CatalogItemTypeInstanceScribePortsInner{
    // Set fields directly
}
```

### Lb (Nullable)

Use the Nullable wrapper methods:
- `obj.Lb.IsSet()` — check if set
- `obj.Lb.Get()` — get the inner value (returns pointer)
- `obj.Lb.Set(&val)` — set the value
- `obj.Lb.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


