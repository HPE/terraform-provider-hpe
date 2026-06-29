# SetInstanceSecurityGroups200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SecurityGroups** | Pointer to [**[]SetInstanceSecurityGroups200ResponseAllOfSecurityGroupsInner**](SetInstanceSecurityGroups200ResponseAllOfSecurityGroupsInner.md) | Array of security group objects | [optional] 
**Success** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &SetInstanceSecurityGroups200Response{
    // Set fields directly
}
```

### SecurityGroups (Nullable)

Use the Nullable wrapper methods:
- `obj.SecurityGroups.IsSet()` — check if set
- `obj.SecurityGroups.Get()` — get the inner value (returns pointer)
- `obj.SecurityGroups.Set(&val)` — set the value
- `obj.SecurityGroups.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


