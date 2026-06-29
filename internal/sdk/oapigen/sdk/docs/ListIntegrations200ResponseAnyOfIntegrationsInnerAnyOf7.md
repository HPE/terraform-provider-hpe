# ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf7

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**IntegrationType** | Pointer to [**ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf7IntegrationType**](ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf7IntegrationType.md) |  | [optional] 
**Url** | Pointer to **string** |  | [optional] 
**ServiceKey** | Pointer to [**NullableListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf7ServiceKey**](ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf7ServiceKey.md) |  | [optional] 
**IsPlugin** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf7Config**](ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf7Config.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**LastSync** | Pointer to **NullableString** |  | [optional] 
**LastSyncDuration** | Pointer to **NullableString** |  | [optional] 
**Credential** | Pointer to [**ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf7Credential**](ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf7Credential.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf7{
    // Set fields directly
}
```

### ServiceKey (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceKey.IsSet()` — check if set
- `obj.ServiceKey.Get()` — get the inner value (returns pointer)
- `obj.ServiceKey.Set(&val)` — set the value
- `obj.ServiceKey.Unset()` — clear the value
### StatusDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusDate.IsSet()` — check if set
- `obj.StatusDate.Get()` — get the inner value (returns pointer)
- `obj.StatusDate.Set(&val)` — set the value
- `obj.StatusDate.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### LastSync (Nullable)

Use the Nullable wrapper methods:
- `obj.LastSync.IsSet()` — check if set
- `obj.LastSync.Get()` — get the inner value (returns pointer)
- `obj.LastSync.Set(&val)` — set the value
- `obj.LastSync.Unset()` — clear the value
### LastSyncDuration (Nullable)

Use the Nullable wrapper methods:
- `obj.LastSyncDuration.IsSet()` — check if set
- `obj.LastSyncDuration.Get()` — get the inner value (returns pointer)
- `obj.LastSyncDuration.Set(&val)` — set the value
- `obj.LastSyncDuration.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


