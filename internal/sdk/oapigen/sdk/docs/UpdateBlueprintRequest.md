# UpdateBlueprintRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the blueprint | 
**Image** | Pointer to **string** | Path to display image. Defaults to an internal Morpheus image. | [optional] 
**Type** | **string** | Blueprint Type | 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Arm** | [**UpdateBlueprintRequestOneOfArm**](UpdateBlueprintRequestOneOfArm.md) |  | 
**CloudFormation** | [**UpdateBlueprintRequestOneOf1CloudFormation**](UpdateBlueprintRequestOneOf1CloudFormation.md) |  | 
**Helm** | [**UpdateBlueprintRequestOneOf2Helm**](UpdateBlueprintRequestOneOf2Helm.md) |  | 
**Kubernetes** | [**UpdateBlueprintRequestOneOf3Kubernetes**](UpdateBlueprintRequestOneOf3Kubernetes.md) |  | 
**Config** | Pointer to [**UpdateBlueprintRequestOneOf5Config**](UpdateBlueprintRequestOneOf5Config.md) |  | [optional] 
**Tiers** | **map[string]interface{}** | Tier definitions - Create in UI to view a baseline for object | 
**Terraform** | [**UpdateBlueprintRequestOneOf5Terraform**](UpdateBlueprintRequestOneOf5Terraform.md) |  | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateBlueprintRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


