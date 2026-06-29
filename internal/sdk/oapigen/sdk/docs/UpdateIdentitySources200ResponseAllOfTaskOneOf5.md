# UpdateIdentitySources200ResponseAllOfTaskOneOf5

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
**Account** | Pointer to [**UpdateIdentitySources200ResponseAllOfTaskOneOf5Account**](UpdateIdentitySources200ResponseAllOfTaskOneOf5Account.md) |  | [optional] 
**DefaultAccountRole** | Pointer to [**UpdateIdentitySources200ResponseAllOfTaskOneOf5DefaultAccountRole**](UpdateIdentitySources200ResponseAllOfTaskOneOf5DefaultAccountRole.md) |  | [optional] 
**Config** | Pointer to [**UpdateIdentitySources200ResponseAllOfTaskOneOf5Config**](UpdateIdentitySources200ResponseAllOfTaskOneOf5Config.md) |  | [optional] 
**RoleMappings** | Pointer to [**[]UpdateIdentitySources200ResponseAllOfTaskOneOf5RoleMappingsInner**](UpdateIdentitySources200ResponseAllOfTaskOneOf5RoleMappingsInner.md) |  | [optional] 
**Subdomain** | Pointer to **string** |  | [optional] 
**LoginURL** | Pointer to **string** |  | [optional] 
**ProviderSettings** | Pointer to [**UpdateIdentitySources200ResponseAllOfTaskOneOf5ProviderSettings**](UpdateIdentitySources200ResponseAllOfTaskOneOf5ProviderSettings.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateIdentitySources200ResponseAllOfTaskOneOf5{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


