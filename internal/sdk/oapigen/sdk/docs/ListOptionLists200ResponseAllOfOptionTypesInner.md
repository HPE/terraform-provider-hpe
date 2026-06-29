# ListOptionLists200ResponseAllOfOptionTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**SourceUrl** | Pointer to **string** |  | [optional] 
**SourceMethod** | Pointer to **string** |  | [optional] 
**ApiType** | Pointer to **NullableString** |  | [optional] 
**IgnoreSSLErrors** | Pointer to **bool** |  | [optional] 
**RealTime** | Pointer to **bool** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Config** | Pointer to [**ListOptionLists200ResponseAllOfOptionTypesInnerConfig**](ListOptionLists200ResponseAllOfOptionTypesInnerConfig.md) |  | [optional] 
**Credential** | Pointer to [**ListOptionLists200ResponseAllOfOptionTypesInnerCredential**](ListOptionLists200ResponseAllOfOptionTypesInnerCredential.md) |  | [optional] 
**ServiceUsername** | Pointer to **NullableString** |  | [optional] 
**ServicePassword** | Pointer to **NullableString** |  | [optional] 
**InitialDataset** | Pointer to **NullableString** |  | [optional] 
**TranslationScript** | Pointer to **string** |  | [optional] 
**RequestScript** | Pointer to **NullableString** |  | [optional] 
**Account** | Pointer to [**NullableListOptionLists200ResponseAllOfOptionTypesInnerAccount**](ListOptionLists200ResponseAllOfOptionTypesInnerAccount.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListOptionLists200ResponseAllOfOptionTypesInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
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
### RequestScript (Nullable)

Use the Nullable wrapper methods:
- `obj.RequestScript.IsSet()` — check if set
- `obj.RequestScript.Get()` — get the inner value (returns pointer)
- `obj.RequestScript.Set(&val)` — set the value
- `obj.RequestScript.Unset()` — clear the value
### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


