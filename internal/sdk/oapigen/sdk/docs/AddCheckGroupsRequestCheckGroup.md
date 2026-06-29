# AddCheckGroupsRequestCheckGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Unique name scoped to your account for the check group | 
**Description** | Pointer to **string** | Optional description field | [optional] 
**MinHappy** | Pointer to **int32** | This specifies the minimum number of checks within the group that must be happy to keep the group from becoming unhealthy. | [optional] [default to 1]
**InUptime** | Pointer to **bool** | Used to determine if check should affect account wide availability calculations | [optional] [default to true]
**Severity** | Pointer to **string** | Determines the maximum severity level this group can incur on an incident when failing | [optional] [default to "critical"]
**Active** | Pointer to **bool** | Used to determine if check group is active | [optional] [default to true]
**Checks** | Pointer to **[]int32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCheckGroupsRequestCheckGroup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


