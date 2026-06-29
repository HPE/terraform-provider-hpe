# RefreshCloudsRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Mode** | Pointer to **string** | Refresh Mode. Run the &#x60;daily&#x60; or &#x60;costing&#x60; job instead of the default &#x60;hourly&#x60; refresh job. | [optional] 
**Rebuild** | Pointer to **string** | Rebuild. Pass &#x60;true&#x60; to purge existing invoices for the period before refreshing. Only applies to mode &#x60;costing&#x60;. | [optional] 
**Period** | Pointer to **string** | Period. Invoice billing period to refresh in the format &#x60;YYYYMM&#x60;. The default period is the current month. Only applies to mode &#x60;costing&#x60;. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &RefreshCloudsRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


