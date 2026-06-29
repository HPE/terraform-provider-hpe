# UpdateSecurityPackagesRequestSecurityPackage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the security package | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Type** | Pointer to **string** | Security Package Type Code | [optional] [default to "scap-package"]
**Description** | Pointer to **string** | A description for the security package | [optional] 
**Url** | Pointer to **string** | URL to download the security package zip file from | [optional] 
**Enabled** | Pointer to **bool** | Can be used to disable the security package | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateSecurityPackagesRequestSecurityPackage{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


