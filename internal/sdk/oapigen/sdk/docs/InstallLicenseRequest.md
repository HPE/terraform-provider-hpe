# InstallLicenseRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**License** | **string** | License Key. This is a long and unique string of your Morpheus license. License keys can be requested from the [Morpheus Hub](https://morpheushub.com). | 
**InstallAction** | Pointer to **string** | Install Action can be passed as &#39;add&#39; to stack the license. By default all currently installed licenses are removed and replaced by the new one. | [optional] [default to "replace"]

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstallLicenseRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


