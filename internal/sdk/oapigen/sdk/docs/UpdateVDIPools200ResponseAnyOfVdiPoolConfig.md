# UpdateVDIPools200ResponseAnyOfVdiPoolConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Group** | Pointer to [**UpdateVDIPools200ResponseAnyOfVdiPoolConfigGroup**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigGroup.md) |  | [optional] 
**Cloud** | Pointer to [**UpdateVDIPools200ResponseAnyOfVdiPoolConfigCloud**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigCloud.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Instance** | Pointer to [**UpdateVDIPools200ResponseAnyOfVdiPoolConfigInstance**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigInstance.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Environment** | Pointer to **string** |  | [optional] 
**Config** | Pointer to [**UpdateVDIPools200ResponseAnyOfVdiPoolConfigConfig**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigConfig.md) |  | [optional] 
**Volumes** | Pointer to [**[]UpdateVDIPools200ResponseAnyOfVdiPoolConfigVolumesInner**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigVolumesInner.md) |  | [optional] 
**HostName** | Pointer to **string** |  | [optional] 
**Layout** | Pointer to [**UpdateVDIPools200ResponseAnyOfVdiPoolConfigLayout**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigLayout.md) |  | [optional] 
**StorageControllers** | Pointer to [**[]UpdateVDIPools200ResponseAnyOfVdiPoolConfigStorageControllersInner**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigStorageControllersInner.md) |  | [optional] 
**Plan** | Pointer to [**UpdateVDIPools200ResponseAnyOfVdiPoolConfigPlan**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigPlan.md) |  | [optional] 
**Version** | Pointer to **string** |  | [optional] 
**NetworkInterfaces** | Pointer to [**[]UpdateVDIPools200ResponseAnyOfVdiPoolConfigNetworkInterfacesInner**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigNetworkInterfacesInner.md) |  | [optional] 
**ExecutionId** | Pointer to **string** |  | [optional] 
**Backup** | Pointer to [**UpdateVDIPools200ResponseAnyOfVdiPoolConfigBackup**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigBackup.md) |  | [optional] 
**LoadBalancer** | Pointer to **[]map[string]interface{}** |  | [optional] 
**HideLock** | Pointer to **bool** |  | [optional] 
**HasNetworks** | Pointer to **bool** |  | [optional] 
**DisplayNetworks** | Pointer to [**[]UpdateVDIPools200ResponseAnyOfVdiPoolConfigDisplayNetworksInner**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigDisplayNetworksInner.md) |  | [optional] 
**Copies** | Pointer to **int64** |  | [optional] 
**ShowScale** | Pointer to **bool** |  | [optional] 
**HasPreview** | Pointer to **bool** |  | [optional] 
**VolumesDisplay** | Pointer to [**[]UpdateVDIPools200ResponseAnyOfVdiPoolConfigVolumesDisplayInner**](UpdateVDIPools200ResponseAnyOfVdiPoolConfigVolumesDisplayInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateVDIPools200ResponseAnyOfVdiPoolConfig{
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


