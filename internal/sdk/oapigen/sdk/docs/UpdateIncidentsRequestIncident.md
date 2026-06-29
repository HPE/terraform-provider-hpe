# UpdateIncidentsRequestIncident

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Resolution** | Pointer to **string** | Description of the resolution to this incident | [optional] 
**Comment** | Pointer to **string** | Comment on this incident, updates summary field | [optional] 
**Status** | Pointer to **string** | Set status | [optional] 
**Severity** | Pointer to **string** | Set severity | [optional] 
**Name** | Pointer to **string** | Set display name | [optional] 
**StartDate** | Pointer to **time.Time** | Set start time | [optional] 
**EndDate** | Pointer to **time.Time** | Set start time | [optional] 
**InUptime** | Pointer to **bool** | Set &#39;In Availability&#39; | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateIncidentsRequestIncident{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


