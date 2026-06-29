# HealthDatabaseInnodbStats

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LargeMemory** | Pointer to **int64** |  | [optional] 
**DictionaryMemory** | Pointer to **int64** |  | [optional] 
**BufferPoolSize** | Pointer to **int64** |  | [optional] 
**FreeBuffers** | Pointer to **int64** |  | [optional] 
**DatabasePages** | Pointer to **int64** |  | [optional] 
**OldPages** | Pointer to **int64** |  | [optional] 
**PendingReads** | Pointer to **int64** |  | [optional] 
**InsertsPerSecond** | Pointer to **float32** |  | [optional] 
**UpdatesPerSecond** | Pointer to **float32** |  | [optional] 
**DeletesPerSecond** | Pointer to **float32** |  | [optional] 
**ReadsPerSecond** | Pointer to **float32** |  | [optional] 
**BufferHitRate** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &HealthDatabaseInnodbStats{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


