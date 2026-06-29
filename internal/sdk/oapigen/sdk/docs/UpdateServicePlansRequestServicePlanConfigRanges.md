# UpdateServicePlansRequestServicePlanConfigRanges

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MinStorage** | Pointer to **int64** | Custom min storage in GB (unless &#x60;storageSizeType&#x60; set to mb) | [optional] 
**MaxStorage** | Pointer to **int64** | Custom max storage in GB (unless &#x60;storageSizeType&#x60; set to mb) | [optional] 
**MinPerDiskSize** | Pointer to **int64** |  | [optional] 
**MaxPerDiskSize** | Pointer to **int64** |  | [optional] 
**MinMemory** | Pointer to **int64** | Custom min memory in bytes | [optional] 
**MaxMemory** | Pointer to **int64** | Custom max memory in bytes | [optional] 
**MinCores** | Pointer to **int64** | Custom min cores | [optional] 
**MaxCores** | Pointer to **int64** | Custom max cores | [optional] 
**MinSockets** | Pointer to **int64** | Custom min sockets | [optional] 
**MaxSockets** | Pointer to **int64** | Custom max sockets | [optional] 
**MinCoresPerSocket** | Pointer to **int64** | Custom min cores allowed per socket | [optional] 
**MaxCoresPerSocket** | Pointer to **int64** | Custom max cores allowed per socket | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateServicePlansRequestServicePlanConfigRanges{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


