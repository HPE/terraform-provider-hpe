# GroupUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A unique name scoped to your account for the group | 
**Code** | Pointer to **string** | Optional code for use with policies | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Location** | Pointer to **string** | Optional location argument for your group | [optional] 
**Config** | Pointer to [**GroupCreateConfig**](GroupCreateConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GroupUpdate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


