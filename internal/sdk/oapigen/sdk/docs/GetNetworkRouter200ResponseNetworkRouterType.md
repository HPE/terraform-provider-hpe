# GetNetworkRouter200ResponseNetworkRouterType

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
**HasNat** | Pointer to **bool** |  | [optional] 
**HasNetworkServer** | Pointer to **bool** |  | [optional] 
**HasFirewallGroups** | Pointer to **bool** |  | [optional] 
**HasSecurityGroupPriority** | Pointer to **bool** |  | [optional] 
**OptionTypes** | Pointer to [**[]GetNetworkRouter200ResponseNetworkRouterTypeOptionTypesInner**](GetNetworkRouter200ResponseNetworkRouterTypeOptionTypesInner.md) |  | [optional] 
**RuleOptionTypes** | Pointer to [**[]GetNetworkRouter200ResponseNetworkRouterTypeRuleOptionTypesInner**](GetNetworkRouter200ResponseNetworkRouterTypeRuleOptionTypesInner.md) |  | [optional] 
**FirewallGroupOptionTypes** | Pointer to [**[]GetNetworkRouter200ResponseNetworkRouterTypeFirewallGroupOptionTypesInner**](GetNetworkRouter200ResponseNetworkRouterTypeFirewallGroupOptionTypesInner.md) |  | [optional] 
**NatOptionTypes** | Pointer to [**[]GetNetworkRouter200ResponseNetworkRouterTypeNatOptionTypesInner**](GetNetworkRouter200ResponseNetworkRouterTypeNatOptionTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRouter200ResponseNetworkRouterType{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


