# NetworkServerCreateNSX

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** | Type Code (nsx-t) | 
**Name** | **string** | Name | 
**Enabled** | Pointer to **bool** | Can be used to enable / disable the network server. | [optional] [default to true]
**ServiceUrl** | **NullableString** | URL | 
**ServiceUsername** | Pointer to **NullableString** | Username | [optional] 
**ServicePassword** | Pointer to **NullableString** | Password | [optional] 
**Config** | Pointer to [**NetworkServerCreateNSXConfig**](NetworkServerCreateNSXConfig.md) |  | [optional] 
**Visibility** | Pointer to **string** | Visibility - Set to public to allow all tenants | [optional] [default to "private"]
**ZoneId** | **int64** | ID of the Cloud that the NSX network server integration should be attached to. | 
**Credential** | Pointer to [**NetworkServerCreateNSXCredential**](NetworkServerCreateNSXCredential.md) |  | [optional] 
**Tenants** | Pointer to [**[]NetworkServerCreateNSXTenantsInner**](NetworkServerCreateNSXTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkServerCreateNSX{
    // Set fields directly
}
```

### ServiceUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceUrl.IsSet()` — check if set
- `obj.ServiceUrl.Get()` — get the inner value (returns pointer)
- `obj.ServiceUrl.Set(&val)` — set the value
- `obj.ServiceUrl.Unset()` — clear the value
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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


