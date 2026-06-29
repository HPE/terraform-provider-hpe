# ListVDIPools200ResponseAllOfVdiPoolsInnerConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Group** | Pointer to [**ListVDIPools200ResponseAllOfVdiPoolsInnerConfigGroup**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigGroup.md) |  | [optional] 
**Cloud** | Pointer to [**NullableListVDIPools200ResponseAllOfVdiPoolsInnerConfigCloud**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigCloud.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Instance** | Pointer to [**ListVDIPools200ResponseAllOfVdiPoolsInnerConfigInstance**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigInstance.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Environment** | Pointer to **string** |  | [optional] 
**Config** | Pointer to [**ListVDIPools200ResponseAllOfVdiPoolsInnerConfigConfig**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigConfig.md) |  | [optional] 
**Volumes** | Pointer to [**[]ListVDIPools200ResponseAllOfVdiPoolsInnerConfigVolumesInner**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigVolumesInner.md) |  | [optional] 
**HostName** | Pointer to **string** |  | [optional] 
**Layout** | Pointer to [**ListVDIPools200ResponseAllOfVdiPoolsInnerConfigLayout**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigLayout.md) |  | [optional] 
**StorageControllers** | Pointer to [**[]ListVDIPools200ResponseAllOfVdiPoolsInnerConfigStorageControllersInner**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigStorageControllersInner.md) |  | [optional] 
**Plan** | Pointer to [**ListVDIPools200ResponseAllOfVdiPoolsInnerConfigPlan**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigPlan.md) |  | [optional] 
**Version** | Pointer to **string** |  | [optional] 
**NetworkInterfaces** | Pointer to [**[]ListVDIPools200ResponseAllOfVdiPoolsInnerConfigNetworkInterfacesInner**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigNetworkInterfacesInner.md) |  | [optional] 
**ExecutionId** | Pointer to **string** |  | [optional] 
**Backup** | Pointer to [**ListVDIPools200ResponseAllOfVdiPoolsInnerConfigBackup**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigBackup.md) |  | [optional] 
**LoadBalancer** | Pointer to **[]map[string]interface{}** |  | [optional] 
**HideLock** | Pointer to **bool** |  | [optional] 
**HasNetworks** | Pointer to **bool** |  | [optional] 
**DisplayNetworks** | Pointer to [**[]ListVDIPools200ResponseAllOfVdiPoolsInnerConfigDisplayNetworksInner**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigDisplayNetworksInner.md) |  | [optional] 
**Copies** | Pointer to **int64** |  | [optional] 
**ShowScale** | Pointer to **bool** |  | [optional] 
**HasPreview** | Pointer to **bool** |  | [optional] 
**VolumesDisplay** | Pointer to [**[]ListVDIPools200ResponseAllOfVdiPoolsInnerConfigVolumesDisplayInner**](ListVDIPools200ResponseAllOfVdiPoolsInnerConfigVolumesDisplayInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListVDIPools200ResponseAllOfVdiPoolsInnerConfig{
    // Set fields directly
}
```

### Cloud (Nullable)

Use the Nullable wrapper methods:
- `obj.Cloud.IsSet()` — check if set
- `obj.Cloud.Get()` — get the inner value (returns pointer)
- `obj.Cloud.Set(&val)` — set the value
- `obj.Cloud.Unset()` — clear the value
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


