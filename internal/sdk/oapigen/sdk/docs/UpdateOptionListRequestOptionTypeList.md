# UpdateOptionListRequestOptionTypeList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **NullableString** | Description | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Type** | Pointer to **string** | Option List Type eg. &#x60;rest&#x60;, &#x60;api&#x60;, &#x60;ldap&#x60; or &#x60;manual&#x60;. | [optional] [default to "rest"]
**SourceUrl** | Pointer to **string** | Source URL, the http(s) URL to request data from. Required when type is rest. | [optional] 
**Visibility** | Pointer to **string** | Visibility | [optional] [default to "private"]
**SourceMethod** | Pointer to **string** | Source Method, the HTTP method. | [optional] [default to "GET"]
**ApiType** | Pointer to **NullableString** | Api Type, The code of the api option list to use, eg. clouds, environments, groups, instances, instance-wiki, networks, servicePlans, resourcePools, securityGroups, servers, server-wiki. Required when type is api. | [optional] 
**IgnoreSSLErrors** | Pointer to **bool** | Ignore SSL Errors. | [optional] [default to false]
**RealTime** | Pointer to **bool** | Real Time. | [optional] [default to false]
**Credential** | Pointer to [**UpdateOptionListRequestOptionTypeListCredential**](UpdateOptionListRequestOptionTypeListCredential.md) |  | [optional] 
**ServiceUsername** | Pointer to **NullableString** | Username for authenticating with Basic Auth when type is rest or ldap username. | [optional] 
**ServicePassword** | Pointer to **NullableString** | Password for authenticating with Basic Auth when type is rest or ldap password. | [optional] 
**InitialDataset** | Pointer to **NullableString** | Initial Dataset. Create an initial JSON or CSV dataset to be used as the collection for this option list. It should be a list containing objects with properties &#39;name&#39;, and &#39;value&#39;. Required when type is manual. | [optional] 
**TranslationScript** | Pointer to **NullableString** | Translation Script. Create a js script to translate the result data object into an Array containing objects with properties &#39;name&#39; and &#39;value&#39;. The input data is provided as data and the result should be put on the global variable results. | [optional] 
**RequestScript** | Pointer to **NullableString** | Request Script. Create a js script to prepare the request. Return a data object as the body for a post, and return an array containing properties &#39;name&#39; and &#39;value&#39; for a get. The input data is provided as data and the result should be put on the global variable results. | [optional] 
**Config** | Pointer to [**UpdateOptionListRequestOptionTypeListConfig**](UpdateOptionListRequestOptionTypeListConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateOptionListRequestOptionTypeList{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### ApiType (Nullable)

Use the Nullable wrapper methods:
- `obj.ApiType.IsSet()` — check if set
- `obj.ApiType.Get()` — get the inner value (returns pointer)
- `obj.ApiType.Set(&val)` — set the value
- `obj.ApiType.Unset()` — clear the value
### ServiceUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceUsername.IsSet()` — check if set
- `obj.ServiceUsername.Get()` — get the inner value (returns pointer)
- `obj.ServiceUsername.Set(&val)` — set the value
- `obj.ServiceUsername.Unset()` — clear the value
### ServicePassword (Nullable)

Use the Nullable wrapper methods:
- `obj.ServicePassword.IsSet()` — check if set
- `obj.ServicePassword.Get()` — get the inner value (returns pointer)
- `obj.ServicePassword.Set(&val)` — set the value
- `obj.ServicePassword.Unset()` — clear the value
### InitialDataset (Nullable)

Use the Nullable wrapper methods:
- `obj.InitialDataset.IsSet()` — check if set
- `obj.InitialDataset.Get()` — get the inner value (returns pointer)
- `obj.InitialDataset.Set(&val)` — set the value
- `obj.InitialDataset.Unset()` — clear the value
### TranslationScript (Nullable)

Use the Nullable wrapper methods:
- `obj.TranslationScript.IsSet()` — check if set
- `obj.TranslationScript.Get()` — get the inner value (returns pointer)
- `obj.TranslationScript.Set(&val)` — set the value
- `obj.TranslationScript.Unset()` — clear the value
### RequestScript (Nullable)

Use the Nullable wrapper methods:
- `obj.RequestScript.IsSet()` — check if set
- `obj.RequestScript.Get()` — get the inner value (returns pointer)
- `obj.RequestScript.Set(&val)` — set the value
- `obj.RequestScript.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


