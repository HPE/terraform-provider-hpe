# InfobloxNetworkPoolServerConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InventoryExisting** | Pointer to **string** | Inventory Existing | [optional] [default to "off"]
**ExtraAttributes** | Pointer to **NullableString** | Extra Attributes | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InfobloxNetworkPoolServerConfig{
    // Set fields directly
}
```

### ExtraAttributes (Nullable)

Use the Nullable wrapper methods:
- `obj.ExtraAttributes.IsSet()` — check if set
- `obj.ExtraAttributes.Get()` — get the inner value (returns pointer)
- `obj.ExtraAttributes.Set(&val)` — set the value
- `obj.ExtraAttributes.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


