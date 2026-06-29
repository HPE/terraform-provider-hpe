# EnableMaintenanceModeRequestServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IgnoreDaemonsets** | Pointer to **bool** | option relevant to kubernetes nodes | [optional] 
**Force** | Pointer to **bool** | option relevant to kubernetes nodes | [optional] 
**DeleteEmptyDir** | Pointer to **bool** | option relevant to kubernetes nodes | [optional] 
**DeleteLocalData** | Pointer to **bool** | option relevant to kubernetes nodes | [optional] 
**MovePoweredOff** | Pointer to **bool** | Move powered-off virtual machines to other hosts | [optional] [default to true]

## Usage

Instantiate with a Go composite literal:

```go
obj := &EnableMaintenanceModeRequestServer{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


