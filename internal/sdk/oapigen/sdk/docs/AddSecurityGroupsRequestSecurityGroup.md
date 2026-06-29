# AddSecurityGroupsRequestSecurityGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name for your security group | 
**Description** | Pointer to **string** | Optional description field | [optional] 
**ZoneId** | **int64** | Scoped Cloud ID | 
**Active** | Pointer to **bool** | Set to &#x60;false&#x60; to disable a security group. | [optional] 
**Visibility** | Pointer to **string** | Visibility for the security group. | [optional] [default to "private"]
**NetworkServerId** | Pointer to **int64** | Network Server ID to scope the security group to a network integration (e.g. NSX-T). Use as an alternative to zoneId for network-server-scoped groups. | [optional] 
**CustomOptions** | Pointer to [**AddSecurityGroupsRequestSecurityGroupCustomOptions**](AddSecurityGroupsRequestSecurityGroupCustomOptions.md) |  | [optional] 
**TenantPermissions** | Pointer to [**AddSecurityGroupsRequestSecurityGroupTenantPermissions**](AddSecurityGroupsRequestSecurityGroupTenantPermissions.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**AddSecurityGroupsRequestSecurityGroupResourcePermissions**](AddSecurityGroupsRequestSecurityGroupResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddSecurityGroupsRequestSecurityGroup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


