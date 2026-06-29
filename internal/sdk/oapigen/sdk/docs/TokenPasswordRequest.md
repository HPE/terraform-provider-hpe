# TokenPasswordRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientId** | **string** | Client ID | [default to "morph-api"]
**GrantType** | **string** | OAuth Grant Type, use &#x60;password&#x60;. | [default to "password"]
**Scope** | **string** | OAuth token scope, use &#x60;write&#x60;. | 
**Username** | **string** | Username Sub-tenant users must format their username as &#x60;subdomain\\username&#x60; with a prefix that is the tenant subdomain or id by default.  | 
**Password** | **string** | Password | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &TokenPasswordRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


