# AddBudgetsRequestBudget

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Scope** | Pointer to **string** |  | [optional] [default to "account"]
**Period** | Pointer to **string** |  | [optional] [default to "year"]
**Year** | Pointer to **int64** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**Interval** | Pointer to **string** |  | [optional] [default to "year"]
**ScopeTenantId** | Pointer to **int64** | The Tenant ID to scope to, for use with &#x60;&#x60;scope&#x60;&#x60;&#x3D;tenant  | [optional] 
**ScopeGroupId** | Pointer to **int64** | The Tenant ID to scope to, for use with &#x60;&#x60;scope&#x60;&#x60;&#x3D;group   | [optional] 
**ScopeCloudId** | Pointer to **int64** | The Tenant ID to scope to, for use with &#x60;&#x60;scope&#x60;&#x60;&#x3D;cloud  | [optional] 
**ScopeUserId** | Pointer to **int64** | The Tenant ID to scope to, for use with &#x60;&#x60;scope&#x60;&#x60;&#x3D;user  | [optional] 
**Costs** | Pointer to **[]int64** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] [default to true]
**ForecastType** | Pointer to [**AddBudgetsRequestBudgetForecastType**](AddBudgetsRequestBudgetForecastType.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddBudgetsRequestBudget{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


