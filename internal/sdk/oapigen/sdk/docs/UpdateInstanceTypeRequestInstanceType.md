# UpdateInstanceTypeRequestInstanceType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Instance type name | [optional] 
**Description** | Pointer to **string** | Instance type description | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Code** | Pointer to **string** | Instance type code | [optional] 
**Category** | Pointer to **string** | Category | [optional] 
**Visibility** | Pointer to **string** | Visibility | [optional] [default to "private"]
**Featured** | Pointer to **bool** | Featured | [optional] 
**HasSettings** | Pointer to **bool** | Enable Settings | [optional] 
**HasAutoScale** | Pointer to **bool** | Enable Scaling (Horizontal) | [optional] 
**HasDeployment** | Pointer to **bool** | Supports Deployments | [optional] 
**EnvironmentPrefix** | Pointer to **string** | Environment Prefix, can be used to make exported evars unique. | [optional] 
**EnvironmentVariables** | Pointer to [**[]UpdateInstanceTypeRequestInstanceTypeEnvironmentVariablesInner**](UpdateInstanceTypeRequestInstanceTypeEnvironmentVariablesInner.md) | Array of instance type env variables. | [optional] 
**PriceSets** | Pointer to [**[]UpdateInstanceTypeRequestInstanceTypePriceSetsInner**](UpdateInstanceTypeRequestInstanceTypePriceSetsInner.md) | Array of price set objects | [optional] 
**OptionTypes** | Pointer to **[]int64** | Array of instance type option type IDs | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateInstanceTypeRequestInstanceType{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


