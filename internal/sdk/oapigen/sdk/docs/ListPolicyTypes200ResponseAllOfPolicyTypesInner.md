# ListPolicyTypes200ResponseAllOfPolicyTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**LoadMethod** | Pointer to **NullableString** |  | [optional] 
**EnforceMethod** | Pointer to **NullableString** |  | [optional] 
**PrepareMethod** | Pointer to **NullableString** |  | [optional] 
**ValidateMethod** | Pointer to **NullableString** |  | [optional] 
**EnforceOnProvision** | Pointer to **bool** |  | [optional] 
**EnforceOnManaged** | Pointer to **bool** |  | [optional] 
**OptionTypes** | Pointer to [**[]ListPolicyTypes200ResponseAllOfPolicyTypesInnerOptionTypesInner**](ListPolicyTypes200ResponseAllOfPolicyTypesInnerOptionTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListPolicyTypes200ResponseAllOfPolicyTypesInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### LoadMethod (Nullable)

Use the Nullable wrapper methods:
- `obj.LoadMethod.IsSet()` — check if set
- `obj.LoadMethod.Get()` — get the inner value (returns pointer)
- `obj.LoadMethod.Set(&val)` — set the value
- `obj.LoadMethod.Unset()` — clear the value
### EnforceMethod (Nullable)

Use the Nullable wrapper methods:
- `obj.EnforceMethod.IsSet()` — check if set
- `obj.EnforceMethod.Get()` — get the inner value (returns pointer)
- `obj.EnforceMethod.Set(&val)` — set the value
- `obj.EnforceMethod.Unset()` — clear the value
### PrepareMethod (Nullable)

Use the Nullable wrapper methods:
- `obj.PrepareMethod.IsSet()` — check if set
- `obj.PrepareMethod.Get()` — get the inner value (returns pointer)
- `obj.PrepareMethod.Set(&val)` — set the value
- `obj.PrepareMethod.Unset()` — clear the value
### ValidateMethod (Nullable)

Use the Nullable wrapper methods:
- `obj.ValidateMethod.IsSet()` — check if set
- `obj.ValidateMethod.Get()` — get the inner value (returns pointer)
- `obj.ValidateMethod.Set(&val)` — set the value
- `obj.ValidateMethod.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


