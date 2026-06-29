# GetPowerSchedules200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Instances** | Pointer to [**[]GetPowerSchedules200ResponseAllOfInstancesInner**](GetPowerSchedules200ResponseAllOfInstancesInner.md) |  | [optional] 
**Servers** | Pointer to [**[]GetPowerSchedules200ResponseAllOfServersInner**](GetPowerSchedules200ResponseAllOfServersInner.md) |  | [optional] 
**Schedule** | Pointer to [**GetPowerSchedules200ResponseAllOfSchedule**](GetPowerSchedules200ResponseAllOfSchedule.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetPowerSchedules200Response{
    // Set fields directly
}
```

### Instances (Nullable)

Use the Nullable wrapper methods:
- `obj.Instances.IsSet()` — check if set
- `obj.Instances.Get()` — get the inner value (returns pointer)
- `obj.Instances.Set(&val)` — set the value
- `obj.Instances.Unset()` — clear the value
### Servers (Nullable)

Use the Nullable wrapper methods:
- `obj.Servers.IsSet()` — check if set
- `obj.Servers.Get()` — get the inner value (returns pointer)
- `obj.Servers.Set(&val)` — set the value
- `obj.Servers.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


