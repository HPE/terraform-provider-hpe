# AddTenant200ResponseAllOfAccount

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Subdomain** | Pointer to **string** |  | [optional] 
**Currency** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**CustomerNumber** | Pointer to **NullableString** |  | [optional] 
**AccountNumber** | Pointer to **NullableString** |  | [optional] 
**AccountName** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Master** | Pointer to **bool** |  | [optional] 
**Parent** | Pointer to [**AddTenant200ResponseAllOfAccountParent**](AddTenant200ResponseAllOfAccountParent.md) |  | [optional] 
**Role** | Pointer to [**AddTenant200ResponseAllOfAccountRole**](AddTenant200ResponseAllOfAccountRole.md) |  | [optional] 
**Stats** | Pointer to [**AddTenant200ResponseAllOfAccountStats**](AddTenant200ResponseAllOfAccountStats.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddTenant200ResponseAllOfAccount{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### CustomerNumber (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomerNumber.IsSet()` — check if set
- `obj.CustomerNumber.Get()` — get the inner value (returns pointer)
- `obj.CustomerNumber.Set(&val)` — set the value
- `obj.CustomerNumber.Unset()` — clear the value
### AccountNumber (Nullable)

Use the Nullable wrapper methods:
- `obj.AccountNumber.IsSet()` — check if set
- `obj.AccountNumber.Get()` — get the inner value (returns pointer)
- `obj.AccountNumber.Set(&val)` — set the value
- `obj.AccountNumber.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


