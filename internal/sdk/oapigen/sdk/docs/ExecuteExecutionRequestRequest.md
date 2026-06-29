# ExecuteExecutionRequestRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Script** | **string** | A script or command to be executed | 
**SendKeys** | Pointer to **bool** | Pass true to send key mappings to the hypervisor console so commands such as &lt;LEFT&gt;, &lt;RIGHT&gt; and &lt;WAIT&gt; can be used. | [optional] [default to false]

## Usage

Instantiate with a Go composite literal:

```go
obj := &ExecuteExecutionRequestRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


