# ClusterCreateServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Config** | [**ClusterCreateServerConfig**](ClusterCreateServerConfig.md) |  | 
**ServerType** | Pointer to [**ClusterCreateServerServerType**](ClusterCreateServerServerType.md) |  | [optional] 
**Name** | **string** | Name to be used for host(s) created in the cluster | 
**Plan** | [**ClusterCreateServerPlan**](ClusterCreateServerPlan.md) |  | 
**ServicePlanOptions** | Pointer to [**ClusterCreateServerServicePlanOptions**](ClusterCreateServerServicePlanOptions.md) |  | [optional] 
**Volumes** | Pointer to [**[]ClusterCreateServerVolumesInner**](ClusterCreateServerVolumesInner.md) | The (optional) volumes parameter is for LV configuration, can create additional LVs at provision It should be passed as an array of Objects | [optional] 
**Network** | Pointer to [**ClusterCreateServerNetwork**](ClusterCreateServerNetwork.md) |  | [optional] 
**NetworkInterfaces** | Pointer to [**[]ClusterCreateServerNetworkInterfacesInner**](ClusterCreateServerNetworkInterfacesInner.md) | The networkInterfaces parameter is for network configuration.  The Options API /api/options/zoneNetworkOptions can be used to see which options are available.  It should be passed as an array of Objects with the following attributes  | [optional] 
**SecurityGroups** | Pointer to **[]string** | Key for security group configuration. | [optional] 
**Visibility** | Pointer to **string** | Visibility for server host | [optional] [default to "private"]
**UserGroup** | Pointer to [**ClusterCreateServerUserGroup**](ClusterCreateServerUserGroup.md) |  | [optional] 
**NetworkDomain** | Pointer to **NullableString** | Network domain | [optional] 
**Hostname** | Pointer to **NullableString** | Hostname for server host | [optional] 
**NodeCount** | Pointer to **int64** | Number of workers or hosts | [optional] 
**Tags** | Pointer to [**[]ClusterCreateServerTagsInner**](ClusterCreateServerTagsInner.md) | Metadata tags, Array of objects having a name and value. | [optional] 
**Labels** | Pointer to **[]string** | Array of strings (keywords). This will set labels on the server and also on the cluster as well by default. | [optional] 
**SshHosts** | Pointer to [**[]ClusterCreateServerSshHostsInner**](ClusterCreateServerSshHostsInner.md) | Array of Host IPs and Names. This is used in conjunction with sshUsername and sshPassword/sshKeyPair to add existing hosts such as with HPE VM clusters. | [optional] 
**SshMasterHosts** | Pointer to **string** | A string consisting of comma-separated master host IP addresses. | [optional] 
**SshWorkerHosts** | Pointer to **string** | A string consisting of comma-separated worker host IP addresses. | [optional] 
**SshPort** | Pointer to **int64** | The port which the worker&#39;s SSH server is listening on. | [optional] 
**SshUsername** | Pointer to **string** | SSH Username | [optional] 
**SshPassword** | Pointer to **NullableString** | SSH Password | [optional] 
**SshKeyPair** | Pointer to [**ClusterCreateServerSshKeyPair**](ClusterCreateServerSshKeyPair.md) |  | [optional] 
**DataDevice** | Pointer to **string** |  | [optional] 
**LvmEnabled** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterCreateServer{
    // Set fields directly
}
```

### NetworkDomain (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkDomain.IsSet()` — check if set
- `obj.NetworkDomain.Get()` — get the inner value (returns pointer)
- `obj.NetworkDomain.Set(&val)` — set the value
- `obj.NetworkDomain.Unset()` — clear the value
### Hostname (Nullable)

Use the Nullable wrapper methods:
- `obj.Hostname.IsSet()` — check if set
- `obj.Hostname.Get()` — get the inner value (returns pointer)
- `obj.Hostname.Set(&val)` — set the value
- `obj.Hostname.Unset()` — clear the value
### SshPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.SshPassword.IsSet()` — check if set
- `obj.SshPassword.Get()` — get the inner value (returns pointer)
- `obj.SshPassword.Set(&val)` — set the value
- `obj.SshPassword.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


