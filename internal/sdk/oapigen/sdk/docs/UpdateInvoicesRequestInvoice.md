# UpdateInvoicesRequestInvoice

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Tags** | Pointer to **[]map[string]interface{}** | This adds or updates the specified Metadata tags and removes any tags not specified. Array of objects having a name and value.  | [optional] 
**AddTags** | Pointer to **[]map[string]interface{}** | Add or update value of Metadata tags. Array of objects having a name and value.  | [optional] 
**RemoveTags** | Pointer to **[]map[string]interface{}** | This removes the specified Metadata tags matching name and optionally value (if provided). Array of objects having a name and value.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateInvoicesRequestInvoice{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


