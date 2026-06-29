# CreateLoadBalancerPoolNode200ResponseAllOfLoadBalancerNode

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**IpAddress** | Pointer to **string** |  | [optional] 
**Port** | Pointer to **int64** |  | [optional] 
**PortType** | Pointer to **NullableString** |  | [optional] 
**MonitorPort** | Pointer to **NullableString** |  | [optional] 
**Weight** | Pointer to **NullableInt64** |  | [optional] 
**NodeState** | Pointer to **NullableString** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**Server** | Pointer to [**CreateLoadBalancerPoolNode200ResponseAllOfLoadBalancerNodeServer**](CreateLoadBalancerPoolNode200ResponseAllOfLoadBalancerNodeServer.md) |  | [optional] 
**InstanceId** | Pointer to **NullableInt64** |  | [optional] 
**ContainerId** | Pointer to **NullableInt64** |  | [optional] 
**NodeSource** | Pointer to **NullableString** |  | [optional] 
**Monitor** | Pointer to [**CreateLoadBalancerPoolNode200ResponseAllOfLoadBalancerNodeMonitor**](CreateLoadBalancerPoolNode200ResponseAllOfLoadBalancerNodeMonitor.md) |  | [optional] 
**MaxConnections** | Pointer to **NullableInt64** |  | [optional] 
**ExternalRefType** | Pointer to **NullableString** |  | [optional] 
**ExternalRefId** | Pointer to **NullableString** |  | [optional] 
**ExternalRefName** | Pointer to **NullableString** |  | [optional] 
**CreatedBy** | Pointer to [**CreateLoadBalancerPoolNode200ResponseAllOfLoadBalancerNodeCreatedBy**](CreateLoadBalancerPoolNode200ResponseAllOfLoadBalancerNodeCreatedBy.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateLoadBalancerPoolNode200ResponseAllOfLoadBalancerNode{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### PortType (Nullable)

Use the Nullable wrapper methods:
- `obj.PortType.IsSet()` — check if set
- `obj.PortType.Get()` — get the inner value (returns pointer)
- `obj.PortType.Set(&val)` — set the value
- `obj.PortType.Unset()` — clear the value
### MonitorPort (Nullable)

Use the Nullable wrapper methods:
- `obj.MonitorPort.IsSet()` — check if set
- `obj.MonitorPort.Get()` — get the inner value (returns pointer)
- `obj.MonitorPort.Set(&val)` — set the value
- `obj.MonitorPort.Unset()` — clear the value
### Weight (Nullable)

Use the Nullable wrapper methods:
- `obj.Weight.IsSet()` — check if set
- `obj.Weight.Get()` — get the inner value (returns pointer)
- `obj.Weight.Set(&val)` — set the value
- `obj.Weight.Unset()` — clear the value
### NodeState (Nullable)

Use the Nullable wrapper methods:
- `obj.NodeState.IsSet()` — check if set
- `obj.NodeState.Get()` — get the inner value (returns pointer)
- `obj.NodeState.Set(&val)` — set the value
- `obj.NodeState.Unset()` — clear the value
### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### StatusDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusDate.IsSet()` — check if set
- `obj.StatusDate.Get()` — get the inner value (returns pointer)
- `obj.StatusDate.Set(&val)` — set the value
- `obj.StatusDate.Unset()` — clear the value
### InstanceId (Nullable)

Use the Nullable wrapper methods:
- `obj.InstanceId.IsSet()` — check if set
- `obj.InstanceId.Get()` — get the inner value (returns pointer)
- `obj.InstanceId.Set(&val)` — set the value
- `obj.InstanceId.Unset()` — clear the value
### ContainerId (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerId.IsSet()` — check if set
- `obj.ContainerId.Get()` — get the inner value (returns pointer)
- `obj.ContainerId.Set(&val)` — set the value
- `obj.ContainerId.Unset()` — clear the value
### NodeSource (Nullable)

Use the Nullable wrapper methods:
- `obj.NodeSource.IsSet()` — check if set
- `obj.NodeSource.Get()` — get the inner value (returns pointer)
- `obj.NodeSource.Set(&val)` — set the value
- `obj.NodeSource.Unset()` — clear the value
### MaxConnections (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxConnections.IsSet()` — check if set
- `obj.MaxConnections.Get()` — get the inner value (returns pointer)
- `obj.MaxConnections.Set(&val)` — set the value
- `obj.MaxConnections.Unset()` — clear the value
### ExternalRefType (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalRefType.IsSet()` — check if set
- `obj.ExternalRefType.Get()` — get the inner value (returns pointer)
- `obj.ExternalRefType.Set(&val)` — set the value
- `obj.ExternalRefType.Unset()` — clear the value
### ExternalRefId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalRefId.IsSet()` — check if set
- `obj.ExternalRefId.Get()` — get the inner value (returns pointer)
- `obj.ExternalRefId.Set(&val)` — set the value
- `obj.ExternalRefId.Unset()` — clear the value
### ExternalRefName (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalRefName.IsSet()` — check if set
- `obj.ExternalRefName.Get()` — get the inner value (returns pointer)
- `obj.ExternalRefName.Set(&val)` — set the value
- `obj.ExternalRefName.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


