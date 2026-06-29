# AddServicePlansRequestServicePlan

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Service plan name | 
**Code** | **string** | Service plan code, must be unique | 
**Description** | Pointer to **string** | Service plan description | [optional] 
**RegionCode** | Pointer to **string** | Service plan region code | [optional] 
**Editable** | Pointer to **bool** | Can be used to enable / disable the editability of the service plan. | [optional] [default to true]
**Visibility** | Pointer to **string** | Can be used to enable / disable the visibility of the service plan, defaults to \&quot;public\&quot; unless the account is not a masterAccount in which the default is \&quot;private\&quot; | [optional] 
**MaxStorage** | **int64** | Max storage size in bytes | 
**MaxMemory** | **int64** | Max memory size in bytes | 
**MaxCores** | Pointer to **int64** | Max number of cores | [optional] 
**MaxCpu** | Pointer to **int64** | Max number of CPUs | [optional] 
**CoresPerSocket** | Pointer to **int64** | Number of cores per CPU | [optional] 
**MaxGpus** | Pointer to **int64** | Max number of GPUs | [optional] 
**MaxDisks** | Pointer to **int64** | Max disks allowed | [optional] 
**ProvisionType** | [**AddServicePlansRequestServicePlanProvisionType**](AddServicePlansRequestServicePlanProvisionType.md) |  | 
**CustomCpu** | Pointer to **bool** | Can be used to enable / disable customizable cpu | [optional] [default to false]
**CustomCores** | Pointer to **bool** | Can be used to enable / disable customizable cores | [optional] [default to false]
**CustomMaxStorage** | Pointer to **bool** | Can be used to enable / disable customizable storage | [optional] [default to false]
**CustomMaxDataStorage** | Pointer to **bool** | Can be used to enable / disable customizable extra volumes. | [optional] [default to false]
**CustomMaxMemory** | Pointer to **bool** | Can be used to enable / disable customizable memory. | [optional] [default to false]
**AddVolumes** | Pointer to **bool** | Can be used to enable / disable ability to add volumes | [optional] [default to false]
**SortOrder** | Pointer to **int64** | Sort order | [optional] 
**PriceSets** | Pointer to [**[]AddServicePlansRequestServicePlanPriceSetsInner**](AddServicePlansRequestServicePlanPriceSetsInner.md) | List of price sets to include in service plan. | [optional] 
**Permissions** | Pointer to [**AddServicePlansRequestServicePlanPermissions**](AddServicePlansRequestServicePlanPermissions.md) |  | [optional] 
**Config** | Pointer to [**AddServicePlansRequestServicePlanConfig**](AddServicePlansRequestServicePlanConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddServicePlansRequestServicePlan{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


