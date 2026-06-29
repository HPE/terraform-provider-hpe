# UpdateClusterRequestCluster

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Cluster name | [optional] 
**Description** | Pointer to **string** | Cluster description | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Enabled** | Pointer to **bool** | Cluster enabled | [optional] 
**UseAgent** | Pointer to **bool** | Use the Agent to relay communications for the Kubernetes API instead of direct. | [optional] 
**ServiceUrl** | Pointer to **string** | Cluster API Url | [optional] 
**ServiceToken** | Pointer to **string** | Cluster API token | [optional] 
**Refresh** | Pointer to **bool** | Queue cluster refresh | [optional] 
**Managed** | Pointer to **bool** | Cluster managed | [optional] 
**AutoRecoverPowerState** | Pointer to **bool** | Automatically Power on VMs | [optional] 
**CpuPlacementMode** | Pointer to **string** | Cluster CPU placement mode | [optional] 
**Integrations** | Pointer to [**[]UpdateClusterRequestClusterIntegrationsInner**](UpdateClusterRequestClusterIntegrationsInner.md) | Cluster integrations | [optional] 
**Config** | Pointer to [**UpdateClusterRequestClusterConfig**](UpdateClusterRequestClusterConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateClusterRequestCluster{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


