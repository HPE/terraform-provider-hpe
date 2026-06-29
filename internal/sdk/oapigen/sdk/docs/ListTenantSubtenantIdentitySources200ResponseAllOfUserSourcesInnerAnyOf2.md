# ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf2

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Deleted** | Pointer to **bool** |  | [optional] 
**AutoSyncOnLogin** | Pointer to **bool** |  | [optional] 
**ExternalLogin** | Pointer to **bool** |  | [optional] 
**AllowCustomMappings** | Pointer to **bool** |  | [optional] 
**ManualRoleAssignment** | Pointer to **bool** |  | [optional] 
**Account** | Pointer to [**ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf2Account**](ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf2Account.md) |  | [optional] 
**DefaultAccountRole** | Pointer to [**ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf2DefaultAccountRole**](ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf2DefaultAccountRole.md) |  | [optional] 
**Config** | Pointer to [**ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf2Config**](ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf2Config.md) |  | [optional] 
**RoleMappings** | Pointer to [**[]ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf2RoleMappingsInner**](ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf2RoleMappingsInner.md) |  | [optional] 
**Subdomain** | Pointer to **string** |  | [optional] 
**LoginURL** | Pointer to **string** |  | [optional] 
**ProviderSettings** | Pointer to **map[string]interface{}** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListTenantSubtenantIdentitySources200ResponseAllOfUserSourcesInnerAnyOf2{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


