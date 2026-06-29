# AddEnvironmentsRequestEnvironment

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A unique name for the environment | 
**Code** | **string** | A unique code for the environment | 
**Description** | Pointer to **string** | A description of the environment | [optional] 
**Visibility** | Pointer to **string** | private or public | [optional] [default to "private"]
**SortOrder** | Pointer to **int64** | Sort order | [optional] [default to 0]

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddEnvironmentsRequestEnvironment{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


