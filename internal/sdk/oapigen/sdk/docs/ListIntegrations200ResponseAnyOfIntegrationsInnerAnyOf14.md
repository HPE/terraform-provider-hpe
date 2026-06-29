# ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf14

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**IntegrationType** | Pointer to [**ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf14IntegrationType**](ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf14IntegrationType.md) |  | [optional] 
**Url** | Pointer to **string** |  | [optional] 
**Port** | Pointer to **string** |  | [optional] 
**Username** | Pointer to **string** |  | [optional] 
**Password** | Pointer to **string** |  | [optional] 
**PasswordHash** | Pointer to **string** |  | [optional] 
**Path** | Pointer to **string** |  | [optional] 
**Version** | Pointer to **string** |  | [optional] 
**WindowsVersion** | Pointer to **string** |  | [optional] 
**RepoUrl** | Pointer to **string** |  | [optional] 
**ServiceMode** | Pointer to **string** |  | [optional] 
**IsPlugin** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf14Config**](ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf14Config.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusDate** | Pointer to **time.Time** |  | [optional] 
**StatusMessage** | Pointer to **string** |  | [optional] 
**LastSync** | Pointer to **NullableString** |  | [optional] 
**LastSyncDuration** | Pointer to **NullableString** |  | [optional] 
**Credential** | Pointer to [**ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf14Credential**](ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf14Credential.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf14{
    // Set fields directly
}
```

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


