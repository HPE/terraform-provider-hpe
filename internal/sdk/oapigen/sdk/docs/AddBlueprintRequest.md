# AddBlueprintRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the blueprint | 
**Image** | Pointer to **string** | Path to display image. Defaults to an internal Morpheus image. | [optional] 
**Type** | **string** | Blueprint Type | 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Arm** | [**AddBlueprintRequestOneOfArm**](AddBlueprintRequestOneOfArm.md) |  | 
**CloudFormation** | [**AddBlueprintRequestOneOf1CloudFormation**](AddBlueprintRequestOneOf1CloudFormation.md) |  | 
**Helm** | [**AddBlueprintRequestOneOf2Helm**](AddBlueprintRequestOneOf2Helm.md) |  | 
**Kubernetes** | [**AddBlueprintRequestOneOf3Kubernetes**](AddBlueprintRequestOneOf3Kubernetes.md) |  | 
**Config** | Pointer to [**AddBlueprintRequestOneOf5Config**](AddBlueprintRequestOneOf5Config.md) |  | [optional] 
**Tiers** | **map[string]interface{}** | Tier definitions - Create in UI to view a baseline for object | 
**Terraform** | [**AddBlueprintRequestOneOf5Terraform**](AddBlueprintRequestOneOf5Terraform.md) |  | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddBlueprintRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


