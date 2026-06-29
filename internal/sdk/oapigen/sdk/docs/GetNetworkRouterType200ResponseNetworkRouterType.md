# GetNetworkRouterType200ResponseNetworkRouterType

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
**OptionTypes** | Pointer to [**[]GetNetworkRouterType200ResponseNetworkRouterTypeOptionTypesInner**](GetNetworkRouterType200ResponseNetworkRouterTypeOptionTypesInner.md) |  | [optional] 
**RuleOptionTypes** | Pointer to [**[]GetNetworkRouterType200ResponseNetworkRouterTypeRuleOptionTypesInner**](GetNetworkRouterType200ResponseNetworkRouterTypeRuleOptionTypesInner.md) |  | [optional] 
**RuleGroupOptionTypes** | Pointer to [**[]GetNetworkRouterType200ResponseNetworkRouterTypeRuleGroupOptionTypesInner**](GetNetworkRouterType200ResponseNetworkRouterTypeRuleGroupOptionTypesInner.md) |  | [optional] 
**NatOptionTypes** | Pointer to [**[]GetNetworkRouterType200ResponseNetworkRouterTypeNatOptionTypesInner**](GetNetworkRouterType200ResponseNetworkRouterTypeNatOptionTypesInner.md) |  | [optional] 
**BgpOptionTypes** | Pointer to [**[]GetNetworkRouterType200ResponseNetworkRouterTypeBgpOptionTypesInner**](GetNetworkRouterType200ResponseNetworkRouterTypeBgpOptionTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRouterType200ResponseNetworkRouterType{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


