# UpdateHostManagedRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Server** | Pointer to [**UpdateHostManagedRequestServer**](UpdateHostManagedRequestServer.md) |  | [optional] 
**InstallAgent** | Pointer to **bool** | Install agent. Set to false to manually install agent instead. | [optional] [default to true]
**InstanceTypeId** | Pointer to **int64** | Instance Type ID for the new Instance | [optional] 
**Layout** | Pointer to **int64** | Layout ID for the new Instance | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateHostManagedRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


