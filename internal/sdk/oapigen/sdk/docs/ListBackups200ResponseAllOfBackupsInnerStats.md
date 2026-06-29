# ListBackups200ResponseAllOfBackupsInnerStats

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TotalSize** | Pointer to **int64** | Total size of all backups in bytes | [optional] 
**AvgSize** | Pointer to **int64** | Average size of each backup in bytes | [optional] 
**TotalCompleted** | Pointer to **int64** | Total completed count | [optional] 
**Success** | Pointer to **int64** | Successful backup count | [optional] 
**Failed** | Pointer to **int64** | Failed backup count | [optional] 
**SuccessRate** | Pointer to **float64** | Success rate 0-1 | [optional] 
**FailRate** | Pointer to **float64** | Failure rate 0-1 | [optional] 
**LastFiveResults** | Pointer to **[]string** | List of the last 5 backup result statuses | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListBackups200ResponseAllOfBackupsInnerStats{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


