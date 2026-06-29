# ApprovalItemApprovalItemDetailsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Category** | Pointer to **string** | The category of the change (e.g. network, volume, plan, planMemory, planCores) | [optional] 
**Type** | Pointer to **string** | The type of change (e.g. increase, decrease, add, remove, change) | [optional] 
**Name** | Pointer to **string** | The name of the item being changed | [optional] 
**FromValue** | Pointer to **NullableString** | The original value before the change | [optional] 
**ToValue** | Pointer to **NullableString** | The new value after the change | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ApprovalItemApprovalItemDetailsInner{
    // Set fields directly
}
```

### FromValue (Nullable)

Use the Nullable wrapper methods:
- `obj.FromValue.IsSet()` — check if set
- `obj.FromValue.Get()` — get the inner value (returns pointer)
- `obj.FromValue.Set(&val)` — set the value
- `obj.FromValue.Unset()` — clear the value
### ToValue (Nullable)

Use the Nullable wrapper methods:
- `obj.ToValue.IsSet()` — check if set
- `obj.ToValue.Get()` — get the inner value (returns pointer)
- `obj.ToValue.Set(&val)` — set the value
- `obj.ToValue.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


