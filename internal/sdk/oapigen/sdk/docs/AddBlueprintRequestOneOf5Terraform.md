# AddBlueprintRequestOneOf5Terraform

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConfigType** | **string** | Configuration Type | 
**Json** | Pointer to **string** | Terraform definition in JSON for &#x60;configType&#x60; &#x60;json&#x60; | [optional] 
**Tf** | Pointer to **string** | Terraform definition for &#x60;configType&#x60; &#x60;tf&#x60; | [optional] 
**Git** | Pointer to [**AddBlueprintRequestOneOf5TerraformGit**](AddBlueprintRequestOneOf5TerraformGit.md) |  | [optional] 
**TfvarSecret** | Pointer to **string** | tfvar secret from Morpheus Cypher | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddBlueprintRequestOneOf5Terraform{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


