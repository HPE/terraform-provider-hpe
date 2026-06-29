# ClusterJobs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Type** | Pointer to [**ClusterJobsType**](ClusterJobsType.md) |  | [optional] 
**JobSummary** | Pointer to **NullableString** |  | [optional] 
**ScheduleMode** | Pointer to **NullableString** |  | [optional] 
**DateTime** | Pointer to **NullableTime** |  | [optional] 
**Status** | Pointer to **NullableString** |  | [optional] 
**Namespace** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**LastRun** | Pointer to **NullableTime** |  | [optional] 
**LastResult** | Pointer to **NullableString** |  | [optional] 
**CreatedBy** | Pointer to [**ClusterJobsCreatedBy**](ClusterJobsCreatedBy.md) |  | [optional] 
**TargetType** | Pointer to **NullableString** |  | [optional] 
**Targets** | Pointer to **[]map[string]interface{}** |  | [optional] 
**CustomConfig** | Pointer to **map[string]interface{}** |  | [optional] 
**CustomOptions** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterJobs{
    // Set fields directly
}
```

### JobSummary (Nullable)

Use the Nullable wrapper methods:
- `obj.JobSummary.IsSet()` — check if set
- `obj.JobSummary.Get()` — get the inner value (returns pointer)
- `obj.JobSummary.Set(&val)` — set the value
- `obj.JobSummary.Unset()` — clear the value
### ScheduleMode (Nullable)

Use the Nullable wrapper methods:
- `obj.ScheduleMode.IsSet()` — check if set
- `obj.ScheduleMode.Get()` — get the inner value (returns pointer)
- `obj.ScheduleMode.Set(&val)` — set the value
- `obj.ScheduleMode.Unset()` — clear the value
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
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### LastRun (Nullable)

Use the Nullable wrapper methods:
- `obj.LastRun.IsSet()` — check if set
- `obj.LastRun.Get()` — get the inner value (returns pointer)
- `obj.LastRun.Set(&val)` — set the value
- `obj.LastRun.Unset()` — clear the value
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
### CustomConfig (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomConfig.IsSet()` — check if set
- `obj.CustomConfig.Get()` — get the inner value (returns pointer)
- `obj.CustomConfig.Set(&val)` — set the value
- `obj.CustomConfig.Unset()` — clear the value
### CustomOptions (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomOptions.IsSet()` — check if set
- `obj.CustomOptions.Get()` — get the inner value (returns pointer)
- `obj.CustomOptions.Set(&val)` — set the value
- `obj.CustomOptions.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


