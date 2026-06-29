# ListHealth200ResponseAllOfHealth

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
**Cpu** | Pointer to [**ListHealth200ResponseAllOfHealthCpu**](ListHealth200ResponseAllOfHealthCpu.md) |  | [optional] 
**Memory** | Pointer to [**ListHealth200ResponseAllOfHealthMemory**](ListHealth200ResponseAllOfHealthMemory.md) |  | [optional] 
**Threads** | Pointer to [**ListHealth200ResponseAllOfHealthThreads**](ListHealth200ResponseAllOfHealthThreads.md) |  | [optional] 
**Database** | Pointer to [**ListHealth200ResponseAllOfHealthDatabase**](ListHealth200ResponseAllOfHealthDatabase.md) |  | [optional] 
**Elastic** | Pointer to [**ListHealth200ResponseAllOfHealthElastic**](ListHealth200ResponseAllOfHealthElastic.md) |  | [optional] 
**Rabbit** | Pointer to [**ListHealth200ResponseAllOfHealthRabbit**](ListHealth200ResponseAllOfHealthRabbit.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListHealth200ResponseAllOfHealth{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


