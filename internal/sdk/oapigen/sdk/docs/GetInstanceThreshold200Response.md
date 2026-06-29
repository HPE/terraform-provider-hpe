# GetInstanceThreshold200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Instance** | Pointer to [**NullableGetInstanceThreshold200ResponseInstance**](GetInstanceThreshold200ResponseInstance.md) |  | [optional] 
**InstanceThreshold** | Pointer to [**GetInstanceThreshold200ResponseInstanceThreshold**](GetInstanceThreshold200ResponseInstanceThreshold.md) |  | [optional] 
**InstanceSchedules** | Pointer to [**[]GetInstanceThreshold200ResponseInstanceSchedulesInner**](GetInstanceThreshold200ResponseInstanceSchedulesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetInstanceThreshold200Response{
    // Set fields directly
}
```

### Instance (Nullable)

Use the Nullable wrapper methods:
- `obj.Instance.IsSet()` — check if set
- `obj.Instance.Get()` — get the inner value (returns pointer)
- `obj.Instance.Set(&val)` — set the value
- `obj.Instance.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


