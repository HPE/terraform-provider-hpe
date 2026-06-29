# UpdateMonitoringSettingsRequestMonitoringSettingsServiceNow

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** | Enabled | [optional] 
**Integration** | Pointer to [**UpdateMonitoringSettingsRequestMonitoringSettingsServiceNowIntegration**](UpdateMonitoringSettingsRequestMonitoringSettingsServiceNowIntegration.md) |  | [optional] 
**NewIncidentAction** | Pointer to **string** | New Incident Action | [optional] 
**CloseIncidentAction** | Pointer to **string** | Close Incident Action | [optional] 
**InfoMapping** | Pointer to **string** | Info Mapping | [optional] 
**WarningMapping** | Pointer to **string** | Warning Mapping | [optional] 
**CriticalMapping** | Pointer to **string** | Critical Mapping | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateMonitoringSettingsRequestMonitoringSettingsServiceNow{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


