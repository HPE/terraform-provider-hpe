# AcknowledgeHealthAlarmsRequestAlarm

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Acknowledged** | **bool** | Pass &#x60;true&#x60; to ackowledge an alarm, or pass &#x60;false&#x60; to unacknowledge it. | 
**Ids** | Pointer to **[]int64** | Array of Alarm ID(s)to be updated. | [optional] [default to []]
**All** | Pointer to **bool** | Pass &#x60;true&#x60; to update all alarms instead of passing ids. This will update any active alarm that is not already acknowledged.  | [optional] [default to false]

## Usage

Instantiate with a Go composite literal:

```go
obj := &AcknowledgeHealthAlarmsRequestAlarm{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


