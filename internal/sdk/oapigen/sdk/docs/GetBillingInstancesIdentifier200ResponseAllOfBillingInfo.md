# GetBillingInstancesIdentifier200ResponseAllOfBillingInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InstanceId** | Pointer to **NullableInt64** |  | [optional] 
**InstanceUUID** | Pointer to **string** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 
**Currency** | Pointer to **string** |  | [optional] 
**Containers** | Pointer to [**[]GetBillingInstancesIdentifier200ResponseAllOfBillingInfoContainersInner**](GetBillingInstancesIdentifier200ResponseAllOfBillingInfoContainersInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetBillingInstancesIdentifier200ResponseAllOfBillingInfo{
    // Set fields directly
}
```

### InstanceId (Nullable)

Use the Nullable wrapper methods:
- `obj.InstanceId.IsSet()` — check if set
- `obj.InstanceId.Get()` — get the inner value (returns pointer)
- `obj.InstanceId.Set(&val)` — set the value
- `obj.InstanceId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


