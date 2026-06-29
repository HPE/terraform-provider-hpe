# ListVdi200ResponseAllOfDesktopsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Logo** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**AllocationStatus** | Pointer to **string** |  | [optional] 
**Allocation** | Pointer to [**ListVdi200ResponseAllOfDesktopsInnerAllocation**](ListVdi200ResponseAllOfDesktopsInnerAllocation.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListVdi200ResponseAllOfDesktopsInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


