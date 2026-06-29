# UpdateHostInstallAgentRequestServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SshUsername** | Pointer to **string** | SSH username to use when provisioning | [optional] 
**SshPassword** | Pointer to **string** | SSH password to use, if not specified the account public key can be used | [optional] 
**ServerOs** | Pointer to [**UpdateHostInstallAgentRequestServerServerOs**](UpdateHostInstallAgentRequestServerServerOs.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateHostInstallAgentRequestServer{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


