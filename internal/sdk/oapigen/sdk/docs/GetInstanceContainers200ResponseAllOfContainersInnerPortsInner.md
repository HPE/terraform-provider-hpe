# GetInstanceContainers200ResponseAllOfContainersInnerPortsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Index** | Pointer to **int64** |  | [optional] 
**External** | Pointer to **int64** |  | [optional] 
**Internal** | Pointer to **int64** |  | [optional] 
**DisplayName** | Pointer to **string** |  | [optional] 
**PrimaryPort** | Pointer to **bool** |  | [optional] 
**Export** | Pointer to **bool** |  | [optional] 
**Visible** | Pointer to **bool** |  | [optional] 
**ExportName** | Pointer to **string** |  | [optional] 
**LoadBalanceProtocol** | Pointer to **string** |  | [optional] 
**LoadBalance** | Pointer to **bool** |  | [optional] 
**Protocol** | Pointer to **string** |  | [optional] 
**Link** | Pointer to **bool** |  | [optional] 
**ExternalIp** | Pointer to **NullableString** |  | [optional] 
**InternalIp** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetInstanceContainers200ResponseAllOfContainersInnerPortsInner{
    // Set fields directly
}
```

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


