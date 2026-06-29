# HealthLogs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TypeCode** | Pointer to **string** |  | [optional] 
**Ts** | Pointer to **time.Time** |  | [optional] 
**Level** | Pointer to **string** |  | [optional] 
**SourceType** | Pointer to **string** |  | [optional] 
**Message** | Pointer to **string** |  | [optional] 
**Hostname** | Pointer to **string** |  | [optional] 
**Title** | Pointer to **string** |  | [optional] 
**LogSignature** | Pointer to **string** |  | [optional] 
**ObjectId** | Pointer to **string** |  | [optional] 
**Seq** | Pointer to **int64** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**SignatureVerified** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &HealthLogs{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


