# InstanceContainerPortsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**External** | Pointer to **int64** |  | [optional] 
**Internal** | Pointer to **int64** |  | [optional] 
**DisplayName** | Pointer to **string** |  | [optional] 
**PrimaryPort** | Pointer to **bool** |  | [optional] 
**Export** | Pointer to **bool** |  | [optional] 
**Visible** | Pointer to **bool** |  | [optional] 
**ExportName** | Pointer to **string** |  | [optional] 
**LoadBalanceProtocol** | Pointer to **NullableString** |  | [optional] 
**LoadBalance** | Pointer to **bool** |  | [optional] 
**Protocol** | Pointer to **NullableString** |  | [optional] 
**Link** | Pointer to **bool** |  | [optional] 
**ExternalIp** | Pointer to **NullableString** |  | [optional] 
**InternalIp** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceContainerPortsInner{
    // Set fields directly
}
```

### LoadBalanceProtocol (Nullable)

Use the Nullable wrapper methods:
- `obj.LoadBalanceProtocol.IsSet()` — check if set
- `obj.LoadBalanceProtocol.Get()` — get the inner value (returns pointer)
- `obj.LoadBalanceProtocol.Set(&val)` — set the value
- `obj.LoadBalanceProtocol.Unset()` — clear the value
### Protocol (Nullable)

Use the Nullable wrapper methods:
- `obj.Protocol.IsSet()` — check if set
- `obj.Protocol.Get()` — get the inner value (returns pointer)
- `obj.Protocol.Set(&val)` — set the value
- `obj.Protocol.Unset()` — clear the value
### ExternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalIp.IsSet()` — check if set
- `obj.ExternalIp.Get()` — get the inner value (returns pointer)
- `obj.ExternalIp.Set(&val)` — set the value
- `obj.ExternalIp.Unset()` — clear the value
### InternalIp (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalIp.IsSet()` — check if set
- `obj.InternalIp.Get()` — get the inner value (returns pointer)
- `obj.InternalIp.Set(&val)` — set the value
- `obj.InternalIp.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


