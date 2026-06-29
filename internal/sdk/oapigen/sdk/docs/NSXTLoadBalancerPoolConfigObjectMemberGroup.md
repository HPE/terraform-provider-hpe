# NSXTLoadBalancerPoolConfigObjectMemberGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Path** | Pointer to **string** | NSX-T member group path. | [optional] 
**IpRevisionFilter** | Pointer to **string** | IP revision filter for the member group. | [optional] 
**MaxIpListSize** | Pointer to **NullableInt64** | Maximum IP list size for the member group. | [optional] 
**Port** | Pointer to **NullableInt64** | Port number for the member group. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NSXTLoadBalancerPoolConfigObjectMemberGroup{
    // Set fields directly
}
```

### MaxIpListSize (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxIpListSize.IsSet()` — check if set
- `obj.MaxIpListSize.Get()` — get the inner value (returns pointer)
- `obj.MaxIpListSize.Set(&val)` — set the value
- `obj.MaxIpListSize.Unset()` — clear the value
### Port (Nullable)

Use the Nullable wrapper methods:
- `obj.Port.IsSet()` — check if set
- `obj.Port.Get()` — get the inner value (returns pointer)
- `obj.Port.Set(&val)` — set the value
- `obj.Port.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


