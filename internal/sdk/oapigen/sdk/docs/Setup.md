# Setup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** |  | [optional] 
**BuildVersion** | Pointer to **string** | Morpheus build version that the server is running. | [optional] 
**Uuid** | Pointer to **string** | The Appliance Unique ID that is auto generated. | [optional] 
**ApplianceUrl** | Pointer to **string** | The Appliance Server URL as defined under Appliance Settings. | [optional] 
**SetupNeeded** | Pointer to **bool** | Flag to determine if the appliance has been setup, only true when appliance is a fresh install and has not been initialized. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &Setup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


