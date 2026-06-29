# TokenRefreshTokenRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientId** | **string** | Client ID | [default to "morph-api"]
**GrantType** | **string** | OAuth Grant Type, use &#x60;refresh_token&#x60;. | [default to "refresh_token"]
**RefreshToken** | **string** | Refresh Token | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &TokenRefreshTokenRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


