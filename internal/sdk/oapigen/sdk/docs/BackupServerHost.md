# BackupServerHost

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LocationType** | **string** |  | 
**Name** | **string** | A name for the backup | 
**ServerId** | **int64** | The ID of the server to backup | 
**BackupType** | **string** | The backup type code for server backups | 
**TargetPath** | Pointer to **string** | The file or directory path on the target host to back up. | [optional] 
**SshUsername** | Pointer to **string** | SSH username for connecting to the server. | [optional] 
**SshPassword** | Pointer to **string** | SSH password for connecting to the server. | [optional] 
**Target** | Pointer to **int64** | The ID of the storage bucket to store the backup in. | [optional] 
**JobAction** | **string** | Create a new backup job, clone an existing job or add the new backup to an existing job. If not specified, defaults to &#x60;new&#x60;. | 
**JobId** | Pointer to **int64** | The ID of the job to clone or add to. Only applies to jobAction &#x60;clone&#x60; and &#x60;addTo&#x60;. | [optional] 
**JobName** | Pointer to **string** | Name for new job. Only applies to jobAction &#x60;new&#x60; and &#x60;clone&#x60;. Defaults to the backup name if not specified. | [optional] 
**JobSchedule** | Pointer to **int64** | The ID of the execute schedule for new job. See Execute Schedules. Only applies to jobAction &#x60;new&#x60; and &#x60;clone&#x60;. | [optional] 
**RetentionCount** | Pointer to **int64** | Retention Count for new job. By default the backup settings value will be used. Only applies to jobAction &#x60;new&#x60; and &#x60;clone&#x60;. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BackupServerHost{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


