# GetVDIApps200ResponseVdiApp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**LaunchPrefix** | Pointer to **string** |  | [optional] 
**IconPath** | Pointer to **NullableString** |  | [optional] 
**Logo** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetVDIApps200ResponseVdiApp{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### IconPath (Nullable)

Use the Nullable wrapper methods:
- `obj.IconPath.IsSet()` — check if set
- `obj.IconPath.Get()` — get the inner value (returns pointer)
- `obj.IconPath.Set(&val)` — set the value
- `obj.IconPath.Unset()` — clear the value
### Logo (Nullable)

Use the Nullable wrapper methods:
- `obj.Logo.IsSet()` — check if set
- `obj.Logo.Get()` — get the inner value (returns pointer)
- `obj.Logo.Set(&val)` — set the value
- `obj.Logo.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


