# GetJobs200ResponseAllOfJobAnyOf1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Type** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf1Type**](GetJobs200ResponseAllOfJobAnyOf1Type.md) |  | [optional] 
**Task** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf1Task**](GetJobs200ResponseAllOfJobAnyOf1Task.md) |  | [optional] 
**JobSummary** | Pointer to **NullableString** |  | [optional] 
**ScheduleMode** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf1ScheduleMode**](GetJobs200ResponseAllOfJobAnyOf1ScheduleMode.md) |  | [optional] 
**DateTime** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **NullableString** |  | [optional] 
**Namespace** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**LastRun** | Pointer to **time.Time** |  | [optional] 
**LastResult** | Pointer to **string** |  | [optional] 
**CreatedBy** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf1CreatedBy**](GetJobs200ResponseAllOfJobAnyOf1CreatedBy.md) |  | [optional] 
**TargetType** | Pointer to **string** |  | [optional] 
**Targets** | Pointer to [**[]GetJobs200ResponseAllOfJobAnyOf1TargetsInner**](GetJobs200ResponseAllOfJobAnyOf1TargetsInner.md) |  | [optional] 
**CustomConfig** | Pointer to **NullableString** |  | [optional] 
**CustomOptions** | Pointer to [**GetJobs200ResponseAllOfJobAnyOf1CustomOptions**](GetJobs200ResponseAllOfJobAnyOf1CustomOptions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetJobs200ResponseAllOfJobAnyOf1{
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
### Targets (Nullable)

Use the Nullable wrapper methods:
- `obj.Targets.IsSet()` — check if set
- `obj.Targets.Get()` — get the inner value (returns pointer)
- `obj.Targets.Set(&val)` — set the value
- `obj.Targets.Unset()` — clear the value
### CustomConfig (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomConfig.IsSet()` — check if set
- `obj.CustomConfig.Get()` — get the inner value (returns pointer)
- `obj.CustomConfig.Set(&val)` — set the value
- `obj.CustomConfig.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


