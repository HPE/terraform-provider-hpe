# AddBlueprintRequestOneOf1CloudFormation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConfigType** | **string** | Configuration Type | 
**Json** | Pointer to **string** | CloudFormation Template in JSON | [optional] 
**Yaml** | Pointer to **string** | CloudFormation Template in YAML | [optional] 
**Git** | Pointer to [**AddBlueprintRequestOneOf1CloudFormationGit**](AddBlueprintRequestOneOf1CloudFormationGit.md) |  | [optional] 
**IAM** | Pointer to **bool** | CloudFormation Attribute CAPABILITY_IAM | [optional] [default to false]
**CAPABILITY_NAMED_IAM** | Pointer to **bool** | CloudFormation Attribute CAPABILITY_NAMED_IAM | [optional] [default to false]
**CAPABILITY_AUTO_EXPAND** | Pointer to **bool** | CloudFormation Attribute CAPABILITY_AUTO_EXPAND | [optional] [default to false]
**InstallAgent** | Pointer to **bool** | Install Morpheus Agent | [optional] [default to false]
**CloudInitEnabled** | Pointer to **bool** | Cloud Init Enabled | [optional] [default to false]

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddBlueprintRequestOneOf1CloudFormation{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


