# NetworkInterfaceUpdateSuccessServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**DisplayName** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ZoneId** | Pointer to **int64** |  | [optional] 
**SiteId** | Pointer to **int64** |  | [optional] 
**SshHost** | Pointer to **string** |  | [optional] 
**SshPort** | Pointer to **int64** |  | [optional] 
**ExternalIp** | Pointer to **string** |  | [optional] 
**InternalIp** | Pointer to **string** |  | [optional] 
**VolumeId** | Pointer to **NullableString** |  | [optional] 
**Platform** | Pointer to **string** |  | [optional] 
**PlatformVersion** | Pointer to **string** |  | [optional] 
**SshUsername** | Pointer to **string** |  | [optional] 
**SshPassword** | Pointer to **string** |  | [optional] 
**OsDevice** | Pointer to **string** |  | [optional] 
**DataDevice** | Pointer to **string** |  | [optional] 
**LvmEnabled** | Pointer to **bool** |  | [optional] 
**ApiKey** | Pointer to **string** |  | [optional] 
**SoftwareRaid** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to **string** |  | [optional] 
**CapacityInfo** | Pointer to [**NetworkInterfaceUpdateSuccessServerCapacityInfo**](NetworkInterfaceUpdateSuccessServerCapacityInfo.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**LastStats** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**ComputeServerType** | Pointer to [**NetworkInterfaceUpdateSuccessServerComputeServerType**](NetworkInterfaceUpdateSuccessServerComputeServerType.md) |  | [optional] 
**Interfaces** | Pointer to [**[]NetworkInterfaceUpdateSuccessServerInterfacesInner**](NetworkInterfaceUpdateSuccessServerInterfacesInner.md) |  | [optional] 
**Zone** | Pointer to [**NetworkInterfaceUpdateSuccessServerZone**](NetworkInterfaceUpdateSuccessServerZone.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkInterfaceUpdateSuccessServer{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### VolumeId (Nullable)

Use the Nullable wrapper methods:
- `obj.VolumeId.IsSet()` — check if set
- `obj.VolumeId.Get()` — get the inner value (returns pointer)
- `obj.VolumeId.Set(&val)` — set the value
- `obj.VolumeId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


