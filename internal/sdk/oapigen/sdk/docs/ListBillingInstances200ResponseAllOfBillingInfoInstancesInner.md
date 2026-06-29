# ListBillingInstances200ResponseAllOfBillingInfoInstancesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InstanceId** | Pointer to **int64** |  | [optional] 
**InstanceUUID** | Pointer to **string** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 
**Currency** | Pointer to **string** |  | [optional] 
**Containers** | Pointer to [**[]ListBillingInstances200ResponseAllOfBillingInfoInstancesInnerContainersInner**](ListBillingInstances200ResponseAllOfBillingInfoInstancesInnerContainersInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListBillingInstances200ResponseAllOfBillingInfoInstancesInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


