# AddSecurityGroupRulesRequestRule

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the rule | [optional] 
**Direction** | Pointer to **string** | Either &#x60;ingress&#x60; or &#x60;egress&#x60; | [optional] [default to "ingress"]
**SourceType** | Pointer to **string** | Either &#x60;cidr&#x60;, &#x60;group&#x60;, &#x60;tier&#x60;, &#x60;all&#x60; | [optional] [default to "cidr"]
**Source** | Pointer to **string** | CIDR representing the source IP(s) which should receive access. Required for &#x60;sourceType&#x60;&#x3D;cidr | [optional] 
**SourceGroup** | Pointer to [**AddSecurityGroupRulesRequestRuleSourceGroup**](AddSecurityGroupRulesRequestRuleSourceGroup.md) |  | [optional] 
**SourceTier** | Pointer to [**AddSecurityGroupRulesRequestRuleSourceTier**](AddSecurityGroupRulesRequestRuleSourceTier.md) |  | [optional] 
**PortRange** | Pointer to **string** | Either a single value (i.e. 55) or a port range (i.e. 1-65535) for which to open access to the source. Required if customRule is true, otherwise, ignored.  | [optional] 
**DestinationPortRange** | Pointer to **string** | Either a single value (i.e. 55) or a port range (i.e. 1-65535) for which to open access to the destination.  | [optional] 
**Protocol** | **string** | Either tcp, udp, icmp. Required if customRule is true, otherwise, ignored. | 
**DestinationType** | Pointer to **string** | Either cidr, group, tier, instance. | [optional] [default to "cidr"]
**Destination** | Pointer to **string** | CIDR representing the destination IP(s) which should receive access. Required for &#x60;destinationType&#x60;&#x3D;cidr. | [optional] 
**DestinationGroup** | Pointer to [**AddSecurityGroupRulesRequestRuleDestinationGroup**](AddSecurityGroupRulesRequestRuleDestinationGroup.md) |  | [optional] 
**DestinationTier** | Pointer to [**AddSecurityGroupRulesRequestRuleDestinationTier**](AddSecurityGroupRulesRequestRuleDestinationTier.md) |  | [optional] 
**RuleType** | **string** | Either &#x60;customRule&#x60; or an &#x60;instance type&#x60; code. | [default to "customRule"]
**Policy** | Pointer to **string** | Either &#x60;accept&#x60; or &#x60;deny&#x60;. | [optional] 
**InstanceTypeId** | Pointer to **int64** | The id of an Instance Type. If specified, the source CIDR will have access to all ports exposed by the particular instance in the cloud, app, or instance. Required if customRule is false, otherwise ignored.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddSecurityGroupRulesRequestRule{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


