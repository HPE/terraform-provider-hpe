# VMWareInstanceConfiguration3

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NoAgent** | Pointer to **NullableBool** | Skipping Agent installation will result in a lack of logging and guest operating system statistics. Automation scripts may also be adversely affected. | [optional] [default to false]
**ResourcePoolId** | Pointer to **string** | id of the resource group to be used, can be prefixed with &#x60;pool-&#x60;. A resource pool group can be specified instead by prefixing its ID with &#x60;poolGroup-&#x60;. | [optional] 
**HostId** | Pointer to **string** | Specific host to deploy to if so desired. | [optional] 
**SmbiosAssetTag** | Pointer to **string** | Sets the asset tag on the SMBIOS for use by the guest operating system. If left blank, the virtual machine name will be used. | [optional] 
**NestedVirtualization** | Pointer to **string** | Enable Nested Virtualization | [optional] [default to "off"]
**VmwareFolderId** | Pointer to **string** | VMWare Folder External ID (as a String) or ID (as an Integer or String) | [optional] 
**CreateUser** | Pointer to **NullableBool** | Create user | [optional] [default to false]
**Template** | Pointer to **int64** | Image ID. This is the ID of a Virtual Image. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &VMWareInstanceConfiguration3{
    // Set fields directly
}
```

### NoAgent (Nullable)

Use the Nullable wrapper methods:
- `obj.NoAgent.IsSet()` — check if set
- `obj.NoAgent.Get()` — get the inner value (returns pointer)
- `obj.NoAgent.Set(&val)` — set the value
- `obj.NoAgent.Unset()` — clear the value
### CreateUser (Nullable)

Use the Nullable wrapper methods:
- `obj.CreateUser.IsSet()` — check if set
- `obj.CreateUser.Get()` — get the inner value (returns pointer)
- `obj.CreateUser.Set(&val)` — set the value
- `obj.CreateUser.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


