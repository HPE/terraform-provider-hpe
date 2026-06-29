# BlueprintKubernetesCreateSuccess

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the blueprint | [optional] 
**Image** | Pointer to **string** | Path to display image. Defaults to an internal Morpheus image. | [optional] 
**Type** | Pointer to **string** | Blueprint Type | [optional] 
**Kubernetes** | Pointer to [**BlueprintKubernetesCreateSuccessKubernetes**](BlueprintKubernetesCreateSuccessKubernetes.md) |  | [optional] 
**Config** | Pointer to [**BlueprintKubernetesCreateSuccessConfig**](BlueprintKubernetesCreateSuccessConfig.md) |  | [optional] 
**Visibility** | Pointer to **string** | Private or Public Access | [optional] [default to "private"]
**ResourcePermission** | Pointer to **map[string]interface{}** | Resource Permission Block | [optional] 
**Owner** | Pointer to **map[string]interface{}** | Owner | [optional] 
**Tenant** | Pointer to **map[string]interface{}** | Tenant | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BlueprintKubernetesCreateSuccess{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


