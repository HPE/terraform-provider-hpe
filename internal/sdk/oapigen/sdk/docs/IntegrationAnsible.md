# IntegrationAnsible

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**IntegrationType** | Pointer to [**AddIntegrations200ResponseAllOfIntegrationOneOfIntegrationType**](AddIntegrations200ResponseAllOfIntegrationOneOfIntegrationType.md) |  | [optional] 
**Url** | Pointer to **string** |  | [optional] 
**ServiceKey** | Pointer to [**AddIntegrations200ResponseAllOfIntegrationOneOfServiceKey**](AddIntegrations200ResponseAllOfIntegrationOneOfServiceKey.md) |  | [optional] 
**IsPlugin** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**AddIntegrations200ResponseAllOfIntegrationOneOfConfig**](AddIntegrations200ResponseAllOfIntegrationOneOfConfig.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusDate** | Pointer to **time.Time** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**LastSync** | Pointer to **NullableString** |  | [optional] 
**LastSyncDuration** | Pointer to **NullableInt64** |  | [optional] 
**Credential** | Pointer to [**AddIntegrations200ResponseAllOfIntegrationOneOfCredential**](AddIntegrations200ResponseAllOfIntegrationOneOfCredential.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &IntegrationAnsible{
    // Set fields directly
}
```

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


