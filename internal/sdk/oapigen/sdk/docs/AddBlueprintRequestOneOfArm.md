# AddBlueprintRequestOneOfArm

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConfigType** | **string** | Configuration Type | 
**Json** | Pointer to **string** | ARM Template in JSON | [optional] 
**Yaml** | Pointer to **string** | ARM Template in YAML | [optional] 
**Git** | Pointer to [**AddBlueprintRequestOneOfArmGit**](AddBlueprintRequestOneOfArmGit.md) |  | [optional] 
**OsType** | Pointer to **string** | OS Type | [optional] 
**InstallAgent** | Pointer to [**AddBlueprintRequestOneOfArmInstallAgent**](AddBlueprintRequestOneOfArmInstallAgent.md) |  | [optional] 
**CloudInitEnabled** | Pointer to [**AddBlueprintRequestOneOfArmCloudInitEnabled**](AddBlueprintRequestOneOfArmCloudInitEnabled.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddBlueprintRequestOneOfArm{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


