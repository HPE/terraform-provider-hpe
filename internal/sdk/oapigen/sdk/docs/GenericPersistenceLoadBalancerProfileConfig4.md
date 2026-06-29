# GenericPersistenceLoadBalancerProfileConfig4

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProfileType** | Pointer to **string** | Profile category. Automatically set from serviceType (persistence-profile). | [optional] 
**SharePersistence** | Pointer to **bool** | Whether persistence is shared across virtual servers that use this profile. | [optional] 
**HaPersistenceMirroring** | Pointer to **bool** | Whether persistence entries are synchronized to the standby node for high availability. | [optional] 
**PersistenceEntryTimeout** | Pointer to **NullableInt64** | Persistence entry timeout in seconds. Defaults to 300. | [optional] [default to 300]
**Tags** | Pointer to [**[]LoadBalancerProfileTag37**](LoadBalancerProfileTag37.md) | NSX-T tags applied to the profile. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GenericPersistenceLoadBalancerProfileConfig4{
    // Set fields directly
}
```

### PersistenceEntryTimeout (Nullable)

Use the Nullable wrapper methods:
- `obj.PersistenceEntryTimeout.IsSet()` — check if set
- `obj.PersistenceEntryTimeout.Get()` — get the inner value (returns pointer)
- `obj.PersistenceEntryTimeout.Set(&val)` — set the value
- `obj.PersistenceEntryTimeout.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


