# ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf

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
**Account** | Pointer to [**ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOfAccount**](ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOfAccount.md) |  | [optional] 
**DefaultAccountRole** | Pointer to [**ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOfDefaultAccountRole**](ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOfDefaultAccountRole.md) |  | [optional] 
**Config** | Pointer to [**ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOfConfig**](ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOfConfig.md) |  | [optional] 
**RoleMappings** | Pointer to [**[]ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOfRoleMappingsInner**](ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOfRoleMappingsInner.md) |  | [optional] 
**Subdomain** | Pointer to **string** |  | [optional] 
**LoginURL** | Pointer to **string** |  | [optional] 
**ProviderSettings** | Pointer to **map[string]interface{}** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


