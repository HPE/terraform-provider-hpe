# UpdateAlertsRequestAlert

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Unique name scoped to your account for the alert | [optional] 
**MinDuration** | Pointer to **int32** | Duration in minutes of the delay before sending notification(s) | [optional] [default to 0]
**MinSeverity** | Pointer to **string** | Severity level threshold for sending notifications. | [optional] [default to "critical"]
**Active** | Pointer to **bool** | Set to false to disable notifications | [optional] [default to true]
**AllChecks** | Pointer to **bool** | Trigger for all checks | [optional] [default to false]
**AllGroups** | Pointer to **bool** | Trigger for all check groups | [optional] [default to false]
**AllApps** | Pointer to **bool** | Trigger for all monitor apps | [optional] [default to false]
**Checks** | Pointer to **[]int32** |  | [optional] 
**Groups** | Pointer to **[]int32** |  | [optional] 
**Apps** | Pointer to **[]int32** |  | [optional] 
**Contacts** | Pointer to [**[]UpdateAlertsRequestAlertContactsInner**](UpdateAlertsRequestAlertContactsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateAlertsRequestAlert{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


