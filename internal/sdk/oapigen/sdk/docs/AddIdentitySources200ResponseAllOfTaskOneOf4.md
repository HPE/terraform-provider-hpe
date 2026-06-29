# AddIdentitySources200ResponseAllOfTaskOneOf4

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
**Account** | Pointer to [**AddIdentitySources200ResponseAllOfTaskOneOf4Account**](AddIdentitySources200ResponseAllOfTaskOneOf4Account.md) |  | [optional] 
**DefaultAccountRole** | Pointer to [**AddIdentitySources200ResponseAllOfTaskOneOf4DefaultAccountRole**](AddIdentitySources200ResponseAllOfTaskOneOf4DefaultAccountRole.md) |  | [optional] 
**Config** | Pointer to [**AddIdentitySources200ResponseAllOfTaskOneOf4Config**](AddIdentitySources200ResponseAllOfTaskOneOf4Config.md) |  | [optional] 
**RoleMappings** | Pointer to [**[]AddIdentitySources200ResponseAllOfTaskOneOf4RoleMappingsInner**](AddIdentitySources200ResponseAllOfTaskOneOf4RoleMappingsInner.md) |  | [optional] 
**Subdomain** | Pointer to **string** |  | [optional] 
**LoginURL** | Pointer to **string** |  | [optional] 
**ProviderSettings** | Pointer to **map[string]interface{}** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddIdentitySources200ResponseAllOfTaskOneOf4{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


