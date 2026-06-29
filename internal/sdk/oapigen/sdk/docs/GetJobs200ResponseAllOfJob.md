# GetJobs200ResponseAllOfJob

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Type** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf2Type**](GetJobs200ResponseAllOfJobAnyOf2Type.md) |  | [optional] 
**Workflow** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf2Workflow**](GetJobs200ResponseAllOfJobAnyOf2Workflow.md) |  | [optional] 
**Task** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf1Task**](GetJobs200ResponseAllOfJobAnyOf1Task.md) |  | [optional] 
**SecurityPackage** | Pointer to [**GetJobs200ResponseAllOfJobAnyOfSecurityPackage**](GetJobs200ResponseAllOfJobAnyOfSecurityPackage.md) |  | [optional] 
**JobSummary** | Pointer to **string** |  | [optional] 
**ScheduleMode** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf2ScheduleMode**](GetJobs200ResponseAllOfJobAnyOf2ScheduleMode.md) |  | [optional] 
**DateTime** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Namespace** | Pointer to **string** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**LastRun** | Pointer to **time.Time** |  | [optional] 
**LastResult** | Pointer to **string** |  | [optional] 
**CreatedBy** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf2CreatedBy**](GetJobs200ResponseAllOfJobAnyOf2CreatedBy.md) |  | [optional] 
**TargetType** | Pointer to **string** |  | [optional] 
**Targets** | Pointer to [**[]GetJobs200ResponseAllOfJobAnyOf2TargetsInner**](GetJobs200ResponseAllOfJobAnyOf2TargetsInner.md) |  | [optional] 
**ScanPath** | Pointer to **string** | Scan Checklist. Only applies to type scap-package. | [optional] 
**SecurityProfile** | Pointer to **string** | Security Profile. Only applies to type scap-package. | [optional] 
**CustomConfig** | Pointer to **string** |  | [optional] 
**CustomOptions** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf2CustomOptions**](GetJobs200ResponseAllOfJobAnyOf2CustomOptions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetJobs200ResponseAllOfJob{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


