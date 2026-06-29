# ListBackupJobs200ResponseAllOfJobsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Backup ID | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**Schedule** | Pointer to [**ListBackupJobs200ResponseAllOfJobsInnerSchedule**](ListBackupJobs200ResponseAllOfJobsInnerSchedule.md) |  | [optional] 
**RetentionCount** | Pointer to **NullableInt64** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**BackupProvider** | Pointer to [**ListBackupJobs200ResponseAllOfJobsInnerBackupProvider**](ListBackupJobs200ResponseAllOfJobsInnerBackupProvider.md) |  | [optional] 
**BackupRespository** | Pointer to [**ListBackupJobs200ResponseAllOfJobsInnerBackupRespository**](ListBackupJobs200ResponseAllOfJobsInnerBackupRespository.md) |  | [optional] 
**CronExpression** | Pointer to **NullableString** | Cron Expression | [optional] 
**NextFire** | Pointer to **NullableTime** | Next Fire is the datetime the job will next occur on according to its schedule | [optional] 
**Source** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Account** | Pointer to [**ListBackupJobs200ResponseAllOfJobsInnerAccount**](ListBackupJobs200ResponseAllOfJobsInnerAccount.md) |  | [optional] 
**Enabled** | Pointer to **bool** | Enabled | [optional] 
**DateCreated** | Pointer to **time.Time** | Date Created | [optional] 
**LastUpdated** | Pointer to **time.Time** | Last Updated | [optional] 
**Backups** | Pointer to [**[]ListBackupJobs200ResponseAllOfJobsInnerBackupsInner**](ListBackupJobs200ResponseAllOfJobsInnerBackupsInner.md) | Backups associated with this job | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListBackupJobs200ResponseAllOfJobsInner{
    // Set fields directly
}
```

### RetentionCount (Nullable)

Use the Nullable wrapper methods:
- `obj.RetentionCount.IsSet()` — check if set
- `obj.RetentionCount.Get()` — get the inner value (returns pointer)
- `obj.RetentionCount.Set(&val)` — set the value
- `obj.RetentionCount.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### CronExpression (Nullable)

Use the Nullable wrapper methods:
- `obj.CronExpression.IsSet()` — check if set
- `obj.CronExpression.Get()` — get the inner value (returns pointer)
- `obj.CronExpression.Set(&val)` — set the value
- `obj.CronExpression.Unset()` — clear the value
### NextFire (Nullable)

Use the Nullable wrapper methods:
- `obj.NextFire.IsSet()` — check if set
- `obj.NextFire.Get()` — get the inner value (returns pointer)
- `obj.NextFire.Set(&val)` — set the value
- `obj.NextFire.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


