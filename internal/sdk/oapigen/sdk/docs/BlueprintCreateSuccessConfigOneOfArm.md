# BlueprintCreateSuccessConfigOneOfArm

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConfigType** | **string** | Configuration Type | 
**Json** | Pointer to **string** | ARM Template in JSON | [optional] 
**Yaml** | Pointer to **string** | ARM Template in YAML | [optional] 
**Git** | Pointer to [**BlueprintCreateSuccessConfigOneOfArmGit**](BlueprintCreateSuccessConfigOneOfArmGit.md) |  | [optional] 
**OsType** | Pointer to **string** | OS Type | [optional] 
**InstallAgent** | Pointer to [**BlueprintCreateSuccessConfigOneOfArmInstallAgent**](BlueprintCreateSuccessConfigOneOfArmInstallAgent.md) |  | [optional] 
**CloudInitEnabled** | Pointer to [**BlueprintCreateSuccessConfigOneOfArmCloudInitEnabled**](BlueprintCreateSuccessConfigOneOfArmCloudInitEnabled.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &BlueprintCreateSuccessConfigOneOfArm{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


