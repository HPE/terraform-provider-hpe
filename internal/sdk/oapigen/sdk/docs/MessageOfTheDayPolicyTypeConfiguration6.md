# MessageOfTheDayPolicyTypeConfiguration6

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MotdTitle** | Pointer to **NullableString** |  | [optional] 
**Motd** | [**MessageOfTheDayPolicyTypeConfiguration6Motd**](MessageOfTheDayPolicyTypeConfiguration6Motd.md) |  | 
**MotdMessage** | Pointer to **string** |  | [optional] 
**MotdType** | Pointer to **string** |  | [optional] 
**MotdFullPage** | Pointer to [**MessageOfTheDayPolicyTypeConfiguration6MotdFullPage**](MessageOfTheDayPolicyTypeConfiguration6MotdFullPage.md) |  | [optional] 
**MotdDate** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &MessageOfTheDayPolicyTypeConfiguration6{
    // Set fields directly
}
```

### MotdTitle (Nullable)

Use the Nullable wrapper methods:
- `obj.MotdTitle.IsSet()` — check if set
- `obj.MotdTitle.Get()` — get the inner value (returns pointer)
- `obj.MotdTitle.Set(&val)` — set the value
- `obj.MotdTitle.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


