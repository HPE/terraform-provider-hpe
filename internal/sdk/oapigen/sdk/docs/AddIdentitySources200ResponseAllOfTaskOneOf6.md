# AddIdentitySources200ResponseAllOfTaskOneOf6

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Deleted** | Pointer to **bool** |  | [optional] 
**AutoSyncOnLogin** | Pointer to **bool** |  | [optional] 
**ExternalLogin** | Pointer to **bool** |  | [optional] 
**AllowCustomMappings** | Pointer to **bool** |  | [optional] 
**ManualRoleAssignment** | Pointer to **bool** |  | [optional] 
**Account** | Pointer to [**AddIdentitySources200ResponseAllOfTaskOneOf6Account**](AddIdentitySources200ResponseAllOfTaskOneOf6Account.md) |  | [optional] 
**DefaultAccountRole** | Pointer to [**AddIdentitySources200ResponseAllOfTaskOneOf6DefaultAccountRole**](AddIdentitySources200ResponseAllOfTaskOneOf6DefaultAccountRole.md) |  | [optional] 
**Config** | Pointer to [**AddIdentitySources200ResponseAllOfTaskOneOf6Config**](AddIdentitySources200ResponseAllOfTaskOneOf6Config.md) |  | [optional] 
**RoleMappings** | Pointer to [**[]AddIdentitySources200ResponseAllOfTaskOneOf6RoleMappingsInner**](AddIdentitySources200ResponseAllOfTaskOneOf6RoleMappingsInner.md) |  | [optional] 
**Subdomain** | Pointer to **string** |  | [optional] 
**LoginURL** | Pointer to **string** |  | [optional] 
**ProviderSettings** | Pointer to [**AddIdentitySources200ResponseAllOfTaskOneOf6ProviderSettings**](AddIdentitySources200ResponseAllOfTaskOneOf6ProviderSettings.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddIdentitySources200ResponseAllOfTaskOneOf6{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


