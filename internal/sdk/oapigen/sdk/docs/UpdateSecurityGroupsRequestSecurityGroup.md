# UpdateSecurityGroupsRequestSecurityGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name for your security group | [optional] 
**Description** | Pointer to **string** | Optional description field | [optional] 
**Active** | Pointer to **bool** | Set to &#x60;false&#x60; to disable a security group. | [optional] 
**Visibility** | Pointer to **string** | Visibility for the security group. | [optional] [default to "private"]
**TenantPermissions** | Pointer to [**UpdateSecurityGroupsRequestSecurityGroupTenantPermissions**](UpdateSecurityGroupsRequestSecurityGroupTenantPermissions.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**UpdateSecurityGroupsRequestSecurityGroupResourcePermissions**](UpdateSecurityGroupsRequestSecurityGroupResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateSecurityGroupsRequestSecurityGroup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


