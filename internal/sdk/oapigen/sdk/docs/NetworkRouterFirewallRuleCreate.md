# NetworkRouterFirewallRuleCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Firewall rule name | 
**Enabled** | Pointer to **bool** | Can be used to enable / disable the rule (true, false). Default is on | [optional] [default to true]
**Priority** | Pointer to **int64** | Priority for rule | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkRouterFirewallRuleCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


