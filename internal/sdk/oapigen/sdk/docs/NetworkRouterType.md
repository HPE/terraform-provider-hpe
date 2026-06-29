# NetworkRouterType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Creatable** | Pointer to **bool** |  | [optional] 
**Selectable** | Pointer to **bool** |  | [optional] 
**HasFirewall** | Pointer to **bool** |  | [optional] 
**HasDhcp** | Pointer to **bool** |  | [optional] 
**HasRouting** | Pointer to **bool** |  | [optional] 
**HasNetworkServer** | Pointer to **bool** |  | [optional] 
**OptionTypes** | Pointer to [**[]NetworkRouterTypeOptionTypesInner**](NetworkRouterTypeOptionTypesInner.md) |  | [optional] 
**RuleOptionTypes** | Pointer to [**[]NetworkRouterTypeRuleOptionTypesInner**](NetworkRouterTypeRuleOptionTypesInner.md) |  | [optional] 
**RuleGroupOptionTypes** | Pointer to [**[]NetworkRouterTypeRuleGroupOptionTypesInner**](NetworkRouterTypeRuleGroupOptionTypesInner.md) |  | [optional] 
**NatOptionTypes** | Pointer to [**[]NetworkRouterTypeNatOptionTypesInner**](NetworkRouterTypeNatOptionTypesInner.md) |  | [optional] 
**BgpOptionTypes** | Pointer to [**[]NetworkRouterTypeBgpOptionTypesInner**](NetworkRouterTypeBgpOptionTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkRouterType{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


