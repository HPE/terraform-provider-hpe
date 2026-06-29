# UpdateHostCloudRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CloudId** | Pointer to **int64** | The cloud/zone ID we are moving the set of servers to | [optional] 
**Servers** | Pointer to [**[]UpdateHostCloudRequestServersInner**](UpdateHostCloudRequestServersInner.md) | A JSON array of source: and target: server ids to be moved. If the target is blank Morpheus will automatically try to match by the servers unique or externalId | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateHostCloudRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


