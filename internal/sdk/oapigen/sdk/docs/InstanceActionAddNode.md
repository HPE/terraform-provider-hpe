# InstanceActionAddNode

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | Pointer to **string** | The instance type code for the node being added during scale operation | [optional] 
**SelectedResourcePoolId** | Pointer to **int64** | The ID of the resource pool where the node will be created | [optional] 
**PreProvisioned** | Pointer to **string** | Set to &#39;on&#39; to use a pre-provisioned server | [optional] 
**SelectedServerId** | Pointer to **int64** | The ID of the pre-provisioned server | [optional] 
**SshUsername** | Pointer to **string** | SSH username for connecting to the pre-provisioned server | [optional] 
**SshPassword** | Pointer to **string** | SSH password for connecting to the pre-provisioned server | [optional] 
**SshHost** | Pointer to **string** | The SSH host IP address for the pre-provisioned server | [optional] 
**SshKeyPair** | Pointer to [**InstanceActionAddNodeSshKeyPair**](InstanceActionAddNodeSshKeyPair.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceActionAddNode{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


