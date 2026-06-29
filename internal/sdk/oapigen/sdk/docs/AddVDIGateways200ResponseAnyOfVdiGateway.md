# AddVDIGateways200ResponseAnyOfVdiGateway

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**GatewayUrl** | Pointer to **NullableString** |  | [optional] 
**ApiKey** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddVDIGateways200ResponseAnyOfVdiGateway{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### GatewayUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.GatewayUrl.IsSet()` — check if set
- `obj.GatewayUrl.Get()` — get the inner value (returns pointer)
- `obj.GatewayUrl.Set(&val)` — set the value
- `obj.GatewayUrl.Unset()` — clear the value
### ApiKey (Nullable)

Use the Nullable wrapper methods:
- `obj.ApiKey.IsSet()` — check if set
- `obj.ApiKey.Get()` — get the inner value (returns pointer)
- `obj.ApiKey.Set(&val)` — set the value
- `obj.ApiKey.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


