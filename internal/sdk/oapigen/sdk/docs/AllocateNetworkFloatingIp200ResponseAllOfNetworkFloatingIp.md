# AllocateNetworkFloatingIp200ResponseAllOfNetworkFloatingIp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Cloud** | Pointer to [**AllocateNetworkFloatingIp200ResponseAllOfNetworkFloatingIpCloud**](AllocateNetworkFloatingIp200ResponseAllOfNetworkFloatingIpCloud.md) |  | [optional] 
**Server** | Pointer to [**AllocateNetworkFloatingIp200ResponseAllOfNetworkFloatingIpServer**](AllocateNetworkFloatingIp200ResponseAllOfNetworkFloatingIpServer.md) |  | [optional] 
**IpStatus** | Pointer to **string** |  | [optional] 
**IpAddress** | Pointer to **string** | IP Address | [optional] 
**IpRange** | Pointer to **NullableString** |  | [optional] 
**PtrId** | Pointer to **NullableString** |  | [optional] 
**NetworkDomain** | Pointer to [**AllocateNetworkFloatingIp200ResponseAllOfNetworkFloatingIpNetworkDomain**](AllocateNetworkFloatingIp200ResponseAllOfNetworkFloatingIpNetworkDomain.md) |  | [optional] 
**CreatedBy** | Pointer to [**AllocateNetworkFloatingIp200ResponseAllOfNetworkFloatingIpCreatedBy**](AllocateNetworkFloatingIp200ResponseAllOfNetworkFloatingIpCreatedBy.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AllocateNetworkFloatingIp200ResponseAllOfNetworkFloatingIp{
    // Set fields directly
}
```

### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### IpRange (Nullable)

Use the Nullable wrapper methods:
- `obj.IpRange.IsSet()` — check if set
- `obj.IpRange.Get()` — get the inner value (returns pointer)
- `obj.IpRange.Set(&val)` — set the value
- `obj.IpRange.Unset()` — clear the value
### PtrId (Nullable)

Use the Nullable wrapper methods:
- `obj.PtrId.IsSet()` — check if set
- `obj.PtrId.Get()` — get the inner value (returns pointer)
- `obj.PtrId.Set(&val)` — set the value
- `obj.PtrId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


