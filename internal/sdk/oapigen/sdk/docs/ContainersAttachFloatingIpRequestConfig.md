# ContainersAttachFloatingIpRequestConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**OsExternalNetworkId** | **string** | The Floating IP identifier in the format: \&quot;ip-ID\&quot; or \&quot;pool-ID\&quot;.  The Options API /api/options/openStack/openstackFloatingIpOptions?containerId&#x3D;{{containerId}} can be used to see which options are available.  | 
**FloatingIpBandwidth** | Pointer to **float32** | Bandwidth (Mbit/s) Only the following cloud types support this parameter: Huawei, OpenTelekom  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ContainersAttachFloatingIpRequestConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


