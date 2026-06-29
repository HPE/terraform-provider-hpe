# NetworkFirewallRuleGroupCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Network firewall rule group name | 
**Description** | Pointer to **string** | Network firewall rule group description | [optional] 
**Priority** | Pointer to **int64** | Network firewall rule group priority | [optional] 
**ExternalType** | **string** | Use SecurityPolicy | 
**GroupLayer** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkFirewallRuleGroupCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


