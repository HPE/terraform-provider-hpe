# AddJobs200ResponseAllOfJob

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Type** | Pointer to [**AddJobs200ResponseAllOfJobType**](AddJobs200ResponseAllOfJobType.md) |  | [optional] 
**Workflow** | Pointer to [**AddJobs200ResponseAllOfJobWorkflow**](AddJobs200ResponseAllOfJobWorkflow.md) |  | [optional] 
**Task** | Pointer to [**AddJobs200ResponseAllOfJobTask**](AddJobs200ResponseAllOfJobTask.md) |  | [optional] 
**SecurityPackage** | Pointer to [**AddJobs200ResponseAllOfJobSecurityPackage**](AddJobs200ResponseAllOfJobSecurityPackage.md) |  | [optional] 
**JobSummary** | Pointer to **NullableString** |  | [optional] 
**ScheduleMode** | Pointer to [**AddJobs200ResponseAllOfJobScheduleMode**](AddJobs200ResponseAllOfJobScheduleMode.md) |  | [optional] 
**DateTime** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **NullableString** |  | [optional] 
**Namespace** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**LastRun** | Pointer to **time.Time** |  | [optional] 
**LastResult** | Pointer to **NullableString** |  | [optional] 
**CreatedBy** | Pointer to [**AddJobs200ResponseAllOfJobCreatedBy**](AddJobs200ResponseAllOfJobCreatedBy.md) |  | [optional] 
**TargetType** | Pointer to **NullableString** |  | [optional] 
**Targets** | Pointer to [**[]AddJobs200ResponseAllOfJobTargetsInner**](AddJobs200ResponseAllOfJobTargetsInner.md) |  | [optional] 
**ScanPath** | Pointer to **NullableString** | Scan Checklist. Only applies to type scap-package. | [optional] 
**SecurityProfile** | Pointer to **NullableString** | Security Profile. Only applies to type scap-package. | [optional] 
**CustomConfig** | Pointer to **NullableString** |  | [optional] 
**CustomOptions** | Pointer to [**AddJobs200ResponseAllOfJobCustomOptions**](AddJobs200ResponseAllOfJobCustomOptions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddJobs200ResponseAllOfJob{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### JobSummary (Nullable)

Use the Nullable wrapper methods:
- `obj.JobSummary.IsSet()` — check if set
- `obj.JobSummary.Get()` — get the inner value (returns pointer)
- `obj.JobSummary.Set(&val)` — set the value
- `obj.JobSummary.Unset()` — clear the value
### DateTime (Nullable)

Use the Nullable wrapper methods:
- `obj.DateTime.IsSet()` — check if set
- `obj.DateTime.Get()` — get the inner value (returns pointer)
- `obj.DateTime.Set(&val)` — set the value
- `obj.DateTime.Unset()` — clear the value
### Status (Nullable)

Use the Nullable wrapper methods:
- `obj.Status.IsSet()` — check if set
- `obj.Status.Get()` — get the inner value (returns pointer)
- `obj.Status.Set(&val)` — set the value
- `obj.Status.Unset()` — clear the value
### Namespace (Nullable)

Use the Nullable wrapper methods:
- `obj.Namespace.IsSet()` — check if set
- `obj.Namespace.Get()` — get the inner value (returns pointer)
- `obj.Namespace.Set(&val)` — set the value
- `obj.Namespace.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### LastResult (Nullable)

Use the Nullable wrapper methods:
- `obj.LastResult.IsSet()` — check if set
- `obj.LastResult.Get()` — get the inner value (returns pointer)
- `obj.LastResult.Set(&val)` — set the value
- `obj.LastResult.Unset()` — clear the value
### TargetType (Nullable)

Use the Nullable wrapper methods:
- `obj.TargetType.IsSet()` — check if set
- `obj.TargetType.Get()` — get the inner value (returns pointer)
- `obj.TargetType.Set(&val)` — set the value
- `obj.TargetType.Unset()` — clear the value
### Targets (Nullable)

Use the Nullable wrapper methods:
- `obj.Targets.IsSet()` — check if set
- `obj.Targets.Get()` — get the inner value (returns pointer)
- `obj.Targets.Set(&val)` — set the value
- `obj.Targets.Unset()` — clear the value
### ScanPath (Nullable)

Use the Nullable wrapper methods:
- `obj.ScanPath.IsSet()` — check if set
- `obj.ScanPath.Get()` — get the inner value (returns pointer)
- `obj.ScanPath.Set(&val)` — set the value
- `obj.ScanPath.Unset()` — clear the value
### SecurityProfile (Nullable)

Use the Nullable wrapper methods:
- `obj.SecurityProfile.IsSet()` — check if set
- `obj.SecurityProfile.Get()` — get the inner value (returns pointer)
- `obj.SecurityProfile.Set(&val)` — set the value
- `obj.SecurityProfile.Unset()` — clear the value
### CustomConfig (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomConfig.IsSet()` — check if set
- `obj.CustomConfig.Get()` — get the inner value (returns pointer)
- `obj.CustomConfig.Set(&val)` — set the value
- `obj.CustomConfig.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


