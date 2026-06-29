# ProvisioningLicense

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**LicenseType** | Pointer to [**GetProvisioningLicense200ResponseLicenseLicenseType**](GetProvisioningLicense200ResponseLicenseLicenseType.md) |  | [optional] 
**LicenseKey** | Pointer to **string** |  | [optional] 
**OrgName** | Pointer to **string** |  | [optional] 
**FullName** | Pointer to **string** |  | [optional] 
**LicenseVersion** | Pointer to **string** |  | [optional] 
**Copies** | Pointer to **int64** |  | [optional] 
**ReservationCount** | Pointer to **int64** |  | [optional] 
**Tenants** | Pointer to **[]map[string]interface{}** |  | [optional] 
**VirtualImages** | Pointer to [**[]GetProvisioningLicense200ResponseLicenseVirtualImagesInner**](GetProvisioningLicense200ResponseLicenseVirtualImagesInner.md) |  | [optional] 
**Account** | Pointer to [**GetProvisioningLicense200ResponseLicenseAccount**](GetProvisioningLicense200ResponseLicenseAccount.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ProvisioningLicense{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


