# SetupRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApplianceName** | **string** | Appliance Name. Choose a name for your Morpheus Appliance.  This is stored in the Appliance Settings. | 
**ApplianceUrl** | Pointer to **string** | Appliance URL. Specify the full URL for your Morpheus Appliance. This is stored in the Appliance Settings. | [optional] 
**AccountName** | **string** | Name of the Master Tenant account being created. | 
**FirstName** | Pointer to **string** | First Name for the System Admin user being created. | [optional] 
**LastName** | Pointer to **string** | Last Name for the System Admin user being created. | [optional] 
**Username** | **string** | Username for the System Admin user being created. | 
**Email** | **string** | Email for the System Admin user being created. | 
**Password** | **string** | Password for the System Admin user being created. | 
**LicenseKey** | Pointer to **string** | License Key to install on setup. By default a trial VME license will be installed if this is not provided. | [optional] 
**Logs** | Pointer to **bool** | Enable Logs | [optional] 
**Monitoring** | Pointer to **bool** | Enable Monitoring | [optional] 
**Backups** | Pointer to **bool** | Enable Backups | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SetupRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


