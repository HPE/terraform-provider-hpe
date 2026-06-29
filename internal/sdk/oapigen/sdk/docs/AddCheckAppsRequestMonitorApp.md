# AddCheckAppsRequestMonitorApp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Unique name scoped to your account for the check app | 
**Description** | Pointer to **string** | Optional description field | [optional] 
**InUptime** | Pointer to **bool** | Used to determine if check should affect account wide availability calculations | [optional] [default to true]
**Severity** | Pointer to **string** | Severity level of incidents that are created when this check fails | [optional] [default to "critical"]
**Active** | Pointer to **bool** | Used to determine if check app is active | [optional] [default to true]
**Checks** | Pointer to **[]int32** |  | [optional] 
**CheckGroups** | Pointer to **[]int32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCheckAppsRequestMonitorApp{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


