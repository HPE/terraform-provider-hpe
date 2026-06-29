# ListIntegrations200ResponseAnyOfIntegrationsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**IntegrationType** | Pointer to [**ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf16IntegrationType**](ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf16IntegrationType.md) |  | [optional] 
**Url** | Pointer to **string** |  | [optional] 
**ServiceKey** | Pointer to [**NullableListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf8ServiceKey**](ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf8ServiceKey.md) |  | [optional] 
**IsPlugin** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusDate** | Pointer to **time.Time** |  | [optional] 
**StatusMessage** | Pointer to **string** |  | [optional] 
**LastSync** | Pointer to **string** |  | [optional] 
**LastSyncDuration** | Pointer to **string** |  | [optional] 
**Credential** | Pointer to [**ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf16Credential**](ListIntegrations200ResponseAnyOfIntegrationsInnerAnyOf16Credential.md) |  | [optional] 
**Version** | Pointer to **string** |  | [optional] 
**Host** | Pointer to **string** |  | [optional] 
**Username** | Pointer to **string** |  | [optional] 
**Password** | Pointer to **string** |  | [optional] 
**PasswordHash** | Pointer to **string** |  | [optional] 
**Token** | Pointer to **string** |  | [optional] 
**TokenHash** | Pointer to **string** |  | [optional] 
**ServiceFlag** | Pointer to **bool** |  | [optional] 
**Port** | Pointer to **string** |  | [optional] 
**Path** | Pointer to **string** |  | [optional] 
**WindowsVersion** | Pointer to **string** |  | [optional] 
**RepoUrl** | Pointer to **string** |  | [optional] 
**ServiceMode** | Pointer to **string** |  | [optional] 
**AuthType** | Pointer to **string** |  | [optional] 
**AuthId** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListIntegrations200ResponseAnyOfIntegrationsInner{
    // Set fields directly
}
```

### ServiceKey (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceKey.IsSet()` — check if set
- `obj.ServiceKey.Get()` — get the inner value (returns pointer)
- `obj.ServiceKey.Set(&val)` — set the value
- `obj.ServiceKey.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


