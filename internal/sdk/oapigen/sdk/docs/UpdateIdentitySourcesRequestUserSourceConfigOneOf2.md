# UpdateIdentitySourcesRequestUserSourceConfigOneOf2

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Url** | Pointer to **string** | AD Server IP/FQDN | [optional] 
**Domain** | Pointer to **string** | Domain | [optional] 
**UseSSL** | Pointer to **string** | Use SSL (\&quot;on\&quot; or \&quot;off\&quot;) | [optional] 
**BindingUsername** | Pointer to **string** | Binding Username | [optional] 
**BindingPassword** | Pointer to **string** | Binding Password | [optional] 
**RequiredGroup** | Pointer to **string** | Required Group | [optional] 
**SearchMemberGroups** | Pointer to **bool** | Include Member Groups | [optional] [default to false]

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateIdentitySourcesRequestUserSourceConfigOneOf2{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


