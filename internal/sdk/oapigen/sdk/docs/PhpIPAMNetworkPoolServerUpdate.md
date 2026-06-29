# PhpIPAMNetworkPoolServerUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**Enabled** | Pointer to **bool** | Can be used to enable / disable the network pool server. | [optional] [default to true]
**ServiceUrl** | Pointer to **NullableString** | URL | [optional] 
**ServiceUsername** | Pointer to **NullableString** | Username | [optional] 
**ServicePassword** | Pointer to **NullableString** | Password | [optional] 
**ServiceThrottleRate** | Pointer to **NullableInt64** | Throttle Rate | [optional] [default to 0]
**IgnoreSsl** | Pointer to **bool** | Disable SSL SNI Verification | [optional] 
**NetworkFilter** | Pointer to **NullableString** | Network Filter | [optional] 
**Config** | Pointer to [**PhpIPAMNetworkPoolServerUpdateConfig**](PhpIPAMNetworkPoolServerUpdateConfig.md) |  | [optional] 
**Credential** | Pointer to [**PhpIPAMNetworkPoolServerUpdateCredential**](PhpIPAMNetworkPoolServerUpdateCredential.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &PhpIPAMNetworkPoolServerUpdate{
    // Set fields directly
}
```

### ServiceUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceUrl.IsSet()` — check if set
- `obj.ServiceUrl.Get()` — get the inner value (returns pointer)
- `obj.ServiceUrl.Set(&val)` — set the value
- `obj.ServiceUrl.Unset()` — clear the value
### ServiceUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceUsername.IsSet()` — check if set
- `obj.ServiceUsername.Get()` — get the inner value (returns pointer)
- `obj.ServiceUsername.Set(&val)` — set the value
- `obj.ServiceUsername.Unset()` — clear the value
### ServicePassword (Nullable)

Use the Nullable wrapper methods:
- `obj.ServicePassword.IsSet()` — check if set
- `obj.ServicePassword.Get()` — get the inner value (returns pointer)
- `obj.ServicePassword.Set(&val)` — set the value
- `obj.ServicePassword.Unset()` — clear the value
### ServiceThrottleRate (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceThrottleRate.IsSet()` — check if set
- `obj.ServiceThrottleRate.Get()` — get the inner value (returns pointer)
- `obj.ServiceThrottleRate.Set(&val)` — set the value
- `obj.ServiceThrottleRate.Unset()` — clear the value
### NetworkFilter (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkFilter.IsSet()` — check if set
- `obj.NetworkFilter.Get()` — get the inner value (returns pointer)
- `obj.NetworkFilter.Set(&val)` — set the value
- `obj.NetworkFilter.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


