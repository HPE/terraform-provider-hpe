# UpdateHostManagedRequestServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SshHost** | Pointer to **string** | Hostname or IP address of the server to be managed | [optional] 
**SshUsername** | Pointer to **string** | SSH username to use when provisioning | [optional] 
**SshPassword** | Pointer to **string** | SSH password to use, if not specified the account public key can be used | [optional] 
**SshKeyPair** | Pointer to [**UpdateHostManagedRequestServerSshKeyPair**](UpdateHostManagedRequestServerSshKeyPair.md) |  | [optional] 
**ServerOs** | Pointer to [**UpdateHostManagedRequestServerServerOs**](UpdateHostManagedRequestServerServerOs.md) |  | [optional] 
**Plan** | Pointer to [**UpdateHostManagedRequestServerPlan**](UpdateHostManagedRequestServerPlan.md) |  | [optional] 
**Account** | Pointer to [**UpdateHostManagedRequestServerAccount**](UpdateHostManagedRequestServerAccount.md) |  | [optional] 
**ProvisionSiteId** | Pointer to **int64** | Specific group to assign the server | [optional] 
**Tags** | Pointer to [**[]UpdateHostManagedRequestServerTagsInner**](UpdateHostManagedRequestServerTagsInner.md) | Metadata tags, Array of objects having a name and value, this adds or updates the specified tags and removes any tags not specified. | [optional] 
**Config** | Pointer to [**UpdateHostManagedRequestServerConfig**](UpdateHostManagedRequestServerConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateHostManagedRequestServer{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


