# BlueprintCreateSuccessConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the blueprint | [optional] 
**Image** | Pointer to **string** | Path to display image. Defaults to an internal Morpheus image. | [optional] 
**Type** | Pointer to **string** | Blueprint Type | [optional] 
**Arm** | Pointer to [**BlueprintCreateSuccessConfigOneOfArm**](BlueprintCreateSuccessConfigOneOfArm.md) |  | [optional] 
**Visibility** | Pointer to **string** | Private or Public Access | [optional] [default to "private"]
**ResourcePermission** | Pointer to **map[string]interface{}** | Resource Permission Block | [optional] 
**Owner** | Pointer to **map[string]interface{}** | Owner | [optional] 
**Tenant** | Pointer to **map[string]interface{}** | Tenant | [optional] 
**CloudFormation** | Pointer to [**BlueprintCreateSuccessConfigOneOf1CloudFormation**](BlueprintCreateSuccessConfigOneOf1CloudFormation.md) |  | [optional] 
**Helm** | Pointer to [**BlueprintCreateSuccessConfigOneOf2Helm**](BlueprintCreateSuccessConfigOneOf2Helm.md) |  | [optional] 
**Kubernetes** | Pointer to [**BlueprintCreateSuccessConfigOneOf3Kubernetes**](BlueprintCreateSuccessConfigOneOf3Kubernetes.md) |  | [optional] 
**Config** | Pointer to [**BlueprintCreateSuccessConfigOneOf5Config**](BlueprintCreateSuccessConfigOneOf5Config.md) |  | [optional] 
**Terraform** | Pointer to [**BlueprintCreateSuccessConfigOneOf5Terraform**](BlueprintCreateSuccessConfigOneOf5Terraform.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BlueprintCreateSuccessConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


