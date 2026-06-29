# HVMInstanceConfiguration1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NoAgent** | Pointer to **NullableBool** | Skipping Agent installation will result in a lack of logging and guest operating system statistics. Automation scripts may also be adversely affected. | [optional] [default to false]
**ResourcePoolId** | Pointer to **string** | id of the resource group to be used, can be prefixed with &#x60;pool-&#x60;. A resource pool group can be specified instead by prefixing its ID with &#x60;poolGroup-&#x60;. | [optional] 
**NestedVirtualization** | Pointer to **string** | Enable Nested Virtualization | [optional] [default to "off"]
**CreateUser** | Pointer to **NullableBool** | Create user | [optional] [default to false]
**PoolProviderType** | Pointer to **string** | The type of pool provider to use for this instance, must be \&quot;mvm\&quot; | [optional] [default to "mvm"]
**KvmHostId** | Pointer to **int64** | The ID of the KVM host to provision the instance on | [optional] 
**ProvisionPoweredOff** | Pointer to **bool** | Whether to provision the instance in a powered off state | [optional] [default to false]

## Usage

Instantiate with a Go composite literal:

```go
obj := &HVMInstanceConfiguration1{
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


