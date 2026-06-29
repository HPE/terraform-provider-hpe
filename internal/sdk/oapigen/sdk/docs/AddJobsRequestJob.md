# AddJobsRequestJob

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the Job | 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Enabled** | Pointer to **bool** | Use this to set enabled state | [optional] [default to true]
**Task** | [**TaskJobPayloadTask**](TaskJobPayloadTask.md) |  | 
**Workflow** | [**WorkflowJobPayloadWorkflow**](WorkflowJobPayloadWorkflow.md) |  | 
**TargetType** | **string** | Target type where job will execute | 
**Targets** | [**[]SecurityScanJobTargetsInner**](SecurityScanJobTargetsInner.md) |  | 
**InstanceLabel** | Pointer to **string** | Instance Label. Only applicable if &#x60;targetType&#x60; is &#x60;instance-label&#x60;. | [optional] 
**ServerLabel** | Pointer to **string** | Server Label. Only applicable if &#x60;targetType&#x60; is &#x60;server-label&#x60;. | [optional] 
**ScheduleMode** | [**SecurityScanJobScheduleMode**](SecurityScanJobScheduleMode.md) |  | 
**CustomOptions** | Pointer to **map[string]interface{}** | Map of options to be used as values in the workflow tasks. These correspond to option types. | [optional] 
**CustomConfig** | Pointer to **string** | Job custom configuration (String in JSON format) | [optional] 
**DateTime** | Pointer to **time.Time** | Date and Time to execute the job. Use UTC time in the format 2020-02-15T05:00:00Z. Required when scheduleMode is &#39;dateTime&#39;. | [optional] 
**Run** | Pointer to **bool** | If true, executes job | [optional] 
**SecurityPackage** | [**SecurityScanJobSecurityPackage**](SecurityScanJobSecurityPackage.md) |  | 
**ScanPath** | Pointer to **string** | Scan Checklist | [optional] 
**SecurityProfile** | Pointer to **string** | Security Profile | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddJobsRequestJob{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


