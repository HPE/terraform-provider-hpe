# WhitelabelSettingsUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** | Can be used to enable / disable whitelabel feature | [optional] 
**ApplianceName** | Pointer to **string** | Appliance name. Master account only | [optional] 
**DisableSupportMenu** | Pointer to **bool** | Can be used to disable support menu | [optional] 
**ResetHeaderLogo** | Pointer to **bool** | Resets header logo to default header logo | [optional] 
**ResetFooterLogo** | Pointer to **bool** | Resets footer logo to default footer logo | [optional] 
**ResetLoginLogo** | Pointer to **bool** | Resets login logo to default login logo | [optional] 
**ResetFavicon** | Pointer to **bool** | Resets favicon to default favicon | [optional] 
**HeaderBgColor** | Pointer to **string** | Header background color | [optional] 
**HeaderFgColor** | Pointer to **string** | Header foreground color | [optional] 
**NavBgColor** | Pointer to **string** | Nav background color | [optional] 
**NavFgColor** | Pointer to **string** | Nav foreground color | [optional] 
**NavHoverColor** | Pointer to **string** | Nav hover color | [optional] 
**PrimaryButtonBgColor** | Pointer to **string** | Primary button background color | [optional] 
**PrimaryButtonFgColor** | Pointer to **string** | Primary button foreground color | [optional] 
**PrimaryButtonHoverBgColor** | Pointer to **string** | Primary button hover background color | [optional] 
**PrimaryButtonHoverFgColor** | Pointer to **string** | Primary button hover foreground color | [optional] 
**FooterBgColor** | Pointer to **string** | Footer background color | [optional] 
**FooterFgColor** | Pointer to **string** | Footer foreground color | [optional] 
**LoginBgColor** | Pointer to **string** | Login background color | [optional] 
**CopyrightString** | Pointer to **string** | Copyright String | [optional] 
**OverrideCss** | Pointer to **string** | Override CSS | [optional] 
**TermsOfUse** | Pointer to **string** | Terms of use content | [optional] 
**PrivacyPolicy** | Pointer to **string** | Privacy policy content | [optional] 
**SupportMenuLinks** | Pointer to [**[]WhitelabelSettingsUpdateSupportMenuLinksInner**](WhitelabelSettingsUpdateSupportMenuLinksInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &WhitelabelSettingsUpdate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


