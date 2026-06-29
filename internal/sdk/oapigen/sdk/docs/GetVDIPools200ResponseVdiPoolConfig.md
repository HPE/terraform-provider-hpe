# GetVDIPools200ResponseVdiPoolConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Group** | Pointer to [**GetVDIPools200ResponseVdiPoolConfigGroup**](GetVDIPools200ResponseVdiPoolConfigGroup.md) |  | [optional] 
**Cloud** | Pointer to [**GetVDIPools200ResponseVdiPoolConfigCloud**](GetVDIPools200ResponseVdiPoolConfigCloud.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Instance** | Pointer to [**GetVDIPools200ResponseVdiPoolConfigInstance**](GetVDIPools200ResponseVdiPoolConfigInstance.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Environment** | Pointer to **string** |  | [optional] 
**Config** | Pointer to [**GetVDIPools200ResponseVdiPoolConfigConfig**](GetVDIPools200ResponseVdiPoolConfigConfig.md) |  | [optional] 
**Volumes** | Pointer to [**[]GetVDIPools200ResponseVdiPoolConfigVolumesInner**](GetVDIPools200ResponseVdiPoolConfigVolumesInner.md) |  | [optional] 
**HostName** | Pointer to **string** |  | [optional] 
**Layout** | Pointer to [**GetVDIPools200ResponseVdiPoolConfigLayout**](GetVDIPools200ResponseVdiPoolConfigLayout.md) |  | [optional] 
**StorageControllers** | Pointer to [**[]GetVDIPools200ResponseVdiPoolConfigStorageControllersInner**](GetVDIPools200ResponseVdiPoolConfigStorageControllersInner.md) |  | [optional] 
**Plan** | Pointer to [**GetVDIPools200ResponseVdiPoolConfigPlan**](GetVDIPools200ResponseVdiPoolConfigPlan.md) |  | [optional] 
**Version** | Pointer to **string** |  | [optional] 
**NetworkInterfaces** | Pointer to [**[]GetVDIPools200ResponseVdiPoolConfigNetworkInterfacesInner**](GetVDIPools200ResponseVdiPoolConfigNetworkInterfacesInner.md) |  | [optional] 
**ExecutionId** | Pointer to **string** |  | [optional] 
**Backup** | Pointer to [**GetVDIPools200ResponseVdiPoolConfigBackup**](GetVDIPools200ResponseVdiPoolConfigBackup.md) |  | [optional] 
**LoadBalancer** | Pointer to **[]map[string]interface{}** |  | [optional] 
**HideLock** | Pointer to **bool** |  | [optional] 
**HasNetworks** | Pointer to **bool** |  | [optional] 
**DisplayNetworks** | Pointer to [**[]GetVDIPools200ResponseVdiPoolConfigDisplayNetworksInner**](GetVDIPools200ResponseVdiPoolConfigDisplayNetworksInner.md) |  | [optional] 
**Copies** | Pointer to **int64** |  | [optional] 
**ShowScale** | Pointer to **bool** |  | [optional] 
**HasPreview** | Pointer to **bool** |  | [optional] 
**VolumesDisplay** | Pointer to [**[]GetVDIPools200ResponseVdiPoolConfigVolumesDisplayInner**](GetVDIPools200ResponseVdiPoolConfigVolumesDisplayInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetVDIPools200ResponseVdiPoolConfig{
    // Set fields directly
}
```

### StorageControllers (Nullable)

Use the Nullable wrapper methods:
- `obj.StorageControllers.IsSet()` — check if set
- `obj.StorageControllers.Get()` — get the inner value (returns pointer)
- `obj.StorageControllers.Set(&val)` — set the value
- `obj.StorageControllers.Unset()` — clear the value
### LoadBalancer (Nullable)

Use the Nullable wrapper methods:
- `obj.LoadBalancer.IsSet()` — check if set
- `obj.LoadBalancer.Get()` — get the inner value (returns pointer)
- `obj.LoadBalancer.Set(&val)` — set the value
- `obj.LoadBalancer.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


