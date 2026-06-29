# GetNetworkPoolServer200ResponseNetworkPoolServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Network Pool Server ID | [optional] 
**Type** | Pointer to [**GetNetworkPoolServer200ResponseNetworkPoolServerType**](GetNetworkPoolServer200ResponseNetworkPoolServerType.md) |  | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**ServiceUrl** | Pointer to **NullableString** | Service URL | [optional] 
**ServiceHost** | Pointer to **NullableString** | Service Host | [optional] 
**ServicePort** | Pointer to **NullableInt32** | Service Port | [optional] 
**ServiceMode** | Pointer to **NullableString** | Service Mode | [optional] 
**ServiceUsername** | Pointer to **NullableString** | Service Username | [optional] 
**ServicePassword** | Pointer to **NullableString** | Service Password | [optional] 
**ServicePasswordHash** | Pointer to **string** |  | [optional] 
**ServiceThrottleRate** | Pointer to **NullableInt64** | Throttle Rate | [optional] [default to 0]
**IgnoreSsl** | Pointer to **NullableBool** | Disable SSL SNI Verification | [optional] [default to true]
**Status** | Pointer to **string** | Status | [optional] 
**StatusMessage** | Pointer to **NullableString** | Status Message | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** | Config object varies with pool server type. | [optional] 
**NetworkFilter** | Pointer to **NullableString** | Network Filter | [optional] 
**ZoneFilter** | Pointer to **NullableString** | Zone Filter | [optional] 
**TenantMatch** | Pointer to **NullableString** | Tenant Match | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Account** | Pointer to [**GetNetworkPoolServer200ResponseNetworkPoolServerAccount**](GetNetworkPoolServer200ResponseNetworkPoolServerAccount.md) |  | [optional] 
**Integration** | Pointer to [**GetNetworkPoolServer200ResponseNetworkPoolServerIntegration**](GetNetworkPoolServer200ResponseNetworkPoolServerIntegration.md) |  | [optional] 
**Pools** | Pointer to [**[]GetNetworkPoolServer200ResponseNetworkPoolServerPoolsInner**](GetNetworkPoolServer200ResponseNetworkPoolServerPoolsInner.md) |  | [optional] 
**Credential** | Pointer to [**GetNetworkPoolServer200ResponseNetworkPoolServerCredential**](GetNetworkPoolServer200ResponseNetworkPoolServerCredential.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkPoolServer200ResponseNetworkPoolServer{
    // Set fields directly
}
```

### ServiceUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceUrl.IsSet()` — check if set
- `obj.ServiceUrl.Get()` — get the inner value (returns pointer)
- `obj.ServiceUrl.Set(&val)` — set the value
- `obj.ServiceUrl.Unset()` — clear the value
### ServiceHost (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceHost.IsSet()` — check if set
- `obj.ServiceHost.Get()` — get the inner value (returns pointer)
- `obj.ServiceHost.Set(&val)` — set the value
- `obj.ServiceHost.Unset()` — clear the value
### ServicePort (Nullable)

Use the Nullable wrapper methods:
- `obj.ServicePort.IsSet()` — check if set
- `obj.ServicePort.Get()` — get the inner value (returns pointer)
- `obj.ServicePort.Set(&val)` — set the value
- `obj.ServicePort.Unset()` — clear the value
### ServiceMode (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceMode.IsSet()` — check if set
- `obj.ServiceMode.Get()` — get the inner value (returns pointer)
- `obj.ServiceMode.Set(&val)` — set the value
- `obj.ServiceMode.Unset()` — clear the value
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
### IgnoreSsl (Nullable)

Use the Nullable wrapper methods:
- `obj.IgnoreSsl.IsSet()` — check if set
- `obj.IgnoreSsl.Get()` — get the inner value (returns pointer)
- `obj.IgnoreSsl.Set(&val)` — set the value
- `obj.IgnoreSsl.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### StatusDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusDate.IsSet()` — check if set
- `obj.StatusDate.Get()` — get the inner value (returns pointer)
- `obj.StatusDate.Set(&val)` — set the value
- `obj.StatusDate.Unset()` — clear the value
### NetworkFilter (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkFilter.IsSet()` — check if set
- `obj.NetworkFilter.Get()` — get the inner value (returns pointer)
- `obj.NetworkFilter.Set(&val)` — set the value
- `obj.NetworkFilter.Unset()` — clear the value
### ZoneFilter (Nullable)

Use the Nullable wrapper methods:
- `obj.ZoneFilter.IsSet()` — check if set
- `obj.ZoneFilter.Get()` — get the inner value (returns pointer)
- `obj.ZoneFilter.Set(&val)` — set the value
- `obj.ZoneFilter.Unset()` — clear the value
### TenantMatch (Nullable)

Use the Nullable wrapper methods:
- `obj.TenantMatch.IsSet()` — check if set
- `obj.TenantMatch.Get()` — get the inner value (returns pointer)
- `obj.TenantMatch.Set(&val)` — set the value
- `obj.TenantMatch.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


