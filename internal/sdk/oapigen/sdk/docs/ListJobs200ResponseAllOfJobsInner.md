# ListJobs200ResponseAllOfJobsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Type** | Pointer to [**ListJobs200ResponseAllOfJobsInnerAnyOf2Type**](ListJobs200ResponseAllOfJobsInnerAnyOf2Type.md) |  | [optional] 
**Workflow** | Pointer to [**ListJobs200ResponseAllOfJobsInnerAnyOf2Workflow**](ListJobs200ResponseAllOfJobsInnerAnyOf2Workflow.md) |  | [optional] 
**Task** | Pointer to [**ListJobs200ResponseAllOfJobsInnerAnyOf1Task**](ListJobs200ResponseAllOfJobsInnerAnyOf1Task.md) |  | [optional] 
**SecurityPackage** | Pointer to [**ListJobs200ResponseAllOfJobsInnerAnyOfSecurityPackage**](ListJobs200ResponseAllOfJobsInnerAnyOfSecurityPackage.md) |  | [optional] 
**JobSummary** | Pointer to **string** |  | [optional] 
**ScheduleMode** | Pointer to [**ListJobs200ResponseAllOfJobsInnerAnyOf2ScheduleMode**](ListJobs200ResponseAllOfJobsInnerAnyOf2ScheduleMode.md) |  | [optional] 
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
**CreatedBy** | Pointer to [**ListJobs200ResponseAllOfJobsInnerAnyOf2CreatedBy**](ListJobs200ResponseAllOfJobsInnerAnyOf2CreatedBy.md) |  | [optional] 
**TargetType** | Pointer to **string** |  | [optional] 
**Targets** | Pointer to [**[]ListJobs200ResponseAllOfJobsInnerAnyOf2TargetsInner**](ListJobs200ResponseAllOfJobsInnerAnyOf2TargetsInner.md) |  | [optional] 
**ScanPath** | Pointer to **string** | Scan Checklist. Only applies to type scap-package. | [optional] 
**SecurityProfile** | Pointer to **string** | Security Profile. Only applies to type scap-package. | [optional] 
**CustomConfig** | Pointer to **string** |  | [optional] 
**CustomOptions** | Pointer to [**ListJobs200ResponseAllOfJobsInnerAnyOf2CustomOptions**](ListJobs200ResponseAllOfJobsInnerAnyOf2CustomOptions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListJobs200ResponseAllOfJobsInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


