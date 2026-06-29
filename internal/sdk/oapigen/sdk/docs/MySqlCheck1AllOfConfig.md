# MySqlCheck1AllOfConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DbPort** | **string** |  | 
**DbName** | **string** |  | 
**DbUser** | **string** |  | 
**DbHost** | **string** |  | 
**CheckOperator** | Pointer to **string** |  | [optional] 
**DbQuery** | **string** |  | 
**CheckResult** | Pointer to **int64** |  | [optional] 
**DbPassword** | **string** |  | 
**DbPasswordHash** | Pointer to **string** |  | [optional] 
**TunnelOn** | Pointer to **string** | Turn &#x60;on&#x60; to enable checks through a proxy host | [optional] [default to "off"]
**SshHost** | Pointer to **string** | Hostname or IP address of the proxy host | [optional] 
**SshPort** | Pointer to **int64** | Port for SSH on the proxy host, defaults to 22 | [optional] 
**SshUser** | Pointer to **string** | SSH user on the proxy host to login as | [optional] 
**SshPassword** | Pointer to **string** | Password for user, if not using key based authentication | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &MySqlCheck1AllOfConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


