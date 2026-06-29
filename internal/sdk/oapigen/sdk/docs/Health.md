# Health

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** |  | [optional] 
**StatusMessage** | Pointer to **string** |  | [optional] 
**ApplianceUrl** | Pointer to **string** |  | [optional] 
**BuildVersion** | Pointer to **string** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**SetupNeeded** | Pointer to **bool** |  | [optional] 
**Date** | Pointer to **time.Time** |  | [optional] 
**Cpu** | Pointer to [**HealthCpu**](HealthCpu.md) |  | [optional] 
**Memory** | Pointer to [**HealthMemory**](HealthMemory.md) |  | [optional] 
**Threads** | Pointer to [**HealthThreads**](HealthThreads.md) |  | [optional] 
**Database** | Pointer to [**HealthDatabase**](HealthDatabase.md) |  | [optional] 
**Elastic** | Pointer to [**HealthElastic**](HealthElastic.md) |  | [optional] 
**Rabbit** | Pointer to [**HealthRabbit**](HealthRabbit.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &Health{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


