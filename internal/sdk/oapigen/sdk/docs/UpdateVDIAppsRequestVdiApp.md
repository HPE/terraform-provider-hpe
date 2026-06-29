# UpdateVDIAppsRequestVdiApp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | VDI App name | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**IconPath** | Pointer to ***os.File** | Icon Path. A relative location of an icon image | [optional] 
**LaunchPrefix** | Pointer to **string** | The RDS App Name Prefix.  Must start with &#39;||&#39; | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateVDIAppsRequestVdiApp{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


