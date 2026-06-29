# AddCloudsRequestZoneConfigAnyOf3

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApplianceUrl** | Pointer to **string** | The URL used by workloads provisioned in the cloud for interacting with the Morpheus appliance. | [optional] 
**DatacenterName** | Pointer to **string** | A custom name used to reference the datacenter for the cloud. | [optional] 
**ExternalId** | Pointer to **NullableString** | The external id of the cloud | [optional] 
**InventoryLevel** | Pointer to **string** | Whether to import existing virtual machines. | [optional] 
**ConsoleKeymap** | Pointer to **string** | The keyboard layout to use for the console | [optional] 
**ApiUrl** | **string** | The SDK URL of the vCenter server. | 
**ApiVersion** | **string** | The SDK version of the vCenter server. | 
**Datacenter** | **string** | The vSphere datacenter to add. | 
**Cluster** | Pointer to **string** | The name of the vSphere cluster | [optional] [default to "all"]
**ConfigManagementId** | Pointer to **string** | The id of the configuration management integration associated with the vSphere cloud. | [optional] 
**ResourcePool** | Pointer to **string** | The name of the vSphere resource pool | [optional] 
**RpcMode** | Pointer to **NullableString** |  | [optional] 
**StorageType** | Pointer to **string** | The default vSphere VMDK type for virtual machines | [optional] [default to "thin"]
**CertificateProvider** | Pointer to **string** | Certificate provider | [optional] [default to "internal"]
**EnableVnc** | Pointer to **NullableString** |  | [optional] 
**HideHostSelection** | Pointer to **NullableString** |  | [optional] 
**EnableDiskTypeSelection** | Pointer to **NullableString** |  | [optional] 
**EnableStorageTypeSelection** | Pointer to **NullableString** |  | [optional] 
**EnableNetworkTypeSelection** | Pointer to **NullableString** |  | [optional] 
**Username** | Pointer to **string** | Username. | [optional] 
**Password** | Pointer to **string** | Password to apply to the user | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCloudsRequestZoneConfigAnyOf3{
    // Set fields directly
}
```

### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### RpcMode (Nullable)

Use the Nullable wrapper methods:
- `obj.RpcMode.IsSet()` — check if set
- `obj.RpcMode.Get()` — get the inner value (returns pointer)
- `obj.RpcMode.Set(&val)` — set the value
- `obj.RpcMode.Unset()` — clear the value
### EnableVnc (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableVnc.IsSet()` — check if set
- `obj.EnableVnc.Get()` — get the inner value (returns pointer)
- `obj.EnableVnc.Set(&val)` — set the value
- `obj.EnableVnc.Unset()` — clear the value
### HideHostSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.HideHostSelection.IsSet()` — check if set
- `obj.HideHostSelection.Get()` — get the inner value (returns pointer)
- `obj.HideHostSelection.Set(&val)` — set the value
- `obj.HideHostSelection.Unset()` — clear the value
### EnableDiskTypeSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableDiskTypeSelection.IsSet()` — check if set
- `obj.EnableDiskTypeSelection.Get()` — get the inner value (returns pointer)
- `obj.EnableDiskTypeSelection.Set(&val)` — set the value
- `obj.EnableDiskTypeSelection.Unset()` — clear the value
### EnableStorageTypeSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableStorageTypeSelection.IsSet()` — check if set
- `obj.EnableStorageTypeSelection.Get()` — get the inner value (returns pointer)
- `obj.EnableStorageTypeSelection.Set(&val)` — set the value
- `obj.EnableStorageTypeSelection.Unset()` — clear the value
### EnableNetworkTypeSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableNetworkTypeSelection.IsSet()` — check if set
- `obj.EnableNetworkTypeSelection.Get()` — get the inner value (returns pointer)
- `obj.EnableNetworkTypeSelection.Set(&val)` — set the value
- `obj.EnableNetworkTypeSelection.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


