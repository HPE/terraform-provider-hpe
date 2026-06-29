# ListIdentitySources200ResponseAllOfUserSourcesInnerAnyOfRoleMappingsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SourceRoleName** | Pointer to **NullableString** |  | [optional] 
**SourceRoleFqn** | Pointer to **string** |  | [optional] 
**MappedRole** | Pointer to [**ListIdentitySources200ResponseAllOfUserSourcesInnerAnyOfRoleMappingsInnerMappedRole**](ListIdentitySources200ResponseAllOfUserSourcesInnerAnyOfRoleMappingsInnerMappedRole.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListIdentitySources200ResponseAllOfUserSourcesInnerAnyOfRoleMappingsInner{
    // Set fields directly
}
```

### SourceRoleName (Nullable)

Use the Nullable wrapper methods:
- `obj.SourceRoleName.IsSet()` — check if set
- `obj.SourceRoleName.Get()` — get the inner value (returns pointer)
- `obj.SourceRoleName.Set(&val)` — set the value
- `obj.SourceRoleName.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


