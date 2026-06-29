# ListProvisioningLicenses200ResponseAllOfLicensesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**LicenseType** | Pointer to [**ListProvisioningLicenses200ResponseAllOfLicensesInnerLicenseType**](ListProvisioningLicenses200ResponseAllOfLicensesInnerLicenseType.md) |  | [optional] 
**LicenseKey** | Pointer to **string** |  | [optional] 
**OrgName** | Pointer to **string** |  | [optional] 
**FullName** | Pointer to **string** |  | [optional] 
**LicenseVersion** | Pointer to **string** |  | [optional] 
**Copies** | Pointer to **int64** |  | [optional] 
**ReservationCount** | Pointer to **int64** |  | [optional] 
**Tenants** | Pointer to **[]map[string]interface{}** |  | [optional] 
**VirtualImages** | Pointer to [**[]ListProvisioningLicenses200ResponseAllOfLicensesInnerVirtualImagesInner**](ListProvisioningLicenses200ResponseAllOfLicensesInnerVirtualImagesInner.md) |  | [optional] 
**Account** | Pointer to [**NullableListProvisioningLicenses200ResponseAllOfLicensesInnerAccount**](ListProvisioningLicenses200ResponseAllOfLicensesInnerAccount.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListProvisioningLicenses200ResponseAllOfLicensesInner{
    // Set fields directly
}
```

### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


