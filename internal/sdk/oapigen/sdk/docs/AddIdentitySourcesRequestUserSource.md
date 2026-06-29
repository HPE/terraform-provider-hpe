# AddIdentitySourcesRequestUserSource

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Account** | Pointer to [**AddIdentitySourcesRequestUserSourceAccount**](AddIdentitySourcesRequestUserSourceAccount.md) |  | [optional] 
**Name** | **string** | A name for the Identity Source | 
**Type** | **string** | IDM type code | 
**Description** | Pointer to **string** | description | [optional] 
**DefaultAccountRole** | Pointer to [**AddIdentitySourcesRequestUserSourceDefaultAccountRole**](AddIdentitySourcesRequestUserSourceDefaultAccountRole.md) |  | [optional] 
**RoleMappings** | Pointer to [**AddIdentitySourcesRequestUserSourceRoleMappings**](AddIdentitySourcesRequestUserSourceRoleMappings.md) |  | [optional] 
**RoleMappingNames** | Pointer to **map[string]string** | Map of Morpheus &#39;&#x60;Role ID&#x60;&#39;:&#39;&#x60;Role Name&#x60;&#39;.  | [optional] 
**AllowCustomMappings** | Pointer to **bool** | Enable Role Mapping Permission | [optional] 
**ManualRoleAssignment** | Pointer to **bool** | Manual Role Assignment | [optional] 
**Config** | Pointer to [**AddIdentitySourcesRequestUserSourceConfig**](AddIdentitySourcesRequestUserSourceConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddIdentitySourcesRequestUserSource{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


