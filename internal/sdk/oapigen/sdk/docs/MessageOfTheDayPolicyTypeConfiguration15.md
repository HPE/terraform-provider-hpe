# MessageOfTheDayPolicyTypeConfiguration15

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MotdTitle** | Pointer to **string** |  | [optional] 
**Motd** | Pointer to [**MessageOfTheDayPolicyTypeConfiguration15Motd**](MessageOfTheDayPolicyTypeConfiguration15Motd.md) |  | [optional] 
**MotdMessage** | Pointer to **string** |  | [optional] 
**MotdType** | Pointer to **string** |  | [optional] 
**MotdFullPage** | Pointer to **NullableBool** |  | [optional] 
**MotdDate** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &MessageOfTheDayPolicyTypeConfiguration15{
    // Set fields directly
}
```

### MotdFullPage (Nullable)

Use the Nullable wrapper methods:
- `obj.MotdFullPage.IsSet()` — check if set
- `obj.MotdFullPage.Get()` — get the inner value (returns pointer)
- `obj.MotdFullPage.Set(&val)` — set the value
- `obj.MotdFullPage.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


