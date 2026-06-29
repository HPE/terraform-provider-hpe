# NSXDHCPServerConfiguration1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EdgeCluster** | Pointer to **NullableString** | Edge Cluster | [optional] 
**PreferredEdgeNode1** | Pointer to **NullableString** | Active Edge Node Options obtained by calling option source with :optionSource &#x3D; nsxtEdgeNodes and networkServerId param | [optional] 
**PreferredEdgeNode2** | Pointer to **NullableString** | Standby Edge Node Options obtained by calling option source with optionSource &#x3D; nsxtEdgeNodes and networkServerId param | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NSXDHCPServerConfiguration1{
    // Set fields directly
}
```

### EdgeCluster (Nullable)

Use the Nullable wrapper methods:
- `obj.EdgeCluster.IsSet()` — check if set
- `obj.EdgeCluster.Get()` — get the inner value (returns pointer)
- `obj.EdgeCluster.Set(&val)` — set the value
- `obj.EdgeCluster.Unset()` — clear the value
### PreferredEdgeNode1 (Nullable)

Use the Nullable wrapper methods:
- `obj.PreferredEdgeNode1.IsSet()` — check if set
- `obj.PreferredEdgeNode1.Get()` — get the inner value (returns pointer)
- `obj.PreferredEdgeNode1.Set(&val)` — set the value
- `obj.PreferredEdgeNode1.Unset()` — clear the value
### PreferredEdgeNode2 (Nullable)

Use the Nullable wrapper methods:
- `obj.PreferredEdgeNode2.IsSet()` — check if set
- `obj.PreferredEdgeNode2.Get()` — get the inner value (returns pointer)
- `obj.PreferredEdgeNode2.Set(&val)` — set the value
- `obj.PreferredEdgeNode2.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


