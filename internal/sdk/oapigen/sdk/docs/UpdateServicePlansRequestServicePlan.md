# UpdateServicePlansRequestServicePlan

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Service plan name | [optional] 
**Code** | Pointer to **string** | Service plan code, must be unique | [optional] 
**Description** | Pointer to **string** | Service plan description | [optional] 
**Editable** | Pointer to **bool** | Can be used to enable / disable the editability of the service plan. | [optional] [default to true]
**MaxStorage** | Pointer to **int64** | Max storage size in bytes | [optional] 
**MaxMemory** | Pointer to **int64** | Max memory size in bytes | [optional] 
**MaxCpu** | Pointer to **int64** |  | [optional] 
**MaxCores** | Pointer to **int64** | Max cores | [optional] 
**MaxDisks** | Pointer to **int64** | Max disks allowed | [optional] 
**ProvisionType** | Pointer to [**UpdateServicePlansRequestServicePlanProvisionType**](UpdateServicePlansRequestServicePlanProvisionType.md) |  | [optional] 
**CoresPerSocket** | Pointer to **int64** |  | [optional] 
**CustomCpu** | Pointer to **bool** | Can be used to enable / disable customizable cpu | [optional] 
**CustomCores** | Pointer to **bool** | Can be used to enable / disable customizable cores | [optional] 
**CustomMaxStorage** | Pointer to **bool** | Can be used to enable / disable customizable storage | [optional] 
**CustomMaxDataStorage** | Pointer to **bool** | Can be used to enable / disable customizable extra volumes. | [optional] 
**CustomMaxMemory** | Pointer to **bool** | Can be used to enable / disable customizable memory. | [optional] 
**AddVolumes** | Pointer to **bool** | Can be used to enable / disable ability to add volumes | [optional] 
**SortOrder** | Pointer to **int64** | Sort order | [optional] 
**PriceSets** | Pointer to [**[]UpdateServicePlansRequestServicePlanPriceSetsInner**](UpdateServicePlansRequestServicePlanPriceSetsInner.md) | List of price sets to include in service plan | [optional] 
**Config** | Pointer to [**UpdateServicePlansRequestServicePlanConfig**](UpdateServicePlansRequestServicePlanConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateServicePlansRequestServicePlan{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


