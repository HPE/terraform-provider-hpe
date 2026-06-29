# AddCloudResourcePool200ResponseResourcePoolAllOfConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CidrBlock** | **string** | Provide the base CIDR Block to use for this VPC (must be between a /16 and /28 Block) | 
**Tenancy** | **string** |  | [default to "default"]
**ProjectId** | Pointer to **string** | Project ID can have lowercase letters, digits, or hyphens. It must start with a lowercase letter and end with a letter or number.  | [optional] 
**Parent** | **interface{}** |  | 
**BillingAccount** | **string** |  | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCloudResourcePool200ResponseResourcePoolAllOfConfig{
    // Set fields directly
}
```

### Parent (Nullable)

Use the Nullable wrapper methods:
- `obj.Parent.IsSet()` — check if set
- `obj.Parent.Get()` — get the inner value (returns pointer)
- `obj.Parent.Set(&val)` — set the value
- `obj.Parent.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


