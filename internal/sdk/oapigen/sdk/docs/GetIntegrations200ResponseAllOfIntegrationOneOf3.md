# GetIntegrations200ResponseAllOfIntegrationOneOf3

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**IntegrationType** | Pointer to [**GetIntegrations200ResponseAllOfIntegrationOneOf3IntegrationType**](GetIntegrations200ResponseAllOfIntegrationOneOf3IntegrationType.md) |  | [optional] 
**Url** | Pointer to **string** |  | [optional] 
**IsPlugin** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**GetIntegrations200ResponseAllOfIntegrationOneOf3Config**](GetIntegrations200ResponseAllOfIntegrationOneOf3Config.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**LastSync** | Pointer to **NullableString** |  | [optional] 
**LastSyncDuration** | Pointer to **NullableString** |  | [optional] 
**Credential** | Pointer to [**GetIntegrations200ResponseAllOfIntegrationOneOf3Credential**](GetIntegrations200ResponseAllOfIntegrationOneOf3Credential.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetIntegrations200ResponseAllOfIntegrationOneOf3{
    // Set fields directly
}
```

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


