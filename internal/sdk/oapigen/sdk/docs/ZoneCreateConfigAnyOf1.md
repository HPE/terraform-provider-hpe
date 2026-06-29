# ZoneCreateConfigAnyOf1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApplianceUrl** | Pointer to **string** | The URL used by workloads provisioned in the cloud for interacting with the Morpheus appliance. | [optional] 
**DatacenterName** | Pointer to **string** | A custom name used to reference the datacenter for the cloud. | [optional] 
**ExternalId** | Pointer to **NullableString** | The external id of the cloud | [optional] 
**CloudType** | Pointer to **string** | The Azure cloud type (global, usgov, german, china). | [optional] [default to "global"]
**InventoryLevel** | Pointer to **string** | Whether to import existing virtual machines. | [optional] 
**ImportExisting** | Pointer to **string** | Whether to import existing resources from the cloud (on, off). | [optional] 
**ConsoleKeymap** | Pointer to **string** | The keyboard layout to use for the console | [optional] 
**SubscriberId** | Pointer to **string** | Azure subscriber id | [optional] 
**TenantId** | Pointer to **string** | Azure tenant id | [optional] 
**ClientId** | Pointer to **string** | Azure client id | [optional] 
**ClientSecret** | Pointer to **string** | Azure client secret | [optional] 
**ResourceGroup** | Pointer to **string** | Azure resource group | [optional] 
**StorageAccount** | Pointer to **string** | The Azure storage account to use. | [optional] 
**RpcMode** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ZoneCreateConfigAnyOf1{
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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


