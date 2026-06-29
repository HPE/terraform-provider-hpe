# SnapshotInstance200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** |  | [optional] 
**ProcessIds** | Pointer to **[]int64** | Process ID(s) to track execution results with, use &#x60;/api/processes/$processId&#x60;. Instances with more than one server will result in multiple processes, one for each snapshot. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SnapshotInstance200Response{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


