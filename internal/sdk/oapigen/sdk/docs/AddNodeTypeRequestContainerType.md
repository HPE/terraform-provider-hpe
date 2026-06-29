# AddNodeTypeRequestContainerType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Node type name | 
**Labels** | Pointer to **[]string** |  | [optional] 
**ShortName** | **string** | The short name is a name with no spaces used for display in your container list. | 
**Description** | Pointer to **string** | Node type description | [optional] 
**ContainerVersion** | **string** | Version of the node type | 
**ProvisionTypeCode** | **string** | Provision type code, eg. &#x60;amazon&#x60;, etc. | 
**Scripts** | Pointer to **[]int64** | Array of script IDs. | [optional] 
**Templates** | Pointer to **[]int64** | Array of file template IDs. | [optional] 
**VirtualImageId** | Pointer to **int64** | Virtual image ID | [optional] 
**OsTypeId** | Pointer to **int64** | OsType ID | [optional] 
**StatTypeCode** | Pointer to **string** | Stat type code.  Varies with node type, see Provision Types (customOptionTypes) for allowed values within selected type. | [optional] 
**LogTypeCode** | Pointer to **string** | Log type code.  Varies with node type, see Provision Types (customOptionTypes) for allowed values within selected type. | [optional] 
**ServerType** | Pointer to **string** | Server type.  Always pass \&quot;vm\&quot;. | [optional] 
**ContainerPorts** | Pointer to [**[]AddNodeTypeRequestContainerTypeContainerPortsInner**](AddNodeTypeRequestContainerTypeContainerPortsInner.md) | List of exposed port definitions in the format NAME&#x3D;PORT|PROTOCOL | [optional] 
**EnvironmentVariables** | Pointer to [**[]AddNodeTypeRequestContainerTypeEnvironmentVariablesInner**](AddNodeTypeRequestContainerTypeEnvironmentVariablesInner.md) | The environmentVariables parameter is array of env objects. | [optional] 
**Config** | Pointer to **map[string]interface{}** | Config object varies with node type.  If using docker, scvmm, ARM, hyperv, or cloudformation, look up provision type details (customOptionTypes) for information. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddNodeTypeRequestContainerType{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


