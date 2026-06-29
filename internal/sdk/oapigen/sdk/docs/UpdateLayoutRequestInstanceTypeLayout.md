# UpdateLayoutRequestInstanceTypeLayout

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Layout name | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**InstanceVersion** | Pointer to **string** | Version of the layout | [optional] 
**Description** | Pointer to **string** | Layout description | [optional] 
**SortOrder** | Pointer to **int64** | Display order of the layout, higher to lower | [optional] 
**Creatable** | Pointer to **bool** | Can be used to enable / disable the creatability of the layout. | [optional] [default to true]
**ProvisionTypeCode** | Pointer to **string** | Provision type code | [optional] 
**MemoryRequirement** | Pointer to **string** | Memory requirement in megabytes | [optional] 
**HasAutoScale** | Pointer to **bool** | Can be used to enable / disable the horizontal scaling. | [optional] [default to false]
**SupportsConvertToManaged** | Pointer to **bool** | Can be used to enable / disable the supports convert to managed. | [optional] [default to false]
**ContainerTypes** | Pointer to **[]int64** | Array of layout node type IDs | [optional] 
**OptionTypes** | Pointer to **[]int64** | Array of layout option type IDs | [optional] 
**SpecTemplates** | Pointer to **[]int64** | Array of layout spec template IDs | [optional] 
**EnvironmentVariables** | Pointer to [**[]UpdateLayoutRequestInstanceTypeLayoutEnvironmentVariablesInner**](UpdateLayoutRequestInstanceTypeLayoutEnvironmentVariablesInner.md) | The environmentVariables parameter is array of env objects | [optional] 
**PriceSets** | Pointer to [**[]UpdateLayoutRequestInstanceTypeLayoutPriceSetsInner**](UpdateLayoutRequestInstanceTypeLayoutPriceSetsInner.md) | Array of price set objects | [optional] 
**Permissions** | Pointer to [**UpdateLayoutRequestInstanceTypeLayoutPermissions**](UpdateLayoutRequestInstanceTypeLayoutPermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateLayoutRequestInstanceTypeLayout{
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


